# Espresso API Documentation Plan

Status: Documentation plan

Last updated: 2026-08-11

Source documents:

- Target contract: [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md)
- Implementation baseline: [ESPRESSO_API_IMPLEMENTATION_PLAN.md](ESPRESSO_API_IMPLEMENTATION_PLAN.md)
- Industry comparison: [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md)

## 1. Purpose and documentation rule

Document Espresso as a read-only intelligence API for two stored resources:

- **Events**: stored `sips` where `kind = "event"`.
- **Signals**: stored `sips` where `kind = "signal"`.

Public examples use gateway paths with the `/espresso` prefix. Backend and
Swagger-generation examples may use the corresponding unprefixed route.

The implementation plan is the source of truth for the first published API
reference. The route proposal defines the intended end state, but a
proposal-only capability must be marked **planned** or omitted until it is
implemented and verified. The industry reference is comparative context; it
does not introduce compatibility promises.

## 2. Release contract and proposal reconciliation

### 2.1 Document in the first published reference

| Area | First published contract |
|---|---|
| Resource kinds | `event` and `signal` only; no Actions. |
| Collection envelope | `success`, `data`, `pagination`, and `meta.as_of`. |
| Detail envelope | `success` and one `data` object. |
| Pagination | Keyset pagination with `limit` (default `20`) and opaque `cursor`. |
| Dates | `from` and `to` are ISO date-only values in the first implementation and bound record creation time (`created_at`), not an event-occurrence or article-publication time. |
| Serialization | `response_type=json`, `yaml`, or `toon`; JSON is canonical. |
| Search/order | Document only implemented `q` behavior. Do not advertise a public `sort` parameter until its cursor semantics are implemented. |
| Resource payloads | `id`, `kind`, `created_at`, `tags`, `summary`, and populated flattened digest fields. Collection records omit provenance; Event and Signal detail records may include it. |
| Relations | Event evidence remains the current narrow response until its full Event projection is implemented. Event-to-Signal and Signal-to-Event collections use their returned resource's collection shape. |
| Discovery and aggregates | Publish only when R11-R14 and R16-R18 are implemented, tested, and present in both OpenAPI documents. |

### 2.2 Proposal items to defer or label as planned

| Proposal item | Documentation treatment before implementation |
|---|---|
| RFC 3339 `from` / `to` input | Do not document; state ISO date-only input in the v1 reference. |
| `sort=recent|relevance` | Do not document as accepted. Keep examples independent of a sort choice. |
| Full provenance fields in collection items | Do not document. `url`, `base_url`, `source_id`, and `source` belong only in detail examples for the first release. |
| Returning both `briefing` and `summary` | Do not document `briefing`. The implementation plan renames it to `summary`. |
| Rich Event evidence projection | Document the current narrow evidence shape only; add the normal Event item projection after R03 changes. |
| Data-quality indicators | Do not promise `data_quality` until it is emitted and tested. |
| Stories, articles, canonical places/entities, lifecycle time, numeric metrics, and Actions | List as unsupported, not as endpoints or filters. |

### 2.3 Internal conflicts requiring one implementation decision

These are not customer-facing documentation decisions, but they must be
resolved before the generated specification is published:

| Conflict | Documentation baseline |
|---|---|
| The proposal's stable core includes provenance on every item; the implementation plan excludes it from collections. | Use the implementation-plan split: collection vs. detail. |
| The proposal preserves `briefing`; the implementation plan replaces it with `summary`. | Use `summary` only. |
| A proposal gap row says to hide `kind`, while its payload section and the implementation plan make `kind` canonical. | Retain `kind`; remove or correct the contradictory gap row before final contract sign-off. |
| The proposal calls for RFC 3339 bounds and public sorting; the implementation plan defers both. | Treat them as future additions. |
| The documents disagree about the intended JSONB scalar predicate. | Keep SQL implementation details out of public docs; verify that exact `event_types` and `impact_levels` filtering has the advertised behavior. |

