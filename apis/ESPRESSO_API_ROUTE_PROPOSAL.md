# Espresso API Target Design and Implementation Specification

Status: Target design and implementation specification
Updated: 2026-08-11
Scope: Espresso API
Comparison baseline: [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md)

## 1. Purpose

This is the target public contract for Espresso. It follows the established
collection, detail, search, filtering, pagination, evidence, source, and
aggregation patterns documented for GDELT Cloud, PredictHQ, and Perigon.

Espresso keeps three deliberate extensions:

- Flattened intelligence digests, with a small stable resource core.
- Semantic vector search for Events and Signals.
- SAME_AS evidence and DERIVED_FROM support relationships.

The contract must not claim fields or routes that current Espresso data cannot
produce truthfully.

## 2. Current data model and capability states

The current database has only two record kinds:

~~~text
event
signal
~~~

There are no event subtypes and no Action records. Event routes therefore
select kind = 'event'; Signal routes select kind = 'signal'.

| State | Meaning |
|---|---|
| Current-data | The stored data exists. Query, handler, response, or documentation work may still be required. |
| Query-change | No ingestion change is required, but SQL or API implementation work is required. |
| Data-quality gap | The data exists but coverage or referential quality prevents an unconditional guarantee. |
| Upstream-data gap | The data required for a route or field is not stored. |
| Reserved | Future route; do not publish. |

Read-only production data assessment:

| Capability | Verified data | State |
|---|---:|---|
| Events | 851,272 rows | Current-data |
| Signals | 104,844 rows | Current-data |
| Actions | No rows | Upstream-data gap |
| Event embeddings | 560,671 rows | Data-quality gap: semantic Event search is partial |
| Signal embeddings | 104,844 rows | Current-data |
| Event briefing | 850,449 rows | Current-data: optional summary |
| Signal briefing | 104,844 rows | Current-data: optional summary |
| Event event_type | 791,545 rows | Current-data: optional classification/filter |
| Event impact_level | 792,385 rows | Current-data: optional classification/filter |
| Event companies, people, products, regions | 399,467; 365,642; 210,535; 558,767 rows | Current-data: exact-string filters/discovery only |
| Source references | 842,232 rows; 725,356 join to a Source | Data-quality gap |
| SAME_AS | 61,337,234 Event-to-Event relations | Current-data |
| DERIVED_FROM | 1,444,495 Signal-to-Event relations | Current-data |
| Created timestamp | 468 Events and 4 Signals use year-0001 | Data-quality gap |

Storage primitives:

~~~text
sips(id, kind, created, source, embedding, tags, digest, url, base_url)
sources(id, base_url, domain_name, site_name, description, favicon, rss_feed)
relations(from_id, to_id, relationship)
~~~

## 3. Industry-compatible contract

### 3.1 Common parameters

All collection routes use:

| Parameter | Type | Contract |
|---|---|---|
| limit | Integer 1-100; default 20 | Maximum returned items. |
| cursor | Opaque string | Keyset pagination. |
| response_type | json, yaml, toon | Equivalent serializations. JSON is canonical. |

Search collections also use:

| Parameter | Type | Contract |
|---|---|---|
| q | String, max 1024 | Semantic search query. |
| from / to | Date or RFC 3339 date-time | Inclusive record-creation bounds. |
| sort | recent, relevance | recent is the default. relevance requires q. |

Date-only values are interpreted in UTC: from starts at 00:00:00Z,
and to ends at 23:59:59.999999999Z. RFC 3339 values preserve their
instant and are normalized to UTC.

Espresso does not accept start, end, active, updated, cancelled, predicted_end,
place, country, latitude, longitude, rank, attendance, spend, or numeric
confidence filters. Those are data gaps, not undocumented filters.

### 3.1.1 Industry compatibility reference

The following names are used by comparison APIs. They are reference material,
not Espresso request parameters.

| Capability | Industry names | Espresso decision |
|---|---|---|
| Pagination | GDELT: cursor; PredictHQ: offset; Perigon: page and size | Maintain cursor only. It is safe for a changing feed and matches GDELT. |
| Search | GDELT: search; PredictHQ and Perigon: q | Maintain q only. |
| Creation-time range | GDELT: date_start/date_end; Perigon: from/to | Maintain from/to only. |
| Total matches | PredictHQ: count; Perigon: numResults; GDELT does not include a total | Do not calculate or return a total by default. |

