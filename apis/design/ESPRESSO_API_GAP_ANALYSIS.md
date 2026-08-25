# Espresso API Capability Gap Analysis

| Field | Value |
|---|---|
| Status | **historical** |
| Authority | Snapshot versus a **superseded** route proposal (2026-08-13 live probes) |
| Audience | Maintainers reading that snapshot |
| Last verified | 2026-08-25 |
| Owner role | API maintainers (archival) |
| Superseded by | [`config/espresso.oas.json`](../../config/espresso.oas.json) |

Last updated: 2026-08-13T17:53:00-04:00. Published collections use `pagination.num_results` and `meta.as_of`.

## References

- [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md): target contract.
- [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md): industry comparison input.

## Method

This compares the current `apis/espresso/` implementation with the target
contract. Actual payloads and failures were captured by running the local
Espresso service and querying each proposed route on 2026-08-13. Payload blocks
are compact projections of the live body; omitted large optional intelligence arrays and prose
are marked with `...`. `meta.as_of` is response-generated.

## V1 Scope Clarification

- `sort` and RFC 3339 date-time values for query `from`/`to` are expected
  and acceptable V1 deferrals. Both will be added later.
- Relation-collection `q` is an accepted V1 deferral.
- Event and Signal collections omit `url`, `base_url`, `source_id`, and
  `source`. This is an accepted V1 payload gap; details expose them when
  usable.
- `/actions/*`, `/events/search`, `/events/summary`, and
  `/events/count` are out of scope for V1. Their absence is not a V1
  readiness gap.
- The scalar JSONB `?|` predicates for `event_type` and `impact_level`
  are valid. They are not a gap and must not be replaced with `->> = ANY`.

## System Constraints

- Signals cannot be filtered by `source_ids` in V1. Their optional source
  reference is not a supported Signal filter dimension, so `/signals` and
  `/events/{id}/signals` intentionally do not bind that parameter.
- Signals and some Events can lack a usable Source. Detail responses must not
  invent one.

## Status

- **Implementation gap:** contradicts the target contract.
- **V1 accepted deferral:** no V1 work required.
- **Future data-quality work:** year-0001 `created_at` means unknown time.
  Collections should exclude those rows; details should return
  `created_at: null` and `data_quality: {"created_at":"unknown"}`.
- **NONE:** expected and actual are materially the same.

## Route Results

### R01 GET /events

**V1 Readiness:** READY

**Implementation check:** `GET /events?limit=1` is supported. The handler binds
`event_types`, the separate `categories` field, exact snake_case company and
people names through `entities`, and the other structured filters described
below.

**Query behavior:** `categories` is a separate exact category filter from
`event_types`. `entities` accepts arbitrary exact company or people names after
snake_case normalization; it is not a type-label or tag filter. Structured
filters (`event_types`, `categories`, `entities`, `impact_levels`, `companies`,
`people`, `products`, and `regions`) use exact matching. `tags` use fuzzy text
matching against persisted tag labels. `sort` and RFC 3339 dates remain V1
deferrals.

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","tags":["tag"],"summary":"...","event_type":"value"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"95376d03-22c1-5b20-a241-04f37c8251a9","kind":"event","created_at":"2026-08-13T02:59:00-04:00","event_type":"benefit_reduction_announcement","impact_level":"high","summary":"On January 1, 2025, Lithuania officially announced a significant reduction in monthly pensions...","tags":null,"drivers":["..."],"impacts":["..."]}],"pagination":{"limit":1,"page":null,"next_page":"eyJ2IjoxLCJpZCI6Ijk1Mzc2ZDAzLTIyYzEtNWIyMC1hMjQxLTA0ZjM3YzgyNTFhOSIsImMiOiIyMDI2LTA4LTEzVDAyOjU5OjAwLTA0OjAwIn0"},"meta":{"as_of":"2026-08-13T21:02:52.865045605Z"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R02 GET /events/{event_id}

**V1 Readiness:** READY

