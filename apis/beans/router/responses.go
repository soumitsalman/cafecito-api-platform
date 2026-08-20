package router

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soumitsalman/cafecito-api-platform/apis/beans/db"
	datautils "github.com/soumitsalman/data-utils"
)

// CursorPagination is the target Beans collection pagination envelope.
type Pagination struct {
	Limit      int     `json:"limit" binding:"required"`
	NumResults int     `json:"num_results" binding:"required"`
	NextCursor *string `json:"next_cursor" binding:"required"`
}

// ResponseMeta contains freshness metadata for changing collections.
type ResponseMeta struct {
	AsOf time.Time `json:"as_of"`
}

// CollectionResponse is the canonical target collection envelope.
type CollectionResponse[T any] struct {
	Pagination Pagination   `json:"pagination" binding:"required"`
	Meta       ResponseMeta `json:"meta" binding:"required"`
	Data       []T          `json:"data" binding:"required"`
}

// DetailResponse is the canonical target detail envelope.
type DetailResponse[T any] struct {
	Data T `json:"data" binding:"required"`
}

// ErrorResponse is the canonical target error envelope.
type ErrorResponse struct {
	Error error `json:"error"`
}

// SourceDocument is the Source collection/detail payload.
type SourceDocument struct {
	ID          uuid.UUID `json:"id" swaggertype:"string" format:"uuid"`
	BaseURL     string    `json:"url"`
	DomainName  string    `json:"domain"`
	SiteName    *string   `json:"name"`
	Description *string   `json:"description,omitempty"`
	Favicon     *string   `json:"favicon_url,omitempty"`
	RSSFeed     *string   `json:"rss_feed_url,omitempty"`
}

// ArticleDocument is the normalized public Article payload.
type ArticleDocument struct {
	ID         uuid.UUID `json:"id" swaggertype:"string" format:"uuid"`
	URL        string    `json:"url"`
	Kind       string    `json:"content_type" enums:"blog,contract,earnings_report,enforcement_action,financial_report,lawsuit,news,official_statement,podcast,post,press_release,research_paper,site,technical_documentation,whitepaper"`
	Created    time.Time `json:"published_at" swaggertype:"string" format:"date-time"`
	Author     *string   `json:"author"`
	ImageURL   *string   `json:"image_url"`
	Title      *string   `json:"title"`
	Summary    *string   `json:"summary"`
	Content    *string   `json:"content,omitempty"`
	Categories []string  `json:"categories"`
	Regions    []string  `json:"regions"`
	Entities   []string  `json:"entities"`
	Sentiments []string  `json:"sentiments"`
	Tags       []string  `json:"tags"`
	// StoryID    *uuid.UUID     `json:"story_id" swaggertype:"string" format:"uuid"` // TODO: enable this later
	StoryID *string         `json:"story_id"`
	Source  *SourceDocument `json:"source"`
	Trend   *Trend          `json:"trend,omitempty"`
}

func toArticleDocument(bean *db.Bean) *ArticleDocument {
	doc := &ArticleDocument{
		ID:         bean.ID,
		URL:        bean.URL,
		Kind:       bean.Kind,
		Created:    bean.Created,
		Author:     nullStringPtr(bean.Author),
		ImageURL:   nullStringPtr(bean.ImageUrl),
		Categories: bean.Categories,
		Sentiments: bean.Sentiments,
		Regions:    bean.Regions,
		Entities:   bean.Entities,
		Tags:       concatArrays(bean.Categories, bean.Regions, bean.Entities),
		Title:      nullStringPtr(bean.Title),
		Summary:    nullStringPtr(bean.Summary),
		Content:    nullStringPtr(bean.Content),
		StoryID:    nullStringPtr(bean.ClusterID),
	}
	if bean.SourceID != uuid.Nil {
		doc.Source = &SourceDocument{
			ID:          bean.SourceID,
			BaseURL:     bean.BaseURL.String,
			DomainName:  bean.Source,
			SiteName:    nullStringPtr(bean.SiteName),
			Description: nullStringPtr(bean.Description),
			Favicon:     nullStringPtr(bean.Favicon),
			RSSFeed:     nullStringPtr(bean.RSSFeed),
		}
	}
	if bean.Likes.Valid || bean.Comments.Valid || bean.Shares.Valid || bean.Subscribers.Valid || bean.Related.Valid {
		doc.Trend = &Trend{
			Likes:       nullInt64Ptr(bean.Likes),
			Comments:    nullInt64Ptr(bean.Comments),
			Shares:      nullInt64Ptr(bean.Shares),
			Subscribers: nullInt64Ptr(bean.Subscribers),
			Related:     nullInt64Ptr(bean.Related),
		}
	}
	return doc
}

