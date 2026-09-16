package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const (
	_BEAN_COLUMNS_BASE              = "id, url, kind, created, author, image_url, language, categories, sentiments, entities, regions, title, source_id, base_url, domain_name, site_name, cluster_id"
	_BEAN_COLUMNS_SUMMARY           = "summary"
	_BEAN_COLUMNS_CONTENT           = "CASE WHEN restricted_content THEN NULL ELSE content END AS content"
	_BEAN_COLUMNS_TREND             = "likes, comments, mentions, subscribers, related, trend_score"
	_BEAN_COLUMNS_ALL               = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY + ", " + _BEAN_COLUMNS_CONTENT + ", " + _BEAN_COLUMNS_TREND
	BEAN_COLUMNS_HEADLINES          = _BEAN_COLUMNS_BASE
	BEAN_COLUMNS_WITHOUT_TREND      = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY
	BEAN_COLUMNS_WITH_TREND         = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY + ", " + _BEAN_COLUMNS_TREND
	BEAN_COLUMNS_MINIMAL            = "id, url, created, title, source_id, base_url, domain_name, site_name, cluster_id"
	BEAN_COLUMNS_MINIMAL_WITH_TREND = BEAN_COLUMNS_MINIMAL + ", " + _BEAN_COLUMNS_TREND
)

const (
	SOURCE_COLUMNS_BASE = "id, base_url, domain_name, site_name"
	SOURCE_COLUMNS_ALL  = SOURCE_COLUMNS_BASE + ", description, favicon, rss_feed"
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
	Language   sql.NullString `db:"language"`
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
	Summary      string    `db:"summary"`
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
	Languages         []string
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
