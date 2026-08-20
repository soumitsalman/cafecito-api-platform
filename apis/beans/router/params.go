package router

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
)

const (
	DEFAULT_CURSOR_LIMIT = 20
	MAX_CURSOR_LIMIT     = 100
	MAX_FILTER_VALUES    = 128
)

// paginationParams is embedded by collection requests for B01, B03-B07, B12,
// and B14-B18. B02, B10, and B13 are detail routes without collection pagination.
type paginationParams struct {
	Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
	Cursor string `form:"cursor"`
}

type itemIDParams struct {
	ID uuid.UUID `uri:"id,parser=encoding.TextUnmarshaler" binding:"required" swaggertype:"string" format:"uuid"`
}

// articleFilterParams contains the non-query, non-identity Article filters
// shared by B01, B03, B04, B05, and B06 target requests.
type articleFilterParams struct {
	ContentType       string      `form:"content_type" binding:"omitempty,oneof=blog contract earnings_report enforcement_action financial_report lawsuit news official_statement podcast post press_release research_paper site technical_documentation whitepaper"`
	Sources           []uuid.UUID `form:"sources,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	ExcludeSources    []uuid.UUID `form:"exclude_sources,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	Domains           []string    `form:"domains" collection_format:"csv" binding:"max=100"`
	ExcludeDomains    []string    `form:"exclude_domains" collection_format:"csv" binding:"max=100"`
	Authors           []string    `form:"authors" collection_format:"csv" binding:"max=100"`
	Tags              []string    `form:"tags" collection_format:"csv" binding:"max=100"`
	Categories        []string    `form:"categories" collection_format:"csv" binding:"max=100"`
	ExcludeCategories []string    `form:"exclude_categories" collection_format:"csv" binding:"max=100"`
	Sentiments        []string    `form:"sentiments" collection_format:"csv" binding:"max=100"`
	Entities          []string    `form:"entities" collection_format:"csv" binding:"max=100"`
	Regions           []string    `form:"regions" collection_format:"csv" binding:"max=100"`
	FullContent       bool        `form:"full_content,default=false"`
}

type vectorSearchParams struct {
	Q              string  `form:"q" binding:"max=512"`
	ScoreThreshold float64 `form:"score_threshold" binding:"min=0,max=1"`
}

// articleFeedParams is shared by B03 GET /articles/latest, B04 GET /articles/top-headlines, and B05 GET /articles/trending target requests.
// getLatestArticles, getTopHeadlines, and getTrendingArticles in routes.go bind this type.
type articleFeedParams struct {
	articleFilterParams
	vectorSearchParams
	paginationParams
}

func (params *articleFeedParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	if len(params.Q) > 0 && params.ScoreThreshold == 0 {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "`score_threshold` is required when `q` is present")
	}
	return nil
}

// articleSearchParams is the B01 GET /articles/search target request.
// searchArticles in routes.go binds this target request.
type articleSearchParams struct {
	articleFeedParams
	IDs  []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=100"`
	URLs []string    `form:"urls" collection_format:"csv" binding:"max=100"`
	From time.Time   `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To   time.Time   `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	paginationParams
}

func (params *articleSearchParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// articleDetailParams is the B02 GET /articles/{id} target request. The
// current scaffold handler is registered in NewRouter.
type articleDetailParams struct {
	itemIDParams
	FullContent bool `form:"full_content,default=false"`
}

func (params *articleDetailParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.itemIDParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// similarArticlesParams is the B06 GET /articles/{id}/similar target request.
// It excludes q, score_threshold, ids, and urls.
type similarArticlesParams struct {
	itemIDParams
	articleFilterParams
	From time.Time `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To   time.Time `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	paginationParams
}

func (params *similarArticlesParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.itemIDParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// articleMentionsParams is the B07 GET /articles/{id}/mentions scaffold
// request registered in NewRouter.
type articleMentionsParams struct {
	itemIDParams
	Platforms []string  `form:"platforms" collection_format:"csv" binding:"max=100"`
	Forums    []string  `form:"forums" collection_format:"csv" binding:"max=100"`
	From      time.Time `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To        time.Time `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	paginationParams
}

func (params *articleMentionsParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.itemIDParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// sourceListParams is the B12 GET /sources target request. The current route
// getSources in routes.go binds this target request.
type sourceListParams struct {
	Q       string      `form:"q" binding:"max=512"`
	IDs     []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=100"`
	Domains []string    `form:"domains" collection_format:"csv" binding:"max=100"`
	paginationParams
}

func (params *sourceListParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// sourceDetailParams is the B13 GET /sources/{id} target request. The current route
// getSource in routes.go binds this target request.
type sourceDetailParams struct {
	itemIDParams
}

func (params *sourceDetailParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.itemIDParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// tagListParams is used by the B14-B18 GET /categories, /entities, /regions, /sentiments, and /tags requests.
// The current route getCategories, getEntities, getRegions, getSentiments, and getTags in routes.go bind this target request.
type tagListParams struct {
	Q string `form:"q" binding:"max=512"`
	paginationParams
}

func (params *tagListParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// storyPathParams binds the B10/B11 path story_id parameter.
// Story IDs are opaque strings (currently a derived cluster key; later a URL).
type storyPathParams struct {
	StoryID string `uri:"story_id" binding:"required"`
}

func (params *storyPathParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	params.StoryID = strings.TrimSpace(params.StoryID)
	if params.StoryID == "" {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "story_id is required")
	}
	return nil
}

// storySearchParams is the B09 GET /stories target request.
type storySearchParams struct {
	articleFilterParams
	vectorSearchParams
	MinArticleCount int       `form:"min_article_count,default=2" binding:"min=2"`
	From            time.Time `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To              time.Time `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	paginationParams
}

func (params *storySearchParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	if _, present := c.GetQuery("score_threshold"); present && strings.TrimSpace(params.Q) == "" {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "score_threshold requires q")
	}
	return nil
}

// storyArticleParams is the B11 GET /stories/{story_id}/articles target request.
type storyArticleParams struct {
	storyPathParams
	articleFilterParams
	From time.Time `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To   time.Time `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	paginationParams
}

func (params *storyArticleParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.storyPathParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

// articleCountParams is the B19 GET /articles/count target request.
type articleCountParams struct {
	articleFilterParams
	IDs     []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	Urls    []string    `form:"urls" collection_format:"csv" binding:"max=128"`
	From    time.Time   `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date" binding:"required"`
	To      time.Time   `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date" binding:"required"`
	GroupBy string      `form:"group_by" binding:"omitempty,oneof=published_day content_type domain category region sentiment"`
}

func (params *articleCountParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}
