# Industry Event and News API Route Reference

Status: External API reference
Reviewed: 2026-08-11

This document lists the documented public routes, principal parameters, response payloads, and purpose of the GDELT Cloud, PredictHQ, and Perigon APIs. It is a reference for comparing the Espresso route proposal with established event and news-intelligence services.

Scope: this covers the public event, story, article, entity, geography, search, aggregation, and related intelligence routes that are relevant to Espresso. It does not attempt to inventory provider account, billing, webhook, administration, private saved-search, or unrelated product routes. Provider paths are shown without a product-specific gateway prefix.

## Common conventions

| Provider | Authentication | Pagination | Typical success envelope |
|---|---|---|---|
| GDELT Cloud | Bearer token in Authorization | Cursor pagination: limit, cursor, next_cursor | success, data, pagination |
| PredictHQ | Bearer token in Authorization | Offset pagination: limit, offset, next, previous | count, overflow, previous, next, results |
| Perigon | apiKey query parameter, Authorization bearer token, or x-api-key | Page pagination: zero-based page and size | status, numResults, articles or results |

All three providers return JSON for these routes. The providers use different identifiers, date semantics, grouping models, and pagination conventions.

## 1. GDELT Cloud API v2

Official reference: https://docs.gdeltcloud.com/api-reference/v2

### Route inventory

| Method and route | Main parameters | What it does |
|---|---|---|
| GET /api/v2/events | date_start, date_end, country, region, continent, admin1, event_family, category, subcategory, domain, has_fatalities, confidence_profile, min_confidence, search, sort, limit, cursor | Searches normalized GDELT events by time, geography, event classification, confidence, fatalities, or free-text/semantic search. |
| GET /api/v2/events/{event_id} | event_id path parameter | Returns one normalized event and its linked story, entity, and article references. |
| GET /api/v2/events/summary | The event filters plus group_by, usually country, region, continent, category, subcategory, or date | Returns aggregate event counts and metrics rather than individual event records. |
| GET /api/v2/stories | date_start, date_end, country, region, continent, admin1, category, event_category, subcategory, domain, has_events, has_fatalities, confidence_profile, article_count_min, article_count_max, search, sort, limit, cursor | Searches clustered stories, which group related articles around a narrative. |
| GET /api/v2/stories/{story_id} | story_id path parameter | Returns one story, its metrics, linked events, entities, and representative articles. |
| GET /api/v2/stories/{story_id}/articles | story_id path parameter, limit, cursor | Paginates through the article evidence belonging to a story. |
| GET /api/v2/stories/summary | The story filters plus group_by | Returns aggregate story counts and metrics. |
| GET /api/v2/entities | Entity search and pagination parameters documented by the provider | Lists canonical entities associated with stories and events. |
| GET /api/v2/entities/{entity_id} | entity_id path parameter | Returns one entity profile and its linked story/event references. |
| GET /api/v2/geo/admin1 | country query parameter | Returns discoverable first-level administrative regions for a country, for use in event and story filters. |

The provider changelog also mentions newer cross-source search and event-to-story traversal routes. Those should be verified against the current reference before being treated as a stable contract.

### Event search parameters

- Time: date_start and date_end, documented as YYYY-MM-DD.
- Geography: country, region, continent, and admin1.
- Classification: event_family, category, and subcategory.
- Source context: domain.
- Fatalities: has_fatalities.
- Confidence: confidence_profile and min_confidence.
- Relevance: search for free-text or semantic retrieval; sort is significance or recent.
- Pagination: limit and cursor.

### Event response payload

The event list and event detail payloads use the same event-card shape. The list is wrapped in a paginated envelope.

~~~json
{
  "success": true,
  "data": [
    {
      "id": "event-id",
      "url": "https://...",
      "primary_story_url": "https://...",
      "family": "conflict",
      "title": "Example event",
      "summary": "Short event summary",
      "event_date": "2026-08-10",
      "category": "protest",
      "subcategory": "civil_unrest",
      "domain": "example.com",
      "event_code": "code",
      "geo": {
        "country": "US",
        "region": "Americas",
        "continent": "North America",
        "admin1": "California",
        "location": "San Francisco",
        "latitude": 37.7749,
        "longitude": -122.4194
      },
      "geo_context": {
        "location_country": "US",
        "actor_origin_countries": ["US"]
      },
      "actors": [
        {"name": "Example actor", "country": "US", "role": "participant"}
      ],
      "metrics": {
        "significance": 0.82,
        "confidence": 0.91
      },
      "has_fatalities": false,
      "fatalities": 0,
      "story_refs": ["story-id"],
      "entity_refs": ["entity-id"],
      "top_articles": [
        {"title": "Example article", "url": "https://...", "source": "example.com"}
      ]
    }
  ],
  "pagination": {
    "limit": 20,
    "cursor": null,
    "next_cursor": "next-cursor"
  }
}
~~~

