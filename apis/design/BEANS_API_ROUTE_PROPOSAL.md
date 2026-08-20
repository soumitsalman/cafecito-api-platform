# Beans API Target Design and Implementation Specification

Status: Target design and implementation specification  
Updated: 2026-08-17  
Scope: Beans News, Blogs, and Publisher Content API only  
Excluded: Espresso Events and Signals API  
Comparison baseline: [NEWS_AND_BLOG_API_MARKET_REPORT.md](NEWS_AND_BLOG_API_MARKET_REPORT.md)

## 1. Purpose and product boundary

This document defines the target public contract for a renovated Beans API.
Beans is an article-first service for news, blogs, reports, analysis, and other
publisher content. It returns the publisher material itself, its source
metadata, related coverage, story groups, and observed audience attention.

Beans does not expose parsed market Events or Signals. Those concepts belong to
Espresso and are outside this proposal.

The stored content types available to Beans V1 are:
```json
[{
  "kind": "blog"
}, {
  "kind": "contract"
}, {
  "kind": "earnings_report"
}, {
  "kind": "enforcement_action"
}, {
  "kind": "financial_report"
}, {
  "kind": "lawsuit"
}, {
  "kind": "news"
}, {
  "kind": "official_statement"
}, {
  "kind": "podcast"
}, {
  "kind": "post"
}, {
  "kind": "press_release"
}, {
  "kind": "research_paper"
}, {
  "kind": "site"
}, {
  "kind": "technical_documentation"
}, {
  "kind": "whitepaper"
}]
```


A report or research publication can be returned when it has been ingested as
one of those types. V1 does not claim a complete report corpus and does not
offer a separate report filter.

The target follows the established provider experience for:

- article search and exact article retrieval;
- latest, headline, and trending feeds;
- source and filter-value discovery;
- optional full-content projection;
- similar-article discovery;
- Story collection, detail, and article-membership routes;
- count and distribution queries; and
- stable collection, detail, and error envelopes.

Beans also retains capabilities supported by its data but not commonly exposed
together by the reviewed providers:

- natural-language semantic retrieval;
- attention-ranked trending results;
- article-level social or forum mentions.

## 2. Verified data capability assessment

### 2.1 Settled identity and source assumptions

This proposal assumes the following work is completed outside this API
renovation:

- every Article has a non-null, unique UUID primary key;
- every Publisher has a non-null, unique UUID primary key;
- an Article's source reference identifies one Publisher UUID;
- the Article UUID is the canonical public Article identity;
- the Publisher UUID is the canonical public Source identity; and
- the Article URL remains a unique public attribute and exact lookup value.

The proposal therefore does not treat Article identity, Publisher identity, or
the Article-to-Publisher reference as capability gaps.

Story identity is separate. The current data can derive groups of related
articles, but the grouping key is not yet a durable public UUID. The Story
routes in this proposal require persistent Story UUIDs and persisted membership
before publication. A Story UUID must not change when its representative
article or ranking changes.

### 2.2 Data profile

A bounded, read-only assessment of production data on 2026-08-17 found:

| Capability | Verified coverage | Target interpretation |
|---|---:|---|
| News Articles | 1,744,608 | V1 content type. |
| Blog Articles | 162,234 | V1 content type. |
| In-scope Articles | 1,906,842 | Searchable V1 corpus. |
| Titles | 1,906,842 | Required Article field. |
| Authors | 1,154,904 | Optional field; sufficient for an author filter. |
| Summaries | 1,887,878 | Optional Article field. |
| Publicly eligible stored content | 785,567 | Optional; still subject to source availability and rights. |
| Semantic-search coverage | 1,082,494 | Useful but incomplete; semantic queries do not cover every Article. |
| Categories | 1,877,108 | Optional filter labels. |
| Sentiments | 1,877,583 | Optional categorical filter labels. |
| Regions | 1,295,193 | Optional extracted region labels, not structured places. |
| Entities | 1,419,908 | Optional extracted strings, not canonical entity profiles. |
| Publishers | 27,188 | Supports Source collection and detail routes. |
| Publisher names | 19,236 | Nullable Source field. |
| Publisher descriptions | 21,636 | Nullable Source detail field. |
| Publisher favicons | 24,565 | Nullable Source detail field. |
| Publisher feeds | 6,748 | Nullable Source detail field. |
| In-scope related-Article pairs | 5,820,256 across 382,456 anchor Articles | Supports similar coverage and Story membership. |
| In-scope mention observations | 1,687,551 across 136,943 Articles | Supports Article mention retrieval. |
| In-scope derived groups | 189,961 | Supports Story renovation after persistent identity is added. |
| Multi-Article groups | 46,288 containing 235,206 Article memberships | Supports useful cross-publisher Story retrieval. |
| Publication range | 2021-08-18 through 2026-08-16 | Supports current and historical search. |

Random samples confirmed that Articles can contain long-form news, blog posts,
publisher analysis, institutional reports, summaries, images, authors, and
multiple enrichment labels. Publisher samples confirmed that name,
description, favicon, and feed coverage varies substantially.

Category, entity, region, and tag values have inconsistent source casing and
spacing. The public contract normalizes these values to lowercase snake_case,
for example united_states and example_company. They remain filter labels, not
canonical taxonomies, geographic identifiers, or knowledge-graph entities.

### 2.3 Market capability coverage after renovation

| Market-report capability | Beans decision | State |
|---|---|---|
| Broad Article search | Support q plus exact Article, source, author, type, taxonomy, and date filters. | V1 |
| Exact Article retrieval | Support one Article by UUID and exact IDs or URLs on search. | V1 |
| News and Blog selection | Support content_type=news or blog. | V1 |
| Latest feed | Keep a dedicated recent-publication feed. | V1 |
| Top headlines | Keep a fixed recent-window attention feed. | V1 |
| Trending feed | Keep an attention-ranked feed with explicit metrics. | V1, Beans extension |
| Full-content control | Support full_content as projection only; return null when unavailable. | V1, partial coverage |
| Source discovery and detail | Support browse/search plus UUID detail. | V1 |
| Category, entity, region, sentiment, and tag discovery | Expose accepted normalized filter values. | V1 |
| Author filtering | Filter by stored byline strings. | V1, optional coverage |
| Include and exclude source filters | Support Publisher UUID and domain allow/deny lists. | V1 |
| Similar or related Articles | Expose known related publisher coverage under an Article. | V1 |
| Story clusters | Expose /stories, Story detail, and Story Article membership. | Renovation prerequisite: persistent Story UUID and membership |
| Semantic retrieval | Use q for natural-language relevance retrieval. | V1, partial corpus coverage |
| Count and distributions | Expose bounded exact counts and approved groupings. | V1 after query-plan review |
| Social/forum mentions | Expose observations attached to a known Article. | V1, Beans extension |
| Boolean or field-targeted lexical search | Not reliably supported by the current search model. | Future gap |
| Language filtering | No reliable Article language value. | Future gap |
| Structured geography and radius search | Region labels do not provide country codes, coordinates, or distance semantics. | Future gap |
| Canonical people, companies, concepts, and tickers | Entity values lack stable identity, type, alias, and ticker metadata. | Future gap |
| Numeric sentiment and confidence | Current sentiment values are categorical and have no calibrated score. | Future gap |
| Story history | No durable Story change history exists. | Future gap |
| Incremental continuation stream | No public continuation contract for newly indexed Articles exists. | Future gap |
| Generated cross-Article summary | Article summaries exist, but Beans does not generate a result-set or Story summary. | Future gap |
| Caller-supplied extraction or analysis | Beans searches indexed content; it does not analyze arbitrary caller URLs or text. | Future gap |
| Webhook or WebSocket delivery | No push-delivery contract exists. | Future gap |
| PR content type | V1 deliberately supports news and blog only. | Future product decision |

