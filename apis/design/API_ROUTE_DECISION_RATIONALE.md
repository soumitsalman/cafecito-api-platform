# API Route Decision Rationale

Status: decision record for the latest Beans and Espresso route proposals  
Prepared: 2026-08-04  
Scope: public route vocabulary, payload boundaries, query contracts, and documentation decisions  
Source documents: [BEANS_API_ROUTE_PROPOSAL.md](BEANS_API_ROUTE_PROPOSAL.md), [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md), [BEANS_ESPRESSO_DOCUMENTATION_PLAN.md](BEANS_ESPRESSO_DOCUMENTATION_PLAN.md)

## Why this document exists

The route proposals went through several revisions. Earlier versions mixed current
handlers, ideal future resources, and speculative analytics fields. This record
explains why the latest decisions were made, what problem each decision solves,
and which assumptions remain deliberately unresolved. It is a rationale document,
not a replacement for the route contracts or an implementation authorization.

The governing product distinction is:

> **Beans finds and verifies what publishers published.**

> **Espresso explains what happened and what it may mean.**

That distinction is the test for every route: if the requested output is a citable
publisher document or media coverage, it belongs to Beans; if it is a structured
activity, development, implication, or cross-event conclusion, it belongs to
Espresso.

## Decision summary

| Decision | Latest choice | Why it was chosen | What it prevents |
|---|---|---|---|
| Beans publisher noun | `source` | NewsAPI and Perigon use source-oriented vocabulary; Beans data represents outlets and media origins. | Treating an outlet as a person or organization model with unsupported metadata. |
| Beans taxonomy routes | `/categories`, `/entities`, `/regions`; legacy `/tags/*` retained | These are different dimensions with different filter semantics; fixed routes are easier for OpenAPI and MCP clients. | A generic tag bag whose values cannot be interpreted consistently. |
| Beans article collection | `/articles` | One resource-oriented collection can support semantic search and publication browsing. | Separate `/search` and `/latest` contracts that drift apart. |
| Beans attention feed | `/articles/trending` | Trend is a distinct ranking question driven by measured attention, not publication recency. | Making `trend_score` a misleading general sort option. |
| Beans briefing feed | `/headlines` | A named briefing surface can have its own recency and eligibility rules. | Claiming that every recent or trending article is a headline. |
| Beans semantic neighbors | `/articles/{article_id}/related` | “Related” is useful for ranked reading discovery, including other sources or later follow-ups. | Equating semantic similarity with same-subject membership or social sharing. |
| Beans external social evidence | `/articles/{article_id}/mentions` | “Mention” describes an external post linking to the exact article URL. | Calling social observations editorial coverage or publisher propagation. |
| Beans same-subject grouping | `/clusters` | Cluster is explicit about grouping and aligns with clustering capabilities in comparable news products. | Using “story” for both a single article and a multi-article group. |
| Beans editorial coverage | `/clusters/{cluster_id}/sources` | Coverage is meaningful at cluster/source level: which outlets are represented. | Adding an ambiguous article-level `/coverage` route. |
| Espresso Event collection | `/events` includes every `kind LIKE 'event%'` record | Stored Event-family records include curated and source-oriented event representations; the route must reflect the storage contract. | Pretending all event records share one curated digest schema. |
| Espresso Event evidence | `/events/{event_id}/evidence` | The caller wants source-specific records that represent/support a development; the route hides `SAME_AS`. | Making `/equivalents` the primary user workflow and exposing graph jargon. |
| Espresso evidence scope | Direct `SAME_AS` neighbors by default | Safe bounded query semantics without assuming a transitive graph or canonical-group table. | Unbounded recursive traversal and false claims of complete equivalence coverage. |
| Espresso Signal traversal | `/signals/{signal_id}/events` and `/events/{event_id}/signals` | Typed routes express the user’s desired resource while the server handles `DERIVED_FROM` direction. | Requiring users or agents to understand `from_id`, `to_id`, or edge orientation. |
| Espresso tags | `/tags` | The existing route is short and adequate for current scalar tag vocabulary. | A longer `/filter-values/tags` abstraction with no additional meaning. |
| Espresso Actions | `/actions` remains required but action-gated | Actions are a core product concept, while their storage/read contract is a separate workstream. | Inventing an Action schema or planning its ingestion in the route document. |
| Analytics | No generic `/analytics/*` routes in the canonical target | Aggregates need named questions, bounded dimensions, and a known consumer. | Routes whose names do not tell an agent what is being counted or summarized. |

## Beans decisions in detail

### `sources` instead of `publishers`

The change is more than a rename. A public `Source` object should contain the
stable source identifier, display name, domain, and nullable metadata. Article
responses embed the same source shape, and article filters use `source_ids`.
This follows the vocabulary used by NewsAPI and Perigon and avoids promising
publisher-specific fields that the current model does not guarantee.

