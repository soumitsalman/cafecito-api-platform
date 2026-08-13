# Beans API Route Proposal

Status: target public contract, target Go types, and target query plan  
Prepared: 2026-08-03  
Scope: Beans only  
Assumption: `beans.id` and `publishers.id` are non-null UUID primary keys before
any UUID route is released. This document does not plan that migration.

## Purpose

Beans is a public media-evidence API. Its atomic resource is an **article**:
one publisher document with a canonical URL, publication time, source, and
optional extracted content.

It is not an event-intelligence API. Espresso owns structured actions, events,
and signals. Beans may group articles into a cluster, but that group remains
media evidence rather than an asserted real-world event.

## Public-shape decisions

The target follows the common public vocabulary used by NewsAPI, Perigon, Event
Registry, and Newscatcher where it does not obscure Beans' product boundary.

- Use `articles` and `sources` as the principal resources. NewsAPI and Perigon
  both call a publisher/outlet/domain a source.
- Use top-level taxonomy resources: `categories`, `entities`, and `regions`.
  Services distinguish these concepts rather than returning one mixed tag list.
  Existing `/tags/*` routes are compatibility paths only.
- Use `clusters` rather than `stories`. Perigon calls clusters Stories; Beans
  deliberately avoids overloading “story,” which can also mean one article.
- Use `related` only for ranked semantic reading. It is not a cluster and is not
  evidence that all results describe the same development.
- Use `mentions` only for external social/forum posts observed linking to the
  anchor article URL. It is not editorial coverage.
- Reserve **coverage** as a documentation term for the distinct sources inside a
  cluster. The public resource is `/clusters/{cluster_id}/sources`, not an
  ambiguous article-level `/coverage` route.
- Use `from` and `to`, as public news APIs do. Each route documents the timestamp
  they constrain.

Market references:

- [NewsAPI endpoints](https://newsapi.org/docs/endpoints) — article search, top
  headlines, and sources.
- [Perigon data model](https://perigon.io/docs/api/data-model) — articles,
  sources, and Story-level clusters.
- [Perigon Sources](https://perigon.io/docs/api/sources) — source discovery and
  filtering.
- [Newscatcher clustering](https://www.newscatcherapi.com/docs/news-api/guides-and-concepts/clustering-news-articles)
  — query result clusters and cluster members.
- [Event Registry terminology](https://newsapi.ai/documentation?tab=terminology)
  — distinct concepts, categories, articles, and events.
- [Meltwater mention search](https://developer.meltwater.com/guides/listening/searching-mentions/)
  — social/listening mention documents and observation windows.

## Target response convention

Use one conventional JSON envelope for every canonical public route:

```json
{
  "data": [],
  "pagination": {
    "limit": 20,
    "next_cursor": null
  },
  "meta": {
    "as_of": "2026-08-03T16:00:00Z"
  }
}
```

Rules:

- Collection routes return `200` and `data: []` when no result matches.
- Detail routes return `404` for a syntactically valid but unknown UUID.
- `next_cursor` is opaque. It encodes the selected sort key plus UUID tie-breaker;
  callers must not construct or interpret it.
- `from` is inclusive and `to` is exclusive. Both are RFC 3339 timestamps; a date
  alone is accepted as midnight UTC.
- `as_of` is the materialized-data freshness timestamp when the response depends
  on a view, aggregate, mention collection, or cluster read model.
- Article lists and specialized article lists share `ArticleSummary`; specialized
  routes add a `trend`, `relation`, or `is_representative` field rather than
  changing article field names.
- Missing enrichment/engagement fields are `null`, never fabricated as zero.

## Target routes

| ID | Canonical public route | Purpose | Current public route / migration |
|---|---|---|---|
| B01 | `GET /beans/articles` | Search or browse individual articles. | Replaces `/articles/search` and `/articles/latest`. |
| B02 | `GET /beans/articles/{article_id}` | Retrieve one article. | New after UUID projection. |
| B03 | `GET /beans/articles/trending` | Rank articles by current attention. | Renames/normalizes existing `/articles/trending`. |
| B04 | `GET /beans/headlines` | Fixed-window high-attention briefing feed. | Replaces `/articles/top-headlines`. |
| B05 | `GET /beans/articles/{article_id}/related` | Ranked semantic reading related to one article. | New typed relation contract; do not expose from raw propagation. |
| B06 | `GET /beans/articles/{article_id}/mentions` | External social/forum posts linking to one article. | Split from `/articles/propagation`. |
| B07 | `GET /beans/sources` | Discover article-producing sources. | Existing `/sources`, normalized. |
| B08 | `GET /beans/sources/{source_id}` | Retrieve one source. | New after UUID projection. |
| B09a | `GET /beans/categories` | Discover valid `categories` values. | `/tags/categories` becomes compatibility path. |
| B09b | `GET /beans/entities` | Discover valid `entities` values. | `/tags/entities` becomes compatibility path. |
| B09c | `GET /beans/regions` | Discover valid `regions` values. | `/tags/regions` becomes compatibility path. |
| B10 | `GET /beans/clusters` | Search durable same-subject article groups. | New read model. |
| B11 | `GET /beans/clusters/{cluster_id}` | Retrieve cluster summary and scope. | New read model. |
| B12 | `GET /beans/clusters/{cluster_id}/articles` | Retrieve attributable cluster-member articles. | New read model. |
| B13 | `GET /beans/clusters/{cluster_id}/sources` | Retrieve editorial source coverage for a cluster. | New read model. |
| B14 | `GET /beans/articles/counts` | Bounded counts for an approved consumer. | Deferred; no release without a named consumer. |

### How similarly named routes differ

| Route | Exact question answered | It does **not** mean |
|---|---|---|
| `/articles` | Which publisher documents match my query or filters? | One underlying development or social attention. |
| `/articles/trending` | Which articles have the strongest current measured attention? | The newest articles or a complete story group. |
| `/headlines` | Which recently published, high-attention articles belong in a briefing? | All latest articles or an editorially curated human feed. |
| `/articles/{id}/related` | What other articles are semantically useful to read after this one? | Same-subject cluster membership or a social share. |
| `/articles/{id}/mentions` | Which external posts linked to this exact article URL? | Other publisher articles or generic discussion of the topic. |
| `/clusters` | Which durable groups represent the same developing subject? | A loose semantic recommendation list. |
| `/clusters/{id}/sources` | Which editorial sources are represented by articles in this cluster? | Social platforms, shares, or social mentions. |

## TARGET ROUTES

### B01 — Search or browse articles

`GET /beans/articles`

**Documentation:** This is the main article-search resource, comparable to
NewsAPI Everything and Perigon Articles. With `q`, it performs semantic search;
without `q`, it is a deterministic publication-date browse feed. It never ranks
by trend score; use `/articles/trending` for that.

| Query parameter | Type | Meaning |
|---|---|---|
| `q` | string, optional | Semantic search text, maximum 512 characters. |
| `source_ids` | CSV UUIDs | Include articles from these sources. |
| `content_types` | CSV enum | `news`, `blog`, or `post`; default is all public article kinds. |
| `categories`, `entities`, `regions` | CSV strings | Exact values returned by the matching taxonomy route. |
| `from`, `to` | RFC 3339 | Inclusive/exclusive `published_at` range. |
| `sort` | enum | `relevance` only with `q`; otherwise `published_at`. Default: `relevance` with `q`, `published_at` without. |
| `include` | CSV enum | Optional `content`; omitted by default. |
| `limit`, `cursor` | pagination | Default 20, max 100; cursor is opaque. |

**Response model:** `ListResponse[ArticleSummary]`

```json
{
  "data": [
    {
      "id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
      "url": "https://example.com/chips",
      "title": "Chip supply tightens after new export controls",
      "summary": "Manufacturers expect longer lead times.",
      "published_at": "2026-08-03T13:10:00Z",
      "content_type": "news",
      "source": {
        "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
        "name": "Example Business",
        "domain": "example.com"
      },
      "categories": ["technology", "business"],
      "entities": ["Example Semiconductor Co"]
    }
  ],
  "pagination": { "limit": 20, "next_cursor": null },
  "meta": { "query_mode": "semantic", "as_of": "2026-08-03T16:00:00Z" }
}
```

### B02 — Retrieve one article

`GET /beans/articles/{article_id}`

**Documentation:** Stable UUID lookup for one article. URL remains an external
attribute and is still accepted by B01 as `urls` during compatibility, but UUID
is the canonical resource identifier.

| Query parameter | Type | Meaning |
|---|---|---|
| `include` | CSV enum | `content` returns full extracted content when available. |

**Response model:** `ObjectResponse[ArticleDetail]`

`ArticleDetail` is `ArticleSummary` plus `author`, `image_url`, optional
`content`, and `content_access.status` (`available` or `unavailable`). Do not
claim a specific unavailability reason unless ingestion records it.

### B03 — Retrieve trending articles

`GET /beans/articles/trending`

**Documentation:** This is an attention-ranked article feed, not a general
article sort. It returns the same base article/source shape as B01 plus the
measured inputs to the trend score. `from` and `to` constrain article
`published_at`; attention freshness is exposed as `meta.as_of`.

| Query parameter | Type | Meaning |
|---|---|---|
| `source_ids`, `content_types`, `categories`, `entities`, `regions` | filters | Same values and semantics as B01. |
| `from`, `to` | RFC 3339 | Article publication window. |
| `limit`, `cursor` | pagination | Ordered by trend score descending; no caller-selected sort. |

**Response model:** `ListResponse[TrendingArticle]`

```json
{
  "data": [
    {
      "id": "725a7f65-f7eb-4643-b7af-78e8980217e9",
      "url": "https://example.com/accelerator",
      "title": "Chipmaker announces a new accelerator",
      "published_at": "2026-08-02T18:20:00Z",
      "source": { "id": "a2f08211-1796-43df-8988-c35a2923a7a0", "name": "Example Technology", "domain": "example.com" },
      "trend": {
        "like_count": 820,
        "comment_count": 146,
        "mention_count": 38,
        "related_article_count": 19,
        "score": 572.4
      }
    }
  ],
  "pagination": { "limit": 20, "next_cursor": null },
  "meta": {
    "ranked_by": "trend_score",
    "score_inputs": ["related_articles", "comments", "mentions", "likes", "recency"],
    "as_of": "2026-08-03T16:00:00Z"
  }
}
```

### B04 — Retrieve headlines

`GET /beans/headlines`

**Documentation:** A fixed recent briefing feed comparable to NewsAPI Top
Headlines. It selects high-attention articles published in the previous 24
hours. It is neither a generic latest feed nor a manually editorially curated
feed. The fixed time window is intentional; `from` and `to` are not accepted.

| Query parameter | Type | Meaning |
|---|---|---|
| `source_ids`, `categories`, `regions` | filters | Narrow the fixed headline feed. |
| `limit`, `cursor` | pagination | Ordered by the headline selection rank. |

**Response model:** `ListResponse[TrendingArticle]`, with
`meta.selection = "highest attention among articles published in the prior 24 hours"`.

### B05 — Retrieve related articles

`GET /beans/articles/{article_id}/related`

**Documentation:** A bounded, relevance-ranked reading list anchored on one
article. Results can be from the same source on another day or from different
sources. They are not asserted to be reports of the same developing subject;
use clusters for that stronger claim.

| Query parameter | Type | Meaning |
|---|---|---|
| `source_ids` | CSV UUIDs | Restrict returned related articles to these sources. |
| `exclude_anchor_source` | boolean | Exclude the source of the anchor article; default `false`. |
| `from`, `to` | RFC 3339 | Related article `published_at` window. |
| `limit`, `cursor` | pagination | Always relevance-ranked; no `q` or `sort`. |

**Response model:** `ListResponse[ArticleSummary]`

`relation_type` is response metadata because every result from this route has
the same relation type. It is not repeated inside every article item. If a
future route can return mixed relation types, introduce a separate mixed-type
response model rather than silently changing this one.

```json
{
  "data": [
    {
      "id": "f9562e73-404c-4eb2-8fa7-c61ad4ef4e0a",
      "url": "https://example.net/suppliers-chip-rules",
      "title": "Suppliers respond to new chip rules",
      "summary": "Suppliers assess the effect of the controls.",
      "published_at": "2026-08-03T13:25:00Z",
      "content_type": "news",
      "source": { "id": "d20b13f5-b0e4-4a9f-ae79-08f824d11fe6", "name": "Example Technology", "domain": "example.net" }
    }
  ],
  "pagination": { "limit": 20, "next_cursor": null },
  "meta": {
    "anchor_article_id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
    "relation_type": "semantic_related"
  }
}
```

### B06 — Retrieve article mentions

`GET /beans/articles/{article_id}/mentions`

**Documentation:** A social/forum evidence feed for one article URL. A result
is an external post that directly references the anchor article URL. It is not
a generic social search and does not return other news publishers’ articles.

| Query parameter | Type | Meaning |
|---|---|---|
| `platforms` | CSV strings | Restrict to normalized platforms/forums, for example `reddit,x`. |
| `from`, `to` | RFC 3339 | Mention `observed_at` window, not article publication time. |
| `sort` | enum | `observed_at` (default) or `engagement`. |
| `limit`, `cursor` | pagination | Mention pagination. |

**Response model:** `ListResponse[ArticleMention]`

```json
{
  "data": [
    {
      "url": "https://reddit.example/post/123",
      "platform": "reddit",
      "forum": "r/technology",
      "observed_at": "2026-08-03T14:02:00Z",
      "engagement": {
        "like_count": 184,
        "comment_count": 37,
        "share_count": null
      }
    }
  ],
  "pagination": { "limit": 20, "next_cursor": null },
  "meta": { "anchor_article_id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75", "as_of": "2026-08-03T16:00:00Z" }
}
```

### B07 — List sources

`GET /beans/sources`

**Documentation:** Source discovery, comparable to NewsAPI Sources and Perigon
Sources. A source is an article-producing publisher, outlet, blog, feed, or
domain. It is not a social platform; social platforms appear only in mentions.

| Query parameter | Type | Meaning |
|---|---|---|
| `q` | string | Case-insensitive source name or domain prefix search. |
| `domains` | CSV strings | Exact source-domain filter. |
| `limit`, `cursor` | pagination | Ordered by normalized source name then ID. |

**Response model:** `ListResponse[Source]`

### B08 — Retrieve one source

`GET /beans/sources/{source_id}`

**Documentation:** Stable source profile lookup. Article responses embed the
same `SourceSummary` so agents generally need this route only for fuller source
metadata.

**Response model:** `ObjectResponse[Source]`

### B09a/B09b/B09c — Discover article taxonomy values

```http
GET /beans/categories
GET /beans/entities
GET /beans/regions
```

**Documentation:** These are separate taxonomy resources because a category is
not an entity and neither is a region. Each returned `value` is directly
accepted by its corresponding B01 filter. Existing `/tags/categories`,
`/tags/entities`, and `/tags/regions` stay as deprecation aliases.

| Query parameter | Type | Meaning |
|---|---|---|
| `q` | string | Case-insensitive prefix search over the displayed value. |
| `limit`, `cursor` | pagination | Alphabetical keyset pagination. |

**Response model:** `ListResponse[TaxonomyValue]`

```json
{
  "data": [
    { "value": "technology", "label": "Technology" },
    { "value": "technology_policy", "label": "Technology policy" }
  ],
  "pagination": { "limit": 20, "next_cursor": null }
}
```

### B10 — Search clusters

`GET /beans/clusters`

**Documentation:** A cluster is a durable, same-subject group of articles. It
is the closest Beans equivalent to a Perigon Story. It is stronger than B05:
cluster membership means the system groups articles under one developing
subject; related means only that an article is useful semantic follow-up.

| Query parameter | Type | Meaning |
|---|---|---|
| `q` | string | Matches cluster name, summary, and aggregate entities. |
| `source_ids`, `categories`, `entities`, `regions` | filters | A cluster matches when at least one member matches each supplied filter group. |
| `from`, `to` | RFC 3339 | A cluster matches when it has at least one member published in the time window. |
| `min_sources` | integer | Minimum distinct editorial sources in the complete cluster. |
| `sort` | enum | `updated_at` (default), `article_count`, or `relevance` with `q`. |
| `limit`, `cursor` | pagination | Cluster pagination. |

**Response model:** `ListResponse[ClusterSummary]`

### B11 — Retrieve one cluster

`GET /beans/clusters/{cluster_id}`

**Documentation:** Returns the current durable state of one cluster. Its
`article_count` and `source_count` are whole-cluster counts as of `meta.as_of`,
not social mention totals.

**Response model:** `ObjectResponse[ClusterDetail]`

### B12 — Retrieve cluster articles

`GET /beans/clusters/{cluster_id}/articles`

**Documentation:** Returns the attributable article evidence for a cluster.
This is the correct route when a client needs the publisher documents behind a
cluster summary.

| Query parameter | Type | Meaning |
|---|---|---|
| `source_ids`, `from`, `to` | filters | Restrict member articles by source and `published_at`. |
| `sort` | enum | `published_at` (default) or `representative_first`. |
| `limit`, `cursor` | pagination | Member pagination. |

**Response model:** `ListResponse[ClusterArticle]`

### B13 — Retrieve cluster source coverage

`GET /beans/clusters/{cluster_id}/sources`

**Documentation:** “Coverage” means the distinct editorial sources represented
in the cluster. This route is intentionally below the cluster resource because
coverage cannot be defined correctly from an arbitrary single article. It does
not return social platforms or social engagement.

| Query parameter | Type | Meaning |
|---|---|---|
| `sort` | enum | `article_count` (default) or `last_published_at`. |
| `limit`, `cursor` | pagination | Source coverage pagination. |

**Response model:** `ListResponse[ClusterSourceCoverage]`

### B14 — Count matching articles

`GET /beans/articles/counts`

**Documentation:** Numeric aggregation only. It is not an analytics namespace
and never returns article records. Release only when a named dashboard, alert,
or reporting consumer requires it.

| Query parameter | Type | Meaning |
|---|---|---|
| Article filters | same as B01 except `include`, `cursor` | Defines the counted article set. Semantic `q` is disallowed until its counting cost is approved. |
| `from`, `to` | RFC 3339 | Article `published_at` window. |
| `group_by` | enum | `published_day`, `source`, `category`, or `region`. |

**Response model:** `ObjectResponse[ArticleCounts]`

## TARGET TYPES

These are target shapes only. They distinguish persistence rows in
`apis/beans/db/types.go` from public HTTP input/output types in
`apis/beans/router/routes.go`. Database types must not carry public JSON
contract tags; router types own JSON names and OpenAPI annotations.

### `apis/beans/db/types.go`

```go
package db

import (
    "context"
    "time"

    "github.com/google/uuid"
)

type CursorPage struct {
    Limit  int
    Cursor string
}

type TimeRange struct {
    From *time.Time
    To   *time.Time
}

type ArticleFilter struct {
    IDs          []uuid.UUID
    URLs         []string
    SourceIDs    []uuid.UUID
    ContentTypes []string
    Categories   []string
    Entities     []string
    Regions      []string
    TimeRange    TimeRange
    QueryVector  []float32
    QueryText    string
}

type ArticleSort string

const (
    ArticleSortPublishedAt ArticleSort = "published_at"
    ArticleSortRelevance   ArticleSort = "relevance"
)

// ArticleRow is the canonical persistence projection. SourceKey remains only
// for joins while the ingestion pipeline transitions from source strings.
type ArticleRow struct {
    ID          uuid.UUID
    URL         string
    SourceID    uuid.UUID
    SourceKey   string
    Kind        string
    Title       string
    Summary     string
    Content     *string
    Author      *string
    ImageURL    *string
    PublishedAt time.Time
    Categories  []string
    Entities    []string
    Regions     []string
}

type SourceRow struct {
    ID          uuid.UUID
    SourceKey   string
    Name        string
    Domain      string
    URL         string
    Description *string
    FaviconURL  *string
    RSSFeedURL  *string
    CollectedAt *time.Time
}

type TrendRow struct {
    ArticleRow
    LikeCount           *int64
    CommentCount        *int64
    MentionCount        *int64
    RelatedArticleCount *int64
    Score               float64
    UpdatedAt           time.Time
}

type ArticleRelationType string

const ArticleRelationSemantic ArticleRelationType = "semantic_related"

type ArticleRelationRow struct {
    AnchorArticleID uuid.UUID
    ArticleID       uuid.UUID
    Type            ArticleRelationType
    Rank            int
    Score           *float64 // nullable until calibrated and public
    ComputedAt      time.Time
}

type MentionRow struct {
    URL         string // external post/comment URL; no fabricated local ID
    ArticleID   uuid.UUID
    Platform    string
    Forum       *string
    ObservedAt  time.Time
    LikeCount   *int64
    CommentCount *int64
    ShareCount  *int64
}

type TaxonomyKind string

const (
    TaxonomyCategory TaxonomyKind = "category"
    TaxonomyEntity   TaxonomyKind = "entity"
    TaxonomyRegion   TaxonomyKind = "region"
)

type TaxonomyValueRow struct {
    Kind  TaxonomyKind
    Value string
    Label *string
}

type ClusterRow struct {
    ID                    uuid.UUID
    Name                  string
    Summary               *string
    FirstPublishedAt      time.Time
    LastPublishedAt       time.Time
    UpdatedAt             time.Time
    ArticleCount          int64
    SourceCount           int64
    RepresentativeArticleID uuid.UUID
    Categories            []string
    Entities              []string
}

type ClusterMemberRow struct {
    ClusterID        uuid.UUID
    ArticleID        uuid.UUID
    IsRepresentative bool
    AddedAt          time.Time
}

type ClusterSourceCoverageRow struct {
    ClusterID       uuid.UUID
    Source          SourceRow
    ArticleCount    int64
    FirstPublishedAt time.Time
    LastPublishedAt  time.Time
}

type ClusterFilter struct {
    QueryText   string
    SourceIDs   []uuid.UUID
    Categories  []string
    Entities    []string
    Regions     []string
    TimeRange   TimeRange
    MinSources  *int
}

type CountGroup string

const (
    CountGroupPublishedDay CountGroup = "published_day"
    CountGroupSource       CountGroup = "source"
    CountGroupCategory     CountGroup = "category"
    CountGroupRegion       CountGroup = "region"
)

type CountRow struct {
    Key   string
    Count int64
}

type Beansack interface {
    ListArticles(context.Context, ArticleFilter, ArticleSort, CursorPage) ([]ArticleRow, string, error)
    GetArticle(context.Context, uuid.UUID, bool) (ArticleRow, error)
    ListTrendingArticles(context.Context, ArticleFilter, CursorPage) ([]TrendRow, string, time.Time, error)
    ListHeadlines(context.Context, ArticleFilter, CursorPage) ([]TrendRow, string, time.Time, error)
    ListRelatedArticles(context.Context, uuid.UUID, ArticleFilter, CursorPage) ([]ArticleRow, string, error)
    ListArticleMentions(context.Context, uuid.UUID, []string, TimeRange, string, CursorPage) ([]MentionRow, string, time.Time, error)
    ListSources(context.Context, string, []string, CursorPage) ([]SourceRow, string, error)
    GetSource(context.Context, uuid.UUID) (SourceRow, error)
    ListTaxonomyValues(context.Context, TaxonomyKind, string, CursorPage) ([]TaxonomyValueRow, string, error)
    ListClusters(context.Context, ClusterFilter, string, CursorPage) ([]ClusterRow, string, time.Time, error)
    GetCluster(context.Context, uuid.UUID) (ClusterRow, time.Time, error)
    ListClusterArticles(context.Context, uuid.UUID, ArticleFilter, string, CursorPage) ([]ClusterMemberRow, []ArticleRow, string, error)
    ListClusterSources(context.Context, uuid.UUID, string, CursorPage) ([]ClusterSourceCoverageRow, string, error)
    CountArticles(context.Context, ArticleFilter, CountGroup) ([]CountRow, time.Time, error)
}
```

### `apis/beans/router/routes.go`

```go
package router

import (
    "time"

    "github.com/google/uuid"
)

type CursorInput struct {
    Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
    Cursor string `form:"cursor" binding:"omitempty,max=2048"`
}

type TimeRangeInput struct {
    From *time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
    To   *time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
}

type ArticleListInput struct {
    Q            string      `form:"q" binding:"omitempty,max=512"`
    SourceIDs    []uuid.UUID `form:"source_ids" collection_format:"csv"`
    ContentTypes []string    `form:"content_types" collection_format:"csv"`
    Categories   []string    `form:"categories" collection_format:"csv"`
    Entities     []string    `form:"entities" collection_format:"csv"`
    Regions      []string    `form:"regions" collection_format:"csv"`
    Sort         string      `form:"sort" binding:"omitempty,oneof=published_at relevance"`
    Include      []string    `form:"include" collection_format:"csv"`
    TimeRangeInput
    CursorInput
}

type TrendingInput struct {
    SourceIDs  []uuid.UUID `form:"source_ids" collection_format:"csv"`
    Categories []string    `form:"categories" collection_format:"csv"`
    Entities   []string    `form:"entities" collection_format:"csv"`
    Regions    []string    `form:"regions" collection_format:"csv"`
    TimeRangeInput
    CursorInput
}

type HeadlinesInput struct {
    SourceIDs  []uuid.UUID `form:"source_ids" collection_format:"csv"`
    Categories []string    `form:"categories" collection_format:"csv"`
    Regions    []string    `form:"regions" collection_format:"csv"`
    CursorInput
}

type RelatedInput struct {
    SourceIDs           []uuid.UUID `form:"source_ids" collection_format:"csv"`
    ExcludeAnchorSource bool        `form:"exclude_anchor_source,default=false"`
    TimeRangeInput
    CursorInput
}

type MentionsInput struct {
    Platforms []string `form:"platforms" collection_format:"csv"`
    Sort      string   `form:"sort" binding:"omitempty,oneof=observed_at engagement"`
    TimeRangeInput // documented as observed_at for this route
    CursorInput
}

type SourcesInput struct {
    Q       string   `form:"q" binding:"omitempty,max=256"`
    Domains []string `form:"domains" collection_format:"csv"`
    CursorInput
}

type TaxonomyInput struct {
    Q string `form:"q" binding:"omitempty,max=256"`
    CursorInput
}

type ClustersInput struct {
    Q          string      `form:"q" binding:"omitempty,max=512"`
    SourceIDs  []uuid.UUID `form:"source_ids" collection_format:"csv"`
    Categories []string    `form:"categories" collection_format:"csv"`
    Entities   []string    `form:"entities" collection_format:"csv"`
    Regions    []string    `form:"regions" collection_format:"csv"`
    MinSources *int        `form:"min_sources" binding:"omitempty,min=1"`
    Sort       string      `form:"sort" binding:"omitempty,oneof=updated_at article_count relevance"`
    TimeRangeInput
    CursorInput
}

type ClusterArticlesInput struct {
    SourceIDs []uuid.UUID `form:"source_ids" collection_format:"csv"`
    Sort      string      `form:"sort" binding:"omitempty,oneof=published_at representative_first"`
    TimeRangeInput
    CursorInput
}

type ClusterSourcesInput struct {
    Sort string `form:"sort" binding:"omitempty,oneof=article_count last_published_at"`
    CursorInput
}

type ArticleCountsInput struct {
    ArticleListInput
    GroupBy string `form:"group_by" binding:"required,oneof=published_day source category region"`
}

type Pagination struct {
    Limit      int     `json:"limit"`
    NextCursor *string `json:"next_cursor"`
}

type Meta struct {
    AsOf             *time.Time `json:"as_of,omitempty"`
    AnchorArticleID  *uuid.UUID `json:"anchor_article_id,omitempty"`
    RelationType     string     `json:"relation_type,omitempty"`
    QueryMode        string     `json:"query_mode,omitempty"`
    RankedBy         string     `json:"ranked_by,omitempty"`
    Selection        string     `json:"selection,omitempty"`
    Definition       string     `json:"definition,omitempty"`
}

type ListResponse[T any] struct {
    Data       []T        `json:"data"`
    Pagination Pagination `json:"pagination"`
    Meta       Meta       `json:"meta,omitempty"`
}

type ObjectResponse[T any] struct {
    Data T    `json:"data"`
    Meta Meta `json:"meta,omitempty"`
}

type SourceSummary struct {
    ID     uuid.UUID `json:"id"`
    Name   string    `json:"name"`
    Domain string    `json:"domain"`
}

type Source struct {
    SourceSummary
    URL         string  `json:"url"`
    Description *string `json:"description"`
    FaviconURL  *string `json:"favicon_url"`
}

type ArticleSummary struct {
    ID          uuid.UUID     `json:"id"`
    URL         string        `json:"url"`
    Title       string        `json:"title"`
    Summary     string        `json:"summary"`
    PublishedAt time.Time     `json:"published_at"`
    ContentType string        `json:"content_type"`
    Source      SourceSummary `json:"source"`
    Categories  []string      `json:"categories,omitempty"`
    Entities    []string      `json:"entities,omitempty"`
    Regions     []string      `json:"regions,omitempty"`
}

type ContentAccess struct {
    Status string `json:"status"` // available or unavailable
}

type ArticleDetail struct {
    ArticleSummary
    Author        *string        `json:"author"`
    ImageURL      *string        `json:"image_url"`
    Content       *string        `json:"content,omitempty"`
    ContentAccess ContentAccess  `json:"content_access"`
}

type Trend struct {
    LikeCount           *int64  `json:"like_count"`
    CommentCount        *int64  `json:"comment_count"`
    MentionCount        *int64  `json:"mention_count"`
    RelatedArticleCount *int64  `json:"related_article_count"`
    Score               float64 `json:"score"`
}

type TrendingArticle struct {
    ArticleSummary
    Trend Trend `json:"trend"`
}

type Engagement struct {
    LikeCount    *int64 `json:"like_count"`
    CommentCount *int64 `json:"comment_count"`
    ShareCount   *int64 `json:"share_count"`
}

type ArticleMention struct {
    URL        string      `json:"url"`
    Platform   string      `json:"platform"`
    Forum      *string     `json:"forum"`
    ObservedAt time.Time   `json:"observed_at"`
    Engagement Engagement  `json:"engagement"`
}

type TaxonomyValue struct {
    Value string  `json:"value"`
    Label *string `json:"label,omitempty"`
}

type ClusterSummary struct {
    ID                  uuid.UUID       `json:"id"`
    Name                string          `json:"name"`
    Summary             *string         `json:"summary"`
    FirstPublishedAt    time.Time       `json:"first_published_at"`
    LastPublishedAt     time.Time       `json:"last_published_at"`
    UpdatedAt           time.Time       `json:"updated_at"`
    ArticleCount        int64           `json:"article_count"`
    SourceCount         int64           `json:"source_count"`
    RepresentativeArticle ArticleSummary `json:"representative_article"`
}

type ClusterDetail struct {
    ClusterSummary
    Categories []string `json:"categories"`
    Entities   []string `json:"entities"`
}

type ClusterArticle struct {
    ArticleSummary
    IsRepresentative bool `json:"is_representative"`
}

type ClusterSourceCoverage struct {
    Source           SourceSummary `json:"source"`
    ArticleCount     int64         `json:"article_count"`
    FirstPublishedAt time.Time     `json:"first_published_at"`
    LastPublishedAt  time.Time     `json:"last_published_at"`
}

type CountValue struct {
    Key   string `json:"key"`
    Count int64  `json:"count"`
}

type ArticleCounts struct {
    GroupBy string       `json:"group_by"`
    Values  []CountValue `json:"values"`
}
```

## TARGET QUERY

The target uses keyset pagination, not `offset`, for canonical routes. Every
sort order includes the UUID as a final tie-breaker. Cursors must encode only
that route’s normalized sort values and UUID.

| Route | Target repository call | Query plan and required conditions | Required data/query change |
|---|---|---|---|
| B01 `/articles` | `ListArticles(filter, sort, page)` | Join `beans b` to `publishers s` on the transitional source string; select article plus source UUID/profile. With `q`, embed query and apply pgvector threshold/order by `distance ASC, b.id`; without `q`, order `b.created DESC, b.id DESC`. Apply source, taxonomy, kind, and publication-time predicates before keyset predicate. | Add IDs to projections; source UUID join; `to`; deterministic sort; keyset cursor. |
| B02 `/articles/{id}` | `GetArticle(id, includeContent)` | `SELECT` one bean by `b.id`, left join its source profile. Content projection is conditional; do not filter out the whole article because content is unavailable. | Article ID projection and single-row error mapping. |
| B03 `/articles/trending` | `ListTrendingArticles(filter, page)` | Query a refreshed trend view that projects `b.id`, `s.id`, raw attention counts, `trend_score`, and view refresh time. Filter B01 scalar fields against article properties and `b.created`; order `trend_score DESC, b.id DESC`. | Trend view must expose source/article UUIDs and reliable refresh timestamp. Rename current derived `shares` count to `mention_count` at HTTP boundary. |
| B04 `/headlines` | `ListHeadlines(filter, page)` | Same trend view, with `b.created >= now() - interval '24 hours'`; apply permitted source/category/region filters; use stable selection ordering. | Make the fixed window and selection rule explicit; do not accept arbitrary time range. |
| B05 `/articles/{id}/related` | `ListRelatedArticles(anchorID, filter, page)` | Resolve anchor UUID to canonical URL. Query a typed relation read model in both stored directions, join the related article and source, exclude anchor, apply returned-article filters, order relation rank/score then related ID. The route fixes the relation type to `semantic_related`, so the repository returns article rows only. | `related_beans` alone is insufficient unless it gains `relationship_type`, rank/score, and bidirectional lookup semantics. |
| B06 `/articles/{id}/mentions` | `ListArticleMentions(id, platforms, observedRange, sort, page)` | Resolve article UUID to canonical URL; query `chatters` where `c.url = anchor.url`; normalize `c.source`/`c.forum` to platform; filter `c.collected`; order `collected DESC, chatter_url` or engagement tuple then `chatter_url`. | Add platform normalization and keyset support. Preserve null metrics rather than coercing zero. |
| B07 `/sources` | `ListSources(q, domains, page)` | Query `publishers`; match `site_name` or normalized domain with prefix search, optional exact domain filter; order normalized name, ID. | Add `publishers.id` projection, domain derivation/index, name/domain search, cursor. |
| B08 `/sources/{id}` | `GetSource(id)` | `SELECT` source profile by `publishers.id`. | Source UUID lookup/index. |
| B09a `/categories` | `ListTaxonomyValues(category, q, page)` | `SELECT DISTINCT unnest(b.categories)` with normalized prefix condition; order normalized value. | Canonical normalization, label mapping if desired, cursor. |
| B09b `/entities` | `ListTaxonomyValues(entity, q, page)` | `SELECT DISTINCT unnest(b.entities)` with normalized prefix condition; order normalized value. | Same; do not add opaque IDs until entity resolution exists. |
| B09c `/regions` | `ListTaxonomyValues(region, q, page)` | `SELECT DISTINCT unnest(b.regions)` with normalized prefix condition; order normalized value. | Same. |
| B10 `/clusters` | `ListClusters(filter, sort, page)` | Query durable `clusters` joined to aggregate cluster taxonomy/source metrics. `q` matches cluster name/summary/entities. Taxonomy/source/time filters use `EXISTS` against `cluster_members` joined to articles. `min_sources` uses whole-cluster `source_count`. | New `clusters` and `cluster_members` read model; aggregate metadata; cursor strategy. |
| B11 `/clusters/{id}` | `GetCluster(id)` | Fetch durable cluster summary plus aggregate categories/entities and representative article/source. | Same cluster read model and 404 mapping. |
| B12 `/clusters/{id}/articles` | `ListClusterArticles(id, filter, sort, page)` | Join `cluster_members` to articles and sources; filter member article source/time; order `published_at DESC, article_id DESC` or representative flag then same. | Stable member relation and pagination index. |
| B13 `/clusters/{id}/sources` | `ListClusterSources(id, sort, page)` | Join cluster members → articles → sources; group by source UUID; calculate member article count and first/last publication; order requested aggregate plus source UUID. | Stable membership, aggregate query/index. |
| B14 `/articles/counts` | `CountArticles(filter, group)` | Reuse B01 scalar predicates only. Allowlist `group_by`; group by date bucket, source UUID, or unnest taxonomy. Do not run vector/semantic counts until query cost is approved. | Named consumer, query plan limits, freshness policy. |

### Target cursor rules

| Route sort | Cursor sort tuple |
|---|---|
| Article `published_at` | `published_at`, `article_id` |
| Article `relevance` | `distance`, `article_id` |
| Trending | `trend_score`, `article_id` |
| Headlines | `selection_score`, `article_id` |
| Related | `relation_rank`, `related_article_id` |
| Mentions `observed_at` | `observed_at`, `mention_url` |
| Mentions `engagement` | normalized engagement score, `observed_at`, `mention_url` |
| Sources/taxonomy | normalized name/value, ID/value |
| Clusters | selected aggregate/time key, `cluster_id` |
| Cluster members | selected member key, `article_id` |
| Cluster source coverage | selected aggregate/time key, `source_id` |

## Compatibility and release order

1. Document existing behavior and publish the new public vocabulary without
   claiming response equivalence.
2. Complete UUID projection for articles and sources, including every list and
   specialized result projection.
3. Publish B01–B04 and B07–B09 with the canonical envelope and cursor contract.
4. Split legacy propagation into B05 only after typed semantic relations exist,
   and B06 after source-to-article UUID translation and platform normalization.
5. Publish B10–B13 only after stable cluster identity, membership, source
   coverage, and refresh rules exist.
6. Publish B14 only after a named aggregation consumer approves bounded cost.
7. Deprecate `/articles/search`, `/articles/latest`, `/articles/top-headlines`,
   `/articles/propagation`, and `/tags/*` only after usage measurement and a
   migration window.

## Definition of done

The public Beans API is complete when a developer or agent can discover a
source or taxonomy value, retrieve attributable articles, distinguish semantic
related reading from same-subject clusters, distinguish editorial source
coverage from social mentions, and know exactly which timestamp every time
filter addresses.