## 3. Industry-aligned contract
Research Doc: [QUERY_PARAM_NAMES.md](QUERY_PARAM_NAMES.md)

### 3.1 Contract decisions

| Area | Beans V1 decision | Industry basis |
|---|---|---|
| Search route | GET /articles/search | World News /search-news, GNews /search, Currents /search, and comparable all/everything routes. |
| Unstructured text query | q only | `q` is the uniform search-text parameter across collection and discovery routes. Its matching mode is route-specific and documented; it may use full-text, fuzzy, prefix, contains, semantic, or hybrid matching. |
| Fuzzy label filter | tags only | `tags` performs tolerant matching over normalized category, region, and entity labels. It does not promise fuzzy matching over Article titles or bodies. |
| Semantic score threshold | score_threshold | Perigon vector search exposes scoreThreshold; Beans uses snake_case and explicit minimum-score semantics. |
| Article detail | GET /articles/{id} | TheNewsAPI UUID lookup and other known-ID or exact-link retrieval routes. |
| Feeds | /articles/latest, /top-headlines, and /articles/trending | NewsData.io/Currents latest plus GNews, NewsAPI.org, and TheNewsAPI headline/top feeds. |
| Story routes | /stories, /stories/{story_id}, /stories/{story_id}/articles | GDELT Story collection/detail/articles and Perigon Story collection plus clusterId membership. |
| Similar Articles | GET /articles/{id}/similar | Article-subresource form of TheNewsAPI similar-by-UUID. |
| Content type | content_type=news,blog,... | NewsAPI.ai's native data-type selection, narrowed to Beans data. |
| Publication range | from and to as YYYY-MM-DD | Perigon, finlight, GNews, and NewsAPI.org use from/to-style ranges. |
| Source filters | sources, exclude_sources, domains, and exclude_domains | Source/domain include and exclude patterns used by Perigon, TheNewsAPI, finlight, and NewsAPI.org. |
| Author filter | authors | World News and Currents explicitly support author filtering. |
| Content projection | full_content | NewsData.io uses full_content; other providers expose equivalent body controls. |
| Pagination | cursor and limit; response next_cursor | GDELT Cloud and Currents V2 use opaque continuation tokens. This is the deliberate V1 choice for a changing news and blog corpus with API-determined ordering. |
| Discovery resources | /categories, /entities, /regions, /sentiments, and /tags | The market report has no shared discovery prefix. Direct plural resources are the neutral industry-aligned choice; /available is Currents-specific. |
| Query-name vocabulary | `q` for unstructured search text across Articles, Stories, Sources, and discovery routes; `tags` for structured label filtering | `QUERY_PARAM_NAMES.md` standardizes `q` across Beans and Espresso. The route contract must document searchable fields and whether matching is lexical, fuzzy, prefix-based, semantic, or hybrid. |
| Exact identity filters | `ids` and `urls` are accepted on Article search only; latest and trending are feed routes and omit them. | Industry services generally separate filtered search, ranked/current feeds, and exact retrieval. World News uses `/retrieve-news?ids=...`, TheNewsAPI uses UUID detail lookup, and GNews, NewsAPI.org, and Currents do not expose article IDs or URLs on their headline/latest feed routes. |
| Ordering | No public sort parameter | Each V1 route defines a deterministic default order. |
| Output | JSON | All reviewed provider routes use JSON. |
| Status fields | No success or status property | HTTP status codes carry request outcome; stable payload fields add more value. |

Provider-specific alternatives such as search, text, keyword, date_start,
date_end, offset, page, pageSize, page_size, page_number, size, number,
articlesCount, and nextPage are comparison references only. They are not Beans
V1 aliases.

### 3.2 Common collection parameters

| Parameter | Type | Contract |
|---|---|---|
| limit | Integer 1-100; default 20 | Maximum items returned in one response. |
| cursor | Opaque string; omitted on the first request | Continuation token returned by `pagination.next_cursor`. The client must preserve it unchanged. |

`cursor` and `next_cursor` follow the explicit cursor convention used by GDELT
Cloud and Currents V2. An opaque continuation is deliberately not named `page`:
clients generally expect `page` to be an integer. `limit` is common across the
reviewed APIs. Beans does not accept `offset`, `page`, `pageSize`, `size`, or
other pagination aliases.

Article collections can accept the following where the route table permits:

| Parameter | Type | Contract |
|---|---|---|
| q | String, maximum 512 characters | Unstructured search text. The route documents whether matching is full-text, fuzzy, prefix, contains, semantic, or hybrid; Boolean syntax and exact phrase behavior are not promised. |
| score_threshold | Number from 0.0 to 1.0 | Optional minimum semantic similarity score. Accepted only when q is present. Omitted means the route applies its normal server-defined relevance behavior. |
| ids | CSV UUIDs, maximum 128 | Exact Article identities. |
| urls | CSV URLs, maximum 128 | Exact Article URLs. |
| content_type | news or blog | Omitted means both types. |
| sources | CSV Publisher UUIDs | Include Articles from any listed Source. |
| exclude_sources | CSV Publisher UUIDs | Exclude Articles from any listed Source. |
| domains | CSV strings | Include Articles from any exact Source domain. |
| exclude_domains | CSV strings | Exclude Articles from any exact Source domain. |
| authors | CSV strings | Case-insensitive exact byline matches. |
| categories | CSV snake_case strings | Match any listed category. |
| exclude_categories | CSV snake_case strings | Exclude Articles containing any listed category. |
| regions | CSV snake_case strings | Match any listed extracted region label. |
| entities | CSV snake_case strings | Match any listed extracted entity label. |
| sentiments | CSV strings | Match any listed categorical sentiment. |
| tags | CSV snake_case strings | Fuzzy match across normalized category, region, and entity labels; all supplied tags must match. This is not fuzzy Article-text matching. |
| from / to | ISO 8601 date, YYYY-MM-DD | Inclusive publication-date bounds. |
| full_content | Boolean; default false | Requests content without changing which Articles match. |

