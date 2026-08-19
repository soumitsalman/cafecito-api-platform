# Beans & Espresso API Capability and Documentation Gap Analysis

**Analysis date:** 2026-07-30  
**Status:** Ideation only; no implementation, OpenAPI, or published documentation changes are included.

## 1. Executive summary

Cafecito should present Beans and Espresso as two layers of one research system rather than as two overlapping search APIs:

- **Beans is the evidence and media-observation layer:** what was published, by whom, when, how sources covered the same story, and what is gaining attention.
- **Espresso is the interpretation layer:** what happened, what it means, what may happen next, and which lower-level observations support that conclusion.

The strongest product opportunity is an explicit evidence chain:

```text
Beans articles and propagation
    -> Espresso actions
        -> Espresso events
            -> Espresso signals
                -> supporting actions, events, and source URLs
```

Beans already has more differentiated backend assets than its public surface suggests: full or summarized article content, semantic search, publisher metadata, social observations, related-article edges, and a transparent trend aggregation. Espresso has a flexible intelligence graph (`sips`, `sources`, and `relations`) and an unusually useful agent-oriented text format. The primary gaps are not “more AI.” They are exposing the graph coherently, preserving provenance, making ranking and time semantics explicit, and turning latent schema capabilities into stable API contracts.

The highest-priority work is:

1. Correct contract contradictions before adding endpoints.
2. Expose Espresso actions and evidence/provenance.
3. Keep stable record IDs in Espresso's text format so agents can traverse relationships.
4. Make Beans story clusters and publisher discovery first-class.
5. Add explicit time ranges, sorting, continuation metadata, and stable response envelopes to both products.
6. Publish one short “Beans vs. Espresso” workflow guide centered on evidence -> interpretation.

## 2. Scope and method

This analysis reviewed:

- Beans and Espresso Go routes, response types, database query layers, tests, service READMEs, and generated Swagger under `apis/`.
- Public guides, MCP documentation, navigation, and product positioning under `docs/`.
- Public gateway OpenAPI documents only where necessary to compare the public contract with backend behavior.
- The attached ingestion/storage sources:
  - `/home/soumitsr/codes/pycoffeemaker/pybeansack/pgsack.sql`
  - `/home/soumitsr/codes/pycoffeemaker/pycupboard/pgcupboard.py`
  - Their imported `models.py` definitions where needed to interpret the stored fields.
- Current official documentation for comparable news, event, and market-data services.

This document evaluates product capability and documentation. It does not prescribe internal implementation details unless the current schema materially constrains a recommendation.

## 3. Recommended product story

### Beans

**One-line promise:** “Find the source material: recent and trending news and blogs, semantically searchable, with coverage and propagation context.”

Beans should answer:

- What has been published about this topic?
- What is new?
- What is gaining attention?
- Which publishers covered the same story?
- Where is the story being discussed?
- Can I retrieve the full article text for selected records?

Beans should not try to become the primary opinionated intelligence product. Its value is broad source coverage, retrieval, enrichment, and transparent media-observation signals.

### Espresso

**One-line promise:** “Understand the development: structured actions, events, and signals with a traceable evidence chain.”

Espresso should answer:

- What atomic action or observation occurred?
- Which actions form a coherent event?
- Which events support a larger signal?
- What are the impact, drivers, affected domains, and outlook?
- How important, novel, or fast-moving is it?
- What source evidence supports the conclusion?

Espresso should not become a second general-purpose news archive. It should expose curated intelligence with better structure, provenance, and synthesis than a raw article API.

### Cross-product workflow

The public narrative should show two canonical workflows:

1. **Research a live topic:** Beans search -> inspect summaries -> fetch selected full content -> inspect same-story coverage/propagation.
2. **Understand a development:** Espresso signals -> supporting events -> underlying actions -> source URLs or linked Beans articles.

The second workflow is not currently complete because actions and source lineage are not publicly retrievable.

## 4. Current-state capability map

### 4.1 Beans public surface

