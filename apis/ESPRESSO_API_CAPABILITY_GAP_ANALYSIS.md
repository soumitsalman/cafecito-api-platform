# Espresso API Capability Gap Analysis

Last updated: 2026-08-12T21:30:57-04:00
Assessment basis: source review of Espresso, then go test ./... in apis/espresso (passes).

## References

- Target: [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md)
- Industry baseline: [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md)
- Reviewed: apis/espresso/router/params.go, routes.go, responses.go; apis/espresso/db/queries.go, types.go.

Gateway paths add /espresso; routes below are backend paths.

## Shared Facts

The proposal conflicts on kind: Section 3.3 requires it; Section 6 says to hide
it. Expected payloads use Section 3.3. This is a proposal correction, not an
implementation gap.

All current collection routes have this parameter gap.

Expected: limit 1-100 (default 20), cursor, response_type=json|yaml|toon;
search routes also accept date or RFC 3339 from/to and sort=recent|relevance.

Actual: limit form-default is 16, maximum is 128, and there is no minimum;
from/to accepts only YYYY-MM-DD; no handler binds sort. json, yaml, and toon
are implemented.

Expected collection payload:
~~~json
{"data":[],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual collection payload:
~~~json
{"data":[],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Runtime correctly omits success, but does not emit pagination.cursor.

Event and Signal collection SQL selects only id, kind, created, tags, digest.
It cannot emit url, base_url, source_id, or source. Its mapper does not project
briefing to summary, reconcile digest tags with sips tags, or convert a
year-0001 timestamp to null plus data_quality.created_at=unknown.

## Route Analysis

### R01 GET /events

Expected parameters: ids, event_types, categories, impact_levels, companies,
people, products, regions, entities, source_ids, tags, q, from, to, sort,
limit, cursor, response_type.

Actual parameters: ids, event_types, impact_levels, companies, people,
products, regions, source_ids, tags, q, extra acc, date-only from/to, limit,
cursor, response_type.

Expected payload:
~~~json
{"data":[{"id":"uuid","kind":"event","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"summary":"optional","briefing":"optional","...digest":"members"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"id":"uuid","kind":"event","created_at":"timestamp","tags":[],"...digest":"members"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: categories, entities, sort absent; acc extra. Stable URL and Source,
summary, quality, and cursor fields are absent. Tags use tags_fts instead of
array overlap. event_type and impact_level use JSONB ?| instead of ->> = ANY.

### R02 GET /events/{event_id}

Expected parameters: response_type.
Actual parameters: response_type.

Expected payload:
~~~json
{"data":{"id":"uuid","kind":"event","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"summary":"optional","links":{"evidence":"/events/uuid/evidence","signals":"/events/uuid/signals"},"counts":{"evidence":3,"signals":2}}}
~~~
Actual payload:
~~~json
{"data":{"id":"uuid","kind":"event","created_at":"timestamp","tags":[],"...digest":"members","source":{"id":"uuid","url":"optional"},"links":{"evidence":"/events/uuid/evidence","signals":"/events/uuid/signals"},"counts":{"coverage":3,"actions":0,"signals":2}}}
~~~
Gap: GetSip inner-joins sources, giving false 404s for null or orphan source
references. Top-level URL, base URL, source ID, summary, and quality are
absent. evidence is named coverage and unsupported actions is added.

### R03 GET /events/{event_id}/evidence

Expected parameters: source_ids, from, to, limit, cursor, response_type; no q
or sort.
Actual parameters: the same names, but date-only from/to and shared pagination.

Expected payload:
~~~json
{"data":[{"id":"uuid","kind":"event","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","tags":[],"...digest":"members"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"event_id":"uuid","created":"timestamp","source_id":"uuid","url":"optional","base_url":"optional"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: source_ids filters the anchor, so the required requested Event can vanish.
The narrow payload uses event_id and created and omits kind, tags, digest,
source object, and quality fields.

### R04 GET /events/{event_id}/signals

Expected parameters: ids, impact_levels, impacted_domains, source_ids, tags,
q, from, to, sort, limit, cursor, response_type.
Actual parameters: impact_levels, impacted_domains, tags, date-only from/to,
limit, cursor, response_type.

Expected payload:
~~~json
{"data":[{"id":"uuid","kind":"signal","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"summary":"optional","...digest":"members"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload: a valid Event ID returns 404 because the handler verifies a
Signal. With a Signal ID on this Event path, it filters results as kind=event
and returns Event items.

Gap: route semantics are inverted. ids, source_ids, q, sort, and SAME_AS
expansion are absent, and the target Signal payload is unreachable.

### R05 GET /signals

Expected parameters: ids, impact_levels, impacted_domains, tags, q, from, to,
sort, limit, cursor, response_type.
Actual parameters: those names except sort, plus acc, date-only from/to, limit,
cursor, response_type.

Expected payload:
~~~json
{"data":[{"id":"uuid","kind":"signal","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"summary":"optional","...digest":"members"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"id":"uuid","kind":"signal","created_at":"timestamp","tags":[],"...digest":"members"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: sort is absent; acc is extra; time and pagination differ. URL, Source,
summary, quality, and cursor are absent. Tags use full-text matching, and
impact_levels uses the wrong scalar JSONB predicate.

### R06 GET /signals/{signal_id}

Expected parameters: response_type.
Actual parameters: response_type.

Expected payload:
~~~json
{"data":{"id":"uuid","kind":"signal","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"summary":"optional","links":{"events":"/signals/uuid/events"},"counts":{"events":4}}}
~~~
Actual payload:
~~~json
{"data":{"id":"uuid","kind":"signal","created_at":"timestamp","tags":[],"...digest":"members","source":{"id":"uuid","url":"optional"},"links":{"events":"/signals/uuid/events"},"counts":{"events":4}}}
~~~
Gap: inner Source join can give a false 404. Top-level URL, base URL, source
ID, summary projection, and sentinel-time quality handling are absent.

### R07 GET /signals/{signal_id}/events

Expected parameters: ids, event_types, categories, impact_levels, companies,
people, products, regions, entities, source_ids, tags, q, from, to, sort,
limit, cursor, response_type.
Actual parameters: event_types, impact_levels, undocumented impacted_domains,
tags, date-only from/to, limit, cursor, response_type.

Expected payload:
~~~json
{"data":[{"id":"uuid","kind":"event","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"summary":"optional","...digest":"members"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"id":"uuid","kind":"event","created_at":"timestamp","tags":[],"...digest":"members"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: direct DERIVED_FROM traversal is correct. ids, categories, entity-array
filters, entities, source_ids, q, and sort are absent; impacted_domains is not
a target Event filter. Relation output omits URL, base URL, source ID, Source,
summary, and quality fields.

### R08 GET /sources

Expected parameters: q, domains, limit, cursor, response_type.
Actual parameters: the same names, with shared actual pagination.

Expected payload:
~~~json
{"data":[{"id":"uuid","domain":"example.com-or-null","name":"example-or-null","url":"https://example.com","description":"text-or-null","favicon_url":"url-or-null","rss_feed_url":"url-or-null"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"id":"uuid","domain":"optional","name":"optional","url":"https://example.com","description":"optional","favicon_url":"optional","rss_feed_url":"optional"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: field names and search match, but unavailable optional fields are omitted,
not null. Pagination differs.

### R09 GET /sources/{source_id}

Expected parameters: response_type.
Actual parameters: response_type.

Expected payload:
~~~json
{"data":{"id":"uuid","domain":"example.com-or-null","name":"example-or-null","url":"https://example.com","description":"text-or-null","favicon_url":"url-or-null","rss_feed_url":"url-or-null"}}
~~~
Actual payload:
~~~json
{"data":{"id":"uuid","domain":"optional","name":"optional","url":"https://example.com","description":"optional","favicon_url":"optional","rss_feed_url":"optional"}}
~~~
Gap: public names and route behavior match. Optional Source values are omitted,
not null.

### R11 GET /tags

Expected parameters: q, resource=event,signal, limit, cursor, response_type.
Actual parameters: q, CSV resource, limit, cursor, response_type.

Expected payload:
~~~json
{"data":[{"value":"tag"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"value":"tag"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: item payload and resource scope conform. Limit and cursor response
contract do not.

### R12 GET /entities

Expected parameters: q, types=company,person, limit, cursor, response_type.
Actual parameters: q, types=company,person,product,stock_ticker, limit,
cursor, response_type.

Expected payload:
~~~json
{"data":[{"value":"name","type":"company-or-person"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload when types is omitted:
~~~json
{"data":[{"value":"name","type":"company-or-people"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: product and stock_ticker are outside target. types=person passes binding
but QueryEventTags recognizes people, not person, so it fails. Person output
uses people, not person. Pagination differs.

### R13 GET /regions

Expected parameters: q, limit, cursor, response_type.
Actual parameters: the same names, with shared actual pagination.

Expected payload:
~~~json
{"data":[{"value":"region","type":"region"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"value":"region","type":"region"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: route, parameters, and item payload conform. Shared limit and
pagination.cursor gaps remain.

### R14 GET /event-types

Expected parameters: q, limit, cursor, response_type.
Actual parameters: the same names, with shared actual pagination.

Expected payload:
~~~json
{"data":[{"value":"stock_decline","type":"event_type"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload:
~~~json
{"data":[{"value":"stock_decline","type":"event_type"}],"pagination":{"limit":16,"next_cursor":"opaque-or-null"},"meta":{"as_of":"RFC3339"}}
~~~
Gap: route, parameters, and item payload conform. Shared limit and
pagination.cursor gaps remain.

### R16 GET /events/count

Expected parameters: Event filters except q; this is not paginated.
Actual parameters: none; route is not registered.

Expected payload:
~~~json
{"data":{"count":42,"event_types":{"stock_decline":12},"impact_levels":{"high":17}},"meta":{"time_field":"created_at","as_of":"RFC3339"}}
~~~
Actual payload: none; path is unregistered.
Gap: route, aggregate query, bindings, and response model are absent.

### R17 GET /events/summary

Expected parameters: bounded from/to, required
group_by=created_day|event_type|impact_level|source|tag|region. This is not paginated.
Actual parameters: none; route is not registered.

Expected payload:
~~~json
{"group_by":"event_type","data":[{"key":"stock_decline","event_count":12}],"meta":{"counted_resource":"event","time_field":"created_at","as_of":"RFC3339"}}
~~~
Actual payload: none; path is unregistered.
Gap: route, bounded-window validation, grouping, aggregate query, and response
model are absent.

### R18 POST /events/search

Expected request body:
~~~json
{"q":"semiconductor demand weakness","filters":{"from":"2026-08-01T00:00:00Z","impact_levels":["high"],"tags":["guidance"]},"limit":20,"cursor":null}
~~~
Actual request handling: no POST route; CORS allows only GET and OPTIONS.

Expected payload:
~~~json
{"data":[{"id":"uuid","kind":"event","created_at":"RFC3339-or-null","url":"optional","base_url":"optional","source_id":"uuid-or-null","source":{},"tags":[],"...digest":"members"}],"pagination":{"limit":20,"cursor":null,"next_cursor":null},"meta":{"as_of":"RFC3339"}}
~~~
Actual payload: none; path is unregistered.
Gap: POST registration, JSON binding, Event filters, vector execution, and
target response are absent.

### R22 GET /actions and Action relation routes

Expected request and payload: no published route or payload; Actions are
reserved because no Action data exists.
Actual request and payload: no Action route is registered.
Gap: none for the reserved publication state. Actions remain an upstream-data
gap.

## Cross-Route Gaps

- Vector SQL does not require embedding IS NOT NULL.
- Generated Swagger and gateway OpenAPI still advertise obsolete text output,
  Event-family kinds, and parameters that request binders do not implement.
- Runtime reads API_KEY, while the target platform configuration is API_KEYS.

## Verification

1. First pass: compared target routes with request binders, handler
   registrations, response writers, and SQL projections.
2. Second pass: removed stale claims that Espresso does not compile, that
   runtime success is a gap, that R03 needs evidence-role metadata, and that
   R13 or R14 item types are wrong.
3. Third pass: ran go test ./... in apis/espresso; all packages passed.
