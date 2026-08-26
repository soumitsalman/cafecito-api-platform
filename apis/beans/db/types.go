package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// // Chatter represents short-form discussion metadata associated with a Bean.
// // @Description Single social or forum mention of a bean URL. Propagation responses use chatter-derived fields to show where an article was discussed and the lower-bound engagement observed at collection time.
// type Chatter struct {
// 	// ChatterURL is the URL of the social post, comment, or discussion item that mentions the Bean URL.
// 	ChatterURL string `db:"chatter_url" bson:"chatter_url" json:"chatter_url"`
// 	// URL is the referenced Bean URL that appeared in the social or forum mention.
// 	URL string `db:"url" bson:"url" json:"url"`
// 	// Source identifies the platform or publisher where the chatter was collected.
// 	Source string `db:"source" json:"source,omitempty"`
// 	// Forum is the community, group, subreddit, page, or forum where the mention was found.
// 	Forum string `db:"forum" bson:"group" json:"forum,omitempty"`
// 	// Collected is when the chatter metrics were collected from the external platform.
// 	Collected time.Time `db:"collected" json:"-" swaggertype:"string" format:"date-time"`
// 	// Likes is the cumulative lower-bound like or upvote count captured for the mention.
// 	Likes int64 `db:"likes" json:"likes,omitempty"`
// 	// Comments is the cumulative lower-bound reply or comment count captured for the mention.
// 	Comments int64 `db:"comments" json:"comments,omitempty"`
// 	// Subscribers is the cumulative lower-bound audience or follower count for the forum/community.
// 	Subscribers int64 `db:"subscribers" json:"subscribers,omitempty"`
// }

// // ChatterAggregate represents aggregated social engagement metrics for a Bean URL.
// // @Description Aggregated social traction for one bean URL. These metrics help rank trending/top-headline results and expose engagement context such as likes, comments, audience size, and shares.
// type ChatterAggregate struct {
// 	// URL is the Bean URL for which aggregate chatter metrics were computed.
// 	URL string `db:"url" json:"url,omitempty"`
// 	// Collected is the latest timestamp when any contributing chatter record was collected.
// 	Collected time.Time `db:"collected" json:"-" swaggertype:"string" format:"date-time"`
// 	// Likes is the aggregate number of likes or upvotes across collected chatter records.
// 	Likes int64 `db:"likes" json:"likes,omitempty"`
// 	// Comments is the aggregate number of replies or comments across collected chatter records.
// 	Comments int64 `db:"comments" json:"comments,omitempty"`
// 	// Subscribers is the aggregate audience size associated with contributing chatter records.
// 	Subscribers int64 `db:"subscribers" json:"subscribers,omitempty"`
// 	// Shares is the aggregate number of reposts, retweets, or share-like actions.
// 	Shares int64 `db:"shares" json:"shares,omitempty"`
// }

// // PropagationCoverage is the same story published by another outlet.
// // @Description One cross-publisher coverage hit for a seed article URL. Use it to see whether a story was republished or covered by another source.
// type PropagationCoverage struct {
// 	URL      string    `json:"url"`
// 	Created  time.Time `json:"created" swaggertype:"string" format:"date-time"`
// 	Source   string    `json:"source"`
// 	SiteName string    `json:"site_name"`
// }

// // PropagationMention is a social/forum mention of an article from chatters.
// // @Description One social or forum mention for a seed article URL, including where it appeared and any available engagement counts.
// type PropagationMention struct {
// 	ShareURL string    `json:"share_url"`
// 	Source   string    `json:"source"`
// 	Forum    string    `json:"forum,omitempty"`
// 	Observed time.Time `json:"observed" swaggertype:"string" format:"date-time"`
// 	Comments int64     `json:"comments,omitempty"`
// 	Likes    int64     `json:"likes,omitempty"`
// }

// // PropagationResult groups publisher coverage and social mentions for one seed URL.
// // @Description Propagation result for one input article URL. `coverage` shows related publisher articles; `mentions` shows social/forum discussion. Empty arrays mean no propagation was found for that URL.
// type PropagationResult struct {
// 	URL      string                `json:"url"`
// 	Coverage []PropagationCoverage `json:"coverage"`
// 	Mentions []PropagationMention  `json:"mentions"`
// }

