# Espresso API Capability Gap Analysis

Last updated: 2026-08-13T17:25:00-04:00

## References

- [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md): target contract.
- [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md): industry comparison input.

## Method

This compares the current \`apis/espresso/\` implementation with the target
contract. Actual payloads and failures were captured by running the local
Espresso service and querying each proposed route on 2026-08-13. Payload blocks
are compact projections of the live body; omitted large digest arrays and prose
are marked with \`...\`. \`meta.as_of\` is response-generated.

## V1 Scope Clarification

- \`sort\` and RFC 3339 date-time values for query \`from\`/\`to\` are expected
  and acceptable V1 deferrals. Both will be added later.
- Relation-collection \`q\` is an accepted V1 deferral.
- Event and Signal collections omit \`url\`, \`base_url\`, \`source_id\`, and
  \`source\`. This is an accepted V1 payload gap; details expose them when
  usable.
- The scalar JSONB \`?|\` predicates for \`event_type\` and \`impact_level\`
  are valid. They are not a gap and must not be replaced with \`->> = ANY\`.

## System Constraints

- Signals and some Events are internally computed and lack a usable Source.
  Detail responses must not invent one.
- \`/signals\` and \`/events/{id}/signals\` intentionally have no
  \`source_id\` query parameter.

## Status

- **Implementation gap:** contradicts the target contract.
- **V1 accepted deferral:** no V1 work required.
- **Future data-quality work:** year-0001 \`created_at\` means unknown time.
  Collections should exclude those rows; details should return
  \`created_at: null\` and \`data_quality: {"created_at":"unknown"}\`.
- **NONE:** expected and actual are materially the same.

## Route Results

### R01 GET /events

**V1 Readiness:** NOT READY: PARAMS GAP

**Actual requests:** \`GET /events?limit=1\` returned 200. \`GET
/events?entities=company&limit=1\` returned 200. \`GET
/events?categories=binary_detection&limit=1\` did not complete within 30
seconds; the server logged 500 after client cancellation.

**Query gap:** \`categories\` binds to \`digest.categories\`, but the target
defines it as an alias for \`event_types\` (\`digest.event_type\`). \`entities\`
is bound as a persisted-tag filter; it must match exact values in the Event
\`companies\` or \`people\` digest arrays. The live \`entities=company\` result
has \`"company"\` in \`tags\`, demonstrating tag filtering rather than entity
membership. \`sort\` and RFC 3339 dates are V1 deferrals.

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","tags":["tag"],"summary":"...","event_type":"value"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"95376d03-22c1-5b20-a241-04f37c8251a9","kind":"event","created_at":"2026-08-13T02:59:00-04:00","event_type":"benefit_reduction_announcement","impact_level":"high","summary":"On January 1, 2025, Lithuania officially announced a significant reduction in monthly pensions...","tags":null,"drivers":["..."],"impacts":["..."]}],"pagination":{"limit":1,"cursor":null,"next_cursor":"eyJ2IjoxLCJpZCI6Ijk1Mzc2ZDAzLTIyYzEtNWIyMC1hMjQxLTA0ZjM3YzgyNTFhOSIsImMiOiIyMDI2LTA4LTEzVDAyOjU5OjAwLTA0OjAwIn0"},"meta":{"as_of":"2026-08-13T21:02:52.865045605Z"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R02 GET /events/{event_id}

**V1 Readiness:** READY

**Actual request:** \`GET /events/95376d03-22c1-5b20-a241-04f37c8251a9\`
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

**V1 Readiness:** NOT READY: RESPONSE GAP

**Actual request:** \`GET
/events/95376d03-22c1-5b20-a241-04f37c8251a9/evidence?limit=1\` returned 200.

**Query gap:** NONE. The current relation query preserves the requested Event
when applying filters, as required.

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","source_id":"source-id","summary":"..."}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"95376d03-22c1-5b20-a241-04f37c8251a9","kind":"event","created_at":"2026-08-13T02:59:00-04:00","tags":null,"source_id":"9bdbfa7f-74ec-5216-9a8a-047725ae0ad0","url":"https://www.lrt.lt/...","base_url":"www.lrt.lt"}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"2026-08-13T21:03:53.633463587Z"}}
~~~

**Payload gap:** Evidence items omit the flattened digest, including
\`summary\`. The target response is an Event payload, not the current
evidence-only projection.

### R04 GET /events/{event_id}/signals

**V1 Readiness:** NOT READY: PARAMS GAP

**Actual request:** \`GET
/events/95376d03-22c1-5b20-a241-04f37c8251a9/signals?limit=1\` returned 200.

**Query gap:** \`ids\` and \`source_ids\` are not accepted. The query follows
only the requested Event; it must first expand direct SAME_AS Event neighbours
before finding derived Signals. \`q\`, \`sort\`, and RFC 3339 dates are V1
deferrals.

**Expected payload:**
~~~json
{"data":[],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[],"pagination":{"limit":1,"cursor":null,"next_cursor":null}}
~~~

**Payload gap:** NONE for this empty capture. The accepted V1 collection core
omission and future sentinel-time policy apply to nonempty results.

### R05 GET /signals

**V1 Readiness:** READY

**Actual request:** \`GET /signals?limit=1\` returned 200.

**Query gap:** NONE. \`impact_levels\` uses the correct scalar JSONB \`?|\`
predicate. \`sort\` and RFC 3339 dates are V1 deferrals.

**Expected payload:**
~~~json
{"data":[{"id":"signal-id","kind":"signal","created_at":"RFC3339","summary":"...","impact_level":"high"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"cbfdd977-7f75-5f23-8a77-09937686c2b3","kind":"signal","created_at":"2026-08-13T02:57:09-04:00","confidence":"high","impact_level":"high","impacted_domains":["crime","justice","public_health"],"summary":"...","tags":["..."]}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R06 GET /signals/{signal_id}

**V1 Readiness:** NOT READY: RESPONSE GAP

**Actual request:** \`GET
/signals/cbfdd977-7f75-5f23-8a77-09937686c2b3\` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":{"id":"signal-id","kind":"signal","created_at":"RFC3339","url":null,"base_url":null,"source_id":null,"summary":"...","links":{"events":"/signals/signal-id/events"},"counts":{"events":48}}}
~~~

**Actual payload:**
~~~json
{"data":{"id":"cbfdd977-7f75-5f23-8a77-09937686c2b3","kind":"signal","created_at":"2026-08-13T02:57:09-04:00","url":null,"base_url":null,"source_id":null,"source":null,"confidence":"high","impact_level":"high","summary":"...","links":{"events":"/signals/cbfdd977-7f75-5f23-8a77-09937686c2b3/events"},"counts":{"events":48}}}
~~~

**Payload gap:** For a null source reference, the target omits \`source\`; the
actual payload emits \`"source": null\`. All other shown stable fields match.

### R07 GET /signals/{signal_id}/events

**V1 Readiness:** NOT READY: PARAMS GAP

**Actual request:** \`GET
/signals/cbfdd977-7f75-5f23-8a77-09937686c2b3/events?limit=1\` returned 200.

**Query gap:** \`ids\` is not accepted. \`categories\` has the same wrong
\`digest.categories\` binding as R01 rather than being an \`event_types\`
alias; \`entities\` is a tag filter rather than companies/people membership.
\`q\`, \`sort\`, and RFC 3339 dates are V1 deferrals.

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","summary":"...","event_type":"value"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"9985a4d6-4577-5dd2-9f81-868c0f1faee1","kind":"event","created_at":"2026-08-13T02:57:09-04:00","event_type":"...","summary":"...","tags":["..."]}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"}}
~~~

**Payload gap:** NONE beyond the V1-accepted collection core omission and the
future sentinel-time policy.

### R08 GET /sources

**V1 Readiness:** NOT READY: RESPONSE GAP

**Actual request:** \`GET /sources?limit=1\` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"id":"source-id","domain":"example","name":null,"url":"example.com","description":null,"favicon_url":"https://example.com/favicon.ico","rss_feed_url":null}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","name":null,"url":"0.30000000000000004.com","favicon_url":"https://0.30000000000000004.com/favicon.ico"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiayI6IjAuMzAwMDAwMDAwMDAwMDAwMDQuY29tIn0"},"meta":{"as_of":"2026-08-13T21:02:52.619574983Z"}}
~~~

**Payload gap:** \`description\` and \`rss_feed_url\` are omitted when null;
the target requires explicit null.

### R09 GET /sources/{source_id}

**V1 Readiness:** NOT READY: RESPONSE GAP

**Actual request:** \`GET /sources/52860efa-d345-58c4-8037-a9f3ae8987c0\`
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

**Payload gap:** \`description\` and \`rss_feed_url\` are omitted when null;
the target requires explicit null.

### R11 GET /tags

**V1 Readiness:** READY

**Actual request:** \`GET /tags?limit=1\` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"value":"tag"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"\n\nmy_wife_quitter_job_episode_script_analysis"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiayI6IlxuXG5teV93aWZlX3F1aXR0ZXJfam9iX2VwaXNvZGVfc2NyaXB0X2FuYWx5c2lzIn0"}}
~~~

**Payload gap:** NONE.

### R12 GET /entities

**V1 Readiness:** NOT READY: PARAMS GAP | RESPONSE GAP

**Actual requests:** \`GET /entities?limit=1\` returned 200. \`GET
/entities?types=person&limit=1\` returned 400.

**Query gap:** The public target spelling is \`types=person\`; current
validation accepts only \`people\`. The service must accept the public spelling
and return the public type value.

**Expected payload:**
~~~json
{"data":[{"value":"person-name","type":"person"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"007_first_light","type":"company"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwMDdfZmlyc3RfbGlnaHQiLCJ0eXBlIjoiY29tcGFueSJ9fQ"},"meta":{"as_of":"2026-08-13T21:07:18.679967901Z"}}
~~~
The target request returns:
~~~json
{"error":{"code":"invalid_request","message":"Key: 'EntitiesParams.Types[0]' Error:Field validation for 'Types[0]' failed on the 'oneof' tag"}}
~~~

**Payload gap:** The default company payload is the same. For the required
\`types=person\` request, no collection payload is returned; the internal
\`people\` spelling emits \`type: "people"\`, not target \`"person"\`.

### R13 GET /regions

**V1 Readiness:** READY

**Actual request:** \`GET /regions?limit=1\` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"value":"region","type":"region"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"0","type":"region"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwIiwidHlwZSI6InJlZ2lvbiJ9fQ"}}
~~~

**Payload gap:** NONE.

### R14 GET /event-types

**V1 Readiness:** READY

**Actual request:** \`GET /event-types?limit=1\` returned 200.

**Query gap:** NONE.

**Expected payload:**
~~~json
{"data":[{"value":"event_type","type":"event_type"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"data":[{"value":"0_60_mph_acceleration","type":"event_type"}],"pagination":{"limit":1,"cursor":null,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwXzYwX21waF9hY2NlbGVyYXRpb24iLCJ0eXBlIjoiZXZlbnRfdHlwZSJ9fQ"}}
~~~

**Payload gap:** NONE.

### R16 GET /events/count

**V1 Readiness:** NOT READY: PARAMS GAP | RESPONSE GAP

**Actual request:** \`GET /events/count\` returned 400 because
\`/events/:id\` captures \`count\` as an event ID.

**Query gap:** The count route and its Event filters are not registered.

**Expected payload:**
~~~json
{"data":{"count":42,"event_types":{"stock_decline":12},"impact_levels":{"high":17,"medium":25}},"meta":{"time_field":"created_at","as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"error":{"code":"invalid_request","message":"invalid UUID length: 5"}}
~~~

**Payload gap:** The required aggregate payload is absent.

### R17 GET /events/summary

**V1 Readiness:** NOT READY: PARAMS GAP | RESPONSE GAP

**Actual request:** \`GET
/events/summary?from=2026-08-01&to=2026-08-10&group_by=event_type\` returned 400
because \`/events/:id\` captures \`summary\` as an event ID.

**Query gap:** The summary route, required bounded dates, \`group_by\`, and
filters are not registered.

**Expected payload:**
~~~json
{"group_by":"event_type","data":[{"key":"stock_decline","event_count":12}],"meta":{"counted_resource":"event","time_field":"created_at","as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~json
{"error":{"code":"invalid_request","message":"invalid UUID length: 7"}}
~~~

**Payload gap:** The required grouped aggregate payload is absent.

### R18 POST /events/search

**V1 Readiness:** NOT READY: PARAMS GAP | RESPONSE GAP

**Actual request:** \`POST /events/search\` with
\`{"q":"wheat","limit":1}\` returned 404.

**Query gap:** No POST route is registered, so the required JSON request body
and nested \`filters\` object are unsupported.

**Expected payload:**
~~~json
{"data":[{"id":"event-id","kind":"event","created_at":"RFC3339","summary":"..."}],"pagination":{"limit":1,"cursor":null,"next_cursor":"opaque"},"meta":{"as_of":"RFC3339"}}
~~~

**Actual payload:**
~~~text
404 page not found
~~~

**Payload gap:** The required collection envelope and Event payload are absent.

### R22 GET /actions and Action relation routes

**V1 Readiness:** READY (reserved; no published V1 endpoint)

**Actual request:** \`GET /actions\` returned 404.

**Query gap:** NONE. R22 is reserved and publishes no Action parameters.

**Expected payload:** none; reserved route.

**Actual payload:**
~~~text
404 page not found
~~~

**Payload gap:** NONE.

## Verification

1. Re-read the target contract, \`apis/AGENTS.md\`, parameter binding, route
   registration, response mapping, and SQL query implementation in
   \`apis/espresso/\`.
2. Ran the local Espresso service and queried every proposed route on
   2026-08-13; representative live bodies and failures are recorded above.
3. Several discovery/detail requests took about 30 seconds because relation
   counts completed before the response; this is observed runtime behavior, not
   a separate target-contract gap.
4. Documentation-only change; no test suite was run.
