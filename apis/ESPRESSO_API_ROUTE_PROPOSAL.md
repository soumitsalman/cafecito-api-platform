# Espresso API Target Contract

Status: target design and implementation instruction  
Updated: 2026-08-10  
Scope: Espresso API only  
Predecessor: [API_CAPABILITY_GAP_ANALYSIS.md](API_CAPABILITY_GAP_ANALYSIS.md)

## 1. Purpose

This document defines the intended public Espresso API, the Go types required
to implement it, and the target query behavior for every route.

It is an implementation instruction, not a claim that every field or route is
already available. Every route is assigned one of these capability states:

| State | Meaning |
|---|---|
| Current-schema | The stored columns or digest fields already exist. Code and query changes may still be required. |
| Query-change | No ingestion schema change is required, but new SQL, joins, cursor handling, or response mapping is required. |
| Action-gated | The route is part of the target surface, but publication waits for the separate Action read contract. This document does not plan Action storage. |
| Upstream-data | The field cannot be produced truthfully until ingestion provides it. |
| Conditional | Implement only after the stated consumer and performance requirements are confirmed. |

This document does not authorize changes to the ingestion schema or the current
Go implementation. It describes the end state an implementation agent should
work toward in stages.

## 2. Product model

Espresso is a structured-intelligence API with four public concepts:

| Public resource | Meaning | Example |
|---|---|---|
| Action | One atomic activity, fact, observation, or measurable change. | A company reduced annual revenue guidance. |
| Event | A record returned by an Event-family route. | A curated development, publisher event extraction, or company post reporting a guidance change. |
| Event evidence | An Event returned through a `SAME_AS` evidence route. | A publisher event extraction reporting the guidance change. |
| Signal | A broader conclusion synthesized from multiple Events, sources, or time periods. | Demand weakness is spreading across the sector. |

The internal storage word `sip` must not appear in public routes or primary
product documentation.

### 2.1 Event-family route rule

Every `/events` route is an Event-family route. Its membership predicate is:

~~~sql
kind LIKE 'event%'
~~~

That includes `event`, `event:news`, `event:blog`, `event:post`, `event:site`,
and `event:social`. `kind` is an internal route-selection predicate only. It is
not a public response field and the API does not derive a public
`representation` field from it.

- `/events` searches all Event-family records.
- `/events/{event_id}` retrieves any Event-family record by UUID.
- `/events/{event_id}/evidence` returns a paginated collection containing the
  requested Event and every direct `SAME_AS` Event. Each item contains only
  `event_id`, `created_at`, `source_id`, `url`, and `base_url`.
- `/signals/{signal_id}/events` returns Event-family records targeted by the
  Signal's `DERIVED_FROM` relations.

Event-family digests may have different fields. The API returns the actual
digest members without requiring them to fit a fixed Event schema.

### 2.2 Relation behavior

`SAME_AS` is bidirectional by definition. The server checks both stored
orientations. Public callers never send `direction`, `from_id`, or `to_id`.

`DERIVED_FROM` is directional internally. The public route name states the
resource the caller wants:

- `/signals/{signal_id}/events` returns Event-family records supporting the Signal.
- `/events/{event_id}/signals` returns Signals derived from the Event or one of
  its equivalent source rows.
- `/events/{event_id}/actions` returns Actions supporting the Event after the
  Action contract is ready.
- `/actions/{action_id}/events` returns Events derived from the Action.

Storage direction is never part of the public API.

## 3. Market-derived design rules

The target contract adapts general patterns from current structured-event
services without copying their domain-specific fields:

- [GDELT Cloud API v2](https://docs.gdeltcloud.com/api-reference/v2): collection
  and detail routes, cursor pagination, evidence previews,
  entity references, and Event summaries.
- [Event Registry](https://newsapi.ai/documentation): Event-to-source drilldown,
  source metadata, filter-value resolution, selectable response information,
  and incremental feeds.
- [PredictHQ Events API](https://docs.predicthq.com/api/events/search-events):
  explicit temporal filters, deterministic ordering, typed entities, update
  filtering, and documented impact ranking.
- [Trading Economics calendar schema](https://docs.tradingeconomics.com/economic_calendar/schema/):
  atomic observations with actual, previous, forecast, revised, unit, reference
  period, source, and update metadata.
- [Perigon Signal Insights](https://perigon.io/docs/api/mcp): evidence that saved
  monitoring definitions have market value, but Perigon Signals are monitors,
  not Espresso's synthesized conclusions.

Do not copy GDELT fatality fields, PredictHQ attendance fields, or the Trading
Economics numeric schema into every Espresso record. Adapt the contract
patterns: stable identity, explicit time, provenance, evidence, ranking
semantics, typed filters, and predictable pagination.

## 4. TARGET ROUTES

Gateway routes include the `/espresso` prefix. Backend Go routes omit it.

### 4.1 Target inventory

| ID | Public route | User question | Priority | Capability |
|---|---|---|---:|---|
| R01 | `GET /espresso/events` | Which Event-family records match my question and filters? | P0 | Current-schema + query-change |
| R02 | `GET /espresso/events/{event_id}` | What Event-family record has this UUID? | P0 | Query-change |
| R03 | `GET /espresso/events/{event_id}/evidence` | Which source-specific records support this Event? | P0 | Query-change |
| R04 | `GET /espresso/events/{event_id}/signals` | Which broader conclusions use this Event as evidence? | P0 | Query-change |
| R05 | `GET /espresso/signals` | Which synthesized conclusions or forecasts match my question? | P0 | Current-schema + query-change |
| R06 | `GET /espresso/signals/{signal_id}` | What is the complete Signal with this UUID? | P0 | Query-change |
| R07 | `GET /espresso/signals/{signal_id}/events` | Which Event-family records support this Signal? | P0 | Query-change |
| R08 | `GET /espresso/sources` | Which source records can I filter by or cite? | P0 | Query-change |
| R09 | `GET /espresso/sources/{source_id}` | What metadata is known about this source? | P0 | Query-change |
| R10 | `GET /espresso/tags` | Which exact tag strings are valid filters? | P0 | Current-schema |
| R11 | `GET /espresso/actions` | Which atomic activities or observations match my filters? | P0 after dependency | Action-gated |
| R12 | `GET /espresso/actions/{action_id}` | What atomic activity does this UUID identify? | P0 after dependency | Action-gated |
| R13 | `GET /espresso/events/{event_id}/actions` | Which Actions compose or support this Event? | P1 after dependency | Action-gated |
| R14 | `GET /espresso/actions/{action_id}/events` | Which Events use this Action? | P1 after dependency | Action-gated |
| R15 | `GET /espresso/events/summary` | How are matching Events distributed across an allowed dimension? | P2 | Conditional |

### 4.2 Routes that sound similar but answer different questions

| Routes | Difference |
|---|---|
| `/events` vs `/events/{id}` | The collection searches all Event-family records. The detail route retrieves exactly one known Event-family record and returns `404` when it does not exist. |
| `/events/{id}/evidence` vs `/signals/{id}/events` | Event evidence follows `SAME_AS` to related Event-family records, usually for source coverage. Signal Events follows `DERIVED_FROM` and returns the Event-family records directly supporting a Signal. |
| `/events/{id}/signals` vs `/signals/{id}/events` | They are inverse user workflows over `DERIVED_FROM`. The first asks for conclusions influenced by an Event; the second asks for Events supporting a conclusion. |
| `/events/{id}/evidence` vs `/events/{id}/actions` | Evidence is reporting or source coverage of the Event. Actions are atomic facts or activities that compose the Event. An article is not an Action. |
| `/actions` vs `/events` | An Action is atomic. An Event combines one or more actions with context and impact. |
| `/events` vs `/signals` | An Event describes a self-contained development. A Signal synthesizes implications across Events, sources, or time. |
| `/events` vs `/events/summary` | `/events` returns individual records. `/events/summary` returns aggregate buckets and never substitutes for Event retrieval. |
| `/tags` vs `/sources` | Tags are scalar filter vocabulary. Sources are UUID-keyed provenance records with metadata. |
| `/events/{id}/evidence` vs retired `/events/{id}/equivalents` | `evidence` answers a customer question and hides `SAME_AS`. `equivalents` exposes an internal graph concept and is not part of the primary target surface. |
| Espresso `/signals` vs Perigon Signals | Espresso Signals are generated conclusions. Perigon Signals are saved monitoring definitions. A future saved-monitoring feature must use `/monitors` or `/saved-queries`, not overload `/signals`. |

### 4.2.1 Response payload separation

| Route shape | Routes | `data` payload | Route-added fields |
|---|---|---|---|
| Event or Signal collection | R01, R04, R05, R07 | Each item is a flattened `Sip.Digest` object. | None. |
| Event or Signal detail | R02, R06 | One flattened `Sip.Digest` object. | `links` and `counts`, parallel to digest fields. |
| Event evidence collection | R03 | Each item is the compact provenance projection `event_id`, `created_at`, `source_id`, `url`, and `base_url`. | None. |
| Source resource | R08, R09 | Separate normalized collection and detail documents. | Only the documented Source fields. |

### 4.3 Common collection parameters

Unless a route overrides a parameter, collections use:

| Parameter | Type | Default | Contract |
|---|---|---:|---|
| `limit` | integer, 1–128 | 20 | Maximum returned records. The query fetches `limit+1` to determine whether a next cursor exists. |
| `cursor` | opaque string | absent | Encodes the last ordering key and UUID. Clients must not construct or inspect it. |
| `response_type` | `json` or `text` | `json` | Text is for MCP/LLM context. It contains the route collection items in deterministic order with a record delimiter. |

Collection success:

```json
{
  "data": [],
  "pagination": {
    "limit": 20,
    "returned_count": 0,
    "next_cursor": null
  },
  "meta": {
    "as_of": "2026-08-03T16:00:00Z"
  }
}
```

Detail success:

```json
{
  "data": {}
}
```

`returned_count` is the number of items in this `data` page. It is not a total
match count. `meta` belongs to collection envelopes only; detail responses do
not contain a `meta` field.

Empty collections return `200` with `data: []`. Missing detail resources return
`404`. Validation failures return `400`. Internal failures return `500` with a
stable error code and no database error text.

No route accepts a `sort` parameter. Event and Signal collections use semantic
distance ordering when `q` is present; otherwise they use stored creation time,
then UUID, descending. R03 uses stored creation time, then Event UUID,
descending. Source and tag collections use the fixed route-specific ordering in
their query contracts.

### 4.3.1 Event and Signal response materialization

R01, R04, R05, and R07 do not have closed public Event or Signal record types.
For every item returned by those routes, the API writes the members of that
row's `Sip.Digest` object directly at the item root. The example name
`any_signal_field` means the actual member name from the digest; it is never a
literal wrapper field.

- `id`, `created_at`, `briefing`, and every other output member are read from
  `Sip.Digest`. The API must not substitute storage-column values for them or
  append storage-column values to the item.
- Flatten every actual digest member at the response-item root. Do not use
  `Extra`, a nested digest field, a fixed field allowlist, or a per-field public
  conversion type.
- If `len(field) == 0`, omit that field. Do not apply another response-field
  omission or rewrite rule.
- Do not add public `kind`, `representation`, or `object` fields. A member
  appears only when it is in the digest, except for the explicit route metadata
  defined for R02 and R06 below.
- R03 is an explicit typed-item exception: each item has only `event_id`,
  `created_at`, `source_id`, `url`, and `base_url`, and the collection contains the requested Event
  plus every Event connected to it by `SAME_AS`.
- R02 adds `links` and `counts` objects alongside the flattened Event digest
  members. `counts` contains `coverage`, `signals`, and `actions`.
- R06 adds `links` and `counts` objects alongside the flattened Signal digest
  members. `counts` contains `events`.
- R02 and R06 write `links` and `counts` inside `data`, parallel to the
  flattened digest members. They are not response-envelope `meta` fields.
- R01, R03, R04, R05, and R07 use the normal collection envelope. R02 and R06
  use the normal detail envelope plus their explicit route metadata.

For example, a digest with `id`, `created_at`, `briefing`, `confidence`, and
`market_context` produces this Event or Signal item:

```json
{
  "id": "event-id",
  "created_at": "2026-07-29T13:00:00Z",
  "briefing": "Example briefing",
  "confidence": 0.92,
  "market_context": { "inventory": "elevated" }
}
```

### 4.4 R01 — Search Event-family records

`GET /espresso/events`

This route returns every record where `kind LIKE 'event%'`. Each result is a
flattened Event digest document as specified in section 4.3.1.

Parameters:

| Parameter | Type | Meaning |
|---|---|---|
| `q` | string, max 1024 | Semantic query. |
| `from` / `to` | RFC 3339 timestamp | Inclusive bounds on the stored creation timestamp. These do not claim Event occurrence time. |
| `ids` | CSV UUIDs, max 128 | Restrict to known Event-family IDs. Prefer the detail route for one ID. |
| `event_types` | CSV strings | Allowlisted match against `digest.event_type`. |
| `impact_levels` | CSV: `low,medium,high` | Match `digest.impact_level`. |
| `companies`, `people`, `products`, `regions` | CSV strings | Match existing digest arrays. Values are not canonical IDs yet. |
| `source_ids` | CSV UUIDs | Match the persisted source UUID on Event-family records. For related Event-family coverage, follow `/events/{id}/evidence`. |
| `tags` | CSV strings | Match persisted tags. |


Response model: a collection of flattened Event digest documents.

Example:

```http
GET /espresso/events?q=semiconductor+demand&impact_levels=high&tags=guidance&limit=2
```

```json
{
  "data": [
    {
      "id": "fe836bc8-d631-4efc-8050-b7b6cf823849",
      "created_at": "2026-07-29T13:00:00Z",
      "briefing": "Example Semiconductor lowered guidance after demand weakened.",
      "event_type": "earnings_guidance",
      "impact_level": "high",
      "actions": ["Annual revenue guidance changed from $10B to $8B."],
      "companies": ["example_semiconductor"],
      "regions": ["asia", "united_states"],
      "market_context": { "inventory": "elevated" }
    }
  ],
  "pagination": { "limit": 2, "returned_count": 1, "next_cursor": null },
  "meta": { "as_of": "2026-08-03T16:00:00Z" }
}
```

`created_at` is returned only when it is a member of the Event digest. The API
does not rename or substitute a storage timestamp into the Event document.

### 4.5 R02 — Retrieve one Event-family record

`GET /espresso/events/{event_id}`

Path parameter: `event_id` must be one RFC 4122 UUID.

Optional parameter: `response_type=json|text`.

Response model: one flattened Event digest document with the R02 route metadata
listed in section 4.5.1.

The detail response contains `url` or source information only when those are
members of the requested digest. Use R03 to inspect other source records that
support it; do not infer that R03 repeats the record returned here.

The route accepts every stored kind satisfying `kind LIKE 'event%'`, including
future Event-family subtypes. That predicate selects rows only; it does not add
fields to the returned digest document.

#### 4.5.1 Event detail metadata

In addition to the flattened digest members, R02 returns `links` and `counts`
objects at the same level as those digest members. `counts.coverage` is the full
direct `SAME_AS` evidence scope, including the requested Event; it is not the R03
page `returned_count`. `counts.signals` counts
Signals derived from the Event's relation scope, and `counts.actions` counts
Actions related to the Event.

```json
{
  "data": {
    "id": "event-id",
    "created_at": "2026-07-29T13:00:00Z",
    "briefing": "Example briefing",
    "links": {
      "evidence": "/events/event-id/evidence",
      "signals": "/events/event-id/signals",
      "actions": "/events/event-id/actions"
    },
    "counts": {
      "coverage": 3,
      "signals": 2,
      "actions": 1
    }
  }
}
```

### 4.6 R03 — Retrieve Event evidence

`GET /espresso/events/{event_id}/evidence`

This is the primary public use of `SAME_AS`. It answers: **which known
publisher or source records are this Event and its related Event records?** The
result includes the requested Event itself and every Event connected to it by
`SAME_AS`.

The API hides relationship direction because `SAME_AS` is bidirectional.

The result contains the requested Event plus every Event connected to it by a
direct `SAME_AS` edge. It does not claim to traverse beyond those direct
neighbours.

Parameters:

| Parameter | Type | Meaning |
|---|---|---|
| `source_ids` | CSV UUIDs | Restrict direct `SAME_AS` neighbours to selected sources. |
| `from` / `to` | RFC 3339 timestamp | Restrict direct `SAME_AS` neighbours by inclusive stored-creation bounds. |

R03 has no `sort` parameter. The server orders its evidence records by stored
creation time and Event UUID, descending.

Response model: a normal paginated collection of `EventEvidenceItem` objects.
Each item contains only `id`, `created_at`, `source_id`, `url`, and
`base_url`. This compact provenance projection is intentionally different from
the flattened Event and Signal digest records returned by other collections.
These fields come from the persisted Event record, not `Sip.Digest`.
`source_id`, `url`, and `base_url` are `null` when their stored values are
unavailable.

The collection envelope is required because a direct `SAME_AS` relation set is
not bounded. It provides the same pagination and freshness contract as the
other Espresso collections without expanding the evidence item payload.

```json
{
  "data": [
    {
      "id": "event-id-for-requested-event",
      "created_at": "2026-07-29T13:00:00Z",
      "source_id": "source-id-for-requested-event",
      "url": "https://example.com/example-semiconductor-guidance",
      "base_url": "https://example.com"
    },
    {
      "id": "event-id-for-related-event",
      "created_at": "2026-07-29T13:08:00Z",
      "source_id": "source-id-for-related-event",
      "url": "https://example.com/example-semiconductor-outlook",
      "base_url": "https://example.com"
    }
  ],
  "pagination": { "limit": 20, "returned_count": 2, "next_cursor": null },
  "meta": { "as_of": "2026-08-03T16:00:00Z" }
}
```

The requested Event remains in the result set even when it has no `SAME_AS`
neighbours or does not match the optional evidence filters.

### 4.7 R04 — Retrieve Signals derived from an Event

`GET /espresso/events/{event_id}/signals`

The server expands the supplied Event-family record's direct `SAME_AS` set and
then finds Signals whose outgoing `DERIVED_FROM` edge targets any member of that
set. The caller does not need to know which stored Event-family row was used as
the relation target.

Parameters: `impact_levels`, `impacted_domains`, `tags`, `from`, `to`, `q`, and
common collection parameters.

Response model: a collection of flattened Signal digest documents.

```json
{
  "data": [
    {
      "id": "signal-id",
      "created_at": "2026-07-30T09:00:00Z",
      "briefing": "Demand weakness is spreading across semiconductor suppliers.",
      "impacted_domains": ["technology", "markets"],
      "forecast": "Supplier earnings revisions are likely over the next quarter."
    }
  ],
  "pagination": { "limit": 20, "returned_count": 1, "next_cursor": null },
  "meta": { "as_of": "2026-08-03T16:00:00Z" }
}
```

The item has no route-added `links` or `counts`; those belong only on R06.

This route is not a generic graph traversal. It returns only `kind=signal`.
The server first verifies that `event_id` identifies an Event-family record. A
missing Event returns `404`; an existing Event with no Signals returns `200`
with an empty collection.

### 4.8 R05 — Search Signals

`GET /espresso/signals`

Parameters:

| Parameter | Type | Meaning |
|---|---|---|
| `q` | string, max 1024 | Semantic query. |
| `from` / `to` | RFC 3339 timestamp | Inclusive bounds on the stored creation timestamp. |
| `ids` | CSV UUIDs, max 128 | Restrict to known Signal IDs. |
| `impact_levels` | CSV: `low,medium,high` | Match Signal impact level. |
| `impacted_domains` | CSV strings | Match existing digest array. |
| `tags` | CSV | Persisted tag filtering. |
| Common collection parameters | — | `limit`, `cursor`, `response_type`. |

Response model: a collection of flattened Signal digest documents.

```json
{
  "data": [
    {
      "id": "5d540490-5ef1-4d61-84fd-1234885b7d97",
      "created_at": "2026-07-30T09:00:00Z",
      "briefing": "Demand weakness is spreading across semiconductor suppliers.",
      "impact_level": "high",
      "drivers": ["export controls", "inventory correction"],
      "impacts": ["lower sector revenue", "margin pressure"],
      "impacted_domains": ["technology", "markets"],
      "forecast": "Supplier earnings revisions are likely over the next quarter."
    }
  ],
  "pagination": { "limit": 20, "returned_count": 1, "next_cursor": null },
  "meta": { "as_of": "2026-08-03T16:00:00Z" }
}
```

Do not describe this route as saved monitoring. A Signal is returned
intelligence, not a monitor configuration.

### 4.9 R06 — Retrieve one Signal

`GET /espresso/signals/{signal_id}`

Path parameter: `signal_id` must be one RFC 4122 UUID.

Optional parameter: `response_type=json|text`.

Response model: one flattened Signal digest document with the R06 route metadata
listed below.

The detail route does not add an `event_refs` field. Use R07 to retrieve
supporting Events.

#### 4.9.1 Signal detail metadata

In addition to the flattened digest members, R06 returns `links` and `counts`
objects at the same level as those digest members. `counts.events` is the
number of Events returned by the Signal's supporting-event relation query.

```json
{
  "data": {
    "id": "signal-id",
    "created_at": "2026-07-30T09:00:00Z",
    "briefing": "Example briefing",
    "links": {
      "events": "/signals/signal-id/events"
    },
    "counts": {
      "events": 4
    }
  }
}
```

### 4.10 R07 — Retrieve Event-family records supporting a Signal

`GET /espresso/signals/{signal_id}/events`

The route follows the Signal's outgoing `DERIVED_FROM` edges and returns target
records whose kind is `event` or starts with `event:`. Use R03 when the caller
specifically wants the `SAME_AS` evidence neighborhood of one Event record.

Parameters: `event_types`, `impact_levels`, `tags`, `from`, `to`, `q`, and
common collection parameters.

Response model: a collection of flattened Event digest documents.

```json
{
  "data": [
    {
      "id": "event-id",
      "created_at": "2026-07-29T13:00:00Z",
      "briefing": "Example Semiconductor lowered guidance after demand weakened.",
      "event_type": "earnings_guidance",
      "market_context": { "inventory": "elevated" }
    }
  ],
  "pagination": { "limit": 20, "returned_count": 1, "next_cursor": null },
  "meta": { "as_of": "2026-08-03T16:00:00Z" }
}
```

The item has no route-added `links` or `counts`; those belong only on R02.

This route returns direct Event-family targets. It does not rewrite source-oriented
records into another document shape.

The server first verifies that `signal_id` identifies a Signal. A missing Signal
returns `404`; an existing Signal with no supporting Events returns `200` with
an empty collection.

### 4.11 R08 — List sources

`GET /espresso/sources`

Parameters:

| Parameter | Type | Meaning |
|---|---|---|
| `q` | string | Case-insensitive match against site name, domain, or base URL. |
| `domains` | CSV strings | Exact domain-name filter. |
| Common collection parameters | — | `limit` and `cursor`. JSON only initially. |

Response model: `SourceCollectionResponse`.

This is the collection-specific Source response. Each collection item has exactly
the normalized discovery fields `id`, `base_url`, `domain_name`, and
`site_name`.

```json
{
  "success": true,
  "data": [
    {
      "id": "1a6023a8-4b26-4b1f-bb2b-ad2d2c260e8e",
      "base_url": "https://example.com",
      "domain_name": "example.com",
      "site_name": "Example Business"
    }
  ],
  "pagination": { "limit": 20, "returned_count": 1, "next_cursor": null },
  "meta": { "as_of": "2026-08-03T16:00:00Z" }
}
```

### 4.12 R09 — Retrieve one source

`GET /espresso/sources/{source_id}`

Path parameter: `source_id` is one UUID.

Response model: `SourceDetailResponse`.

This is the individual-source response. Its document has the same discovery fields plus
`description`, `favicon`, and `rss_feed`; optional values are explicitly
`null` when unavailable.

```json
{
  "success": true,
  "data": {
    "id": "1a6023a8-4b26-4b1f-bb2b-ad2d2c260e8e",
    "base_url": "https://example.com",
    "domain_name": "example.com",
    "site_name": "Example Business",
    "description": null,
    "favicon": null,
    "rss_feed": null
  }
}
```

This is provenance metadata. It does not return Events published by the source.
Use `/events?source_ids={source_id}` for that question.

### 4.13 R10 — Discover tags

`GET /espresso/tags`

Parameters:

| Parameter | Type | Meaning |
|---|---|---|
| `q` | string | Case-insensitive substring or prefix match. |
| `resource` | `event,signal,action,evidence` | Optional kind scope. `action` remains gated. |
| Common collection parameters | — | `limit` and `cursor`. |

Response model: `CollectionResponse[string]`.

Do not return `record_count` until a measured aggregate query and freshness
contract exist. Do not rename this route to `/filter-values/tags`.

### 4.14 R11 and R12 — Search and retrieve Actions

`GET /espresso/actions`  
`GET /espresso/actions/{action_id}`

These are required target routes but remain action-gated.

Target collection parameters: `q`, `from`, `to`, `action_types`,
`subject_types`, `subject_ids`, `tags`, `source_ids`,
and common collection parameters.

Response models:

- `GET /actions` → `CollectionResponse[Action]`
- `GET /actions/{id}` → `DetailResponse[Action]`

The Action payload must be discriminated by `action_type`. Numeric observations
may contain actual, previous, forecast, revised, currency, unit, and reference
period. Appointments, launches, transactions, policy decisions, and lawsuits
must not be forced into the numeric observation shape.

No Action sample in this document should be treated as the final ingestion
schema. The types below define the API adapter boundary only.

### 4.15 R13 and R14 — Event/Action navigation

`GET /espresso/events/{event_id}/actions`  
`GET /espresso/actions/{action_id}/events`

Both routes are action-gated.

- Event Actions answers which atomic facts compose the Event.
- Action Events answers which canonical developments use the Action.
- Neither route accepts direction or relationship parameters.

Parameters: resource-appropriate filters plus common collection parameters.

Response models:

- Event Actions → `CollectionResponse[Action]`
- Action Events → a collection of flattened Event digest documents

### 4.16 R15 — Summarize Events

`GET /espresso/events/summary`

This conditional route returns aggregate buckets, not Events.

Parameters:

| Parameter | Type | Meaning |
|---|---|---|
| Event collection filters | — | Same structured filters as `/events`, excluding semantic `q` in the first version. |
| `group_by` | `created_day,event_type,impact_level,source,tag,region` | Required allowlisted grouping. |
| `from` / `to` | RFC 3339 timestamp | Required bounded range. |
| `limit` | integer, max 100 | Maximum buckets. |

Response model: `EventSummaryResponse`.

```json
{
  "group_by": "created_day",
  "data": [
    {
      "key": "2026-07-29",
      "event_count": 42,
      "coverage_count": 137
    }
  ],
  "meta": {
    "counted_resource": "event",
    "time_field": "created_at",
    "as_of": "2026-08-03T16:00:00Z"
  }
}
```

Implement only after a named consumer, bounded ranges, query-plan review, and
the meaning of each grouping field are documented.

## 5. TARGET TYPES

These blocks describe the intended direction for
`apis/espresso/db/types.go` and `apis/espresso/router/types.go`. They are not
code changes in this document.

### 5.1 Target `apis/espresso/db/types.go`

The database layer must stop discarding persisted `kind`, `source`, `tags`,
`url`, and `base_url`. It should return storage-shaped records and leave public
JSON naming to the router.

```go
package db

import (
    "encoding/json"
    "strings"
    "time"

    "github.com/google/uuid"
)

type SipKind string

const (
    SipKindAction      SipKind = "action"
    SipKindEvent       SipKind = "event"
    SipKindEventBlog   SipKind = "event:blog"
    SipKindEventNews   SipKind = "event:news"
    SipKindEventPost   SipKind = "event:post"
    SipKindEventSite   SipKind = "event:site"
    SipKindEventSocial SipKind = "event:social"
    SipKindSignal      SipKind = "signal"
)

var EventFamilyKinds = []SipKind{
    SipKindEvent,
    SipKindEventBlog,
    SipKindEventNews,
    SipKindEventPost,
    SipKindEventSite,
    SipKindEventSocial,
}

func IsEventKind(kind SipKind) bool {
    return strings.HasPrefix(string(kind), "event")
}

type Sip struct {
    ID       uuid.UUID       `db:"id"`
    Kind     SipKind         `db:"kind"`
    Created  time.Time       `db:"created"`
    SourceID *uuid.UUID      `db:"source"`
    Tags     []string        `db:"tags"`
    Digest   json.RawMessage `db:"digest"`
    URL      *string         `db:"url"`
    BaseURL  *string         `db:"base_url"`
}

type Source struct {
    ID          uuid.UUID `db:"id"`
    BaseURL     string    `db:"base_url"`
    DomainName  *string   `db:"domain_name"`
    SiteName    *string   `db:"site_name"`
    Description *string   `db:"description"`
    Favicon     *string   `db:"favicon"`
    RSSFeed     *string   `db:"rss_feed"`
}

type Relation struct {
    FromID       uuid.UUID `db:"from_id"`
    ToID         uuid.UUID `db:"to_id"`
    Relationship string    `db:"relationship"`
}

type EvidenceRow struct {
    EventID  uuid.UUID  `db:"event_id"`
    Created  time.Time  `db:"created"`
    SourceID *uuid.UUID `db:"source_id"`
    URL      *string    `db:"url"`
    BaseURL  *string    `db:"base_url"`
}

type TagMode string

const (
    TagAny TagMode = "any"
    TagAll TagMode = "all"
)

type PageRequest struct {
    Limit  int
    Cursor *Cursor
}

type Cursor struct {
    Version       int
    ID            uuid.UUID
    Created       *time.Time
    Distance      *float64
    TextKey       *string
}

type SipFilters struct {
    IDs             []uuid.UUID
    Kinds           []SipKind
    CreatedFrom     *time.Time
    CreatedTo       *time.Time
    SourceIDs       []uuid.UUID
    Tags            []string
    TagMode         TagMode
    EventTypes      []string
    ImpactLevels    []string
    Companies       []string
    People          []string
    Products        []string
    Regions         []string
    ImpactedDomains []string
    ActionTypes     []string
    SubjectTypes    []string
    SubjectIDs      []string
    Embedding       []float32
}

type Page[T any] struct {
    Items      []T
    NextCursor *Cursor
}

type SummaryGroup string

const (
    SummaryCreatedDay SummaryGroup = "created_day"
    SummaryEventType  SummaryGroup = "event_type"
    SummaryImpact     SummaryGroup = "impact_level"
    SummarySource     SummaryGroup = "source"
    SummaryTag        SummaryGroup = "tag"
    SummaryRegion     SummaryGroup = "region"
)

type EventSummaryRow struct {
    Key           string `db:"key"`
    EventCount    int64  `db:"event_count"`
    CoverageCount int64  `db:"coverage_count"`
}
```

Type rules:

- Use pointers for nullable persisted columns.
- Keep `Digest` as `json.RawMessage` in the database layer. It is the complete
  source of Event and Signal response fields; do not decode it into a fixed
  schema there.
- Do not flatten or filter digest fields inside the database package.
- Do not expose embedding vectors in public response types.
- Cursor decoding and validation happen before the repository query.
- Query filters are typed and allowlisted. Do not restore a public catch-all
  SQL `Extra` field.

### 5.2 Target `apis/espresso/router/types.go`

The router owns response envelopes, direct Event/Signal digest materialization,
text rendering, normalized Source documents, and error models.

```go
package router

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

type Pagination struct {
    Limit         int     `json:"limit"`
    ReturnedCount int     `json:"returned_count"`
    NextCursor    *string `json:"next_cursor"`
}

type ResponseMeta struct {
    AsOf time.Time `json:"as_of"`
}

type CollectionResponse[T any] struct {
    Data       []T        `json:"data"`
    Pagination Pagination `json:"pagination"`
    Meta       ResponseMeta `json:"meta"`
}

type DetailResponse[T any] struct {
    Data T `json:"data"`
}

type APIError struct {
    Code    string         `json:"code"`
    Message string         `json:"message"`
    Details map[string]any `json:"details,omitempty"`
}

type ErrorResponse struct {
    Error APIError `json:"error"`
}

type SourceReference struct {
    ID      uuid.UUID `json:"id"`
    Name    *string   `json:"name,omitempty"`
    Domain  *string   `json:"domain,omitempty"`
    BaseURL *string   `json:"base_url,omitempty"`
}

type EventEvidenceItem struct {
    EventID  uuid.UUID  `json:"event_id"`
    Created  time.Time  `json:"created_at"`
    SourceID *uuid.UUID `json:"source_id"`
    URL      *string    `json:"url"`
    BaseURL  *string    `json:"base_url"`
}

type SourceCollectionResponse struct {
    Success    bool                   `json:"success"`
    Data       []SourceCollectionItem `json:"data"`
    Pagination Pagination             `json:"pagination"`
    Meta       ResponseMeta           `json:"meta"`
}

type SourceDetailResponse struct {
    Success bool         `json:"success"`
    Data    SourceDetail `json:"data"`
}

type SourceCollectionItem struct {
    ID         uuid.UUID `json:"id"`
    BaseURL    string    `json:"base_url"`
    DomainName *string   `json:"domain_name"`
    SiteName   *string   `json:"site_name"`
}

type SourceDetail struct {
    ID          uuid.UUID `json:"id"`
    BaseURL     string    `json:"base_url"`
    DomainName  *string   `json:"domain_name"`
    SiteName    *string   `json:"site_name"`
    Description *string   `json:"description"`
    Favicon     *string   `json:"favicon"`
    RSSFeed     *string   `json:"rss_feed"`
}

type ActionSubject struct {
    Type string  `json:"type"`
    ID   *string `json:"id,omitempty"`
    Name string  `json:"name"`
}

type ActionObservation struct {
    Metric          string     `json:"metric"`
    Actual          any        `json:"actual,omitempty"`
    Previous        any        `json:"previous,omitempty"`
    Forecast        any        `json:"forecast,omitempty"`
    Revised         any        `json:"revised,omitempty"`
    Currency        *string    `json:"currency,omitempty"`
    Unit            *string    `json:"unit,omitempty"`
    ReferencePeriod *string    `json:"reference_period,omitempty"`
    ReferenceDate   *time.Time `json:"reference_date,omitempty"`
}

type ActionLinks struct {
    Self   string `json:"self"`
    Events string `json:"events"`
}

type Action struct {
    ID          uuid.UUID          `json:"id"`
    Object      string             `json:"object" enums:"action"`
    ActionType  string             `json:"action_type"`
    CreatedAt   time.Time          `json:"created_at"`
    OccurredAt  *time.Time         `json:"occurred_at,omitempty"`
    Subject     *ActionSubject     `json:"subject,omitempty"`
    Observation *ActionObservation `json:"observation,omitempty"`
    Details     json.RawMessage    `json:"details,omitempty"`
    Source      *SourceReference   `json:"source,omitempty"`
    Tags        []string           `json:"tags"`
    Links       ActionLinks        `json:"links"`
}

type EventSummaryBucket struct {
    Key           string `json:"key"`
    EventCount    int64  `json:"event_count"`
    CoverageCount int64  `json:"coverage_count"`
}

type EventSummaryMeta struct {
    CountedResource string    `json:"counted_resource"`
    TimeField       string    `json:"time_field"`
    AsOf            time.Time `json:"as_of"`
}

type EventSummaryResponse struct {
    GroupBy string               `json:"group_by"`
    Data    []EventSummaryBucket `json:"data"`
    Meta    EventSummaryMeta     `json:"meta"`
}
```

Router mapping rules:

1. Use stored `kind` only to select and validate the route's rows. Never write
   it to an Event or Signal response item and never derive `representation`.
2. Require the selected Event or Signal `Digest` to be a JSON object.
3. Write every actual digest member directly into the response item. Do not
   decode to `Event`, `Signal`, `EventEvidence`, an allowlist, or `Extra`.
4. Read `id`, `created_at`, and all other flattened Event/Signal members from
   the digest; R02/R06 relationship metadata is the explicit route exception.
5. Omit a digest member if `len(field) == 0`; otherwise retain it unchanged.
6. R03 evidence items use `EventEvidenceItem`: a collection item with only
   `event_id`, `created_at`, `source_id`, `url`, and `base_url`. It includes the
   requested Event and every direct `SAME_AS` Event.
7. R02 adds the documented `links` and `counts` objects to the flattened Event
   digest item, parallel to its digest members and not inside response `meta`.
8. R06 adds the documented `links` and `counts` objects to the flattened Signal
   digest item, parallel to its digest members and not inside response `meta`.
9. Use `SourceCollectionResponse` for R08 and `SourceDetailResponse` for R09.
   Do not use one Source response type for both routes.
10. Action decoding remains behind an adapter until the external Action contract
   is finalized.
11. Text rendering serializes the same flattened digest members deterministically
   without dropping them.

### 5.3 Upstream-dependent types not yet claimable

Event and Signal routes have no field allowlist: a defined member of a digest is
returned according to section 4.3.1. The remaining gated public type is the
Action contract, including complete Action occurrence, revision, and value
semantics.

## 6. TARGET QUERY

### 6.1 Shared query rules

1. Event queries include every kind satisfying `s.kind LIKE 'event%'`.
2. Event-family evidence queries use the same prefix predicate and may use the
   documented source, time, tag, and cursor filters. They do not have a
   representation filter.
3. Signal queries always include `s.kind = 'signal'`.
4. Action queries always include `s.kind = 'action'` after the dependency gate.
5. Select at least:
   `id, kind, created, source, tags, digest, url, base_url`.
6. `from` and `to` filter `s.created` until another public time field exists.
7. Tag `any` uses `s.tags && @tags`. Tag `all` uses `s.tags @> @tags`.
8. Typed digest filters are allowlisted expressions. Never interpolate a
   client-supplied JSON key or SQL expression.
9. Scalar pagination is keyset pagination. Offset is retained only on deprecated
   compatibility routes.
10. Every collection ordering has a deterministic UUID tie-breaker.
11. Fetch `limit+1` rows, return at most `limit`, and encode the last returned
    ordering key in an opaque versioned cursor.
12. Semantic-query cursor data includes distance and UUID. Default-order cursor
    data includes created timestamp and UUID.
13. Relationship queries omit orphaned targets from the public result and
    increment an internal integrity metric.
14. `SAME_AS` is resolved in both orientations.
15. `DERIVED_FROM` is queried in its stored direction, while route handlers
    present the appropriate forward or inverse user view.
16. Every subresource route verifies its parent first. Missing parent returns
    `404`; an existing parent with no related rows returns `200` and `data: []`.

### 6.2 Route-to-repository contract

| Route | Target repository method | Return |
|---|---|---|
| R01 `GET /events` | `QueryEvents(ctx, filters, page)` | `Page[Sip]` |
| R02 `GET /events/{id}` | `GetEvent(ctx, id)` | `Sip` |
| R03 `GET /events/{id}/evidence` | `QueryEventEvidence(ctx, id, filters, page)` | `Page[EvidenceRow]` |
| R04 `GET /events/{id}/signals` | `QueryEventSignals(ctx, id, filters, page)` | `Page[Sip]` |
| R05 `GET /signals` | `QuerySignals(ctx, filters, page)` | `Page[Sip]` |
| R06 `GET /signals/{id}` | `GetSignal(ctx, id)` | `Sip` |
| R07 `GET /signals/{id}/events` | `QuerySignalEvents(ctx, id, filters, page)` | `Page[Sip]` |
| R08 `GET /sources` | `QuerySources(ctx, query, domains, page)` | `Page[Source]` |
| R09 `GET /sources/{id}` | `GetSource(ctx, id)` | `Source` |
| R10 `GET /tags` | `QueryTags(ctx, query, kinds, page)` | `Page[string]` |
| R11 `GET /actions` | `QueryActions(ctx, filters, page)` | `Page[Sip]` |
| R12 `GET /actions/{id}` | `GetAction(ctx, id)` | `Sip` |
| R13 `GET /events/{id}/actions` | `QueryEventActions(ctx, id, filters, page)` | `Page[Sip]` |
| R14 `GET /actions/{id}/events` | `QueryActionEvents(ctx, id, filters, page)` | `Page[Sip]` |
| R15 `GET /events/summary` | `SummarizeEvents(ctx, filters, group)` | `[]EventSummaryRow` |

### 6.3 R01 query — Event-family collection

Scalar shape:

```sql
SELECT
    s.id, s.kind, s.created, s.source, s.tags, s.digest, s.url, s.base_url
FROM sips AS s
WHERE s.kind LIKE 'event%'
  AND (@from IS NULL OR s.created >= @from)
  AND (@to IS NULL OR s.created <= @to)
  AND (@ids_empty OR s.id = ANY(@ids))
  AND (
      @source_ids_empty
      OR s.source = ANY(@source_ids)
      OR EXISTS (
          SELECT 1
          FROM relations AS source_same
          JOIN sips AS evidence
            ON evidence.id = CASE
                WHEN source_same.from_id = s.id THEN source_same.to_id
                ELSE source_same.from_id
            END
           AND evidence.kind LIKE 'event%'
          WHERE source_same.relationship = 'SAME_AS'
            AND (source_same.from_id = s.id OR source_same.to_id = s.id)
            AND evidence.source = ANY(@source_ids)
      )
  )
  AND (@event_types_empty OR s.digest->>'event_type' = ANY(@event_types))
  AND (@impact_levels_empty OR s.digest->>'impact_level' = ANY(@impact_levels))
  AND (@tags_empty OR <allowlisted any/all tag predicate>)
  AND (@cursor_created IS NULL OR (s.created, s.id) < (@cursor_created, @cursor_id))
ORDER BY s.created DESC, s.id DESC
LIMIT @limit_plus_one;
```

Companies, people, products, and regions use allowlisted JSON-array predicates.
If these filters become frequent, measure and add appropriate expression/GIN
indexes rather than accepting unbounded JSON scans.

Semantic shape:

```sql
SELECT
    s.id, s.kind, s.created, s.source, s.tags, s.digest, s.url, s.base_url,
    s.embedding <=> @embedding AS distance
FROM sips AS s
WHERE s.kind LIKE 'event%'
  AND <same structured filters>
  AND (
      @cursor_distance IS NULL
      OR (s.embedding <=> @embedding, s.id) > (@cursor_distance, @cursor_id)
  )
ORDER BY distance ASC, s.id ASC
LIMIT @limit_plus_one;
```


### 6.4 R02 query — one Event-family record

```sql
SELECT id, kind, created, source, tags, digest, url, base_url
FROM sips
WHERE id = @event_id
  AND kind LIKE 'event%'
LIMIT 1;
```

No row maps to `404 event_not_found`. A database failure maps to `500` and must
not be confused with not-found. The R02 query also obtains `counts.coverage`, `counts.signals`, and
`counts.actions` using the corresponding relationship sets. `counts.coverage`
counts the R03 `evidence_ids` set before evidence filters and pagination, so an
Event with no `SAME_AS` neighbour reports `1`.

### 6.5 R03 query — Event evidence

The result includes the requested record and every direct `SAME_AS` neighbour.
R03 returns only `event_id`, `created_at`, `source_id`, `url`, and `base_url`.

For a one-hop star relation:

```sql
WITH evidence_ids AS (
    SELECT @event_id::uuid AS id
    UNION
    SELECT CASE
        WHEN r.from_id = @event_id THEN r.to_id
        ELSE r.from_id
    END AS id
    FROM relations AS r
    WHERE r.relationship = 'SAME_AS'
      AND (r.from_id = @event_id OR r.to_id = @event_id)
)
SELECT
    s.id AS event_id,
    s.created,
    s.source AS source_id,
    s.url,
    s.base_url
FROM evidence_ids AS e
JOIN sips AS s ON s.id = e.id
WHERE s.kind LIKE 'event%'
  AND (
      s.id = @event_id
      OR (
          (@source_ids_empty OR s.source = ANY(@source_ids))
          AND (@from IS NULL OR s.created >= @from)
          AND (@to IS NULL OR s.created <= @to)
      )
  )
  AND (@cursor_created IS NULL OR (s.created, s.id) < (@cursor_created, @cursor_id))
ORDER BY s.created DESC, s.id DESC
LIMIT @limit_plus_one;
```

The query is intentionally one-hop: it includes the requested Event and the
Events directly connected to it by `SAME_AS`.

### 6.6 R04 query — Signals derived from an Event

```sql
WITH event_scope AS (
    SELECT @event_id::uuid AS id
    UNION
    SELECT CASE
        WHEN r.from_id = @event_id THEN r.to_id
        ELSE r.from_id
    END
    FROM relations AS r
    WHERE r.relationship = 'SAME_AS'
      AND (r.from_id = @event_id OR r.to_id = @event_id)
)
SELECT DISTINCT
    s.id, s.kind, s.created, s.source, s.tags, s.digest, s.url, s.base_url
FROM event_scope AS scope
JOIN relations AS derived
  ON derived.relationship = 'DERIVED_FROM'
 AND derived.to_id = scope.id
JOIN sips AS s
  ON s.id = derived.from_id
 AND s.kind = 'signal'
WHERE <allowlisted Signal filters and cursor>
ORDER BY s.created DESC, s.id DESC
LIMIT @limit_plus_one;
```

The query assumes the established meaning `signal DERIVED_FROM event`. Public
callers do not know or supply that direction.

### 6.7 R05 query — Signal collection

Use the R01 scalar and semantic shapes with:

```sql
WHERE s.kind = 'signal'
  AND (@impact_levels_empty OR s.digest->>'impact_level' = ANY(@impact_levels))
  AND (
      @impacted_domains_empty
      OR (s.digest->'impacted_domains') ?| @impacted_domains
  )
```

All cursor, tag, date, ID, semantic, and deterministic ordering rules remain
the same.

### 6.8 R06 query — one Signal

```sql
SELECT id, kind, created, source, tags, digest, url, base_url
FROM sips
WHERE id = @signal_id
  AND kind = 'signal'
LIMIT 1;
```

The R06 query also obtains `counts.events` from the Signal's `DERIVED_FROM`
targets.

### 6.9 R07 query — Event-family records supporting a Signal

```sql
SELECT DISTINCT
    e.id, e.kind, e.created, e.source, e.tags, e.digest, e.url, e.base_url
FROM relations AS r
JOIN sips AS e
  ON e.id = r.to_id
 AND e.kind LIKE 'event%'
WHERE r.relationship = 'DERIVED_FROM'
  AND r.from_id = @signal_id
  AND <allowlisted Event-family filters and cursor>
ORDER BY e.created DESC, e.id DESC
LIMIT @limit_plus_one;
```

This is deliberately not the same query as R03. R03 follows `SAME_AS` from a
requested Event and returns its related Event-family records. R07 follows the
Signal's direct `DERIVED_FROM` targets and returns Event-family records without
rewriting their digest documents.

### 6.10 R08 query — source collection

```sql
SELECT id, base_url, domain_name, site_name
FROM sources
WHERE (
    @q = ''
    OR site_name ILIKE '%' || @q || '%'
    OR domain_name ILIKE '%' || @q || '%'
    OR base_url ILIKE '%' || @q || '%'
)
  AND (@domains_empty OR domain_name = ANY(@domains))
  AND (
      @cursor_key IS NULL
      OR (COALESCE(site_name, domain_name, base_url), id) > (@cursor_key, @cursor_id)
  )
ORDER BY COALESCE(site_name, domain_name, base_url) ASC, id ASC
LIMIT @limit_plus_one;
```

Measure search performance before adding trigram indexes. Do not hide a slow
unindexed contains search behind undocumented behavior.

### 6.11 R09 query — one source

```sql
SELECT id, base_url, domain_name, site_name, description, favicon, rss_feed
FROM sources
WHERE id = @source_id
LIMIT 1;
```

### 6.12 R10 query — tag discovery

```sql
WITH values AS (
    SELECT DISTINCT unnest(s.tags) AS tag
    FROM sips AS s
    WHERE s.tags IS NOT NULL
      AND (@kinds_empty OR s.kind = ANY(@kinds))
)
SELECT tag
FROM values
WHERE (@q = '' OR tag ILIKE '%' || @q || '%')
  AND (@cursor_tag IS NULL OR tag > @cursor_tag)
ORDER BY tag ASC
LIMIT @limit_plus_one;
```

The response is a string list. Counting every matching record per tag is a
different aggregate query and is intentionally excluded.

### 6.13 R11 query — Action collection

After the Action dependency is ready, use the same collection framework:

```sql
SELECT id, kind, created, source, tags, digest, url, base_url
FROM sips AS s
WHERE s.kind = 'action'
  AND <allowlisted Action filters>
  AND <cursor predicate for deterministic ordering>
ORDER BY deterministic ordering, s.id
LIMIT @limit_plus_one;
```

The adapter must receive the final Action digest variants and filter mapping
from the separate Action workstream. This document does not define storage.

### 6.14 R12 query — one Action

```sql
SELECT id, kind, created, source, tags, digest, url, base_url
FROM sips
WHERE id = @action_id
  AND kind = 'action'
LIMIT 1;
```

### 6.15 R13 query — Actions supporting an Event

The target semantic convention is that an Event is derived from its Actions.
If the Action workstream exposes another physical representation, the database
adapter must translate it without changing the public route.

```sql
SELECT DISTINCT
    a.id, a.kind, a.created, a.source, a.tags, a.digest, a.url, a.base_url
FROM relations AS r
JOIN sips AS a
  ON a.id = r.to_id
 AND a.kind = 'action'
WHERE r.relationship = 'DERIVED_FROM'
  AND r.from_id = @event_id
  AND <allowlisted Action filters and cursor>
ORDER BY a.created DESC, a.id DESC
LIMIT @limit_plus_one;
```

Do not publish this query until the Action relation contract is confirmed by
the separate workstream.

### 6.16 R14 query — Events using an Action

```sql
SELECT DISTINCT
    e.id, e.kind, e.created, e.source, e.tags, e.digest, e.url, e.base_url
FROM relations AS r
JOIN sips AS e
  ON e.id = r.from_id
 AND e.kind LIKE 'event%'
WHERE r.relationship = 'DERIVED_FROM'
  AND r.to_id = @action_id
  AND <allowlisted Event filters and cursor>
ORDER BY e.created DESC, e.id DESC
LIMIT @limit_plus_one;
```

This is the inverse user view of R13, not a bidirectional interpretation of
`DERIVED_FROM`.

### 6.17 R15 query — Event summary

Start with the same Event-family filter CTE used by R01:

```sql
WITH filtered_events AS (
    SELECT s.id, s.created, s.source, s.tags, s.digest
    FROM sips AS s
    WHERE s.kind LIKE 'event%'
      AND s.created >= @from
      AND s.created <= @to
      AND <same allowlisted structured filters as GET /events>
),
coverage_counts AS (
    SELECT
        f.id,
        COUNT(DISTINCT evidence.id) AS coverage_count
    FROM filtered_events AS f
    LEFT JOIN relations AS r
      ON r.relationship = 'SAME_AS'
     AND (r.from_id = f.id OR r.to_id = f.id)
    LEFT JOIN sips AS evidence
      ON evidence.id = CASE
          WHEN r.from_id = f.id THEN r.to_id
          ELSE r.from_id
      END
     AND evidence.kind LIKE 'event%'
    GROUP BY f.id
)
SELECT
    <allowlisted group expression> AS key,
    COUNT(DISTINCT f.id) AS event_count,
    COALESCE(SUM(ec.coverage_count), 0) AS coverage_count
FROM filtered_events AS f
JOIN coverage_counts AS ec ON ec.id = f.id
<allowlisted lateral expansion only for tag or region grouping>
GROUP BY key
ORDER BY event_count DESC, key ASC
LIMIT @limit;
```

Each `group_by` value maps to a hard-coded SQL expression. Never interpolate
the raw client value. Query-plan review is required because the bidirectional
relation join can be expensive.

For `group_by=source`, source means evidence source. One Event may therefore
appear in more than one source bucket, and source-bucket Event counts are not
additive. State that behavior in the route documentation.

## 7. Response and documentation requirements

### 7.1 JSON and text parity

- JSON is the canonical contract.
- Text is a deterministic projection of the route's JSON response.
- R03 text renders the same compact `event_id`, `created_at`, `source_id`,
  `url`, and `base_url` collection items.
- R02 and R06 text include their route metadata counts and links.
- Event/Signal digest fields otherwise render without a fixed field list.
- Arrays retain their order and use an escaped delimiter.
- Records use an explicit delimiter such as `---`.
- Empty text collections return `200` with an empty body only if the OpenAPI
  contract says so; otherwise render the empty collection envelope as JSON.

### 7.2 OpenAPI requirements

For every published route:

1. Give the operation a unique verb-first `operationId`.
2. Begin the summary with the user question it answers.
3. Explain Event-family selection versus Signal selection, and state that
   Event and Signal response items are flattened digest documents with no public
   `kind` or `representation` fields.
4. State exactly which timestamp `from` and `to` filter.
5. State `tags` any/all behavior.
6. State the fixed server ordering and that the route does not accept a `sort`
   parameter.
7. Include executable request and response examples, including the standard
   collection envelope for every collection route.
8. Document `200`, `400`, `404`, `429`, and `500`.
9. Keep backend Swagger and gateway OpenAPI schemas aligned.
10. Do not advertise action routes as available before the dependency gate
    closes.

### 7.3 MCP tool names

Primary tools:

- `search_events`
- `get_event`
- `get_event_evidence`
- `get_event_signals`
- `search_signals`
- `get_signal`
- `get_signal_events`
- `list_intelligence_sources`
- `get_intelligence_source`
- `list_intelligence_tags`
- `search_actions`, `get_action`, `get_event_actions`, and
  `get_action_events` after action readiness
- `summarize_events` only after R15 is published

Do not expose `get_related_sips` in the primary MCP catalog.

## 8. Compatibility and rejected public routes

| Existing or proposed design | Decision | Replacement or reason |
|---|---|---|
| `GET /related/{relationship}` | Deprecate, retain temporarily | Dedicated evidence routes replace it. |
| `GET /events/{id}/equivalents` | Do not include in primary target | Use `/events/{id}/evidence`. A diagnostic route may exist separately for operators. |
| Public `direction`, `from_id`, `to_id` | Reject | Storage orientation is internal. |
| `GET /analytics/events` | Reject | It does not identify a measure. Conditional `/events/summary` has a defined aggregate contract. |
| `GET /analytics/signals` | Reject | No defined user question or measure. |
| `GET /signals/counts` | Reject for now | Counting conclusions does not itself explain intelligence. |
| `GET /filter-values/tags` | Reject | Keep the direct `/tags` resource. |
| `GET /sips` | Reject | Internal persistence abstraction. |
| `GET /events/trending` | Reject until score exists | No trending contract is defined. |
| `GET /events/similar` using `SAME_AS` | Reject | `SAME_AS` is identity equivalence, not semantic similarity. |

## 9. Execution plan

### Stage 0 — Truthful current contract

| Task | Change type | Priority | Acceptance |
|---|---|---:|---|
| Document current bare arrays, offset pagination, `created_at` semantics, empty `204` behavior, and tag matching before changing them. | Documentation/spec | P0 | Docs match deployed behavior. |
| Correct Event and Signal flattened-digest response schemas in backend Swagger and gateway OpenAPI. | Documentation/spec/generated | P0 | Generated and gateway specs agree. |
| Make text rendering deterministic from the flattened digest document. | Code/test/docs | P0 | Golden tests are stable across runs. |
| Add detail routes backed by existing UUID filtering. | Code/query/spec/test/docs | P0 | Correct `404` versus `500` behavior. |

### Stage 1 — Required read capabilities

| Task | Change type | Priority | Acceptance |
|---|---|---:|---|
| Select the stored fields needed for selection and the complete raw digest. | Query/types/test | P0 | Event and Signal response data is not discarded or retyped. |
| Define `/events` as the complete Event-family surface. | Query/test/docs | P0 | Every row satisfying `kind LIKE 'event%'` is eligible; its digest document is returned unchanged except for the length-zero omission rule. |
| Implement bidirectional Event evidence query. | Query/code/test | P0 | Both stored `SAME_AS` orientations return the same evidence. |
| Implement direct Signal→Event-family and equivalence-aware Event→Signal queries. | Query/code/test | P0 | Routes return only their declared resource's flattened digest documents. |
| Add source list/detail and source filtering. | Query/code/test/docs | P0 | Source IDs copy directly into Event filters. |
| Add `to`, typed digest filters, deterministic ordering, and cursor pagination. | Query/code/spec/test/docs | P0 | Every collection has stable, non-duplicating page traversal. |
| Define orphan relation behavior and internal integrity metrics. | Query/test/operations | P0 | Orphans do not corrupt responses and are measurable. |

### Stage 2 — Publish target routes

| Task | Change type | Priority | Acceptance |
|---|---|---:|---|
| Publish R01–R10 with the documented response envelopes and examples. | Code/query/spec/generated/test/docs | P0 | Examples execute against the API. |
| Deprecate generic `/related/{relationship}` and publish a migration table. | Spec/docs | P1 | Existing clients have a typed replacement. |
| Close the Action adapter gate and publish R11–R14 without planning Action storage here. | Adapter/code/query/spec/test/docs | Dependency P0/P1 | Final Action types and relation direction are contract-tested. |
| Measure and, if justified, publish R15. | Query/spec/test/docs | P2 | Named consumer, bounded range, reviewed query plan. |

### Stage 3 — Improve quality and parity

| Task | Change type | Priority | Acceptance |
|---|---|---:|---|
| Add explicit confidence/significance only after upstream methodology exists. | Upstream/API/docs | P1 | Scores have definitions, ranges, and calibration. |
| Normalize entity IDs and discovery routes. | Upstream/query/API/docs | P1 | Entity filter values are stable and unambiguous. |
| Add occurrence/update/lifecycle fields. | Upstream/API/docs | P1 | Every timestamp has one documented meaning. |
| Evaluate polling cursors, streams, or webhooks. | Product/infrastructure | P2 | Concrete freshness consumer and delivery SLO. |
| Monitor latency, result count, DB/embedder time, relation integrity, source completeness, response bytes, and freshness by kind. | Operations | P0 | Dashboards and alerts exist without logging sensitive queries. |

## 10. Definition of done

The target Espresso contract is complete when a new developer or agent can:

1. distinguish Action, Event-family records, Event evidence routes, and Signal;
2. search and retrieve each available resource by UUID;
3. trace an Event to source evidence;
4. trace a Signal to Event-family records and an Event to derived Signals;
5. trace Events to Actions after the Action contract is ready;
6. filter using documented time, tag, source, and digest semantics;
7. paginate deterministically without duplicates or omissions;
8. identify source and URL provenance when the database contains it;
9. use JSON or compact text without losing identity or timestamp precision; and
10. use the API without understanding `SAME_AS` or `DERIVED_FROM` storage
    direction.
