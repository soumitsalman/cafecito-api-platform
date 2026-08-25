# Espresso Bruno collection

Maintained executable examples for the Espresso REST contract. This collection is not the public portal; portal pages link here as runnable requests.

**Audience:** maintainers and integrators  
**Authority:** live routes in `apis/espresso/router/` and `config/espresso.oas.json`  
**Last verified:** 2026-08-25

## Setup

1. Open `apis/espresso/tests/bruno` in Bruno.
2. Set collection variables:

| Variable | Example | Notes |
| --- | --- | --- |
| `baseUrl` | `https://api.cafecito.tech/espresso` | Gateway prefix included. Local backend is typically `http://localhost:8080` without `/espresso`. |
| API key | Bearer token | Put `Authorization: Bearer <key>` on protected requests. Health does not require it. |

3. Run Health first, then a collection search such as Find concrete Events.

## Contract reminders

- Send `pagination.next_cursor` unchanged as `cursor`. Collections do not return `pagination.cursor`.
- Empty collections are HTTP 200 with `data: []`. Missing detail is HTTP 404.
- Errors are `{ "error": { "code", "message" } }`.
- Event/Signal stable core: `id`, `kind`, `created_at`, `tags`. Ignore unknown extension fields.
- `response_type` is `json`, `yaml`, or `toon` (same logical payload).

## Requests in this collection

| Request | REST path | MCP tool (if exported) |
| --- | --- | --- |
| Check Espresso service availability | GET /health | REST-only |
| Find concrete Events | GET /events | searchEvents |
| Inspect one Event | GET /events/{id} | getEvent |
| Inspect evidence for an Event | GET /events/{id}/evidence | getEventEvidence |
| Find Signals connected to an Event | GET /events/{id}/signals | getEventSignals |
| Find synthesized Signals | GET /signals | searchSignals |
| Inspect one Signal | GET /signals/{id} | getSignal |
| Inspect Events supporting a Signal | GET /signals/{id}/events | getSignalEvents |
| Find intelligence Sources | GET /sources | listIntelligenceSources |
| Inspect one Source | GET /sources/{id} | getIntelligenceSource |
| Discover fuzzy tag vocabulary | GET /tags | listIntelligenceTags |
| Discover exact Event entity filters | GET /entities | listIntelligenceEntities |
| Discover exact Event region filters | GET /regions | listIntelligenceRegions |
| Discover exact Event type filters | GET /event-types | listIntelligenceEventTypes |

Do not add Beans requests to this folder.