The UUID assumption remains a prerequisite: the proposal assumes article and
source IDs are stable UUID primary keys, but it does not plan the separate UUID
backfill/primary-key workstream.

### Fixed taxonomy routes instead of a generic tag endpoint

Categories, entities, and regions are not interchangeable labels:

- `/categories` returns values usable by the `categories` article filter.
- `/entities` returns values usable by the `entities` article filter.
- `/regions` returns values usable by the `regions` article filter.

Each response should expose a machine value and a display label. Counts are not
included by default because their freshness and filter scope would need a separate
contract. The existing `/tags/categories`, `/tags/entities`, and `/tags/regions`
remain migration aliases.

### `related`, `mentions`, and `clusters` are intentionally non-overlapping

| Route | Input anchor | Query meaning | Response meaning |
|---|---|---|---|
| `/articles/{id}/related` | One article UUID | Rank semantically useful reading; publication `from`/`to` constrain returned articles. | Plain `ArticleSummary` records; invariant `relation_type` appears once in `meta`, not on every item. |
| `/articles/{id}/mentions` | One article UUID/canonical URL | Find external social/forum posts linking that exact URL; `from`/`to` constrain `observed_at`. | Mention URL, platform/forum, observation time, and only engagement fields actually supplied. |
| `/clusters` | Subject/search/filter, not necessarily one article | Group articles that describe the same developing subject; `from`/`to` constrain member publication time. | Durable cluster identity, summary, article/source counts, and links to member articles/sources. |
| `/clusters/{id}/sources` | One cluster UUID | Determine which editorial sources are represented by cluster members. | Source records with article counts and publication span; no social metrics. |

The payload differences are part of the design. A client should be able to choose
a route from the question and then validate the returned object without reading
implementation details.

### Why `/clusters` is not `/stories`

Comparable products use “stories” for grouped coverage, but “story” is ambiguous
in a media API because a user may call one article a story. `clusters` states the
operation: grouping records by subject. The route is only justified once the
backend can provide a stable cluster identity, membership rules, representative
article, merge/split behavior, and bounded pagination. Until then, it is a
future read model rather than a cosmetic alias for `related_beans`.

### Why `/articles/trending` and `/headlines` remain separate

These routes answer different questions:

- Trending: “Which articles have the strongest measured attention?”
- Headlines: “Which recently published high-attention articles belong in a
  current briefing?”

Trending is attention-ranked and exposes trend inputs plus `meta.as_of`. Headlines
has a fixed-window selection contract. Neither route is a synonym for `/articles`
with a different undocumented sort.

## Espresso decisions in detail

### Why `/events` means the Event family

The latest proposal deliberately defines membership using the stored predicate:

```sql
kind LIKE 'event%'
```

This includes `event`, `event:news`, `event:blog`, `event:post`, `event:site`, and
`event:social`. The public response preserves both the exact `kind` and a derived
`representation`, so an agent can distinguish a curated Event from a source-oriented
Event record. Optional richer fields may be present for `kind=event`, but the API
must not imply that all Event-family digests are structurally identical.

This decision was made after identifying that a canonical-only `/events` contract
would silently hide source-level event records that existing data and relations
already use.

### Why `evidence` is preferred to `equivalents`

`SAME_AS` is an implementation relationship. A caller’s question is usually one
of these:

- “Which source records represent this development?”
- “What URL and source support this record?”

`GET /events/{event_id}/evidence` names that customer value and keeps relation
direction invisible. `/events/{event_id}/equivalents` may be retained as a
compatibility or diagnostic view, but it should not be the primary route.

The event detail response carries the requested record’s direct `url` and `source`.
The evidence response returns the other direct `SAME_AS` Event-family records,
each with its own URL, source, kind/representation, timestamp, and compact
briefing. The default scope is explicitly `direct_same_as`; recursive/transitive
expansion is not promised until the data model guarantees a bounded, maintained
equivalence group.

### Why relation direction is hidden

`SAME_AS` is bidirectional by design. `DERIVED_FROM` is directional in storage,
but users do not ask for “outgoing edges.” They ask for resources:

- Signal → supporting Events: `/signals/{signal_id}/events`.
- Event → Signals that use it: `/events/{event_id}/signals`.
- Event → supporting Actions: `/events/{event_id}/actions`.
- Action → Events that use it: `/actions/{action_id}/events`.

The API performs the correct direct or inverse query internally. This makes the
contract stable if storage orientation changes and keeps MCP tool descriptions
human-readable.

### Why sources are P0

A structured conclusion without provenance is difficult to trust or cite. Espresso
therefore treats `/sources` and source fields on Event/Evidence responses as part
of the core contract, not an optional enrichment. Missing source metadata remains
nullable; the API does not fabricate names, domains, or URLs.

### Why analytics are deferred and named only as `/events/summary`

