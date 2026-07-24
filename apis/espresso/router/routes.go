// @title 			Espresso API & MCP
// @version 		0.1
// @description 	MCP-ready business intelligence over curated "sips" for agents, dashboards, and automated research workflows.
// @description 	A **sip** is one unit of intelligence. `event` records describe observed business or market developments. `signal` records synthesize higher-level implications from related events and actions. `action` records are lower-level source facts used by the ingestion pipeline.
// @description 	Agent workflow: (1) listTags to discover filter vocabulary; (2) searchEvents for time-ordered developments; (3) searchSignals for synthesized implications; (4) getRelatedSips to follow `same_as` duplicates or `derived_from` intelligence chains.
// @description 	Conventions: Auth is optional at the backend but API-key protected through the gateway. Pagination uses `limit` default 16 max 128 and `offset` default 0. Empty result sets return HTTP 204, not an error. All sip IDs are UUID strings such as `339366bc-464d-582f-8132-6875ccc814d2`.
// @description 	Response formats: use `response_type=json` for structured application data. Use `response_type=text` for MCP/LLM context; it returns the same underlying records as compact field-per-line plain text with lower token overhead.
// @schemes 		https
// @license.name 	MIT
// @contact.name 	Project Cafecito
// @contact.url  	http://cafecito.tech
// @contact.email 	soumitsrah@cafecito.tech
package router

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maypok86/otter/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/embedding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	MIN_WINDOW          = 1
	DEFAULT_WINDOW      = 7 // DAYS
	DEFAULT_ACCURACY    = 0.5
	DEFAULT_LIMIT       = 16
	MAX_LIMIT           = 128
	FAVICON_PATH        = "images/espresso-insta-dark.png"
	DEFAULT_CONCURRENCY = 512
)

const (
	_CACHE_SIZE = 1000
	_CACHE_TTL  = 30 * time.Minute
)

const (
	_EMBEDDER_ERROR = "Embedder just died. Retry in a bit."
	_DB_ERROR       = "DB just died. Retry in a bit."
)

// baseQueryParams holds shared pagination query parameters for list endpoints.
type baseQueryParams struct {
	ResponseType string `form:"response_type,default=json" binding:"oneof=json text"`
	Limit        int    `form:"limit,default=16" binding:"min=1,max=128"`
	Offset       int    `form:"offset" binding:"min=0"`
}

// sipsQueryParams holds shared filter and search parameters for /events and /signals.
type sipsQueryParams struct {
	IDs  []string  `form:"ids" collection_format:"csv"`
	From time.Time `form:"from" time_format:"2006-01-02" swaggertype:"string" format:"date"`
	Q    string    `form:"q" binding:"max=1024"`
	Acc  float64   `form:"acc,default=0.5" binding:"min=0,max=1"`
	Tags []string  `form:"tags" collection_format:"csv"`
	baseQueryParams
}

type relatedURIParams struct {
	Relationship string `uri:"relationship" binding:"required,oneof=same_as derived_from"`
}

// relatedQueryParams holds query parameters for GET /related/{relationship}.
type relatedQueryParams struct {
	IDs []string `form:"ids" collection_format:"csv" binding:"required"`
	baseQueryParams
}