Different filter dimensions are combined with AND. Multiple values within one
include filter are OR. Exclusion filters remove matching Articles.

Date query values are UTC calendar dates. For example, from=2026-02-01 begins
at 00:00:00Z and to=2026-02-01 includes that complete UTC date. V1 accepts only
the ISO 8601 YYYY-MM-DD form for date query parameters. RFC3339 date-time values
are not accepted as query parameters.

### 3.3 Collection, detail, and error responses

Collections use:

~~~json
{
  "pagination": {
    "limit": 20,
    "next_cursor": "opaque-token-or-null"
  },
  "data": []
}
~~~

Rules:

- `pagination.limit` echoes the effective request value.
- `pagination.next_cursor` is an opaque string when another page is available;
  clients send it back unchanged as `cursor`. It is `null` when traversal is
  complete.
- Beans does not return `found`, `returned`, `page`, `next_page`, or an exact
  total in collection responses. The number of records in the current response
  is `data.length`.
- Empty collections return HTTP 200 with `data: []` and `next_cursor: null`.
- `meta.as_of` is included when a response depends on changing trend, Story, or
  mention observations. It is omitted for direct Article and Source reads.

The cursor is bound to the result ordering, filters, and effective limit. A
client that changes any of these begins a new traversal without a cursor. This
avoids the skip and duplicate risk of integer pagination over a changing corpus
while V1 has no public `sort` parameter.

Detail routes use:

~~~json
{
  "data": {}
}
~~~

Errors use an HTTP error status and:

~~~json
{
  "error": {
    "code": "invalid_request",
    "message": "The limit must be between 1 and 100."
  }
}
~~~

The contract does not use success, success: false, or a duplicate HTTP status
field in the payload.

### 3.4 Article payload

Article collection and detail routes use the same field names. Detail retrieval
can add links, but it does not use a different Article vocabulary.

| Field | Rule |
|---|---|
| id | Required Article UUID and canonical public identity. |
| url | Required canonical Article URL. |
| content_type | Required; news or blog. |
| title | Required. |
| summary | Nullable. |
| content | Included only when full_content=true; string when available, otherwise null. |
| author | Nullable byline. |
| image_url | Nullable. |
| published_at | Required ISO 8601 UTC timestamp. |
| source | Required nested Source summary. |
| story_id | Nullable stable Story UUID after Story publication. |
| categories, regions, entities, tags | Arrays of lowercase snake_case values. |
| sentiments | Array of categorical values. |

~~~json
{
  "id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
  "url": "https://publisher.example/article",
  "content_type": "news",
  "title": "Article headline",
  "summary": "Article summary",
  "author": "Author Name",
  "image_url": "https://publisher.example/image.jpg",
  "published_at": "2026-08-17T12:00:00Z",
  "source": {
    "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
    "domain": "publisher.example",
    "name": "Publisher",
    "url": "https://publisher.example"
  },
  "story_id": "330e5a9f-ee65-4781-b8ad-409a8bb33ab4",
  "categories": ["technology"],
  "regions": ["united_states"],
  "entities": ["example_company"],
  "sentiments": ["analytical"],
  "tags": ["technology", "united_states", "example_company"]
}
~~~

full_content=true changes projection only. An unavailable or restricted body
returns content: null; it does not remove the Article from the result.

### 3.5 Trend payload

Trending items add a nested trend object to the Article:

~~~json
{
  "id": "725a7f65-f7eb-4643-b7af-78e8980217e9",
  "url": "https://publisher.example/article",
  "content_type": "news",
  "title": "Article headline",
  "summary": "Article summary",
  "content": null,
  "author": "Author Name",
  "image_url": null,
  "published_at": "2026-08-17T12:00:00Z",
  "source": {
    "id": "a2f08211-1796-43df-8988-c35a2923a7a0",
    "domain": "publisher.example",
    "name": "Publisher",
    "url": "https://publisher.example"
  },
  "story_id": null,
  "categories": [],
  "regions": [],
  "entities": [],
  "sentiments": [],
  "tags": [],
  "trend": {
    "likes": 120,
    "comments": 42,
    "mentions": 8,
    "audiences": 25000,
    "related": 4,
    "trend_score": 91.7
  }
}
~~~

Missing engagement values are null, not fabricated as zero. mentions is
the number of observed external mentions, not a publisher-reported share total.

### 3.5.1 Article route response profiles

The Article data item uses one canonical field set across B01, B03, B04, B05, and B02. The route-specific differences are:

| Route | Response data item | Route-specific response behavior |
|---|---|---|
| B01 — `/articles/search` | Canonical Article object | Collection envelope; `content` follows the `full_content` projection rule. |
| B03 — `/articles/latest` | Canonical Article object | Collection envelope; recent-publication ordering; no `trend` object. |
| B04 — `/top-headlines` | Canonical Article object | Collection envelope; server-defined headline ordering; no `trend` object. |
| B05 — `/articles/trending` | Canonical Article object plus `trend` | Collection envelope; trend-score ordering; `meta.as_of` identifies the changing observation time. |
| B02 — `/articles/{id}` | One canonical Article object | Detail envelope; may add `links.similar`, `links.mentions`, and optional `links.story`; no `pagination` or `meta.as_of`. |

The canonical Article fields are `id`, `url`, `content_type`, `title`,
`summary`, `content`, `author`, `image_url`, `published_at`, `source`,
`story_id`, `categories`, `regions`, `entities`, `sentiments`, and `tags`.
Trending does not replace or omit these fields; it adds `trend`. The `content`
field follows the same `full_content` rule on every Article route.

B05 uses the same collection envelope with conditional trend metadata:

~~~json
{
  "pagination": {
    "limit": 20,
    "next_cursor": null
  },
  "meta": {
    "as_of": "2026-08-17T15:00:00Z"
  },
  "data": []
}
~~~

B01, B03, and B04 omit `meta.as_of` unless a future implementation makes their
results depend on changing trend, Story, or mention observations.

### 3.5.2 NewsAPI.ai detail-enrichment comparison

The market report documents a richer Article data model for NewsAPI.ai. Its
known-Article operation can request a body and optional enrichment fields, but
these fields are not necessarily returned for every Article and are not all
detail-only fields.

Documented NewsAPI.ai fields that Beans does not currently expose include:

| Field or group | Meaning | Closest Beans field or status |
|---|---|---|
| `uri` | Provider Article URI | `id` is the Beans canonical identity. |
| `body` | Full Article body | `content`; controlled by `full_content`. |
| `date`, `time`, `dateTime`, `dateTimePub` | Publication and related timestamps | `published_at`; Beans does not expose the additional timestamp set. |
| `lang` | Provider language code | Future gap. |
| `dataType` | `news`, `blog`, or `pr` classification | `content_type` supports `news` and `blog`; `pr` is out of scope. |
| `isDuplicate`, `duplicateList` | Duplicate status and related Articles | Future gap. |
| `eventUri`, `storyUri` | Provider Event and Story references | `story_id` is the closest Beans field; no Event equivalent. |
| `relevance` | Provider relevance value | Not exposed in Beans V1. |
| `sentiment` | Numeric sentiment score | `sentiments` contains categorical values only. |
| `concepts` | Disambiguated concepts | `entities` and labels are the closest available fields. |
| `authors` | Structured author records | `author` is a nullable byline string. |
| `links`, `videos` | Related links and video references | Future gap. |
| `shares` | Social-sharing metrics | `trend` contains selected attention metrics, not this object. |
| `extractedDates` | Dates extracted from Article text | Future gap. |
| `location` | Extracted geographic information | Future gap; Beans regions are labels, not structured locations. |
| `source.uri`, `source.title` | Provider Source identity and display fields | Beans exposes a normalized Source object with `id`, `domain`, `name`, and `url`. |

NewsAPI.ai supports these fields through `includeArticle` options and body-size
parameters on its Article operations. Beans must not add them to
`/articles/{id}` merely because another provider exposes them. A Beans detail
response uses the same canonical Article fields as collection responses, adds
only the documented `links` object, and applies the same `full_content` rule.

### 3.6 Story payload

A Story is a persistent grouping of related Articles. It is a news-coverage
resource. Its identity remains stable as Articles enter
or leave the group.

The market comparison leads to the following contract:

- GDELT is the closest route and payload model. It uses one Story card for list
  and detail responses, includes a small `top_articles` preview, and exposes the
  complete membership through `/stories/{story_id}/articles`.
- Perigon exposes additional generated summaries, lifecycle timestamps,
  reprint metrics, locations, topics, and taxonomy objects. Beans V1 does not
  return those fields because the current data cannot support their documented
  meanings.
- World News returns an inline cluster containing a `news` array. NewsAPI.ai
  exposes a Story reference on an Article. Neither defines a more suitable
  standalone Story object for Beans.

Beans therefore uses one canonical Story object. The following table is the
complete field allowlist; fields not listed here are not returned.

| Field | B09 Story item | B10 Story detail | Rule |
|---|---|---|---|
| id | Required | Required | Stable Story UUID. It does not change when the representative Article changes. |
| title | Required | Required | Headline selected from the current top Article. It is not a generated editorial title and can change. |
| first_published_at | Required | Required | Earliest `published_at` among current member Articles. |
| last_published_at | Required | Required | Latest `published_at` among current member Articles. |
| articles | Required | Required | Number of currently known member Articles before request filters or pagination. |
| sources | Required | Required | Number of distinct member Sources before request filters or pagination. |
| categories | Required array | Required array | Up to ten most frequent normalized categories across current members. |
| regions | Required array | Required array | Up to ten most frequent normalized regions across current members. |
| entities | Required array | Required array | Up to ten most frequent normalized entities across current members. |
| tags | Required array | Required array | Up to ten most frequent normalized combined labels across current members. |
| top_articles | Required array | Required array | One to three compact Article previews selected by the server-defined Story ranking. |
| links | Not returned | Required | Detail-only object containing only `articles`. |

`categories`, `regions`, `entities`, and `tags` contain lowercase snake_case
strings. These arrays are always present and use `[]` when no value is
available. Frequency ties use lexical order.

Each `top_articles` item contains exactly these fields:

| Field | Rule |
|---|---|
| id | Required Article UUID. |
| url | Required Article URL. |
| title | Required Article title. |
| published_at | Required ISO 8601 UTC timestamp. |
| source | Required compact Source object containing exactly `id`, `domain`, `name`, and `url`. Missing Source display metadata is `null`; the keys are not omitted. |

A `top_articles` item does not contain `content_type`, `summary`, `content`,
`author`, `image_url`, `story_id`, enrichment arrays, `trend`, or `links`.

~~~json
{
  "id": "330e5a9f-ee65-4781-b8ad-409a8bb33ab4",
  "title": "Representative Story headline",
  "first_published_at": "2026-08-17T10:00:00Z",
  "last_published_at": "2026-08-17T14:20:00Z",
  "articles": 12,
  "sources": 9,
  "categories": ["technology"],
  "regions": ["united_states"],
  "entities": ["example_company"],
  "tags": ["technology", "united_states", "example_company"],
  "top_articles": [
    {
      "id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
      "url": "https://publisher.example/article",
      "title": "Article headline",
      "published_at": "2026-08-17T12:00:00Z",
      "source": {
        "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
        "domain": "publisher.example",
        "name": "Publisher",
        "url": "https://publisher.example"
      }
    }
  ]
}
~~~

Story title and top_articles are selected from member Articles. V1 does not
invent a generated editorial summary. articles and sources describe
the full currently known membership; they do not claim complete internet
coverage.

The Story object does not include a Story `url`, `summary`, `description`,
`content`, `key_points`, generated narrative, Story lifecycle timestamps,
rank, significance, confidence, geographic objects, canonical entity objects,
event references, reprint metrics, taxonomy objects, or the full member
Article list. These are provider-specific fields from GDELT or Perigon whose
semantics Beans V1 cannot support. The full member list belongs only in B11.

### 3.7 Source payload

Source collection item:

~~~json
{
  "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
  "domain": "publisher.example",
  "name": "Publisher",
  "url": "https://publisher.example"
}
~~~

Source detail adds nullable metadata:

~~~json
{
  "data": {
    "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
    "domain": "publisher.example",
    "name": "Publisher",
    "url": "https://publisher.example",
    "description": "Publisher description",
    "favicon_url": "https://publisher.example/favicon.ico",
    "rss_feed_url": "https://publisher.example/feed.xml"
  }
}
~~~

### 3.8 Mention payload

~~~json
{
  "url": "https://social.example/post/123",
  "platform": "social.example",
  "forum": "technology",
  "observed_at": "2026-08-17T14:02:00Z",
  "engagement": {
    "likes": 120,
    "comments": 42,
    "audience": 25000
  }
}
~~~

## 4. Target routes

Gateway paths add /beans. Backend routes omit that prefix.

