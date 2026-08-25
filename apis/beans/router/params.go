package router

import (
	"fmt"
	"reflect"
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

func (params *itemIDParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

func bindQuery(c *gin.Context, params any) error {
	if err := rejectUnknownQuery(c, params); err != nil {
		return err
	}
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

func rejectUnknownQuery(c *gin.Context, params any) error {
	allowed := formQueryNames(params)
	for key := range c.Request.URL.Query() {
		if _, ok := allowed[key]; !ok {
			return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, fmt.Sprintf("Unknown or unsupported query parameter: %s", key))
		}
	}
	return nil
}

func formQueryNames(params any) map[string]struct{} {
	names := map[string]struct{}{}
	collectFormNames(reflect.TypeOf(params), names)
	return names
}

func collectFormNames(t reflect.Type, names map[string]struct{}) {
	if t == nil {
		return
	}
	if t.Kind() == reflect.Pointer {
		collectFormNames(t.Elem(), names)
		return
	}
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			collectFormNames(field.Type, names)
			continue
		}
		tag := field.Tag.Get("form")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name != "" {
			names[name] = struct{}{}
		}
	}
}

func requireScoreThresholdNeedsQ(c *gin.Context, q string) error {
	if _, present := c.GetQuery("score_threshold"); present && strings.TrimSpace(q) == "" {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "score_threshold requires q")
	}
	return nil
}

// articleScopeParams is the Article filter set shared by feeds and search, excluding content_type.
type articleScopeParams struct {
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

// articleFilterParams contains the non-query, non-identity Article filters
// shared by B01, B03, B05, and B06. content_type is request-filterable except post,
// which is response-only.
type articleFilterParams struct {
	ContentType string `form:"content_type" binding:"omitempty,oneof=blog contract earnings_report enforcement_action financial_report lawsuit news official_statement podcast press_release research_paper site technical_documentation whitepaper"`
	articleScopeParams
}

type vectorSearchParams struct {
	Q              string  `form:"q" binding:"max=512"`
	ScoreThreshold float64 `form:"score_threshold" binding:"min=0,max=1"`
}

// articleFeedParams is shared by GET /articles/latest and GET /articles/trending.
// Feed routes reject ids, urls, from, and to.
type articleFeedParams struct {
	articleFilterParams
	vectorSearchParams
	paginationParams
}

func (params *articleFeedParams) shouldBind(c *gin.Context) error {
	if err := bindQuery(c, params); err != nil {
		return err
	}
	return requireScoreThresholdNeedsQ(c, params.Q)
}

// topHeadlinesParams is GET /top-headlines. It rejects ids, urls, from, to, and content_type.
type topHeadlinesParams struct {
	articleScopeParams
	vectorSearchParams
	paginationParams
}

func (params *topHeadlinesParams) shouldBind(c *gin.Context) error {
	if err := bindQuery(c, params); err != nil {
		return err
	}
	return requireScoreThresholdNeedsQ(c, params.Q)
}

// articleSearchParams is the B01 GET /articles/search target request.
type articleSearchParams struct {
	articleFeedParams
	IDs  []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=100"`
	URLs []string    `form:"urls" collection_format:"csv" binding:"max=100"`
	From time.Time   `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To   time.Time   `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
}

func (params *articleSearchParams) shouldBind(c *gin.Context) error {
	if err := bindQuery(c, params); err != nil {
		return err
	}
	return requireScoreThresholdNeedsQ(c, params.Q)
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
	return bindQuery(c, params)
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
	return bindQuery(c, params)
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
	return bindQuery(c, params)
}

// sourceSearchParams is the B12 GET /sources target request. The current route
// getSources in routes.go binds this target request.
type sourceSearchParams struct {
	Q       string      `form:"q" binding:"max=512"`
	IDs     []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=100"`
	Domains []string    `form:"domains" collection_format:"csv" binding:"max=100"`
	paginationParams
}

func (params *sourceSearchParams) shouldBind(c *gin.Context) error {
	return bindQuery(c, params)
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
	return bindQuery(c, params)
}

// tagListParams is used by the B14-B18 GET /categories, /entities, /regions, /sentiments, and /tags requests.
// The current route getCategories, getEntities, getRegions, getSentiments, and getTags in routes.go bind this target request.
type tagListParams struct {
	Q string `form:"q" binding:"max=512"`
	paginationParams
}

func (params *tagListParams) shouldBind(c *gin.Context) error {
	return bindQuery(c, params)
}

// // storyPathParams binds the B10/B11 path story_id parameter.
// // Story IDs are opaque strings (currently a derived cluster key; later a URL).
// type storyPathParams struct {
// 	StoryID string `uri:"story_id" binding:"required"`
// }

// func (params *storyPathParams) shouldBind(c *gin.Context) error {
// 	if err := c.ShouldBindUri(params); err != nil {
// 		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
// 	}
// 	params.StoryID = strings.TrimSpace(params.StoryID)
// 	if params.StoryID == "" {
// 		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "story_id is required")
// 	}
// 	return nil
// }

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
	if err := bindQuery(c, params); err != nil {
		return err
	}
	return requireScoreThresholdNeedsQ(c, params.Q)
}

// storyArticleParams is the B11 GET /stories/{story_id}/articles target request.
type storyArticleParams struct {
	itemIDParams
	articleFilterParams
	From time.Time `form:"from" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	To   time.Time `form:"to" time_format:"2006-01-02" time_utc:"true" swaggertype:"string" format:"date"`
	paginationParams
}

func (params *storyArticleParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.itemIDParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return bindQuery(c, params)
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
	return bindQuery(c, params)
}
