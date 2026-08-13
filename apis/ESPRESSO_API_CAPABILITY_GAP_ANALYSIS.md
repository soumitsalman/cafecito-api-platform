# Espresso API Capability Gap Analysis

Last updated: 2026-08-13T09:52:53-04:00
Live capture window: 2026-08-12T21:51:04-04:00 through 2026-08-12T21:54:44-04:00.

## References

- Target: [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md)
- Industry reference: [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md)
- Follow strictly [AGENTS.md](AGENTS.md)

## Method

This is a historical live-run assessment. Espresso was started locally on
http://127.0.0.1:18081 using apis/espresso/.env, and every target path below
was requested. The expected payload is the target-contract projection of the
same live record; the actual payload is a field-preserving excerpt of the HTTP
body received. The source tree has changed since this capture, so current-code findings override any conflicting historical payload excerpt. Arbitrary digest arrays and long briefing text are elided only
where they do not affect the gap. All stable IDs, timestamps, status codes, and shown response fields are from the live capture; as_of is response-generated and therefore shown as RFC3339 in target projections.

The target proposal has a contradiction about kind: Section 3.3 requires it in
the stable resource core and Section 6 says to hide it. Expected payloads
follow Section 3.3. This is a proposal decision, not an Espresso gap.

All collection excerpts show a shared envelope difference: target pagination
includes cursor and uses default limit 20; live Espresso returns no cursor and
the binder default is limit 16. Individual live requests intentionally used
limit=1.

## V1 Scope Clarification

- Routes lack `sort` and RFC 3339 date-time values for `from`/`to` query params
- `/events/:id/signals`, `/signals/:id/events`, `/events/:id/evidence` lack `q` query param
- Events and Signals collection response payload omits url, base_url, source_id and source fields. They will exist ONLY in the details payload such as `/events/:id` and `/signals/:id`
- `/events/:id/signals`, `/signals/:id/events` will NOT expand search scope using `SAME_AS` of the anchor event id

These will all be considered for future extension

## System Constraints
- Signals and some events are computed internally; Lacks source info
- Signals search (`/signals`, `/events/:id/signals`) lack `source_id` query param. 

## Finding Status

- **V1 accepted deferral:** no implementation work is required for V1. This applies to sort, RFC 3339 date-time input, relation q, and Event/Signal collection URL and Source fields.
- **Implementation gap:** behavior contradicts the published V1 contract and needs a code or contract change.
- **Future data-quality work:** only applies to the known year-0001 created_at sentinel. The required decision is to exclude those records from collections and return null plus data_quality.created_at=unknown on details, as specified by the proposal.
- **No gap:** the scalar JSONB filters for event_type and impact_level are settled, valid filters. Do not replace them with ->> / = ANY. The contradictory proposal text must be reconciled separately.

## Route Results

### R01 GET /events

**Live request:** GET /events?event_types=binary_detection&limit=1 returned
200 and Event b95db348-0de3-5759-86c5-843deb2aaf63.

**Query gap:** categories is the event_types compatibility alias and entities is a convenience filter across companies and people. Neither binds today, so callers must send event_types and separate companies/people filters instead. sort and RFC 3339 from/to are V1-accepted deferrals; the live 400 is expected V1 behavior. acc and impacted_domains are accepted implementation extensions.

**Expected payload:**
~~~json
{"data":[{"id":"b95db348-0de3-5759-86c5-843deb2aaf63","kind":"event","created_at":"2026-08-09T13:39:55-04:00","tags":["paris","earth_sciences_and_natural_resources"],"summary":"On August 9, 2026, astronomers detected evidence supporting Betelgeuse...","briefing":"On August 9, 2026, astronomers detected evidence supporting Betelgeuse...","event_type":"binary_detection","companies":["observatoire_de_paris","american_association_of_variable_star_observers"],"regions":["paris","chile"]}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"id":"b95db348-0de3-5759-86c5-843deb2aaf63","kind":"event","created_at":"2026-08-09T13:39:55-04:00","briefing":"On August 9, 2026, astronomers detected evidence supporting Betelgeuse...","event_type":"binary_detection","companies":["observatoire_de_paris","american_association_of_variable_star_observers"],"regions":["paris","chile"],"tags":["paris","earth_sciences_and_natural_resources","stellar_evolution","satellite_systems_and_space_operations","observatoire_de_paris","chile","american_association_of_variable_star_observers"]}],"pagination":{"limit":1,"next_cursor":null},"meta":{"as_of":"2026-08-13T01:53:24.973088031Z"}}
~~~
**Payload gap:** Collection URL and Source omission is V1-accepted. The only data-quality follow-up is the year-0001 sentinel policy; currently a sentinel would be returned as a literal timestamp rather than excluded or marked. Tag overlap and scalar JSONB filters are valid current behavior.

