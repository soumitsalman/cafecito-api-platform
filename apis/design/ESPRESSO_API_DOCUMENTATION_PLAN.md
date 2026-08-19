# Espresso API Documentation Plan

Status: Documentation implementation plan

Updated: 2026-08-14

References:

- Target contract: [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md)
- V1 readiness: [ESPRESSO_API_GAP_ANALYSIS.md](ESPRESSO_API_GAP_ANALYSIS.md)
- Industry comparison: [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md)

## 1. Documentation authority and V1 scope

The target contract defines Espresso's public vocabulary. The capability gap
analysis decides whether that contract can be presented as available in V1.
Documentation must not present an implementation gap, accepted deferral, or
future data capability as available.

Gateway examples use /espresso/...; backend Swagger annotations use the
unprefixed route. The public resource names are Event, Signal, and Source.
"sip" is storage terminology. Actions are not a public resource.

| Area | V1 documentation rule |
|---|---|
| Success and errors | Collections return data, pagination, and meta; details return data. Use HTTP status codes, never success or success=false. |
| Pagination | limit is 1-100 with default 20. cursor is the opaque continuation token and pagination.next_cursor identifies the following position. `pagination.num_results` is the number of records in the current page; Espresso does not return a separate total-count query. |
| Serialization | response_type is json, yaml, or toon. JSON is canonical; YAML and TOON are token-optimized serializations for MCP and AI-agent clients. Do not document text. |
| Time | from and to accept ISO date-only input and bound created_at. They are not Event occurrence, article publication, or lifecycle time. RFC 3339 date-time input is deferred. |
| Search and ordering | Document q and score_threshold only on Event and Signal collection routes. Do not document public sort or relation-collection q in V1. |
| Payloads | Event and Signal collections omit url, base_url, source_id, and source. Details expose them only when usable. summary is the public narrative field; do not document briefing. |
| Timestamp quality | Treat year-0001 created_at values as data-quality signals; they are not occurrence or publication timestamps. |
| Deferred routes | Do not publish /actions/*, POST /events/search, GET /events/summary, or GET /events/count in V1. |

V1 restrictions from the gap analysis:

- GET /events documents `event_types` and `categories` as separate exact
  filters, and `entities` as an exact match against Event company or people
  names after snake_case normalization.
- Do not document source_ids for Signal collections or
  GET /events/{event_id}/signals. It is intentionally unsupported in V1.
- Relation collections document ids as CSV UUID filters for the returned
  resource.
- Document R03 as its current narrow evidence projection until it returns
  normal flattened Event items with summary.
- Document R04 as direct derivation until SAME_AS scope expansion is implemented.
- Source description and rss_feed_url may be absent when no value is stored.

## 2. Documentation delivery sequence

### 1. apis/espresso/router/routes.go

  - Replace stale Swaggo descriptions and parameter annotations. They are the
    source for generated backend OpenAPI documents.
  - Use Event and Signal, not "Event-family" or kind LIKE event%. Public kinds
    are event and signal only.
  - Remove text, tag_mode, public sort, offset pagination, fabricated
    reported/site_name fields, and Action references.
  - Normalize Swagger placeholders to event_id, signal_id, and source_id so
    backend Swagger and the gateway match.
  - Document categories as a separate exact category field, entities as exact
    Event company/person names, and relation ids as CSV UUID filters.

  Sample annotation:

  ~~~go
  // @Param from query string false "Inclusive created_at lower date bound (YYYY-MM-DD)."
  // @Param cursor query string false "Opaque cursor token. Follow pagination.next_cursor."
  // @Param response_type query string false "Output serialization." Enums(json, yaml, toon) default(json)
  // @Success 200 {object} EventCollectionResponse
  ~~~

  | Routes | Required annotation content |
  |---|---|
  | GET /events and GET /signals | Date-only created_at bounds, fuzzy tag matching, exact structured filters, cursor pagination, q and score_threshold only if active, and optional flattened intelligence fields. |
  | Event and Signal detail | Detail-only provenance, nullable/omitted Source behavior, links, and counts. Do not call created_at an event date. |
  | GET /events/{event_id}/evidence | Show the directly related records that make up an Event's evidence trail, helping clients assess source coverage. It is not article content or story membership. |
  | GET /events/{event_id}/signals | Return the higher-level Signals that were derived from the Event's evidence trail. |
  | GET /signals/{signal_id}/events | Return the Events that were used to derive the Signal, so clients can inspect the basis of its conclusion. |
  | Sources | Source q is case-insensitive metadata matching, not semantic search. |
  | Discovery | Exact public filter names in normalized snake_case, not canonical entities or geography. Tags are fuzzy text matched; entity types are company and people. |

### 2. apis/espresso/docs/swagger.yaml, apis/espresso/docs/swagger.json, and apis/espresso/docs/docs.go

  - Regenerate all artifacts from the revised annotations; do not manually
    repair generated output.
  - Model distinct schemas for Event/Signal collection items, Event/Signal
    details, Sources, discovery values, relation evidence, pagination,
    num_results, meta, and errors. Intelligence-field schemas allow optional flattened intelligence members
    without exposing a nested internal object, embeddings, or relation direction.
  - Encode envelopes without success. Collections require data, num_results,
    and pagination.cursor/next_cursor; document meta.as_of where emitted. Details use:

  ~~~json
  {"data": {"id": "event-id"}}
  ~~~

  - Keep url, base_url, source_id, and source absent from collection schemas
    and available only in detail schemas. Define summary, not briefing, and
    preserve persisted tags as canonical.
  - Set default 20, minimum 1, and maximum 100 only after implementation
    enforces this public limit.
  - Describe empty collections as HTTP 200 with data: []. Define 400, 401,
    404, 429, and 500 with the standard error body where applicable.

  Regenerate with:

  ~~~bash
  cd apis/espresso
  go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal
  ~~~

### 3. config/espresso.oas.json

  - Mirror every published Swagger operation under /espresso and retain the
    Zuplo rewrite target without that prefix.
  - Copy operation IDs, parameter names, enums, descriptions, success schemas,
    and error schemas from generated Swagger. Remove Event-family and text,
    tag_mode, offset, and success-envelope language.
  - Update shared components before path operations so all routes use the same
    Event, Signal, Source, discovery, pagination, num_results, meta, and error
    definitions.
  - Retain existing gateway security, inbound policies, and rewrites. This plan
    changes documentation only; it does not change authorization or routing.
  - Do not add aggregate, vector, or Action paths. Add a future route's path,
    policy block, rewrite target, operation, and schema with its Swagger work.

  | Backend Swagger path | Gateway OpenAPI path | Zuplo rewrite target |
  |---|---|---|
  | /events | /espresso/events | <ESPRESSO_BASE_URL>/events |
  | /events/{event_id} | /espresso/events/{event_id} | <ESPRESSO_BASE_URL>/events/<event_id> |
  | /events/{event_id}/evidence | /espresso/events/{event_id}/evidence | <ESPRESSO_BASE_URL>/events/<event_id>/evidence |
  | /events/{event_id}/signals | /espresso/events/{event_id}/signals | <ESPRESSO_BASE_URL>/events/<event_id>/signals |
  | /signals and /signals/{signal_id} | prefixed equivalents | matching unprefixed Signal paths |
  | /signals/{signal_id}/events | /espresso/signals/{signal_id}/events | <ESPRESSO_BASE_URL>/signals/<signal_id>/events |
  | /sources and /sources/{source_id} | prefixed equivalents | matching unprefixed Source paths |
  | /tags, /entities, /regions, /event-types | prefixed equivalents | matching unprefixed discovery paths |

### 4. docs/pages/howtos/espresso-howto.mdx

  - Rewrite the page around the V1 public contract rather than the old sip,
    Action, plain-text, offset, and reported model.
  - Explain two resource types first, then discovery, Event search, Event
    detail/evidence, Signals/supporting Events, and Sources. Keep health as a
    short operational check separate from the data API.
  - Add a common parameter table for limit, cursor, response_type, from, to,
    and route-specific q. Explain that clients follow pagination.next_cursor
    rather than use an offset.
  - Replace stale examples with the JSON collection/detail envelope and
    created_at, summary, id, and kind. Only show a second request when
    next_cursor is non-null.
  - Include one short YAML or TOON response after canonical JSON. Explain that
    YAML and TOON are token-optimized output formats for MCP and AI agents; do
    not present response_type=text as an MCP optimization.
  - Add a limitations callout: no Actions, Stories, articles, canonical entity
    profiles, structured places, lifecycle dates, numeric metrics, or complete
    Event embedding coverage.

  Required user scenarios:

  1. Discover tags, Event types, entities, and regions before filtering.
  2. Find recent Events by tags and date-only created_at; paginate with
     next_cursor.
  3. Retrieve an Event for detail-only Source provenance, then follow evidence
     or Signal links.
  4. Search Signals and inspect direct supporting Events.
  5. Find a Source and use its ID as an Event source_ids filter.

  Sample V1 collection example:

  ~~~json
  {
    "data": [{
      "id": "event-id",
      "kind": "event",
      "created_at": "2026-08-13T02:59:00Z",
      "tags": ["policy"],
      "summary": "A normalized briefing.",
      "event_type": "policy_reform"
    }],
    "pagination": {
      "limit": 20,
      "num_results": 20,
      "cursor": "provided-string or null",
      "next_cursor": "opaque internally calculated string or null"
    },
    "meta": {"as_of": "2026-08-14T00:00:00Z"}
  }
  ~~~

### 5. Cross-surface validation

  - Regenerate Swagger after annotation changes, then compare operations and
    schemas with config/espresso.oas.json.
  - Validate every how-to request against its gateway path and generated OpenAPI
    schema; do not retain legacy endpoint names or fields.
  - Check examples use date-only V1 inputs, cursor/next_cursor pagination,
    num_results, canonical JSON envelopes, and fields returned by the current
    collection/detail shape.
  - Re-check the capability analysis before publishing a target-complete route.
    A deferred route stays absent from reference and how-to.

## 3. Endpoint publication matrix

| ID | Endpoint | V1 documentation action | Three-surface requirement |
|---|---|---|---|
| R01 | GET /events | Publish separate exact `event_types` and `categories` filters, plus exact Event company/person names for `entities`. | Keep parameters identical in annotations, Swagger, Zuplo, and examples. |
| R02 | GET /events/{event_id} | Publish. | Detail schema/how-to show optional provenance, links, counts, and summary. |
| R03 | GET /events/{event_id}/evidence | Publish an Event evidence trail. | Use a collection envelope, never a bare list or text; document the current V1 item projection. |
| R04 | GET /events/{event_id}/signals | Publish Signals supported by an Event's evidence trail. | Document returned Signal filters, including ids; do not list Signal source_ids or relation q in V1. |
| R05 | GET /signals | Publish. | Describe Signal filters and shared collection envelope. |
| R06 | GET /signals/{signal_id} | Publish. | Show nullable/omitted Source provenance and Event link/count. |
| R07 | GET /signals/{signal_id}/events | Publish the Events supporting a Signal. | Document returned Event filters, including ids, separate categories, and entities. |
| R08-R09 | Source list/detail | Publish. | Use domain, name, url, description, favicon_url, and rss_feed_url; describe V1 omission of unavailable values. |
| R11-R14 | Discovery routes | Publish. | Public filter vocabulary: exact snake_case names for event types, entities, and regions; fuzzy-match labels for tags; types=company,people; update limits after implementation. |
| R16-R18 | Count, summary, Event vector search | Do not publish in V1. | No Swagger annotations, gateway paths, or how-to examples until implemented. |
| R22 | Actions and Action relations | Do not publish. | Exclude from schemas, route inventory, parameter scopes, and examples. |

## 4. Industry query-parameter mapping

This comparison supports migration guidance only. It does not make Espresso a
drop-in replacement for GDELT Cloud, PredictHQ, or Perigon.

| Capability | GDELT Cloud | PredictHQ | Perigon | Espresso V1 | Documentation message |
|---|---|---|---|---|---|
| Page size | limit | limit | size | limit | Espresso returns a numbered page. |
| Position | cursor | offset | page | cursor | Follow the returned Espresso cursor. |
| Continuation | next_cursor | next URL | next page | pagination.next_cursor | Follow next_cursor while it is present. |
| Text/semantic query | search | q | q; vector prompt body | q where implemented | Source q is text matching. Event vector POST is deferred. |
| Date bounds | date_start/date_end | start/end/lifecycle | from/to and add/refresh | from/to | Bounds record creation only and accept a date in V1. |
| Classification | family/category/subcategory | category/label | category/topic/taxonomy/label | event_types, categories, impact_levels, tags | Event types and categories are separate exact snake_case filters; tags use fuzzy text matching. |
| Entity filter | entity references | entity.id | people/company identifiers | Event companies, people, products, entities | Exact snake_case names, not entity IDs or profiles. |
| Geography | country/region/admin1 | country/place/within | country/city/coordinates | Event regions | Exact normalized snake_case region name; no structured geography. |
| Source filter | domain | source context | source/domain/group/paywall | Event source_ids; Source q/domains | No domain, group, paywall, or source-geography Event filter. |
| Ordering | sort | sort | sortBy | no V1 equivalent | Do not translate a provider sort choice. |
| Page result count | absent | count (total) | numResults (total) | num_results | `num_results` is the current page count; no separate total is returned. |

## 5. Response-field mapping

### Envelope mapping

| Meaning | GDELT Cloud | PredictHQ | Perigon | Espresso |
|---|---|---|---|---|
| Success indication | success | HTTP status/payload | status | HTTP 2xx; no body success field |
| Records | data | results | articles or results | data |
| Continuation | pagination.next_cursor | next URL | page calculation | pagination.next_cursor |
| Page result count | not standard | count (total) | numResults (total) | num_results |
| Retrieval metadata | provider-specific | provider-specific | provider-specific | meta.as_of when emitted |

### Record mapping

| Provider field or concept | Espresso field | Boundary to state |
|---|---|---|
| Record ID | id | kind identifies an Espresso Event or Signal. |
| GDELT summary; PredictHQ description; Perigon summary/description | summary | Optional normalized narrative. Espresso has no stable title. |
| GDELT/Perigon URL | url | Detail-only provenance when stored; not an article-detail payload. |
| Provider source object/domain | source_id, source | Detail-only when a usable Source exists. Do not invent a Source for computed Signals or orphaned references. |
| Event, article, or lifecycle timestamp | created_at | Record creation only; not occurrence, publication, refresh, cancellation, or update time. |
| Categories, labels, topics | event_type, categories, impact_level, tags | Optional flattened values. Event types and categories are separate exact fields; tags use fuzzy text matching; no topic/taxonomy hierarchy. |
| People, companies, entities, locations | people, companies, products, regions | Optional exact-string arrays, not normalized objects. |
| Metrics, rank, confidence, attendance, fatalities, spend | impact_level; optional confidence | impact_level is categorical and confidence is non-numeric. |
| Story/article or cluster links | links, counts, relation routes | Relations are not stories or article clusters. |
| Provider-specific enrichment | optional intelligence members | drivers, impacts, forecast, activities, macro_context, and impacted_domains are optional. |

## 6. Espresso-specific value to explain

| Capability | Documentation position |
|---|---|
| First-class Signals | Signals are synthesized intelligence separate from Events, not saved monitors or generic news stories. |
| DERIVED_FROM traversal | Explain direct Event/Signal drilldown and stored direction. |
| SAME_AS evidence | Explain direct bidirectional source-record evidence, not a story cluster or full transitive graph. |
| Flattened intelligence payload | Small stable core plus optional root-level intelligence fields; no nested internal object. |
| Exact-value discovery | Event types, categories, company/people entity names, and regions are discoverable exact snake_case filter values; tags are discoverable labels used with fuzzy matching. |
| JSON, YAML, and TOON | JSON is canonical; YAML and TOON are token-optimized output formats for MCP and AI agents. |
| meta.as_of | It identifies the response snapshot and is distinct from created_at. |

## 7. Known limitations and publication gates

Every route page and the how-to limitations callout must state:

- No Actions, Stories, articles, stable titles, article bodies, or citation roles.
- No occurrence/lifecycle time, structured places, country codes, coordinates,
  canonical entities, aliases, topic taxonomy, or numeric industry metrics.
- Event embeddings and Source references are incomplete; provenance can be
  absent and semantic Event coverage is partial.

Before publishing or materially revising a route, verify:

1. The handler accepts every documented parameter and returns the documented
   JSON, YAML, and TOON shape.
2. Swaggo annotations, generated apis/espresso/docs/ artifacts,
   config/espresso.oas.json, and the how-to have identical route names,
   enums, defaults, limits, envelopes, and field names.
3. The gateway path has /espresso; its Zuplo rewrite target does not.
4. Each example uses a real V1 field, cursor/next_cursor pagination, date-only bounds, and the
   route's collection/detail projection.
5. The capability analysis marks behavior ready, or documentation labels it
   planned rather than available.