## 3. Documentation deliverables and order

1. `apis/espresso/README.md`
   - Replace stale offset, text serialization, and `/related/{relationship}` examples.
   - Add a short quickstart using `GET /espresso/events`, cursor traversal, and an Event detail lookup.
   - Link to the generated OpenAPI specification and gateway API reference.

2. Espresso API reference and generated Swagger
   - Define security, common query parameters, envelopes, Event/Signal collection and detail schemas, Source schemas, and error responses.
   - Keep Go route annotations, `apis/espresso/docs/swagger.yaml`, and `apis/espresso/docs/swagger.json` aligned.

3. Gateway OpenAPI and developer portal
   - Mirror all published backend routes in `config/espresso.oas.json` with the `/espresso` prefix.
   - Add conceptual pages for Events, Signals, evidence/provenance, Sources, discovery values, and aggregates as they become available.
   - Use one authoritative parameter/response table rather than duplicating route-specific field lists that can drift from OpenAPI.

4. Examples and verification
   - Provide JSON examples for every published route. Add YAML and TOON examples only for the shared serialization behavior, not for every endpoint.
   - Validate examples against generated schemas and route tests before publishing.
   - Regenerate Swagger and update gateway OpenAPI in the same change whenever a route annotation changes.

## 4. Information architecture

The published documentation should use this order:

1. **Overview** — resource model, authentication, base path, supported formats, and the meaning of `created_at`.
2. **Common behavior** — pagination, date bounds, `q`, errors, empty collections, and response envelopes.
3. **Events** — list, detail, evidence, related Signals, count, summary, and vector search when available.
4. **Signals** — list, detail, and supporting Events.
5. **Sources and discovery** — Sources, tags, entities, regions, and event types.
6. **Relationships and provenance** — `SAME_AS`, `DERIVED_FROM`, link objects, and relationship limitations.
7. **Compatibility guide** — concise mapping to GDELT Cloud, PredictHQ, and Perigon, including non-equivalences.
8. **Known limits** — data gaps and intentionally unsupported provider-like features.

Every endpoint page should state: purpose, availability, authentication,
parameters, response shape, an example, pagination behavior where applicable,
and the relevant data limitation. Discovery and aggregate pages must explicitly
say whether values are exact stored strings rather than normalized taxonomy
values.

## 5. First-release route coverage

| Espresso route | Documentation status | Required documentation notes |
|---|---|---|
| `GET /espresso/events` | Publish after R01 | Event collection shape; supported Event filters; cursor semantics; `created_at` scope. |
| `GET /espresso/events/{event_id}` | Publish after R02 | Detail-only provenance fields, `links`, and `counts`. |
| `GET /espresso/events/{event_id}/evidence` | Publish after R03 | Current narrow evidence response; direct bidirectional `SAME_AS` only; not article evidence. |
| `GET /espresso/events/{event_id}/signals` | Publish after R04 | Signals derived from the Event and direct `SAME_AS` neighbours; returned records are Signal collections. |
| `GET /espresso/signals` | Publish after R05 | Signal collection shape and supported filters. |
| `GET /espresso/signals/{signal_id}` | Publish after R06 | Detail-only provenance fields, supporting Event link, and count. |
| `GET /espresso/signals/{signal_id}/events` | Publish after R07 | Direct `DERIVED_FROM` Event support; returned records are Event collections. |
| `GET /espresso/sources` and `/{source_id}` | Publish after R08-R09 | Normalized Source fields; Source `q` is text matching, not vector search. |
| `GET /espresso/tags`, `/entities`, `/regions`, `/event-types` | Publish after R11-R14 | Exact stored-value discovery, not canonical entity/place metadata. |
| `GET /espresso/events/count` | Publish after R16 | Count and categorical distributions, not pagination total. |
| `GET /espresso/events/summary` | Publish after R17 | Requires bounded `from` and `to`; approved `group_by` values only; tag/region buckets are non-additive. |
| `POST /espresso/events/search` | Publish after R18 | JSON vector-search body; response is the Event collection envelope. |

