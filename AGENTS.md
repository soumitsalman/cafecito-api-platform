# Cafecito API Manager Repository Knowledge Base

Last refreshed: 2026-07-01

## System map

Project Cafecito is a monorepo: Zuplo gateway, Zudoku developer portal, and backend Go services live in `cafecito-api-manager`.

- **Root** (`config/`, `modules/`, `docs/`): Zuplo gateway and Zudoku developer portal. Owns public API paths, API-key creation, Clerk integration, rate limits, quotas, OpenAPI specs, and docs.
- **`services/beans/`**: Beans service — read-only news/blog aggregation API with article search, trends, source metadata, and propagation tracking.
- **`services/espresso/`**: Espresso service — read-only business intelligence API over "sips" with events, signals, tags, related records, and token-efficient text responses for MCP/agent workflows.
- **`services/TEI-Dockerfile`** + root **`docker-compose.yml`**: shared local embedder and combined Docker stack.

Gateway paths use product prefixes such as `/beans/...` and `/espresso/...`. Backend Go services expose their routes without those prefixes locally, typically on `:8080`.

## Gateway (root)

Stack:

- TypeScript, Zuplo runtime, Zudoku docs portal
- Workspaces: root gateway plus `docs/`
- Main scripts: `npm run dev`, `npm run test`, `npm run docs`, `npm run lint`

Key files:

- `config/beans.oas.json`: Beans public OpenAPI routes under `/beans`
- `config/espresso.oas.json`: Espresso public OpenAPI routes under `/espresso`
- `config/developer.oas.json`: developer API for key creation and Clerk webhooks
- `config/policies.json`: Zuplo auth, quota, backend-key, and rate-limit policies
- `docs/zudoku.config.tsx`: docs navigation, OpenAPI mounting, Clerk auth, API-key creation flow
- `modules/create-api-keys.ts`: creates Zuplo key-bucket consumers with user metadata
- `modules/clerk-webhook.ts`: syncs/deletes consumers for Clerk user/subscription events
- `modules/consumer-ops.ts`: lists, patches, and deletes Zuplo consumers by `tags.sub`
- `modules/gate-auth.ts`: rejects requests without authenticated `request.user.sub`
- `modules/tiered-rate-limit.ts`: per-plan user rate limit logic

Auth/gateway behavior:

- API keys are created from the Zudoku portal through `POST /v1/developer/api-key`.
- Zudoku signs the create-key request using Clerk auth context.
- Created Zuplo consumers carry `tags.sub`, `tags.email`, and metadata such as `subscription_plan` and `subscription_status`.
- `gate-auth` currently requires an authenticated user but has inactive-subscription blocking commented out.
- Gateway injects backend auth with `X-API-KEY: $env(BACKEND_API_KEY)`.
- Published docs say one API key works across product APIs and MCP endpoints.

Docs/products:

- Live products in docs: Beans and Espresso.
- Future/reserved products: Cortado and Latte.
- MCP docs point to hosted endpoints such as `https://api.cafecito.tech/beans/mcp` and `https://api.cafecito.tech/espresso/mcp`.

CI/deploy:

- `.github/workflows/deploy-gateway.yml`: lint/test; `paths-ignore: services/**`
- Zuplo deploys via GitHub integration — configure path filters to exclude `services/**`

## services/beans

Stack:

- Go module `github.com/soumitsalman/beansapi`
- Gin, pgx, pgvector, zerolog, godotenv, swaggo, gRPC TEI embedder client
- Tests under `tests/`, OpenAPI generated under `docs/`

Core domain:

- A "bean" is an article/post keyed by canonical `url`.
- PostgreSQL schema is expected from the Beans ingestion pipeline in `pycoffeemaker/pybeansack/pgsack.sql`.
- Tables/entities include `beans`, `publishers`, `chatters`, `related_beans`, and `trend_aggregates`.
- Semantic search uses 384-dimensional pgvector embeddings.

Important files:

- `main.go`: loads env, creates `beansack.PGSack`, remote embedder, router, and server.
- `router/routes.go`: HTTP routes, validation, auth/concurrency middleware, swagger annotations.
- `beansack/types.go`: response/domain types such as `Bean`, `Publisher`, `Chatter`, `BeanTrend`, `PropagationResult`.
- `beansack/pgsack.go`: SQL builder and query execution.
- `nlp/embedder.go`: gRPC client to TEI-compatible embedding service.
- `docs/swagger.yaml` and `docs/swagger.json`: generated API spec.

Local/backend routes:

- `GET /health`
- `GET /tags/categories`
- `GET /tags/entities`
- `GET /tags/regions`
- `GET /sources`
- `GET /articles/search`
- `GET /articles/latest`
- `GET /articles/trending`
- `GET /articles/top-headlines`
- `GET /articles/propagation`
- `POST /articles/propagation`
- `GET /swagger/*any`

Common query concepts:

- Pagination: `limit` default 16, max 128; `offset` default 0.
- Article filters: `q`, `acc`, `content_type`, `urls`, `tags`, `categories`, `regions`, `entities`, `sources`, `from`, `full_content`.
- Propagation accepts up to 128 URLs and returns coverage plus social mentions per input URL.

Runtime config:

- Required: `PG_CONNECTION_STRING`, `EMBEDDER_BASE_URL`
- Optional: `EMBEDDER_API_KEY`, `EMBEDDER_MODEL`, `PORT`, `MAX_CONCURRENT_REQUESTS`, `API_KEYS`
- `API_KEYS` format is semicolon-separated `Header=Value`, for example `X-API-KEY=secret;Authorization=Bearer token`.
- If `API_KEYS` is unset, backend auth is disabled.