**Actual request:** `GET /events/95376d03-22c1-5b20-a241-04f37c8251a9`
returned 200 after 30.4 seconds.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":{"id":"event-id","kind":"event","created_at":"RFC3339","url":"https://example.com","base_url":"example.com","source_id":"source-id","source":{"id":"source-id","domain":"example","name":null,"url":"example.com"},"summary":"...","links":{"evidence":"/events/event-id/evidence","signals":"/events/event-id/signals"},"counts":{"evidence":1,"actions":0,"signals":0}}}
~~~

**Actual payload:**
~~~json
{"data":{"id":"95376d03-22c1-5b20-a241-04f37c8251a9","kind":"event","created_at":"2026-08-13T02:59:00-04:00","url":"https://www.lrt.lt/...","base_url":"www.lrt.lt","source_id":"9bdbfa7f-74ec-5216-9a8a-047725ae0ad0","source":{"id":"9bdbfa7f-74ec-5216-9a8a-047725ae0ad0","domain":"lrt","name":null,"url":"www.lrt.lt"},"summary":"On January 1, 2025, Lithuania officially announced a significant reduction in monthly pensions...","links":{"evidence":"/events/95376d03-22c1-5b20-a241-04f37c8251a9/evidence","signals":"/events/95376d03-22c1-5b20-a241-04f37c8251a9/signals"},"counts":{"evidence":1,"actions":0,"signals":0}}}
~~~

**Payload gap:** NONE, except the future sentinel-time policy.

### R03 GET /events/{event_id}/evidence

**V1 Readiness:** READY

**Actual request:** `GET
/events/95376d03-22c1-5b20-a241-04f37c8251a9/evidence?limit=1` returned 200.

**Query gap:** NONE. The current relation query preserves the requested Event
when applying filters, as required.

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","source_id":"source-id","summary":"..."}],"pagination":{"limit":1,"page":null,"next_page":null},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"95376d03-22c1-5b20-a241-04f37c8251a9","kind":"event","created_at":"2026-08-13T02:59:00-04:00","tags":null,"source_id":"9bdbfa7f-74ec-5216-9a8a-047725ae0ad0","url":"https://www.lrt.lt/...","base_url":"www.lrt.lt"}],"pagination":{"limit":1,"page":null,"next_page":null},"meta":{"as_of":"2026-08-13T21:03:53.633463587Z"}}
~~~

**Payload gap:** Evidence items omit the flattened intelligence fields, including
`summary`. The target response is an Event payload, not the current
evidence-only projection. This is acceptable for V1.

### R04 GET /events/{event_id}/signals

**V1 Readiness:** READY

**Actual requests:** `GET
/events/9985a4d6-4577-5dd2-9f81-868c0f1faee1/signals?limit=1` returned the
derived Signal. Adding `source_ids=00000000-0000-0000-0000-000000000000`
returned the same Signal; the parameter is ignored by design.

**Query gap:** `ids` is not accepted. The query follows only the requested
Event; it must first expand direct SAME_AS Event neighbours before finding
derived Signals. `source_ids` is intentionally unsupported for Signals and
is not a gap. `q`, `sort`, and RFC 3339 dates are V1 deferrals. This is acceptable for V1

