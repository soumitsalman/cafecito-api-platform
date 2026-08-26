// @title             Espresso API & MCP
// @version           0.5
// @description       Espresso is a business and market intelligence data API for AI agents, automated research, and analytical applications.
// @description       **Events** are concrete developments involving an organization, person, product, market, or region. **Signals** are higher-level conclusions synthesized from supporting Events.
// @description       **Choose a route by user intent**: What happened? Search Events. What does it mean or what is the outlook? Search Signals. What supports a conclusion? Retrieve a Signal, then list its supporting Events. What evidence or source coverage exists? Retrieve an Event, then inspect its evidence. Which exact filter value should I use? Use a discovery route only when the value is not already known.
// @description       **Recommended agent workflow**: (1) search the appropriate collection with the smallest useful filter set; (2) select IDs from `data`; (3) retrieve detail only for selected IDs; (4) traverse evidence, related Signals, or supporting Events only when explanation, provenance, or context is needed.
// @description       **Collections** return `{data, pagination, meta}`. Pagination contains `limit`, `num_results` (this page only), and `next_cursor`. To continue, send `pagination.next_cursor` unchanged as the next request `cursor`; never construct or decode cursor tokens. Empty collections return HTTP 200 with `data: []`. Detail routes return `{data}`; missing detail resources return HTTP 404. Errors use `{ "error": { "code", "message" } }`. Backend authentication uses the `X-API-KEY` header (or other headers listed in `API_KEY`); `/health` does not. Public clients send Bearer keys to the gateway, not this service.
// @description       **Filtering**: `tags` use fuzzy text matching. `event_types`, `categories`, `entities`, `impact_levels`, `companies`, `people`, `products`, and `regions` use exact matching after snake_case normalization. `categories` and `event_types` are separate fields. `from` and `to` bound record `created_at`, not occurrence, publication, lifecycle, or forecast time.
// @description       **Formats**: JSON is canonical. YAML and TOON represent the same public payload in token-optimized forms for MCP and AI-agent context. Public payloads never expose embeddings, relation direction, or internal storage objects.
// @schemes           https
// @securityDefinitions.apikey BackendAPIKey
// @in header
// @name X-API-KEY
// @license.name      MIT
// @contact.name      Project Cafecito
// @contact.url       http://cafecito.tech
// @contact.email     soumitsrah@cafecito.tech
package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/toon-format/toon-go"

	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/embedding"
	datautils "github.com/soumitsalman/data-utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	MIN_WINDOW              = 1
	DEFAULT_WINDOW          = 7 // DAYS
	DEFAULT_SCORE_THRESHOLD = 0.5
	DEFAULT_LIMIT           = 20
	MAX_LIMIT               = 100
)

const (
	API_ERROR_MSG_OUR_BAD              = "It's not you, it's us. Retry in a bit."
	API_ERROR_MSG_TRY_DIFFERENT_FORMAT = "Serialization issue. Fallback to `response_type=json`."
	API_ERROR_MSG_EVENT_NOT_FOUND      = "EVENT with this ID not found"
	API_ERROR_MSG_SIGNAL_NOT_FOUND     = "SIGNAL with this ID not found"
	API_ERROR_MSG_ACTION_NOT_FOUND     = "ACTION with this ID not found"
	API_ERROR_MSG_SOURCE_NOT_FOUND     = "SOURCE with this ID not found"
	API_ERROR_MSG_MISSING_API_KEY      = "Missing API Key"
)

// Configuration wires database, embedding, auth, and caching dependencies into HTTP handlers.
type Configuration struct {
	DB       *db.Cupboard
	Embedder embedding.Embedder
	APIKeys  map[string]string
}

