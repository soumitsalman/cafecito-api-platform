package router

import (
	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	datautils "github.com/soumitsalman/data-utils"
)

// These routes are for internal use ONLY.
// They are not part of Public API.
// Exclude them from swagger docs and `beans.oas.json` definitions

type privateStoryDocument struct {
	StoryDocument
	Trend *Trend `json:"trend,omitempty"`
}

func toPrivateStoryDocument(story *db.Cluster) privateStoryDocument {
	return privateStoryDocument{
		StoryDocument: toStoryDocument(story),
		Trend:         nullTrendPtr(&story.TrendProperties),
	}
}

// handler for GET /preview/stories
// query params are used to filter the stories
// similar to getStories, but with more trend properties fields
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
	story_docs := datautils.Transform(page_out.Items, func(story *db.Cluster) privateStoryDocument {
		return toPrivateStoryDocument(story)
	})
	writeCollection(c, story_docs, page_req.Limit, page_out.NextCursor)
}

// handler for GET /preview/stories/:id
// query params are used to filter the top_articles and the representative article contents
func (r *Configuration) previewGetStory(c *gin.Context) {
	var params linkedArticlesParams
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
	writeDetail(c, toPrivateStoryDocument(&story))
}

// handler for GET /preview/stories/:id/articles
// query params are used to filter the articles
// equivalent to getStoryArticles, but with less columns.
func (r *Configuration) previewGetStoryArticles(c *gin.Context) {
	var params linkedArticlesParams
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
	writeCollection(c, toArticlePreviewDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}

// handler for GET /preview/articles/:id/similar
// equivalent to getSimilarArticles, but with less columns.
func (r *Configuration) previewGetSimilarArticles(c *gin.Context) {
	var params linkedArticlesParams
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
	writeCollection(c, toArticlePreviewDocuments(page_out.Items), page_req.Limit, page_out.NextCursor)
}
