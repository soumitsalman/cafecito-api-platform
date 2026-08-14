// @title 			Espresso API & MCP
// @version 		0.2
// @description 	MCP-ready business intelligence over curated intelligence records for agents, dashboards, and automated research workflows.
// @description  An **Event** is a concrete intelligence record with kind `event`. A **Signal** is a synthesized conclusion derived from Events. The internal storage word `sip` is not part of the public vocabulary.
// @description 	Agent workflow: (1) listTags to discover filter vocabulary; (2) searchEvents for developments; (3) getEvent and getEventEvidence to trace source coverage; (4) searchSignals and getSignalEvents to trace synthesized conclusions.
// @description  Conventions: Auth is optional at the backend but API-key protected through the gateway. Collections use page pagination: `limit` default 20 max 100, and an opaque `page` returned as `next_page`. Empty collections return HTTP 200 with `data: []`. Missing detail resources return 404. All IDs are RFC 4122 UUID strings.
// @description Response formats: use `response_type=json` for canonical structured application data. Use `response_type=yaml` or `response_type=toon` for token-optimized output to MCP and AI-agent clients. Public Event and Signal payloads expose flattened intelligence fields without embeddings, relation direction, or a nested internal object.
// @schemes 		https
// @license.name 	MIT
// @contact.name 	Project Cafecito
// @contact.url  	http://cafecito.tech
// @contact.email 	soumitsrah@cafecito.tech
package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/stoewer/go-strcase"
	"github.com/toon-format/toon-go"

	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/embedding"
	datautils "github.com/soumitsalman/data-utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	MIN_WINDOW       = 1
	DEFAULT_WINDOW   = 7 // DAYS
	DEFAULT_ACCURACY = 0.5
	DEFAULT_LIMIT    = 20
	MAX_LIMIT        = 100
)

const (
	API_ERROR_INVALID_REQUEST = "invalid_request"
	API_ERROR_DB_ERROR        = "db_unavailable"
	API_ERROR_EMBEDDING_ERROR = "embedder_unavailable"
	API_ERROR_ENCODING_ERROR  = "encoding_error"
	API_ERROR_INVALID_DATA    = "invalid_data"
	API_ERROR_NOT_FOUND       = "not_found"
	API_ERROR_UNAUTHORIZED    = "unauthorized"
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
		return nil, APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return &db.PageRequest{Limit: p.Limit, Cursor: cursor}, nil
}

func normalizeTags(items []string) []string {
	for i, item := range items {
		items[i] = strcase.SnakeCase(item)
	}
	return items
}

func (p *sipQueryParams) bindFilters(c *gin.Context, config *Configuration, filters *db.Filters) error {
	filters.CreatedFrom = p.From
	filters.CreatedTo = p.To
	filters.Tags = normalizeTags(p.Tags)
	filters.Entities = normalizeTags(p.Entities)
	filters.Categories = normalizeTags(p.Categories)
	filters.Companies = normalizeTags(p.Companies)
	filters.People = normalizeTags(p.People)
	filters.Products = normalizeTags(p.Products)
	filters.Regions = normalizeTags(p.Regions)
	filters.ImpactedDomains = normalizeTags(p.ImpactedDomains)
	filters.ImpactLevels = normalizeTags(p.ImpactLevels)
	return nil
}

func (p *vectorSearchParams) bindFilters(c *gin.Context, config *Configuration, filters *db.Filters) error {
	if len(p.Q) > 0 {
		filters.Embedding = config.Embedder.EmbedQuery(c, p.Q)
		if len(filters.Embedding) == 0 {
			return APIError{Code: API_ERROR_EMBEDDING_ERROR, Message: API_ERROR_MSG_OUR_BAD}
		}
		distance := (1 - p.Acc) * 2
		filters.Distance = &distance
	}
	return nil
}