Do not document `GET /actions`, Action relations, Stories, article routes, or
generic `/related/{relationship}` routes.

## 6. Query parameter mapping

### 6.1 Common collection and search parameters

| Capability | GDELT Cloud | PredictHQ | Perigon | Espresso | Documentation instruction |
|---|---|---|---|---|---|
| Page size | `limit` | `limit` | `size` | `limit` (`1-100`, default `20`) | Explain that Espresso uses a maximum result count, not a page number. |
| Position | `cursor` | `offset` | `page` | `cursor` | Map GDELT cursor conceptually. Do not offer offset/page aliases or numeric cursor assumptions. |
| Continuation | `next_cursor` | `next` URL | increment `page` | `pagination.next_cursor` | Tell clients to replay the returned opaque cursor with the same filters. |
| Text/semantic search | `search` | `q` | `q`; vector `prompt` body | `q`; planned `POST /events/search` body uses `q` | Describe `q` only where the route supports it. Source `q` is ordinary case-insensitive text matching. |
| Time range | `date_start`, `date_end` | `start.*`, `end.*`, lifecycle ranges | `from`, `to`, ingestion ranges | `from`, `to` | State that Espresso bounds **record creation time**, not occurrence, publication, update, or lifecycle time. V1 accepts date-only values. |
| Ordering | `sort` | `sort` | `sortBy` | No v1 public equivalent | Do not map a provider sort value into an Espresso request. |
| Result total | No standard total | `count` | `numResults` | No collection total | State that `data.length` is the returned-page length. Use `/events/count` only when that route is available. |
| Response format | JSON | JSON | JSON | `response_type=json|yaml|toon` | Call this an Espresso serialization extension, with JSON as the canonical schema. |

### 6.2 Event and Signal filter mapping

| Provider filter family | GDELT Cloud | PredictHQ | Perigon | Espresso mapping | Important limit to document |
|---|---|---|---|---|---|
| Exact resource IDs | `event_id` path | `id` | article/story IDs | `ids` CSV for a collection; `{event_id}` / `{signal_id}` paths for detail | IDs identify Espresso records, not a provider-compatible article or story. |
| Classification | `event_family`, `category`, `subcategory` | `category`, `label`, `phq_label` | `category`, `topic`, `taxonomy`, `label` | `event_types`; `categories` is an alias; `impact_levels`; `tags` | `categories` is not a separate taxonomy. There is no topic, taxonomy, subcategory, or label hierarchy. |
| Entities | entity routes/references | `entity.id`, embedded entities | `people`, `companies`, entity filters | Event `companies`, `people`, `products`, `entities` (company/person union) | Values are exact strings. No entity IDs, aliases, profiles, or generic entity lookup. |
| Geography | country, region, continent, admin1 | country, place, within | country/state/city/coordinates/distance | Event `regions`; discovery through `/regions` | Region strings are not canonical geography and do not support country, coordinates, hierarchy, or radius queries. |
| Source | `domain` | source context in event metadata | source/domain/source group/paywall | `source_ids` for Event/Signal records; Source `domains` and text `q` on `/sources` | Event/Signal filtering does not use provider-style domain, publisher group, paywall, or source geography filters. |
| Impact/quality | fatalities, confidence profile, minimum confidence | rank, local rank, confidence, attendance, spend | sentiment and source metrics | `impact_levels`; optional flattened `confidence` may be returned | `impact_level` is categorical. `confidence` is not a numeric filter or comparable score. |
| Event timing/state | event date | start/end/active/cancelled/state | publication/add/refresh/updated time | none beyond `created_at` bounds | Do not present `created_at` as an event date, article time, or lifecycle state. |
| Relation scope | story/article links | parent-child and embedded metadata | `clusterId`, story/article relationships | Event evidence and Event/Signal relation routes | Espresso relations express record identity evidence and derivation, not narrative clustering. |