func toArticleDocuments(beans []db.Bean) []ArticleDocument {
	return datautils.Transform(beans, func(bean *db.Bean) ArticleDocument {
		return *toArticleDocument(bean)
	})
}

// Trend contains nullable attention metrics for headline and trending Articles.
type Trend struct {
	Likes       *int64 `json:"likes,omitempty"`
	Comments    *int64 `json:"comments,omitempty"`
	Shares      *int64 `json:"mentions,omitempty"`
	Subscribers *int64 `json:"audience,omitempty"`
	Related     *int64 `json:"related,omitempty"`
}

// MentionEngagement contains nullable observed engagement values.
type MentionEngagement struct {
	Likes    *int64 `json:"likes"`
	Comments *int64 `json:"comments"`
	Audience *int64 `json:"audience"`
}

// MentionDocument is the B07 external mention payload.
type MentionDocument struct {
	URL        string            `json:"url"`
	Platform   string            `json:"platform"`
	Forum      *string           `json:"forum"`
	ObservedAt time.Time         `json:"observed_at" swaggertype:"string" format:"date-time"`
	Engagement MentionEngagement `json:"engagement"`
}

// TagDocument is the normalized item returned by B14-B18.
type TagDocument struct {
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// ArticleCollectionResponse names the Article collection schema for Swagger.
type ArticleCollectionResponse struct {
	Pagination Pagination        `json:"pagination" binding:"required"`
	Meta       ResponseMeta      `json:"meta" binding:"required"`
	Data       []ArticleDocument `json:"data" binding:"required"`
}

// ArticleDetailResponse names the Article detail schema for Swagger.
type ArticleDetailResponse struct {
	Data ArticleDocument `json:"data" binding:"required"`
}

// SourceCollectionResponse names the Source collection schema for Swagger.
type SourceCollectionResponse struct {
	Pagination Pagination       `json:"pagination" binding:"required"`
	Meta       ResponseMeta     `json:"meta" binding:"required"`
	Data       []SourceDocument `json:"data" binding:"required"`
}

// SourceDetailResponse names the Source detail schema for Swagger.
type SourceDetailResponse struct {
	Data SourceDocument `json:"data" binding:"required"`
}

// TagCollectionResponse names the tag collection schema for Swagger.
type TagCollectionResponse struct {
	Pagination Pagination    `json:"pagination" binding:"required"`
	Meta       ResponseMeta  `json:"meta" binding:"required"`
	Data       []TagDocument `json:"data" binding:"required"`
}

// MentionCollectionResponse names the mention collection schema for Swagger.
type MentionCollectionResponse struct {
	Pagination Pagination        `json:"pagination" binding:"required"`
	Meta       ResponseMeta      `json:"meta" binding:"required"`
	Data       []MentionDocument `json:"data" binding:"required"`
}

func nullStringPtr(n sql.NullString) *string {
	if n.Valid {
		return &n.String
	}
	return nil
}

func nullInt64Ptr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func nullFloat64Ptr(n sql.NullFloat64) *float64 {
	if n.Valid {
		return &n.Float64
	}
	return nil
}

func concatArrays(arrays ...[]string) []string {
	var result []string
	for _, arr := range arrays {
		result = append(result, arr...)
	}
	return result
}

func toSourceDocument(source *db.Source) *SourceDocument {
	doc := &SourceDocument{
		ID:         source.ID,
		BaseURL:    source.BaseURL,
		DomainName: source.DomainName,
	}
	if source.SiteName.Valid {
		doc.SiteName = &source.SiteName.String
	}
	if source.Description.Valid {
		doc.Description = &source.Description.String
	}
	if source.Favicon.Valid {
		doc.Favicon = &source.Favicon.String
	}
	if source.RSSFeed.Valid {
		doc.RSSFeed = &source.RSSFeed.String
	}
	return doc
}

func toSourceDocuments(sources []db.Source) []SourceDocument {
	return datautils.Transform(sources, func(source *db.Source) SourceDocument {
		return *toSourceDocument(source)
	})

}

func toTagDocuments(tags []string, TAG_TYPE string) []TagDocument {
	return datautils.Transform(tags, func(tag *string) TagDocument {
		return TagDocument{Value: *tag, Type: TAG_TYPE}
	})
}

func emptyStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func storyArticlesPath(story_id string) string {
	return "/stories/" + story_id + "/articles"
}

func toCompactSource(bean *db.Bean) CompactSource {
	return CompactSource{
		ID:     bean.SourceID,
		Domain: bean.Source,
		Name:   nullStringPtr(bean.SiteName),
		URL:    bean.BaseURL.String,
	}
}

func toStoryArticlePreview(bean *db.Bean) StoryArticlePreview {
	title := ""
	if bean.Title.Valid {
		title = bean.Title.String
	}
	return StoryArticlePreview{
		ID:          bean.ID,
		URL:         bean.URL,
		Title:       title,
		PublishedAt: bean.Created,
		Source:      toCompactSource(bean),
	}
}

func toStoryItem(story *db.Story) StoryItem {
	previews := datautils.Transform(story.TopArticles, func(bean *db.Bean) StoryArticlePreview {
		return toStoryArticlePreview(bean)
	})
	if previews == nil {
		previews = []StoryArticlePreview{}
	}
	return StoryItem{
		ID:               story.ID,
		Title:            story.Title,
		FirstPublishedAt: story.FirstPublishedAt,
		LastPublishedAt:  story.LastPublishedAt,
		ArticleCount:     story.ArticleCount,
		SourceCount:      story.SourceCount,
		Categories:       emptyStrings(story.Categories),
		Regions:          emptyStrings(story.Regions),
		Entities:         emptyStrings(story.Entities),
		Tags:             emptyStrings(story.Tags),
		TopArticles:      previews,
	}
}

func toStoryItems(stories []db.Story) []StoryItem {
	return datautils.Transform(stories, func(story *db.Story) StoryItem {
		return toStoryItem(story)
	})
}

func toStoryDetailItem(story *db.Story) StoryDetailItem {
	return StoryDetailItem{
		StoryItem: toStoryItem(story),
		Links:     StoryLinks{Articles: storyArticlesPath(story.ID)},
	}
}

func toMentionDocument(mention *db.Mention) *MentionDocument {
	return &MentionDocument{
		URL:        mention.URL,
		Platform:   mention.Platform,
		Forum:      nullStringPtr(mention.Forum),
		ObservedAt: mention.ObservedAt,
		Engagement: MentionEngagement{
			Likes:    nullInt64Ptr(mention.Likes),
			Comments: nullInt64Ptr(mention.Comments),
			Audience: nullInt64Ptr(mention.Subscribers),
		},
	}
}

func toMentionDocuments(mentions []db.Mention) []MentionDocument {
	return datautils.Transform(mentions, func(mention *db.Mention) MentionDocument {
		return *toMentionDocument(mention)
	})
}

func toArticleDetailItem(bean *db.Bean, full_content bool) ArticleDetailItem {
	id := bean.ID.String()
	return ArticleDetailItem{
		ArticleDocument: *toArticleDocument(bean),
		Links: ArticleLinks{
			Similar:  "/articles/" + id + "/similar",
			Mentions: "/articles/" + id + "/mentions",
		},
		includeContent: full_content,
	}
}

// ArticleLinks contains sub-resource links for B02 Article detail.
type ArticleLinks struct {
	Similar  string `json:"similar,omitempty"`
	Mentions string `json:"mentions,omitempty"`
	Story    string `json:"story,omitempty"`
}

// ArticleDetailItem wraps an ArticleDocument with detail-only links.
type ArticleDetailItem struct {
	ArticleDocument
	Links          ArticleLinks `json:"links"`
	includeContent bool         `json:"-"`
}

func (item ArticleDetailItem) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(item.ArticleDocument)
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	if item.includeContent {
		payload["content"] = item.Content
	}
	payload["links"] = item.Links
	return json.Marshal(payload)
}