| ID | Route | What it does | Closest market-report route or feature | State |
|---|---|---|---|---|
| B01 | GET /articles/search | Searches News and Blog Articles. | World News /search-news; GNews /search; NewsAPI.ai getArticles; NewsAPI.org /everything; Currents /search. | Existing + Renovation |
| B02 | GET /articles/{id} | Retrieves one Article by UUID. | TheNewsAPI /news/uuid/{uuid}; World News /retrieve-news; NewsAPI.ai getArticle. | Current-data + Renovation |
| B03 | GET /articles/latest | Returns recently published Articles. | NewsData.io /latest; Currents /latest-news. | Existing + Renovation |
| B04 | GET /top-headlines | Returns current headline Articles from a fixed recent window. | GNews and NewsAPI.org /top-headlines; TheNewsAPI /news/headlines. | Existing + Renovation |
| B05 | GET /articles/trending | Returns Articles ranked by measured attention. | Closest to World News /top-news and TheNewsAPI /news/top; trend metrics are a Beans extension. | Existing + Renovation |
| B06 | GET /articles/{id}/similar | Returns known related publisher coverage. | TheNewsAPI /news/similar/{uuid}. | Current-data + Renovation |
| B07 | GET /articles/{id}/mentions | Returns observed social or forum links to an Article. | Beans extension. | Current-data + Renovation |
| B09 | GET /stories | Searches and browses Story groups. | GDELT /stories; Perigon /stories/all. | Current-data + Persistent Story prerequisite |
| B10 | GET /stories/{story_id} | Retrieves one Story with top Article evidence. | GDELT /stories/{story_id}. | Current-data + Persistent Story prerequisite |
| B11 | GET /stories/{story_id}/articles | Paginates all known member Articles. | GDELT /stories/{story_id}/articles; Perigon clusterId Article filtering. | Current-data + Persistent Story prerequisite |
| B12 | GET /sources | Searches and lists Publisher Sources. | World News /search-news-sources; TheNewsAPI /news/sources; finlight /sources; NewsAPI.org /top-headlines/sources. | Existing + Renovation |
| B13 | GET /sources/{id} | Retrieves one Source profile by UUID. | Beans extension over industry Source collections. | Current-data + Renovation |
| B14 | GET /categories | Discovers category filter values. | Direct category resource; Currents exposes the equivalent under /available/categories. | Current-data + Renovation |
| B15 | GET /entities | Discovers extracted entity filter values. | Direct entity resource pattern used by GDELT and supporting metadata APIs. | Current-data + Renovation |
| B16 | GET /regions | Discovers region filter values. | Direct geography resource; comparable to GDELT geography and Currents /available/regions. | Current-data + Renovation |
| B17 | GET /sentiments | Discovers categorical sentiment filter values. | Direct filter-value resource adapted from provider sentiment filters. | Current-data + Renovation |
| B18 | GET /tags | Discovers normalized combined filter labels. | Direct filter-value resource and Beans extension. | Current-data + Renovation |
| B19 | GET /articles/count | Returns a bounded Article count or approved distribution. | NewsData.io news-count and provider total-result fields. | Current-data + Query-plan review |

### 4.1 Article routes

#### B01 — GET /articles/search

Accepted parameters:

| Parameter | Contract |
|---|---|
| q | Optional natural-language relevance query. |
| score_threshold | Optional number from 0.0 to 1.0. Requires q and removes results below the requested semantic similarity score. |
| ids, urls | Optional exact Article filters. |
| content_type | news or blog; omitted means both. |
| sources, exclude_sources, domains, exclude_domains | Source include/exclude filters. |
| authors | Case-insensitive exact byline matches. |
| categories, exclude_categories | Category include/exclude filters. |
| regions, entities, sentiments, tags | Enrichment-label filters. |
| from, to | Inclusive published_at dates. |
| full_content | Optional content projection. |
| cursor, limit | Common cursor pagination. |

The route permits unfiltered browsing and source-only, date-only, author-only,
content-type-only, exact-ID, and exact-URL requests.

V1 does not accept sort. Results use relevance descending when q is present and
published_at descending otherwise. Article UUID is the stable tie-breaker.

#### B02 — GET /articles/{id}

Returns one Article or HTTP 404.

| Parameter | Contract |
|---|---|
| id | Required Article UUID path parameter. |
| full_content | Requests content when publicly available. |

Detail can add:

~~~json
{
  "links": {
    "similar": "/articles/e4270d59-2fd8-4a2f-90ce-b1a669bedb75/similar",
    "mentions": "/articles/e4270d59-2fd8-4a2f-90ce-b1a669bedb75/mentions",
    "story": "/stories/330e5a9f-ee65-4781-b8ad-409a8bb33ab4"
  }
}
~~~

The story link is omitted when the Article has no published Story membership.

#### B03 — GET /articles/latest

Accepts the B01 filters except ids, urls, from, and to. The route always uses
the previous seven UTC dates and orders results by published_at descending, then
Article UUID. q narrows the candidate set but does not replace the feed's
recent-first ordering. score_threshold is accepted only when q is present.

#### B04 — GET /top-headlines

Accepts the B01 filters except ids, urls, from, and to. Its previous-24-hour
publication window is fixed. q and score_threshold narrow the candidate set but
do not replace the server-defined headline ranking. published_at descending and
Article UUID provide deterministic tie-breakers.

#### B05 — GET /articles/trending

Accepts the B01 filters except ids and urls. By default, the route ranks
Articles using attention observations from the previous 24 hours. Optional from
and to define the inclusive, bounded UTC observation window used to calculate
trend_score; they are not merely Article publication filters. q and
score_threshold narrow the candidate set but do not replace trend ordering.
Results use trend_score descending, published_at descending, and then Article
UUID.

Exact identity filters are intentionally not accepted on B03, B04, or B05.
These routes are feed operations: B03 defines a recent-publication window, B04
defines a headline-ranked result set, and B05 defines an attention-ranked result
set. They are not exact retrieval operations.

This follows the market pattern documented in
`NEWS_AND_BLOG_API_MARKET_REPORT.md`: World News separates filtered search from
`/retrieve-news?ids=...`; GNews, NewsAPI.org, and Currents use separate search
and headline/latest routes without article-ID or URL filters; and TheNewsAPI
uses a UUID detail route for known-article retrieval.

Beans therefore uses the following decision criteria:

- `ids` and `urls` are useful on B01 when a caller wants to constrain a general
  search result set to known Articles.
- B02 is the canonical operation when a caller wants one known Article.
- B03 and B05 keep their feed semantics and do not accept `ids` or `urls`.
- `urls` is a convenience exact-match filter, not a second Article identity.
  Article UUIDs remain the canonical identity. URL matching must be exact and
  must not imply fuzzy URL, domain, or alias matching.

This is a deliberate product decision, not an unresolved legacy omission.

#### B06 — GET /articles/{id}/similar

Returns Article items from known related publisher coverage. It does not claim
that every related Article belongs to the same persisted Story.

Accepted parameters:

| Parameter | Contract |
|---|---|
| content_type | Filter returned Articles. |
| sources, exclude_sources, domains, exclude_domains | Filter returned Articles by Source. |
| authors, categories, regions, entities, sentiments, tags | Filter returned Articles. |
| from, to | Returned Article publication bounds. |
| full_content | Optional content projection. |
| cursor, limit | Common cursor pagination. |