Commands:

- Build: `cd services/beans && go build -o beansapi .`
- Run: `cd services/beans && make run` or `./beansapi`
- Docker (from repo root): `docker compose up --build tei beansapi` (beans on `:8080`)
- Regenerate docs: `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g main.go -o docs`
- Tests: `go test ./tests/...` with a reachable database and `.env`.

Deploy: `.github/workflows/deploy-beans.yml` → Azure Container App `cafecito-beans-api` (paths: `services/beans/**`).

## services/espresso

Stack:

- Go module `github.com/soumitsalman/espressoapi`
- Gin, pgx, pgvector, zerolog, godotenv, swaggo, gRPC TEI embedder client
- Tests under `tests/`, OpenAPI generated under `docs/`

Core domain:

- A "sip" is a UUID-keyed unit of intelligence with kind `action`, `event`, or `signal`.
- The router flattens each sip `digest` JSON and merges `id` and `created` into response objects.
- PostgreSQL schema is expected from the Espresso ingestion pipeline in `pycoffeemaker/pycupboard/pgcupboard.py`.
- Tables/entities include `sips`, `sources`, and `relations`.
- Relations support `same_as` and `derived_from`.
- Semantic search uses 384-dimensional pgvector embeddings.

Important files:

- `main.go`: loads env, creates `cupboard.Cupboard`, remote embedder, router, and server.
- `router/routes.go`: HTTP routes, validation, response selection, auth/concurrency middleware, swagger annotations.
- `router/types.go`: flattened JSON response shapes and `response_type=text` rendering.
- `cupboard/types.go`: persistence types for `Sip`, `Source`, and `Relation`.
- `cupboard/database.go`: SQL builder and query execution.
- `nlp/embedder.go`: gRPC client to TEI-compatible embedding service.
- `docs/swagger.yaml` and `docs/swagger.json`: generated API spec.

Local/backend routes:

- `GET /health`
- `GET /tags`
- `GET /events`
- `GET /signals`
- `GET /related/:relationship`
- `GET /swagger/*any`

Common query concepts:

- Pagination: `limit` default 16, max 128; `offset` default 0.
- Event/signal filters: `ids`, `tags`, `q`, `acc`, `from`, `response_type`, `limit`, `offset`.
- `response_type=json` is default.
- `response_type=text` returns compact plain-text records for MCP/LLM context.

Runtime config:

- Required: `PG_CONNECTION_STRING`, `EMBEDDER_BASE_URL`
- Optional: `EMBEDDER_API_KEY`, `EMBEDDER_MODEL`, `PORT`, `MAX_CONCURRENT_REQUESTS`, `API_KEYS`
- `API_KEYS` format is semicolon-separated `Header=Value`.
- If `API_KEYS` is unset, backend auth is disabled.

Commands:

- Build: `cd services/espresso && go build -o espressoapi .`
- Run: `cd services/espresso && make run` or `./espressoapi`
- Docker (from repo root): `docker compose up --build tei espressoapi` (espresso on `:8081`)
- Regenerate docs: `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal`
- Tests: `go test ./tests/...` with a reachable database and `.env`.

Deploy: `.github/workflows/deploy-espresso.yml` → Azure Container App `cafecito-espresso-api` (paths: `services/espresso/**`).

## Local Docker stack

Root `docker-compose.yml` runs:

- `tei` — shared embedder from `services/TEI-Dockerfile`, port `10000`
- `beansapi` — port `8080`, env from `services/beans/.env`
- `espressoapi` — port `8081` (container `8080`), env from `services/espresso/.env`

Production Azure deploys use pre-built `docker.io/soumitsr/tei-static:latest` sidecar, not the local TEI Dockerfile.

## Cross-component implementation notes

- Public OpenAPI specs in `config/` should stay aligned with generated Swagger specs in `services/*/docs/`.
- Gateway paths add `/beans` or `/espresso`; backend service route files do not.
- The gateway authenticates public traffic and forwards a backend API key header. The Go services can also enforce `API_KEYS` directly.
- Both Go services use the same basic runtime pattern: env loading, DB pool, remote embedder, Gin router, CORS, optional API key middleware, and in-memory concurrency queue.
- Both Go services rely on external ingestion pipelines for database schema and data population; schema changes may require coordinating with `pycoffeemaker`.
- For MCP/agent ergonomics, Espresso explicitly supports `response_type=text`; Beans exposes MCP docs and JSON APIs but does not currently mirror Espresso's text response format in the backend routes.

## First files to open by task

- Gateway route/auth/rate-limit issue: `config/policies.json`, then `modules/gate-auth.ts`, `modules/tiered-rate-limit.ts`, and the relevant `config/*.oas.json`.
- Developer portal/API-key issue: `docs/zudoku.config.tsx`, `modules/create-api-keys.ts`, `modules/clerk-webhook.ts`, `modules/consumer-ops.ts`.
- Beans endpoint behavior: `services/beans/router/routes.go`, then `services/beans/beansack/pgsack.go` and `services/beans/beansack/types.go`.
- Beans API spec/docs drift: `services/beans/docs/swagger.yaml`, `config/beans.oas.json`, and `docs/pages/howtos/beans-howto.mdx`.
- Espresso endpoint behavior: `services/espresso/router/routes.go`, `services/espresso/router/types.go`, then `services/espresso/cupboard/database.go`.
- Espresso API spec/docs drift: `services/espresso/docs/swagger.yaml`, `config/espresso.oas.json`, and `docs/pages/howtos/espresso-howto.mdx`.