// createEventFilters converts EventSearchParams into typed Filters.
func (p *EventSearchParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs:        p.IDs,
		SourceIDs:  p.SourceIDs,
		EventTypes: normalizeTags(p.EventTypes),
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
		EventTypes: normalizeTags(params.EventTypes),
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
func writePage[T any](c *gin.Context, items []T, limit int, cursor string, next_cursor *string, response_type string) {
	if items == nil {
		items = []T{}
	}
	var cursor_val *string = nil
	if cursor != "" {
		cursor_val = &cursor
	}
	response := PageResponse[T]{
		Data:       items,
		Pagination: Pagination{Limit: limit, NumResults: len(items), Cursor: cursor_val, NextCursor: next_cursor},
		Meta:       ResponseMeta{AsOf: time.Now().UTC()},
	}
	switch response_type {
	case "yaml":
		c.YAML(http.StatusOK, response)
	case "toon":
		if encoded, err := toon.MarshalString(response); err != nil {
			writeError(c, APIError{Code: API_ERROR_ENCODING_ERROR, Message: API_ERROR_MSG_TRY_DIFFERENT_FORMAT})
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
			writeError(c, APIError{Code: API_ERROR_ENCODING_ERROR, Message: API_ERROR_MSG_TRY_DIFFERENT_FORMAT})
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
	if api_err, ok := err.(APIError); ok {
		switch api_err.Code {
		case API_ERROR_INVALID_REQUEST:
			status = http.StatusBadRequest
		case API_ERROR_NOT_FOUND:
			status = http.StatusNotFound
		}
		c.AbortWithStatusJSON(status, ErrorResponse{Error: api_err})
	} else {
		c.AbortWithStatusJSON(status, ErrorResponse{Error: APIError{Message: API_ERROR_MSG_OUR_BAD}})
	}
}

// health godoc
// @Summary Check API health
// @Description Lightweight liveness probe. Use it before other tools to confirm the Espresso backend is reachable. This endpoint does not require query parameters and returns only service status.
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
// @Summary Search Events
// @Description Find Event records that match optional filters and an optional semantic query. Use this route to find concrete developments before following an Event's detail, evidence, or Signal links.
// @Description **Time**: `from` and `to` are inclusive ISO date-only bounds on record `created_at`; they are not occurrence or publication timestamps.
// @Description **Filters**: tags use fuzzy text matching. `event_types`, `categories`, `entities`, `impact_levels`, `companies`, `people`, `products`, and `regions` use exact matching after normalizing names to snake_case. `categories` is a separate category filter from `event_types`.
// @Description **Output**: Event collections contain flattened intelligence fields plus `id`, `created_at`, and `kind`; provenance fields are available on the detail route when usable. YAML and TOON are token-optimized serializations for MCP and AI-agent clients.
// @Tags Events
// @Produce json
// @Param q query string false "Optional natural-language semantic query. Max 1024 characters." maxlength(1024)
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
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} EventCollectionResponse "Event collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid query parameters, malformed UUID, or malformed page token"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database or embedder unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEvent godoc
// @Summary Retrieve one Event
// @Description Retrieve an Event by UUID. The detail payload contains flattened intelligence fields and detail-only provenance, optional Source metadata, relation links, and relation counts. `created_at` is the record creation time, not an Event occurrence date.
// @Tags Events
// @Produce json
// @Param event_id path string true "Event UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} EventDetailResponse "Event detail envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID"
// @Failure 404 {object} ErrorResponse "No Event with this UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	if event.IsZero() {
		writeError(c, APIError{Code: API_ERROR_NOT_FOUND, Message: API_ERROR_MSG_EVENT_NOT_FOUND})
		return
	}
	counts, err := r.DB.CountRelations(c.Request.Context(), params.ID)
	if err != nil {
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writeItem(c, NewDigestDocumentForExtendedSip(&event).addEventDetails(counts), params.ResponseType)
}

// getEventEvidence godoc
// @Summary Retrieve an Event evidence trail
// @Description Show the directly related records that make up an Event's evidence trail, helping clients assess source coverage. It is not article content or story membership. The item projection exposes record identity, creation time, tags, source ID, and available URLs.
// @Tags Events
// @Produce json
// @Param event_id path string true "Event UUID (RFC 4122)." format(uuid)
// @Param ids query []string false "Restrict returned evidence records to Event UUIDs (CSV)." collectionFormat(csv)
// @Param source_ids query []string false "Restrict direct evidence records to Source UUIDs (CSV)." collectionFormat(csv)
// @Param from query string false "Inclusive evidence created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive evidence created_at upper bound (YYYY-MM-DD)." format(date)
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} EventEvidenceCollectionResponse "Event evidence collection envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID or invalid parameters"
// @Failure 404 {object} ErrorResponse "No Event with this UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	if !exists {
		writeError(c, APIError{Code: API_ERROR_NOT_FOUND, Message: API_ERROR_MSG_EVENT_NOT_FOUND})
		return
	}
	page_out, err := r.DB.QuerySameSips(c.Request.Context(), params.ID, *filters, *page_req)
	if err != nil {
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	evidence := datautils.Transform(page_out.Items, func(sip *db.Sip) EventEvidence {
		return NewEventEvidence(sip)
	})
	writePage(c, evidence, page_req.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEventSignals godoc
// @Summary Retrieve Signals derived from an Event
// @Description Return the higher-level Signals that were derived from the Event's evidence trail. Use this relation route after retrieving an Event to inspect synthesized conclusions; an existing Event with no matching Signals returns an empty collection.
// @Tags Events
// @Produce json
// @Param event_id path string true "Event UUID (RFC 4122)." format(uuid)
// @Param ids query []string false "Restrict returned Signals to Signal UUIDs (CSV)." collectionFormat(csv)
// @Param from query string false "Inclusive Signal created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive Signal created_at upper bound (YYYY-MM-DD)." format(date)
// @Param impact_levels query []string false "Exact impact-level names (CSV), for example high,medium." collectionFormat(csv)
// @Param impacted_domains query []string false "Exact impacted-domain names in snake_case (CSV), for example public_health,climate." collectionFormat(csv)
// @Param tags query []string false "Fuzzy text match against persisted Signal tag labels (CSV)." collectionFormat(csv)
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} SignalCollectionResponse "Signals derived from this Event"
// @Failure 400 {object} ErrorResponse "Malformed UUID, invalid page token, or invalid parameters"
// @Failure 404 {object} ErrorResponse "No Event with this UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database or embedder unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	if !exists {
		writeError(c, APIError{Code: API_ERROR_NOT_FOUND, Message: API_ERROR_MSG_EVENT_NOT_FOUND})
		return
	}
	page_out, err := r.DB.QueryDerivedSips(c.Request.Context(), params.ID, *filters, *page_req)
	if err != nil {
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getSignals godoc
// @Summary Search Signals
// @Description Find synthesized Signal records with optional filters and an optional semantic query. Use this route for higher-level conclusions, then follow a Signal detail or supporting-Events link to inspect its basis.
// @Description **Time**: `from` and `to` are inclusive ISO date-only bounds on Signal `created_at`.
// @Description **Filters**: tags use fuzzy text matching. Impact levels and impacted domains use exact snake_case matching. Collection items omit provenance fields; use Signal detail when Source metadata is usable.
// @Tags Signals
// @Produce json
// @Param q query string false "Optional natural-language semantic query. Max 1024 characters." maxlength(1024)
// @Param from query string false "Inclusive Signal created_at lower bound (YYYY-MM-DD)." format(date)
// @Param to query string false "Inclusive Signal created_at upper bound (YYYY-MM-DD)." format(date)
// @Param ids query []string false "Restrict to Signal UUIDs (CSV)." collectionFormat(csv)
// @Param impact_levels query []string false "Exact impact-level names (CSV), for example high,medium." collectionFormat(csv)
// @Param impacted_domains query []string false "Exact impacted-domain names in snake_case (CSV), for example public_health,climate." collectionFormat(csv)
// @Param tags query []string false "Fuzzy text match against persisted Signal tag labels (CSV)." collectionFormat(csv)
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} SignalCollectionResponse "Signal collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid query parameters, malformed UUID, or malformed page token"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database or embedder unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// @Summary Retrieve one Signal
// @Description Retrieve a Signal by UUID. The detail payload contains flattened intelligence fields, optional provenance when a usable Source exists, and links/counts for supporting Events. `created_at` is the Signal record creation time.
// @Tags Signals
// @Produce json
// @Param signal_id path string true "Signal UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SignalDetailResponse "Signal detail envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID"
// @Failure 404 {object} ErrorResponse "No Signal with this UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	if signal.IsZero() {
		writeError(c, APIError{Code: API_ERROR_NOT_FOUND, Message: API_ERROR_MSG_SIGNAL_NOT_FOUND})
		return
	}
	counts, err := r.DB.CountRelations(c.Request.Context(), params.ID)
	if err != nil {
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writeItem(c, NewDigestDocumentForExtendedSip(&signal).addSignalDetails(counts), params.ResponseType)
}