// ArticleDetailItemResponse names the Article detail schema for Swagger.
type ArticleDetailItemResponse struct {
	Data ArticleDetailItem `json:"data" binding:"required"`
}

// CompactSource is the nested source object on Story top_articles.
// Keys are always present; missing display metadata is null.
type CompactSource struct {
	ID     uuid.UUID `json:"id" swaggertype:"string" format:"uuid"`
	Domain string    `json:"domain"`
	Name   *string   `json:"name"`
	URL    string    `json:"url"`
}

// StoryArticlePreview is a compact Article preview for Story top_articles.
type StoryArticlePreview struct {
	ID          uuid.UUID     `json:"id" swaggertype:"string" format:"uuid"`
	URL         string        `json:"url"`
	Title       string        `json:"title"`
	PublishedAt time.Time     `json:"published_at" swaggertype:"string" format:"date-time"`
	Source      CompactSource `json:"source"`
}

// StoryItem is the canonical Story payload for B09 and B10.
type StoryItem struct {
	ID               string                `json:"id"`
	Title            string                `json:"title"`
	FirstPublishedAt time.Time             `json:"first_published_at" swaggertype:"string" format:"date-time"`
	LastPublishedAt  time.Time             `json:"last_published_at" swaggertype:"string" format:"date-time"`
	ArticleCount     int                   `json:"article_count"`
	SourceCount      int                   `json:"source_count"`
	Categories       []string              `json:"categories"`
	Regions          []string              `json:"regions"`
	Entities         []string              `json:"entities"`
	Tags             []string              `json:"tags"`
	TopArticles      []StoryArticlePreview `json:"top_articles"`
}