There is no industry-standard field name for the number of records in the
current page. GDELT omits it, while PredictHQ count and Perigon numResults mean
total matching records, not page length. Espresso therefore omits
returned_count. Clients obtain the page length from data.length.

### 3.2 Collection response

Espresso uses a data/pagination envelope with cursor pagination. It intentionally does not include a success field. PredictHQ and Perigon do not use a success field either; Espresso relies on the HTTP status code and the response body.

~~~json
{
  "data": [],
  "pagination": {
    "limit": 20,
    "cursor": null,
    "next_cursor": null
  },
  "meta": {
    "as_of": "2026-08-11T15:00:00Z"
  }
}
~~~

- Empty collections return 200 with data: [].

Detail responses are:

~~~json
{"data": {}}
~~~

Successful responses use HTTP 2xx status codes. Errors use HTTP status codes such as 400, 404, 429, or 500 and an error body; the Espresso contract does not return success=false.

### 3.3 Event and Signal payload

Every Event or Signal item has a stable core plus actual flattened digest
members. The stable core is:

| Field | Source | Rule |
|---|---|---|
| id | sips.id | Always present. |
| kind | sips.kind | Always present; its only current values are event and signal. |
| created_at | sips.created | Null with data_quality.created_at="unknown" for a year-0001 sentinel. |
| url | sips.url | Present when stored. |
| base_url | sips.base_url | Present when stored. |
| source_id | sips.source | Null when no source exists. |
| source | sources join | Optional object with id, domain, name, and url. Omitted when the source reference is null or orphaned. |
| tags | sips.tags | Present when stored. |
| summary | digest.briefing | Optional normalized projection. |
| data_quality | Derived | Present only for a known quality condition. |

The API also returns every non-empty digest member at the item root. briefing is
therefore preserved even when summary is present. Canonical core fields win on
a name conflict. The API exposes kind as a flattened core field; it never exposes embedding, relation direction, or a nested digest object.

~~~json
{
  "id": "event-id",
  "kind": "event",
  "created_at": "2026-08-11T05:37:59Z",
  "url": "https://example.com/article",
  "base_url": "https://example.com",
  "source_id": "source-id",
  "source": {
    "id": "source-id",
    "domain": "example.com",
    "name": null,
    "url": "https://example.com"
  },
  "tags": ["investment_and_capital_markets"],
  "summary": "A normalized briefing.",
  "briefing": "A normalized briefing.",
  "event_type": "stock_decline",
  "impact_level": "medium",
  "companies": ["example_company"]
}
~~~

There is no stored title. Clients use summary/briefing until ingestion provides a
stable title.


### 3.4 Digest field model

The response item is the union of flattened core fields (id, kind, created_at,
url, base_url, source_id, source, and tags) and the actual fields in that row's
digest. The complete field set varies by Event or Signal; clients must not
require every digest field to be present.

The following digest fields are common intelligence fields, not a closed schema:

| Field | Meaning | Availability |
|---|---|---|
| briefing | Narrative explanation of the Event or Signal. Also projected as optional summary. | Common in Events; current on all Signals. |
| drivers | Factors causing the development. | Optional string array. |
| impacts | Observed or expected effects. | Optional string array. |
| forecast | Forward-looking expectation. | Optional string. |
| activities | Dated or undated supporting activities. | Optional string array. |
| event_type | Classification of an Event. | Optional; not expected on Signals. |
| impact_level | Categorical impact assessment. | Optional string. |
| macro_context | Broader condition or theme surrounding an Event. | Optional string. |
| impacted_domains | Domains affected by the development. | Optional string array. |
| companies, people, products, regions | Extracted exact-string filter values. | Optional arrays; not canonical entity or geography IDs. |
| confidence | Source/ingestion confidence where available. | Optional and currently non-numeric; do not treat it as an industry-style score. |