// Configuration wires database, embedding, auth, and caching dependencies into HTTP handlers.
type Configuration struct {
	DB       *db.Cupboard
	Embedder embedding.Embedder
	APIKeys  map[string]string
	queue    chan int
	cache    *otter.Cache[string, []float32]
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
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// getTags godoc
// @Summary Discover tag filters for Espresso intelligence
// @Description Returns a paginated, alphabetically sorted list of unique tag strings extracted from event and signal sips.
// @Description **When to use**: call this before searchEvents or searchSignals when an agent needs valid tag vocabulary instead of guessing filter values.
// @Description **Filter behavior**: tags returned here can be passed to `tags` on `/events` and `/signals`; multiple tag values are treated as an inclusive AND by those search endpoints.
// @Description **Response formats**: `response_type=json` returns a JSON string array. `response_type=text` returns one comma-separated plain-text string for lower-token MCP context.
// @Description **Pagination**: use `offset` to walk the full vocabulary when `limit` is smaller than the total number of tags.
// @Tags Tags
// @Produce json
// @Produce plain
// @Param response_type query string false "Output format. json returns a JSON string array; text returns the same tags as comma-separated plain text for lower token cost." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param offset query int false "Number of tags to skip. Default 0." minimum(0)
// @Success 200 {array} string "Tag strings when response_type=json; comma-separated tags when response_type=text"
// @Success 204 "No tags found (empty result, not an error)"
// @Failure 400 {object} ErrorResponse "Invalid limit, offset, or response_type"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID listTags
// @Router /tags [get]
func (r *Configuration) getTags(c *gin.Context) {
	var params baseQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := r.DB.GetTags(c.Request.Context(), db.Pagination{Limit: params.Limit, Offset: params.Offset})
	returnResponse(c, data, err, params.ResponseType)
}

func convertStringsToUUIDs(strings []string) ([]uuid.UUID, error) {
	errs := make([]error, 0, len(strings))
	uuids := make([]uuid.UUID, 0, len(strings))
	for _, raw := range strings {
		if id, err := uuid.Parse(raw); err != nil {
			errs = append(errs, errors.New("invalid id: "+raw))
		} else {
			uuids = append(uuids, id)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return uuids, nil
}

func (config *Configuration) extractSipsParams(c *gin.Context) (*db.Condition, *db.Pagination, string) {
	var input sipsQueryParams
	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil, ""
	}
	conditions := db.Condition{
		Created: input.From,
		Tags:    input.Tags,
	}
	if input.From.IsZero() {
		conditions.Created = time.Now().AddDate(0, 0, -DEFAULT_WINDOW) // default to last 7 days if no published/trending filter provided
	}
	if len(input.IDs) > 0 {
		ids, err := convertStringsToUUIDs(input.IDs)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil, nil, ""
		}
		conditions.IDs = ids
	}
	if input.Q != "" {
		conditions.Embedding = config.Embedder.EmbedQuery(c, input.Q)
		if len(conditions.Embedding) == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": _EMBEDDER_ERROR})
			return nil, nil, ""
		}
		config.cache.Set(input.Q, conditions.Embedding)
		distance := (1 - input.Acc) * 2
		conditions.Distance = &distance
	}
	return &conditions, &db.Pagination{Limit: input.Limit, Offset: input.Offset}, input.ResponseType
}

// getEvents godoc
// @Summary Search event intelligence
// @Description Returns event-kind sips. Searches without `q` are sorted by `created` descending; semantic searches with `q` are ranked by cosine distance, most similar first.
// @Description **When to use**: retrieve concrete developments, incidents, company actions, policy changes, market moves, or other observed business events before moving to higher-level signals.
// @Description **Search modes**: use `ids` for exact UUID lookup, `tags` for inclusive-AND tag filtering, `q` + `acc` for semantic search, and `from` to set the oldest creation date. These filters can be combined.
// @Description **Default time window**: when `from` is omitted, the service uses its default recent window, currently about the last 7 days.
// @Description **Response shape**: JSON responses are flattened digest objects with `id` and `reported` added by the router. Stable fields include `briefing`, `event_type`, `actions`, `people`, `regions`, `cross_domain_impacts`, `future_outlook`, `impact_level`, and `tags`; additional pipeline-specific keys may appear.
// @Description **Agent format**: use `response_type=text` for compact field-per-line records when feeding an LLM or MCP client. Use JSON when the caller needs structured parsing.
// @Tags Events
// @Produce json
// @Produce plain
// @Param ids query []string false "Exact event sip UUIDs to fetch (CSV). Use when following references or retrieving known records." collectionFormat(csv)
// @Param tags query []string false "Tag filters (CSV). Multiple values are inclusive AND, so every supplied tag must match." collectionFormat(csv)
// @Param q query string false "Natural-language semantic search query. Max 1024 characters; requires the embedder." maxlength(1024)
// @Param acc query number false "Match strictness for q. 0.0=broad, 1.0=strict. Default 0.5." default(0.5) minimum(0) maximum(1)
// @Param from query string false "Only include events created on or after this date (YYYY-MM-DD). Defaults to the recent window when omitted." format(date)
// @Param response_type query string false "Output format. json returns flattened digest objects; text returns compact plain-text records for LLM/MCP context." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param offset query int false "Number of events to skip. Default 0." minimum(0)
// @Success 200 {array} Event "Event digests when response_type=json; plain-text event blocks when response_type=text"
// @Success 204 "No matching events (empty result, not an error)"
// @Failure 400 {object} ErrorResponse "Invalid query parameters or malformed UUID in ids"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database or embedder unavailable; retry"
// @ID searchEvents
// @Router /events [get]
func (r *Configuration) getEvents(c *gin.Context) {
	conditions, page, response_type := r.extractSipsParams(c)
	if conditions == nil || page == nil {
		return
	}
	conditions.Kinds = db.EVENTS
	items, err := r.DB.QuerySips(c.Request.Context(), *conditions, *page)
	returnResponse(c, items, err, response_type)
}