### R02 GET /events/{event_id}

**Live request:** GET /events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3 returned 200.

**Query gap:** NONE. Both expect and accept only response_type.

**Expected payload:**
~~~json
{"data":{"id":"6fe7144a-aa51-5fa1-a2dd-dd906c583fe3","kind":"event","created_at":"2026-08-10T01:07:00-04:00","url":"https://www.thecooldown.com/sustainable-food/us-home-bakers-stock-up-flour-prices/","base_url":"www.thecooldown.com","source_id":null,"tags":["and_reuse_food","home"],"summary":"On August 10, 2026, the USDA revised its wheat supply outlooks...","briefing":"On August 10, 2026, the USDA revised its wheat supply outlooks...","links":{"evidence":"/events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3/evidence","signals":"/events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3/signals"},"counts":{"evidence":1,"signals":0}}}
~~~
**Actual payload:**
~~~json
{"data":{"id":"6fe7144a-aa51-5fa1-a2dd-dd906c583fe3","kind":"event","created_at":"2026-08-10T01:07:00-04:00","briefing":"On August 10, 2026, the USDA revised its wheat supply outlooks...","forecast":"Supply decline continues into next year","tags":["and_reuse_food","home","change_the_way_you_buy"],"source":{"url":"www.thecooldown.com"},"links":{"evidence":"/events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3/evidence","signals":"/events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3/signals"},"counts":{"coverage":1}}}
~~~
**Payload gap:** The displayed body is historical. Current code now projects url, base_url, source_id, and summary on details. coverage replaces evidence. GetSip inner-joins sources, so a null or orphan Source is falsely 404.

### R03 GET /events/{event_id}/evidence

**Live request:** GET /events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3/evidence?limit=1
returned 200.

**Query gap:** source_ids is intended to restrict only SAME_AS neighbours. The current SQL applies the common source predicate after it combines the anchor and neighbours, so an anchor from a different source can be removed. Fix the query to preserve the requested id unconditionally and apply source_ids only to the neighbour branch. Date-only from/to is a V1-accepted deferral.

**Expected payload:**
~~~json
{"data":[{"id":"6fe7144a-aa51-5fa1-a2dd-dd906c583fe3","kind":"event","created_at":"2026-08-10T01:07:00-04:00","url":"https://www.thecooldown.com/sustainable-food/us-home-bakers-stock-up-flour-prices/","base_url":"www.thecooldown.com","source_id":null,"tags":["and_reuse_food","home"],"briefing":"On August 10, 2026, the USDA revised its wheat supply outlooks..."}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"event_id":"6fe7144a-aa51-5fa1-a2dd-dd906c583fe3","created":"2026-08-10T01:07:00-04:00","url":"https://www.thecooldown.com/sustainable-food/us-home-bakers-stock-up-flour-prices/","base_url":"www.thecooldown.com"}],"pagination":{"limit":1,"next_cursor":null},"meta":{"as_of":"2026-08-13T01:53:07.690795687Z"}}
~~~
**Payload gap:** Live evidence uses event_id and created and omits kind, tags, digest fields, and Source. The specific future data-quality work is year-0001 handling; no generic quality object is required for ordinary evidence rows.

### R04 GET /events/{event_id}/signals

**Live request:** GET /events/6fe7144a-aa51-5fa1-a2dd-dd906c583fe3/signals?limit=1
returned 200 with an empty collection.

**Query gap:** ids and source_ids are missing target filters. q and sort are V1-accepted deferrals, and date-only from/to is expected V1 behavior. The remaining functional gap is relation scope: query direct SAME_AS neighbours before finding DERIVED_FROM Signals.