Current backend routes:

| Capability | Endpoint(s) | Current assessment |
|---|---|---|
| Health | `GET /health` | Complete liveness probe |
| Tag discovery | `GET /tags/categories`, `/tags/entities`, `/tags/regions` | Useful exact-filter vocabulary, but no sentiment or source discovery |
| Publisher resolution | `GET /sources` | Backend can list or filter, but public contract presents it primarily as known-ID resolution |
| Full-corpus search | `GET /articles/search` | Semantic search, URL lookup, and scalar filtering |
| Recent feed | `GET /articles/latest` | Reverse chronological, default seven-day window |
| Trending feed | `GET /articles/trending` | Trend-ranked, default seven-day activity window |
| Daily headlines | `GET /articles/top-headlines` | Fixed 24-hour trend-ranked window |
| Story propagation | `GET`/`POST /articles/propagation` | Cross-publisher related coverage plus social/forum mentions |
| Agent access | Hosted MCP tools mirroring operations | Broad tool coverage, but no compact text response |

Important backend assets from the attached schema:

- `beans.url` is the canonical primary key; the UUID `id` is not the public identity and is not constrained as unique in the attached SQL.
- `beans` stores summary, optional content, restricted-content state, publish and collection timestamps, categories, sentiments, regions, entities, and a 320-dimensional embedding.
- `publishers` stores display metadata and RSS feed URLs.
- `chatters` stores repeated social observations and engagement counters.
- `related_beans` stores directed article-to-article edges.
- `trend_aggregates` combines related coverage, comments, observed shares, likes, and recency.
- `aggregated_beans_view` already includes related URLs, a cluster representative, trend fields, and publisher metadata, though not all of these are exposed by the API.

### 4.2 Espresso public surface

Current backend routes:

| Capability | Endpoint(s) | Current assessment |
|---|---|---|
| Health | `GET /health` | Complete liveness probe |
| Tag discovery | `GET /tags` | Simple global vocabulary |
| Event retrieval | `GET /events` | ID/tag/date filtering and semantic search |
| Signal retrieval | `GET /signals` | Same retrieval model as events |
| Relationship traversal | `GET /related/{relationship}` | Supports `same_as` and `derived_from`, but direction is not represented in the response contract |
| Agent format | `response_type=text` | Valuable differentiator, but currently loses IDs and some precision |
| Agent access | Hosted MCP tools mirroring operations | Good starting point, incomplete hierarchy |

Important backend assets from the attached schema:

- `sips` stores `action`, `event`, and `signal` records with UUID, kind, source UUID, URL, base URL, timestamp, tags, JSON digest, and a required 320-dimensional embedding.
- Sip IDs are deterministically derived from URLs during ingestion.
- `sources` stores source identity and display metadata.
- `relations` stores typed directed edges but does not enforce foreign keys.
- Ingestion is immutable on ID conflict and has a cleanup window, so retention is an operational property that should be made explicit.
- The public router returns only flattened digest data plus `id` and `reported`; stored `kind`, `source`, `url`, and `base_url` are not exposed.
- Scalar list queries sort by newest first and then by relation count. Semantic queries sort by vector distance.

### 4.3 Schema and operations drift

Both attached schemas declare `vector(320)`. The repository knowledge note describes 384-dimensional embeddings. The service READMEs agree with 320. This should be resolved as an operational source-of-truth issue before any embedding migration or client-facing documentation is changed.

Beans and Espresso also depend on ingestion-owned retention, refresh, clustering, and enrichment behavior. Public capability work therefore needs an explicit contract between the ingestion repository and this API repository for:

- Embedding model and dimension.
- Retention and archive windows.
- Refresh cadence and observable freshness.
- Stable taxonomy and digest schema versions.
- Relation direction and meaning.
- Full-content availability and restriction semantics.

## 5. Competitive benchmark

This is a capability benchmark, not a recommendation to clone every competitor.

### 5.1 Comparable services

