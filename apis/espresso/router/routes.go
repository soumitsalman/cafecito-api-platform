// @title 			Espresso API & MCP
// @version 		0.2
// @description 	MCP-ready business intelligence over curated intelligence records for agents, dashboards, and automated research workflows.
// @description 	An **Event** is any record whose kind starts with `event` (canonical `event`, plus `event:news`, `event:blog`, `event:post`, `event:site`, `event:social`). A **Signal** is a synthesized conclusion derived from Events. The internal storage word `sip` is not part of the public vocabulary.
// @description 	Agent workflow: (1) listTags to discover filter vocabulary; (2) searchEvents for developments; (3) getEvent and getEventEvidence to trace source coverage; (4) searchSignals and getSignalEvents to trace synthesized conclusions.
// @description 	Conventions: Auth is optional at the backend but API-key protected through the gateway. Collections use cursor pagination: `limit` default 20 max 128, and an opaque `cursor` returned as `next_cursor`. Empty collections return HTTP 200 with `data: []`. Missing detail resources return 404. All IDs are RFC 4122 UUID strings.
// @description Response formats: use `response_type=json` for structured application data. Use `response_type=text` for MCP/LLM context; it returns the same non-empty digest members as compact field-per-line plain text with `---` record delimiters. Public Event and Signal payloads do not synthesize storage id, created, kind, representation, or object fields.
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

	"github.com/alpkeskin/gotoon"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	MAX_LIMIT        = 128
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

// createEventFilters converts EventSearchParams into typed Filters.
func (p *EventSearchParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs:          p.IDs,
		SourceIDs:    p.SourceIDs,
		EventTypes:   p.EventTypes,
		ImpactLevels: p.ImpactLevels,
		Companies:    p.Companies,
		People:       p.People,
		Products:     p.Products,
		Regions:      p.Regions,
		Tags:         p.Tags,
		CreatedFrom:  p.From,
		CreatedTo:    p.To,
	}
	if len(p.Q) > 0 {
		filters.Embedding = config.Embedder.EmbedQuery(c, p.Q)
		if len(filters.Embedding) == 0 {
			return nil, APIError{Code: API_ERROR_EMBEDDING_ERROR, Message: API_ERROR_MSG_OUR_BAD}
		}
		distance := (1 - p.Acc) * 2
		filters.Distance = &distance
	}
	return filters, nil
}

// buildSignalFilters converts signalSearchParams into typed SipFilters.
func (p *SignalSearchParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		IDs:             p.IDs,
		ImpactLevels:    p.ImpactLevels,
		ImpactedDomains: p.ImpactedDomains,
		Tags:            p.Tags,
		CreatedFrom:     p.From,
		CreatedTo:       p.To,
	}
	if len(p.Q) > 0 {
		filters.Embedding = config.Embedder.EmbedQuery(c, p.Q)
		if len(filters.Embedding) == 0 {
			return nil, APIError{Code: API_ERROR_EMBEDDING_ERROR, Message: API_ERROR_MSG_OUR_BAD}
		}
		distance := (1 - p.Acc) * 2
		filters.Distance = &distance
	}
	return filters, nil
}

func (params *EventEvidenceParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		SourceIDs:   params.SourceIDs,
		CreatedFrom: params.From,
		CreatedTo:   params.To,
	}
	return filters, nil
}

func (params *EventSignalsParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		ImpactLevels:    params.ImpactLevels,
		ImpactedDomains: params.ImpactedDomains,
		Tags:            params.Tags,
		CreatedFrom:     params.From,
		CreatedTo:       params.To,
	}
	return filters, nil
}