### 6.3 Route-specific Espresso parameters

| Route family | Espresso parameters to document |
|---|---|
| Event collections | `ids`, `event_types`, `categories`, `impact_levels`, `companies`, `people`, `products`, `regions`, `entities`, `source_ids`, `tags`, plus applicable common parameters. |
| Signal collections | `ids`, `impact_levels`, `impacted_domains`, `source_ids`, `tags`, plus applicable common parameters. |
| Event evidence | `source_ids`, `from`, `to`, `limit`, `cursor`, `response_type`; no semantic `q` or public sort. |
| Event-to-Signal and Signal-to-Event collections | Filters for the returned resource, plus applicable `q`, time, cursor, and serialization parameters. |
| Sources | `q`, `domains`, `limit`, `cursor`, `response_type`. |
| Discovery | `q`, `limit`, `cursor`, `response_type`, with `resource` for tags and `types` for entities. |
| Summary | Event filters, bounded `from` and `to`, and approved `group_by`; publish the exact list with the endpoint. |
| Vector Event search | JSON `q`, `filters`, `limit`, and `cursor`; no page/offset fields. |

## 7. Response-field mapping

### 7.1 Envelope mapping

| Meaning | GDELT Cloud | PredictHQ | Perigon | Espresso |
|---|---|---|---|---|
| Success indicator | `success` | HTTP status plus payload | `status` | `success` |
| Records | `data` | `results` | `articles` or `results` | `data` |
| Page size | `pagination.limit` | request `limit` | request `size` | `pagination.limit` |
| Continuation | `pagination.next_cursor` | `next` URL | next `page` calculation | `pagination.next_cursor` |
| Total matches | not standard | `count` | `numResults` | omitted from normal collections |
| Snapshot metadata | provider-specific | provider-specific | provider-specific | `meta.as_of` |

### 7.2 Event and Signal item mapping

| Provider field or concept | Espresso field | Mapping and boundary |
|---|---|---|
| GDELT/PredictHQ/Perigon record ID | `id` | Direct record identifier. `kind` distinguishes Espresso Event from Signal. |
| GDELT `summary`; PredictHQ `description`; Perigon `summary`/`description` | `summary` | Normalized from `digest.briefing`; it is optional. Espresso has no stable `title` in the current data. |
| GDELT `url`; Perigon article `url` | `url` | Detail-only provenance field when stored. It is not an article-detail guarantee. |
| Provider source/domain object | `source_id`, `source` | Detail-only mapping from a linked Source. `source` may be omitted when a reference is absent or orphaned. |
| GDELT `event_date`; PredictHQ lifecycle dates; Perigon publication/ingestion dates | `created_at` | Record creation timestamp only. It must not be relabeled as occurrence, publication, refresh, or update time. |
| GDELT category/subcategory; PredictHQ category/labels; Perigon categories/topics | `event_type`, `impact_level`, `tags` | Flattened values when populated. `event_type` is an Espresso classification; tags are persisted tags. |
| Provider people, companies, or entities | `people`, `companies`, `products`, `regions` | Optional flattened exact-string arrays; not normalized entity or place objects. |
| Provider scores and metrics | optional `confidence`; `impact_level` | Do not imply numeric comparability. Espresso does not supply rank, significance, attendance, spend, fatality, or metric methodology. |
| GDELT story/article references; Perigon cluster linkage | `links`, `counts`, relation endpoints | Event detail links to evidence and Signals; Signal detail links to Events. These are relation traversals, not stories or clusters. |
| Provider-specific additional enrichment | other populated digest members | Espresso may expose non-empty digest values such as `drivers`, `impacts`, `forecast`, `activities`, `macro_context`, and `impacted_domains` at the item root. Consumers must treat them as optional. |