Results use published_at descending and Article UUID as the tie-breaker.

#### B07 — GET /articles/{id}/mentions

Returns external posts or forum records observed linking to the Article URL.

| Parameter | Contract |
|---|---|
| platforms | CSV exact normalized platform values. |
| forums | CSV exact forum or community values. |
| from, to | Inclusive observed_at dates, not Article publication dates. |
| cursor, limit | Common cursor pagination. |

Results use observed_at descending and mention URL as the stable tie-breaker.

### 4.2 Story routes

#### B09 — GET /stories

Accepted parameters:

| Parameter | Contract |
|---|---|
| q | Unstructured search text applied to member Articles; matching behavior is documented for the Story route and may include semantic relevance retrieval. |
| content_type | Requires at least one matching member Article of this type. |
| sources, exclude_sources, domains, exclude_domains | Match through member Article Sources. |
| authors, categories, regions, entities, sentiments, tags | Match through member Articles. |
| from, to | Bounds member Article publication dates. |
| cursor, limit | Common cursor pagination. |

Results use latest member publication descending by default. When q is present,
best member relevance is the primary order. Story UUID is the stable
tie-breaker.

The filters determine whether a Story is included. They do not recalculate its
article_count, source_count, label arrays, or top_articles preview. Those fields
always describe the full current Story membership.

The B09 top-level response contains exactly these fields:

| Field | Rule |
|---|---|
| data | Required array of the canonical Story objects defined in Section 3.6. B09 items do not contain `links`. |
| pagination | Required object containing exactly `limit` and `next_cursor`. |
| meta | Required object containing exactly `as_of`. |

The response does not contain `success`, `status`, `results`, `news`,
`numResults`, `found`, `returned`, `page`, or `next_page`.

Response:

~~~json
{
  "pagination": {
    "limit": 20,
    "next_cursor": "opaque-token-or-null"
  },
  "meta": {
    "as_of": "2026-08-17T15:00:00Z"
  },
  "data": [
    {
      "id": "330e5a9f-ee65-4781-b8ad-409a8bb33ab4",
      "title": "Representative Story headline",
      "first_published_at": "2026-08-17T10:00:00Z",
      "last_published_at": "2026-08-17T14:20:00Z",
      "article_count": 12,
      "source_count": 9,
      "categories": ["technology"],
      "regions": ["united_states"],
      "entities": ["example_company"],
      "tags": ["technology", "united_states", "example_company"],
      "top_articles": [
        {
          "id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
          "url": "https://publisher.example/article",
          "title": "Article headline",
          "published_at": "2026-08-17T12:00:00Z",
          "source": {
            "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
            "domain": "publisher.example",
            "name": "Publisher",
            "url": "https://publisher.example"
          }
        }
      ]
    }
  ]
}
~~~

An empty match returns HTTP 200 with `data: []` and `next_cursor: null`.

#### B10 — GET /stories/{story_id}

Returns one Story or HTTP 404. The response contains at most three top_articles,
following the GDELT Story-card pattern. It does not include full Article content.

This route accepts no query parameters. Its data member uses the complete Story
payload from Section 3.6 and adds a link to the paginated membership route:

| Top-level field | Rule |
|---|---|
| data | Required canonical Story object plus the detail-only `links` object. |

B10 has no collection page metadata because it returns one Story. Its
`data.links` object contains exactly `articles`; it does not contain
similar-Article, mention, or Source links.

~~~json
{
  "data": {
    "id": "330e5a9f-ee65-4781-b8ad-409a8bb33ab4",
    "title": "Representative Story headline",
    "first_published_at": "2026-08-17T10:00:00Z",
    "last_published_at": "2026-08-17T14:20:00Z",
    "article_count": 12,
    "source_count": 9,
    "categories": ["technology"],
    "regions": ["united_states"],
    "entities": ["example_company"],
    "tags": ["technology", "united_states", "example_company"],
    "top_articles": [
      {
        "id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
        "url": "https://publisher.example/article",
        "title": "Article headline",
        "published_at": "2026-08-17T12:00:00Z",
        "source": {
          "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
          "domain": "publisher.example",
          "name": "Publisher",
          "url": "https://publisher.example"
        }
      }
    ],
    "links": {
      "articles": "/stories/330e5a9f-ee65-4781-b8ad-409a8bb33ab4/articles"
    }
  }
}
~~~

article_count can be greater than the number of top_articles because
top_articles is only a maximum-three preview. Use B11 for complete paginated
membership.

#### B11 — GET /stories/{story_id}/articles

Returns all currently known member Articles through a paginated Article
collection.

Accepted parameters:

| Parameter | Contract |
|---|---|
| content_type | Filter member Articles. |
| sources, exclude_sources, domains, exclude_domains | Filter member Article Sources. |
| authors, categories, regions, entities, sentiments, tags | Filter member Articles. |
| from, to | Member Article publication bounds. |
| full_content | Optional content projection. |
| cursor, limit | Common cursor pagination. |

Results use published_at descending and Article UUID as the stable tie-breaker.

The B11 top-level response contains exactly these fields:

| Field | Rule |
|---|---|
| data | Required array of member Article objects. |
| pagination | Required object containing exactly `limit` and `next_cursor`. |
| meta | Required object containing exactly `story_id` and `as_of`. |

Each B11 `data` item contains exactly the following Article fields:

| Field | Presence and rule |
|---|---|
| id | Required Article UUID. |
| url | Required canonical Article URL. |
| content_type | Required; `news` or `blog`. |
| title | Required. |
| summary | Required key; string or `null`. |
| content | Omitted when `full_content=false`. When `full_content=true`, required key containing a string or `null`. |
| author | Required key; byline string or `null`. |
| image_url | Required key; URL string or `null`. |
| published_at | Required ISO 8601 UTC timestamp. |
| source | Required compact Source object containing exactly `id`, `domain`, `name`, and `url`. Missing Source display metadata is `null`; the keys are not omitted. |
| story_id | Required and equal to the requested `{story_id}` path value. |
| categories | Required array of lowercase snake_case strings. |
| regions | Required array of lowercase snake_case strings. |
| entities | Required array of lowercase snake_case strings. |
| sentiments | Required array of categorical sentiment strings. |
| tags | Required array of lowercase snake_case strings. |

The enrichment arrays are always present and use `[]` when empty. A B11 Article
item does not contain the Story `title`, Story counts, `top_articles`, Story
`links`, a `trend` object, similar Articles, mentions, or Source-detail fields
such as `description`, `favicon_url`, and `rss_feed_url`.

Filters apply only to the returned member Articles; they do not modify the
Story or its membership.

