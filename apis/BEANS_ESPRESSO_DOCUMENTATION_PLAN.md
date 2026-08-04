# Beans and Espresso Documentation Separation Plan

Status: proposed documentation contract and execution plan  
Prepared: 2026-07-30  
Scope: developer portal, OpenAPI descriptions, examples, and MCP guidance

## The distinction the documentation must teach

Use these two sentences everywhere the products are introduced:

> **Beans finds and verifies what publishers published.** It returns articles,
> sources, clusters, related reading, and external mentions.

> **Espresso explains what happened and what it may mean.** It returns structured
> actions, Event-family records, and synthesized Signals with context, impact,
> forecasts, and transparent evidence relationships.

A reader should understand this before seeing authentication, parameter lists,
or endpoint tables.

## The current documentation problem

Both products currently look like search APIs over overlapping news-derived
data, and the route proposals now provide the canonical resource vocabulary. Terms such as related, event, story, source, signal, trending, and semantic
search appear without a single boundary model. This creates predictable
confusion:

- An article about an earnings release can look interchangeable with an
  Espresso event about the same release.
- A Beans article cluster can be mistaken for an Espresso event.
- A Beans trend score can be mistaken for an Espresso signal.
- “Related” can mean semantic similarity, same underlying event, publisher
  propagation, or an intelligence dependency.
- MCP clients see operations before learning which product owns the user
  question.

The documentation must define the returned object, not merely describe the
search technology.

## Product boundary

| Dimension | Beans | Espresso |
|---|---|---|
| Primary object | Publisher article | Action, Event-family record, or Signal |
| One record means | One document published by one source | One interpreted activity, development, or synthesized conclusion |
| Primary user question | “What did publishers publish?” | “What happened, why, and what may follow?” |
| Evidence | Article URL, source metadata, content/summary, cluster coverage, external mentions | Source metadata plus Event/Signal/Action relationships |
| Search result | Citable media records | Structured intelligence records |
| Grouping unit | Cluster of articles about the same development | Event-family records connected by `SAME_AS` or `DERIVED_FROM` |
| Higher-order output | Trending article ranking and source coverage | Signal synthesized across Events, sources, and time |
| Best for | News retrieval, citations, source monitoring, publisher coverage | Market/business/political interpretation and impact analysis |
| Not for | Claiming a synthesized conclusion not present in source material | Returning the complete article corpus or article full text |



## Terminology contract

| Term | Product | Required definition | Must not mean |
|---|---|---|---|
| Article | Beans | One source document with UUID, URL, title, publication time, and content metadata. | An Event, Cluster, or Signal. |
| Source | Beans and Espresso | Beans: publisher/outlet metadata. Espresso: provenance for an intelligence record. | Automatically the same entity across databases. |
| Cluster | Beans | Durable group of articles covering the same underlying development. | An Espresso Event. |
| Related article | Beans | Ranked semantic reading recommendation for an article. | Guaranteed cluster membership or social propagation. |
| Mention | Beans | External social/forum post linking to an anchor article URL. | Another publisher article. |
| Cluster source coverage | Beans | Distinct sources represented by articles in a Cluster. | Article-level mentions. |
| Trending article | Beans | Article ranked by available attention/engagement and recency inputs. | An Espresso Signal or forecast. |
| Action | Espresso | One atomic activity, observation, or measurable change. | A publisher article. |
| Event-family record | Espresso | A record with `kind` beginning `event`, preserving its representation (curated or source-oriented). | A Beans Cluster. |
| Signal | Espresso | A conclusion synthesized across Events, sources, and time. | A social trend score or alert notification. |
| Evidence | Espresso | Other Event-family records representing the same development through bidirectional `SAME_AS` relationships. | A directional dependency. |
| Supporting Event | Espresso | Event-family record connected to a Signal through `DERIVED_FROM`. | An undirected equivalent. |



## The decision guide that must appear near the top of the portal