func (params *SignalEventsParams) createFilters(c *gin.Context, config *Configuration) (*db.Filters, error) {
	filters := &db.Filters{
		ImpactLevels:    params.ImpactLevels,
		ImpactedDomains: params.ImpactedDomains,
		Tags:            params.Tags,
		EventTypes:      params.EventTypes,
		CreatedFrom:     params.From,
		CreatedTo:       params.To,
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
		Pagination: Pagination{Limit: limit, NextCursor: next_cursor},
		Meta:       ResponseMeta{AsOf: time.Now().UTC()},
	}
	switch response_type {
	case "yaml":
		c.YAML(http.StatusOK, response)
	case "toon":
		if encoded, err := gotoon.Encode(response); err != nil {
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
		if encoded, err := gotoon.Encode(response); err != nil {
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
// @Summary Search Event-family records
// @Description Which Event-family records match my question and filters? Returns every record where kind LIKE event% (canonical event plus event:news, event:blog, event:post, event:site, event:social). Each public item is the raw non-empty flattened digest object; storage fields such as id, created, kind, representation, and object are not synthesized.
// @Description **When to use**: retrieve concrete developments, incidents, company actions, policy changes, or market moves before moving to higher-level signals.
// @Description **Search modes**: `q` + `acc` for semantic search (default sort becomes `relevance`); without `q` records are sorted by `created_at` descending (`recent`).
// @Description **Time**: `from`/`to` are inclusive bounds on record `created_at` until an occurrence-time field exists. They do not claim Event occurrence time.
// @Description **Tags**: match on persisted tags uses overlap (any supplied tag).
// @Description **Response shape**: raw non-empty flattened digest members. Use the detail route for follow-up links and relation counts.
// @Description **Agent format**: use `response_type=text` for compact field-per-line records with `---` delimiters when feeding an LLM or MCP client.
// @Tags Events
// @Produce json
// @Produce plain
// @Param q query string false "Natural-language semantic search query. Max 1024 characters; requires the embedder. When present, default sort is relevance." maxlength(1024)
// @Param acc query number false "Match strictness for q. 0.0=broad, 1.0=strict. Default 0.5." default(0.5) minimum(0) maximum(1)
// @Param from query string false "Only include records created on or after this date (RFC3339 timestamp)." format(date-time)
// @Param to query string false "Only include records created on or before this date (RFC3339 timestamp)." format(date-time)
// @Param ids query []string false "Restrict to known Event-family IDs (CSV). Prefer the detail route for one ID." collectionFormat(csv)
// @Param event_types query []string false "Allowlisted match against digest.event_type (CSV)." collectionFormat(csv)
// @Param impact_levels query []string false "Match digest.impact_level (CSV): low, medium, high." collectionFormat(csv)
// @Param companies query []string false "Match digest.companies array (CSV)." collectionFormat(csv)
// @Param people query []string false "Match digest.people array (CSV)." collectionFormat(csv)
// @Param products query []string false "Match digest.products array (CSV)." collectionFormat(csv)
// @Param regions query []string false "Match digest.regions array (CSV)." collectionFormat(csv)
// @Param source_ids query []string false "Match the persisted source UUID, including direct SAME_AS evidence coverage (CSV)." collectionFormat(csv)
// @Param tags query []string false "Match persisted tags (CSV). tag_mode=any uses overlap; tag_mode=all requires every supplied tag." collectionFormat(csv)
// @Param response_type query string false "Output format. json returns the collection envelope; text returns compact plain-text records." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Page size. Default 20, max 128." default(20) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor returned as next_cursor. Clients must not construct or inspect it."
// @Success 200 {object} SipCollectionResponse "Event-family records when response_type=json; plain-text event blocks when response_type=text"
// @Failure 400 {object} APIError "Invalid query parameters, malformed UUID, or malformed cursor"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database or embedder unavailable; retry"
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
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEvent godoc
// @Summary Retrieve one Event-family record
// @Description What Event-family record has this UUID? Retrieves any record whose kind starts with `event` by UUID and returns its raw non-empty flattened digest members plus detail metadata links and relation counts. Returns 404 when no such record exists.
// @Tags Events
// @Produce json
// @Produce plain
// @Param event_id path string true "Event-family record UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Output format. json returns the detail envelope; text returns a compact plain-text record." Enums(json, text) default(json)
// @Success 200 {object} SipItemResponse "The Event-family record"
// @Failure 400 {object} APIError "Malformed UUID"
// @Failure 404 {object} APIError "No Event-family record with this UUID"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database unavailable; retry"
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
// @Summary Retrieve Event evidence
// @Description Which source-specific records support this Event? Returns a bare JSON list containing the requested Event and every direct SAME_AS Event-family neighbour in both relation orientations. Each item contains only event_id, created, source_id, url, and base_url. The default scope is direct_same_as: a bounded claim about current data, not a complete transitive equivalence closure.
// @Tags Events
// @Produce json
// @Produce plain
// @Param id path string true "Event-family record UUID (RFC 4122)." format(uuid)
// @Param source_ids query []string false "Restrict evidence to selected sources (CSV)." collectionFormat(csv)
// @Param from query string false "Inclusive evidence created_at lower bound (RFC3339 timestamp)." format(date-time)
// @Param to query string false "Inclusive evidence created_at upper bound (RFC3339 timestamp)." format(date-time)
// @Param response_type query string false "Output format. json returns the bare evidence list; text returns the same records as compact plain text." Enums(json, text) default(json)
// @Success 200 {object} EventEvidenceCollectionResponse "Event evidence collection"
// @Failure 400 {object} APIError "Malformed UUID or invalid parameters"
// @Failure 404 {object} APIError "No Event-family record with this UUID"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database unavailable; retry"
// @ID getEventEvidence
// @Router /events/{id}/evidence [get]
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
	writePage(c, evidence, page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEventSignals godoc
// @Summary Retrieve Signals derived from an Event
// @Description Which broader conclusions use this Event as evidence? Returns Signals whose DERIVED_FROM edges target the requested Event or any of its direct SAME_AS equivalents. The caller does not need to know which Event-family record was used as the relation target. Returns 404 when the Event does not exist; 200 with an empty collection when it exists but has no Signals.
// @Tags Events
// @Produce json
// @Produce plain
// @Param id path string true "Event-family record UUID (RFC 4122)." format(uuid)
// @Param from query string false "Inclusive Signal created_at lower bound (RFC3339 timestamp)." format(date-time)
// @Param to query string false "Inclusive Signal created_at upper bound (RFC3339 timestamp)." format(date-time)
// @Param impact_levels query []string false "Match Signal digest.impact_level (CSV)." collectionFormat(csv)
// @Param impacted_domains query []string false "Match Signal digest.impacted_domains array (CSV)." collectionFormat(csv)
// @Param tags query []string false "Match persisted tags (CSV). tag_mode=any uses overlap; tag_mode=all requires every supplied tag." collectionFormat(csv)
// @Param response_type query string false "Output format. json returns the collection envelope; text returns compact plain-text records." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 20, max 128." default(20) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor returned as next_cursor."
// @Success 200 {object} SipCollectionResponse "Signals derived from this Event"
// @Failure 400 {object} APIError "Malformed UUID, invalid cursor, or invalid parameters"
// @Failure 404 {object} APIError "No Event-family record with this UUID"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database or embedder unavailable; retry"
// @ID getEventSignals
// @Router /events/{id}/signals [get]
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
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getSignals godoc
// @Summary Search synthesized Signals
// @Description Which synthesized conclusions or forecasts match my question? Returns signal-kind records. A Signal is returned intelligence, not a saved monitoring definition.
// @Description **Search modes**: `q` + `acc` for semantic search (default sort becomes `relevance`); without `q` records are sorted by `created_at` descending (`recent`).
// @Description **Time**: `from`/`to` are inclusive bounds on Signal `created_at`.
// @Description **Tags**: match on persisted tags uses overlap (any supplied tag).
// @Tags Signals
// @Produce json
// @Produce plain
// @Param q query string false "Natural-language semantic search query. Max 1024 characters." maxlength(1024)
// @Param acc query number false "Match strictness for q. 0.0=broad, 1.0=strict. Default 0.5." default(0.5) minimum(0) maximum(1)
// @Param from query string false "Inclusive Signal created_at lower bound (RFC3339 timestamp)." format(date-time)
// @Param to query string false "Inclusive Signal created_at upper bound (RFC3339 timestamp)." format(date-time)
// @Param ids query []string false "Restrict to known Signal IDs (CSV)." collectionFormat(csv)
// @Param impact_levels query []string false "Match Signal digest.impact_level (CSV)." collectionFormat(csv)
// @Param impacted_domains query []string false "Match Signal digest.impacted_domains array (CSV)." collectionFormat(csv)
// @Param tags query []string false "Match persisted tags (CSV). tag_mode=any uses overlap; tag_mode=all requires every supplied tag." collectionFormat(csv)
// @Param tag_mode query string false "Tag matching mode." Enums(any, all) default(any)
// @Param sort query string false "Sort order. recent=created_at desc (default), relevance=semantic distance asc (requires q)." Enums(recent, relevance) default(recent)
// @Param response_type query string false "Output format. json returns the collection envelope; text returns compact plain-text records." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 20, max 128." default(20) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor returned as next_cursor."
// @Success 200 {object} SipCollectionResponse "Signal records"
// @Failure 400 {object} APIError "Invalid query parameters, malformed UUID, or malformed cursor"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database or embedder unavailable; retry"
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
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// @Summary Retrieve one Signal
// @Description What is the complete Signal with this UUID? Retrieves one signal-kind record by UUID. The detail payload links to supporting Events instead of inventing inline event references from unstructured digest strings. Returns 404 when no such Signal exists.
// @Tags Signals
// @Produce json
// @Produce plain
// @Param signal_id path string true "Signal record UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Output format. json returns the detail envelope; text returns a compact plain-text record." Enums(json, text) default(json)
// @Success 200 {object} SipItemResponse "The Signal record"
// @Failure 400 {object} APIError "Malformed UUID"
// @Failure 404 {object} APIError "No Signal with this UUID"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database unavailable; retry"
// @ID getSignal
// @Router /signals/{id} [get]
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
// @Summary Retrieve Event-family records supporting a Signal
// @Description Which Event-family records support this Signal? Follows DERIVED_FROM edges in the stored signal-to-event direction and returns raw non-empty flattened digest members from direct Event-family targets. Returns 404 when the Signal does not exist; 200 with an empty collection when it exists but has no supporting Events.
// @Tags Signals
// @Produce json
// @Produce plain
// @Param id path string true "Signal record UUID (RFC 4122)." format(uuid)
// @Param from query string false "Inclusive Event created_at lower bound (RFC3339 timestamp)." format(date-time)
// @Param to query string false "Inclusive Event created_at upper bound (RFC3339 timestamp)." format(date-time)
// @Param event_types query []string false "Match digest.event_type (CSV)." collectionFormat(csv)
// @Param impact_levels query []string false "Match digest.impact_level (CSV)." collectionFormat(csv)
// @Param tags query []string false "Match persisted tags (CSV). tag_mode=any uses overlap; tag_mode=all requires every supplied tag." collectionFormat(csv)
// @Param response_type query string false "Output format. json returns the collection envelope; text returns compact plain-text records." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 20, max 128." default(20) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor returned as next_cursor."
// @Success 200 {object} SipCollectionResponse "Event-family records supporting this Signal"
// @Failure 400 {object} APIError "Malformed UUID, invalid cursor, or invalid parameters"
// @Failure 404 {object} APIError "No Signal with this UUID"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database or embedder unavailable; retry"
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
	writePage(c, NewDigestDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// @Summary List intelligence sources
// @Description Which source records can I filter by or cite? Returns provenance records keyed by UUID. Optional source metadata may be null; missing optional metadata is not an error and the API does not fabricate names, domains, or URLs.
// @Tags Sources
// @Produce json
// @Param q query string false "Case-insensitive match against domain, site name, or base URL."
// @Param domains query []string false "Exact domain-name filter (CSV)." collectionFormat(csv)
// @Param limit query int false "Page size. Default 20, max 128." default(20) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor returned as next_cursor."
// @Param response_type query string false "Output format: json, yaml, or toon." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SourceCollectionResponse "Source records"
// @Failure 400 {object} APIError "Invalid limit, cursor, or parameters"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database unavailable; retry"
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
	writePage(c, NewSourceDocuments(page_out.Items), page_req.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getSource godoc
// @Summary Retrieve one intelligence source
// @Description What metadata is known about this source? This is provenance metadata. It does not return Events published by the source; use GET /events?source_ids={source_id} for that question. Returns 404 when no such source exists.
// @Tags Sources
// @Produce json
// @Param source_id path string true "Source record UUID (RFC 4122)." format(uuid)
// @Param response_type query string false "Output format: json, yaml, or toon." Enums(json, yaml, toon) default(json)
// @Success 200 {object} SourceItemResponse "The source record"
// @Failure 400 {object} APIError "Malformed UUID"
// @Failure 404 {object} APIError "No source with this UUID"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database unavailable; retry"
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
// @Description Which exact tag strings are valid filters? Returns a paginated, alphabetically sorted list of unique tag value objects extracted from Event and Signal records.
// @Description **When to use**: call this before searchEvents or searchSignals when an agent needs valid tag vocabulary instead of guessing filter values.
// @Description **Filter behavior**: tags returned here can be passed to `tags` on `/events` and `/signals`; matching uses overlap (any supplied tag).
// @Description **Response formats**: `response_type` accepts json, yaml, or toon. JSON returns a collection of `{ "value": "tag" }` objects.
// @Description **Pagination**: cursor-based. Use `next_cursor` to continue; `limit` default 20, max 128.
// @Tags Tags
// @Produce json
// @Produce plain
// @Param q query string false "Case-insensitive substring or prefix match."
// @Param resource query []string false "Optional kind scope (CSV): event, signal." collectionFormat(csv)
// @Param response_type query string false "Output format: json, yaml, or toon." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 20, max 128." default(20) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor returned as next_cursor. Clients must not construct or inspect it."
// @Success 200 {object} DiscoveryValueCollectionResponse "Tag value objects"
// @Failure 400 {object} APIError "Invalid limit, cursor, or response_type"
// @Failure 401 {object} APIError "Missing or invalid API key"
// @Failure 429 {object} APIError "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} APIError "Database unavailable; retry"
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
	writePage(c, items, page.Limit, next, params.ResponseType)
}

// getEntities godoc
// @Summary Discover Event entities
// @Description Returns distinct exact company and person strings stored in Event digests.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter."
// @Param types query []string false "Entity types: company, person." collectionFormat(csv)
// @Param response_type query string false "Output format." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor."
// @Success 200 {object} DiscoveryValueCollectionResponse
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
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
	writePage(c, page_out.Items, page.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getRegions godoc
// @Summary Discover Event regions
// @Description Returns distinct exact region strings stored in Event digests.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter."
// @Param response_type query string false "Output format." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor."
// @Success 200 {object} DiscoveryValueCollectionResponse
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
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
	writePage(c, page_out.Items, page.Limit, encodeNextCursor(page_out.NextCursor), params.ResponseType)
}

// getEventTypes godoc
// @Summary Discover Event types
// @Description Returns distinct exact event_type values stored in Event digests.
// @Tags Discovery
// @Produce json
// @Param q query string false "Case-insensitive substring filter."
// @Param response_type query string false "Output format." Enums(json, yaml, toon) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param cursor query string false "Opaque pagination cursor."
// @Success 200 {object} DiscoveryValueCollectionResponse
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
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

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", config.health)

	protected := router.Group("/")
	protected.Use(config.apiKeyMiddleware)
	protected.GET("/tags", config.getTags)
	protected.GET("/entities", config.getEntities)
	protected.GET("/regions", config.getRegions)
	protected.GET("/event-types", config.getEventTypes)
	protected.GET("/events", config.getEvents)
	protected.GET("/events/search", config.getEvents)
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