A detail response uses data as one object instead of an array. The provider documents top_articles as an inline convenience set; the story articles route is used when the complete article set is needed.

### Story response payload

~~~json
{
  "success": true,
  "data": [
    {
      "id": "story-id",
      "url": "https://...",
      "title": "Example story",
      "story_date": "2026-08-10",
      "category": "politics",
      "subcategory": "elections",
      "geo": {
        "country": "US",
        "region": "Americas",
        "continent": "North America",
        "admin1": "California"
      },
      "geo_context": {
        "location_country": "US"
      },
      "metrics": {
        "significance": 0.76,
        "article_count": 42,
        "linked_event_count": 3,
        "max_linked_event_significance": 0.88
      },
      "has_events": true,
      "has_fatalities": false,
      "fatalities": 0,
      "linked_events": ["event-id"],
      "entity_refs": ["entity-id"],
      "top_articles": [
        {"title": "Example article", "url": "https://...", "source": "example.com"}
      ]
    }
  ],
  "pagination": {
    "limit": 20,
    "cursor": null,
    "next_cursor": "next-cursor"
  }
}
~~~

### Summary, entity, and geography payloads

- Summary routes return success, group_by, and data. Each data item is an aggregate bucket containing the grouping value plus event or story counts and provider metrics.
- Entity routes return canonical entity identity fields, a public URL, and linked story/event references. The exact profile fields vary by entity type.
- The admin1 route returns the available first-level administrative-region values for the requested country.

## 2. PredictHQ Events API

Official references:
- Events search: https://docs.predicthq.com/api/events/search-events
- Event counts: https://docs.predicthq.com/api/events/get-event-counts
- Places: https://docs.predicthq.com/api/places
- Entities: https://docs.predicthq.com/getting-started/predicthq-data/entities

### Route inventory

| Method and route | Main parameters | What it does |
|---|---|---|
| GET /v1/events/ | q, id, entity.id, category, label, phq_label, country, place.scope, place.exact, within, start.*, end.*, active.*, first_seen.*, updated.*, cancelled.*, predicted_end.*, impact.*, rank.*, local_rank.*, state, parent.include, private.include, brand_unsafe.exclude, phq_attendance.*, predicted_event_spend.*, sort, limit, offset | Searches and filters predicted or observed events. |
| GET /v1/events/count/ | The same event filters as /v1/events/ | Returns aggregate event counts and rank/category/label distributions without returning the full event list. |
| GET /v1/places/ | q, id, country, location, type, limit | Searches the canonical place database used by event filters. |
| GET /v1/broadcasts/ | Event/broadcast filters documented by the provider, including time, country, category, and pagination filters | Lists broadcast or scheduled media events that can be used as event signals. |
| GET /v1/broadcasts/count/ | The broadcast filters | Returns counts for matching broadcasts. |

PredictHQ supports exact event lookup through the id filter on GET /v1/events/. The reviewed event API documentation does not define a separate GET /v1/events/{id} route.

### Event search parameters

- Text and identity: q, id, and entity.id.
- Classification: category, legacy label, phq_label, phq_label.op, and phq_label.exclude.
- Geography: country, place.scope, place.exact, within, and saved_location.location_id.
- Time: start, end, active, first_seen, updated, cancelled, predicted_end, and impact. Each supports comparison operators such as gte, gt, lte, and lt; time-zone variants use .tz.
- Rank and quality: rank, rank_level, local_rank, local_rank_level, location_confidence_score, and start_date_confidence_score.
- State and ownership: state, parent.include, private.include, private.user_id, and private.org_review.
- Impact and attendance: phq_attendance, impact.industry, predicted_event_spend, beam.analysis_id, and beam.group_id.
- Ordering and pagination: sort, limit, and offset. The documented default limit is 10 and the default offset is 0.
- Date input: the API documents YYYY-MM-DD and timestamp forms such as YYYY-MM-DDThh:mm:ss; UTC is the default when no time zone is supplied.

### Event response payload