The stored tags column is canonical. If a digest also contains tags, the
response emits one reconciled top-level tags value rather than two fields.

## 4. Target routes

Gateway paths add /espresso. Backend routes omit that prefix.

| ID | Route | What it does | Origin | State |
|---|---|---|---|---|
| R01 | GET /events | Searches Event records. | Industry route: GDELT Events; PredictHQ Events search. | Current-data + Query-change |
| R02 | GET /events/{event_id} | Retrieves one Event. | Industry route: GDELT Event detail; PredictHQ uses an ID filter instead. | Current-data + Query-change |
| R03 | GET /events/{event_id}/evidence | Returns direct SAME_AS Event evidence. | Industry pattern adapted from GDELT Story article evidence; Espresso-specific relation semantics. | Current-data + Query-change |
| R04 | GET /events/{event_id}/signals | Returns Signals derived from the Event or its direct evidence neighbours. | Espresso extension. | Current-data + Query-change |
| R05 | GET /signals | Searches Signal records. | Espresso extension; not Perigon saved-monitor Signals. | Current-data + Query-change |
| R06 | GET /signals/{signal_id} | Retrieves one Signal. | Espresso extension. | Current-data + Query-change |
| R07 | GET /signals/{signal_id}/events | Returns direct Event support for one Signal. | Espresso extension. | Current-data + Query-change |
| R08 | GET /sources | Searches Source records. | Industry route/pattern: Perigon Sources. | Current-data + Query-change |
| R09 | GET /sources/{source_id} | Retrieves one Source. | Espresso extension over industry Source discovery. | Current-data + Query-change |
| R11 | GET /tags | Discovers stored tag values. | Espresso extension; comparable APIs expose filter values differently. | Current-data + Query-change |
| R12 | GET /entities | Discovers company and person strings in Event digests. | Industry pattern adapted from GDELT Entities and PredictHQ entity filtering. | Current-data + Query-change |
| R13 | GET /regions | Discovers region strings in Event digests. | Industry pattern adapted from GDELT geography discovery and PredictHQ Places. | Current-data + Query-change |
| R14 | GET /event-types | Discovers Event event_type values. | Industry pattern adapted from category/filter discovery. | Current-data + Query-change |
| R16 | GET /events/count | Returns Event count distributions. | Industry route: PredictHQ Event Counts. | Current-data + Query-change |
| R17 | GET /events/summary | Returns grouped Event aggregates. | Industry route: GDELT Event Summary. | Current-data + Query-change |
| R18 | POST /events/search | Runs Event vector search with JSON filters. | Industry pattern adapted from Perigon Vector News Search. | Current-data + Query-change |
| R22 | GET /actions and Action relation routes | Not published. | Espresso extension; reserved because no Action data exists. | Reserved / Upstream-data gap |

### 4.1 Event routes

GET /events accepts:

| Parameter | Meaning |
|---|---|
| ids | CSV UUIDs, maximum 128. |
| event_types | Exact digest.event_type values. |
| categories | Alias for event_types. Espresso has no separate category taxonomy. |
| impact_levels | Exact digest.impact_level values. |
| companies, people, products, regions | Exact-string membership in the corresponding digest array. |
| entities | Exact-string membership in companies or people. |
| source_ids | Direct source UUID filter. It does not expand evidence. |
| tags | Persisted tag overlap. |
| search, time, sort, and collection parameters | Section 3. |

GET /events/{event_id} returns one Event plus:

~~~json
{
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


#### GET /events/{event_id}/evidence

Returns the requested Event plus its direct SAME_AS Event neighbours.

This is Espresso's closest equivalent to a GDELT Story article-evidence route.
It does not claim article title, publication time, content, or citation-role
data that Espresso does not store.

Accepted evidence parameters:

| Parameter | Contract |
|---|---|
| source_ids | CSV UUIDs. Restricts direct SAME_AS neighbours to those Source IDs; the requested Event remains included. |
| from / to | Restrict direct neighbours by creation time. |
| limit, cursor, response_type | Common collection parameters. |

Evidence does not accept search or sort. It has fixed recent ordering by
created_at, then id, because it is a provenance drilldown rather than a
semantic-search collection.

#### GET /events/{event_id}/signals

Expands direct SAME_AS Event neighbours before
finding derived Signals. It accepts the filters of its returned Signal resource:

| Parameter | Contract |
|---|---|
| ids | CSV Signal UUIDs, maximum 128. |
| impact_levels | Exact Signal digest.impact_level values. |
| impacted_domains | Exact membership in the Signal digest array. |
| source_ids | Direct Signal Source UUIDs. |
| tags | Persisted Signal tag overlap. |
| q | Semantic query. |
| from / to | Signal creation-time bounds. |
| sort | recent or relevance; relevance requires q. |
| limit, cursor, response_type | Common collection parameters. |

### 4.2 Signal routes

GET /signals accepts ids, impact_levels, impacted_domains, tags, search,
time, sort, and the common collection parameters.

GET /signals/{signal_id} returns one Signal plus:

~~~json
{
  "links": {
    "events": "/signals/signal-id/events"
  },
  "counts": {
    "events": 4
  }
}
~~~

#### GET /signals/{signal_id}/events

Returns direct DERIVED_FROM Event targets. It
accepts the filters of its returned Event resource:

| Parameter | Contract |
|---|---|
| ids | CSV Event UUIDs, maximum 128. |
| event_types / categories | Exact digest.event_type values; categories is the compatibility alias. |
| impact_levels | Exact digest.impact_level values. |
| companies, people, products, regions | Exact membership in the corresponding Event digest arrays. |
| entities | Exact membership in Event companies or people. |
| source_ids | Direct Event Source UUIDs. |
| tags | Persisted Event tag overlap. |
| q | Semantic query. |
| from / to | Event creation-time bounds. |
| sort | recent or relevance; relevance requires q. |
| limit, cursor, response_type | Common collection parameters. |

### 4.3 Source routes

GET /sources accepts q, domains, limit, cursor, and response_type.
Its q parameter is case-insensitive Source text matching, not vector search.

GET /sources/{source_id} returns:

~~~json
{
  "data": {
    "id": "source-id",
    "domain": "example.com",
    "name": null,
    "url": "https://example.com",
    "description": null,
    "favicon_url": null,
    "rss_feed_url": null
  }
}
~~~

The public Source mapping is:

| Stored field | Public field |
|---|---|
| sources.id | id |
| sources.domain_name | domain |
| sources.site_name | name |
| sources.base_url | url |
| sources.description | description |
| sources.favicon | favicon_url |
| sources.rss_feed | rss_feed_url |

### 4.4 Discovery routes

All discovery routes accept q, limit, cursor, and response_type.

| Route | Extra parameter | Item payload |
|---|---|---|
| GET /tags | resource=event,signal | {"value":"tag"} |
| GET /entities | types=company,person | {"value":"name","type":"company"} |
| GET /regions | none | {"value":"region", "type": "region"} |
| GET /event-types | none | {"value":"stock_decline", "type": "event_type"} |

These routes expose exact stored values. They do not provide canonical entity
profiles, aliases, countries, coordinates, or geography hierarchy.

### 4.5 Count and summary routes

GET /events/count accepts Event filters except semantic search in version one.

~~~json
{
  "data": {
    "count": 42,
    "event_types": {"stock_decline": 12},
    "impact_levels": {"high": 17, "medium": 25}
  },
  "meta": {
    "time_field": "created_at",
    "as_of": "2026-08-11T15:00:00Z"
  }
}
~~~

GET /events/summary requires bounded from/to and group_by:

~~~text
created_day | event_type | impact_level | source | tag | region
~~~

~~~json
{
  "group_by": "event_type",
  "data": [
    {"key": "stock_decline", "event_count": 12}
  ],
  "meta": {
    "counted_resource": "event",
    "time_field": "created_at",
    "as_of": "2026-08-11T15:00:00Z"
  }
}
~~~

Tag and region buckets expand arrays and are not additive. Source buckets use
direct source IDs. These routes require a bounded time range, query-plan review,
and measured expression/GIN indexes before publication.

### 4.6 Vector extensions

POST /events/search provides a Perigon-style JSON Event-search body. Its
response is identical to GET /events.

~~~json
{
  "q": "semiconductor demand weakness",
  "filters": {
    "from": "2026-08-01T00:00:00Z",
    "impact_levels": ["high"],
    "tags": ["guidance"]
  },
  "limit": 20,
  "cursor": null
}
~~~

## 5. Query rules

1. Event queries use s.kind = 'event'. Signal queries use s.kind = 'signal'.
2. Select id, kind, created, source, tags, digest, url, and base_url; left join
   Source for the optional source projection.
3. Semantic queries require s.embedding IS NOT NULL and order by distance, id.
   Recent queries order by created, id.
4. Exclude year-0001 timestamps from default and time-filtered collections.
5. Fetch limit+1 items for cursor pagination.
6. SAME_AS is bidirectional. DERIVED_FROM uses its stored direction internally.
7. Verify a subresource parent before querying children.
8. Never interpolate a client-supplied JSON key or SQL expression.

Use correct PostgreSQL operators:

~~~sql
-- Scalar JSON string fields
(s.digest->>'event_type') = ANY(@event_types::text[])
(s.digest->>'impact_level') = ANY(@impact_levels::text[])

-- JSON arrays of strings
COALESCE(s.digest->'companies', '[]'::jsonb) ?| @companies::text[]
COALESCE(s.digest->'people', '[]'::jsonb) ?| @people::text[]
COALESCE(s.digest->'regions', '[]'::jsonb) ?| @regions::text[]

-- Persisted SQL array
s.tags && @tags::text[]
~~~

The current scalar use of ?| for event_type and impact_level is incorrect.
Those comparisons must use ->> followed by = ANY(...).

## 6. Current implementation gaps

| Area | Current state | Target correction |
|---|---|---|
| Pagination | Default limit 16; cursor only; no page item count. | Default 20; cursor pagination only. Clients use data.length for page length. |
| Time | Date-only from/to | RFC 3339 support. |
| Search | Current parameter does not use the target name. | Keep q. |
| Ordering | No public sort | Add recent and search-only relevance. |
| Relation semantic search | Relation inputs do not bind the semantic query | Bind q and apply vector ordering. |
| Payload | kind is exposed; canonical source/URL/summary mapping is incomplete | Hide kind; emit stable core plus flattened digest. |
| Evidence | event_id and created fields | id and created_at with full Event evidence projection. |
| Source response | Storage-shaped/inconsistent collection and detail mapping | Use the normalized Source payload. |
| Serialization docs | Older docs advertise text | Document json/yaml/toon only. |
| Auth config | Runtime reads API_KEY; platform docs use API_KEYS | Correct runtime configuration. |
| Additional routes | Discovery, count/summary, and Event JSON vector search are absent. | Implement the published additional routes after query-plan and response tests. |

## 7. Future gaps

These industry capabilities remain gaps until ingestion provides the required
data:

| Capability | Gap |
|---|---|
| Stories and article clusters | SAME_AS is identity evidence, not narrative clustering or article membership. |
| Article detail fields | No stable title, published timestamp, body, or citation role. |
| Event occurrence/lifecycle time | created is record creation only. |
| Places and structured geography | No country codes, latitude/longitude, or hierarchy. |
| Canonical entities | Companies and people are strings without IDs, aliases, or profiles. |
| Topic taxonomy | No reliable topics/taxonomy store. |
| Quantitative metrics | No rank, significance, numeric confidence, attendance, spend, fatalities, or methodology. |
| Complete semantic Event coverage | 290,601 Events have no embedding. |
| Complete Source enrichment | Source references can be orphaned; metadata coverage is partial. |
| Actions | No Action records or relations. |

## 8. Publication order

1. Correct existing Event, Signal, evidence, Source, tag, date, parameter,
   scalar-filter, pagination, sort, serialization, and API_KEYS behavior.
2. Publish discovery, count, and summary routes after aggregate query review.
3. Publish Event JSON vector search.
4. Backfill Event embeddings, repair Source references, and replace sentinel
   timestamps.
5. Add Stories, articles, occurrence/lifecycle time, places, canonical entities,
   topics, metrics, and Actions only after upstream data exists.