~~~json
{
  "pagination": {
    "limit": 20,
    "next_cursor": "opaque-token-or-null"
  },
  "meta": {
    "story_id": "330e5a9f-ee65-4781-b8ad-409a8bb33ab4",
    "as_of": "2026-08-17T15:00:00Z"
  },
  "data": [
    {
      "id": "e4270d59-2fd8-4a2f-90ce-b1a669bedb75",
      "url": "https://publisher.example/article",
      "content_type": "news",
      "title": "Article headline",
      "summary": "Article summary",
      "author": "Author Name",
      "image_url": "https://publisher.example/image.jpg",
      "published_at": "2026-08-17T12:00:00Z",
      "source": {
        "id": "0af3e96e-912c-4e51-95b7-e44b53e5e07b",
        "domain": "publisher.example",
        "name": "Publisher",
        "url": "https://publisher.example"
      },
      "story_id": "330e5a9f-ee65-4781-b8ad-409a8bb33ab4",
      "categories": ["technology"],
      "regions": ["united_states"],
      "entities": ["example_company"],
      "sentiments": ["analytical"],
      "tags": ["technology", "united_states", "example_company"]
    }
  ]
}
~~~

If the Story exists but the supplied member filters match nothing, the route
returns HTTP 200 with data: []. If the Story does not exist, it returns HTTP
404. full_content follows the ordinary Article projection rule: when requested,
content is a string where available and null otherwise.

### 4.3 Source routes

#### B12 — GET /sources

| Parameter | Contract |
|---|---|
| q | Case-insensitive fuzzy text search over Publisher name, domain, aliases, and description where available. This is lexical metadata matching, not semantic Article search. |
| ids | CSV Publisher UUIDs. |
| domains | CSV exact domains. |
| cursor, limit | Common cursor pagination. |

An unfiltered request browses all Sources. Results use name ascending, then
domain ascending, then Publisher UUID.

#### B13 — GET /sources/{id}

Returns one Source detail object or HTTP 404. id is the Publisher UUID.

### 4.4 Discovery routes

B14-B18 accept `q`, `cursor`, and `limit`.

| Route | Item payload | Filter populated |
|---|---|---|
| GET /categories | {"value":"technology"} | categories |
| GET /entities | {"value":"example_company"} | entities |
| GET /regions | {"value":"united_states"} | regions |
| GET /sentiments | {"value":"analytical"} | sentiments |
| GET /tags | {"value":"technology"} | tags |

`q` is an optional case-insensitive fuzzy text query over the public value. The route may use prefix, contains, token, or tolerant matching as documented. Category, entity,
region, and tag values are returned as lowercase snake_case exactly as clients
send them to Article filters. Values use ascending lexical order.

No common discovery prefix exists across the reviewed services. Currents uses
/available for a subset of its filter values, while other providers expose
sources, entities, people, companies, topics, and geography as direct resource
collections. Beans therefore uses direct plural resource names.

The parameter name is intentionally uniform across route families:

- `q` means “search this route’s public searchable content.” On Articles and
  Stories it may invoke semantic relevance retrieval; on Sources it searches
  source metadata; on discovery routes it searches the returned canonical values.
- `tags` filters Articles or Stories by tolerant matching against normalized
  labels; it is not a free-text title or body search.

This follows `QUERY_PARAM_NAMES.md` and the market research. The reviewed
services use different provider-specific names, including GDELT’s `search` and
Perigon’s vector-search `prompt`, while PredictHQ and Perigon use `q` for
ordinary text search. Beans normalizes these variations to `q` so users and AI
clients can use the same input field across resources. The route contract must
state the searchable fields and matching behavior explicitly; `q` does not by
itself promise a particular algorithm.

### 4.5 Count route

#### B19 — GET /articles/count

The route accepts the B01 scalar and CSV filters except full_content,
score_threshold, cursor, and limit. from and to are required to keep the count
bounded. q remains disabled until its cost and exact semantics are reviewed.

Without group_by:

~~~json
{
  "data": {
    "count": 42
  },
  "meta": {
    "counted_resource": "article",
    "time_field": "published_at",
    "as_of": "2026-08-17T15:00:00Z"
  }
}
~~~

Approved group_by values are published_day, content_type, domain, category,
region, and sentiment.

~~~json
{
  "data": {
    "group_by": "content_type",
    "buckets": [
      {"key": "news", "count": 39},
      {"key": "blog", "count": 3}
    ]
  },
  "meta": {
    "counted_resource": "article",
    "time_field": "published_at",
    "as_of": "2026-08-17T15:00:00Z"
  }
}
~~~

Buckets for multi-valued fields are not additive because one Article can
contribute to several values.

## 5. Query and behavior rules

1. Public Article queries return only news and blog.
2. Article and Publisher UUIDs are canonical identities. URLs and domains are
   public attributes and exact filter values.
3. q is unstructured search text. Each route documents its searchable fields and matching mode. score_threshold is valid only when that route performs semantic retrieval and excludes results below the requested similarity score.
4. full_content=true changes projection only. Unavailable content is null.
5. from and to filter published_at except on mentions, where they filter
   observed_at.
6. Include-list values are OR within a dimension. Different dimensions are AND.
   Exclusion lists remove a match after the include conditions are applied.
7. V1 has no sort parameter. Every route defines its default order and a stable
   UUID tie-breaker.
8. Story identity is persistent. A changing representative Article must not
   change Story ID.
9. Story filters match through member Articles. Story detail is not an
   Article-detail alias.
10. Source identity and Source display metadata come from the referenced
    Publisher. Optional Source detail values remain nullable.
11. Category, entity, region, and tag values are lowercase snake_case in both
    filters and responses.
12. Pagination uses opaque cursor and limit. Clients send `next_cursor` back
    unchanged as `cursor` and must not parse or construct token values.
13. Empty collections and empty relationship members are [], not null or HTTP
    204.
14. A missing Article, Story, or Source returns HTTP 404 before a subresource
    query runs.

## 6. Renovation delta from the current Beans API

The identity and Article-to-Publisher migrations stated in Section 2.1 are
settled prerequisites owned by another workstream and are not counted here.