// StoryLinks contains sub-resource links for B10 Story detail.
type StoryLinks struct {
	Articles string `json:"articles"`
}

// StoryDetailItem wraps a StoryItem with detail-only links.
type StoryDetailItem struct {
	StoryItem
	Links StoryLinks `json:"links"`
}

// StoryDetailResponse names the Story detail schema for Swagger.
type StoryDetailResponse struct {
	Data StoryDetailItem `json:"data" binding:"required"`
	Meta ResponseMeta    `json:"meta" binding:"required"`
}

// StoryCollectionResponse names the Story collection schema for Swagger.
type StoryCollectionResponse struct {
	Pagination Pagination   `json:"pagination" binding:"required"`
	Meta       ResponseMeta `json:"meta" binding:"required"`
	Data       []StoryItem  `json:"data" binding:"required"`
}

// StoryArticleMeta contains freshness metadata for B11 Story article collections.
type StoryArticleMeta struct {
	StoryID string    `json:"story_id"`
	AsOf    time.Time `json:"as_of"`
}

// StoryArticleCollectionResponse names the B11 Story article collection schema for Swagger.
type StoryArticleCollectionResponse struct {
	Pagination Pagination        `json:"pagination" binding:"required"`
	Meta       StoryArticleMeta  `json:"meta" binding:"required"`
	Data       []ArticleDocument `json:"data" binding:"required"`
}

// CountData is the B19 count response without group_by.
type CountData struct {
	Count int64 `json:"count"`
}

// CountBucket is a single bucket in a B19 group_by distribution.
type CountBucket struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// CountGroupData is the B19 count response with group_by.
type CountGroupData struct {
	GroupBy string        `json:"group_by"`
	Buckets []CountBucket `json:"buckets"`
}

// CountMeta contains metadata for B19 count responses.
type CountMeta struct {
	CountedResource string    `json:"counted_resource"`
	TimeField       string    `json:"time_field"`
	AsOf            time.Time `json:"as_of"`
}

// CountResponse names the B19 count schema for Swagger.
type CountResponse struct {
	Data CountData `json:"data" binding:"required"`
	Meta CountMeta `json:"meta" binding:"required"`
}

// CountGroupResponse names the B19 group_by count schema for Swagger.
type CountGroupResponse struct {
	Data CountGroupData `json:"data" binding:"required"`
	Meta CountMeta      `json:"meta" binding:"required"`
}