// getSignalEvents godoc
// @Summary Retrieve Events supporting a Signal
// @Description Return the Events that were used to derive the Signal, so clients can inspect the basis of its conclusion. An existing Signal with no matching Events returns an empty collection.
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
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} EventCollectionResponse "Events supporting this Signal"
// @Failure 400 {object} ErrorResponse "Malformed UUID, invalid page token, or invalid parameters"
// @Failure 404 {object} ErrorResponse "No Signal with this UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database or embedder unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	if !exists {
		writeError(c, APIError{Code: API_ERROR_NOT_FOUND, Message: API_ERROR_MSG_SIGNAL_NOT_FOUND})
		return
	}
	page_out, err := r.DB.QuerySupportingSips(c.Request.Context(), params.ID, *filters, *page_req)
	if err != nil {
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// @Summary List intelligence sources
// @Description Find provenance Source records for filtering and citation. `q` is case-insensitive metadata matching across domain, name, and URL; it is not semantic search. Optional metadata may be omitted when unavailable.
// @Tags Sources
// @Produce json
// @Param q query string false "Case-insensitive Source metadata search." maxlength(1024)
// @Param domains query []string false "Exact Source domain-name filter (CSV)." collectionFormat(csv)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SourceCollectionResponse "Source collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, page token, or parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, NewSourceDocuments(page_out.Items), page_req.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getSource godoc
// @Summary Retrieve one intelligence Source
// @Description Retrieve Source provenance metadata by UUID. It does not return Events published by the Source; use `GET /events?source_ids={source_id}` for that use case. Description, favicon, and RSS feed metadata may be omitted when unavailable.
// @Tags Sources
// @Produce json
// @Param source_id path string true "Source UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Output serialization. JSON is canonical; YAML and TOON are token-optimized for MCP and AI-agent clients." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SourceItemResponse "Source detail envelope"
// @Failure 400 {object} ErrorResponse "Malformed UUID"
// @Failure 404 {object} ErrorResponse "No Source with this UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	if source.IsZero() {
		writeError(c, APIError{Code: API_ERROR_NOT_FOUND, Message: API_ERROR_MSG_SOURCE_NOT_FOUND})
		return
	}
	writeItem(c, NewSourceDocument(&source), item.ResponseType)
}