**Expected payload:**
~~~json
{"data":[{"id":"signal-id","kind":"signal","created_at":"RFC3339","summary":"..."}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"cbfdd977-7f75-5f23-8a77-09937686c2b3","kind":"signal","created_at":"2026-08-13T02:57:09-04:00","summary":"..."}],"pagination":{"limit":1,"page":null,"next_page":null},"meta":{"as_of":"2026-08-13T21:52:46.111423976Z"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R05 GET /signals

**V1 Readiness:** READY

**Actual request:** `GET /signals?limit=1` returned 200.

**Query gap:** NONE. `impact_levels` uses the correct scalar JSONB `?|`
predicate. `sort` and RFC 3339 dates are V1 deferrals.

**Expected payload:**
~~~json
{"data":[{"id":"signal-id","kind":"signal","created_at":"RFC3339","summary":"...","impact_level":"high"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"cbfdd977-7f75-5f23-8a77-09937686c2b3","kind":"signal","created_at":"2026-08-13T02:57:09-04:00","confidence":"high","impact_level":"high","impacted_domains":["crime","justice","public_health"],"summary":"...","tags":["..."]}],"pagination":{"limit":1,"page":null,"next_page":"opaque"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R06 GET /signals/{signal_id}

**V1 Readiness:** READY

**Actual request:** `GET
/signals/cbfdd977-7f75-5f23-8a77-09937686c2b3` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":{"id":"signal-id","kind":"signal","created_at":"RFC3339","url":null,"base_url":null,"source_id":null,"summary":"...","links":{"events":"/signals/signal-id/events"},"counts":{"events":48}}}
~~~

**Actual payload:**
~~~json
{"data":{"id":"cbfdd977-7f75-5f23-8a77-09937686c2b3","kind":"signal","created_at":"2026-08-13T02:57:09-04:00","url":null,"base_url":null,"source_id":null,"confidence":"high","impact_level":"high","summary":"...","links":{"events":"/signals/cbfdd977-7f75-5f23-8a77-09937686c2b3/events"},"counts":{"events":48}}}
~~~

**Payload gap:** NONE

### R07 GET /signals/{signal_id}/events

**V1 Readiness:** READY

**Implementation check:** `GET
/signals/{signal_id}/events?limit=1` is supported and accepts the returned
Event filters, including `ids`, separate `categories`, and exact `entities`.

**Query behavior:** `categories` is a separate exact category filter from
`event_types`, and `entities` accepts arbitrary exact company or people names
after snake_case normalization. Structured Event filters use exact matching;
Event tags use fuzzy text matching against persisted labels. `q`, `sort`, and
RFC 3339 dates are V1 deferrals. This is acceptable for V1

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","summary":"...","event_type":"value"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"9985a4d6-4577-5dd2-9f81-868c0f1faee1","kind":"event","created_at":"2026-08-13T02:57:09-04:00","event_type":"...","summary":"...","tags":["..."]}],"pagination":{"limit":1,"page":null,"next_page":"opaque"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R08 GET /sources

**V1 Readiness:** READY

**Actual request:** `GET /sources?limit=1` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"id":"source-id","domain":"example","name":null,"url":"example.com","description":null,"favicon_url":"https://example.com/favicon.ico","rss_feed_url":null}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","name":null,"url":"0.30000000000000004.com","favicon_url":"https://0.30000000000000004.com/favicon.ico"}],"pagination":{"limit":1,"page":null,"next_page":"eyJ2IjoxLCJpZCI6bnVsbCwiayI6IjAuMzAwMDAwMDAwMDAwMDAwMDQuY29tIn0"},"meta":{"as_of":"2026-08-13T21:02:52.619574983Z"}}
~~~

**Payload gap:** `description` and `rss_feed_url` are omitted when null; the target requires explicit null. This is acceptable for V1

### R09 GET /sources/{source_id}

**V1 Readiness:** READY

**Actual request:** `GET /sources/52860efa-d345-58c4-8037-a9f3ae8987c0`
returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":{"id":"source-id","domain":"example","name":null,"url":"example.com","description":null,"favicon_url":"https://example.com/favicon.ico","rss_feed_url":null}}
~~~

**Actual payload:**
~~~json
{"data":{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","name":null,"url":"0.30000000000000004.com","favicon_url":"https://0.30000000000000004.com/favicon.ico"}}
~~~

**Payload gap:** `description` and `rss_feed_url` are omitted when null; the target requires explicit null. This is acceptable for V1

### R11 GET /tags

**V1 Readiness:** READY

**Actual request:** `GET /tags?limit=1` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"value":"tag"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"\n\nmy_wife_quitter_job_episode_script_analysis"}],"pagination":{"limit":1,"page":null,"next_page":"eyJ2IjoxLCJpZCI6bnVsbCwiayI6IlxuXG5teV93aWZlX3F1aXR0ZXJfam9iX2VwaXNvZGVfc2NyaXB0X2FuYWx5c2lzIn0"}}
~~~

**Payload gap:** NONE.

### R12 GET /entities

**V1 Readiness:** READY

**Actual requests:** `GET /entities?limit=1` returned 200. `GET
/entities?types=people&limit=1` returned 200.

**Query gap:** NONE. The proposal requires `types=company,people`; singular
`person` is not a target request value.

**Expected payload:**
~~~json
{"data":[{"value":"entity-name","type":"people"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"007_first_light","type":"people"}],"pagination":{"limit":1,"page":null,"next_page":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwMDdfZmlyc3RfbGlnaHQiLCJ0eXBlIjoicGVvcGxlIn19"},"meta":{"as_of":"2026-08-13T21:51:33.70590324Z"}}
~~~

**Payload gap:** NONE.

### R13 GET /regions

**V1 Readiness:** READY

**Actual request:** `GET /regions?limit=1` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"value":"region","type":"region"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"0","type":"region"}],"pagination":{"limit":1,"page":null,"next_page":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwIiwidHlwZSI6InJlZ2lvbiJ9fQ"}}
~~~

**Payload gap:** NONE.

### R14 GET /event-types

**V1 Readiness:** READY

**Actual request:** `GET /event-types?limit=1` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"value":"event_type","type":"event_type"}],"pagination":{"limit":1,"page":null,"next_page":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"0_60_mph_acceleration","type":"event_type"}],"pagination":{"limit":1,"page":null,"next_page":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwXzYwX21waF9hY2NlbGVyYXRpb24iLCJ0eXBlIjoiZXZlbnRfdHlwZSJ9fQ"}}
~~~

**Payload gap:** NONE.

### R16 GET /events/count

**V1 Readiness:** READY (out of scope for V1)

**Actual request:** `GET /events/count` returned 400 because
`/events/:id` captures `count` as an event ID.

**Query gap:** NONE for V1. The target explicitly defers this route.

**Expected payload:** none in V1; route is deferred.

**Actual payload:**
~~~json
{"error":{"code":"invalid_request","message":"invalid UUID length: 5"}}
~~~

**Payload gap:** NONE for V1.

### R17 GET /events/summary

**V1 Readiness:** READY (out of scope for V1)

**Actual request:** `GET
/events/summary?from=2026-08-01&to=2026-08-10&group_by=event_type` returned 400
because `/events/:id` captures `summary` as an event ID.

**Query gap:** NONE for V1. The target explicitly defers this route.

**Expected payload:** none in V1; route is deferred.

**Actual payload:**
~~~json
{"error":{"code":"invalid_request","message":"invalid UUID length: 7"}}
~~~

**Payload gap:** NONE for V1.

### R18 POST /events/search

**V1 Readiness:** READY (out of scope for V1)

**Actual request:** `POST /events/search` with
`{"q":"wheat","limit":1}` returned 404.

**Query gap:** NONE for V1. The target explicitly defers this route.

**Expected payload:** none in V1; route is deferred.

**Actual payload:**
~~~text
404 page not found
~~~

**Payload gap:** NONE for V1.

### R22 GET /actions and Action relation routes

**V1 Readiness:** READY (reserved; no published V1 endpoint)

**Actual request:** `GET /actions` returned 404.

**Query gap:** NONE. R22 is reserved and publishes no Action parameters.

**Expected payload:** none; reserved route.

**Actual payload:**
~~~text
404 page not found
~~~

**Payload gap:** NONE.

## Verification

1. Re-read the target contract, `apis/AGENTS.md`, parameter binding, route
   registration, response mapping, and SQL query implementation in
   `apis/espresso/`.
2. Ran the local Espresso service and queried every proposed route on
   2026-08-13; representative live bodies and failures are recorded above.
3. Several discovery/detail requests took about 30 seconds because relation
   counts completed before the response; this is observed runtime behavior, not
   a separate target-contract gap.
4. Documentation-only change; no test suite was run.