| Service | Relevant current capabilities | Product lesson |
|---|---|---|
| [NewsAPI](https://newsapi.org/docs/endpoints) | Full-corpus “Everything,” top headlines, and source listing | Baseline news APIs make source discovery, country/category filtering, and simple feed choice obvious |
| [Perigon](https://perigon.io/docs/api/intro) | Articles, semantic/vector search, clustered stories, sources, people, companies, topics, journalists, and summarization | Treat articles, stories, and entities as separate resources with stable identities |
| [Perigon data model](https://perigon.io/docs/api/data-model) | Story clusters evolve, aggregate entities, track size/scope, and link articles, sources, and journalists | Beans' latent related-article graph should become a first-class story resource |
| [Newscatcher clustering](https://www.newscatcherapi.com/docs/news-api/guides-and-concepts/clustering-news-articles) | Query-time semantic clustering and explicit clustering controls | Users need to choose between raw articles, deduplicated results, and story clusters |
| [Event Registry / NewsAPI.ai](https://newsapi.ai/documentation?tab=searchArticles) | Boolean and exclusion queries, events, concepts, autosuggest, streams, breaking events, social ranking, and aggregate/facet results | Advanced research APIs expose sort, include/exclude logic, facets, and streams rather than only row retrieval |
| [GDELT Cloud API v2](https://docs.gdeltcloud.com/api-reference/v2) | Events, clustered stories, entities, evidence articles, normalized geography, significance/confidence, summary endpoints, stable detail routes, and cursor pagination | Espresso needs explicit evidence, confidence/significance, summaries, typed filters, and detail resources |
| [Alpha Vantage News & Sentiment](https://www.alphavantage.co/documentation/) | Ticker/topic filters, start and end times, explicit sorting, relevance, and per-entity sentiment | Market intelligence is more useful when entity identifiers, time bounds, relevance, and sentiment are queryable |
| [Trading Economics](https://docs.tradingeconomics.com/) | Economic indicators, calendars, forecasts, market data, and structured country/indicator resources | Espresso should link intelligence to canonical market/macro identifiers rather than trying to reproduce raw market-data feeds |
| [Trading Economics streaming](https://docs.tradingeconomics.com/economic_calendar/streaming/) | Live calendar releases over a streaming connection | Alerts/streams are a later-stage parity item for monitoring workflows |

### 5.2 Feature-parity summary

Legend: **Strong** = clearly competitive; **Partial** = backend capability or limited public form; **Gap** = expected capability not currently public.

| Capability | Beans | Espresso | Competitive implication |
|---|---|---|---|
| Recent retrieval | Strong | Strong | Baseline parity |
| Semantic search | Strong | Strong | Competitive |
| Scalar tags | Strong | Partial | Espresso lacks typed facets |
| Source discovery | Partial | Gap | Both schemas have sources; neither presents a strong discovery workflow |
| Full text | Strong, conditional | Not applicable | Beans differentiator if rights/availability are documented |
| Story clustering | Partial | Partial | Relations exist, but there is no stable story resource |
| Actions -> events -> signals | Not applicable | Partial | Core Espresso promise is incomplete without actions |
| Evidence/provenance | Partial via propagation | Gap in public Espresso output | Highest-value differentiation opportunity |
| Trend/significance ranking | Strong but opaque | Hidden/weak | Espresso needs explicit importance versus recency |
| Aggregations/facets | Gap | Gap | Major analyst/dashboard gap |
| Time start + end | Start only | Start only | Below common parity |
| Explicit sorting | Endpoint-implied only | Query-dependent only | Below common parity |
| Typed company/ticker/place filters | Flat entities | Flat tags | Below finance/intelligence parity |
| Stable detail endpoints | URL batch lookup only | ID batch lookup only | Awkward for chaining and citations |
| Cursor/continuation metadata | Gap | Gap | Offset pagination is fragile for live data |
| Streaming/alerts/webhooks | Gap | Gap | Later-stage monitoring parity |
| Token-efficient output | Gap | Strong but lossy | Espresso differentiator; Beans can use field selection before adding text |
| MCP | Strong | Strong | Competitive distribution advantage |

## 6. Beans capability gaps

### 6.1 First-class story clusters

`related_beans`, cluster IDs, related URL arrays, and publisher joins already exist in the storage/view layer. The public API exposes only a propagation lookup from seed URLs. This makes the graph useful only after the caller already has an article URL and does not provide a canonical “story” object.

Recommended resource:

```text
GET /stories
GET /stories/{story_id}
GET /stories/{story_id}/articles
```

A story should expose:

- Stable story/cluster ID.
- Generated or representative title.
- First and latest publish times.
- Article count and unique source count.
- Representative and newest articles.
- Categories, regions, entities, and sentiments aggregated from articles.
- Coverage velocity and trend score with an `as_of` timestamp.
- Optional propagation/social summary.

Do not expose the current representative URL as a permanent cluster identity until stability through merges/splits is defined.

### 6.2 Publisher discovery

The database layer already has `DistinctSources`; the current `/sources` handler can query without source IDs, but its public contract describes source resolution and marks `sources` as required.

Choose and document one behavior:

- Preferably make `GET /sources` a discoverable catalog with optional `ids`, `q`, and metadata filters.
- Or keep it as a strict resolver and add `GET /sources/search`.

Useful future source fields require ingestion support: country, language, source type, topical focus, paywall/full-content availability, update cadence, and quality/coverage metadata.

### 6.3 Query completeness

Parity-oriented additions:

- `to` timestamp/date.
- Explicit `sort=relevance|published_desc|trend`.
- `language`, publisher country, and article geography as distinct concepts.
- Include/exclude source and taxonomy filters.
- Sentiment discovery and filtering, since sentiment is already returned.
- Optional facets such as counts by source, category, region, entity, or sentiment.
- Stable cursor pagination for live feeds.

Avoid adding dozens of query knobs at once. Start with `to`, `sort`, source discovery, and one aggregation endpoint.

### 6.4 Trend transparency

The attached SQL computes:

```text
(100 * related + 50 * comments + 10 * shares + likes)
-------------------------------------------------------
                 age / recency term
```

The exact SQL uses `CURRENT_DATE + 2 - updated` as the denominator. It also calculates `shares` as a count of chatter records, even though the `chatters` table has a separate `shares` column. Public descriptions currently imply repost/share totals. Before promoting trend analytics:

- Define whether “shares” means observed posts, platform share counters, or both.
- Publish `trend_score_version`, `trend_updated_at`, and the observation window.
- Explain that metrics are lower bounds from collected platforms.
- Decide whether subscriber/audience size belongs in the score; it is returned but not used by the attached formula.
- State how and when the materialized view refreshes.

### 6.5 Full-content contract

`full_content=true` is a meaningful differentiator, but callers cannot tell whether content is absent, restricted, not collected, or intentionally omitted. Add a compact availability contract:

```json
{
  "content_status": "available|restricted|unavailable|not_requested",
  "content": "..."
}
```

Document licensing/citation expectations, freshness, extraction format, and whether Markdown conversion preserves links and headings. Prefer `include=content` or field selection over returning large bodies by default.

### 6.6 Trend and facet summaries

Article ranking does not answer “what topics are rising?” Add a summary surface only after score semantics are stable:

```text
GET /trends?group_by=entity|category|region|source&from=...&to=...
```

Return current volume, previous-period volume, delta, representative articles, source diversity, and `as_of`. This is more useful to dashboards and agents than asking them to download 128 articles and count tags.

## 7. Espresso capability gaps

### 7.1 Actions are missing from the public hierarchy

This is the most direct capability gap. Actions are described as the foundation of events and signals and are stored as sips, but the Go database constants and router expose only event and signal kinds.

Recommended:

```text
GET /actions
GET /sips/{id}
```

Actions should support the same basic retrieval contract as events/signals: IDs, semantic query, typed filters, time range, sort, pagination, and text/JSON output.

The action schema should be explicit enough to represent:

- Subject/entity and canonical identifiers.
- Action type.
- Observed value and unit where applicable.
- Event/market timestamp versus ingestion/report timestamp.
- Source URL and source metadata.
- Confidence and extraction method.

### 7.2 Evidence and provenance

Espresso stores `source`, `url`, and `base_url`, but the public response drops them. A signal can make a strong claim without returning the evidence needed to verify it.

Every action/event/signal detail response should include:

- `id`, `kind`, `reported_at`, and when relevant `observed_at`.
- Source URL(s), source identity, and publisher/site metadata.
- Supporting sip references grouped by kind.
- Relation type and direction.
- Evidence count and source diversity.
- Confidence, coverage, and generation/extraction version where meaningful.
- A short citation-ready evidence preview, with a paginated route for full evidence.

Recommended shape:

```text
GET /sips/{id}
GET /sips/{id}/evidence
GET /sips/{id}/relations?relationship=derived_from&direction=upstream
```

Cross-link source URLs to Beans records when available. This is the most defensible product differentiator in the roadmap.

### 7.3 Relationship semantics and direction

The attached schema stores directed `from_id -> to_id` edges. The current related-sip SQL matches either side of an edge, so the API behaves as undirected traversal. Documentation describes `derived_from` as downstream traversal. Those are not equivalent.

Before adding more relationship types:

- Define the direction of `DERIVED_FROM` in one sentence and one diagram.
- Expose `from_id`, `to_id`, relationship, and direction or provide `direction=upstream|downstream|both`.
- Return edge metadata rather than only flattened neighboring sips.
- Define behavior for missing endpoints because relations do not have foreign keys.
- Consider `supported_by`, `contradicts`, `updates`, and `supersedes` only after the two existing relations are unambiguous.

### 7.4 Importance, novelty, and trending

Scalar Espresso queries currently sort primarily by recency and secondarily by relation count; semantic queries sort by similarity. The API does not expose why one record is more important or more observed than another.

Add explicit:

- `sort=recent|relevance|significance|trending`.
- A versioned `significance` or `trend_score`.
- `novelty`, evidence volume, source diversity, confidence, and `as_of`.
- A documented observation window.

Do not label relation count alone as “trending.” A defensible signal should combine evidence volume, source diversity, recency, novelty, and confidence.

### 7.5 Typed intelligence filters

One flat tag list is insufficient for market and political intelligence. Promote stable dimensions:

- Companies and canonical company IDs.
- Stock tickers and exchange-qualified instruments.
- Products/commodities.
- People and organizations.
- Countries, regions, and locations.
- Domains/sectors.
- Event type and action type.
- Impact level.
- Forecast horizon.

The JSON digests already contain some of these keys, but filtering them through a single tags array loses type and creates collisions. Add entity discovery/resolution before adding dozens of string filters.

### 7.6 Summary and change-detection endpoints

Espresso should answer “what changed?” without returning raw lists:

```text
GET /signals/summary?group_by=domain|company|region
GET /events/summary?group_by=event_type|region
GET /emerging?dimension=entity|topic
```

Useful outputs include counts, deltas versus a prior window, significance distribution, representative records, and source diversity.

### 7.7 Text format must remain chainable and deterministic

`response_type=text` is a real advantage, but the current implementation has several contract problems:

- It omits `id`, so an agent cannot take a text result and call `/related/{relationship}`.
- It renders timestamps as dates, losing time precision relative to JSON.
- It iterates ordinary map keys without deterministic ordering.
- It labels extracted entity values as `tags:` and may also render the digest's own `tags`, producing two semantically different `tags:` lines.
- Tests and READMEs expect `related:` or `related_to:`, while current code emits different labels.
- It omits `site_name`.

Publish and test a stable grammar. At minimum every record needs:

```text
id:
kind:
reported:
source:
tags:
briefing:
...
```

Use deterministic field ordering and escaping. If token cost is the goal, a compact JSON projection or `fields=` parameter may be safer than a loosely specified bespoke format.

## 8. Cross-cutting API gaps

### 8.1 Pagination and empty results

Offset pagination is easy but unstable for feeds that change while a client pages through them. Neither API reports a total, next offset, cursor, or `has_more`.

Recommended future envelope:

```json
{
  "data": [],
  "page": {
    "limit": 16,
    "next_cursor": null,
    "has_more": false
  },
  "meta": {
    "as_of": "2026-07-30T12:00:00Z"
  }
}
```

For a future version, prefer `200` with an empty `data` array over `204`; it is easier for SDKs and agents to parse consistently. Preserve existing behavior until a versioned migration plan exists.

### 8.2 Stable detail resources

Batch filters (`urls=` or `ids=`) can retrieve known records but make linking, caching, citations, and OpenAPI examples awkward.

- Espresso can safely use UUID detail routes now.
- Beans should not expose its nullable/unconstrained ingestion UUID as a stable public ID without first defining it. A query-by-canonical-URL detail route is safer in the short term; stable story IDs need their own lifecycle contract.

### 8.3 Field selection and response profiles

Agent and dashboard clients need different payloads. Add one consistent mechanism:

- `fields=` for projection, or
- `view=summary|detail|evidence`, or
- `include=content,publisher,evidence`.

This is preferable to making every endpoint return all enrichment and reduces pressure to duplicate Espresso's text renderer in Beans.

### 8.4 Time and freshness semantics

Both products need:

- `from` and `to`.
- Explicit timezone behavior.
- Clear separation of published/observed/reported/collected/updated timestamps.
- `as_of` on aggregates and ranks.
- Documented default time windows.
- Retention and historical coverage by plan.

### 8.5 Reliability and operational contract

Publish:

- Rate-limit and quota headers.
- `Retry-After` behavior for 429.
- Request/correlation IDs.
- Timeout recommendations.
- A status page and incident link.
- Data freshness targets.
- Changelog and deprecation policy.
- Versioning policy for fields, taxonomies, scores, and digest schemas.

Do not prioritize SDKs until the response envelopes and field contracts are stable. The existing cURL, JavaScript, Python, OpenAPI, MCP, `llms.txt`, and interactive reference are sufficient for the current stage.

## 9. Documentation gap analysis

### 9.1 What is already strong

- Public how-tos provide cURL, JavaScript, and Python examples.
- The Beans guide explains exact versus fuzzy tags and progressive filter discovery.
- The Espresso guide introduces the sip hierarchy and text mode.
- Both references have operation IDs that map naturally to MCP tools.
- Empty-result behavior and pagination limits are documented.
- The portal publishes OpenAPI reference pages plus `llms.txt`/`llms-full.txt`, which is valuable for agents.
- Beans documents endpoint-specific response shapes more clearly than many early-stage APIs.

### 9.2 Immediate truth/contract corrections

| Issue | Current conflict | Recommended decision |
|---|---|---|
| Beans search `trend_score` | Go godoc/generated Swagger say search returns it; query field selection omits it; public how-to correctly says it is absent | Either select and return it consistently or remove it from schema/description |
| Beans `acc` default | Service README says `0.75`; router and public guide use `0.5` | Make `0.5` the single documented default unless behavior changes |
| Beans query length | Public docs say 3–512 characters; binding enforces only max 512 | Enforce the minimum or remove the claim |
| Beans sources | Public contract says `sources` is required; router does not require it and can return a catalog page | Decide resolver versus discovery semantics |
| Beans trend `shares` | Type/docs imply share totals; attached SQL uses chatter-record count in the score/view | Rename or redefine the metric |
| Espresso `site_name` | Public how-to, response type comment, and examples say it is merged; `enrichSipDigest` adds only `id` and `reported` | Join/expose source metadata or remove the claim |
| Espresso text grammar | How-to shows `related:`; service README shows `related_to:`; current renderer uses `tags:`; tests expect another variant | Define one stable grammar and update implementation/tests/docs together |
| Espresso text chainability | Text output omits IDs although related traversal requires IDs | Always include `id` and `kind` |
| Espresso related direction | Docs describe downstream traversal; query traverses both sides | Define and expose direction |
| Espresso tag catalog | Docs say event/signal tags; SQL selects tags from all sip kinds, including actions | Scope the query or state that action tags may appear |
| Espresso response timestamp | Public API uses `reported`; parts of service README still show `created` or `date` | Standardize timestamp names and meanings |
| Espresso index settings | Attached schema uses HNSW `m=24`, `ef_construction=128`; service README says `m=16`, `64` | Generate operational docs from the schema or remove volatile tuning details |
| Health auth | API overview says every endpoint requires a key; health is public at the gateway/backend route level | Document health as unauthenticated if that remains intentional |
| MCP metering | MCP guide says REST and MCP share rate/quota behavior; public route policies need to be verified against that statement | Make behavior match or qualify the docs |
| Concurrency 429 | Swagger describes concurrency rejection; middleware currently waits on a bounded channel | Distinguish gateway rate limiting from backend queueing |
| Stale README paths/runtime | Service READMEs reference old package paths and an older standalone TEI topology | Refresh or reduce implementation-specific README content |

### 9.3 Missing conceptual documentation

Create one concise “Choose an API” page:

| Question | Use |
|---|---|
| What was published? | Beans latest/search |
| What is gaining media/social attention? | Beans trending/top headlines |
| Who else covered or discussed this story? | Beans propagation/story |
| What concrete development occurred? | Espresso event |
| What atomic observation supports it? | Espresso action |
| What does it mean across sources/domains? | Espresso signal |
| What evidence supports this conclusion? | Espresso evidence -> Beans/source URLs |

Also document:

- Exact definitions of action, event, signal, story, related article, propagation, and trend.
- Required versus optional fields.
- Score formulas or at least score inputs, version, window, and `as_of`.
- Data freshness, historical coverage, retention, and source coverage.
- Missing/restricted content behavior.
- Query recovery recipes for empty results.
- Relationship direction with an example graph.
- Text versus JSON tradeoffs and a machine-readable text grammar.

### 9.4 Documentation information architecture

Recommended top-level structure:

```text
Start
  - Choose Beans or Espresso
  - Authentication and limits
  - Quickstarts

Beans
  - Concepts: article, story, source, propagation, trend
  - Search and filtering
  - Story/propagation workflow
  - Full-content availability
  - API reference

Espresso
  - Concepts: action -> event -> signal
  - Evidence and relations
  - Search, ranking, and filters
  - JSON and text formats
  - API reference

Agents and MCP
  - Client configuration
  - Recommended tool sequences
  - Token-efficient response profiles
  - Retry and empty-result behavior

Operations
  - Errors, limits, pagination
  - Freshness and retention
  - Changelog and deprecations
  - Status
```

## 10. Prioritized roadmap

### P0: Make the current contract true

1. Resolve every contradiction in section 9.2.
2. Define the source of truth for embedding dimension and HNSW settings.
3. Define Espresso relation direction.
4. Define and test one deterministic Espresso text grammar that includes IDs.
5. Publish the Beans/Espresso decision guide.
6. Document score windows, freshness, retention, and full-content states.

### P1: Complete each product's core promise

**Espresso**

1. Expose `GET /actions`.
2. Add `GET /sips/{id}` with kind, source, URL, and evidence summary.
3. Add directional evidence/relation traversal.
4. Expose source metadata and Beans/source links.

**Beans**

1. Make source discovery first-class.
2. Add stable story cluster list/detail/article surfaces.
3. Expose cluster and trend freshness metadata.

**Both**

1. Add `to` and explicit `sort`.
2. Add stable pagination envelopes/continuation metadata in a versioned response.
3. Add field selection or response profiles.

### P2: Improve analyst usefulness

**Espresso**

- Typed entity/ticker/domain/action/event filters.
- Significance, novelty, confidence, evidence volume, and source diversity.
- Event/signal summary and emerging-topic endpoints.

**Beans**

- Sentiment filtering/discovery.
- Language and source-country support.
- Trend/facet summaries.
- Include/exclude and source-group filters.

### P3: Monitoring and ecosystem

- Webhooks or saved-query alerts.
- Streaming for breaking events/signals where demand justifies it.
- Bulk export/asynchronous jobs for research datasets.
- Official SDKs after contracts stabilize.
- User-defined monitors, watchlists, and scheduled briefs.

## 11. Proposed capability acceptance criteria

### Product clarity

- A new user can choose Beans or Espresso in under one minute.
- Every Espresso signal can be traced to events/actions and source evidence.
- Every documented score states its purpose, inputs, window, version, and `as_of`.

### API ergonomics

- List responses have deterministic continuation metadata.
- Empty results are handled consistently.
- JSON and text responses preserve IDs needed for the next workflow step.
- Time filters can reproduce a closed historical window.
- Detail resources have stable identities and citation-ready URLs.

### Documentation quality

- Backend behavior, generated Swagger, gateway OpenAPI, service README, tests, and public how-to agree.
- Every endpoint includes a “when to use,” a minimal request, a representative response, and recovery guidance.
- Agent examples start with small limits and summary views.
- Beans examples request full content only after selecting records.
- Espresso examples begin with signals for synthesis and drill down to evidence.

### Product metrics

Track:

- Search-to-detail and signal-to-evidence follow-through rates.
- Empty-result rate by endpoint and filter.
- Percentage of signals with complete source lineage.
- Percentage of stories with two or more unique publishers.
- Median evidence source diversity.
- MCP tool-call retries and invalid follow-up IDs.
- Payload and token size by response profile.
- Freshness lag from ingestion to API availability.

## 12. Recommended first release slice

The smallest coherent release is not a broad set of new endpoints. It is:

1. Contract cleanup and the “Beans vs. Espresso” guide.
2. Espresso `actions`, sip detail, and source/evidence lineage.
3. A deterministic text format containing `id`, `kind`, and source reference.
4. Beans source discovery.
5. One stable Beans story-cluster detail surface.
6. `to`, `sort`, and continuation metadata across the affected list routes.

That slice completes the product narrative and creates a defensible agent workflow without trying to match every news, finance, event, and market-data platform.

## 13. Explicit non-goals

- Do not turn Beans into a derived market-opinion API.
- Do not turn Espresso into a duplicate raw news archive.
- Do not reproduce price feeds, filings, economic calendars, or fundamentals already available from specialist providers; link to canonical identifiers and data sources instead.
- Do not add streaming, alerts, SDKs, or dozens of filters before stable IDs, evidence, envelopes, and versioning exist.
- Do not expose ingestion UUIDs or cluster representatives as permanent public identifiers until stability is guaranteed.
- Do not claim full-content, sentiment, significance, or propagation precision without clear availability and methodology contracts.

## 14. Research sources

Official product/API documentation consulted:

- [NewsAPI endpoints](https://newsapi.org/docs/endpoints)
- [Perigon API overview](https://perigon.io/docs/api/intro)
- [Perigon data model](https://perigon.io/docs/api/data-model)
- [Perigon pagination](https://perigon.io/docs/api/pagination)
- [Newscatcher article clustering](https://www.newscatcherapi.com/docs/news-api/guides-and-concepts/clustering-news-articles)
- [Event Registry / NewsAPI.ai documentation](https://newsapi.ai/documentation?tab=searchArticles)
- [GDELT Cloud API v2](https://docs.gdeltcloud.com/api-reference/v2)
- [Alpha Vantage API documentation](https://www.alphavantage.co/documentation/)
- [Trading Economics API documentation](https://docs.tradingeconomics.com/)
- [Trading Economics calendar streaming](https://docs.tradingeconomics.com/economic_calendar/streaming/)