| Area | Current behavior | Target correction |
|---|---|---|
| Article scope | Runtime defaults to News and Blog, while types and documentation describe additional kinds. | Publicly define and enforce news and blog only. |
| Response envelope | Collections return bare arrays and can return HTTP 204 when empty. | Return HTTP 200 with TheNewsAPI-aligned meta and data. |
| Pagination | Runtime uses limit and offset. | Publish opaque cursor and limit; return `pagination.next_cursor` and no exact collection total. |
| Article detail | No canonical UUID Article route. | Add GET /articles/{id}. |
| Search eligibility | Search rejects source-only, date-only, type-only, and browse requests. | Permit every documented filter combination and unfiltered browse. |
| Exact retrieval | Search accepts URLs but not canonical UUIDs. | Add ids and retain urls as exact filters. |
| Time | Only from is supported. | Add inclusive to and keep both values as YYYY-MM-DD. |
| Authors | Author is returned but cannot be filtered. | Add authors to Article and Story-member queries. |
| Source exclusions | Only a Source include filter exists. | Add exclude_sources, domains, and exclude_domains. |
| Search controls | The current threshold has an unclear public name. | Publish score_threshold with explicit 0.0-1.0 minimum-similarity semantics and require q. |
| Full content | full_content can remove restricted or missing-content Articles. | Preserve the matched set and return content as string or null. |
| Payload consistency | Search and feeds return different flat field sets. | Use one Article vocabulary and add nested trend only where relevant. |
| Source projection | Source metadata is flattened or omitted depending on route. | Return the normalized nested Source summary. |
| Trend metrics | Flat fields obscure metric meaning and missing values can become zero. | Return accurately named nullable fields under trend. |
| Sources | /sources primarily resolves legacy source strings. | Make it browseable/searchable by Publisher UUID, name, and domain; add detail. |
| Discovery | Category/entity/region routes live under /tags and return bare strings. | Publish direct plural resources with {value} items and normalized envelopes; add sentiments and a top-level tags resource. |
| Similar Articles | Related coverage is available only through URL propagation. | Add GET /articles/{id}/similar. |
| Mentions | Mention observations are available only through URL propagation. | Add GET /articles/{id}/mentions. |
| Stories | A derived grouping exists, but no persistent public Story identity or membership resource exists. | Persist stable Story UUIDs/membership, then add B09-B11. |
| Count | Count support is not a public route. | Add bounded B19 after plan review. |
| Semantic coverage | Natural-language relevance does not cover the complete Article corpus. | Improve coverage or keep the limitation explicit. |
| Label quality | Case, spacing, and snake_case variants coexist. | Normalize category, entity, region, and tag values at the public boundary. |

## 7. V1 publication scope

### 7.1 Request limitations

- V1 has no sort query parameter. Every collection route uses its documented
  server-defined ordering.
- V1 collection routes use `limit` and an opaque `cursor`; clients receive
  `pagination.next_cursor` and send it back unchanged. Beans does not accept
  integer `page` or `offset` parameters.
- Every V1 date query parameter accepts only an ISO 8601 calendar date in
  YYYY-MM-DD form, for example 2026-02-01.
- RFC3339 date-time values are not accepted for date query parameters.
- Response fields that represent an instant, such as published_at, observed_at,
  and meta.as_of, remain ISO 8601 UTC timestamps.

### 7.2 Included routes

The renovated V1 target includes:

- B01-B07 Article, feed, similar, and mention routes;
- B12-B13 Source routes;
- B14-B18 discovery routes; and
- the normalized Article, Source, trend, collection, detail, and error payloads.

B09-B11 are part of the V1 design but must not be published until Story UUIDs
and memberships are persistent. The API should publish the industry-aligned
/stories resource family rather than a made-up
/articles/{id}/cluster endpoint.

B19 remains gated on bounded filters, reviewed query plans, and response tests.

Existing /tags/* routes can remain documented compatibility routes during a
versioned migration. The canonical and deployed Article paths overlap, so bare
arrays and the target envelope cannot coexist through undocumented behavior;
the change requires an explicit version or a coordinated breaking release.

## 8. Provider-backed future capability gaps

Every item below is explicitly present in at least one service documented in
the market report. This table does not include speculative capabilities.

| Capability | Beans gap | Market-report evidence |
|---|---|---|
| Boolean and field-targeted lexical search | q does not promise Boolean grammar, phrases, exclusions, or title/body field selection. | World News, TheNewsAPI, GNews, NewsAPI.ai, NewsAPI.org, and NewsData.io expose Boolean or field-specific search controls. |
| Language filtering | No reliable Article language value is available. | World News, TheNewsAPI, finlight, GNews, NewsAPI.ai, NewsData.io, NewsAPI.org, and Currents expose language controls. |
| Structured country and geography filtering | Region labels are not country codes or structured Article/Source locations. | World News, GDELT, Perigon, finlight, NewsData.io, NewsAPI.org, and Currents expose country or geography filters. |
| Radius-based location filtering | No coordinates or place geometry are available. | World News documents location-filter radius discovery. |
| Canonical people, companies, concepts, and tickers | Entity strings lack stable identifiers, types, aliases, and market identifiers. | Perigon exposes people and companies; finlight exposes tickers; NewsAPI.ai exposes conceptUri. |
| Numeric sentiment and confidence | Sentiment is categorical and has no calibrated numeric score or confidence. | World News exposes sentiment ranges; finlight exposes sentiment and confidence. |
| Journalist profiles | Beans has optional bylines but no author identity or profile resource. | Perigon exposes /journalists and journalist detail. |
| Article revision timestamps | Beans does not expose reliable publisher revision history. | finlight returns publication, indexing, and revision timestamps. |
| Story history | Beans has no public Story change log. | Perigon exposes /stories/history. |
| Incremental Article stream | Beans has no continuation values for newly indexed News and Blogs. | NewsAPI.ai exposes minuteStreamArticles with separate continuation values. |
| Cross-Article or Story summarization | Beans returns stored Article summaries but does not generate summaries over a result set or Story. | Perigon exposes /articles/summarize. |
| Caller-supplied extraction and analysis | Beans does not extract arbitrary URLs or analyze caller-provided text. | World News exposes /extract-news; NewsAPI.ai exposes extraction, annotation, categorization, similarity, and sentiment operations. |
| Webhook or WebSocket delivery | Beans has no subscription or push contract. | finlight documents webhook and WebSocket delivery. |
| PR content type | V1 has no PR discriminator. | NewsAPI.ai dataType explicitly supports pr. |
| Rich Source classification | Source records do not reliably provide country, language, paywall, rank, traffic, or topic metadata. | Perigon and finlight expose richer Source metadata. |

## 9. Publication order

1. Complete the settled Article/Publisher identity and Source-reference
   prerequisite.
2. Enforce the News and Blog product boundary and normalize Article, Source,
   error, and empty-collection payloads.
3. Add UUID Article and Source detail routes plus exact ids and urls filters.
4. Add to, author and Source-exclusion filters, normalized snake_case label
   values, and deterministic route ordering.
5. Correct full_content so it changes projection rather than result membership.
6. Add similar-Article and mentions routes.
7. Introduce the canonical collection envelope through an explicit version or
   coordinated breaking release.
8. Persist stable Story UUIDs and memberships, then publish /stories,
   /stories/{story_id}, and /stories/{story_id}/articles.
9. Publish bounded Article counts after query-plan and response-contract review.
10. Add only provider-backed future capabilities after the required data and
    product contracts exist.