~~~json
{
  "count": 1,
  "overflow": false,
  "previous": null,
  "next": "https://api.predicthq.com/v1/events/?offset=10",
  "results": [
    {
      "id": "event-id",
      "title": "Example event",
      "description": "Event description",
      "category": "concerts",
      "labels": ["concert", "music"],
      "rank": 73,
      "local_rank": 87,
      "phq_attendance": 13833,
      "entities": [
        {"entity_id": "entity-id", "name": "Example artist", "type": "person"}
      ],
      "start": "2026-08-10T19:00:00Z",
      "end": "2026-08-10T22:00:00Z",
      "active": "2026-08-10T19:00:00Z",
      "location": {
        "name": "Example venue",
        "formatted_address": "Example address",
        "geometry": {"coordinates": [-122.4194, 37.7749]}
      },
      "state": "active",
      "first_seen": "2026-01-01T00:00:00Z",
      "updated": "2026-08-01T00:00:00Z"
    }
  ]
}
~~~

This is a representative payload. PredictHQ event objects contain additional category-specific fields, parent/child relationships, recurrence or grouping metadata, impact data, and source information.

### Event count response

~~~json
{
  "count": 1,
  "top_rank": 73,
  "top_local_rank": 87,
  "rank_levels": {"1": 0, "2": 0, "3": 0, "4": 1, "5": 0},
  "local_rank_levels": {"1": 0, "2": 0, "3": 0, "4": 1, "5": 0},
  "categories": {"concerts": 1},
  "phq_labels": {"concert": 1}
}
~~~

The exact bucket values depend on the requested filters. phq_labels is the current field; labels is a legacy field in the provider documentation.

### Places response

A places response is a collection of canonical place records. A place record contains fields such as id, type, name, county, region, country, country_alpha2, country_alpha3, and location as a longitude/latitude pair.

~~~json
{
  "count": 1,
  "results": [
    {
      "id": "place-id",
      "type": "locality",
      "name": "San Francisco",
      "county": "San Francisco County",
      "region": "California",
      "country": "United States",
      "country_alpha2": "US",
      "country_alpha3": "USA",
      "location": [-122.4194, 37.7749]
    }
  ]
}
~~~

PredictHQ entities are normally embedded in event records and queried through event filters such as entity.id. They are not exposed as a generic, separately documented entity CRUD route in the reviewed event API surface.

## 3. Perigon News and Intelligence API

Official references:
- API overview: https://perigon.io/docs/api/intro
- Getting started: https://docs.perigon.io/docs/getting-started
- Article data: https://docs.perigon.io/docs/article-data
- Story data: https://docs.perigon.io/docs/story-data
- Vector endpoint: https://docs.perigon.io/docs/vector-endpoint
- Summarizer: https://perigon.io/docs/api/summarizer
- Wikipedia: https://perigon.io/docs/api/wikipedia
- Pagination: https://perigon.io/docs/api/pagination

Perigon documentation currently uses both /v1/articles/all and /v1/all for the article search surface. It also uses both /v1/sources and /v1/sources/all in different pages. The route aliases should be confirmed against the provider account's current OpenAPI or endpoint reference before implementation.

### Route inventory

| Method and route | Main parameters or body | What it does |
|---|---|---|
| GET /v1/articles/all | q, title, desc, content, url, linkTo, from, to, addDateFrom, addDateTo, refreshDateFrom, refreshDateTo, source, excludeSource, sourceGroup, excludeSourceGroup, category, excludeCategory, topic, excludeTopic, taxonomy, prefixTaxonomy, medium, label, excludeLabel, language, location filters, source geography, entity filters, sentiment filters, sortBy, page, size | Searches and filters individual news articles. |
| GET /v1/stories/all | q, name, source, clusterId, category, topic, taxonomy, from, to, initializedFrom, initializedTo, updatedFrom, updatedTo, minClusterSize, maxClusterSize, sortBy, state, location filters, showDuplicates, page, size | Searches article clusters representing stories or narratives. |
| GET /v1/stories/history | story and time filters documented by Perigon, plus page and size | Returns changes to a story over time. |
| POST /v1/vector/news/all | JSON body containing prompt, showReprints, source/category/topic/date filters, page, size, and scoreThreshold | Performs semantic news search using a prompt and optional structured filters. |
| POST /v1/articles/summarize | JSON body containing the normal article filters plus prompt, maxArticleCount, method, and summarizeFields | Produces an LLM-generated summary over matching articles or clusters. |
| GET /v1/wikipedia/all | q, title, summary, text, reference, id, category, wikiRevisionFrom, wikiRevisionTo, scrapedAtFrom, scrapedAtTo, pageviewsFrom, pageviewsTo, withPageviews, sortBy, page, size | Searches Perigon's Wikipedia data. |
| POST /v1/vector/wikipedia/all | JSON body containing prompt, filter, revision/pageview filters, page, size, and scoreThreshold | Performs semantic search over Wikipedia pages. |
| GET /v1/journalists | q, name, twitter, topic, country, updatedAtFrom, page, size | Searches journalists and authors. |
| GET /v1/journalists/{journalist_id} | journalist_id path parameter | Returns one journalist profile. |
| GET /v1/sources or /v1/sources/all | domain, name, sourceGroup, country, sourceCountry, sourceState, sourceCity, minMonthlyVisits, maxMonthlyVisits, minMonthlyPosts, category, topic, sortBy, paywall, page, size | Searches source and publisher metadata. |
| GET /v1/people/all | name, wikidataId, occupationId, occupationLabel, page, size | Searches people and public biographical metadata. |
| GET /v1/companies/all | name, domain, ticker, wikidataId, page, size | Searches company metadata used to enrich news and entities. |
| GET /v1/topics/all | category and pagination parameters documented by Perigon | Lists or searches Perigon topic metadata. |

