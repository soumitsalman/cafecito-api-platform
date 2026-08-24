package router

import (
	// "fmt" // Legacy route validation dependency.
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared/embedding"

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
func (p *articleFilterParams) createFilters(c *gin.Context, r *Configuration) (*db.BeanFilters, error) {
	filters := db.BeanFilters{
		Kind:              strings.ToLower(p.ContentType),
		Sources:           p.Sources,
		ExcludeSources:    p.ExcludeSources,
		Domains:           utils.NormalizeTexts(p.Domains),
		ExcludeDomains:    utils.NormalizeTexts(p.ExcludeDomains),
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

func (p *articleFeedParams) createFilters(c *gin.Context, r *Configuration) (*db.BeanFilters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (p *articleSearchParams) createFilters(c *gin.Context, r *Configuration) (*db.BeanFilters, error) {
	filters, _ := p.articleFeedParams.createFilters(c, r)
	filters.IDs = p.IDs
	filters.URLs = p.URLs
	filters.CreatedFrom = p.From
	filters.CreatedTo = utils.NormalizeEndOfDay(p.To)
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return filters, nil
}

func (p *similarArticlesParams) createFilters(c *gin.Context, r *Configuration) (*db.BeanFilters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	filters.CreatedFrom = p.From
	filters.CreatedTo = utils.NormalizeEndOfDay(p.To)
	return filters, nil
}

func (p *articleMentionsParams) createFilters(c *gin.Context, r *Configuration) (*db.MentionFilters, error) {
	return &db.MentionFilters{
		Platforms:    utils.NormalizeTexts(p.Platforms),
		Forums:       utils.NormalizeTexts(p.Forums),
		ObservedFrom: p.From,
		ObservedTo:   utils.NormalizeEndOfDay(p.To),
	}, nil
}

func (p *sourceSearchParams) createFilters(c *gin.Context, r *Configuration) (*db.SourceFilters, error) {
	return &db.SourceFilters{
		Q:       utils.NormalizeText(p.Q),
		Domains: utils.NormalizeTexts(p.Domains),
		IDs:     p.IDs,
	}, nil
}

func (p *storySearchParams) createFilters(c *gin.Context, r *Configuration) (*db.ClusterFilters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	filters.FullContent = false
	filters.CreatedFrom = p.From
	filters.CreatedTo = utils.NormalizeEndOfDay(p.To)
	if err := p.vectorSearchParams.attachToFilters(c, r, filters); err != nil {
		return nil, err
	}
	return &db.ClusterFilters{BeanFilters: *filters, MinBeanCount: p.MinArticleCount}, nil
}

func (p *storyArticleParams) createFilters(c *gin.Context, r *Configuration) (*db.BeanFilters, error) {
	filters, _ := p.articleFilterParams.createFilters(c, r)
	filters.ClusterID = p.ID
	filters.CreatedFrom = p.From
	filters.CreatedTo = utils.NormalizeEndOfDay(p.To)
	return filters, nil
}

func (p *vectorSearchParams) attachToFilters(c *gin.Context, config *Configuration, filters *db.BeanFilters) error {
	q := strings.TrimSpace(p.Q)
	if q != "" {
		filters.Embedding = config.Embedder.EmbedQuery(c, q)
		if len(filters.Embedding) == 0 {
			return utils.NewAPIError(utils.API_ERROR_EMBEDDING_ERROR, API_ERROR_MSG_OUR_BAD)
		}
		if p.ScoreThreshold > 0 {
			distance := (1 - p.ScoreThreshold) * 2
			filters.Distance = distance
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

func writeStoryArticles(c *gin.Context, items []ArticleDocument, limit int, next_cursor *db.Cursor, story_id uuid.UUID) {
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

// searchArticles godoc
// @Summary Search Articles
// @Description Returns Articles matching an optional relevance query, exact Article IDs or URLs, and filters. Without q, results are newest first.
// @Tags Articles
// @Produce json
// @Param q query string false "Optional relevance query." maxlength(512)
// @Param score_threshold query number false "Optional relevance threshold used with q." minimum(0) maximum(1)
// @Param ids query []string false "Exact Article UUIDs (CSV)." collectionFormat(csv)
// @Param urls query []string false "Exact Article URLs (CSV)." collectionFormat(csv)
// @Param content_type query string false "Stored Article type." Enums(blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,post,press_release,research_paper,site,technical_documentation,whitepaper)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param from query string false "UTC lower timestamp bound." format(date)
// @Param to query string false "UTC upper timestamp bound." format(date)
// @Param full_content query bool false "Include content when available." default(false)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} ArticleCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID searchArticles
// @Router /articles/search [get]
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

// getLatestArticles godoc
// @Summary List latest Articles
// @Description Returns Articles ordered newest first. Date bounds are not accepted.
// @Tags Articles
// @Produce json
// @Param q query string false "Optional relevance query. Requires score_threshold greater than zero." maxlength(512)
// @Param score_threshold query number false "Required when q is supplied." minimum(0) maximum(1)
// @Param content_type query string false "Stored Article type." Enums(blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,press_release,research_paper,site,technical_documentation,whitepaper)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param full_content query bool false "Include content when available." default(false)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} ArticleCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID getLatestArticles
// @Router /articles/latest [get]
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

	// if params.From.IsZero() {
	// 	filters.CreatedFrom = time.Now().AddDate(0, 0, -DEFAULT_WINDOW)
	// } else {
	// 	filters.CreatedFrom = params.From
	// }
	// if !params.To.IsZero() {
	// 	filters.CreatedTo = params.To
	// }

	page_out, err := r.DB.QueryLatestBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_WITHOUT_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// getTrendingArticles godoc
// @Summary List trending Articles
// @Description Returns attention-ranked Articles with trend metrics when available. Date bounds are not accepted.
// @Tags Articles
// @Produce json
// @Param q query string false "Optional relevance query. Requires score_threshold greater than zero." maxlength(512)
// @Param score_threshold query number false "Required when q is supplied." minimum(0) maximum(1)
// @Param content_type query string false "Stored Article type." Enums(blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,post,press_release,research_paper,site,technical_documentation,whitepaper)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param full_content query bool false "Include content when available." default(false)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} ArticleCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID getTrendingArticles
// @Router /articles/trending [get]
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

	// if params.From.IsZero() {
	// 	filters.UpdatedFrom = time.Now().AddDate(0, 0, -DEFAULT_WINDOW)
	// } else {
	// 	filters.UpdatedFrom = params.From
	// }
	// if !params.To.IsZero() {
	// 	filters.UpdatedTo = params.To
	// }

	page_out, err := r.DB.QueryTrendingBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_WITH_TREND)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryTrendingBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// getTopHeadlinesArticles is the B04 GET /articles/top-headlines target scaffold.
// Primary difference between this and getTrendingArticles and getTopHeadlines is the window of time and content type.
// getTopHeadlines is always fixed within the last 24 hours and `news` content type. The returned news are created and trending in the last 24 hours.
// The result excludes content unless explicitly requested.
// getTopHeadlines godoc
// @Summary List top headlines
// @Description Returns news Articles from the recent 24-hour window, ordered by attention. content_type and date bounds are not accepted.
// @Tags Articles
// @Produce json
// @Param q query string false "Optional relevance query. Requires score_threshold greater than zero." maxlength(512)
// @Param score_threshold query number false "Required when q is supplied." minimum(0) maximum(1)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param full_content query bool false "Include content when available." default(false)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} ArticleCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID getTopHeadlines
// @Router /top-headlines [get]
func (r *Configuration) getTopHeadlines(c *gin.Context) {
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

	filters.CreatedFrom = time.Now().AddDate(0, 0, -MIN_WINDOW)
	filters.ObservedFrom = time.Now().AddDate(0, 0, -MIN_WINDOW)
	filters.Kind = "news"

	page_out, err := r.DB.QueryTrendingBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_HEADLINES)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryTrendingBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// getArticle godoc
// @Summary Get an Article
// @Description Returns one Article selected by UUID. Set full_content=true to request content when available.
// @Tags Articles
// @Produce json
// @Param id path string true "Article UUID." format(uuid)
// @Param full_content query bool false "Include content when available." default(false)
// @Success 200 {object} ArticleDetailResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @Failure 404 {object} ErrorResponse "Article not found"
// @ID getArticle
// @Router /articles/{id} [get]
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
	writeDetail(c, *toArticleDetail(&bean))
}

// getSimilarArticles godoc
// @Summary List related Articles
// @Description Returns related Articles for an Article UUID, ordered newest first. It is not a relevance-ranked search route.
// @Tags Articles
// @Produce json
// @Param id path string true "Article UUID." format(uuid)
// @Param content_type query string false "Stored Article type." Enums(blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,post,press_release,research_paper,site,technical_documentation,whitepaper)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param from query string false "UTC lower timestamp bound." format(date)
// @Param to query string false "UTC upper timestamp bound." format(date)
// @Param full_content query bool false "Include content when available." default(false)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} ArticleCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @Failure 404 {object} ErrorResponse "Article not found"
// @ID getSimilarArticles
// @Router /articles/{id}/similar [get]
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

// getArticleMentions godoc
// @Summary List Article mentions
// @Description Returns external platform or forum observations for an Article UUID, ordered by observation time.
// @Tags Articles
// @Produce json
// @Param id path string true "Article UUID." format(uuid)
// @Param platforms query []string false "Mention platforms (CSV)." collectionFormat(csv)
// @Param forums query []string false "Mention forums (CSV)." collectionFormat(csv)
// @Param from query string false "UTC lower observation timestamp." format(date)
// @Param to query string false "UTC upper observation timestamp." format(date)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} MentionCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @Failure 404 {object} ErrorResponse "Article not found"
// @ID getArticleMentions
// @Router /articles/{id}/mentions [get]
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
	filters, err := params.createFilters(c, r)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QueryMentions(c.Request.Context(), params.ID, *filters, *page_req)
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

// getSources godoc
// @Summary List Sources
// @Description Returns publisher Sources. q matches the beginning of Source metadata; domains narrows results.
// @Tags Sources
// @Produce json
// @Param q query string false "Optional Source prefix query." maxlength(512)
// @Param ids query []string false "Source ids (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains (CSV)." collectionFormat(csv)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} SourceCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID listSources
// @Router /sources [get]
func (r *Configuration) getSources(c *gin.Context) {
	var params sourceSearchParams
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
	page_out, err := r.DB.QuerySources(c.Request.Context(), *filters, *page_req, db.SOURCE_COLUMNS_BASE)
	if err != nil {
		shared.LogError(err, "[ERROR] QuerySources")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toSourceDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// getSource godoc
// @Summary Get a Source
// @Description Returns one publisher Source selected by UUID.
// @Tags Sources
// @Produce json
// @Param id path string true "Source UUID." format(uuid)
// @Success 200 {object} SourceDetailResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @Failure 404 {object} ErrorResponse "Source not found"
// @ID getSource
// @Router /sources/{id} [get]
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
// @Summary List entity labels
// @Description Lists values accepted by the corresponding Article filter.
// @Tags Discovery
// @Produce json
// @Param q query string false "Optional case-insensitive prefix query." maxlength(512)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} TagCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID listEntities
// @Router /entities [get]
func (r *Configuration) getEntities(c *gin.Context) {
	getTags(r, c, "entities", "entity")
}

// getRegions godoc
// @Summary List region labels
// @Description Lists values accepted by the corresponding Article filter.
// @Tags Discovery
// @Produce json
// @Param q query string false "Optional case-insensitive prefix query." maxlength(512)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} TagCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID listRegions
// @Router /regions [get]
func (r *Configuration) getRegions(c *gin.Context) {
	getTags(r, c, "regions", "region")
}

// getCategories godoc
// @Summary List category labels
// @Description Lists values accepted by the corresponding Article filter.
// @Tags Discovery
// @Produce json
// @Param q query string false "Optional case-insensitive prefix query." maxlength(512)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} TagCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID listCategories
// @Router /categories [get]
func (r *Configuration) getCategories(c *gin.Context) {
	getTags(r, c, "categories", "category")
}

// getSentiments godoc
// @Summary List sentiment labels
// @Description Lists values accepted by the corresponding Article filter.
// @Tags Discovery
// @Produce json
// @Param q query string false "Optional case-insensitive prefix query." maxlength(512)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} TagCollectionResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID listSentiments
// @Router /sentiments [get]
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

// getStories godoc
// @Summary List Stories
// @Description Returns Stories identified by stable UUIDs. Use Story filters to narrow coverage; a Story collection does not include every member Article.
// @Tags Stories
// @Produce json
// @Param q query string false "Optional Story relevance query." maxlength(512)
// @Param score_threshold query number false "Optional Story relevance threshold; requires q." minimum(0) maximum(1)
// @Param content_type query string false "Stored Article type." Enums(blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,post,press_release,research_paper,site,technical_documentation,whitepaper)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param min_article_count query int false "Minimum Story Article count. Default 2." default(2) minimum(2)
// @Param from query string false "UTC lower publication timestamp." format(date)
// @Param to query string false "UTC upper publication timestamp." format(date)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} StoryCollectionResponse "Story collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @ID listStories
// @Router /stories [get]
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

	page_out, err := r.DB.QueryClusters(c.Request.Context(), *filters, *page_req)
	if err != nil {
		shared.LogError(err, "[ERROR] QueryClusters")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toStoryDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// func (r *Configuration) dispatchStory(c *gin.Context) {
// 	raw := strings.TrimSpace(strings.TrimPrefix(c.Param("story_id"), "/"))
// 	if raw == "" {
// 		writeError(c, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "story_id is required"))
// 		return
// 	}
// 	if strings.HasSuffix(raw, "/articles") {
// 		candidate := strings.TrimSuffix(raw, "/articles")
// 		exists, err := r.DB.StoryExists(c.Request.Context(), raw)
// 		if err != nil {
// 			shared.LogError(err, "[ERROR] StoryExists")
// 			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
// 			return
// 		}
// 		if !exists && candidate != "" {
// 			c.Params = gin.Params{{Key: "story_id", Value: candidate}}
// 			r.getStoryArticles(c)
// 			return
// 		}
// 	}
// 	c.Params = gin.Params{{Key: "story_id", Value: raw}}
// 	r.getStory(c)
// }

// getStory godoc
// @Summary Get a Story
// @Description Returns one Story selected by its stable UUID, including a link to its paginated member Articles.
// @Tags Stories
// @Produce json
// @Param id path string true "Story UUID." format(uuid)
// @Success 200 {object} StoryDetailResponse "Story detail envelope"
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @Failure 404 {object} ErrorResponse "Story not found"
// @ID getStory
// @Router /stories/{id} [get]
func (r *Configuration) getStory(c *gin.Context) {
	var params itemIDParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	story, err := r.DB.GetCluster(c.Request.Context(), params.ID)
	if err != nil {
		shared.LogError(err, "[ERROR] GetCluster")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_STORY_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeDetail(c, toStoryDetail(&story))
}

// getStoryArticles godoc
// @Summary List Story Articles
// @Description Returns the member Articles for one stable Story UUID. Article filters narrow the returned members.
// @Tags Stories
// @Produce json
// @Param id path string true "Story UUID." format(uuid)
// @Param content_type query string false "Stored Article type." Enums(blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,post,press_release,research_paper,site,technical_documentation,whitepaper)
// @Param sources query []string false "Source UUIDs to include (CSV)." collectionFormat(csv)
// @Param exclude_sources query []string false "Source UUIDs to exclude (CSV)." collectionFormat(csv)
// @Param domains query []string false "Source domains to include (CSV)." collectionFormat(csv)
// @Param exclude_domains query []string false "Source domains to exclude (CSV)." collectionFormat(csv)
// @Param authors query []string false "Author text filters (CSV)." collectionFormat(csv)
// @Param categories query []string false "Category values (CSV)." collectionFormat(csv)
// @Param exclude_categories query []string false "Excluded category values (CSV)." collectionFormat(csv)
// @Param regions query []string false "Region values (CSV)." collectionFormat(csv)
// @Param entities query []string false "Entity values (CSV)." collectionFormat(csv)
// @Param sentiments query []string false "Sentiment values (CSV)." collectionFormat(csv)
// @Param tags query []string false "Normalized tag terms (CSV)." collectionFormat(csv)
// @Param from query string false "UTC lower publication timestamp." format(date)
// @Param to query string false "UTC upper publication timestamp." format(date)
// @Param full_content query bool false "Include content when available." default(false)
// @Param limit query int false "Maximum records per page. Default 20, max 100." default(20) minimum(1) maximum(100)
// @Param cursor query string false "Opaque continuation token from pagination.next_cursor. Send it unchanged."
// @Success 200 {object} StoryArticleCollectionResponse "Story Article collection envelope"
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 401 {object} ErrorResponse "Missing or invalid API key"
// @Failure 429 {object} ErrorResponse "Request limit reached"
// @Failure 500 {object} ErrorResponse "Service unavailable"
// @Failure 404 {object} ErrorResponse "Story not found"
// @ID listStoryArticles
// @Router /stories/{id}/articles [get]
func (r *Configuration) getStoryArticles(c *gin.Context) {
	var params storyArticleParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	exists, err := r.DB.ClusterExists(c.Request.Context(), params.ID)
	if err != nil {
		shared.LogError(err, "[ERROR] ClusterExists")
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
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeStoryArticles(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor, params.ID)
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
	protected.GET("/articles/:id", config.getArticle)
	protected.GET("/articles/:id/similar", config.getSimilarArticles)
	protected.GET("/articles/:id/mentions", config.getArticleMentions)

	// HEADLINES routes
	protected.GET("/top-headlines", config.getTopHeadlines)

	// STORIES routes. Wildcard captures URL-like story IDs (slashes); trailing /articles is membership.
	protected.GET("/stories", config.getStories)
	protected.GET("/stories/:id", config.getStory)
	protected.GET("/stories/:id/articles", config.getStoryArticles)

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
