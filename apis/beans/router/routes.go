package router

import (
	// "fmt" // Legacy route validation dependency.
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/embedding"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/* Legacy route constants retained for reference; excluded from the target Beans route proposal.
const (
	MIN_WINDOW          = 1
	DEFAULT_WINDOW      = 7 // DAYS
	DEFAULT_ACCURACY    = 0.5
	DEFAULT_LIMIT       = 16
	MAX_LIMIT           = 128
	FAVICON_PATH        = "./images/beans.png"
)
*/

// const (
// 	DEFAULT_CONCURRENCY = 512
// )

/* Legacy cache, search, and trend constants retained for reference; excluded from the target Beans route proposal.
const (
	_CACHE_SIZE = 1000
	_CACHE_TTL  = 30 * time.Minute
)

const (
	_EMBEDDER_ERROR     = "Embedder just died. Retry in a bit."
	_DB_ERROR           = "DB just died. Retry in a bit."
	_NEEDS_SEARCH_PARAM = "At least one search parameter is required (q, tags, categories, regions, entities)."
)

const (
	_BEAN_TREND_FIELDS = "likes, comments, shares, related, trend_score"
)
*/

/* Legacy route parameter types and URL propagation binding retained for reference; excluded from the target Beans route proposal.
// PaginationInput holds shared list-endpoint pagination query params.
// form default=16 and max=128 must stay in sync with DEFAULT_LIMIT and MAX_LIMIT.
type PaginationInput struct {
	Limit  int `form:"limit,default=16" binding:"min=1,max=128"`
	Offset int `form:"offset" binding:"min=0"`
}

func (p PaginationInput) ToDB() db.Pagination {
	return db.Pagination{Limit: p.Limit, Offset: p.Offset}
}

func normalizePagination(p *PaginationInput) error {
	if p.Limit == 0 {
		p.Limit = DEFAULT_LIMIT
	}
	if p.Limit > MAX_LIMIT {
		return fmt.Errorf("limit must be <= %d", MAX_LIMIT)
	}
	if p.Offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}
	return nil
}

// TagsInput describes pagination for tag discovery endpoints (/tags/*).
type TagsInput struct {
	PaginationInput
}

// PublishersInput describes query parameters for getPublishers (/sources).
type PublishersInput struct {
	// Sources: Publisher source IDs to resolve (CSV). Required for meaningful results.
	Sources []string `form:"sources" collection_format:"csv"`
	PaginationInput
}

// ArticlesInput describes shared query parameters for article list and search endpoints.
type ArticlesInput struct {
	// URLs: Optional list of article URLs to fetch directly (CSV).
	URLs []string `form:"urls" collection_format:"csv"`
	// Q: Free-text semantic/vector search query (max 512 chars).
	Q string `form:"q" binding:"max=512"`
	// Acc: Similarity accuracy threshold (0.0-1.0). Higher => stricter match.
	// Used to compute vector distance (distance = 1 - Acc).
	Acc float64 `form:"acc,default=0.5" binding:"min=0,max=1"`
	// ContentType: Optional content type filter (stored kind).
	ContentType string `form:"content_type" binding:"omitempty,oneof=blog contract earnings_report enforcement_action financial_report lawsuit news official_statement podcast post press_release research_paper site technical_documentation whitepaper"`
	// Categories: Filter results to one or more categories/topics (CSV).
	Categories []string `form:"categories" collection_format:"csv"`
	// Regions: Filter results to one or more geographic regions (CSV).
	Regions []string `form:"regions" collection_format:"csv"`
	// Entities: Filter results to one or more named entities (CSV).
	Entities []string `form:"entities" collection_format:"csv"`
	// Tags: Tag/keyword filters (CSV). Combined into a full-text query for tag matching.
	Tags []string `form:"tags" collection_format:"csv"`
	// Sources: Publisher/source IDs to include (CSV).
	Sources []string `form:"sources" collection_format:"csv"`
	// From: Start date for published/updated filtering (format YYYY-MM-DD).
	From time.Time `form:"from" time_format:"2006-01-02" swaggertype:"string" format:"date"`
	// FullContent: If true, include full article content in results (larger payload).
	FullContent bool `form:"full_content,default=false"`
	// PaginationInput: Embeds common pagination params (Limit, Offset).
	PaginationInput
}

type PropagationInput struct {
	// URLs lists seed article URLs to analyze for cross-outlet coverage and social mentions (1–128 items).
	URLs []string `form:"urls" json:"urls" collection_format:"csv" binding:"required,min=1,dive,url"`
}

func bindPropagationInput(c *gin.Context) (PropagationInput, error) {
	var input PropagationInput
	switch c.Request.Method {
	case http.MethodGet:
		if err := c.ShouldBindQuery(&input); err != nil {
			return PropagationInput{}, err
		}
	case http.MethodPost:
		if err := c.ShouldBindJSON(&input); err != nil {
			return PropagationInput{}, err
		}
	default:
		return PropagationInput{}, fmt.Errorf("method not allowed")
	}
	if len(input.URLs) > MAX_LIMIT {
		return PropagationInput{}, fmt.Errorf("urls must contain at most %d items", MAX_LIMIT)
	}
	return input, nil
}
*/

const (
	MIN_WINDOW              = 1
	DEFAULT_WINDOW          = 7 // DAYS
	DEFAULT_SCORE_THRESHOLD = 0.5
	DEFAULT_LIMIT           = 20
	MAX_LIMIT               = 100
)

const (
	API_ERROR_MSG_OUR_BAD            = "Our bad. Please try again later."
	API_ERROR_MSG_SOURCE_NOT_FOUND   = "Source not found."
	API_ERROR_MSG_TAG_NOT_FOUND      = "Tag not found."
	API_ERROR_MSG_ENTITY_NOT_FOUND   = "Entity not found."
	API_ERROR_MSG_REGION_NOT_FOUND   = "Region not found."
	API_ERROR_MSG_CATEGORY_NOT_FOUND = "Category not found."
	API_ERROR_MSG_COMPANY_NOT_FOUND  = "Company not found."
	API_ERROR_MSG_PRODUCT_NOT_FOUND  = "Product not found."
	API_ERROR_MSG_ARTICLE_NOT_FOUND  = "Article not found."
	API_ERROR_MSG_STORY_NOT_FOUND    = "Story not found."
)

type Configuration struct {
	DB       *db.PGSack
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

// toDBFilters converts router-level article filter params to db-level Filters.
func (p *articleFilterParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters := db.Filters{
		Kind:              p.ContentType,
		Sources:           p.Sources,
		ExcludeSources:    p.ExcludeSources,
		Domains:           p.Domains,
		ExcludeDomains:    p.ExcludeDomains,
		Authors:           p.Authors,
		Tags:              utils.NormalizeTags(p.Tags),
		Categories:        utils.NormalizeTags(p.Categories),
		ExcludeCategories: utils.NormalizeTags(p.ExcludeCategories),
		Sentiments:        utils.NormalizeTags(p.Sentiments),
		Entities:          utils.NormalizeTags(p.Entities),
		Regions:           utils.NormalizeTags(p.Regions),
		FullContent:       p.FullContent,
	}
	return &filters, nil
}

func (p *articleFeedParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (p *articleSearchParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters, _ := p.articleFeedParams.createFilters(c, r)
	filters.IDs = p.IDs
	filters.URLs = p.URLs
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (p *headlinesParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (p *storySearchParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	filters.FullContent = false
	if !p.From.IsZero() {
		filters.CreatedFrom = p.From
	}
	if !p.To.IsZero() {
		filters.CreatedTo = p.To
	}
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (p *storyArticleParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	filters.ClusterID = p.StoryID
	if !p.From.IsZero() {
		filters.CreatedFrom = p.From
	}
	if !p.To.IsZero() {
		filters.CreatedTo = p.To
	}
	return filters, nil
}

func (p *similarArticlesParams) createFilters(c *gin.Context, r *Configuration) (*db.Filters, error) {
	filters, err := p.articleFilterParams.createFilters(c, r)
	if err != nil {
		return nil, err
	}
	if !p.From.IsZero() {
		filters.CreatedFrom = p.From
	}
	if !p.To.IsZero() {
		filters.CreatedTo = p.To
	}
	return filters, nil
}

func (p *articleMentionsParams) createFilters() db.MentionFilters {
	return db.MentionFilters{
		Platforms:    p.Platforms,
		Forums:       p.Forums,
		ObservedFrom: p.From,
		ObservedTo:   p.To,
	}
}

func (p *vectorSearchParams) attachToFilters(c *gin.Context, config *Configuration, filters *db.Filters) error {
	q := strings.TrimSpace(p.Q)
	if q != "" {
		filters.Embedding = config.Embedder.EmbedQuery(c, q)
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

// writeCollection writes a typed collection envelope as JSON, or compact text when requested.
func writeCollection[T any](c *gin.Context, items []T, limit int, next_cursor *db.Cursor) {
	if items == nil {
		items = []T{}
	}
	response := CollectionResponse[T]{
		Data:       items,
		Pagination: Pagination{Limit: limit, NumResults: len(items)},
		Meta:       ResponseMeta{AsOf: time.Now().UTC()},
	}
	if next_cursor != nil {
		if cur, err := next_cursor.Encode(); err == nil {
			response.Pagination.NextCursor = &cur
		}
	}
	c.JSON(http.StatusOK, response)
}

// writeDetail writes a typed detail envelope as JSON, or compact text when requested.
func writeDetail[T any](c *gin.Context, item T) {
	c.JSON(http.StatusOK, DetailResponse[T]{Data: item})
}

func writeStoryDetail(c *gin.Context, item StoryDetailItem) {
	c.JSON(http.StatusOK, StoryDetailResponse{
		Data: item,
		Meta: ResponseMeta{AsOf: time.Now().UTC()},
	})
}

func writeStoryArticles(c *gin.Context, items []ArticleDocument, limit int, next_cursor *db.Cursor, story_id string) {
	if items == nil {
		items = []ArticleDocument{}
	}
	response := StoryArticleCollectionResponse{
		Data:       items,
		Pagination: Pagination{Limit: limit, NumResults: len(items)},
		Meta:       StoryArticleMeta{StoryID: story_id, AsOf: time.Now().UTC()},
	}
	if next_cursor != nil {
		if cur, err := next_cursor.Encode(); err == nil {
			response.Pagination.NextCursor = &cur
		}
	}
	c.JSON(http.StatusOK, response)
}

// writeError writes an APIError to the response.
// Uses InternalServerError for DB, Embedding, Encoding errors and default cases
func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if api_err, ok := err.(utils.APIError); ok {
		switch api_err.Code {
		case utils.API_ERROR_INVALID_REQUEST:
			status = http.StatusBadRequest
		case utils.API_ERROR_NOT_FOUND:
			status = http.StatusNotFound
		}
	}
	c.AbortWithStatusJSON(status, ErrorResponse{Error: err})
}

// writeScaffoldNotImplemented writes a 501 Not Implemented response for scaffold routes.
func writeScaffoldNotImplemented(c *gin.Context, route string) {
	writeError(c, utils.NewAPIError("not_implemented", route+" is not yet implemented"))
}

// writeScaffoldInvalidRequest writes a 400 Bad Request response from a binding error.
func writeScaffoldInvalidRequest(c *gin.Context, err error) {
	writeError(c, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error()))
}

// health godoc
// @Summary Check API health
// @Description Lightweight liveness probe. Use before other tools to confirm the service is reachable. No authentication required when API keys are disabled.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "status alive"
// @ID healthCheck
// @Router /health [get]
func (r *Configuration) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// B01 GET /articles/search
func (r *Configuration) searchArticles(c *gin.Context) {
	var params articleSearchParams
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

	page_out, err := r.DB.QueryBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_WITHOUT_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// B03 GET /articles/latest
func (r *Configuration) getLatestArticles(c *gin.Context) {
	var params articleFeedParams
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

	if params.From.IsZero() {
		filters.CreatedFrom = time.Now().AddDate(0, 0, -DEFAULT_WINDOW)
	} else {
		filters.CreatedFrom = params.From
	}
	if !params.To.IsZero() {
		filters.CreatedTo = params.To
	}

	page_out, err := r.DB.QueryBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_WITHOUT_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// B05 GET /articles/trending
func (r *Configuration) getTrendingArticles(c *gin.Context) {
	var params articleFeedParams
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

	if params.From.IsZero() {
		filters.UpdatedFrom = time.Now().AddDate(0, 0, -DEFAULT_WINDOW)
	} else {
		filters.UpdatedFrom = params.From
	}
	if !params.To.IsZero() {
		filters.UpdatedTo = params.To
	}

	page_out, err := r.DB.QueryTrendingBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_WITH_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryTrendingBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// getTopHeadlinesArticles is the B04 GET /articles/top-headlines target scaffold.
// Primary difference between this and getTrendingArticles and getTopHeadlines is the window of time.
// getTopHeadlines is always fixed within the last 24 hours. The returned article must be created and received social media traction in the last 24 hours.
// The result excluded summary and content
func (r *Configuration) getTopHeadlines(c *gin.Context) {
	var params headlinesParams
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

	filters.CreatedFrom = time.Now().AddDate(0, 0, -MIN_WINDOW)
	filters.UpdatedFrom = time.Now().AddDate(0, 0, -MIN_WINDOW)

	page_out, err := r.DB.QueryTrendingBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_HEADLINES)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryTrendingBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

func (r *Configuration) getArticle(c *gin.Context) {
	var params articleDetailParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	bean, err := r.DB.GetBean(c.Request.Context(), params.ID, params.FullContent)
	if err != nil {
		shared.LogError(err, "[ERROR] GetBean")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_ARTICLE_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeDetail(c, toArticleDetailItem(&bean, params.FullContent))
}

func (r *Configuration) getSimilarArticles(c *gin.Context) {
	var params similarArticlesParams
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

	page_out, err := r.DB.QuerySimilarBeans(c.Request.Context(), params.ID, *filters, *page_req, db.BEAN_COLUMNS_WITHOUT_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QuerySimilarBeans")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_ARTICLE_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

func (r *Configuration) getArticleMentions(c *gin.Context) {
	var params articleMentionsParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}

	page_out, err := r.DB.QueryMentions(c.Request.Context(), params.ID, params.createFilters(), *page_req)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryMentions")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_ARTICLE_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeCollection(c, toMentionDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

func (r *Configuration) getStories(c *gin.Context) {
	var params storySearchParams
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

	page_out, err := r.DB.QueryStories(c.Request.Context(), *filters, *page_req, params.MinArticleCount)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryStories")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toStoryItems(page_out.Items), page_req.Limit, page_out.NextCursor)
}

func (r *Configuration) dispatchStory(c *gin.Context) {
	raw := strings.TrimSpace(strings.TrimPrefix(c.Param("story_id"), "/"))
	if raw == "" {
		writeError(c, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "story_id is required"))
		return
	}
	if strings.HasSuffix(raw, "/articles") {
		candidate := strings.TrimSuffix(raw, "/articles")
		exists, err := r.DB.StoryExists(c.Request.Context(), raw)
		if err != nil {
			shared.LogError(err, "[ERROR] StoryExists")
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
			return
		}
		if !exists && candidate != "" {
			c.Params = gin.Params{{Key: "story_id", Value: candidate}}
			r.getStoryArticles(c)
			return
		}
	}
	c.Params = gin.Params{{Key: "story_id", Value: raw}}
	r.getStory(c)
}

func (r *Configuration) getStory(c *gin.Context) {
	var params storyPathParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	story, err := r.DB.GetStory(c.Request.Context(), params.StoryID)
	if err != nil {
		shared.LogError(err, "[ERROR] GetStory")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_STORY_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeStoryDetail(c, toStoryDetailItem(&story))
}

func (r *Configuration) getStoryArticles(c *gin.Context) {
	var params storyArticleParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	exists, err := r.DB.StoryExists(c.Request.Context(), params.StoryID)
	if err != nil {
		shared.LogError(err, "[ERROR] StoryExists")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if !exists {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_STORY_NOT_FOUND))
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

	page_out, err := r.DB.QueryBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_WITHOUT_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeStoryArticles(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor, params.StoryID)
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
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID listIntelligenceSources
// @Router /sources [get]
func (r *Configuration) getSources(c *gin.Context) {
	var params sourceListParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	page_req, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QuerySources(c.Request.Context(), params.Q, params.Domains, *page_req, db.SOURCE_COLUMNS_BASE)
	if err != nil {
		shared.LogError(err, "[ERROR] QuerySources")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toSourceDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
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
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Concurrency limit exceeded; retry shortly"
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID getIntelligenceSource
// @Router /sources/{source_id} [get]
func (r *Configuration) getSource(c *gin.Context) {
	var params sourceDetailParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}

	source, err := r.DB.GetSource(c.Request.Context(), params.ID)
	if err != nil {
		shared.LogError(err, "[ERROR] GetSource")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_SOURCE_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeDetail(c, toSourceDocument(&source))
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
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID listIntelligenceEntities
// @Router /entities [get]
func (r *Configuration) getEntities(c *gin.Context) {
	getTags(r, c, "entities", "entity")
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
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID listIntelligenceRegions
// @Router /regions [get]
func (r *Configuration) getRegions(c *gin.Context) {
	getTags(r, c, "regions", "region")
}

// getCategories godoc
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
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID listIntelligenceEventTypes
// @Router /event-types [get]
func (r *Configuration) getCategories(c *gin.Context) {
	getTags(r, c, "categories", "category")
}

// getSentiments godoc
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
// @Failure 500 {object} ErrorResponse "Database unavailable; retry"
// @ID listIntelligenceEventTypes
// @Router /event-types [get]
func (r *Configuration) getSentiments(c *gin.Context) {
	getTags(r, c, "sentiments", "sentiment")
}

func getTags(r *Configuration, c *gin.Context, db_tag_type string, response_tag_type string) {
	var params tagListParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	page, err := params.createPageRequest(c, r)
	if err != nil {
		writeError(c, err)
		return
	}

	page_out, err := r.DB.QueryTags(c.Request.Context(), strings.ToLower(strings.TrimSpace(params.Q)), db_tag_type, *page)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryTags")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toTagDocuments(page_out.Items, response_tag_type), page.Limit, page_out.NextCursor)
}

func NewRouter(db *db.PGSack, embedder embedding.Embedder, api_keys map[string]string) *gin.Engine {
	config := &Configuration{
		DB:       db,
		Embedder: embedder,
		APIKeys:  api_keys,
	}

	router := gin.New()
	// JSON logs and recovery using zerolog
	router.Use(
		requestLogger,
		gin.Recovery(),
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: false,
			MaxAge:           24 * time.Hour,
		}),
	)

	// Swagger / OpenAPI endpoints
	// NOTE: run `swag init` to generate docs (package `docs`) before using the UI.
	// Serve Swagger UI and point it at the generated spec in assets/docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", config.health)

	// Authenticated group
	protected := router.Group("/")
	protected.Use(config.apiKeyMiddleware)

	// TAGS discovery routes
	// protected.GET("/tags", config.getTags)
	protected.GET("/categories", config.getCategories)
	protected.GET("/entities", config.getEntities)
	protected.GET("/regions", config.getRegions)
	protected.GET("/sentiments", config.getSentiments)

	// SOURCES routes
	protected.GET("/sources", config.getSources)
	protected.GET("/sources/:id", config.getSource)

	// ARTICLES routes
	protected.GET("/articles/search", config.searchArticles)
	protected.GET("/articles/latest", config.getLatestArticles)
	protected.GET("/articles/trending", config.getTrendingArticles)
	protected.GET("/articles/:id/similar", config.getSimilarArticles)
	protected.GET("/articles/:id/mentions", config.getArticleMentions)
	protected.GET("/articles/:id", config.getArticle)

	// HEADLINES routes
	protected.GET("/headlines", config.getTopHeadlines)

	// STORIES routes. Wildcard captures URL-like story IDs (slashes); trailing /articles is membership.
	protected.GET("/stories", config.getStories)
	protected.GET("/stories/*story_id", config.dispatchStory)

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
	writeError(c, utils.NewAPIError(utils.API_ERROR_UNAUTHORIZED, "Missing API Key"))
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