| User asks... | Start with | Why this route exists | Typical follow-up |
|---|---|---|---|
| “Find articles about a company.” | `GET /beans/articles` | Returns source documents that can be cited, filtered, and opened. | `GET /beans/articles/{article_id}` or `/related`. |
| “What are today's high-attention headlines?” | `GET /beans/headlines` | Produces a recent briefing of high-attention articles; it is not the full latest feed. | Open an article or inspect its trend metadata. |
| “Which outlets covered this development?” | `GET /beans/clusters/{cluster_id}/sources` | Counts distinct editorial sources represented in a Cluster. | Retrieve `/clusters/{cluster_id}/articles`. |
| “Show other reading on this article.” | `GET /beans/articles/{article_id}/related` | Returns ranked semantic neighbors for reading discovery. | Compare with the article’s Cluster when known. |
| “Which social/forum posts linked to this article?” | `GET /beans/articles/{article_id}/mentions` | Returns external posts linking the exact anchor URL. | Review engagement and observed time. |
| “Group coverage of the same development.” | `GET /beans/clusters` | Returns durable same-subject groups of Beans articles. | `GET /beans/clusters/{cluster_id}`. |
| “What concrete development happened?” | `GET /espresso/events` | Returns structured Event-family records with context and impact. | `GET /espresso/events/{event_id}` or `/evidence`. |
| “What exact activity or metric changed?” | `GET /espresso/actions` | Returns atomic observations; expose after the action read contract ships. | Follow `/espresso/actions/{action_id}/events`. |
| “What broader pattern or risk is emerging?” | `GET /espresso/signals` | Returns synthesized conclusions across Events and time. | `GET /espresso/signals/{signal_id}/events`. |
| “Which records support this Event?” | `GET /espresso/events/{event_id}/evidence` | Expands bidirectional `SAME_AS` evidence without exposing graph direction. | Open each Event-family record. |
| “Which Events feed this Signal?” | `GET /espresso/signals/{signal_id}/events` | Resolves stored `DERIVED_FROM` relationships in user language. | Retrieve source metadata. |
| “Summarize a bounded set of Events.” | `GET /espresso/events/summary` (conditional) | Provides an explicit aggregation only when the implementation and contract exist. | Use `/events` for record retrieval. |


## A three-question chooser for humans and agents

1. Need a publisher URL, title, summary/full text, article trend, source list, Cluster, or social/forum mention? Use **Beans**.
2. Need an atomic Action, structured Event-family record, Signal, impact, driver, forecast, or relation evidence? Use **Espresso**.
3. Need both? Start with Espresso to identify the structured development or Signal, then use its source URLs/subjects to query Beans. Treat that as a search handoff, not a shared-ID join.

Never choose by the word “related” alone: Beans `/related` means semantic reading; Espresso `/evidence` means bidirectional same-development evidence.



## Documentation information architecture

### Landing and overview pages

#### `docs/pages/introduction.mdx`

Change:

- Add the two-sentence product distinction before any endpoint list.
- Add two product cards titled “Publisher evidence” and “Structured
  intelligence.”
- Give each card three example user questions and one “not for” statement.
- Link to the product chooser in `api-overview.mdx`.

Documentation change only.

#### `docs/pages/api-overview.mdx`

Change:

- Make the product-boundary table the first substantive section.
- Add the decision guide and three-question chooser.
- Show the hierarchy:

```text
Beans
  Source -> Article -> Related reading
                    -> Mentions
  Cluster -> Member articles -> Cluster source coverage

Espresso
  Action -> Event-family record -> Signal
  Event-family record <-> Evidence (`SAME_AS`)
  Signal -> supporting Events (`DERIVED_FROM`)
```

- State that the diagram is conceptual. It does not claim cross-database foreign
  keys.
- Put authentication and shared conventions after product selection.

Documentation change only.

### Product how-to pages

#### `docs/pages/howtos/beans-howto.mdx`

Required opening:

> Use Beans when the answer should contain publisher articles or source evidence. Each article is a citable document; a Cluster groups coverage of one development. Beans does not synthesize an Event or forecast.

Required order:

1. Search `GET /beans/articles`; retrieve one with `GET /beans/articles/{article_id}`.
2. Use `GET /beans/headlines` for a current high-attention briefing and `GET /beans/articles/trending` for ranked article attention.
3. Discover `GET /beans/sources`, `/beans/sources/{source_id}`, `/beans/categories`, `/beans/entities`, and `/beans/regions`.
4. Use `GET /beans/articles/{article_id}/related` for semantic reading recommendations.
5. Use `GET /beans/articles/{article_id}/mentions` for external posts linking the exact URL.
6. Use `GET /beans/clusters`, `/beans/clusters/{cluster_id}`, `/articles`, and `/sources` for same-development grouping and editorial source coverage.
7. Treat `GET /beans/articles/counts` as deferred until a named aggregation consumer and cost contract exist.

For every route include a literal “Use when” question, a “Do not use when” alternative, exact defaults/time fields, request and success/empty/error examples, and the next route. Use `source`/`source_id` in public docs; retain `publisher` only when describing legacy compatibility.

Documentation change; OpenAPI descriptions must be updated with each public route packet.



#### `docs/pages/howtos/espresso-howto.mdx`

Required opening:

> Use Espresso when the answer should contain a structured activity, Event-family development, impact, or cross-event Signal. Espresso is not the article/full-text API.

Required order:

1. Choose Action, Event-family record, or Signal with a three-row comparison.
2. Search `GET /espresso/events` and retrieve `GET /espresso/events/{event_id}`.
3. Use `GET /espresso/events/{event_id}/evidence` for bidirectional `SAME_AS` evidence; never ask the user to reason about edge direction.
4. Use `GET /espresso/signals` and `GET /espresso/signals/{signal_id}`.
5. Follow `GET /espresso/signals/{signal_id}/events` for supporting Events and `GET /espresso/events/{event_id}/signals` for Signals using an Event.
6. Use `GET /espresso/sources`, `/espresso/sources/{source_id}`, and `GET /espresso/tags` for provenance and filters.
7. Document `GET /espresso/actions`, `/actions/{action_id}`, and both action/Event traversals only as planned until the action read contract is available.
8. Document `GET /espresso/events/summary` only as a conditional bounded aggregation; do not call it generic analytics.

Say “other records representing the same development,” “Events supporting this Signal,” and “Signals using this Event.” Keep storage direction in an internal ADR, not the user guide.

Documentation change; route/spec work remains separate.



### MCP guide

#### `docs/pages/howtos/mcp-howto.mdx`

Put the product chooser before connection setup and keep catalogs separate.

Beans — publisher evidence: `search_articles`, `get_article`, `get_trending_articles`, `get_headlines`, `find_related_articles`, `get_article_mentions`, `search_clusters`, `get_cluster`, `get_cluster_articles`, `get_cluster_sources`, `list_sources`, `get_source`, `list_categories`, `list_entities`, and `list_regions`.

Espresso — structured intelligence: `search_events`, `get_event`, `get_event_evidence`, `get_event_signals`, `search_signals`, `get_signal`, `get_signal_events`, `list_intelligence_sources`, `get_intelligence_source`, and `list_tags`. Planned action tools are `search_actions`, `get_action`, `get_event_actions`, and `get_action_events`; `summarize_events` is conditional.

Include prompts that demonstrate first-tool selection: full article text → Beans; broader emerging risk → Espresso; “same development evidence” → Espresso `/evidence`; “which outlets covered it” → Beans Cluster sources. Document multi-product handoffs as search/URL transfers, not shared UUID joins.

Documentation and gateway OpenAPI/tool-description change.



## OpenAPI documentation changes

The route proposal documents define future paths. For every route that becomes
public, update both the backend annotations/generated Swagger and the gateway
OpenAPI in one work packet.

### Product-level descriptions

Beans OpenAPI description must say:

> Search and retrieve source articles, article trends, semantic related reading,
> external mentions, and same-development Clusters. Responses are citable media
> records; use Espresso for structured Events and Signals.

Espresso OpenAPI description must say:

> Search and retrieve structured Actions, Event-family records, and Signals.
> Responses encode activities, context, impacts, forecasts, and transparent
> evidence relationships; they are not full publisher articles.

### Operation-description template

Every operation description must answer, in this order:

1. **Answers:** one literal user question.
2. **Returns:** the exact resource and one-record meaning.
3. **Use when:** a positive selection example.
4. **Do not use when:** the closest product/route alternative.
5. **Search/filter behavior:** searched fields and combination rules.
6. **Time:** filtered timestamp, default window, and timezone.
7. **Sort/pagination:** default sort, enum, maximum page size, and cursor.
8. **Missing fields:** nullable/restricted behavior.
9. **Next step:** the most likely detail/evidence route.

Bad description:

> Returns semantic results with scalar filtering and analytics.

Good Beans description:

> Answers “Which publisher articles discuss this subject?” Returns one record
> per article. Use Espresso `/espresso/events` instead when the desired output is a
> structured Event-family development with impact and context.

Good Espresso description:

> Answers “Which self-contained developments happened?” Returns one structured
> Event-family record per result. Use Beans `/beans/articles` instead when the
> desired output is publisher URLs, article text, or citations.

## Example documentation flows

### Beans-only flow: source research

Question: “What did publishers report about export controls, and which article received the most attention?”

```text
1. GET /beans/articles -> source documents
2. GET /beans/articles/trending -> ranked article attention
3. GET /beans/articles/{article_id}/mentions -> external posts for the exact URL
4. GET /beans/clusters/{cluster_id}/sources -> distinct editorial sources, when clustered
```

Expected answer: cited URLs, article summaries/full text where permitted, and transparent attention/coverage fields. No synthesized forecast.

### Espresso-only flow: intelligence research