// getTags godoc
// @Summary Discover tag filters for Espresso intelligence
// @Description List persisted tag labels available for Event and Signal filtering. Tag filters use fuzzy text matching, so clients can discover useful label vocabulary here without requiring a canonical taxonomy.
// @Description **Pagination**: follow `pagination.next_page`; `limit` defaults to 20 and is capped at 100.
// @Description **Formats**: JSON is canonical; YAML and TOON are token-optimized serializations for MCP and AI-agent clients.
// @Tags Tags
// @Produce json
// @Param q query string false "Case-insensitive substring or prefix match." maxlength(1024)
// @Param resource query []string false "Optional resource scope (CSV): event, signal." collectionFormat(csv)
// @Param response_type query string false "Output serialization." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} DiscoveryValueCollectionResponse "Tag value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, page token, or response_type"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	next := encodeNextCursor(page_out.NextCursor)
	items := datautils.Transform(page_out.Items, func(tag *string) db.TagValue {
		return db.TagValue{Value: *tag}
	})
	writePage(c, items, page.Limit, params.Cursor, next, params.ResponseType)
}

// getEntities godoc
// @Summary Discover Event entities
// @Description List exact company and people names available for Event filtering. Values are normalized snake_case filter names, not canonical entity IDs or profiles.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter." maxlength(1024)
// @Param types query []string false "Entity types (CSV): company, people." collectionFormat(csv)
// @Param response_type query string false "Output serialization." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} DiscoveryValueCollectionResponse "Entity value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, page token, or parameters"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, page_out.Items, page.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getRegions godoc
// @Summary Discover Event regions
// @Description List exact region names available for Event filtering. Values are normalized snake_case names, not structured geography or canonical place records.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter." maxlength(1024)
// @Param response_type query string false "Output serialization." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} DiscoveryValueCollectionResponse "Region value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, page token, or parameters"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, page_out.Items, page.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEventTypes godoc
// @Summary Discover Event types
// @Description List exact Event type names available for Event filtering. These values can be supplied to `event_types`; `categories` is a separate exact category filter.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter." maxlength(1024)
// @Param response_type query string false "Output serialization." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param page query string false "Opaque page token. Follow pagination.next_page for the next page."
// @Success 200 {object} DiscoveryValueCollectionResponse "Event type value collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid limit, page token, or parameters"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
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
		writeError(c, APIError{Code: API_ERROR_DB_ERROR, Message: API_ERROR_MSG_OUR_BAD})
		return
	}
	writePage(c, page_out.Items, page.Limit, params.Cursor, encodeNextCursor(page_out.NextCursor), params.ResponseType)
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

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", config.health)

	protected := router.Group("/")
	protected.Use(config.apiKeyMiddleware)
	protected.GET("/tags", config.getTags)
	protected.GET("/entities", config.getEntities)
	protected.GET("/regions", config.getRegions)
	protected.GET("/event-types", config.getEventTypes)
	protected.GET("/events", config.getEvents)
	protected.GET("/events/:id", config.getEvent)
	protected.GET("/events/:id/signals", config.getEventSignals)
	protected.GET("/events/:id/evidence", config.getEventEvidence)
	protected.GET("/signals", config.getSignals)
	protected.GET("/signals/:id", config.getSignal)
	protected.GET("/signals/:id/events", config.getSignalEvents)
	protected.GET("/sources", config.getSources)
	protected.GET("/sources/:id", config.getSource)

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
	writeError(c, APIError{Code: API_ERROR_UNAUTHORIZED, Message: API_ERROR_MSG_MISSING_API_KEY})
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
