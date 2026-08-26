# Backend APIs (maintainer only)

**Status:** current  
**Audience:** service maintainers  
**Authority:** live contract is each service's `router/` bindings, handlers, and response structs; generated Swagger is produced from those annotations. Public docs live in the portal (`docs/pages/`) and gateway OpenAPI (`config/beans.oas.json`, `config/espresso.oas.json`).  
**Last verified:** 2026-08-25

This file is **not** the public API contract. Do not copy route tables, auth headers, or storage details from here into portal pages.

Public product docs: [Beans overview](../docs/pages/products/beans/overview.mdx), [Espresso overview](../docs/pages/products/espresso/overview.mdx). Gateway paths use `/beans` and `/espresso`. Local Go processes expose the same routes **without** those prefixes, typically on `:8080`.

## Local run

Required: `PG_CONNECTION_STRING`. Optional: `EMBEDDER_BASE_URL`, `EMBEDDER_API_KEY`, `EMBEDDER_MODEL`, `PORT`, `API_KEY`.

`API_KEY` is semicolon-separated `Header=Value` pairs for the **backend** process. Public clients authenticate with a Bearer API key at the gateway. Do not document private backend headers on public surfaces. If `API_KEY` is unset, backend auth is disabled.

```bash
cd apis/beans && go run .     # :8080
cd apis/espresso && go run .  # :8080 (run one at a time natively)
```

`GET /health` is unauthenticated. Other REST routes use `apiKeyMiddleware` when `API_KEY` is set.

## Tests and CI fixtures

Integration tests live under each service's `tests/` directory. Env is loaded from that service's `.env` (at least `PG_CONNECTION_STRING`; Beans also needs embedder vars for some tests).

The manual **API Contract & Docs** workflow (`.github/workflows/api-contract-docs.yml`, `workflow_dispatch` only) starts `pgvector/pgvector:pg16`, applies `.github/ci/schema-beans.sql` or `schema-espresso.sql`, seeds Espresso from `.github/ci/seed-espresso.sql`, and relies on Beans `tests/fixtures_test.go` for deterministic rows. It does not run on PRs or gate deploys. Router vector-search tests that call a live embedder are skipped in that workflow (`-skip 'TestRouterVectorSearch'`). Stress tests skip when no process is listening on `:8080`.

```bash
# Beans (seeds CI fixtures; fake embedder)
cd apis/beans && go test ./tests/...

# Espresso (hermetic contract tests always; DB tests need PG + fixtures)
cd apis/espresso && go test ./tests/...
```

Bruno collections: `apis/beans/bruno/` and `apis/espresso/bruno/`.

## Regenerate Swagger

After changing annotations in a service `router/`, from that service directory:

```bash
# Beans
cd apis/beans
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal

# Espresso
cd apis/espresso
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal
```

Never hand-edit `docs/docs.go`, `docs/swagger.json`, or `docs/swagger.yaml`. Run the command twice; the second run must have no diff.

## Public contract pointers

Do not duplicate the portal here. After a public behavior change, close the cascade in root `AGENTS.md`: router annotations → generated Swagger → `config/<product>.oas.json` → `docs/pages/` → Bruno.

Shared collection envelope (both products): `{ data, pagination, meta }` with `pagination.limit`, `pagination.num_results` (this page only), and `pagination.next_cursor`. Neither product serializes `pagination.cursor`. Empty collections are HTTP 200. Missing detail is HTTP 404. Errors are `{ "error": { "code", "message" } }`. REST except health uses Bearer at the gateway.

Product-specific notes:

- **Beans:** `meta.as_of` is required on collections. Feed routes (`/articles/latest`, `/articles/trending`, `/news/top-headlines`) reject `ids`, `urls`, `from`, and `to`. `/news/top-headlines` additionally rejects `content_type` and is fixed to news in a 24-hour window. Keep `GET /news/top-headlines` on the backend. `content_type=post` as a request filter returns 400; `post` may still appear on Article responses. Unknown or route-inapplicable query parameters return 400.
- **Espresso:** Formats via `response_type`: `json` (default), `yaml`, `toon` (same logical payload). Event/Signal stable core: `id`, `kind`, `created_at`, `tags`. Conditional: `summary`, `source`, `links`, `counts`. Ignore unknown extension fields. There is no public Actions route.

## Layout

| Path | Purpose |
|------|---------|
| `beans/main.go`, `espresso/main.go` | Process entry |
| `*/router/` | HTTP contract, annotations, handlers |
| `beans/db/`, `espresso/db/` | Persistence access |
| `*/docs/` | Generated Swagger |
| `*/tests/` | Contract, DB, and stress tests |
| `*/bruno/` | Maintained executable HTTP examples (OpenCollection) |
| `design/` | Internal design records (not public contract) |

Internal persistence and retrieval details stay in each service `db/` package and `apis/AGENTS.md`. They are not public documentation.