const (
	_BEAN_COLUMNS_BASE         = "id, url, kind, created, author, image_url, categories, sentiments, entities, regions, title, source_id, base_url, domain_name, site_name, cluster_id"
	_BEAN_COLUMNS_SUMMARY      = "summary"
	_BEAN_COLUMNS_CONTENT      = "CASE WHEN restricted_content THEN NULL ELSE content END AS content"
	_BEAN_COLUMNS_TREND        = "likes, comments, mentions, subscribers, related, trend_score"
	_BEAN_COLUMNS_ALL          = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY + ", " + _BEAN_COLUMNS_CONTENT + ", " + _BEAN_COLUMNS_TREND
	BEAN_COLUMNS_HEADLINES     = _BEAN_COLUMNS_BASE
	BEAN_COLUMNS_WITHOUT_TREND = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY
	BEAN_COLUMNS_WITH_TREND    = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY + ", " + _BEAN_COLUMNS_TREND
)

const (
	SOURCE_COLUMNS_BASE = "id, base_url, domain_name, site_name"
	SOURCE_COLUMNS_ALL  = SOURCE_COLUMNS_BASE + ", description, favicon, rss_feed"
)

const (
	_CLUSTER_BEAN_COLUMNS_MINIMAL = "id, url, created, title, source_id, base_url, domain_name, site_name, cluster_id"
)

const (
	SORT_RECENT   = "created"
	SORT_TRENDING = "trend_score"
	SORT_RELEVANT = "relevance"
)

type TrendProperties struct {
	Likes       sql.NullInt64   `db:"likes"`
	Comments    sql.NullInt64   `db:"comments"`
	Mentions    sql.NullInt64   `db:"mentions"`
	Subscribers sql.NullInt64   `db:"subscribers"`
	Related     sql.NullInt64   `db:"related"`
	ClusterID   uuid.UUID       `db:"cluster_id"`
	TrendScore  sql.NullFloat64 `db:"trend_score"`
}

type SourceProperties struct {
	BaseURL     string         `db:"base_url"`
	DomainName  sql.NullString `db:"domain_name"`
	SiteName    sql.NullString `db:"site_name"`
	Description sql.NullString `db:"description"`
	Favicon     sql.NullString `db:"favicon"`
	RSSFeed     sql.NullString `db:"rss_feed"`
}

func (source *SourceProperties) IsZero() bool {
	return source.BaseURL == "" || !source.DomainName.Valid
}

// Bean represents a single article or post indexed by Beansack.
// @Description Primary article/post object returned by Beans article endpoints. Agents should treat `url` as the stable identifier, `source` as the publisher id, `summary` as the compact context field, and `content` as optional full text only present when requested. `categories`, `regions`, `entities`, `sentiments`, and `tags` are inferred enrichment fields for filtering and grounding responses. Internal embedding and gist fields are used for search but omitted from JSON.
type Bean struct {
	// URL is the canonical URL of the article or post.
	ID         uuid.UUID      `db:"id"`
	URL        string         `db:"url"`
	Kind       string         `db:"kind"`
	Created    time.Time      `db:"created"`
	Author     sql.NullString `db:"author"`
	ImageUrl   sql.NullString `db:"image_url"`
	Categories []string       `db:"categories"`
	Sentiments []string       `db:"sentiments"`
	Entities   []string       `db:"entities"`
	Regions    []string       `db:"regions"`
	Title      sql.NullString `db:"title"`
	Summary    sql.NullString `db:"summary"`
	Content    sql.NullString `db:"content"`
	SourceID   uuid.UUID      `db:"source_id"`
	BaseURL    string         `db:"base_url"`
	SourceProperties
	TrendProperties
	Distance sql.NullFloat64 `db:"distance"`
}

func (b *Bean) IsZero() bool {
	return b.ID == uuid.Nil || b.URL == ""
}

type Trend struct {
	ID       uuid.UUID    `db:"id"`
	URL      string       `db:"url"`
	Observed sql.NullTime `db:"observed"`
	TrendProperties
}

// Mention is one observed social or forum post linking an Article URL.
type Mention struct {
	URL         string         `db:"chatter_url"`
	Platform    string         `db:"platform"`
	Forum       sql.NullString `db:"forum"`
	Observed    time.Time      `db:"observed"`
	Likes       sql.NullInt64  `db:"likes"`
	Comments    sql.NullInt64  `db:"comments"`
	Subscribers sql.NullInt64  `db:"subscribers"`
}

// MentionFilters constrains B07 mention collection queries.
type MentionFilters struct {
	Platforms    []string
	Forums       []string
	ObservedFrom time.Time
	ObservedTo   time.Time
}