Question: “What developments indicate demand weakness, and what broader risk is emerging?”

```text
1. GET /espresso/events -> structured developments
2. GET /espresso/signals -> cross-event conclusions
3. GET /espresso/signals/{signal_id}/events -> supporting Events
4. GET /espresso/events/{event_id}/evidence -> same-development evidence when needed
```

Expected answer: actions/context/impact, Signal drivers and forecast, and provenance. Do not return full article content as though Espresso were Beans.

### Cross-product flow

Question: “Explain the emerging risk and give me publisher coverage I can cite.”

```text
1. GET /espresso/signals -> choose the conclusion
2. GET /espresso/signals/{signal_id}/events -> identify supporting developments and source URLs
3. GET /beans/articles -> retrieve citable publisher documents using subjects, sources, or URLs
4. GET /beans/clusters/{cluster_id}/sources -> inspect editorial coverage when a Cluster is available
```

Call step 3 a search/URL handoff, not a guaranteed shared-ID join.



## Error and empty-result documentation

Both products should use the same error envelope, but examples must remain
product-specific.

```json
{
  "error": {
    "code": "invalid_parameter",
    "message": "sort must be one of: published_at, relevance, trend_score",
    "parameter": "sort",
    "request_id": "req_123"
  }
}
```

Rules to document:

- `200` plus `data: []` means no matching records.
- `400` means an input cannot be interpreted; name the parameter and allowed
  values.
- `404` applies to UUID detail routes.
- `429` includes retry guidance.
- `500/503` distinguishes database/embedder dependency failures where useful.
- Errors never suggest switching products unless the requested output clearly
  belongs to the other product.

## Navigation changes

In `docs/zudoku.config.tsx`, organize navigation as:

```text
Start
  Introduction
  Choose Beans or Espresso
  Authentication

Beans — Publisher Evidence
  Overview
  Articles and headlines
  Sources and taxonomy
  Related reading and mentions
  Clusters and source coverage
  MCP workflows

Espresso — Structured Intelligence
  Overview
  Events and evidence
  Signals and supporting Events
  Sources and tags
  Actions (planned until released)
  Event summaries (conditional)
  MCP workflows
```

Do not create navigation links to planned API-reference routes until the OpenAPI operation exists. Label conditional and planned capabilities visibly.



## Artifact-by-artifact change plan

| Artifact | Required change | Change type | Priority | Acceptance |
|---|---|---|---:|---|
| `docs/pages/introduction.mdx` | Add two-sentence distinction and product cards. | Documentation | P0 | A reader states the difference after the first screen. |
| `docs/pages/api-overview.mdx` | Add boundary table, decision guide, chooser, and hierarchy. | Documentation | P0 | Ten intent fixtures select the intended product. |
| `docs/pages/howtos/beans-howto.mdx` | Reorder around article workflows and add use/do-not-use examples. | Documentation | P0 | Every live Beans route has request, response, empty, error, and next step. |
| `docs/pages/howtos/espresso-howto.mdx` | Reorder around action/event/signal and typed evidence. | Documentation | P0 | No public explanation requires graph direction. |
| `docs/pages/howtos/mcp-howto.mdx` | Split tools by product and add selection fixtures. | Documentation/spec | P0 | Tool evaluation chooses correct product and first tool. |
| `config/beans.oas.json` | Add product boundary and operation template content as routes ship. | Specification | P0 | Gateway contract matches Beans behavior and generated Swagger. |
| `config/espresso.oas.json` | Add product boundary and operation template content as routes ship. | Specification | P0 | Gateway contract matches Espresso behavior and generated Swagger. |
| `apis/beans/router/routes.go` annotations | Mirror route documentation when behavior changes. | Code comment/spec source | P0 with route work | Regenerated Swagger matches gateway path minus prefix. |
| `apis/espresso/router/routes.go` annotations | Mirror route documentation when behavior changes. | Code comment/spec source | P0 with route work | Regenerated Swagger matches gateway path minus prefix. |
| `docs/zudoku.config.tsx` | Group navigation by product purpose. | Documentation configuration | P1 | No planned route appears as live reference. |
| `BEANS_API_ROUTE_PROPOSAL.md` | Route vocabulary and Beans examples are the source of truth for planned/public route wording. | Documentation contract | P0 | Plan uses `/articles`, `/sources`, `/clusters`, `/mentions`, and top-level taxonomy names. |
| `ESPRESSO_API_ROUTE_PROPOSAL.md` | Route vocabulary, relation semantics, and conditional capabilities are the source of truth for Espresso docs. | Documentation contract | P0 | Plan uses `/evidence`, `/tags`, `/events/summary`, and action gating; no generic analytics routes. |
| New glossary page | Publish the terminology contract. | Documentation | P1 | Terms are linked from both how-tos and MCP guide. |
| Documentation examples/fixtures | Make examples executable in CI. | Test/documentation | P1 | Quickstart examples fail CI when contracts drift. |

