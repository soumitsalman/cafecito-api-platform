# Cafecito API Product Surface and Execution Roadmap

Status: proposed execution instruction  
Prepared: 2026-07-30  
Scope: Beans API, Espresso API, gateway OpenAPI, developer documentation, and MCP-facing contracts

## Relationship to the predecessor

This document is the execution-oriented successor to
[API_CAPABILITY_GAP_ANALYSIS.md](API_CAPABILITY_GAP_ANALYSIS.md). The predecessor
explains the capability and documentation gaps. This document converts those
findings into:

1. an ideal target API surface;
2. a parity trace to comparable services;
3. a staged delivery order;
4. small work packets that can be given to an implementation agent one at a
   time; and
5. explicit separation of code, query, documentation, generated-spec, and
   operations work.

The predecessor remains authoritative background and must not be replaced by
this roadmap.

## Fixed constraints

Every agent executing a work packet from this document must preserve these
constraints.

- Keep the gateway product prefixes `/beans` and `/espresso`.
- Do not combine unrelated work packets.
- Keep the backend schema and ingestion pipelines outside the API repository
  unless a packet explicitly says otherwise.
- Treat the ongoing Beans UUID backfill and `id` primary-key conversion as an
  external dependency. Design article routes as though `beans.id` is the stable
  primary key, but do not plan or implement that conversion here.
- Treat the ongoing Espresso action-storage effort as an external dependency.
  `/actions` is a required product route, but this roadmap does not define how
  action records are produced or stored.
- `SAME_AS` is bidirectional. Clients must never supply or interpret edge
  direction for it.
- `DERIVED_FROM` is directional internally. Clients must select an intuitive
  typed subresource, such as `/signals/{id}/events`, and must never supply
  `from_id`, `to_id`, or a direction flag.
- Preserve existing public routes as compatibility aliases until a measured
  deprecation period is complete.
- A route or Swaggo annotation change is incomplete until the generated backend
  Swagger and the matching gateway OpenAPI file are synchronized.
- All APIs remain read-only from the client's perspective.

## Evidence from the read-only data scan

The Neon project `cafecito-apps-v2` was inspected read-only. No branch, schema,
row, role, or project setting was changed.

### Beans data that the API can already exploit

- `beansdb.beans` contains about 2.0 million records across `news`, `blog`,
  `post`, and `site`.
- Articles already have URL, title, source, content, timestamps, semantic
  vectors, categories, sentiments, regions, entities, and full-text tags, with
  varying field completeness.
- `related_beans` contains about 1.8 million URL-to-URL relationships.
- `chatters` contains about 2.8 million social or forum observations.
- `trend_aggregates` contains about 338,000 aggregate records.
- `publishers` contains about 28,000 source records, although publisher UUID
  coverage and optional metadata completeness are low.

This is enough to expose article search, latest content, headlines, related
articles, source discovery, and propagation today. Stable story resources,
strong taxonomy filtering, and analytical summaries need additional read
models or query work.

### Espresso data that the API can already exploit

- `espressodb.sips` contains about 826,000 records.
- The stored kinds currently include curated `event`, `signal`, and
  source-level `event:news`, `event:blog`, `event:post`, and `event:site`.
- No stored `action` rows were observed. This agrees with the separate action
  workstream and does not change the requirement for an `/actions` route.
- Curated event digests already contain fields such as briefing, event type,
  impact level, actions, macro context, future outlook, affected regions, and
  named companies, people, and products.
- Signal digests already contain briefing, linked events, forecast, impacts,
  drivers, impacted domains, tags, and sometimes confidence.
- `sources` contains about 47,000 records with stable UUIDs, though descriptive
  metadata is incomplete.
- `relations` contains about 61 million edges: approximately 98% `SAME_AS` and
  2% `DERIVED_FROM`.
- Sampled `DERIVED_FROM` edges are stored from a signal to its supporting
  source event. The public API should therefore expose both
  `/signals/{id}/events` and `/events/{id}/signals`, translating the requested
  resource into the correct outgoing or incoming lookup internally.

### Data-contract risks to address before claiming strong guarantees

- Both databases contain sentinel or invalid year-0001 timestamps.
- Beans taxonomy values contain casing and naming variants, including
  Title Case and snake_case forms for the same apparent concept.
- Beans `restricted_content` is not a simple complete Boolean field. The API
  needs a documented content-availability contract rather than implying that
  `full_content=true` always returns unrestricted content.
- Some Espresso relation endpoints no longer resolve to a current `sips` row.
  The API must define whether to omit unresolved edges or expose them as
  unavailable references.
- Espresso source descriptions, site names, and feeds are optional.
- Social observations measure different engagement units by forum. The API
  should use neutral terms such as `mention_count` and expose forum-specific
  metrics rather than describing all observations as shares.

## Product-wide route and response rules

These rules are the target contract for both products.

### Route rules

1. Use plural resource nouns for collections: `/articles`, `/events`,
   `/signals`, `/actions`, and `/sources`.
2. Use `/{id}` for one stable resource.
3. Search and filtering belong on a collection `GET`; `/search` is retained
   only as a deprecated compatibility alias.
4. Relationships are named subresources whose result type is apparent:
   `/signals/{id}/events`, not `/related/derived_from?direction=out`.
5. Use `from` and `to` for time bounds and document both as RFC 3339 timestamps.
6. Use `sort` with a documented enum. Never silently change sort semantics when
   `q` is supplied.
7. Use cursor pagination for newly introduced high-volume routes. Continue
   accepting `offset` on legacy routes during migration.
8. Use `include` or a named view for optional expensive fields. Do not overload
   a Boolean until its exact behavior is defined.

### Response rules

Every canonical JSON collection should converge on:

```json
{
  "data": [],
  "pagination": {
    "next_cursor": null,
    "limit": 20
  },
  "meta": {
    "as_of": "2026-07-30T12:00:00Z",
    "query_mode": "semantic"
  }
}
```

Every canonical error should converge on:

```json
{
  "error": {
    "code": "invalid_parameter",
    "message": "sort must be one of: relevance, published_at, trend_score",
    "parameter": "sort",
    "request_id": "..."
  }
}
```

Rules for agent-friendly responses:

- IDs and timestamps are never omitted.
- Field ordering in `response_type=text` is deterministic.
- Text responses include a stable record delimiter and escaping rules.
- Empty results are successful responses with `data: []`, not ambiguous
  messages.
- Every response documents freshness, pagination, and whether expensive fields
  were included.