**Expected payload:**
~~~json
{"data":[],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[],"pagination":{"limit":1,"next_cursor":null},"meta":{"as_of":"2026-08-13T01:53:07.727441715Z"}}
~~~
**Payload gap:** Collection URL and Source omission is V1-accepted. The remaining future data-quality work is to apply the year-0001 sentinel policy consistently when this route returns a Signal.

### R05 GET /signals

**Live request:** GET /signals?limit=1 returned 200 and Signal
ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6.

**Query gap:** NONE for V1. sort and RFC 3339 from/to are accepted deferrals; acc is an implementation extension. Tags use persisted-array overlap, and impact_levels uses the settled valid scalar JSONB filter.

**Expected payload:**
~~~json
{"data":[{"id":"ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6","kind":"signal","created_at":"2026-08-09T13:39:55-04:00","tags":["chile","portugal"],"summary":"On July 29, 2026, ESA/VLT achieved six sigma visibility breakthrough...","briefing":"On July 29, 2026, ESA/VLT achieved six sigma visibility breakthrough...","impact_level":"high","impacted_domains":["astronomy","astrobiology","theoretical_model_validation"]}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"id":"ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6","kind":"signal","created_at":"2026-08-09T13:39:55-04:00","briefing":"On July 29, 2026, ESA/VLT achieved six sigma visibility breakthrough...","confidence":"high","impact_level":"high","impacted_domains":["astronomy","astrobiology","theoretical_model_validation"],"tags":["chile","portugal","paris","la_silla","eso"]}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6ImVmOGUwMzU4LTJiZjQtNTRlNy05MWQ1LWMyYTFhMTUwZjZmNiIsImMiOiIyMDI2LTA4LTA5VDEzOjM5OjU1LTA0OjAwIn0"},"meta":{"as_of":"2026-08-13T01:51:04.08525812Z"}}
~~~
**Payload gap:** Collection URL and Source omission is V1-accepted. The only future quality work is year-0001 sentinel handling. pagination.cursor remains a V1 contract gap.

### R06 GET /signals/{signal_id}

**Live request:** GET /signals/ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6 returned 200.

**Query gap:** NONE. Both expect and accept only response_type.

**Expected payload:**
~~~json
{"data":{"id":"ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6","kind":"signal","created_at":"2026-08-09T13:39:55-04:00","source_id":"ee2e48b1-98b1-5019-895c-1c56f8eb9db2","source":{"id":"ee2e48b1-98b1-5019-895c-1c56f8eb9db2"},"tags":["chile","portugal"],"summary":"On July 29, 2026, ESA/VLT achieved six sigma visibility breakthrough...","briefing":"On July 29, 2026, ESA/VLT achieved six sigma visibility breakthrough...","links":{"events":"/signals/ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6/events"},"counts":{"events":7}}}
~~~
**Actual payload:**
~~~json
{"data":{"id":"ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6","kind":"signal","created_at":"2026-08-09T13:39:55-04:00","briefing":"On July 29, 2026, ESA/VLT achieved six sigma visibility breakthrough...","confidence":"high","impact_level":"high","source":{"id":"ee2e48b1-98b1-5019-895c-1c56f8eb9db2"},"tags":["chile","portugal","paris","la_silla","eso"],"links":{"events":"/signals/ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6/events"},"counts":{"events":7}}}
~~~
**Payload gap:** The displayed body is historical. Current code projects url, base_url, source_id, and summary on details, but deletes briefing; the target requires both summary and briefing. The remaining data-quality work is year-0001 sentinel handling. The inner Source join can falsely 404 an orphan reference.

### R07 GET /signals/{signal_id}/events

**Live request:** GET /signals/ef8e0358-2bf4-54e7-91d5-c2a1a150f6f6/events?limit=1
returned 200.

**Query gap:** Direct DERIVED_FROM traversal works. ids, categories, and entities do not bind. companies, people, products, regions, and source_ids bind but SignalEventsParams.createFilters drops them before SQL; add them to db.Filters there. q and sort are V1-accepted deferrals; date-only from/to is expected V1 behavior. impacted_domains is a non-target Event filter.

**Expected payload:**
~~~json
{"data":[{"id":"b95db348-0de3-5759-86c5-843deb2aaf63","kind":"event","created_at":"2026-08-09T13:39:55-04:00","tags":["paris","earth_sciences_and_natural_resources"],"summary":"On August 9, 2026, astronomers detected evidence supporting Betelgeuse...","briefing":"On August 9, 2026, astronomers detected evidence supporting Betelgeuse...","event_type":"binary_detection"}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"id":"b95db348-0de3-5759-86c5-843deb2aaf63","kind":"event","created_at":"2026-08-09T13:39:55-04:00","briefing":"On August 9, 2026, astronomers detected evidence supporting Betelgeuse...","event_type":"binary_detection","companies":["observatoire_de_paris","american_association_of_variable_star_observers"],"regions":["paris","chile"],"tags":["paris","earth_sciences_and_natural_resources","stellar_evolution","satellite_systems_and_space_operations","observatoire_de_paris","chile","american_association_of_variable_star_observers"]}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6ImI5NWRiMzQ4LTBkZTMtNTc1OS04NmM1LTg0M2RlYjJhYWY2MyIsImMiOiIyMDI2LTA4LTA5VDEzOjM5OjU1LTA0OjAwIn0"},"meta":{"as_of":"2026-08-13T01:53:24.973088031Z"}}
~~~
**Payload gap:** Collection URL and Source omission is V1-accepted. The only future quality work is year-0001 sentinel handling. pagination.cursor remains a V1 contract gap.

### R08 GET /sources

**Live request:** GET /sources?limit=1 returned 200.

**Query gap:** q, domains, limit, cursor, and response_type are bound. Action: change default limit from 16 to 20, enforce 1-100 instead of only max=128, and echo the request cursor in pagination.cursor.

**Expected payload:**
~~~json
{"data":[{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","name":null,"url":"0.30000000000000004.com","description":null,"favicon_url":"https://0.30000000000000004.com/favicon.ico","rss_feed_url":null}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","url":"0.30000000000000004.com","favicon_url":"https://0.30000000000000004.com/favicon.ico"}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiayI6IjAuMzAwMDAwMDAwMDAwMDAwMDQuY29tIn0"},"meta":{"as_of":"2026-08-13T01:51:04.109728759Z"}}
~~~
**Payload gap:** Target null Source fields are omitted by omitempty.

### R09 GET /sources/{source_id}

**Live request:** GET /sources/52860efa-d345-58c4-8037-a9f3ae8987c0 returned 200.

**Query gap:** NONE. Both expect and accept only response_type.

**Expected payload:**
~~~json
{"data":{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","name":null,"url":"0.30000000000000004.com","description":null,"favicon_url":"https://0.30000000000000004.com/favicon.ico","rss_feed_url":null}}
~~~
**Actual payload:**
~~~json
{"data":{"id":"52860efa-d345-58c4-8037-a9f3ae8987c0","domain":"30000000000000004","url":"0.30000000000000004.com","favicon_url":"https://0.30000000000000004.com/favicon.ico"}}
~~~
**Payload gap:** Target null Source fields are omitted.

### R11 GET /tags

**Live request:** GET /tags?limit=1 returned 200.

**Query gap:** q, resource, limit, cursor, and response_type are bound. Action: apply the common pagination correction: default 20, range 1-100, and echo pagination.cursor.

**Expected payload:**
~~~json
{"data":[{"value":"\n\nmy_wife_quitter_job_episode_script_analysis"}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"value":"\n\nmy_wife_quitter_job_episode_script_analysis"}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiayI6IlxuXG5teV93aWZlX3F1aXR0ZXJfam9iX2VwaXNvZGVfc2NyaXB0X2FuYWx5c2lzIn0"},"meta":{"as_of":"2026-08-13T01:51:31.989277294Z"}}
~~~
**Payload gap:** Only shared limit/cursor envelope differences.

### R12 GET /entities

**Live requests:** GET /entities?limit=1 returned 200.

**Query gap:** TBD

**Expected payload:**
~~~json
{"data":[{"value":"007_first_light","type":"company"}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"value":"007_first_light","type":"company"}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwMDdfZmlyc3RfbGlnaHQiLCJ0eXBlIjoiY29tcGFueSJ9fQ"},"meta":{"as_of":"2026-08-13T01:52:24.247286918Z"}}
~~~
**Payload gap:** Default company response matches. Person output is type people,
not person; the target request spelling produces
{"error":{"code":"db_unavailable","message":"It's not you, it's us. Retry in a bit."}}.

### R13 GET /regions

**Live request:** GET /regions?limit=1 returned 200.

**Query gap:** All target parameters are bound. Action: apply the common pagination correction: default 20, range 1-100, and echo pagination.cursor.

**Expected payload:**
~~~json
{"data":[{"value":"0","type":"region"}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"value":"0","type":"region"}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwIiwidHlwZSI6InJlZ2lvbiJ9fQ"},"meta":{"as_of":"2026-08-13T01:52:29.751667796Z"}}
~~~
**Payload gap:** Only shared limit/cursor envelope differences.

### R14 GET /event-types

**Live request:** GET /event-types?limit=1 returned 200.

**Query gap:** All target parameters are bound. Action: apply the common pagination correction: default 20, range 1-100, and echo pagination.cursor.

**Expected payload:**
~~~json
{"data":[{"value":"0_60_mph_acceleration","type":"event_type"}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"data":[{"value":"0_60_mph_acceleration","type":"event_type"}],"pagination":{"limit":1,"next_cursor":"eyJ2IjoxLCJpZCI6bnVsbCwiZXQiOnsidmFsdWUiOiIwXzYwX21waF9hY2NlbGVyYXRpb24iLCJ0eXBlIjoiZXZlbnRfdHlwZSJ9fQ"},"meta":{"as_of":"2026-08-13T01:54:44.523440127Z"}}
~~~
**Payload gap:** Only shared limit/cursor envelope differences.

### R16 GET /events/count

**Live request:** GET /events/count returned 400 because registered /events/:id
captures count.

**Query gap:** Expected Event filters except q. Actual has no count binding;
the path is an invalid UUID.

**Expected payload:**
~~~json
{"data":{"count":42,"event_types":{"stock_decline":12},"impact_levels":{"high":17,"medium":25}},"meta":{"time_field":"created_at","as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"error":{"code":"invalid_request","message":"invalid UUID length: 5"}}
~~~
**Payload gap:** Aggregate payload is absent.

### R17 GET /events/summary

**Live request:** GET /events/summary?from=2026-08-01&to=2026-08-10&group_by=event_type
returned 400 because /events/:id captures summary.

**Query gap:** Expected bounded from/to and required group_by. Actual has no
summary binding and treats summary as an invalid UUID.

**Expected payload:**
~~~json
{"group_by":"event_type","data":[{"key":"stock_decline","event_count":12}],"meta":{"counted_resource":"event","time_field":"created_at","as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~json
{"error":{"code":"invalid_request","message":"invalid UUID length: 7"}}
~~~
**Payload gap:** Grouped aggregate payload is absent.

### R18 POST /events/search

**Live request:** POST /events/search with {"q":"wheat","limit":1} returned 404.

**Query gap:** Expected JSON q, filters, limit, cursor. Actual has no POST
registration or JSON binding; CORS allows only GET and OPTIONS.

**Expected payload:**
~~~json
{"data":[{"id":"6fe7144a-aa51-5fa1-a2dd-dd906c583fe3","kind":"event","created_at":"2026-08-10T01:07:00-04:00","tags":["and_reuse_food","home"],"summary":"On August 10, 2026, the USDA revised its wheat supply outlooks...","briefing":"On August 10, 2026, the USDA revised its wheat supply outlooks..."}],"pagination":{"limit":1,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
**Actual payload:**
~~~text
404 page not found
~~~
**Payload gap:** No route payload exists.

### R22 GET /actions and Action relation routes

**Live request:** GET /actions returned 404.

**Query gap:** NONE. R22 is reserved and specifies no published Action request parameters.

**Expected payload:** none. R22 is reserved until Action data exists.

**Actual payload:**
~~~text
404 page not found
~~~
**Payload gap:** NONE. R22 is reserved; the live 404 is consistent with no published Action route. Actions remain an upstream-data capability gap.

## Verification

1. Fresh source review of route binding, response mapping, SQL selection, and
   route registration.
2. Live HTTP capture of every target route, including discovery,
   aggregate-name collisions, POST, and Actions.
3. Targeted probes for category alias, sort, RFC 3339 time, and person entity
   type. go test ./... was run in apis/espresso earlier in this session and
   passed; this revision is documentation-only.