### 7.3 Espresso payload examples to include

Use separate examples for a collection item and a detail item. The collection
example must not contain `url`, `base_url`, `source_id`, or `source`.

~~~json
{
  "id": "event-id",
  "kind": "event",
  "created_at": "2026-08-11T05:37:59Z",
  "tags": ["investment_and_capital_markets"],
  "summary": "A normalized briefing.",
  "event_type": "stock_decline",
  "impact_level": "medium",
  "companies": ["example_company"]
}
~~~

The Event detail example adds available provenance and relation navigation:

~~~json
{
  "url": "https://example.com/article",
  "base_url": "https://example.com",
  "source_id": "source-id",
  "source": {
    "id": "source-id",
    "domain": "example.com",
    "name": null,
    "url": "https://example.com"
  },
  "links": {
    "evidence": "/events/event-id/evidence",
    "signals": "/events/event-id/signals"
  },
  "counts": {
    "evidence": 3,
    "signals": 2
  }
}
~~~

## 8. Espresso capabilities that the comparison APIs do not combine

Espresso should distinguish these capabilities without overstating data
coverage:

| Espresso capability | Why it is distinct |
|---|---|
| First-class Signal resource | GDELT, PredictHQ, and Perigon do not expose the same paired Event/Signal resource model in the reviewed routes. |
| `DERIVED_FROM` traversal | `/events/{event_id}/signals` and `/signals/{signal_id}/events` expose explicit intelligence derivation, rather than only embedded entities, article membership, or cluster linkage. |
| Bidirectional `SAME_AS` evidence | Event evidence represents equivalent source records and can expand the Event scope before finding derived Signals. It is not a story cluster. |
| Flexible flattened intelligence digest | Optional fields such as drivers, impacts, forecast, activities, macro context, and impacted domains are emitted alongside a small stable core rather than constrained to one fixed event schema. |
| Unified serialization choice | The same API response can be requested as JSON, YAML, or TOON; the reviewed comparison surfaces are JSON-oriented. |
| Explicit exact-value discovery | Tags, Event types, entity strings, and region strings can be discovered from Espresso's own stored values, with clear non-canonical semantics. |
| Signal-aware semantic retrieval | Semantic `q` can operate for Signal collections and relation collections; planned Event vector search uses a structured JSON body and the same Event envelope. |
| Snapshot metadata | `meta.as_of` lets clients identify the API snapshot time independently of each record's creation timestamp. |

## 9. Known limits to state plainly

- Espresso has no Story or Article resource. `SAME_AS` evidence is not a
  narrative cluster, and `url` does not make an Event a complete article.
- `created_at` is ingestion/record creation time. Espresso has no event
  occurrence, publication, lifecycle, cancellation, refresh, or update time.
- `companies`, `people`, `products`, and `regions` are exact strings, not IDs,
  canonical profiles, aliases, or geographic hierarchy.
- There are no structured place, country, coordinate, radius, language,
  source-group, paywall, topic, taxonomy, sentiment, rank, attendance, spend,
  fatality, or numerical-confidence filters.
- Event semantic-search coverage is incomplete until missing embeddings are
  backfilled. Source references and metadata can also be incomplete.
- Actions are reserved and must not appear in public route lists, accepted
  resource scopes, examples, or response enums.

## 10. Publication gates

Before publishing or revising an Espresso endpoint page:

1. Confirm the route exists and its behavior is covered by service tests.
2. Verify the route annotation, generated Swagger, and `config/espresso.oas.json`
   have identical paths, parameters, schemas, and format enums.
3. Confirm the gateway path uses `/espresso` while the backend route does not.
4. Validate every example against the implemented collection/detail projection.
5. Recheck date semantics, pagination cursor behavior, relationship direction,
   and unsupported filters against the current implementation plan.
6. Mark an endpoint planned rather than published when its data prerequisite,
   query plan, indexes, or response tests remain incomplete.