// Source holds metadata about a content source (publisher).
// ID is the canonical publisher identifier and matches Bean.Source values.
type Source struct {
	ID uuid.UUID `db:"id"`
	SourceProperties
}

func (source *Source) IsZero() bool {
	return source.ID == uuid.Nil || source.SourceProperties.IsZero()
}

type clusterBase struct {
	ID          uuid.UUID       `db:"id"`
	LastCreated time.Time       `db:"last_created"`
	Distance    sql.NullFloat64 `db:"distance"`
}

// Cluster is a derived group of related articles sharing a cluster_id.
type Cluster struct {
	clusterBase
	Title        string    `db:"title"`
	FirstCreated time.Time `db:"first_created"`
	BeanCount    int       `db:"bean_count"`
	SourceCount  int       `db:"source_count"`
	Categories   []string  `db:"categories"`
	Regions      []string  `db:"regions"`
	Entities     []string  `db:"entities"`
	Tags         []string  `db:"tags"`
	TopArticles  []Bean    `db:"-"`
}

func (s *Cluster) IsZero() bool {
	return s.ID == uuid.Nil || s.FirstCreated.IsZero() || s.LastCreated.IsZero()
}

type BeanFilters struct {
	IDs               []uuid.UUID
	URLs              []string
	Sources           []uuid.UUID
	ExcludeSources    []uuid.UUID
	Domains           []string
	ExcludeDomains    []string
	Kind              string
	CreatedFrom       time.Time
	CreatedTo         time.Time
	ObservedFrom      time.Time
	ObservedTo        time.Time
	Authors           []string
	Tags              []string
	Categories        []string
	ExcludeCategories []string
	Sentiments        []string
	Entities          []string
	Regions           []string
	FullContent       bool
	Embedding         []float32
	Distance          float64
	ClusterID         uuid.UUID
	Extra             []string
}

type SourceFilters struct {
	Q       string
	IDs     []uuid.UUID
	Domains []string
}

type ClusterFilters struct {
	BeanFilters
	MinBeanCount int
}

// DEPRECATED - ONLY APPLICABLE TO DUCKDB

// type sqlVector []float32
// type sqlStringArray []string

// func (vec sqlStringArray) Value() (driver.Value, error) {
// 	bytes, err := json.Marshal(vec)
// 	return driver.Value(string(bytes)), err
// }

// func (vec *sqlStringArray) Scan(value interface{}) error {
// 	if value == nil {
// 		*vec = nil
// 		return nil
// 	}

// 	switch value := value.(type) {
// 	case []interface{}:
// 		converted := make([]string, len(value))
// 		for i, val := range value {
// 			converted[i] = val.(string)
// 		}
// 		*vec = converted
// 	case []byte:
// 	case string:
// 		return json.Unmarshal([]byte(value), vec)
// 	default:
// 		return fmt.Errorf("unsupported type: %T", value)
// 	}
// 	return nil
// }

// func (vec sqlVector) Value() (driver.Value, error) {
// 	if len(vec) == 0 {
// 		return driver.Value(nil), fmt.Errorf("vector cannot be nil or empty")
// 	}
// 	bytes, err := json.Marshal(vec)
// 	return driver.Value(string(bytes)), err
// }

// func (vec *sqlVector) Scan(value interface{}) error {
// 	if value == nil {
// 		*vec = nil
// 		return nil
// 	}

// 	switch value := value.(type) {
// 	case []interface{}:
// 		converted := make([]float32, len(value))
// 		for i, val := range value {
// 			switch v := val.(type) {
// 			case float64:
// 				converted[i] = float32(v)
// 			case float32:
// 				converted[i] = v
// 			case int:
// 				converted[i] = float32(v)
// 			default:
// 				return fmt.Errorf("unsupported array element type: %T", val)
// 			}
// 		}
// 		*vec = converted
// 		return nil
// 	case []float32:
// 		*vec = value
// 		return nil
// 	case []float64:
// 		converted := make([]float32, len(value))
// 		for i, v := range value {
// 			converted[i] = float32(v)
// 		}
// 		*vec = converted
// 		return nil
// 	case []int:
// 		converted := make([]float32, len(value))
// 		for i, val := range value {
// 			converted[i] = float32(val)
// 		}
// 		*vec = converted
// 		return nil
// 	case []byte:
// 		return json.Unmarshal(value, vec)
// 	case string:
// 		return json.Unmarshal([]byte(value), vec)
// 	default:
// 		return fmt.Errorf("unsupported type: %T", value)
// 	}
// }