// createPageRequest creates a PageRequest for db from given pagination params and validates the cursor.
func (p *paginationParams) createPageRequest(c *gin.Context, config *Configuration) (*db.PageRequest, error) {
	if p.Limit == 0 {
		p.Limit = DEFAULT_LIMIT
	}
	cursor, err := db.DecodeCursor(p.Cursor)
	if err != nil {
		return nil, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return &db.PageRequest{Limit: p.Limit, Cursor: cursor}, nil
}

func (p *sipQueryParams) bindFilters(c *gin.Context, config *Configuration, filters *db.Filters) error {
	filters.CreatedFrom = p.From
	filters.CreatedTo = p.To
	filters.Tags = utils.NormalizeTags(p.Tags)
	filters.Entities = utils.NormalizeTags(p.Entities)
	filters.Categories = utils.NormalizeTags(p.Categories)
	filters.Companies = utils.NormalizeTags(p.Companies)
	filters.People = utils.NormalizeTags(p.People)
	filters.Products = utils.NormalizeTags(p.Products)
	filters.Regions = utils.NormalizeTags(p.Regions)
	filters.ImpactedDomains = utils.NormalizeTags(p.ImpactedDomains)
	filters.ImpactLevels = utils.NormalizeTags(p.ImpactLevels)
	return nil
}

func (p *vectorSearchParams) bindFilters(c *gin.Context, config *Configuration, filters *db.Filters) error {
	if len(p.Q) > 0 {
		filters.Embedding = config.Embedder.EmbedQuery(c, p.Q)
		if len(filters.Embedding) == 0 {
			return utils.NewAPIError(utils.API_ERROR_EMBEDDING_ERROR, API_ERROR_MSG_OUR_BAD)
		}
		if p.ScoreThreshold > 0 {
			distance := (1 - p.ScoreThreshold) * 2
			filters.Distance = &distance
		}
	}
	return nil
}

// createEventFilters converts EventSearchParams into typed Filters.
func (p *EventSearchParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs:        p.IDs,
		SourceIDs:  p.SourceIDs,
		EventTypes: utils.NormalizeTags(p.EventTypes),
	}
	if err := p.sipQueryParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	if err := p.vectorSearchParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

// buildSignalFilters converts signalSearchParams into typed SipFilters.
func (p *SignalSearchParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs: p.IDs,
	}
	if err := p.sipQueryParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	if err := p.vectorSearchParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (params *EventEvidenceParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs:       params.IDs,
		SourceIDs: params.SourceIDs,
	}
	if err := params.sipQueryParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (params *EventSignalsParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs: params.IDs,
	}
	if err := params.sipQueryParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (params *SignalEventsParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs:        params.IDs,
		SourceIDs:  params.SourceIDs,
		EventTypes: utils.NormalizeTags(params.EventTypes),
	}
	if err := params.sipQueryParams.bindFilters(c, config, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func createSipKinds(resources []string) []string {
	var kinds []string
	for _, res := range resources {
		switch strings.ToLower(strings.TrimSpace(res)) {
		case "event", "evidence":
			kinds = append(kinds, db.SIP_KIND_EVENT)
		case "signal":
			kinds = append(kinds, db.SIP_KIND_SIGNAL)
		}
	}
	return kinds
}

// writeCollection writes a typed collection envelope as JSON, or compact text when requested.
func writePage[T any](c *gin.Context, items []T, limit int, next_cursor *string, response_type string) {
	if items == nil {
		items = []T{}
	}
	response := PageResponse[T]{
		Data:       items,
		Pagination: NewPagination(limit, len(items), next_cursor),
		Meta:       ResponseMeta{AsOf: time.Now().UTC()},
	}
	switch response_type {
	case "yaml":
		c.YAML(http.StatusOK, response)
	case "toon":
		if encoded, err := toon.MarshalString(response); err != nil {
			writeError(c, utils.NewAPIError(utils.API_ERROR_ENCODING_ERROR, API_ERROR_MSG_TRY_DIFFERENT_FORMAT))
		} else {
			c.String(http.StatusOK, encoded)
		}
	default:
		c.JSON(http.StatusOK, response)
	}
}

// writeDetail writes a typed detail envelope as JSON, or compact text when requested.
func writeItem[T any](c *gin.Context, item T, response_type string) {
	response := ItemResponse[T]{Data: item}
	switch response_type {
	case "yaml":
		c.YAML(http.StatusOK, response)
	case "toon":
		if encoded, err := toon.MarshalString(response); err != nil {
			writeError(c, utils.NewAPIError(utils.API_ERROR_ENCODING_ERROR, API_ERROR_MSG_TRY_DIFFERENT_FORMAT))
		} else {
			c.String(http.StatusOK, encoded)
		}
	default:
		c.JSON(http.StatusOK, response)
	}
}

// writeError writes an APIError to the response.
// Uses InternalServerError for DB, Embedding, Encoding errors and default cases
func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	body := APIError{Code: utils.API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD}
	if api_err, ok := err.(utils.APIError); ok {
		body = APIError{Code: api_err.Code, Message: api_err.Message}
		switch api_err.Code {
		case utils.API_ERROR_INVALID_REQUEST:
			status = http.StatusBadRequest
		case utils.API_ERROR_NOT_FOUND:
			status = http.StatusNotFound
		case utils.API_ERROR_UNAUTHORIZED:
			status = http.StatusUnauthorized
		}
	}
	c.AbortWithStatusJSON(status, ErrorResponse{Error: body})
}

// health godoc
// @Summary Check Espresso service availability
// @Description Use this lightweight endpoint to confirm that the Espresso service is reachable before making intelligence requests. It returns service status only; it does not validate an API key, search data, or report dependency health.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Service is alive"
// @ID healthCheck
// @Router /health [get]
func (r *Configuration) health(c *gin.Context) {
	// TODO: tie with embedder health and db health check
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// getEvents godoc
// @Summary Find concrete Events
// @Description Use when the user asks what happened to a company, person, product, region, or topic. Returns concrete Event records, not article bodies and not synthesized conclusions.
// @Description Carry a selected `data[].id` into Event detail, evidence, or related-Signals routes. Use `tags` for fuzzy concepts; use structured filters for exact normalized values. `categories` is not an alias for `event_types`.
// @Description Search Signals instead when the user asks for meaning, implication, or outlook. `from` and `to` bound record `created_at`, not occurrence or publication time.
// @Tags Events
// @Produce json
// @Param q query string false "Optional natural-language semantic query. Max 1024 characters." maxlength(1024)
// @Param score_threshold query number false "Minimum semantic similarity threshold for q. 0.0 is broad and 1.0 is strict. Default 0.5." default(0.5) minimum(0) maximum(1)
// @Param from query string false "Inclusive created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive created_at upper bound (YYYY-MM-DD)." format(date)
// @Param ids query []string false "Restrict to Event UUIDs (CSV)." collectionFormat(csv)
// @Param event_types query []string false "Exact Event type names in snake_case (CSV), for example policy_change,market_entry." collectionFormat(csv)
// @Param categories query []string false "Exact category names in snake_case (CSV), for example regulation,technology. This is separate from event_types." collectionFormat(csv)
// @Param entities query []string false "Exact entity names in snake_case, matched against company or people names (CSV), for example microsoft,nvidia." collectionFormat(csv)
// @Param impact_levels query []string false "Exact impact-level names (CSV), for example high,medium." collectionFormat(csv)
// @Param companies query []string false "Exact company names in snake_case (CSV), for example microsoft,nvidia." collectionFormat(csv)
// @Param people query []string false "Exact person names in snake_case (CSV), for example sam_altman,elon_musk." collectionFormat(csv)
// @Param products query []string false "Exact product names in snake_case (CSV), for example windows,geforce." collectionFormat(csv)
// @Param regions query []string false "Exact region names in snake_case (CSV), for example north_america,europe." collectionFormat(csv)
// @Param source_ids query []string false "Restrict to direct Event source UUIDs (CSV)." collectionFormat(csv)
// @Param tags query []string false "Fuzzy text match against persisted tag labels (CSV)." collectionFormat(csv)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} EventCollectionResponse "Event collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid query parameters, malformed UUID, or malformed cursor token"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID searchEvents
// @Router /events [get]
func (r *Configuration) getEvents(c *gin.Context) {
	var params EventSearchParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters, err := params.createFilters(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters.Kind = db.SIP_KIND_EVENT

	page_out, err := r.DB.QuerySips(c.Request.Context(), *filters, *page_req)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEvent godoc
// @Summary Inspect one Event
// @Description Use after selecting an Event from a collection. Returns its complete public view, optional Source provenance, and links/counts for available evidence and related Signals.
// @Description Carry the Event ID to `/events/{event_id}/evidence` for supporting context or source coverage, or to `/events/{event_id}/signals` for associated higher-level conclusions. `created_at` is record creation time, not Event occurrence time.
// @Tags Events
// @Produce json
// @Param event_id path string true "Event UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} EventDetailResponse "Event detail envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID"
// @Failure 404 {object} ErrorResponse "No Event with this UUID"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID getEvent
// @Router /events/{event_id} [get]
func (r *Configuration) getEvent(c *gin.Context) {
	var params itemParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	event, err := r.DB.GetSip(c.Request.Context(), params.ID, db.SIP_KIND_EVENT)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if event.IsZero() {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_EVENT_NOT_FOUND))
		return
	}
	counts, err := r.DB.CountRelations(c.Request.Context(), params.ID)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeItem(c, NewDigestDocumentForExtendedSip(&event).addEventDetails(counts), params.ResponseType)
}

// getEventEvidence godoc
// @Summary Inspect evidence for an Event
// @Description Use after selecting an Event when the user needs supporting context, available source coverage, or traceability. Returns directly related evidence records with identity, creation time, tags, Source IDs, and available URLs.
// @Description This is not an article-body endpoint, a story-clustering endpoint, or a complete record-history export. An empty collection means no evidence records are available for this Event under the supplied filters.
// @Tags Events
// @Produce json
// @Param event_id path string true "Event UUID (RFC 4122)." format(uuid)
// @Param ids query []string false "Restrict returned evidence records to Event UUIDs (CSV)." collectionFormat(csv)
// @Param source_ids query []string false "Restrict direct evidence records to Source UUIDs (CSV)." collectionFormat(csv)
// @Param from query string false "Inclusive evidence created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive evidence created_at upper bound (YYYY-MM-DD)." format(date)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} EventEvidenceCollectionResponse "Event evidence collection envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID or invalid parameters"
// @Failure 404 {object} ErrorResponse "No Event with this UUID"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID getEventEvidence
// @Router /events/{event_id}/evidence [get]
func (r *Configuration) getEventEvidence(c *gin.Context) {
	var params EventEvidenceParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters, err := params.createFilters(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters.Kind = db.SIP_KIND_EVENT

	exists, err := r.DB.SipExists(c.Request.Context(), params.ID, db.SIP_KIND_EVENT)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if !exists {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_EVENT_NOT_FOUND))
		return
	}
	page_out, err := r.DB.QuerySameSips(c.Request.Context(), params.ID, *filters, *page_req)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	evidence := datautils.Transform(page_out.Items, func(sip *db.Sip) EventEvidence {
		return NewEventEvidence(sip)
	})
	writePage(c, evidence, page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEventSignals godoc
// @Summary Find Signals connected to an Event
// @Description Use after selecting an Event to find higher-level conclusions associated with that development. Returns Signal records that can be inspected individually or followed to their supporting Events.
// @Description An empty collection means Espresso has no available Signal connected to this Event; it does not invalidate the Event. This route narrows associated Signals; use `/signals` for a new Signal search.
// @Tags Events
// @Produce json
// @Param event_id path string true "Event UUID (RFC 4122)." format(uuid)
// @Param ids query []string false "Restrict returned Signals to Signal UUIDs (CSV)." collectionFormat(csv)
// @Param from query string false "Inclusive Signal created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive Signal created_at upper bound (YYYY-MM-DD)." format(date)
// @Param impact_levels query []string false "Exact impact-level names (CSV), for example high,medium." collectionFormat(csv)
// @Param impacted_domains query []string false "Exact impacted-domain names in snake_case (CSV), for example public_health,climate." collectionFormat(csv)
// @Param tags query []string false "Fuzzy text match against persisted Signal tag labels (CSV)." collectionFormat(csv)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} SignalCollectionResponse "Signals derived from this Event"
// @Failure 400 {object} ErrorResponse "Malformed UUID, invalid cursor token, or invalid parameters"
// @Failure 404 {object} ErrorResponse "No Event with this UUID"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID getEventSignals
// @Router /events/{event_id}/signals [get]
func (r *Configuration) getEventSignals(c *gin.Context) {
	var params EventSignalsParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters, err := params.createFilters(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters.Kind = db.SIP_KIND_SIGNAL

	exists, err := r.DB.SipExists(c.Request.Context(), params.ID, db.SIP_KIND_EVENT)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if !exists {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_EVENT_NOT_FOUND))
		return
	}
	page_out, err := r.DB.QueryDerivedSips(c.Request.Context(), params.ID, *filters, *page_req)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getSignals godoc
// @Summary Find synthesized Signals
// @Description Use when the user asks what a set of developments means, what impact is expected, or what broader conclusion Espresso has produced. Returns synthesized Signals, not raw observations or article content.
// @Description Carry a selected `data[].id` into Signal detail, then list supporting Events when the conclusion needs substantiation. Search Events instead when the user needs a concrete development rather than an interpretation.
// @Description `from` and `to` bound Signal `created_at`; tags are fuzzy text matches, while impact levels and impacted domains are exact normalized values.
// @Tags Signals
// @Produce json
// @Param q query string false "Optional natural-language semantic query. Max 1024 characters." maxlength(1024)
// @Param score_threshold query number false "Minimum semantic similarity threshold for q. 0.0 is broad and 1.0 is strict. Default 0.5." default(0.5) minimum(0) maximum(1)
// @Param from query string false "Inclusive Signal created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive Signal created_at upper bound (YYYY-MM-DD)." format(date)
// @Param ids query []string false "Restrict to Signal UUIDs (CSV)." collectionFormat(csv)
// @Param impact_levels query []string false "Exact impact-level names (CSV), for example high,medium." collectionFormat(csv)
// @Param impacted_domains query []string false "Exact impacted-domain names in snake_case (CSV), for example public_health,climate." collectionFormat(csv)
// @Param tags query []string false "Fuzzy text match against persisted Signal tag labels (CSV)." collectionFormat(csv)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} SignalCollectionResponse "Signal collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid query parameters, malformed UUID, or malformed cursor token"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID searchSignals
// @Router /signals [get]
func (r *Configuration) getSignals(c *gin.Context) {
	var params SignalSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters, err := params.createFilters(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters.Kind = db.SIP_KIND_SIGNAL

	page_out, err := r.DB.QuerySips(c.Request.Context(), *filters, *page_req)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// @Summary Inspect one Signal
// @Description Use after selecting a Signal from a collection. Returns its complete public fields, optional Source provenance, and a link/count for Events that support the conclusion.
// @Description Carry the Signal ID to `/signals/{signal_id}/events` when an answer needs concrete supporting developments. `created_at` is record creation time.
// @Tags Signals
// @Produce json
// @Param signal_id path string true "Signal UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SignalDetailResponse "Signal detail envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID"
// @Failure 404 {object} ErrorResponse "No Signal with this UUID"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID getSignal
// @Router /signals/{signal_id} [get]
func (r *Configuration) getSignal(c *gin.Context) {
	var params itemParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	signal, err := r.DB.GetSip(c.Request.Context(), params.ID, db.SIP_KIND_SIGNAL)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if signal.IsZero() {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_SIGNAL_NOT_FOUND))
		return
	}
	counts, err := r.DB.CountRelations(c.Request.Context(), params.ID)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeItem(c, NewDigestDocumentForExtendedSip(&signal).addSignalDetails(counts), params.ResponseType)
}

// getSignalEvents godoc
// @Summary Inspect Events supporting a Signal
// @Description Use after selecting a Signal when an agent must explain, verify, or cite the concrete developments behind its conclusion. Returns Events that support the Signal.
// @Description Apply Event filters only to narrow this existing support set. This route does not perform a new semantic search and does not return unrelated Events. An empty collection means no supporting Events match the supplied filters.
// @Tags Signals
// @Produce json
// @Param signal_id path string true "Signal UUID (RFC 4122)." format(uuid)
// @Param ids query []string false "Restrict returned Events to Event UUIDs (CSV)." collectionFormat(csv)
// @Param event_types query []string false "Exact Event type names in snake_case (CSV), for example policy_change,market_entry." collectionFormat(csv)
// @Param categories query []string false "Exact category names in snake_case (CSV), for example regulation,technology. This is separate from event_types." collectionFormat(csv)
// @Param entities query []string false "Exact entity names in snake_case, matched against company or people names (CSV), for example microsoft,nvidia." collectionFormat(csv)
// @Param impact_levels query []string false "Exact impact-level names (CSV), for example high,medium." collectionFormat(csv)
// @Param companies query []string false "Exact company names in snake_case (CSV), for example microsoft,nvidia." collectionFormat(csv)
// @Param people query []string false "Exact person names in snake_case (CSV), for example sam_altman,elon_musk." collectionFormat(csv)
// @Param products query []string false "Exact product names in snake_case (CSV), for example windows,geforce." collectionFormat(csv)
// @Param regions query []string false "Exact region names in snake_case (CSV), for example north_america,europe." collectionFormat(csv)
// @Param source_ids query []string false "Restrict returned Events to Source UUIDs (CSV)." collectionFormat(csv)
// @Param tags query []string false "Fuzzy text match against persisted Event tag labels (CSV)." collectionFormat(csv)
// @Param from query string false "Inclusive Event created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive Event created_at upper bound (YYYY-MM-DD)." format(date)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} EventCollectionResponse "Events supporting this Signal"
// @Failure 400 {object} ErrorResponse "Malformed UUID, invalid cursor token, or invalid parameters"
// @Failure 404 {object} ErrorResponse "No Signal with this UUID"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID getSignalEvents
// @Router /signals/{signal_id}/events [get]
func (r *Configuration) getSignalEvents(c *gin.Context) {
	var params SignalEventsParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters, err := params.createFilters(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	filters.Kind = db.SIP_KIND_EVENT

	exists, err := r.DB.SipExists(c.Request.Context(), params.ID, db.SIP_KIND_SIGNAL)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if !exists {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_SIGNAL_NOT_FOUND))
		return
	}
	page_out, err := r.DB.QuerySupportingSips(c.Request.Context(), params.ID, *filters, *page_req)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// @Summary Find intelligence Sources
// @Description Use to discover or resolve provenance Sources before filtering Events by `source_ids`, or when an answer needs source metadata for citation. `q` performs case-insensitive metadata matching across source domain, name, and URL; it is not semantic search.
// @Description Carry a selected `data[].id` into Source detail or `GET /events?source_ids={source_id}`. Source results describe publishers and provenance; they do not contain Events.
// @Tags Sources
// @Produce json
// @Param q query string false "Case-insensitive Source metadata search." maxlength(1024)
// @Param domains query []string false "Exact Source domain-name filter (CSV)." collectionFormat(csv)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SourceCollectionResponse "Source collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, cursor token, or parameters"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID listIntelligenceSources
// @Router /sources [get]
func (r *Configuration) getSources(c *gin.Context) {
	var params SourcesParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QuerySources(c.Request.Context(), params.Q, params.Domains, *page_req)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, NewSourceDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getSource godoc
// @Summary Inspect one Source
// @Description Returns provenance metadata for one Source. Use this route to enrich a citation or inspect publisher metadata, not to retrieve published Events.
// @Description To find Events from this Source, call `GET /events?source_ids={source_id}`. Optional description, favicon, and RSS feed fields can be absent when unavailable.
// @Tags Sources
// @Produce json
// @Param source_id path string true "Source UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SourceItemResponse "Source detail envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID"
// @Failure 404 {object} ErrorResponse "No Source with this UUID"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID getIntelligenceSource
// @Router /sources/{source_id} [get]
func (r *Configuration) getSource(c *gin.Context) {
	var item itemParams
	if err := item.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	source, err := r.DB.GetSource(c.Request.Context(), item.ID)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if source.IsZero() {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_SOURCE_NOT_FOUND))
		return
	}
	writeItem(c, NewSourceDocument(&source), item.ResponseType)
}

// getTags godoc
// @Summary Discover fuzzy tag vocabulary
// @Description Use only when an agent needs vocabulary for a fuzzy `tags` query. Returns persisted labels for Event and Signal filtering; tags are not a fixed taxonomy or exact-only values.
// @Description If a useful tag is already known, search Events or Signals directly instead of making a discovery request. Use `resource` to limit discovery to Event or Signal labels.
// @Tags Tags
// @Produce json
// @Param q query string false "Case-insensitive substring or prefix match." maxlength(1024)
// @Param resource query []string false "Optional resource scope (CSV): event, signal." collectionFormat(csv)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} DiscoveryValueCollectionResponse "Tag value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, cursor token, or response_type"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID listIntelligenceTags
// @Router /tags [get]
func (r *Configuration) getTags(c *gin.Context) {
	var params TagsParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	page, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}

	kinds := createSipKinds(params.Resource)
	if len(kinds) == 0 {
		kinds = []string{db.SIP_KIND_EVENT, db.SIP_KIND_SIGNAL}
	}
	page_out, err := r.DB.QueryTags(c.Request.Context(), params.Q, kinds, *page)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	next := encodeNextCursor(page_out.NextCursor)
	items := datautils.Transform(page_out.Items, func(tag *string) db.Tag {
		return db.Tag{Value: *tag}
	})
	writePage(c, items, page.Limit, next, params.ResponseType)
}

// getEntities godoc
// @Summary Discover exact Event entity filters
// @Description Use only when an agent needs available company or people names before applying an exact Event filter. Returned values are normalized snake_case filter strings, not canonical entity IDs or profiles.
// @Description If a known normalized value is already available, query Events directly. Use `types=company` or `types=people` to reduce the returned vocabulary.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter." maxlength(1024)
// @Param types query []string false "Entity types (CSV): company, people." collectionFormat(csv)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} DiscoveryValueCollectionResponse "Entity value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, cursor token, or parameters"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID listIntelligenceEntities
// @Router /entities [get]
func (r *Configuration) getEntities(c *gin.Context) {
	var params EntitiesParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	page, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	types := params.Types
	if len(types) == 0 {
		types = []string{db.EVENT_TAG_TYPE_COMPANY, db.EVENT_TAG_TYPE_PEOPLE}
	}
	page_out, err := r.DB.QueryEventTags(c.Request.Context(), params.Q, types, *page)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, page_out.Items, page.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getRegions godoc
// @Summary Discover exact Event region filters
// @Description Use only when an agent needs available region values before applying the exact `regions` Event filter. Returned values are normalized snake_case filter strings, not canonical places, coordinates, or structured geography.
// @Description If a known normalized value is already available, query Events directly.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter." maxlength(1024)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} DiscoveryValueCollectionResponse "Region value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, cursor token, or parameters"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID listIntelligenceRegions
// @Router /regions [get]
func (r *Configuration) getRegions(c *gin.Context) {
	var params DiscoveryParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	page, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QueryEventTags(c.Request.Context(), params.Q, []string{db.EVENT_TAG_TYPE_REGION}, *page)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, page_out.Items, page.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEventTypes godoc
// @Summary Discover exact Event type filters
// @Description Use only when an agent needs available values before applying the exact `event_types` Event filter. Returned values are normalized snake_case filter strings.
// @Description `event_types` and `categories` filter different Event fields. If a known normalized value is already available, query Events directly.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter." maxlength(1024)
// @Param response_type query string false "Response serialization: JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token. Send pagination.next_cursor from a previous response unchanged as cursor; never construct or decode it."
// @Success 200 {object} DiscoveryValueCollectionResponse "Event type value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, cursor token, or parameters"
// @Failure 500 {object} ErrorResponse "Service unavailable; retry."
// @Security BackendAPIKey
// @ID listIntelligenceEventTypes
// @Router /event-types [get]
func (r *Configuration) getEventTypes(c *gin.Context) {
	var params DiscoveryParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	page, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QueryEventTags(c.Request.Context(), params.Q, []string{db.EVENT_TAG_TYPE_EVENT_TYPE}, *page)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writePage(c, page_out.Items, page.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

func NewRouter(db *db.Cupboard, embedder embedding.Embedder, api_keys map[string]string) *gin.Engine {
	config := &Configuration{
		DB:       db,
		Embedder: embedder,
		APIKeys:  api_keys,
		// cache: otter.Must(&otter.Options[string, []float32]{
		// 	MaximumSize:      _CACHE_SIZE,
		// 	ExpiryCalculator: otter.ExpiryAccessing[string, []float32](_CACHE_TTL),
		// }),
	}

	router := gin.New()
	router.Use(
		requestLogger,
		gin.Recovery(),
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: false,
			MaxAge:           24 * time.Hour,
		}),
	)

	router.GET("/health", config.health)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authenticated group
	protected := router.Group("/")
	protected.Use(config.apiKeyMiddleware)

	// TAGS discovery routes
	protected.GET("/tags", config.getTags)
	protected.GET("/entities", config.getEntities)
	protected.GET("/regions", config.getRegions)
	protected.GET("/event-types", config.getEventTypes)

	// SOURCES routes
	protected.GET("/sources", config.getSources)
	protected.GET("/sources/:id", config.getSource)

	// EVENTS routes
	protected.GET("/events", config.getEvents)
	protected.GET("/events/:id", config.getEvent)
	protected.GET("/events/:id/signals", config.getEventSignals)
	protected.GET("/events/:id/evidence", config.getEventEvidence)

	// SIGNALS routes
	protected.GET("/signals", config.getSignals)
	protected.GET("/signals/:id", config.getSignal)
	protected.GET("/signals/:id/events", config.getSignalEvents)

	return router
}

// Middleware
func (r *Configuration) apiKeyMiddleware(c *gin.Context) {
	if len(r.APIKeys) == 0 {
		c.Next()
		return
	}
	for header, expected := range r.APIKeys {
		if strings.TrimSpace(c.GetHeader(header)) == expected {
			c.Next()
			return
		}
	}
	writeError(c, utils.NewAPIError(utils.API_ERROR_UNAUTHORIZED, API_ERROR_MSG_MISSING_API_KEY))
}

// requestLogger logs request path, query parameters, status and latency in JSON via zerolog
func requestLogger(c *gin.Context) {
	start := time.Now()
	c.Next()

	status := c.Writer.Status()

	var evt *zerolog.Event
	if len(c.Errors) > 0 || status >= 500 {
		evt = log.Error()
	} else if status >= 400 {
		evt = log.Warn()
	} else {
		evt = log.Info()
	}
	evt.Str("module", "ROUTER").Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Interface("query", c.Request.URL.Query()).
		Int("status", status).
		Float64("latency", time.Since(start).Seconds())

	if len(c.Errors) > 0 {
		evt.Str("error", c.Errors.String())
	}
	evt.Msg("incoming")
}

// _WIRING_END_
