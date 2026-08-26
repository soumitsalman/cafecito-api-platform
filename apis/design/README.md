# Design records index

Internal Cafecito design and research notes. These files are **not** the public API contract. Published behavior lives in gateway OpenAPI, the developer portal, and generated backend Swagger.

**Last verified:** 2026-08-25  
**Owner role:** documentation steward

## Public contract authority

| Fact | Authority (do not invent from this folder) |
|---|---|
| Beans routes, schemas, envelopes | [`config/beans.oas.json`](../../config/beans.oas.json), portal `/api/beans` |
| Espresso routes, schemas, envelopes | [`config/espresso.oas.json`](../../config/espresso.oas.json), portal `/api/espresso` |
| Human guides | [`docs/pages/`](../../docs/pages/) (see [`docs/README.md`](../../docs/README.md)) |
| Service-local OpenAPI | regenerated `apis/<service>/docs/` from router annotations ([`apis/README.md`](../README.md)) |

**Published Beans collection envelope (V1):** `{ data, pagination, meta }` with `pagination.limit`, `pagination.num_results` (this page only), `pagination.next_cursor`, and `meta.as_of`. Empty collections are HTTP 200 with `data: []`. Do not treat older proposal examples that omit `num_results` or restrict `meta.as_of` as live contract.

## Product status (for contributors)

- **Live:** Beans, Espresso (REST and MCP).
- **Future / reserved:** Cortado, Latte (not live APIs).
- **SDK:** there is **no official SDK**; clients use REST or MCP.

## Current documents

| Document | Role |
|---|---|
| [DOCUMENTATION_GAP_REPORT.md](DOCUMENTATION_GAP_REPORT.md) | Current repository documentation audit |
| [QUERY_PARAM_NAMES.md](QUERY_PARAM_NAMES.md) | Current naming decision for unstructured search (`q`) |
| [NEWS_AND_BLOG_API_MARKET_REPORT.md](NEWS_AND_BLOG_API_MARKET_REPORT.md) | Current external news/blog API comparison (not Cafecito contract) |
| [INDUSTRY_EVENT_API_ROUTE_REFERENCE.md](INDUSTRY_EVENT_API_ROUTE_REFERENCE.md) | Current external event/news API comparison (not Cafecito contract) |
| [UNUSUAL_WHALES_API_INVENTORY.md](UNUSUAL_WHALES_API_INVENTORY.md) | Current competitor route inventory (research) |

## Historical or superseded documents

Older route proposals and execution plans remain for history. They do **not** override gateway OpenAPI.

| Document | Status | Superseded by |
|---|---|---|
| [BEANS_API_ROUTE_PROPOSAL.md](BEANS_API_ROUTE_PROPOSAL.md) | superseded | `config/beans.oas.json` |
| [ESPRESSO_API_ROUTE_PROPOSAL.md](ESPRESSO_API_ROUTE_PROPOSAL.md) | superseded | `config/espresso.oas.json` |

## Public documentation boundary

Do not copy internal storage, retrieval, or operations detail into `docs/pages/`, gateway OAS, or generated Swagger. Describe client-visible resources (Articles, Events, Signals, Sources, pagination, errors, API keys). Never publish table names, embeddings, HNSW, gateway policy internals, or backend credentials. Full rule: root [`AGENTS.md`](../../AGENTS.md) and [`docs/CONTRIBUTING.md`](../../docs/CONTRIBUTING.md).