## Execution sequence

### Documentation Stage 0 — Fix selection before adding routes

1. Update introduction and API overview.
2. Publish terminology and the decision guide.
3. Rewrite current Beans/Espresso how-to openings and current-route examples.
4. Split the MCP catalog and add product-selection fixtures.

These are documentation changes. They do not require a query or schema change.

### Documentation Stage 1 — Define future contracts

1. Treat `BEANS_API_ROUTE_PROPOSAL.md` and `ESPRESSO_API_ROUTE_PROPOSAL.md` as the route vocabulary source of truth.
2. Document Beans UUID article/source identifiers and Espresso UUID Event/Signal/Action identifiers without exposing storage internals.
3. Document `SAME_AS` through Event evidence and `DERIVED_FROM` through supporting Event lookups; keep direction out of user-facing prose.
4. Label action routes and conditional Event summaries as planned until their read contracts exist.

Mostly documentation/specification work; route behavior is delivered in the
product execution plans.

### Documentation Stage 2 — Ship docs with each route

For each route packet:

1. update handler annotations;
2. regenerate backend Swagger;
3. update gateway OpenAPI;
4. update the product how-to;
5. update MCP tool descriptions if exposed;
6. add executable examples and selection fixtures.

A route is not complete if only one of these artifacts changes.

### Documentation Stage 3 — Measure comprehension

Track:

- documentation searches containing “difference,” “story,” “event,” “signal,”
  “source,” and “related”;
- API validation errors by operation/parameter;
- MCP wrong-product selection rate;
- abandoned quickstarts;
- follow-up sequences between search/detail/evidence operations;
- support questions that reveal terminology confusion.

Use those signals to simplify wording. Do not add more route families to solve a
documentation problem.

## Documentation acceptance tests

Create fixtures that assert the expected product and first route/tool:

| Prompt | Expected product/tool |
|---|---|
| “Find the BBC articles about the policy change.” | Beans `search_articles` → `/beans/articles` |
| “Give me the full text for this article UUID.” | Beans `get_article` → `/beans/articles/{article_id}` |
| “Which outlets covered this development?” | Beans `get_cluster_sources` → `/beans/clusters/{cluster_id}/sources` |
| “Which social posts linked this article?” | Beans `get_article_mentions` → `/beans/articles/{article_id}/mentions` |
| “What happened to the company's guidance?” | Espresso `search_events` → `/espresso/events` |
| “What exact value changed?” | Espresso `search_actions` when available → `/espresso/actions` |
| “What risks are emerging across suppliers?” | Espresso `search_signals` → `/espresso/signals` |
| “Which Events support this conclusion?” | Espresso `get_signal_events` → `/espresso/signals/{signal_id}/events` |
| “Show other records of the same development.” | Espresso `get_event_evidence` → `/espresso/events/{event_id}/evidence` |
| “Show related reading for this article.” | Beans `find_related_articles` → `/beans/articles/{article_id}/related` |
| “Summarize these Events.” | Espresso `summarize_events` only when conditional route is released |

Definition of done:

- At least 90% of reviewed human testers choose the correct product from the
  overview alone.
- MCP selection fixtures choose the correct product and first tool.
- No live example references a planned route.
- No Espresso relation example exposes direction.
- No Beans example describes a trend score as an intelligence signal.
- No cross-product example claims a shared UUID or foreign-key join that does
  not exist.


## Canonical Beans documentation strings for OpenAPI and MCP

This section supersedes older Beans references in this plan to `publishers`,
article-level `coverage`, `/article-clusters`, and `/tags/*`. The canonical
public vocabulary is **sources**, **mentions**, **clusters**, and separate
**categories**, **entities**, and **regions** routes.

Use these strings verbatim in gateway OpenAPI, generated Swagger annotations,
portal pages, and MCP tool descriptions. For unreleased routes, show the same
text only in a visibly **planned** capability section; do not expose it as a
live operation.

### Beans API description (`info.description`)

