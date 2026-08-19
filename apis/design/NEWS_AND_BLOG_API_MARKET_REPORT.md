# Industry News and Blog API Route Reference

Status: External API reference
Reviewed: 2026-08-17

This document lists the documented public routes, principal query parameters, response payloads, and capabilities of multi-source news and blog APIs. It is a route and contract reference for comparing Beans with established news-intelligence services.

## Scope

This reference covers article search, current feeds, source/taxonomy discovery, article detail lookup, full-content access, enrichment, related-article discovery, and provider analysis endpoints where they materially affect a Beans-style News and Blogs API. It excludes event-only APIs, publisher-only APIs, account, billing, webhook, saved-search, administration, and unrelated product routes. Provider paths are shown without a product-specific gateway prefix.

A searchable body, a response snippet, and full article text are materially different. "Full content" below means that provider documentation explicitly says a full body can be returned. It does not mean Markdown, structured content, a canonical publisher copy, or universal source coverage.

## Common Conventions

| Provider | Primary collection envelope | Pagination / continuation | Content contract |
| --- | --- | --- | --- |
| World News API | offset, number, available, news | offset and number | text is documented as full article text; no Markdown or structured-body guarantee. |
| GDELT Cloud | success, data, pagination | Cursor: limit and cursor | Story cards include only top article evidence; story-article retrieval is not documented as full body. |
| Perigon | status, numResults, articles or results | Zero-based page and size | Article content is source-, license-, and plan-dependent; no separate blog type is documented. |
| TheNewsAPI | meta, data | page and limit; meta returns found and returned | Metadata plus a 60-character snippet, not article body. |
| finlight | status, page, pageSize, articles | page and pageSize | List summaries; by-link detail can include content where entitled. |
| GNews | totalArticles, articles | page with max result count; no cursor in the documented envelope | Free content is truncated; paid plans document full content. |
| NewsAPI.ai / Event Registry | Provider-configurable article result; item schema documented | articlesPage and articlesCount | body length is controlled with articleBodyLen; -1 requests full body. |
| NewsData.io | status, results, nextPage | Opaque nextPage passed back as page | full_content=1 requests source-available content. |
| NewsAPI.org | status, totalResults, articles | page and pageSize | content is always truncated to 200 characters. |
| Currents | status, news, page, next_cursor (V2 search) | V1 page_number/page_size; V2 cursor continuation | Current example contains no article body field. |

All providers return JSON, but their identifiers, date formats, content rights, and response envelopes are incompatible. A client should normalize only the stable public fields it needs rather than assuming that fields named content or summary have equivalent meaning.

## Beans Baseline

| Capability | Current Beans surface |
| --- | --- |
| Feeds | GET /articles/latest, /articles/trending, and /articles/top-headlines. |
| Search | GET /articles/search accepts q, urls, content_type (news or blog), categories, regions, entities, tags, sources, from, full_content, limit, and offset. |
| Discovery | GET /tags/categories, /tags/entities, /tags/regions, and /sources expose exact filter values. |
| Cross-source context | GET and POST /articles/propagation accept article URLs and return coverage and mention data. |
| Payload | Article objects include title, URL, author, source, summary, optional content, image, publication date, category/region/entity/sentiment labels, and tags. Search/trend results can add engagement and trend fields. |
| Pagination | List routes use limit and offset and return arrays, not a total-count envelope. |

## 1. Capability and Commercial Snapshot