// getSignals godoc
// @Summary Search synthesized signals
// @Description Returns signal-kind sips. Searches without `q` are sorted by `created` descending; semantic searches with `q` are ranked by cosine distance, most similar first.
// @Description **When to use**: retrieve synthesized business implications, forecasts, drivers, and cross-event patterns after or instead of searching raw events.
// @Description **Search modes**: use `ids` for exact UUID lookup, `tags` for inclusive-AND tag filtering, `q` + `acc` for semantic search, and `from` to set the oldest creation date. These filters can be combined.
// @Description **Default time window**: when `from` is omitted, the service uses its default recent window, currently about the last 7 days.
// @Description **Response shape**: JSON responses are flattened digest objects with `id` and `reported` added by the router. Stable fields include `briefing`, `events`, `drivers`, `impacts`, `impacted_domains`, `forecast`, `impact_level`, and `tags`; additional pipeline-specific keys may appear.
// @Description **Agent format**: use `response_type=text` for compact field-per-line records when feeding an LLM or MCP client. Use JSON when the caller needs structured parsing.
// @Tags Signals
// @Produce json
// @Produce plain
// @Param ids query []string false "Exact signal sip UUIDs to fetch (CSV). Use when following references or retrieving known records." collectionFormat(csv)
// @Param tags query []string false "Tag filters (CSV). Multiple values are inclusive AND, so every supplied tag must match." collectionFormat(csv)
// @Param q query string false "Natural-language semantic search query. Max 1024 characters; requires the embedder." maxlength(1024)
// @Param acc query number false "Match strictness for q. 0.0=broad, 1.0=strict. Default 0.5." default(0.5) minimum(0) maximum(1)
// @Param from query string false "Only include signals created on or after this date (YYYY-MM-DD). Defaults to the recent window when omitted." format(date)
// @Param response_type query string false "Output format. json returns flattened digest objects; text returns compact plain-text records for LLM/MCP context." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param offset query int false "Number of signals to skip. Default 0." minimum(0)
// @Success 200 {array} Signal "Signal digests when response_type=json; plain-text signal blocks when response_type=text"
// @Success 204 "No matching signals (empty result, not an error)"
// @Failure 400 {object} ErrorResponse "Invalid query parameters or malformed UUID in ids"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database or embedder unavailable; retry"
// @ID searchSignals
// @Router /signals [get]
func (r *Configuration) getSignals(c *gin.Context) {
	conditions, page, response_type := r.extractSipsParams(c)
	if conditions == nil || page == nil {
		return
	}
	conditions.Kinds = db.SIGNALS
	items, err := r.DB.QuerySips(c.Request.Context(), *conditions, *page)
	returnResponse(c, items, err, response_type)
}

