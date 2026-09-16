package router

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
)

// These routes are for internal use ONLY.
// They are not part of Public API.
// Exclude them from swagger docs and `beans.oas.json` definitions

const (
	PRIVATE_SORT_RECENT   = "recent"
	PRIVATE_SORT_TREND    = "trend"
	PRIVATE_SORT_RELEVANT = "relevant"
)

// privateUniqueArticleParams is the GET /private/articles/unique target request.
// It extends articleSearchParams with a sort selector and reuses from/to for
// created bounds (all sorts) and observed bounds (sort=trend).
type privateUniqueArticleParams struct {
	articleSearchParams
	Sort string `form:"sort,default=recent" binding:"omitempty,oneof=recent trend relevant"`
}

// StoryArticlePreviewDocument is a compact Article preview for Story top_articles.
type privateStoryArticlePreviewDocument struct {
	StoryArticlePreviewDocument
	Trend *Trend `json:"trend,omitempty"`
}

func toPrivateStoryArticlePreview(bean *db.Bean) *privateStoryArticlePreviewDocument {
	return &privateStoryArticlePreviewDocument{
		StoryArticlePreviewDocument: toStoryArticlePreview(bean),
		Trend:                       nullArticleTrendPtr(bean),
	}
}

func (params *privateUniqueArticleParams) shouldBind(c *gin.Context) error {
	if err := bindQuery(c, params); err != nil {
		return err
	}
	if err := requireScoreThresholdNeedsQ(c, params.Q); err != nil {
		return err
	}
	if params.dbSort() == db.SORT_RELEVANT && strings.TrimSpace(params.Q) == "" {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "sort=relevant requires q")
	}
	return nil
}

// dbSort maps the request sort value to the db cursor sort constant.
func (params *privateUniqueArticleParams) dbSort() string {
	switch strings.ToLower(strings.TrimSpace(params.Sort)) {
	case PRIVATE_SORT_TREND:
		return db.SORT_TRENDING
	case PRIVATE_SORT_RELEVANT:
		return db.SORT_RELEVANT
	default:
		return db.SORT_RECENT
	}
}

func (params *privateUniqueArticleParams) createFilters(c *gin.Context, r *Configuration) (*db.BeanFilters, error) {
	filters, err := params.articleSearchParams.createFilters(c, r)
	if err != nil {
		return nil, err
	}
	// sort=trend scopes observed attention metrics to the same from/to window.
	if params.dbSort() == db.SORT_TRENDING {
		filters.ObservedFrom = params.From
		filters.ObservedTo = utils.NormalizeEndOfDay(params.To)
	}
	return filters, nil
}

func (params *privateUniqueArticleParams) createPageRequest(c *gin.Context, r *Configuration) (*db.PageRequest, error) {
	page_req, err := params.paginationParams.createPageRequest(c, r)
	if err != nil {
		return nil, err
	}
	// Reject cursors from a different sort mode to keep the keyset scan stable.
	if page_req.Cursor != nil && page_req.Cursor.Sort != params.dbSort() {
		return nil, utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, "cursor sort does not match requested sort")
	}
	return page_req, nil
}

func (r *Configuration) privateGetUniqueArticles(c *gin.Context) {
	var params privateUniqueArticleParams
	filters, page_req, err := extractBeanFiltersAndPage(r, c, &params)
	if err != nil {
		writeError(c, err)
		return
	}
	page_out, err := r.DB.QueryUniqueBeans(c.Request.Context(), *filters, *page_req, params.dbSort(), db.BEAN_COLUMNS_WITH_TREND)
	if err != nil {
		utils.LogError(err, "[ERROR] QueryUniqueBeans")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	writeCollection(c, toArticleDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// privateStoryParams is the GET /private/stories/:id target request.
// The optional languages filter only influences which member article supplies the
// story title/summary; stats, categories, entities, and timestamps are unaffected.
type privateStoryParams struct {
	itemIDParams
	Languages []string `form:"languages" collection_format:"csv" binding:"max=100"`
}

func (params *privateStoryParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(&params.itemIDParams); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return bindQuery(c, params)
}

func (r *Configuration) privateGetStory(c *gin.Context) {
	var params privateStoryParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	story, err := r.DB.GetClusterInLanguages(c.Request.Context(), params.ID, utils.NormalizeTexts(params.Languages))
	if err != nil {
		utils.LogError(err, "[ERROR] GetClusterInLanguages")
		if errors.Is(err, db.ErrNonExistentID) {
			writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_STORY_NOT_FOUND))
		} else {
			writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		}
		return
	}
	writeDetail(c, toStoryDetail(&story))
}

// privateGetStoryArticles is the same shape as getStoryArticles: storyArticleParams,
// extractBeanFiltersAndPage, ClusterExists, QueryBeans. It silently drops limit/cursor
// so the full cluster membership is returned, and selects CLUSTER_BEAN_COLUMNS_MINIMAL.
func (r *Configuration) privateGetStoryArticles(c *gin.Context) {
	var params storyArticleParams
	filters, page_req, err := extractBeanFiltersAndPage(r, c, &params)
	if err != nil {
		writeError(c, err)
		return
	}

	exists, err := r.DB.ClusterExists(c.Request.Context(), params.ID)
	if err != nil {
		utils.LogError(err, "[ERROR] ClusterExists")
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	if !exists {
		writeError(c, utils.NewAPIError(utils.API_ERROR_NOT_FOUND, API_ERROR_MSG_STORY_NOT_FOUND))
		return
	}

	page_out, err := r.DB.QueryBeans(c.Request.Context(), *filters, *page_req, db.BEAN_COLUMNS_MINIMAL_WITH_TREND)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	previews := make([]privateStoryArticlePreviewDocument, 0, len(page_out.Items))
	for i := range page_out.Items {
		previews = append(previews, *toPrivateStoryArticlePreview(&page_out.Items[i]))
	}
	writeCollection(c, previews, len(previews), nil)
}