| Provider | Full-content position | Closest Beans comparison |
| --- | --- | --- |
| [World News API](https://worldnewsapi.com/docs/) | Search and retrieve documentation call text the full article text. No Markdown or structured-body guarantee. | General search plus entities, sentiment, source country, geographic filtering, and clustered top news. |
| [GDELT Cloud](https://docs.gdeltcloud.com/api-reference/v2) | Story list/detail responses expose only top-three article evidence; the story-articles route paginates the source set. | Story cluster and article-evidence service, not a general raw-article or explicit blog-search API. |
| [Perigon](https://docs.perigon.io/docs/getting-started) | Article responses include content when available; access depends on source, licensing, and plan. | Direct article/news search, story clusters, semantic retrieval, source discovery, and enrichment; no dedicated blog type documented. |
| [TheNewsAPI](https://www.thenewsapi.com/documentation) | No full body. snippet is the first 60 characters of the article body. | Lightweight article cards, total counts, source catalog, and related-article route. |
| [finlight](https://docs.finlight.me/en/v2/rest-endpoints/) | Detail lookup can return content; availability is source/tier-dependent and its includeContent switch is deprecated. | Financial specialist: companies/tickers, sentiment confidence, revision updates, WebSockets/webhooks. |
| [NewsDataHub](https://newsdatahub.com/) | Not applicable; shut down. | Exclude from a current vendor comparison. |
| [GNews](https://docs.gnews.io/) | content is truncated on Free; paid plans include full content. | Compact paid full-text search plus ranked headlines. |
| [NewsAPI.ai / Event Registry](https://www.newsapi.ai/documentation) | body is available; articleBodyLen=-1 requests the full body. | Most complete analysis comparison: native news/blog/pr selection, concepts, categories, duplicates, social metrics, and text analytics. |
| [NewsData.io](https://newsdata.io/) | full_content=1 requests full content where available; paid docs state full-content access. | Continuation token, AI tags/summary/sentiment, count routes, and credits. |
| [NewsAPI.org](https://newsapi.org/docs) | No full body. content is truncated to 200 characters. | Commodity search/headlines/sources baseline. |
| [Currents](https://currentsapi.services/en/docs/endpoint) | Current response docs list headline metadata; paid pricing says available article-content fields without guaranteeing full body. | Clean latest-feed/search split and taxonomy discovery. |

## 2. Content and Pricing Scope

| Provider | Content behavior | Price/limit snapshot reviewed 2026-08-17 |
| --- | --- | --- |
| World News API | Documentation says search/retrieve return full text. It does not substantiate an RSS-only ingestion guarantee, nor Markdown/structured text. | Journalist: $379/month, 5,000 points/day, then $0.0025/point. Search costs 1 point plus 0.01/result. |
| GDELT Cloud | Story endpoints expose article evidence and a paginated story-source list, not a documented full-body article contract. | Commercial pricing was not reviewed in this update; the current public docs require bearer API authentication. |
| Perigon | Article content is a standard response field but availability depends on source, licensing, and provider enrichment. | Commercial pricing was not reviewed in this update; plan determines endpoint and enrichment access. |
| TheNewsAPI | title, meta description, keywords, and 60-character snippet only. | Standard: $49/month, 10,000 requests/day, 100 articles/request. |
| finlight | List responses contain summaries; by-link detail can include content. | Pro Standard: $99/month, 50,000 requests/month; sentiment, entities/tickers, historical data. |
| NewsDataHub | Service shut down on 2026-06-05. | Former $39.99 plan is not actionable. |
| GNews | content truncated on Free, full content on paid tiers. | Essential: EUR 49.99/month, 1,000 requests/day, 25 articles/request. |
| NewsAPI.ai | Full body plus optional analysis. | 5K plan: $90/month. A recent article search uses 1 token and returns up to 100 articles. |
| NewsData.io | full_content=1 where content is available. | Provider blog lists Basic at $199.99/month, 20,000 credits/month, 50 articles/credit. This conflicts with the initial $149.99 note; verify the client-rendered price page before purchase. |
| NewsAPI.org | content is always truncated to 200 characters; no plan offers full article text. | Business: $449/month, 250,000 requests/month. |
| Currents | No body field in current documented example. | Builder: $69/month, 75,000 requests/month, 50 results/request. |

## 3. Provider Route Reference

### World News API

Official references: [search](https://worldnewsapi.com/docs/search-news/), [retrieve](https://worldnewsapi.com/docs/retrieve-news/), [top news](https://worldnewsapi.com/docs/top-news/), [pricing](https://worldnewsapi.com/pricing).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /search-news | text, text-match-indexes, source-country, language, min-sentiment, max-sentiment, publish dates, news-sources, authors, categories, entities, location-filter, sort, offset, number | Broad article search with phrases, AND, OR, parentheses, and - exclusions. Supports entity, radius, source, category, country, and sentiment filters. |
| GET /retrieve-news | ids | Retrieves selected article IDs from search/top-news. |
| GET /top-news | source-country, language, date, headline/cluster-size controls | Country/language top stories clustered across sources; article count drives rank. |
| GET /search-news-sources | name | Discovers monitored source domains. |
| POST /extract-news | Provider-documented URL/page-content inputs | Extracts an article from supplied web content separately from its indexed corpus. |

#### Query Parameter Semantics

- Text matching: text supports keywords, phrases, Boolean AND/OR, parentheses, and exclusions. text-match-indexes controls which indexed text fields are searched.
- Time and ordering: publish dates bound the publication window; sort selects the ranking; offset and number control the result window.
- Source and classification: source-country, language, news-sources, authors, and categories narrow the provider corpus.
- Enrichment and geography: entities, min-sentiment, max-sentiment, and location-filter support entity, sentiment, and radius-based discovery.
- Top-news is a separate cluster feed. Its source-country, language, date, headline, and cluster-size controls determine the ranked cluster set.

Search returns {offset, number, available, news}. A news item contains id, title, text, summary, URL, image/video URLs, publish date, authors, category, language, source country, and sentiment. top-news returns clusters containing a news array.

#### Response Payload

Live documentation example for GET /retrieve-news:

~~~json
{
  "news": [
    {
      "id": 2352,
      "title": "Article headline",
      "text": "Full article text",
      "summary": "Article summary",
      "url": "https://publisher.example/article",
      "image": "https://publisher.example/image.jpg",
      "video": null,
      "source_country": "mx",
      "sentiment": -0.449
    }
  ]
}
~~~

Search adds the offset, number, and available counters around a news array. The search documentation also lists publication date, authors, category, language, and source-country fields on article objects.

#### Beans Comparison

The closest unified search/enrichment competitor. Its main additions are geo radius, source-country, sentiment-range filters, and clustered stories. Beans has an explicit news|blog filter and propagation.

### GDELT Cloud

Official references: [API v2](https://docs.gdeltcloud.com/api-reference/v2), [documentation guide](https://docs.gdeltcloud.com/API_DOCUMENTATION_GUIDE).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /api/v2/stories | date_start, date_end, country, region, continent, admin1, category, event_category, subcategory, domain, has_events, has_fatalities, confidence_profile, article_count_min, article_count_max, search, sort, limit, cursor | Searches generated, clustered news stories and returns a small inline article-evidence preview. |
| GET /api/v2/stories/{story_id} | story_id path parameter | Fetches one story cluster with its article-count, linked-event, entity, geography, and top-article evidence. |
| GET /api/v2/stories/{story_id}/articles | story_id path parameter, limit, cursor | Paginates the full source-article list attached to a known story. |
| GET /api/v2/stories/summary | Story filters plus group_by | Returns aggregate story counts and metrics rather than article records. |

#### Query Parameter Semantics

- Text and time: search performs story-level retrieval; date_start and date_end bound the story window.
- Geography and classification: country, region, continent, admin1, category, event_category, subcategory, and domain filter a story through its generated context and linked events.
- Evidence and confidence: has_events, has_fatalities, article_count_min, article_count_max, confidence_profile, and sort control analytical scope and ranking.
- Pagination: limit and cursor use a cursor envelope. Story list/detail responses return only top_articles as an inline preview; the story-articles route is the complete evidence traversal.
- Scope limit: the reviewed v2 surface has no general all-article search route and no explicit blog type or blog filter.

#### Response Payload

Live v2 story-card response shape:

~~~json
{
  "success": true,
  "data": [
    {
      "id": "story-id",
      "url": "https://gdeltcloud.com/story/example-story",
      "title": "Example story",
      "story_date": "2026-04-17",
      "category": "infrastructure",
      "subcategory": null,
      "geo": {
        "country": "Japan",
        "region": "East Asia",
        "continent": "Asia",
        "admin1": "Tokyo"
      },
      "metrics": {
        "significance": 0.82,
        "article_count": 12,
        "linked_event_count": 2
      },
      "entity_refs": [
        { "id": "Tokyo", "name": "Tokyo", "type": "LOCATION" }
      ],
      "top_articles": [
        {
          "url": "https://publisher.example/article",
          "title": "Example article",
          "domain": "publisher.example",
          "rank": 1
        }
      ]
    }
  ],
  "pagination": {
    "limit": 25,
    "cursor": null,
    "next_cursor": "next-cursor"
  }
}
~~~

The article preview is limited to the top three source articles. GET /api/v2/stories/{story_id}/articles is the documented route for the complete paginated evidence set. The public v2 documentation does not describe a guaranteed article-body field for this route.

#### Beans Comparison

GDELT Cloud is relevant for clustered narratives, article-evidence traversal, normalized geography, and cursor pagination. It is not a direct replacement for Beans article search because its public v2 news surface is story-first, links article evidence to generated event context, and does not document native blog selection or a general raw-article corpus.

### Perigon

Official references: [getting started](https://docs.perigon.io/docs/getting-started), [article data](https://docs.perigon.io/docs/article-data), [story data](https://docs.perigon.io/docs/story-data), [vector news](https://docs.perigon.io/docs/vector-endpoint), [sources](https://docs.perigon.io/docs/searching-sources).

Perigon documentation uses both /v1/articles/all and /v1/all for the article-search surface, and both /v1/sources and /v1/sources/all for source discovery. Confirm the current account's endpoint reference before implementation.

#### Route Inventory

| Route | Main parameters or body | Capability |
| --- | --- | --- |
| GET /v1/articles/all | q, title, desc, content, url, linkTo, from, to, source/excludeSource, sourceGroup/excludeSourceGroup, category, topic, taxonomy, medium, label, language, location, entity, sentiment, sortBy, page, size | Searches and filters individual news articles with content and enrichment fields where entitled. |
| GET /v1/stories/all | q, name, source, clusterId, category, topic, taxonomy, time ranges, cluster-size filters, sortBy, state, showDuplicates, page, size | Searches narrative clusters of related articles and reports aggregate story metrics. |
| GET /v1/stories/history | Story and time filters, page, size | Returns story-change records rather than a current article list. |
| POST /v1/vector/news/all | JSON body with prompt, reprint/source/category/topic/date filters, page, size, scoreThreshold | Performs semantic retrieval over the real-time news corpus. |
| POST /v1/articles/summarize | JSON body with article filters, prompt, maxArticleCount, method, summarizeFields | Generates a summary over matching articles or story clusters. |
| GET /v1/sources or /v1/sources/all | domain, name, source group, geography, traffic, category, topic, paywall, sortBy, page, size | Searches publisher/source metadata. |
| GET /v1/journalists, /v1/people/all, /v1/companies/all, /v1/topics/all | Route-specific identity and pagination filters | Discovers author, entity, company, and topic records used to enrich news results. |

#### Query Parameter Semantics

- Text: q, title, desc, content, name, and route-specific terms support article or story matching. The provider documents advanced Boolean search on relevant search surfaces.
- Time: from/to filter article publication time. addDate, refreshDate, initialized, and updated ranges expose ingestion or story lifecycle windows.
- Classification and source: category, topic, taxonomy, medium, label, source, excludeSource, sourceGroup, excludeSourceGroup, domain, and paywall define the publisher and editorial corpus.
- Geography and entities: article/source location filters, people, companies, entity fields, and sentiment filters support enriched retrieval.
- Pagination: page is zero-based; the reviewed documentation states a default size of 10, a maximum of 100, and a 10,000-result pagination ceiling.
- Scope limit: Perigon describes news and web-content retrieval, but the reviewed public documentation does not define a distinct blog data type or a blog-only filter.

#### Response Payload

Representative article-search response:

~~~json
{
  "status": 200,
  "numResults": 1,
  "articles": [
    {
      "url": "https://publisher.example/article",
      "articleId": "article-id",
      "clusterId": "story-id",
      "title": "Example headline",
      "description": "Article description",
      "content": "Article body when available",
      "authorsByline": "Example author",
      "source": {
        "domain": "publisher.example",
        "paywall": false,
        "location": { "country": "US" }
      },
      "imageUrl": "https://publisher.example/image.jpg",
      "country": "US",
      "language": "en",
      "pubDate": "2026-08-17T12:00:00Z",
      "links": [],
      "people": [],
      "companies": [],
      "entities": [],
      "sentiment": {},
      "categories": [],
      "topics": [],
      "taxonomies": []
    }
  ]
}
~~~

GET /v1/stories/all returns status, numResults, and results with story-cluster fields such as id, name, summary, createdAt, updatedAt, uniqueCount, reprintCount, totalCount, countries, topics, categories, people, companies, and taxonomies. Use clusterId on article search to retrieve the articles in one story.

#### Beans Comparison

Perigon is a direct high-capability news API comparison: individual article search, story clustering, semantic retrieval, generated summarization, source discovery, and entity enrichment. Its article payload is much larger than Beans and its content availability remains source-, license-, and plan-dependent. It does not document a native blog discriminator, whereas Beans exposes content_type news|blog directly.

### TheNewsAPI

Official references: [documentation](https://www.thenewsapi.com/documentation), [pricing](https://www.thenewsapi.com/pricing).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /v1/news/all | search, search_fields, locale, category/domain/source include and exclude lists, language, dates, sort, limit, page | Live/historical article search. Search supports +, |, -, quoted phrases, prefixes, and parentheses. |
| GET /v1/news/top | Same primary filters as all | Ranked top-story feed. |
| GET /v1/news/headlines | locale, domains/source IDs, language, published_on, headlines_per_category, include_similar | Category-grouped headline feed with optional similar articles. |
| GET /v1/news/uuid/{uuid} | api_token | Retrieves a discovered article card. |
| GET /v1/news/similar/{uuid} | Category/domain/source/language/date, limit, page | Finds articles similar to one known UUID. |
| GET /v1/news/sources | category/language, page | Source catalog. |

#### Query Parameter Semantics

- Text: search accepts Boolean operators, quoted phrases, prefixes, and grouped expressions; search_fields chooses the article fields matched.
- Publisher scope: locale, domains, sources, and their exclude variants provide country- and publisher-level filtering.
- Classification and language: category and language accept include/exclude lists. published_on and the documented date controls bound recency.
- Ranking and pagination: sort controls ranking; limit and page drive the list. The meta envelope provides found, returned, limit, and page.
- Related discovery: include_similar is available on category headlines, while /similar/{uuid} applies the same principal filters to a known article.

List envelopes are {meta, data}. meta contains found, returned, limit, page. Each article card includes uuid, title, description, keywords, 60-character snippet, URL, image_url, language, published_at, source, categories, locale, and relevance where applicable.

#### Response Payload

Field-preserving excerpt from the live list-response documentation:

~~~json
{
  "meta": {
    "found": 1234,
    "returned": 3,
    "limit": 3,
    "page": 1
  },
  "data": [
    {
      "uuid": "article-uuid",
      "title": "Article headline",
      "description": "Publisher description",
      "keywords": "keyword, another",
      "snippet": "First 60 characters of the article body...",
      "url": "https://publisher.example/article",
      "image_url": "https://publisher.example/image.jpg",
      "language": "en",
      "published_at": "2026-08-17T12:00:00.000000Z",
      "source": "Publisher",
      "categories": ["business"],
      "relevance_score": 0.89,
      "locale": "us"
    }
  ]
}
~~~

The optional relevance_score is returned for relevant search/similar results; the 60-character snippet is not a full article body.

#### Beans Comparison

Compact payloads, total results, a source catalog, and related articles are useful patterns; it is not a full-content reference.

### finlight

Official references: [REST endpoints](https://docs.finlight.me/en/v2/rest-endpoints/), [pricing](https://finlight.me/pricing).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| POST /v2/articles | query, sources/excludeSources, tickers, countries, categories, includeEntities, from, to, language, orderBy, order, pageSize, page | Financial-news search with ticker/company filters, categories, sentiment, and confidence. |
| GET /v2/articles/by-link | Required link; includeContent, includeEntities | Known-article lookup with optional body and company entities. includeContent is deprecated in current docs. |
| GET /v2/sources | No main filters documented | Source domain, origin country, languages, and availability/default status. |

#### Query Parameter Semantics

- Search and time: query, from, and to support financial-news search over a publication window.
- Publisher, issuer, and topic: sources/excludeSources, tickers, countries, and categories bound the corpus by publisher, company, geography, and classification.
- Entity enrichment: includeEntities asks for company/ticker-oriented enrichment where the caller's entitlement permits it.
- Ordering and pagination: orderBy, order, pageSize, and page define the ranked result set.
- Detail lookup: /v2/articles/by-link takes the exact article link. includeContent is documented but deprecated; the endpoint may return content when available.

List responses are {status, page, pageSize, articles}. An article contains link, source, title, nullable summary, publication/index/revision timestamps, language, sentiment, confidence, images, countries, categories, and tier-gated companies. Detail uses {status, article} and can contain content.

#### Response Payload

Live list-response shape:

~~~json
{
  "status": "ok",
  "page": 1,
  "pageSize": 20,
  "articles": [
    {
      "link": "https://publisher.example/article",
      "source": "publisher.example",
      "title": "Article headline",
      "summary": "Article summary",
      "publishDate": "2026-08-17T12:00:00Z",
      "language": "en",
      "sentiment": "positive",
      "confidence": "0.95",
      "images": ["https://publisher.example/image.jpg"],
      "countries": ["US"],
      "categories": ["markets", "technology"],
      "companies": ["Example Corp"]
    }
  ]
}
~~~

GET /v2/articles/by-link changes the inner key to article and may add content. The current docs treat company/entity fields and content availability as entitlement- and source-dependent.

#### Beans Comparison

Vertical but valuable as a model for tickers, confidence, and revision-aware delivery. It does not document universally available full text.

### GNews

Official references: [search](https://docs.gnews.io/endpoints/search-endpoint), [top headlines](https://docs.gnews.io/endpoints/top-headlines-endpoint), [response schema](https://docs.gnews.io/json-response), [pricing](https://gnews.io/).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /api/v4/search | Required q; lang, country, max, in, nullable, from, to, sortby, page, truncate | Keyword search with phrase, AND, OR, NOT syntax and title/description/content selection. |
| GET /api/v4/top-headlines | category, lang, country, max, nullable, from, to, q, page, truncate | Google News-ranked current headlines. |

#### Query Parameter Semantics

- Text: q is required for /search. Boolean and phrase syntax are accepted; in limits matching to title, description, content, or a documented combination.
- Geography and language: lang and country scope the corpus. country is not an article-level source-country filter.
- Time and ranking: from, to, and sortby control the publication window and result ordering.
- Result shape: max limits returned articles; page selects the result window; nullable and truncate alter fields in the returned article objects.
- Headlines: /top-headlines adds category and optional q filtering against the provider-ranked current feed.

Responses are {totalArticles, articles}. Article fields are id, title, description, content, URL, image, publishedAt, lang, and source; source.country is Search-only.

#### Response Payload

Live response example shared by the search and top-headlines documentation:

~~~json
{
  "totalArticles": 54904,
  "articles": [
    {
      "id": "article-id",
      "title": "Article headline",
      "description": "Article description",
      "content": "Article content",
      "url": "https://publisher.example/article",
      "image": "https://publisher.example/image.jpg",
      "publishedAt": "2026-08-17T12:00:00Z",
      "lang": "en",
      "source": {
        "id": "publisher-id",
        "name": "Publisher",
        "url": "https://publisher.example",
        "country": "us"
      }
    }
  ]
}
~~~

content is plan-sensitive: the provider documents truncated content on Free and full content on paid tiers. source.country is documented only for search.

#### Beans Comparison

The simplest paid full-content design, but without a source catalog, story model, or propagation route.

### NewsAPI.ai / Event Registry

Official references: [article search](https://www.newsapi.ai/documentation?tab=searchArticles), [data model](https://www.newsapi.ai/documentation?tab=data_models), [plans](https://www.newsapi.ai/plans).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET or POST /api/v1/article/getArticles | keyword, keywordOper, keywordSearchMode, keywordLoc, conceptUri, sourceUri, categoryUri, dates, nested query, articlesPage, articlesCount, sort, articleBodyLen, dataType, includeArticle flags | Searches news, blog, and pr independently or together. Supports Boolean/nested filters, 100 article pages, and time/source/sentiment/category aggregations. |
| GET or POST /api/v1/article/getArticle | articleUri, infoArticleBodyLen, includeArticle flags | Retrieves known article URIs with body and selected enrichments. |
| GET or POST /api/v1/minuteStreamArticles | Recent-activity continuation values, dataType, filters, articleBodyLen, includeArticle flags | Incremental stream with separate continuation values for news and blogs. |
| GET or POST /api/v1/annotate, /categorize, /semanticSimilarity, /sentiment, /extractArticleInfo | Text plus endpoint options | Optional analysis of caller-supplied content. |

#### Query Parameter Semantics

- Text and Boolean logic: keyword, keywordOper, keywordSearchMode, and keywordLoc control lexical query construction and where matching occurs. Nested query expressions combine text, source, concept, category, and date constraints.
- Corpus selection: dataType selects news, blog, pr, or the provider's documented combination. This is the clearest native news-versus-blog control in the reviewed set.
- Source and knowledge filters: sourceUri, conceptUri, categoryUri, language, dates, and sort provide source, entity/concept, classification, time, and ordering control.
- Page and body size: articlesPage and articlesCount paginate up to the documented page size; articleBodyLen=-1 requests the full body.
- Optional analysis: includeArticle flags determine enrichments. Separate annotate, categorize, semanticSimilarity, sentiment, and extractArticleInfo endpoints accept caller-supplied text.

Article payloads can include uri, URL, title, body, dates, lang, dataType, sentiment, image, source, categories, disambiguated concepts, authors, links, videos, social shares, duplicate list, extracted dates, and event/story URIs. Most enrichment is optional.

#### Response Payload

The live data-model documentation exposes the resultType=articles item schema. Its outer list envelope is rendered client-side and was not available to the crawler, so this is the documented article item rather than an inferred wrapper:

~~~json
{
  "uri": "article-uri",
  "url": "https://publisher.example/article",
  "title": "Article headline",
  "body": "Article body",
  "date": "2026-08-17",
  "time": "12:00:00",
  "dateTime": "2026-08-17T12:00:00Z",
  "dateTimePub": "2026-08-17T11:45:00Z",
  "lang": "eng",
  "isDuplicate": false,
  "dataType": "news",
  "sentiment": -0.2,
  "eventUri": "event-uri",
  "relevance": 34,
  "storyUri": "story-uri",
  "image": "https://publisher.example/image.jpg",
  "source": { "uri": "source-uri", "title": "Publisher" },
  "categories": [],
  "concepts": [],
  "authors": [],
  "links": [],
  "videos": [],
  "shares": {},
  "duplicateList": [],
  "extractedDates": [],
  "location": null
}
~~~

GET/POST /api/v1/article/getArticles paginates with articlesPage and articlesCount. articleBodyLen=-1 requests the full body; includeArticle flags control optional enrichments.

#### Beans Comparison

The most complete direct comparison. Its additional primitives are canonical detail lookup, dataType news|blog|pr, full-body control, duplicate/similarity context, aggregation, and analysis APIs. Its contract is much larger than Beans.

### NewsData.io

Official references: [latest endpoint](https://newsdata.io/blog/latest-news-endpoint/), [response fields](https://newsdata.io/blog/news-api-response-object/), [pricing model](https://newsdata.io/blog/pricing-plan-in-newsdata-io/).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /api/1/latest | apikey, id, q, qInTitle, qInMeta, country/language/domain filters, timeframe, timezone, full_content, image, size, page | Latest feed, normally covering the last 48 hours, with content, image, timeframe, and field filters. |
| GET /api/1/archive | Search/timeframe/content/image/page controls | Historical article search with the same full-content switch. |
| Latest/news-count, market, and crypto endpoints | Endpoint-specific filters | Adds aggregates and domain-specific feeds. |

#### Query Parameter Semantics

- Text: q, qInTitle, and qInMeta select general, title-only, or metadata-focused matching.
- Corpus controls: country, language, domain, id, image, and timeframe narrow the set by geography, language, publisher, known record, asset availability, and recency.
- Content delivery: full_content=1 asks for content when the provider has it. It does not convert publisher HTML into a structured document.
- Pagination: size controls the list size. nextPage is an opaque continuation string and must be passed unchanged as page.
- Historical scope: /api/1/latest is a recent-feed route; /api/1/archive is the historical search surface.

The response includes status, results, and nextPage. Result fields include IDs, title/description, content, publication date, image/source, category/tags, sentiment, and AI region. Send nextPage unchanged as page to continue.

#### Response Payload

The provider's live response-object documentation lists this field layout, but its full JSON example is not accessible to the crawler. This is a field-preserving schema excerpt, not a captured HTTP response:

~~~json
{
  "status": "<provider status>",
  "totalResults": 123,
  "results": [
    {
      "article_id": "article-id",
      "title": "Article headline",
      "link": "https://publisher.example/article",
      "keywords": ["keyword"],
      "creator": ["Author"],
      "video_url": null,
      "description": "Article description",
      "content": "Article content when requested and available",
      "pubDate": "2026-08-17 12:00:00",
      "image_url": "https://publisher.example/image.jpg",
      "source_id": "publisher-id",
      "source_url": "https://publisher.example",
      "source_icon": "https://publisher.example/icon.png",
      "source_priority": 1,
      "country": ["united states"],
      "category": ["business"],
      "language": "english",
      "ai_tag": "topic",
      "sentiment": "positive",
      "sentiment_stats": {},
      "ai_region": "North America"
    }
  ],
  "nextPage": "opaque-continuation"
}
~~~

The content field is governed by full_content=1 and source availability. Use nextPage unchanged as the next request's page value.

#### Beans Comparison

Full-content selection, AI enrichment, a count route, and opaque continuation are worthwhile comparison points. Beans has more direct propagation support.

### NewsAPI.org

Official references: [Everything](https://newsapi.org/docs/endpoints/everything), [Top Headlines](https://newsapi.org/docs/endpoints/top-headlines), [Sources](https://newsapi.org/docs/endpoints/sources), [pricing](https://newsapi.org/pricing).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /v2/everything | q, searchIn, sources, domains, excludeDomains, from, to, language, sortBy, pageSize, page | Broad news/blog search with phrase/Boolean operators, domain inclusion/exclusion, and sorting. |
| GET /v2/top-headlines | country, category, sources, q, pageSize, page | Current breaking/top-headline feed. |
| GET /v2/top-headlines/sources | category, language, country | Source discovery. |

#### Query Parameter Semantics

- Text: q supports keyword, phrase, and Boolean-style expressions. searchIn selects title, description, or content matching.
- Publisher scope: sources is an allowlist of source IDs; domains and excludeDomains control publisher domains.
- Time and ranking: from, to, language, and sortBy control the publication range, locale, and relevance/popularity/date ordering.
- Feed selection: /top-headlines uses country, category, sources, and optional q; source discovery is separate at /top-headlines/sources.
- Pagination: page and pageSize control list retrieval; totalResults reports the matching count.

Responses are {status, totalResults, articles}; article fields include source, author, title, description, URL, image URL, publication time, and truncated content.

#### Response Payload

Live response example for /v2/everything and /v2/top-headlines:

~~~json
{
  "status": "ok",
  "totalResults": 38,
  "articles": [
    {
      "source": { "id": "publisher-id", "name": "Publisher" },
      "author": "Author",
      "title": "Article headline",
      "description": "Article description",
      "url": "https://publisher.example/article",
      "urlToImage": "https://publisher.example/image.jpg",
      "publishedAt": "2026-08-17T12:00:00Z",
      "content": "First 200 characters of article content..."
    }
  ]
}
~~~

The provider documents content as a 200-character truncation, including on paid plans; it is not a full-body field.

#### Beans Comparison

Baseline search/headline/source contract. Beans extends it with content selection, enrichment, trends, and propagation.

### Currents

Official references: [endpoint overview](https://currentsapi.services/en/docs/endpoint), [latest news](https://currentsapi.services/en/docs/latest_news), [pricing](https://currentsapi.services/en/product/price).

#### Route Inventory

| Route | Main parameters | Capability |
| --- | --- | --- |
| GET /v1/latest-news and /v2/latest-news | language, country, category, type, domain, domain_not, author, page_number, page_size | Latest global feed. V2 exposes a canonical category taxonomy. |
| GET /v1/search and /v2/search | keywords or query, language/category/date, domain, page_number/page_size, cursor (V2) | Historical/news search complementing latest. V2 adds keyset cursor pagination. |
| GET /v1/available/categories, /available/regions, /available/languages and v2 equivalents | None | Exact filter-value discovery. |

#### Query Parameter Semantics

- Latest feed scope: language, country, category, type, domain, domain_not, and author filter the current feed.
- Search: /search accepts keywords for standard term search or query for Boolean expressions, plus date, language, category, publisher, and pagination controls. When both are sent, keywords takes precedence.
- Taxonomy discovery: the available/categories, available/regions, and available/languages routes return documented filter values instead of requiring callers to infer them.
- Pagination: V1 uses page_number and page_size. V2 search supports cursor and returns next_cursor for keyset continuation; preserve the token unchanged. The latest response returns page but does not include a total count in the documented example.
- Versioning: V2 publishes a canonical category taxonomy; use version-specific field/value documentation rather than assuming V1 and V2 are interchangeable.

Latest responses are {status, news, page}; V2 search responses also return next_cursor. The documented news fields are id, title, description, URL, author, image, language, category array, and publication time. No body field is listed.

#### Response Payload

Live latest-news response example:

~~~json
{
  "status": "ok",
  "news": [
    {
      "id": "article-id",
      "title": "Article headline",
      "description": "Article description",
      "url": "https://publisher.example/article",
      "author": "Author",
      "image": "https://publisher.example/image.jpg",
      "language": "en",
      "category": ["business"],
      "published": "2026-08-17 12:00:00 +0000"
    }
  ],
  "page": 1,
  "next_cursor": "opaque-cursor-or-null"
}
~~~

No article-body or content field appears in the provider's current documented response example. next_cursor is a V2 search continuation field.

#### Beans Comparison

Useful for a latest/search separation and published taxonomy values, but not a documented full-content alternative.

### Former Provider: NewsDataHub

[NewsDataHub](https://newsdatahub.com/) says it shut down on 2026-06-05 and no longer accepts new requests or subscriptions. It should not be used as a production dependency or current contract comparison.

#### Response Payload

Unavailable. The provider's only current public page is its shutdown notice, not live API documentation or a response example. A historical payload would be stale and should not be treated as an active contract.

## 4. Text, Fuzzy, and Semantic Search Inputs

A free-text parameter is not necessarily fuzzy or semantic. This reference uses three distinct labels:

- Lexical: terms, phrases, field targeting, Boolean operators, or relevance-ranked keyword matching over indexed text.
- Fuzzy: the provider explicitly documents tolerant matching over labels or text.
- Semantic: the provider accepts natural language and ranks records by meaning/similarity rather than only token or phrase matching.

### Current Cafecito API Inputs

| API | Lexical / fuzzy input | Semantic input | Contract distinction |
| --- | --- | --- | --- |
| Beans | tags is a fuzzy CSV filter across category, region, and entity labels. | q is the natural-language semantic query; acc controls match strictness. | q and tags are combinable. Exact categories, regions, entities, and sources are separate structured filters. |
| Espresso | tags is a fuzzy CSV match over persisted Event or Signal tag labels. Source q is case-insensitive metadata matching only. | q is an optional natural-language semantic query on Event and Signal collections. | q on source discovery is not semantic; only Event and Signal q use the semantic collection-search contract. |

### External Provider Inputs

| Provider | Lexical or fuzzy-style input | Semantic retrieval input | Classification |
| --- | --- | --- | --- |
| World News API | text with text-match-indexes supports terms, quoted phrases, AND/OR, parentheses, and exclusions. | None documented on the reviewed news-search route. | Lexical. Wildcards are not supported; no fuzzy or semantic parameter is documented. |
| GDELT Cloud | Structured story filters narrow the candidate set. | search on v2 Story list is a natural-language semantic query; the provider documents similarity ranking on list endpoints. | Semantic, without a separate score-threshold parameter. |
| Perigon | q, title, desc, content, name, and field filters perform article/story text retrieval. | POST /v1/vector/news/all accepts prompt; scoreThreshold controls the semantic result cutoff. | Separate lexical and semantic routes. |
| TheNewsAPI | search and search_fields support phrase, Boolean, prefix, and grouped lexical matching. | None documented. | Lexical. |
| finlight | query supports keyword logic, exact phrases, inclusion/exclusion operators, and documented field-level filters. | None documented on the reviewed article-search route. | Advanced lexical. |
| GNews | q with in selects title, description, content, or a documented combination; phrases and Boolean operators are supported. | None documented. | Lexical. |
| NewsAPI.ai / Event Registry | keyword, keywordOper, keywordSearchMode, keywordLoc, and the advanced query object control article matching. keywordSearchMode=simple is relevance-ranked keyword matching. | No corpus semantic-search parameter is documented for getArticles. semanticSimilarity is a separate text-analysis endpoint for supplied content. | Lexical/relevance-ranked retrieval, plus standalone semantic analysis. |
| NewsData.io | q, qInTitle, and qInMeta target keywords/phrases across article, title, or metadata fields; Boolean operators are documented. | None documented. | Lexical. |
| NewsAPI.org | q and searchIn provide keyword, phrase, and Boolean-style matching. | None documented. | Lexical. |
| Currents | keywords performs term search; query accepts Boolean syntax. | None documented. | Lexical. |

### Comparison Notes

- Explicit fuzzy matching is unusual in the reviewed provider set. Beans tags and Espresso tags are the only public inputs documented as fuzzy label matching.
- GDELT Cloud search, Perigon prompt, Beans q, and Espresso Event/Signal q are the corpus-retrieval inputs documented as semantic.
- NewsAPI.ai simple mode is broad, relevance-ranked keyword search, not semantic retrieval. Its semanticSimilarity endpoint analyzes caller-supplied text and should not be represented as an alternative to corpus semantic search.
- Do not infer typo tolerance, vector retrieval, or meaning-based ranking from a provider merely accepting free text, Boolean expressions, phrases, or a relevance sort.

## 5. Pagination Strategy: Integer Page vs Opaque Cursor

The providers use three related but incompatible navigation models: integer page/offset values, an opaque cursor token, and a provider continuation token named page. A client must select the model from the route contract, not from the parameter name alone.

### Provider Grouping

| Request style | Providers | Continuation rule |
| --- | --- | --- |
| Integer page or offset | Beans: offset/limit; World News API: offset/number; Perigon: page/size; TheNewsAPI: page/limit; finlight: page/pageSize; GNews: page/max; NewsAPI.ai article search: articlesPage/articlesCount; NewsAPI.org: page/pageSize; Currents V1: page_number/page_size | Increase the integer after a successful page. Perigon explicitly documents zero-based page numbering; use each provider's documented starting value and stop condition. |
| Opaque cursor | GDELT Cloud: cursor/limit; Currents V2 search: cursor/page_size | Send pagination.next_cursor or next_cursor back as cursor. The token is provider-generated and must not be parsed or calculated by the client. |
| Opaque continuation named page | NewsData.io: page/size | Send the response nextPage value back unchanged as page. Despite its name, page is not an integer on continuation requests. |
| Stream continuation | NewsAPI.ai minuteStreamArticles | Use the provider's stream continuation values. This is separate from integer articlesPage for article search. |

### Tradeoffs

| Concern | Integer page / offset | Opaque cursor / continuation |
| --- | --- | --- |
| Best fit | Bounded historical search, administration tables, and UIs that need arbitrary page jumps. | Latest feeds, incremental ingestion, and deep traversal of changing article corpora. |
| Random access | A client can request a known page or offset directly. | Normally sequential only; a client needs saved prior cursors to navigate backward. |
| Deep retrieval cost | Can degrade at high offsets, depending on provider implementation. | Usually supports efficient forward traversal without a large calculated offset. |
| New articles while paging | A changing result set can create duplicates or skipped records between pages. | Can offer more stable continuation when the provider binds the token to a consistent order or result snapshot. |
| Total results | Natural to expose as totalResults, found, or available, though counting can be expensive or approximate. | Often unavailable; the terminal condition is an absent next token or empty data array. |
| URL and caching behavior | Short, inspectable, shareable URLs. | Tokens can be long, expire, be account-bound, and be unsafe to reuse after filter changes. |
| Client behavior | Increment the requested number and apply a documented stop condition. | Preserve and URL-encode the token exactly; do not decode, mutate, or synthesize it. |

### Practical Contract Rules

- Keep query filters, sort order, and page size stable across a traversal. Changing them creates a new result set; cursor and continuation tokens commonly become invalid.
- For integer pagination over live news, add a bounded time range and stable sort where the provider supports them. This reduces, but does not eliminate, duplicates and omissions while new articles arrive.
- Stop integer pagination on an empty or short page, or when a documented total is reached. Beans has no total-count envelope, so a client must use the returned array length.
- Stop cursor/continuation pagination when next_cursor or nextPage is absent or the response has no records.
- Do not send offset/page and cursor/continuation parameters together unless a provider explicitly documents that combination.

## 6. Cross-Provider Contract Comparison

| Provider | Primary retrieval model | Detail / relationship route | Discovery surface | Full-content control | Pagination / response convention |
| --- | --- | --- | --- | --- | --- |
| World News API | General news search and clustered top news | GET /retrieve-news by IDs; top-news clusters related coverage | GET /search-news-sources | text is documented full text | offset/number; offset, number, available, news |
| GDELT Cloud | Generated clustered stories with article evidence | Story detail and /stories/{story_id}/articles | GET /geo/admin1; entity references in stories | No documented full-body article contract | Cursor; success, data, pagination |
| Perigon | Individual articles, story clusters, and semantic news search | Story clusters connect through clusterId; no reviewed canonical article-detail route | Source, journalist, people, company, and topic routes | Article content where source, license, and plan permit | Zero-based page/size; status, numResults, articles or results |
| TheNewsAPI | Article cards and category headlines | UUID lookup and /similar/{uuid} | GET /news/sources | No full body; snippet only | page/limit; meta plus data |
| finlight | Financial article search | GET /articles/by-link | GET /sources | Detail content may be available | page/pageSize; status plus articles |
| GNews | Keyword search and ranked top headlines | No reviewed canonical detail or relationship route | No separate reviewed source catalog | Paid content; Free truncates | page/max; totalArticles plus articles |
| NewsAPI.ai / Event Registry | News, blog, and PR article search plus stream | GET/POST getArticle; event/story/duplicate references | Source, concept, and category filters; supporting metadata endpoints | articleBodyLen=-1 | articlesPage/articlesCount; provider-configurable result |
| NewsData.io | Recent and archive article feeds | No reviewed canonical article-detail route | Endpoint-specific count/domain/taxonomy features | full_content=1 where available | Opaque nextPage; status plus results |
| NewsAPI.org | Everything search and top headlines | No reviewed canonical detail/relationship route | GET /top-headlines/sources | Always 200-character truncation | page/pageSize; status, totalResults, articles |
| Currents | Latest feed and historical search | No reviewed canonical detail/relationship route | Available category, region, and language routes | No current body field | V1 page_number/page_size; V2 cursor; status, news, page, next_cursor |

The strongest general comparison is World News API for broad indexed search and clustered coverage. Perigon is the highest-capability direct article/news comparison. GNews is the compact paid-content model. NewsAPI.ai is the widest route surface for native news/blog selection and analysis. GDELT Cloud is the story-first, article-evidence comparison, while finlight is the relevant financial vertical. The other providers are useful baselines for compact envelopes, source discovery, continuation tokens, and taxonomy routes.

## 7. Implications for Beans Route Comparison

For route-name and contract comparison, the closest established patterns are:

- Collection search: World News API /search-news, Perigon /v1/articles/all, TheNewsAPI /v1/news/all, GNews /api/v4/search, NewsAPI.ai /api/v1/article/getArticles, NewsData.io /api/1/archive, NewsAPI.org /v2/everything, and Currents /v2/search. GDELT Cloud /api/v2/stories is story-cluster retrieval rather than direct article search.
- Current feeds: World News API /top-news, TheNewsAPI /v1/news/top and /headlines, GNews /top-headlines, NewsData.io /latest, NewsAPI.org /top-headlines, and Currents /latest-news.
- Detail lookup: World News API /retrieve-news, TheNewsAPI /uuid/{uuid}, finlight /by-link, NewsAPI.ai /getArticle, and GDELT Cloud /stories/{story_id} are the clearest provider patterns. Perigon's reviewed route surface connects individual articles to a story through clusterId rather than documenting a canonical article-detail route. Beans currently relies on article URLs in collection and propagation workflows; a canonical detail route would be a deliberate surface expansion.
- Blog selection: NewsAPI.ai dataType has separate news, blog, and PR values. Perigon and GDELT Cloud do not document a dedicated blog type in their reviewed public route surfaces. Beans already exposes content_type news|blog, which is a smaller, direct contract that does not adopt a PR category.
- Source and taxonomy discovery: Perigon, TheNewsAPI, finlight, World News API, NewsAPI.org, and Currents expose provider-defined sources or filter values. Beans already has category, entity, region, and source discovery routes; exact filter values should remain discoverable rather than hard-coded by clients.
- Related coverage: GDELT Cloud offers story clusters with paginated article evidence; Perigon exposes story clusters linked through article clusterId; World News API clusters top news; TheNewsAPI has similar-by-UUID; and NewsAPI.ai exposes event/story/duplicate context. Beans propagation is a distinct URL-driven comparison workflow; it should not be described as provider-native story clustering.
- Pagination and count semantics: GDELT Cloud and Currents V2 search use cursor envelopes; Perigon uses zero-based page and size; TheNewsAPI, NewsAPI.org, and GNews offer count-like fields; NewsData.io uses an opaque continuation token; World News API uses an available count with offset. Beans uses limit and offset with arrays and no total-count envelope. Adding a total or continuation field would be a public contract decision, not a requirement for parity.
- Content semantics: external content fields vary from 60-character snippets to Perigon's source/license/plan-dependent bodies. GDELT Cloud's story-article evidence does not establish a full-body contract. Beans optional content must retain the same source-availability caveat and must not imply Markdown, structured extraction, canonical publisher copies, or complete publisher coverage.
- Enrichment and analysis: Perigon, World News API, finlight, NewsData.io, and NewsAPI.ai add sentiment, entities, concepts, AI tags, semantic retrieval, summaries, or standalone analysis. GDELT Cloud adds generated story/event context rather than a general article-analysis endpoint. Beans has practical enrichment labels; adding complex analysis endpoints would materially broaden the product rather than complete an existing route family.
- Response compatibility: a stable Beans envelope should preserve its documented public article fields and pagination behavior. Provider field names, identifiers, and content rights should be normalized at the public contract boundary only when a new capability requires it.

Pricing remains a dated procurement input. Providers charge by points, requests, credits, or tokens; result limits and content rights determine effective cost.

## 8. Source and Version Notes

This reference summarizes provider documentation and representative payloads reviewed on 2026-08-17. Provider schemas, route aliases, parameters, limits, pricing, and content rights can change by plan, account, and source. Confirm the current provider documentation and account-specific OpenAPI/schema before implementing a client or declaring exact compatibility.

Primary documentation:

- World News API: [search](https://worldnewsapi.com/docs/search-news/), [retrieve](https://worldnewsapi.com/docs/retrieve-news/), [pricing](https://worldnewsapi.com/pricing)
- GDELT Cloud: [API v2](https://docs.gdeltcloud.com/api-reference/v2), [documentation guide](https://docs.gdeltcloud.com/API_DOCUMENTATION_GUIDE)
- Perigon: [getting started](https://docs.perigon.io/docs/getting-started), [article data](https://docs.perigon.io/docs/article-data), [story data](https://docs.perigon.io/docs/story-data), [vector news](https://docs.perigon.io/docs/vector-endpoint), [sources](https://docs.perigon.io/docs/searching-sources)
- TheNewsAPI: [documentation](https://www.thenewsapi.com/documentation), [pricing](https://www.thenewsapi.com/pricing)
- finlight: [REST reference](https://docs.finlight.me/en/v2/rest-endpoints/), [advanced query building](https://docs.finlight.me/en/v2/advanced-query-building/), [pricing](https://finlight.me/pricing)
- GNews: [search](https://docs.gnews.io/endpoints/search-endpoint), [response schema](https://docs.gnews.io/json-response), [pricing](https://gnews.io/)
- NewsAPI.ai: [article search](https://www.newsapi.ai/documentation?tab=searchArticles), [data model](https://www.newsapi.ai/documentation?tab=data_models), [plans](https://www.newsapi.ai/plans)
- NewsData.io: [latest](https://newsdata.io/blog/latest-news-endpoint/), [response fields](https://newsdata.io/blog/news-api-response-object/), [pricing](https://newsdata.io/blog/pricing-plan-in-newsdata-io/)
- NewsAPI.org: [Everything](https://newsapi.org/docs/endpoints/everything), [pricing](https://newsapi.org/pricing)
- Currents: [search](https://currentsapi.services/en/docs/search), [latest news](https://currentsapi.services/en/docs/latest_news), [pricing](https://currentsapi.services/en/product/price)
- NewsDataHub: [shutdown notice](https://newsdatahub.com/)
