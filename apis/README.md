# Backend APIs (maintainer only)

**Status:** current  
**Audience:** service maintainers  
**Authority:** live contract is each service's `router/` bindings, handlers, and response structs; generated Swagger is produced from those annotations. Public docs live in the portal (`docs/pages/`) and gateway OpenAPI (`config/beans.oas.json`, `config/espresso.oas.json`).  
**Last verified:** 2026-08-25

This file is **not** the public API contract. Do not copy route tables, auth headers, or storage details from here into portal pages.

Public product docs: [Beans overview](../docs/pages/products/beans/overview.mdx), [Espresso overview](../docs/pages/products/espresso/overview.mdx). Gateway paths use `/beans` and `/espresso`. Local Go processes expose the same routes **without** those prefixes, typically on `:8080`.

## Local run

Required: `PG_CONNECTION_STRING`. Optional: `EMBEDDER_BASE_URL`, `EMBEDDER_API_KEY`, `EMBEDDER_MODEL`, `PORT`, `API_KEY`.

`API_KEY` is semicolon-separated `Header=Value` pairs for the **backend** process. Public clients authenticate with a Bearer API key at the gateway. Do not document preview backend headers on public surfaces. If `API_KEY` is unset, backend auth is disabled.

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

- **Beans:** `meta.as_of` is required on collections. Feed routes (`/articles/latest`, `/articles/trending`, `/news/latest`, `/news/trending`, `/news/top-headlines`) reject `ids`, `urls`, `from`, and `to`. `/news/latest`, `/news/trending`, and `/news/top-headlines` additionally reject `content_type` and are fixed to news. `/news/top-headlines` uses a 24-hour window. Keep the `/news/...` feed routes on the backend. `content_type=post` as a request filter returns 400; `post` may still appear on Article responses. Unknown or route-inapplicable query parameters return 400.
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

## Backend response cache

Beans and Espresso use `gin-contrib/cache` v1.4.4 directly for their public data
GET handlers: `CachePage` captures/replays responses and the library's persistence
stores handle storage and expiry. Authentication runs before caching. Health,
Swagger, and `/preview` are excluded. Existing three-argument `NewRouter` callers
remain uncached; the optional fourth argument is a `*shared.ResponseCache` wrapper.

| Environment variable | Default | Behavior |
|---|---|---|
| `RESPONSE_CACHE_ENABLED` | `true` | Set `false` to bypass store initialization and route caching |
| `RESPONSE_CACHE_PATH` | blank | `redis://` or `rediss://` connection string; absent/blank selects in-memory storage |
| `RESPONSE_CACHE_TTL` | `10m` | Positive Go duration such as `30s`, `10m`, or `1h`, applied to every cached route |

The shared initializer is `shared.NewResponseCache()`. Its wrapper holds
the selected store, path, and TTL; `Wrap(handler)` delegates to `CachePage` using
that initialized configuration. Missing/blank TTL uses ten minutes. Invalid,
zero, or negative TTL fails initialization when caching is enabled. Disabled
caching skips configuration validation. Configuration changes require restart.
Provision a separate Redis instance or database number for each service through
`RESPONSE_CACHE_PATH` Fly secrets. Connection pooling, serialization, expiry,
and database-wide `Flush` use the library defaults; do not share the configured
database with unrelated data. There is no custom Redis breaker, operation budget,
or runtime fallback to memory.

`RESPONSE_CACHE_DIR`, `RESPONSE_CACHE_DISK_MAX_BYTES`,
`RESPONSE_CACHE_DISK_MAX_ENTRIES`, and `RESPONSE_CACHE_MAX_ENTRY_BYTES` are removed
and ignored. Memory is local to each process/Machine, disappears on restart, and
has no application-enforced capacity limit. No volume or disk cleanup is needed.

Keys include the concrete path and raw query string via `RequestURI()`. Equivalent
parameter orders/defaults create different entries. Query values, cursors,
`full_content`, and `response_type` remain part of the key. Hits skip the handler,
including binding, embedding, database work, and serialization. Cached bodies
retain their original `meta.as_of`, pagination cursor, and format.

All cached routes use `RESPONSE_CACHE_TTL`, defaulting to ten minutes per process.
There is no TTL jitter or shared-fill coordination. The library caches responses
below HTTP 300, captures headers and bytes during writes, and removes entries
when the handler aborts. Its header replay and incremental/concurrent write
behavior are used unchanged. Strict completed-200-only, size, cookie/private,
and header-allowlist policies are intentionally removed.

Disable caching with `RESPONSE_CACHE_ENABLED=false` and restart/redeploy. Existing
Redis entries from the custom implementation use different keys and expire
naturally. Old disk cache files are no longer read or cleaned up by the APIs.

### Cache verification

Use an isolated Redis database for `RESPONSE_CACHE_TEST_REDIS`; never production.
The shared suite uses native memory by default. Service tests use their normal
Postgres fixtures and verify all public data routes still serve warmed responses
after database access is closed.

```bash
cd apis
RESPONSE_CACHE_TEST_REDIS=redis://127.0.0.1:6379/0 \
  GOCACHE=/tmp/cafecito-go-build go test -race ./shared/tests -run TestResponseCache -count=1
# Set PG_CONNECTION_STRING to the appropriate service's isolated fixture database:
GOCACHE=/tmp/cafecito-go-build go test ./beans/tests -run TestResponseCache -count=1
GOCACHE=/tmp/cafecito-go-build go test ./espresso/tests -run TestResponseCache -count=1
```

The service route suites are sequential: per-key coalescing and concurrent
response capture are now library behavior. See the
[design record](design/RESPONSE_CACHE_DESIGN.md) for the accepted tradeoffs.
Cache implementation details and credentials remain excluded from public docs.

Verified on 2026-10-16: builds, focused static checks, and race-enabled shared and
service cache tests pass. Broader suites retain the same 10 Beans/four Espresso
failures reproduced in the pre-replacement source snapshot. Verification used
isolated Redis/Postgres fixtures; Beans needed the existing database-only
`ideology` fixture compatibility overlay. Stress tests against the running local
service and the Espresso live-embedder router test were excluded. See the design
record's validation section for the precise limits.