### Shared Perigon parameters

- Text: q, title, desc, content, name, and route-specific text fields.
- Time: from and to for article publication time; addDateFrom/To and refreshDateFrom/To for ingestion lifecycle; story routes additionally expose initialized and updated ranges.
- Classification: category, topic, taxonomy, medium, and label.
- Source: source, excludeSource, sourceGroup, excludeSourceGroup, domain, and paywall.
- Location: article and source country, state, city, area, county, latitude, longitude, and distance filters.
- Entity enrichment: people, companies, and related entity identifiers.
- Ordering: sortBy is route-specific, commonly date or relevance for articles and createdAt, updatedAt, count, or relevance for stories.
- Pagination: page is zero-based, size defaults to 10, and size has a documented maximum of 100. The API limits paginated result sets to 10,000 records.

### Article search response

~~~json
{
  "status": 200,
  "numResults": 1,
  "articles": [
    {
      "url": "https://example.com/article",
      "articleId": "article-id",
      "clusterId": "story-id",
      "title": "Example headline",
      "description": "Article description",
      "content": "Article body when available",
      "authorsByline": "Example Author",
      "source": {
        "domain": "example.com",
        "paywall": false,
        "location": {"country": "US"}
      },
      "imageUrl": "https://example.com/image.jpg",
      "country": "US",
      "language": "en",
      "pubDate": "2026-08-10T12:00:00Z",
      "addDate": "2026-08-10T12:05:00Z",
      "refreshDate": "2026-08-10T12:10:00Z",
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

The article object is substantially richer than this representative payload. Content availability depends on the source, licensing, and provider enrichment.

### Story search response

~~~json
{
  "status": 200,
  "numResults": 1,
  "results": [
    {
      "id": "story-id",
      "name": "Example story",
      "summary": "Cluster summary",
      "summaryReference": "https://example.com/reference",
      "keyPoints": ["Key point one"],
      "createdAt": "2026-08-10T12:00:00Z",
      "updatedAt": "2026-08-10T14:00:00Z",
      "initializedAt": "2026-08-10T12:00:00Z",
      "duplicateOf": null,
      "uniqueCount": 12,
      "reprintCount": 30,
      "totalCount": 42,
      "countries": ["US"],
      "topCountries": [{"country": "US", "count": 30}],
      "topics": [],
      "topTopics": [],
      "categories": [],
      "people": [],
      "companies": [],
      "locations": [],
      "taxonomies": [],
      "topTaxonomies": []
    }
  ]
}
~~~

The stories/history route returns story-change records rather than a current story snapshot. The exact history fields are provider-version dependent.

### Vector search response

POST /v1/vector/news/all accepts a prompt and structured filters in JSON. It returns article-style records in a Perigon success envelope, normally with relevance or similarity information in the returned result. POST /v1/vector/wikipedia/all returns semantically ranked Wikipedia page records. The exact score field should be taken from the provider's current schema.

### Summarizer response

~~~json
{
  "status": 200,
  "numResults": 50,
  "summary": "Generated summary text",
  "results": [
    {
      "url": "https://example.com/article",
      "articleId": "article-id",
      "clusterId": "story-id",
      "title": "Example headline"
    }
  ]
}
~~~

The summarize request supports prompt, maxArticleCount, method (ARTICLES or CLUSTERS), and summarizeFields (TITLE, CONTENT, or SUMMARY), in addition to the normal article filters.

### Entity, source, journalist, and Wikipedia responses

Perigon supplemental entity routes generally use this envelope:

~~~json
{
  "status": 200,
  "numResults": 1,
  "results": [
    {
      "id": "provider-id",
      "name": "Example name",
      "description": "Provider metadata",
      "topics": [],
      "categories": []
    }
  ]
}
~~~

The route-specific result fields are:

- Journalists: id, name/fullName, headline, description, title, locations, updatedAt, topTopics, topSources, topCategories, topLabels, topCountries, social handles, and profile URLs.
- Sources: domain, name, alternate names, location, global rank, monthly visits, monthly posts, paywall, categories, and topics.
- People: Wikidata identifier, name, gender, dates of birth/death, occupation, employer, image, and abstract.
- Companies: company identity, domain, ticker or market identifiers, industry, headquarters, and knowledge-graph identifiers.
- Topics: topic name and category metadata.
- Wikipedia: page identifiers, title, summary, text, references, categories, Wikidata identifiers, and optional pageview metrics.

## 4. Cross-provider comparison

| Capability | GDELT Cloud | PredictHQ | Perigon |
|---|---|---|---|
| Primary unit | Event, story, entity, article evidence | Event, place, embedded entity | Article, story cluster, entity, source |
| Search route style | Separate /events, /stories, /entities | One primary /events/ search route plus /count/ | Separate article, story, vector, entity, and summarizer routes |
| Semantic search | search parameter on event/story surfaces | q is text search; event ontology and filters are primary | Dedicated vector POST routes |
| Aggregation | /events/summary and /stories/summary | /events/count/ and /broadcasts/count/ | Result counts; no equivalent general event-summary route in this surface |
| Evidence traversal | story articles route and linked references | Event object includes source/entity data; no equivalent generic story/article traversal in reviewed surface | Article and story routes, with clusterId linking |
| Entity model | Separate entity routes and references | Entities embedded in events and filtered by entity.id | Supplemental entity routes plus enriched article/story fields |
| Time model | Day-oriented event_date filters | Event lifecycle timestamps and comparison operators | Article/story ingestion and publication timestamps |
| Pagination | Cursor | Offset | Page and size |
| Response shape | success/data/pagination | count/results/next/previous | status/numResults/articles or results |
| Record shape | Provider-defined event/story cards | Strongly typed event objects | Rich article/story/entity objects |

## 5. Implications for Espresso route comparison

For route-name and contract comparison, the closest existing patterns are:

- Collection search: GDELT /events and /stories, PredictHQ /v1/events/, and Perigon /v1/articles/all or /v1/stories/all.
- Detail lookup: GDELT /events/{id}, /stories/{id}, and /entities/{id}; PredictHQ uses an id filter on the collection route; Perigon exposes a journalist detail route but commonly returns articles and stories through collections.
- Aggregation: GDELT /events/summary and PredictHQ /events/count/.
- Semantic search: GDELT search and Perigon dedicated vector routes. PredictHQ primarily combines q with structured event filters.
- Related evidence: GDELT story articles and linked references, Perigon clusterId/article-story relationships, and PredictHQ embedded entities and event metadata.
- Response compatibility: Espresso should define stable top-level pagination and count fields while preserving provider-independent event fields such as id, created_at, kind, title/summary, timestamps, geography, entities, sources, and relationships.
- The provider comparison does not justify replacing event identity with event_id or created with created_at. Those names should be normalized at the Espresso contract boundary if the target contract uses id and created_at.
- A flexible digest can remain an Espresso-specific extension, but its stable envelope should expose the common fields needed by all three provider families.

## 6. Source and version notes

This reference summarizes the providers' public documentation and representative payloads. Provider schemas, route aliases, limits, and authentication rules can change. Before implementing a client or declaring exact compatibility, obtain the provider's current OpenAPI/schema for the account and region being used.

Primary documentation:
- GDELT Cloud API v2: https://docs.gdeltcloud.com/api-reference/v2
- GDELT Cloud API guide: https://docs.gdeltcloud.com/API_DOCUMENTATION_GUIDE
- PredictHQ Search Events: https://docs.predicthq.com/api/events/search-events
- PredictHQ Event Counts: https://docs.predicthq.com/api/events/get-event-counts
- PredictHQ Places: https://docs.predicthq.com/api/places
- PredictHQ Entities: https://docs.predicthq.com/getting-started/predicthq-data/entities
- Perigon API overview: https://perigon.io/docs/api/intro
- Perigon Getting Started: https://docs.perigon.io/docs/getting-started
- Perigon Article Data: https://docs.perigon.io/docs/article-data
- Perigon Story Data: https://docs.perigon.io/docs/story-data
- Perigon Vector Endpoint: https://docs.perigon.io/docs/vector-endpoint
- Perigon Summarizer: https://perigon.io/docs/api/summarizer
- Perigon Wikipedia: https://perigon.io/docs/api/wikipedia
- Perigon Pagination: https://perigon.io/docs/api/pagination

