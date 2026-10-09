package router

import (
	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	datautils "github.com/soumitsalman/data-utils"
)

// These routes are for internal use ONLY.
// They are not part of Public API.
// Exclude them from swagger docs and `beans.oas.json` definitions

const (
	PRIVATE_SORT_RECENT   = "recent"
	PRIVATE_SORT_TREND    = "trend"
	PRIVATE_SORT_RELEVANT = "relevant"
)

// StoryArticlePreviewDocument is a compact Article preview for Story top_articles.
type privateArticlePreviewDocument struct {
	ArticlePreviewDocument
	Trend *Trend `json:"trend,omitempty"`
}

func toPrivateArticlePreview(bean *db.Bean) *privateArticlePreviewDocument {
	return &privateArticlePreviewDocument{
		ArticlePreviewDocument: toArticlePreview(bean),
		Trend:                  nullArticleTrendPtr(bean),
	}
}

// // dbSort maps the request sort value to the db cursor sort constant.
// func (params *storySearchParams) dbSort() string {
// 	switch strings.ToLower(strings.TrimSpace(params.Sort)) {
// 	case PRIVATE_SORT_TREND:
// 		return db.SORT_TRENDING
// 	case PRIVATE_SORT_RELEVANT:
// 		return db.SORT_RELEVANT
// 	default:
// 		return db.SORT_RECENT
// 	}
// }
// func (params *privateStorySearchParams) createFilters(c *gin.Context, r *Configuration) (*db.ClusterFilters, error) {
// 	filters, err := params.storySearchParams.createFilters(c, r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if params.Q != "" && params.ScoreThreshold <= 0 {
// 		return nil, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, API_ERROR_MSG_SCORE_THRESHOLD_REQUIRED)
// 	}
// 	// sort=trend scopes observed attention metrics to the same from/to window.
// 	if params.dbSort() == db.SORT_TRENDING {
// 		filters.ObservedFrom = params.From
// 		filters.ObservedTo = utils.NormalizeEndOfDay(params.To)
// 	}
// 	return filters, nil
// }

// func (params *privateStorySearchParams) createPageRequest(c *gin.Context, r *Configuration) (*db.PageRequest, error) {
// 	page_req, err := params.paginationParams.createPageRequest(c, r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Reject cursors from a different sort mode to keep the keyset scan stable.
// 	if page_req.Cursor != nil && page_req.Cursor.Sort != params.dbSort() {
// 		return nil, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "cursor sort does not match requested sort")
// 	}
// 	return page_req, nil
// }

// handler for GET /preview/stories
// query params are used to filter the stories
func (r *Configuration) previewGetStories(c *gin.Context) {
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

	page_out, err := r.DB.QueryClusterPreviews(c.Request.Context(), *filters, *page_req)
	if err != nil {
		writeError(c, err)
		return
	}
	writeCollection(c, toStoryDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// handler for GET /preview/stories/:id
// query params are used to filter the top_articles and the representative article contents
func (r *Configuration) previewGetStory(c *gin.Context) {
	var params similarArticlesParams
	filters, _, err := extractBeanFiltersAndPage(r, c, &params)
	if err != nil {
		writeError(c, err)
		return
	}

	story, err := r.DB.GetCluster(c.Request.Context(), params.ID, *filters)
	if err != nil {
		writeError(c, err)
		return
	}
	writeDetail(c, toStoryDetail(&story))
}

// handler for GET /preview/stories/:id/articles
// query params are used to filter the articles
// equivalent to getStoryArticles, but with less columns.
// TODO: for future allow mulitple story ids
func (r *Configuration) previewGetStoryArticles(c *gin.Context) {
	var params similarArticlesParams
	filters, page_req, err := extractBeanFiltersAndPage(r, c, &params)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QueryClusterMembers(c.Request.Context(), params.ID, *filters, *page_req, db.BEAN_COLUMNS_PREVIEW)
	if err != nil {
		writeError(c, err)
		return
	}
	writeStoryArticles(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor, params.ID)
}

// handler for GET /preview/articles/:id/similar
// equivalent to getSimilarArticles, but includes trend columns.
func (r *Configuration) previewGetSimilarArticles(c *gin.Context) {
	var params similarArticlesParams
	filters, page_req, err := extractBeanFiltersAndPage(r, c, &params)
	if err != nil {
		writeError(c, err)
		return
	}

	page_out, err := r.DB.QuerySimilarBeans(c.Request.Context(), params.ID, *filters, *page_req, db.BEAN_COLUMNS_PREVIEW)
	if err != nil {
		writeError(c, err)
		return
	}
	previews := datautils.Transform(page_out.Items, func(item *db.Bean) privateArticlePreviewDocument {
		return *toPrivateArticlePreview(item)
	})
	writeCollection(c, previews, page_req.Limit, page_out.NextCursor)
}