- OpenAPI schemas use concrete objects rather than unconstrained JSON wherever
  the digest shape is stable.

## Ideal Beans product surface

The canonical Beans surface should be small enough to understand from route
names alone while retaining advanced search and propagation features.

| Canonical route | Purpose and key inputs | Existing Cafecito route/capability | Comparable-service parity | Delivery |
|---|---|---|---|---|
| `GET /beans/articles` | Search or browse articles. Inputs: `q`, typed filters, `from`, `to`, `sort`, pagination, and content view. Default is recent content; `q` enables semantic relevance. | Consolidates `/beans/articles/search` and `/beans/articles/latest`. | [NewsAPI Everything](https://newsapi.org/docs/endpoints/everything), Perigon article and vector search, Event Registry article search. | Stage 0 alias, Stage 2 full contract |
| `GET /beans/articles/{id}` | Retrieve one stable article with explicit content-availability metadata. | Can initially delegate to the existing ID filter after the external UUID dependency is ready. | Perigon and Event Registry article resources. | Stage 2, dependency-gated |
| `GET /beans/headlines` | A deliberately narrow, high-value current feed. Inputs should include region, category, source, and `from`; default sort is trend/relevance within a documented recency window. | Renames `/beans/articles/top-headlines`. | [NewsAPI Top Headlines](https://newsapi.org/docs/endpoints/top-headlines). | Stage 0 alias |
| `GET /beans/articles/{id}/related` | Return semantically or graph-related articles, with a relationship or similarity explanation. | Packages `related_beans` and current related-search behavior around article ID. | Perigon story membership, Event Registry event clusters, GDELT story articles. | Stage 2 |
| `GET /beans/articles/{id}/propagation` | Return publisher coverage and social/forum mentions for one article. The server resolves the article URL. | Repackages current propagation GET behavior. | Event Registry social signals and Perigon/GDELT story coverage concepts. | Stage 2 |
| `POST /beans/articles/propagation` | Batch form for up to the documented maximum number of IDs or URLs. This is an advanced workload route, not the first route in the quickstart. | Repackages current propagation POST capability. | Batch analysis patterns in commercial news-intelligence APIs; no parity claim should imply resource creation. | Stage 0 documentation, Stage 2 IDs |
| `GET /beans/sources` | Browse and filter publishers by category, language, region, or domain when available. | Renames `/beans/sources` only if public shape needs normalization; backend already has source listing/resolution. | [NewsAPI Sources](https://newsapi.org/docs/endpoints/sources), [Perigon Sources](https://perigon.io/docs/api/sources). | Stage 0 documentation, Stage 2 filters |
| `GET /beans/sources/{source}` | Retrieve source metadata using the stable source string until publisher UUID coverage supports an ID migration. | Uses current publisher data. | NewsAPI source IDs and Perigon source records. | Stage 2 |
| `GET /beans/taxonomies/categories` | Discover valid category filters and canonical values. | Alias for `/beans/tags/categories`. | Event Registry category suggestions; Perigon topics. | Stage 0 alias |
| `GET /beans/taxonomies/entities` | Discover indexed entities and canonical values. | Alias for `/beans/tags/entities`. | GDELT entities; Event Registry concept suggestions. | Stage 0 alias |
| `GET /beans/taxonomies/regions` | Discover valid geographic filters. | Alias for `/beans/tags/regions`. | Event Registry location suggestions; news-source country filters. | Stage 0 alias |
| `GET /beans/taxonomies/sentiments` | Discover or document the sentiment vocabulary and filtering rules. | Data exists, but no matching public discovery route is currently prominent. | Alpha Vantage news sentiment and Event Registry sentiment filtering. | Stage 1 capability, Stage 2 route |
| `GET /beans/stories` | Search durable cross-publisher story clusters: one real-world narrative represented by many publisher articles. Filter by story-level time range, source count, entities, regions, or topic; return one row per story, not one row per article. | Related-article graph exists, but no durable story identity/read model exists. | Perigon stories, GDELT stories, Event Registry events. | Stage 2 after Stage 1 read model |
| `GET /beans/stories/{id}` | Retrieve cluster summary, representative article, time range, sources, and entities. | Requires a story identity and cluster summary capability. | Perigon story detail and GDELT story detail. | Stage 2 after Stage 1 |
| `GET /beans/stories/{id}/articles` | Retrieve the articles and coverage timeline for a story. | Can eventually package `related_beans` behind a stable story. | [GDELT story articles](https://docs.gdeltcloud.com/api-reference/v2). | Stage 2 after Stage 1 |
| `GET /beans/analytics/trends` | Explain what is gaining attention, not merely how many records were ingested. Return time-bucketed trend scores for articles or stories, grouped by an allowlisted dimension such as category, entity, region, source, or topic. | `trend_aggregates` exists but lacks a clear public analytical contract. | Event Registry aggregates, Perigon analytics tools, GDELT summaries. | Stage 2 |
| `GET /beans/analytics/article-volume` | Answer “how much coverage was published?” Return counts of articles over time, grouped into documented buckets and optionally filtered by source, category, region, entity, or content type. This is volume, not popularity or trend strength. | Requires aggregate query packaging. | Perigon article counts and GDELT summary routes. | Stage 2 |

### Beans stories versus articles

Use `GET /beans/articles` when the user wants individual documents: headlines,
titles, URLs, summaries/content, authors, publishers, publication timestamps,
or article-level filters. The result contains one row per article, so two
publishers covering the same event can produce two or more rows.

Use `GET /beans/stories` when the user wants the underlying narrative across
publishers: a canonical story identity, representative article, coverage
timeline, source count, related entities, and the complete article membership.
The result contains one row per story. It is therefore the right starting point
for “how did this story propagate?” and the wrong starting point for “give me
the latest 20 article URLs.”

Example article search:

```json
{
  "request": "GET /beans/articles?q=central-bank-rate-cut&sort=relevance&limit=2",
  "data": [
    {
      "id": "8d6c...",
      "title": "Central bank signals a possible rate cut",
      "source": "example-news",
      "published_at": "2026-07-30T13:10:00Z",
      "story_id": "story_42"
    },
    {
      "id": "21ab...",
      "title": "Markets react to the central bank announcement",
      "source": "example-business",
      "published_at": "2026-07-30T13:22:00Z",
      "story_id": "story_42"
    }
  ],
  "pagination": { "next_cursor": "..." }
}
```

Example story search:

```json
{
  "request": "GET /beans/stories?q=central-bank-rate-cut&limit=2",
  "data": [
    {
      "id": "story_42",
      "headline": "Central bank signals a possible rate cut",
      "first_seen_at": "2026-07-30T13:10:00Z",
      "last_seen_at": "2026-07-30T16:45:00Z",
      "article_count": 27,
      "source_count": 11,
      "top_entities": ["Example Central Bank", "Exampleland"],
      "representative_article_id": "8d6c..."
    }
  ],
  "pagination": { "next_cursor": null }
}
```

The `story_id` shown on an article is a navigation link into the story
surface; it is not a substitute for the article ID. Story identity, membership,
and merge/split behavior remain Stage 1 read-model work.

### Beans sort and content contract

`GET /beans/articles` should use these explicit sort values:

- `published_at`: newest publication time first;
- `relevance`: semantic/text relevance, valid when `q` is supplied;
- `trend_score`: current social/coverage trend score;
- `collected_at`: ingestion recency, primarily for operational consumers.

Do not keep separate canonical "latest" and "search" implementations. They are
views over the same article collection. `/headlines` remains separate because
it represents a curated product promise rather than a sort synonym.

Replace the ambiguous `full_content` promise with a documented content view,
for example `view=summary|content`. Each returned article should include:

```json
{
  "content": null,
  "content_access": {
    "status": "restricted",
    "reason": "publisher_restriction"
  }
}
```

The exact enum may change during contract review, but the API must distinguish
missing, restricted, unavailable, and included content.

## Ideal Espresso product surface

Espresso should present actions, events, and signals as first-class resources.
The word "sip" can remain an internal persistence abstraction; clients should
not need it.

| Canonical route | Purpose and relation behavior | Existing Cafecito route/capability | Comparable-service parity | Delivery |
|---|---|---|---|---|
| `GET /espresso/actions` | Search and browse atomic activities. Must exist in the target contract even though storage is handled by a separate workstream. | No current route or stored action rows. | Economic-calendar observations, market-news facts, and event evidence in Trading Economics/GDELT-like products. | Stage 2, external dependency |
| `GET /espresso/actions/{id}` | Retrieve one action with source, time, entities, and typed facts. | Depends on action read availability, not planned storage work. | Individual observation/fact resources in structured intelligence APIs. | Stage 2, external dependency |
| `GET /espresso/events` | Search and browse self-contained events with typed filters, source, impact, entities, and time bounds. | Current `/espresso/events`. | [GDELT Events](https://docs.gdeltcloud.com/api-reference/v2), Event Registry events. | Stage 0 documentation, Stage 2 contract |
| `GET /espresso/events/{id}` | Retrieve one event and its structured digest. | Can be implemented using the existing ID-filtered event query before adding a specialized repository method. | GDELT event detail. | Stage 0/2 |
| `GET /espresso/signals` | Search and browse cross-event intelligence with forecast, drivers, impacts, confidence when present, and time bounds. | Current `/espresso/signals`. | Perigon Signal Insights and derived analytical outputs in market-intelligence APIs. | Stage 0 documentation, Stage 2 contract |
| `GET /espresso/signals/{id}` | Retrieve one signal and its structured digest. | Can be implemented using the existing ID-filtered signal query. | Signal/insight detail patterns in intelligence platforms. | Stage 0/2 |
| `GET /espresso/actions/{id}/similar` | Return `SAME_AS` peers. Query both stored edge orientations. | Generic `/espresso/related/same_as`. | Duplicate-event and cluster-member retrieval in GDELT/Event Registry. | Stage 2, action dependency |
| `GET /espresso/events/{id}/similar` | Return `SAME_AS` peers without a direction parameter. | Generic `/espresso/related/same_as` already queries both orientations. | Event clustering/deduplication in GDELT and Event Registry. | Stage 1 query contract, Stage 2 route |
| `GET /espresso/signals/{id}/similar` | Return `SAME_AS` peers without a direction parameter. | Generic relation capability. | Similar-insight discovery patterns. | Stage 2 |
| `GET /espresso/signals/{id}/events` | Return supporting events by following outgoing `DERIVED_FROM` edges internally. | Current generic relation query incorrectly treats `DERIVED_FROM` as undirected. | Evidence articles/events attached to GDELT stories and intelligence signals. | Stage 1 direction fix, Stage 2 route |
| `GET /espresso/events/{id}/signals` | Return signals derived from this event by following incoming `DERIVED_FROM` edges internally. | Requires the inverse directional query. | Reverse evidence-to-insight navigation in intelligence graphs. | Stage 1 direction fix, Stage 2 route |
| `GET /espresso/events/{id}/actions` | Return actions summarized by an event using the correct internally stored direction. | Depends on action data and relation semantics. | Event evidence and constituent fact navigation. | Stage 2, action dependency |
| `GET /espresso/actions/{id}/events` | Return events that incorporate the action using the inverse internal lookup. | Depends on action data and relation semantics. | Reverse fact-to-event navigation. | Stage 2, action dependency |
| `GET /espresso/sources` | Browse source records and their available metadata. | Table and stable IDs exist; repository/list route is needed. | NewsAPI/Perigon source discovery and GDELT article evidence. | Stage 1 capability, Stage 2 route |
| `GET /espresso/sources/{id}` | Retrieve one source. Missing optional metadata remains nullable. | Stable source UUIDs already exist. | Source-detail routes in news intelligence APIs. | Stage 1 capability, Stage 2 route |
| `GET /espresso/taxonomies/tags` | Discover canonical filter values and counts when inexpensive. | Renames `/espresso/tags`. | Event Registry suggestions and Perigon topic discovery. | Stage 0 alias |
| `GET /espresso/analytics/events` | Summarize the event layer: how many events occurred, where, which sources/types/tags/entities are active, and how impact levels change over time. It is an aggregate view over `/events`, not a replacement for retrieving event digests. | Digest data exists; typed extraction/aggregation is needed. | [GDELT event summaries](https://docs.gdeltcloud.com/api-reference/v2), Perigon analytics. | Stage 2 |
| `GET /espresso/analytics/signals` | Summarize the intelligence layer: how many signals are active, their impact/confidence distribution, dominant drivers or domains, and how forecasts change across a time window. It is an aggregate view over `/signals`, not a list of signal records. | Digest data exists; typed extraction/aggregation is needed. | Intelligence insight dashboards and aggregation APIs. | Stage 2 |

### Analytics route examples

The analytics routes answer distribution and change questions. They should
return compact aggregate rows with a documented bucket, dimension, and
freshness timestamp; they should not return the full event, signal, or article
digests.

Example Beans article-volume response:

```json
{
  "request": "GET /beans/analytics/article-volume?from=2026-07-01T00:00:00Z&to=2026-07-30T00:00:00Z&group_by=category&bucket=day&categories=technology",
  "data": [
    { "bucket_start": "2026-07-29T00:00:00Z", "group": "technology", "article_count": 1842 },
    { "bucket_start": "2026-07-30T00:00:00Z", "group": "technology", "article_count": 1961 }
  ],
  "meta": {
    "measure": "published_article_count",
    "group_by": "category",
    "bucket": "day",
    "as_of": "2026-07-30T16:00:00Z"
  }
}
```

Example Beans trend response:

```json
{
  "request": "GET /beans/analytics/trends?from=2026-07-01T00:00:00Z&to=2026-07-30T00:00:00Z&group_by=entity&limit=2",
  "data": [
    {
      "group": "Example Central Bank",
      "trend_score": 0.87,
      "article_count": 12940,
      "mention_count": 48320,
      "direction": "rising"
    },
    {
      "group": "Exampleland",
      "trend_score": 0.71,
      "article_count": 8940,
      "mention_count": 27110,
      "direction": "steady"
    }
  ],
  "meta": {
    "measure": "attention",
    "group_by": "entity",
    "as_of": "2026-07-30T16:00:00Z"
  }
}
```

Example Espresso event analytics response:

```json
{
  "request": "GET /espresso/analytics/events?from=2026-07-01T00:00:00Z&to=2026-07-30T00:00:00Z&group_by=impact_level&bucket=week",
  "data": [
    { "bucket_start": "2026-07-27T00:00:00Z", "group": "high", "event_count": 318 },
    { "bucket_start": "2026-07-27T00:00:00Z", "group": "medium", "event_count": 1240 }
  ],
  "meta": {
    "kind": "event",
    "group_by": "impact_level",
    "bucket": "week",
    "as_of": "2026-07-30T16:00:00Z"
  }
}
```

Example Espresso signal analytics response:

```json
{
  "request": "GET /espresso/analytics/signals?from=2026-07-01T00:00:00Z&to=2026-07-30T00:00:00Z&group_by=domain",
  "data": [
    {
      "group": "markets",
      "signal_count": 86,
      "high_impact_count": 21,
      "average_confidence": 0.78,
      "top_drivers": ["rates", "inflation"]
    },
    {
      "group": "technology",
      "signal_count": 63,
      "high_impact_count": 14,
      "average_confidence": 0.74,
      "top_drivers": ["chips", "cloud"]
    }
  ],
  "meta": {
    "kind": "signal",
    "group_by": "domain",
    "as_of": "2026-07-30T16:00:00Z"
  }
}
```

The field names above are target-contract examples. Before implementation,
each aggregate must define its authoritative source field, null behavior,
normalization rules, freshness, maximum time range, and whether counts are
exact or estimated.

### Espresso relation contract

The generic route `/espresso/related/{relationship}` should be deprecated after
typed subresources reach parity.

The public contract is:

| User asks for | Canonical route | Internal edge behavior |
|---|---|---|
| records equivalent to an event | `/events/{id}/similar` | Match `SAME_AS` where the ID is either endpoint. |
| events supporting a signal | `/signals/{id}/events` | Follow signal-to-event `DERIVED_FROM`. |
| signals using an event | `/events/{id}/signals` | Follow event-to-signal inverse lookup. |
| actions summarized by an event | `/events/{id}/actions` | Follow the canonical event-to-action direction defined with the action workstream. |
| events using an action | `/actions/{id}/events` | Perform the inverse lookup internally. |

No public request accepts `direction`, `from_id`, or `to_id`. If a relation
points to a missing sip, the default route omits it and increments an internal
orphan metric. An optional diagnostic representation should not be added until
there is a concrete consumer.

## AI-agent and MCP documentation contract

REST route names should be resource-oriented. MCP tool names should be
task-oriented. The MCP layer should not mechanically expose every HTTP
operation as an equally prominent tool.

### Recommended MCP tools

Beans:

- `search_articles`: semantic/text search plus scalar filters.
- `get_article`: retrieve one article by UUID.
- `get_headlines`: retrieve the current curated feed.
- `find_related_articles`: retrieve related coverage for one article.
- `trace_article_propagation`: retrieve publisher coverage and social/forum
  mentions.
- `list_news_sources`: discover valid sources.
- `discover_article_filters`: discover categories, entities, regions, and
  sentiment values.
- `search_stories` and `get_story`: add only when stable story resources exist.
- `get_article_trends`: add only when the analytical route has a stable
  contract.

Espresso:

- `search_actions` and `get_action`: publish when action reads are available.
- `search_events` and `get_event`.
- `search_signals` and `get_signal`.
- `find_similar_events`: hides `SAME_AS` orientation.
- `get_signal_evidence`: returns the events supporting a signal.
- `get_event_signals`: returns the signals derived from an event.
- `get_event_actions`: add when actions are available.
- `list_intelligence_sources`.
- `discover_intelligence_tags`.

This follows the task-oriented tool catalog used by
[Perigon's MCP server](https://perigon.io/docs/api/mcp), while the REST layer
retains predictable resources.

### Required documentation for each route and tool

Every route and MCP tool description must answer:

1. What user question should invoke this operation?
2. What should not invoke it?
3. Which fields are searched by `q`?
4. Is `q` semantic, lexical, or hybrid?
5. What is the default time window and sort?
6. Which timestamp does `from` or `to` filter?
7. What are the maximum page size and pagination method?
8. Which fields may be absent?
9. What does an empty result mean?
10. What is the cost or latency implication of content, embeddings, relations,
    or broad time windows?
11. Which operation should be called next for detail or evidence?
12. Is the operation read-only and idempotent?

### Agent decision guide

The docs must include a compact decision table similar to:

| User intent | First tool | Follow-up |
|---|---|---|
| "What was published about X?" | Beans `search_articles` | `get_article` or `find_related_articles` |
| "What is breaking now?" | Beans `get_headlines` | `trace_article_propagation` |
| "How widely did this story spread?" | Beans `trace_article_propagation` | `find_related_articles` |
| "What happened to company X?" | Espresso `search_events` | `get_event`, then `get_event_signals` |
| "What broader conclusion is emerging?" | Espresso `search_signals` | `get_signal_evidence` |
| "What observations compose this event?" | Espresso `get_event_actions` | `get_action` |

Also provide:

- one minimal request;
- one filtered request;
- one semantic-search request;
- one pagination continuation;
- one empty result;
- one validation error;
- one JSON response;
- one compact-text response; and
- one complete multi-tool workflow per product.

For MCP descriptors, declare read-only/idempotent annotations where supported,
but do not treat annotations as an authorization mechanism. The MCP project's
[tool-annotation guidance](https://blog.modelcontextprotocol.io/posts/2026-03-16-tool-annotations/)
describes these as behavioral hints.

## Priority model

- P0: correctness, route comprehension, contract alignment, or migration safety.
- P1: significant parity or agent-usability improvement.
- P2: analytical depth, operational maturity, or optimization after the core
  surface is stable.

Change types used below:

- **DOC**: authored Markdown/MDX or examples only.
- **SPEC**: gateway OpenAPI or Swaggo annotations.
- **CODE**: Go/TypeScript behavior without a database query change.
- **QUERY**: SQL/repository behavior; still no schema mutation unless stated.
- **GENERATED**: regenerated Swagger files.
- **OPS**: metrics, dashboards, alerts, or runbooks.
- **TEST**: unit, contract, integration, or compatibility tests.

## How to execute this roadmap piece by piece

Give an agent exactly one work-packet ID and this document. The agent must:

1. restate the packet's objective and exclusions;
2. inspect the named artifacts and current tests;
3. make only the listed change types;
4. keep existing public routes working unless the packet explicitly authorizes
   a removal;
5. update Swaggo, generated Swagger, and gateway OpenAPI together for any route
   contract change;
6. add tests only under `apis/beans/tests/` or `apis/espresso/tests/`;
7. run the packet's acceptance checks; and
8. report changed files, validation evidence, remaining dependencies, and any
   discovered contract conflict.

An agent must stop rather than silently expand a packet into ingestion, schema,
authentication, billing, or unrelated UI work.

# Stage 0 — Use current capabilities better

Goal: improve comprehension and consistency now, without a schema change or a
major SQL/query redesign.

## Packet S0-01 — Freeze the canonical route contract

Priority: P0  
Product surface: Beans + Espresso  
Change type: DOC  
Primary artifact: a new route-contract ADR under `docs/`

Instructions:

1. Copy the canonical route tables and product-wide rules from this roadmap
   into a concise ADR.
2. Mark every route as canonical, compatibility alias, dependency-gated, or
   future.
3. Record the `SAME_AS` and `DERIVED_FROM` rules verbatim.
4. Record that Beans article IDs and Espresso action storage are external
   dependencies.
5. Do not edit route code in this packet.

Acceptance:

- Every current public route maps to a target route or an explicit retained
  compatibility route.
- The ADR contains no plan for the Beans ID conversion or action ingestion.
- A reviewer can determine relation direction without reading SQL.

## Packet S0-02 — Correct and align the current API specifications

Priority: P0  
Product surface: Beans + Espresso  
Change type: SPEC, GENERATED, TEST  
Artifacts:

- `apis/beans/router/routes.go`
- `apis/beans/docs/`
- `config/beans.oas.json`
- `apis/espresso/router/routes.go`
- `apis/espresso/docs/`
- `config/espresso.oas.json`

Instructions:

1. Audit current parameters, defaults, maximums, response status codes, and
   response schemas against the actual handlers.
2. Correct inaccurate descriptions, including any Beans response field that is
   documented but not selected by its query.
3. Give every operation a stable, task-oriented `operationId`.
4. Add parameter enums and examples wherever behavior already has a finite
   set.
5. Regenerate backend Swagger and update the gateway OpenAPI in the same
   change.
6. Do not add a new database query.

Acceptance:

- Generated Swagger and gateway OpenAPI describe the same paths and parameters,
  except for the expected gateway product prefix.
- JSON specs parse and the existing service tests pass.
- No response schema promises a field the handler cannot return.

## Packet S0-03 — Publish route aliases backed by existing handlers

Priority: P0  
Product surface: Beans + Espresso  
Change type: CODE, SPEC, GENERATED, TEST, DOC  
Artifacts: service routers, both OpenAPI pairs, both how-to pages

Add compatibility-safe aliases:

- `GET /beans/articles` → existing article search/list capability;
- `GET /beans/headlines` → existing top-headlines handler;
- `GET /beans/taxonomies/categories`;
- `GET /beans/taxonomies/entities`;
- `GET /beans/taxonomies/regions`;
- `GET /espresso/taxonomies/tags`;
- `GET /espresso/events/{id}` → existing ID-filtered event capability;
- `GET /espresso/signals/{id}` → existing ID-filtered signal capability.

Instructions:

1. Reuse current handler and repository behavior; do not duplicate SQL.
2. Define how zero and multiple results are handled for detail aliases.
3. Retain every old route and label it as a compatibility alias in the specs.
4. For Beans `/articles`, document the initial default behavior exactly as
   implemented; do not claim the final sort contract until Stage 1.

Acceptance:

- Old and new routes return equivalent records for equivalent requests.
- Detail routes return one record or a documented not-found error.
- Tests prove alias parity.

## Packet S0-04 — Make compact text deterministic and identifiable

Priority: P0  
Product surface: Espresso  
Change type: CODE, SPEC, TEST, DOC  
Artifacts:

- `apis/espresso/router/types.go`
- `apis/espresso/tests/`
- `apis/espresso/router/routes.go`
- Espresso Swagger, gateway OpenAPI, and how-to

Instructions:

1. Include `id`, exact timestamp, resource kind, and a stable record delimiter.
2. Use a fixed field order.
3. Define escaping for newlines and delimiters.
4. Remove ambiguous or duplicate labels.
5. Preserve `response_type=json` as the default.
6. Do not alter database queries.

Acceptance:

- The same record renders byte-for-byte consistently across runs.
- An agent can retrieve a detail route using the ID in the text response.
- Tests cover missing optional digest fields and embedded newlines.

## Packet S0-05 — Rewrite the AI/MCP usage guides around user intent

Priority: P0  
Product surface: Beans + Espresso + MCP  
Change type: DOC, SPEC  
Artifacts:

- `docs/pages/howtos/beans-howto.mdx`
- `docs/pages/howtos/espresso-howto.mdx`
- `docs/pages/howtos/mcp-howto.mdx`
- `docs/pages/api-overview.mdx`
- gateway OpenAPI descriptions

Instructions:

1. Add the agent decision guide and route-selection rules from this roadmap.
2. Explain semantic search versus scalar filtering and when to combine them.
3. Explain latest versus headlines versus trending without using the terms as
   interchangeable synonyms.
4. Explain event versus signal versus action and clearly mark actions as
   dependency-gated until available.
5. Document relation traversal with user questions, never graph orientation.
6. Add exact defaults, limits, time semantics, field-availability notes,
   examples, empty responses, and errors.
7. Add one multi-step MCP workflow for each product.

Acceptance:

- A reader can select the correct first tool without inspecting OpenAPI.
- Every example uses a currently available route or is visibly labeled
  "planned".
- No example exposes `from_id`, `to_id`, or relation direction.

# Stage 1 — Build the underlying capabilities

Goal: implement the query and contract foundations required by the ideal
surface. This stage may change repository queries but should avoid schema
changes unless separately approved.

## Packet S1-01 — Correct Espresso relation semantics

Priority: P0  
Product surface: Espresso relations  
Change type: QUERY, CODE, TEST  
Artifacts:

- `apis/espresso/db/pgcupboard.go`
- `apis/espresso/db/types.go`
- `apis/espresso/router/routes.go`
- `apis/espresso/tests/`

Instructions:

1. Split the current generic relation lookup into explicit repository
   operations:
   - bidirectional `SAME_AS` neighbors;
   - outgoing `DERIVED_FROM` children/evidence;
   - incoming `DERIVED_FROM` parents/consumers.
2. Keep `SAME_AS` symmetric regardless of stored orientation.
3. Preserve `DERIVED_FROM` direction.
4. Define and test behavior for missing related sip rows.
5. Do not publish final typed routes in this packet; expose capability through
   tests and internal interfaces first.

Acceptance:

- Reversing a stored `SAME_AS` edge does not change the public neighbor set.
- Reversing a `DERIVED_FROM` lookup changes the result as expected.
- Tests model signal-to-event and event-to-signal navigation.
- The generic legacy route remains compatible until Stage 2.

## Packet S1-02 — Establish stable detail retrieval

Priority: P0  
Product surface: Beans articles; Espresso actions/events/signals  
Change type: QUERY, CODE, TEST  
Artifacts: both repository packages and both test directories

Instructions:

1. Add a single-record repository method for Espresso by kind and UUID.
2. Prepare a single-record Beans method by UUID, but merge or release it only
   after the external Beans ID workstream declares the primary-key contract
   ready.
3. Define not-found separately from query failure.
4. Return stable IDs in all projections.
5. Add an action retrieval interface compatible with the agreed action read
   model; do not implement storage or ingestion.

Acceptance:

- Event and signal detail retrieval performs a bounded indexed lookup.
- Beans detail work has an explicit dependency check and no fallback to
  ambiguous URLs after UUID cutover.
- Action route packaging can consume the interface without knowing ingestion
  details.

## Packet S1-03 — Define time, sort, and pagination semantics

Priority: P0  
Product surface: Beans + Espresso collections  
Change type: QUERY, CODE, SPEC, TEST, DOC  
Artifacts: repositories, router validation/types, OpenAPI pairs, API guides

Instructions:

1. Define the canonical timestamp for each resource:
   - Beans publication time versus ingestion time;
   - Espresso event/action occurrence time versus record creation time.
2. Exclude or explicitly quarantine sentinel year-0001 values from normal
   time-range results.
3. Add `to`.
4. Add stable sort enums and deterministic tie-breakers.
5. Design cursor pagination from the sort tuple and stable ID.
6. Retain offset pagination on compatibility routes.

Acceptance:

- Two pages cannot duplicate or skip a record when the underlying dataset is
  unchanged.
- Time filters state which database concept they constrain.
- Sentinel timestamps cannot appear as plausible current data.
- Invalid cursor, sort, or time inputs return the standard error envelope.

## Packet S1-04 — Normalize filter vocabularies at the API boundary

Priority: P1  
Product surface: Beans taxonomies; Espresso tags and digest filters  
Change type: QUERY, CODE, TEST, DOC  
Artifacts: repository filter builders, router validation, taxonomy docs

Instructions:

1. Inventory duplicate Beans category/entity/region values caused by casing,
   spacing, or snake_case differences.
2. Choose canonical external values and document alias resolution.
3. Normalize input at the API boundary while preserving stored source data.
4. Define typed Espresso filters for stable digest fields such as impact level,
   event type, company, region, and domain.
5. Use allowlisted JSON paths and values; do not accept arbitrary SQL/JSON path
   expressions.
6. Do not perform an ingestion rewrite in this packet.

Acceptance:

- Equivalent documented aliases produce equivalent result sets.
- Taxonomy discovery returns canonical values.
- Unknown values produce helpful validation or a clearly documented empty
  result.
- Typed Espresso filters have indexes or measured query plans before release.

## Packet S1-05 — Add source repository capabilities

Priority: P1  
Product surface: Beans + Espresso sources  
Change type: QUERY, CODE, TEST, DOC  
Artifacts: source types/repositories and tests

Instructions:

1. Add bounded source list and detail methods.
2. Beans uses the stable `source` string until publisher UUID coverage supports
   a separate migration.
3. Espresso uses the existing source UUID.
4. Make incomplete metadata nullable and document field coverage.
5. Support only filters backed by reliable stored values.

Acceptance:

- List and detail results use stable identifiers.
- Missing description, site name, favicon, or RSS feed is not treated as a
  server error.
- Query plans are bounded and pagination is deterministic.

## Packet S1-06 — Create the story and analytical read models

Priority: P1  
Product surface: Beans stories/analytics; Espresso analytics  
Change type: QUERY, CODE, TEST, DOC  
Artifacts: design ADR first; repository implementation only after review

Instructions:

1. Define a durable Beans story identity independent of one member URL.
2. Define representative article, cluster start/end, source count, article
   count, and merge/split behavior.
3. Determine whether `related_beans` plus derived summaries is sufficient or a
   maintained read model is necessary.
4. Define allowlisted dimensions and bucket sizes for Beans and Espresso
   analytics.
5. Measure plans on realistic cardinalities before route work begins.

Acceptance:

- Story IDs remain stable when a new article joins a cluster.
- Analytics have explicit freshness and maximum-range guarantees.
- No route is published from this packet.
- Any required schema or ingestion work becomes a separately approved
  dependency, not an implicit change.

## Packet S1-07 — Close external dependency gates

Priority: P0 before dependent Stage 2 packets  
Product surface: Beans detail; Espresso actions  
Change type: DOC, TEST; CODE only for integration adapters  
Artifacts: dependency checklist and contract tests

Instructions:

1. For Beans, verify externally that all public article records have UUIDs,
   `id` is the stable primary key, and uniqueness/index expectations hold.
2. For Espresso, obtain the action read contract: kind value, ID, timestamps,
   digest fields, source, tags, embedding behavior, and relation orientation.
3. Add read-only contract tests against those declared shapes.
4. Do not implement either external workstream.

Acceptance:

- The checklist identifies an owner and readiness signal for each dependency.
- Stage 2 detail/action route packets can run without guessing storage
  semantics.

# Stage 2 — Package capabilities into the canonical routes

Goal: publish the ideal resource surface, migrate callers safely, and expose a
small task-oriented MCP catalog.

## Packet S2-01 — Complete the canonical Beans collection and detail routes

Priority: P0  
Product surface: Beans articles and headlines  
Change type: CODE, QUERY, SPEC, GENERATED, TEST, DOC  
Artifacts: Beans router/repository/tests, both Beans OpenAPI specs, Beans guide

Instructions:

1. Complete `GET /beans/articles` with the Stage 1 time, sort, cursor, filter,
   and content contracts.
2. Publish `GET /beans/articles/{id}` after S1-07 passes.
3. Complete `GET /beans/headlines` with a documented recency window and
   eligibility definition.
4. Keep search/latest/top-headlines aliases and add deprecation metadata plus a
   migration table.

Acceptance:

- Canonical and legacy-route equivalence is tested where semantics overlap.
- Article detail always has a stable UUID.
- Headlines are explainable as a product selection, not merely an unexplained
  query constant.

## Packet S2-02 — Publish Beans relation, propagation, and source routes

Priority: P1  
Product surface: Beans related articles, propagation, sources  
Change type: CODE, QUERY, SPEC, GENERATED, TEST, DOC

Instructions:

1. Publish `/articles/{id}/related`.
2. Publish single-article and batch propagation routes.
3. Resolve IDs to URLs internally.
4. Separate publisher coverage from forum/social mentions in the schema.
5. Publish source list/detail using S1-05.
6. Retain URL batch inputs for migration and external URLs where explicitly
   supported.

Acceptance:

- Clients do not need to know that database relations are URL-keyed.
- Mention metrics preserve forum-specific meaning.
- Batch limits, partial failures, and per-input results are explicit.

## Packet S2-03 — Publish Beans taxonomy, story, and analytics routes

Priority: P1  
Product surface: Beans discovery and intelligence  
Change type: CODE, QUERY, SPEC, GENERATED, TEST, DOC

Instructions:

1. Publish canonical taxonomy routes using S1-04.
2. Publish story routes only after S1-06 acceptance passes.
3. Publish analytical routes with freshness and maximum-range metadata.
4. Add story and analytics MCP tools only after the REST contracts stabilize.

Acceptance:

- Story IDs and pagination are stable.
- Taxonomy values are canonical.
- Analytical queries have measured latency limits and bounded ranges.

## Packet S2-04 — Complete Espresso action, event, and signal resources

Priority: P0  
Product surface: Espresso core resources  
Change type: CODE, QUERY, SPEC, GENERATED, TEST, DOC

Instructions:

1. Complete collection and detail routes for events and signals.
2. Publish actions and action detail only after the S1-07 action dependency
   passes.
3. Apply the same time, sort, cursor, error, and response-envelope conventions
   across all three kinds.
4. Return typed stable digest fields and an extension area for less stable
   fields; do not flatten arbitrary JSON without a contract.

Acceptance:

- Actions, events, and signals share predictable navigation.
- The API does not claim actions before readable data exists.
- JSON and compact text preserve IDs, kinds, and exact timestamps.

## Packet S2-05 — Publish typed Espresso relationship routes

Priority: P0  
Product surface: Espresso graph navigation  
Change type: CODE, QUERY, SPEC, GENERATED, TEST, DOC

Instructions:

1. Package S1-01 as the typed relationship subresources in the target table.
2. Never accept relation direction from a client.
3. Add link objects or documented next operations between detail responses and
   their relation routes.
4. Deprecate `/related/{relationship}` but keep it operational during the
   migration window.

Acceptance:

- `SAME_AS` results are orientation-independent.
- Signal evidence and reverse event-to-signal navigation return the expected
  opposing sets.
- OpenAPI descriptions never teach users the storage direction.

## Packet S2-06 — Publish Espresso sources, taxonomies, and analytics

Priority: P1  
Product surface: Espresso discovery and analytics  
Change type: CODE, QUERY, SPEC, GENERATED, TEST, DOC

Instructions:

1. Publish source list/detail from S1-05.
2. Publish taxonomy/tag discovery from S1-04.
3. Publish event and signal analytics from S1-06.
4. State metadata completeness and analytical freshness in every relevant
   schema.

Acceptance:

- Source and taxonomy IDs/values can be copied directly into collection
  filters.
- Analytics are bounded, paginated where appropriate, and cacheable.

## Packet S2-07 — Curate the MCP tool surface

Priority: P0  
Product surface: Beans MCP + Espresso MCP  
Change type: SPEC, DOC, TEST

Instructions:

1. Expose the recommended task-oriented tools, not every compatibility route.
2. Write descriptions using the twelve required documentation questions.
3. Add read-only/idempotent hints.
4. Include positive and negative tool-selection examples.
5. Test representative natural-language questions against expected tool
   selection and argument shape.
6. Measure the token size of JSON and text results.

Acceptance:

- The common user intents require no generic relation tool.
- Deprecated routes are absent from the primary MCP catalog.
- Evaluation fixtures show correct product and tool selection.

# Stage 3 — Improve, monitor, and operate

Goal: make usability, correctness, freshness, and parity measurable.

## Packet S3-01 — Add request and query observability

Priority: P0  
Product surface: Beans + Espresso runtime  
Change type: CODE, OPS, TEST  
Artifacts: middleware, repository instrumentation, dashboard definition

Record:

- request ID and `operationId`;
- route, status, latency, result count, and response bytes;
- database duration and embedding duration;
- query mode, sort, page size, and whether content/relations were included;
- rate-limit outcome and cache hit where applicable.

Do not record API keys, raw authorization, full article content, or raw search
queries by default.

Acceptance:

- A single request can be traced across gateway, API, embedder, and database.
- High-cardinality or sensitive values are excluded.
- P50/P95/P99 latency and error rate are visible by canonical operation.

## Packet S3-02 — Add data freshness and integrity metrics

Priority: P0  
Product surface: Beans + Espresso data health  
Change type: QUERY, OPS, DOC

Measure:

- latest publication, collection, event, signal, and action timestamps;
- year-0001/sentinel timestamp counts;
- UUID completeness for Beans articles;
- embedding coverage by kind;
- missing source, URL, digest, and key digest-field rates;
- orphan relation rate by relationship;
- taxonomy alias/duplicate rate;
- source metadata completeness;
- last trend-aggregate refresh;
- counts by article/sip kind.

Acceptance:

- Every user-visible freshness promise has a corresponding metric.
- Orphan and sentinel thresholds alert before they materially affect results.
- Metrics are aggregate and do not expose article content.

## Packet S3-03 — Define SLOs and alerts

Priority: P1  
Product surface: Gateway + APIs + dependencies  
Change type: OPS, DOC

Initial SLO candidates:

- availability by canonical operation;
- P95 latency for scalar search, semantic search, detail, relation, and
  analytics route classes;
- ingestion-to-API freshness by product kind;
- successful embedding rate;
- valid relation resolution rate;
- MCP tool-call success and response-size targets.

Instructions:

1. Establish baselines before choosing final thresholds.
2. Separate API availability from upstream freshness.
3. Write an owner, severity, and runbook link for each alert.
4. Use burn-rate alerts for availability SLOs.

Acceptance:

- Every alert has a user-impact statement and actionable runbook.
- A stale pipeline does not masquerade as healthy merely because HTTP returns
  200.

## Packet S3-04 — Add compatibility, contract, and parity checks

Priority: P1  
Product surface: All public routes and MCP tools  
Change type: TEST, OPS, DOC

Instructions:

1. Diff backend Swagger against gateway OpenAPI after normalizing product
   prefixes.
2. Test deprecated alias equivalence.
3. Add consumer-style contract tests for documented examples.
4. Run MCP tool-selection and response-token evaluations.
5. Quarterly, compare the route/parity matrix against NewsAPI, Perigon,
   GDELT, Event Registry, and relevant structured-intelligence providers.
6. Treat competitor parity as user-capability parity, not route-name copying.

Acceptance:

- CI fails on undocumented public route drift.
- Every quickstart example is executable.
- Parity reviews produce explicit keep/add/reject decisions.

## Packet S3-05 — Use telemetry to simplify the surface

Priority: P2  
Product surface: Developer experience  
Change type: DOC, SPEC, CODE, OPS

Instructions:

1. Review route/tool usage, validation failures, empty-result rates, and
   follow-up sequences.
2. Improve descriptions or defaults when a recurring error indicates contract
   confusion.
3. Remove a compatibility alias only after its documented migration period and
   observed usage threshold.
4. Promote or demote MCP tools based on successful user workflows, not raw call
   count alone.

Acceptance:

- Every removal has usage evidence, a migration notice, and a rollback plan.
- Documentation changes are tied to observed friction.

## Recommended execution order

Run packets in this order unless a dependency blocks them:

1. S0-01, S0-02, S0-04, S0-05
2. S0-03
3. S1-01, S1-03, S1-04, S1-05
4. S1-02 and S1-07 as external dependencies become ready
5. S1-06
6. S2-01, S2-04, S2-05
7. S2-02, S2-06, S2-07
8. S2-03
9. S3-01 and S3-02 before setting final SLO thresholds
10. S3-03, S3-04, S3-05

S0-01, S0-02, S0-04, and S0-05 can proceed immediately because they do not
need schema or major query changes. S1-01 is the first correctness-sensitive
query change. S2-05 must not precede it.

## Artifact map

| Change | Beans artifacts | Espresso artifacts | Shared/public artifacts |
|---|---|---|---|
| Router or validation code | `apis/beans/router/routes.go` | `apis/espresso/router/routes.go`, `apis/espresso/router/types.go` | — |
| Query/repository code | `apis/beans/db/pgsack.go`, `apis/beans/db/types.go` | `apis/espresso/db/pgcupboard.go`, `apis/espresso/db/types.go` | — |
| Tests | `apis/beans/tests/` | `apis/espresso/tests/` | root gateway tests where applicable |
| Backend generated API spec | `apis/beans/docs/` | `apis/espresso/docs/` | — |
| Public gateway API spec | — | — | `config/beans.oas.json`, `config/espresso.oas.json` |
| Product guide | `docs/pages/howtos/beans-howto.mdx` | `docs/pages/howtos/espresso-howto.mdx` | `docs/pages/api-overview.mdx`, `docs/pages/howtos/mcp-howto.mdx` |
| Portal navigation | — | — | `docs/zudoku.config.tsx` only when a new guide/page must be mounted |

## Explicit non-goals

- Designing or executing the Beans UUID backfill/primary-key conversion.
- Designing or executing Espresso action ingestion or storage.
- Replacing the gateway, authentication, quota, or subscription model.
- Copying competitors' proprietary schemas.
- Exposing raw relation edge orientation to users.
- Performing a wholesale taxonomy rewrite during route packaging.
- Adding write APIs.
- Removing current routes before compatibility evidence exists.

## Definition of roadmap completion

This roadmap is complete when:

- the canonical route tables are implemented or explicitly dependency-gated;
- current routes have safe compatibility mappings;
- Beans supports intuitive article, headline, related, propagation, source,
  taxonomy, story, and analytical workflows at the stages described;
- Espresso exposes actions, events, and signals as first-class resources, with
  relation direction entirely hidden by typed navigation;
- REST and MCP documentation let an agent select and sequence tools without
  reading implementation code;
- backend Swagger, gateway OpenAPI, how-to documentation, and behavior remain
  synchronized; and
- usability, freshness, latency, errors, relation integrity, and MCP response
  size are observable.

## External references used for parity

- [NewsAPI endpoints](https://newsapi.org/docs/endpoints)
- [Perigon API introduction](https://perigon.io/docs/api/intro)
- [Perigon data model](https://perigon.io/docs/api/data-model)
- [Perigon sources](https://perigon.io/docs/api/sources)
- [Perigon MCP server](https://perigon.io/docs/api/mcp)
- [GDELT Cloud API v2](https://docs.gdeltcloud.com/api-reference/v2)
- [Event Registry / NewsAPI.ai article search](https://newsapi.ai/documentation?tab=searchArticles)
- [Alpha Vantage API documentation](https://www.alphavantage.co/documentation/)
- [Trading Economics API documentation](https://docs.tradingeconomics.com/)
- [Model Context Protocol tool annotations](https://blog.modelcontextprotocol.io/posts/2026-03-16-tool-annotations/)