// getRelated godoc
// @Summary Follow related intelligence records
// @Description Returns sips linked to one or more source UUIDs through the requested relationship.
// @Description **When to use**: after searchEvents or searchSignals, call this endpoint to expand context around a known sip, deduplicate equivalent records, or trace derived intelligence.
// @Description **Relationships**: `same_as` returns equivalent or duplicate records. `derived_from` returns downstream records generated from, or based on, the supplied source sip IDs.
// @Description **Input**: `ids` is required and must contain one or more RFC 4122 UUID strings. The `relationship` path value must be exactly `same_as` or `derived_from`.
// @Description **Response shape**: each result is a flattened event or signal digest with `id` and `reported` added by the router. Additional pipeline-specific keys may appear.
// @Description **Agent format**: use `response_type=text` for compact field-per-line related records; use JSON when the caller needs structured parsing.
// @Tags Related
// @Produce json
// @Produce plain
// @Param relationship path string true "Relationship to traverse. same_as finds equivalent records; derived_from follows generated intelligence." Enums(same_as, derived_from)
// @Param ids query []string true "Source sip UUIDs (CSV). Example: b07049b5-54c0-50b0-a620-d3aea3f8a173" collectionFormat(csv)
// @Param response_type query string false "Output format. json returns flattened digest objects; text returns compact plain-text records for LLM/MCP context." Enums(json, text) default(json)
// @Param limit query int false "Page size. Default 16, max 128." default(16) minimum(1) maximum(128)
// @Param offset query int false "Number of related records to skip. Default 0." minimum(0)
// @Success 200 {array} Event "Related event/signal digests when response_type=json; plain-text related record blocks when response_type=text"
// @Success 204 "No related sips found (empty result, not an error)"
// @Failure 400 {object} ErrorResponse "Missing ids, invalid relationship, invalid response_type, or malformed UUID"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID getRelatedSips
// @Router /related/{relationship} [get]
func (r *Configuration) getRelated(c *gin.Context) {
	var uri relatedURIParams
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var query relatedQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(query.IDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ids is required"})
		return
	}
	ids, err := convertStringsToUUIDs(query.IDs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := r.DB.QueryRelatedSips(
		c.Request.Context(),
		db.Condition{
			IDs:          ids,
			Relationship: strings.ToUpper(uri.Relationship),
		},
		db.Pagination{Limit: query.Limit, Offset: query.Offset},
	)
	returnResponse(c, items, err, query.ResponseType)
}

func returnResponse[T any](c *gin.Context, items []T, err error, response_type string) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": _DB_ERROR})
		return
	}
	if len(items) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	if sips, ok := any(items).([]db.Sip); ok {
		if response_type == "text" {
			c.String(http.StatusOK, SipsToText(sips))
			return
		}
		c.JSON(http.StatusOK, sipsToDigest(sips))
		return
	} else if tags, ok := any(items).([]string); ok {
		if response_type == "text" {
			c.String(http.StatusOK, strings.Join(tags, ", "))
			return
		}
	}
	c.JSON(http.StatusOK, items)
}

func NewRouter(db *db.Cupboard, embedder embedding.Embedder, api_keys map[string]string, max_concurrent_requests int) *gin.Engine {
	if max_concurrent_requests <= 0 {
		max_concurrent_requests = DEFAULT_CONCURRENCY // default to 100 if not set or invalid
	}
	config := &Configuration{
		DB:       db,
		Embedder: embedder,
		APIKeys:  api_keys,
		queue:    make(chan int, max_concurrent_requests),
		cache: otter.Must(&otter.Options[string, []float32]{
			MaximumSize:      _CACHE_SIZE,
			ExpiryCalculator: otter.ExpiryAccessing[string, []float32](_CACHE_TTL),
		}),
	}

	router := gin.New()
	// JSON access logs and recovery using zerolog
	router.Use(
		// logger
		requestLogger,
		// recovery
		gin.Recovery(),
		// cors
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: false,
			MaxAge:           24 * time.Hour,
		}),
	)

	// Swagger / OpenAPI endpoints
	// NOTE: run `swag init` to generate docs (package `docs`) before using the UI.
	// Serve Swagger UI and point it at the generated spec in assets/docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", config.health)
	router.StaticFile("favicon.ico", FAVICON_PATH)

	// protected group
	protected := router.Group("/")
	protected.Use(config.apiKeyMiddleware, config.concurrencyMiddleware)
	protected.GET("/tags", config.getTags)
	protected.GET("/events", config.getEvents)
	protected.GET("/signals", config.getSignals)
	protected.GET("/related/:relationship", config.getRelated)
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
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API Key"})
}

func (r *Configuration) concurrencyMiddleware(c *gin.Context) {
	if r.queue != nil {
		r.queue <- 1
		defer func() { <-r.queue }()
	}
	c.Next()
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