> Beans is a read-only media-evidence API for finding and verifying what news
> sources, publishers, blogs, and other outlets published. Its main resource is
> an article: one attributable source document with a UUID, canonical URL, title,
> publication time, source, and optional extracted content.
>
> Use Beans when the answer needs article URLs, citations, source identity,
> article summaries/full text, semantic related reading, social/forum posts
> linking to an article, or multi-source article clusters. Use Espresso instead
> when the answer needs a structured action, event, impact, driver, or forecast.
> Beans does not assert a synthesized real-world conclusion from the articles it
> returns.
>
> Choose `/articles` for individual documents, `/articles/{id}/related` for
> semantic follow-up reading, `/articles/{id}/mentions` for external posts that
> linked to one article URL, and `/clusters` for a durable same-subject group.
> Editorial coverage means the sources represented in a cluster and is returned
> by `/clusters/{id}/sources`; it is not an article-level social-share route.

### Shared documentation strings

| Subject | OpenAPI / MCP description |
|---|---|
| Read-only behavior | `Beans is read-only. Repeating a successful GET request does not alter articles, sources, clusters, or mention records.` |
| `limit` and `cursor` | `limit is the maximum records returned; default 20, maximum 100. cursor is an opaque continuation token from pagination.next_cursor. Do not create, modify, decode, or reuse it with different filters or sort order.` |
| `from` and `to` | `from is inclusive and to is exclusive. Both accept RFC 3339 timestamps; a date-only value means midnight UTC. Read the route description to learn which timestamp is filtered.` |
| Empty result | `200 with data: [] means the request was valid but matched no records. It is not an error and does not prove that an input UUID or filter value is invalid.` |

### Route operation-description strings

#### `GET /beans/articles`

> **Answers:** “Which source-published articles match this subject or filter?”
> **Returns:** one attributable source document per result, including UUID,
> canonical URL, publication time, source summary, and extracted taxonomy.
> **Use when:** you need article URLs, summaries, source evidence, or semantic
> search. **Do not use when:** you need a synthesized conclusion (use Espresso)
> or a durable same-subject group (use `/clusters`). With `q`, results are
> semantic-search matches; without `q`, results are a publication-date browse
> feed. Use `/articles/trending` for attention ranking.

#### `GET /beans/articles/{article_id}`

> **Answers:** “What is the complete API record for this known article UUID?”
> **Returns:** one source-published article. Full extracted content is returned
> only when `include=content` is requested and available. **Use when:** you have
> an article UUID from a list, related, trending, headline, or cluster response.
> **Do not use when:** you need an article group; use `/clusters`.

#### `GET /beans/articles/trending`

> **Answers:** “Which individual articles have the strongest current measured
> attention?” **Returns:** article summaries plus observed attention inputs and a
> ranking score. **Use when:** you need high-attention publisher coverage.
> **Do not use when:** you need newest-first articles (use `/articles`) or a
> multi-source same-subject group (use `/clusters`).

#### `GET /beans/headlines`

> **Answers:** “Which recently published, high-attention articles belong in a
> current briefing?” **Returns:** a fixed 24-hour attention-selected feed.
> **Use when:** you need a concise current headline feed. **Do not use when:**
> you need a caller-selected time range or human editorial curation.

#### `GET /beans/articles/{article_id}/related`

> **Answers:** “What other articles are semantically useful to read after this
> one?” **Returns:** relevance-ranked article summaries. Every result shares
> `meta.relation_type=semantic_related`; the relation is not repeated per item.
> **Use when:** you need adjacent, follow-up, or contextually similar reading.
> **Do not use when:** you need a complete same-subject group (use `/clusters`)
> or social posts linking to the article (use `/articles/{id}/mentions`). Results
> may be from the same source on another day or from another source.

#### `GET /beans/articles/{article_id}/mentions`

> **Answers:** “Which observed social or forum posts linked to this exact article
> URL?” **Returns:** external post URLs, normalized platform/forum, observation
> time, and available engagement metrics. **Use when:** you need article-specific
> social/forum evidence. **Do not use when:** you need other source reporting
> (use `/clusters/{id}/articles`) or semantic follow-up reading (use `/related`).

#### `GET /beans/sources` and `GET /beans/sources/{source_id}`

> **Answers:** “Which article-producing sources are available?” and “What is
> known about this source UUID?” **Returns:** profiles for publishers, outlets,
> blogs, feeds, or domains. **Use when:** you need source discovery or UUIDs for
> `source_ids`. **Do not use when:** you need social platforms; platforms appear
> only on article mention records. To retrieve a source’s articles, call
> `/articles?source_ids={source_id}`.

#### `GET /beans/categories`, `GET /beans/entities`, `GET /beans/regions`

> **Answers:** “Which exact values can I pass to this article filter?”
> **Returns:** one taxonomy family only: categories, extracted entities, or
> regions. **Use when:** you need valid filters before calling `/articles` or
> `/clusters`. **Do not use when:** you need a mixed tag list; the three families
> intentionally have different meanings and filter behavior.