Earlier `/analytics/events` and `/analytics/signals` proposals described outputs
without establishing the question, grouping dimensions, timestamp semantics, or
consumer. The latest design keeps only a conditional `/events/summary` because it
can be defined around a specific bounded question such as “How many Event-family
records occurred per created day?” It is not a generic analytics namespace, and
there is no equivalent generic `/signals/summary` route until a concrete consumer
and schema exist.

### Why Actions remain in the target surface

Actions are necessary to express atomic facts and to make Event composition useful.
However, a separate workstream owns how Actions are stored and clarified in
EspressoDB. The route proposal therefore reserves `/actions` and its traversals,
but blocks publication until the read contract defines IDs, timestamps, typed
observations, provenance, and relation behavior. No Action storage work is hidden
inside the API route plan.

## Payload and query choices

### Common response envelope

The proposals converge on a conventional envelope:

```json
{
  "data": [],
  "pagination": { "limit": 20, "next_cursor": null },
  "meta": { "as_of": "2026-08-04T00:00:00Z" }
}
```

This was chosen because agents need to distinguish records, pagination state,
and freshness. It is a target contract, not the current bare-array/offset behavior.
Compatibility routes must remain explicitly labeled until the canonical contract
ships.

### Time parameters

`from` and `to` are retained because they are familiar across public APIs. Their
meaning is route-specific and documented rather than inferred:

- Beans articles, trending, clusters, and related: article `published_at`.
- Beans mentions: mention `observed_at`.
- Espresso Event-family and Signal collections: record `created_at` until an
  occurrence-time field is defined.

This avoids overloading a single time parameter with undocumented meanings while
preserving familiar public syntax.

### Provenance and nullable fields

The latest payload designs include stable IDs, source/URL provenance, and explicit
nullable fields. Missing enrichment is represented as `null`, never as fabricated
zero counts or guessed metadata. This is necessary for trustworthy citations and
agent reasoning.

## Corrections from earlier proposals

The latest decisions intentionally reverse or narrow several earlier suggestions:

1. `article-clusters` was shortened to `/clusters` after comparing route vocabulary,
   but the route remains gated on durable cluster identity.
2. `/coverage` was removed as an article-level route; social evidence is `/mentions`,
   while editorial coverage is `/clusters/{id}/sources`.
3. `/filter-values/{dimension}` was removed in favor of fixed taxonomy routes.
4. Espresso `/equivalents` was demoted in favor of `/evidence`.
5. Generic `/analytics/events` and `/analytics/signals` were rejected in favor of
   a conditional, named `/events/summary`.
6. Espresso `/events` was changed from canonical-only to the complete `event%`
   family, with `kind` preserved.
7. Repeated relation metadata was removed from Beans related-item payloads and
   placed once in response `meta`.
8. Speculative fields such as unverified content-restriction reasons, confidence,
   occurrence timestamps, and Action numeric schemas were removed or marked
   upstream/action-gated.
9. Direct `SAME_AS` evidence became the default Espresso evidence scope instead
   of promising a transitive graph expansion.

## What these decisions deliberately do not claim

- Beans does not infer real-world Events or Signals from articles.
- Espresso does not replace Beans as the article/full-text/citation API.
- A Cluster is not automatically an Espresso Event.
- A semantic related result is not guaranteed same-subject coverage.
- A social Mention is not editorial source coverage.
- An Event-family record is not necessarily a curated Event digest.
- A direct `SAME_AS` evidence response is not a complete transitive equivalence
  closure.
- `/actions` does not become available until its separate read contract exists.
- UUID migration and primary-key conversion are dependencies, not work items here.

## Decision tests for future changes

Before accepting a new route or payload field, ask:

1. What exact user question does it answer?
2. Which product owns that question: publisher evidence (Beans) or structured
   intelligence (Espresso)?
3. Does a comparable service expose the capability, and is the comparison about
   user capability rather than copied naming?
4. Does the current schema/query support the response, or is it explicitly marked
   query-change, conditional, upstream-data, or action-gated?
5. Are similarly named routes separated by input anchor, timestamp semantics,
   response object, and follow-up route?
6. Can an agent select the route without understanding SQL or relation direction?
7. Is the payload truthful about missing, nullable, and freshness fields?
8. Can the query be bounded and paginated at the expected data volume?

A proposal that cannot answer these questions should remain a design hypothesis,
not a canonical public route.

## Reference trail

- [BEANS_API_ROUTE_PROPOSAL.md](BEANS_API_ROUTE_PROPOSAL.md)
- [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md)
- [BEANS_ESPRESSO_DOCUMENTATION_PLAN.md](BEANS_ESPRESSO_DOCUMENTATION_PLAN.md)
- [API_CAPABILITY_GAP_ANALYSIS.md](API_CAPABILITY_GAP_ANALYSIS.md)
- [Beans API renovation thread](thread://019fc7ec-345f-7273-a649-67ca28bed241)
- [Espresso API renovation thread](thread://019fc805-bb9e-76e0-9f3a-03f1bd8e8cd1)