#### `GET /beans/clusters`, `GET /beans/clusters/{cluster_id}`

> **Answers:** “Which durable article groups describe the same developing
> subject?” and “What is the current scope of this group?” **Returns:** cluster
> summaries with member-article count, editorial-source count, publication
> bounds, representative article, and aggregate taxonomy. **Use when:** you need
> to reduce duplicate reporting and follow one development across sources.
> **Do not use when:** you only need semantic reading (use `/related`) or social
> posts (use `/mentions`).

#### `GET /beans/clusters/{cluster_id}/articles`

> **Answers:** “Which source-published articles are evidence for this cluster?”
> **Returns:** paginated attributable article summaries. **Use when:** you need
> URLs and sources behind a cluster claim. **Do not use when:** you need only a
> source-level coverage summary; use `/clusters/{id}/sources`.

#### `GET /beans/clusters/{cluster_id}/sources`

> **Answers:** “Which editorial sources are represented in this cluster, and how
> much evidence does each contribute?” **Returns:** one source aggregate per
> source with article count and first/last publication times. **Use when:** you
> need multi-source editorial coverage. **Do not use when:** you need social
> mentions or engagement; use `/articles/{id}/mentions`.

#### `GET /beans/articles/counts`

> **Answers:** “How many matching articles exist in each allowed group?”
> **Returns:** numeric aggregate values only. **Use when:** an approved dashboard,
> alert, or report needs bounded counts. **Do not use when:** you need article
> records, trends, or semantic results. This route is planned until a named
> consumer and cost policy exist.

### Query-parameter documentation strings

| Parameter | Route(s) | OpenAPI description |
|---|---|---|
| `q` | `/articles` | `Optional semantic search text. Searches article meaning, not an exact keyword index. Default sort is relevance when q is supplied. Maximum 512 characters.` |
| `q` | `/sources` | `Optional case-insensitive source-name or source-domain prefix search. This is not article-content search.` |
| `q` | `/categories`, `/entities`, `/regions` | `Optional case-insensitive prefix search over valid filter values. Copy returned values into the matching article filter.` |
| `q` | `/clusters` | `Optional search over durable cluster name, summary, and aggregate entities. This does not search every member article body directly.` |
| `source_ids` | article, trending, headline, related, cluster routes | `Optional comma-separated source UUIDs. A source is the publisher/outlet/domain that produced an article. Multiple values are ORed within this filter.` |
| `content_types` | `/articles`, `/articles/trending` | `Optional comma-separated public article kinds: news, blog, or post. Omit for all public kinds.` |
| `categories`, `entities`, `regions` | matching article/cluster routes | `Optional comma-separated exact values from GET /categories, GET /entities, or GET /regions. Values are ORed within one filter family.` |
| `from`, `to` | article, trending, related, cluster, cluster-member, count routes | `Inclusive/exclusive article published_at range. On /clusters, a cluster matches when at least one member falls in the range.` |
| `from`, `to` | `/articles/{id}/mentions` | `Inclusive/exclusive mention observed_at range, not the anchor article publication time.` |
| `sort` | `/articles` | `published_at or relevance. relevance is valid only with q. Default: relevance with q, published_at without q.` |
| `sort` | `/articles/{id}/mentions` | `observed_at or engagement. Default observed_at. Engagement ordering uses only available normalized metrics.` |
| `sort` | `/clusters` | `updated_at, article_count, or relevance. relevance is valid only with q.` |
| `sort` | cluster-member/source routes | `For members: published_at or representative_first. For source coverage: article_count or last_published_at.` |
| `include` | `/articles`, `/articles/{id}` | `Optional comma-separated expansions. content requests extracted full text when available and can substantially increase response size.` |
| `exclude_anchor_source` | `/articles/{id}/related` | `When true, omit the anchor article source. Default false; same-source follow-up reporting is otherwise allowed.` |
| `platforms` | `/articles/{id}/mentions` | `Optional normalized social/forum platforms, such as reddit or x. Filters mention records, not article sources.` |
| `min_sources` | `/clusters` | `Minimum distinct editorial sources in the complete cluster. This is not a social-platform count.` |
| `domains` | `/sources` | `Optional exact source domains. Use q for source-name/domain prefix discovery.` |
| `group_by` | `/articles/counts` | `Required aggregation dimension: published_day, source, category, or region. The response contains numbers only.` |
| `limit`, `cursor` | all collections | `limit defaults to 20 and cannot exceed 100. cursor is the opaque pagination.next_cursor token; do not modify it or change filters/sort between pages.` |

### Response-field documentation strings

| Field | OpenAPI description |
|---|---|
| `data` | `Requested records. An empty array means the request was valid but matched no records.` |
| `pagination.limit` | `Effective page size used for this response.` |
| `pagination.next_cursor` | `Opaque next-page token, or null when no additional page exists. Reuse only with identical filters and sort.` |
| `meta.as_of` | `Freshness time for aggregate, trend, mention, or cluster data. It is not necessarily article publication time.` |
| `article.id` | `Stable UUID for this source-published document. Use it for article detail, related-article, and article-mention routes.` |
| `article.url` | `Canonical external source URL. Use it for citation/browser navigation; use article.id for API resource lookup.` |
| `article.title` | `Publisher-provided or extracted headline.` |
| `article.summary` | `Compact article summary. It is not an Espresso event summary or a system conclusion.` |
| `article.published_at` | `When the source published the article, in UTC.` |
| `article.content_type` | `Public article kind, such as news, blog, or post.` |
| `article.content` | `Extracted full text when requested and available. Null/absence does not mean the article does not exist.` |
| `article.content_access.status` | `available when extracted content can be returned; unavailable when article metadata exists but full text cannot be returned.` |
| `article.source` | `Compact publisher/outlet/blog/domain profile. It is not a social platform.` |
| `source.id`, `source.name`, `source.domain`, `source.url` | `Stable API identifier, human-readable name, normalized primary domain, and homepage URL for an article-producing source.` |
| `categories`, `entities`, `regions` | `Exact normalized taxonomy values on an article. Discover valid values through the corresponding taxonomy route before filtering.` |
| `trend.like_count`, `trend.comment_count` | `Observed aggregate engagement from collected mention records; null when unavailable.` |
| `trend.mention_count` | `Collected social/forum mention-record count contributing to attention. It is not necessarily a platform-reported share total.` |
| `trend.related_article_count` | `Qualifying related-article relationship count used by trend calculation; it does not prove cluster membership.` |
| `trend.score` | `Internal ranking score. Compare only within one response/as_of time, not as universal popularity.` |
| `meta.relation_type` | `Relation shared by every related-article item. Current value semantic_related means relevance-ranked reading, not same-subject cluster membership.` |
| `meta.anchor_article_id` | `UUID named in the route path; included so stored agent results remain self-describing.` |
| `mention.url` | `External post URL observed linking to the anchor article URL.` |
| `mention.platform`, `mention.forum` | `Normalized social/forum platform and optional community/group context. Neither is an article source.` |
| `mention.observed_at` | `When Beans observed/collected the mention; it can differ from post creation time.` |
| `mention.engagement.*` | `Observed lower-bound platform engagement at collection time; null when the metric is unavailable.` |
| `cluster.id` | `Stable UUID for a durable same-subject article group.` |
| `cluster.name`, `cluster.summary` | `Current label and media-coverage summary for the group; not an Espresso event conclusion.` |
| `cluster.first_published_at`, `cluster.last_published_at` | `Earliest and latest publication time among current cluster members.` |
| `cluster.updated_at` | `When cluster summary or membership was last refreshed.` |
| `cluster.article_count` | `Total current member articles, not only the current response page.` |
| `cluster.source_count` | `Distinct editorial sources in the cluster; excludes social/forum platforms.` |
| `cluster.representative_article` | `One attributable article selected to represent the cluster. Retrieve member articles for full evidence.` |
| `coverage.source`, `coverage.article_count`, `coverage.first_published_at`, `coverage.last_published_at` | `One source contribution to cluster evidence: source profile, member count, and earliest/latest member publication time.` |
| `taxonomy.value`, `taxonomy.label` | `Exact normalized filter string and optional human-readable display form. Copy value, not label, into filters when they differ.` |
| `counts.group_by`, `counts.values[].key`, `counts.values[].count` | `Requested aggregate dimension, group value, and matching article count. Counts never return article records.` |

### Agent selection and follow-up strings

Use these exact sentences in Beans MCP tool descriptions and how-to callouts:

- `After /articles, call /articles/{id} for one record, /related for semantic follow-up reading, /mentions for social/forum post evidence, or /clusters for the same developing subject across sources.`
- `Do not call /related when the user asks which outlets covered the same development; call /clusters, then /clusters/{id}/sources or /clusters/{id}/articles.`
- `Do not call /mentions for generic social discussion of a topic. It is limited to posts that linked to the anchor article URL.`
- `Do not call Beans when the requested answer is an interpreted action, event, impact, driver, or forecast. Start with Espresso.`
