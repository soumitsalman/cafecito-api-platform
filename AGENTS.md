# Cafecito API Manager Repository Knowledge Base

Last refreshed: 2026-06-25

## System map

Project Cafecito is split into a gateway/developer-portal repository and product API repositories.

- `cafecito-api-manager` is the Zuplo gateway and Zudoku developer portal. It owns public API paths, API-key creation, Clerk integration, rate limits, quotas, OpenAPI specs, and docs.
- `go-beans-api` is the Beans service: read-only news/blog aggregation API with article search, trends, source metadata, and propagation tracking.
- `go-espresso-api` is the Espresso service: read-only business intelligence API over "sips" with events, signals, tags, related records, and token-efficient text responses for MCP/agent workflows.

Gateway paths use product prefixes such as `/beans/...` and `/espresso/...`. Backend Go services expose their routes without those prefixes locally, typically on `:8080`.

## cafecito-api-manager

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

## go-beans-api

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

- Build: `go build -o beansapi .`
- Run: `make run` or `./beansapi`
- Docker: `docker compose up --build`
- Regenerate docs: `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g main.go -o docs`
- Tests: `go test ./tests/...` with a reachable database and `.env`.

## go-espresso-api

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

- Build: `go build -o espressoapi .`
- Run: `make run` or `./espressoapi`
- Docker: `docker compose up --build`
- Regenerate docs: `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal`
- Tests: `go test ./tests/...` with a reachable database and `.env`.

## Cross-repo implementation notes

- Public OpenAPI specs in `cafecito-api-manager/config` should stay aligned with generated Swagger specs in each Go service.
- Gateway paths add `/beans` or `/espresso`; backend service route files do not.
- The gateway authenticates public traffic and forwards a backend API key header. The Go services can also enforce `API_KEYS` directly.
- Both Go services use the same basic runtime pattern: env loading, DB pool, remote embedder, Gin router, CORS, optional API key middleware, and in-memory concurrency queue.
- Both Go services rely on external ingestion pipelines for database schema and data population; schema changes may require coordinating with `pycoffeemaker`.
- For MCP/agent ergonomics, Espresso explicitly supports `response_type=text`; Beans exposes MCP docs and JSON APIs but does not currently mirror Espresso's text response format in the backend routes.

## First files to open by task

- Gateway route/auth/rate-limit issue: `cafecito-api-manager/config/policies.json`, then `modules/gate-auth.ts`, `modules/tiered-rate-limit.ts`, and the relevant `config/*.oas.json`.
- Developer portal/API-key issue: `docs/zudoku.config.tsx`, `modules/create-api-keys.ts`, `modules/clerk-webhook.ts`, `modules/consumer-ops.ts`.
- Beans endpoint behavior: `go-beans-api/router/routes.go`, then `beansack/pgsack.go` and `beansack/types.go`.
- Beans API spec/docs drift: `go-beans-api/docs/swagger.yaml`, `cafecito-api-manager/config/beans.oas.json`, and `docs/pages/howtos/beans-howto.mdx`.
- Espresso endpoint behavior: `go-espresso-api/router/routes.go`, `router/types.go`, then `cupboard/database.go`.
- Espresso API spec/docs drift: `go-espresso-api/docs/swagger.yaml`, `cafecito-api-manager/config/espresso.oas.json`, and `docs/pages/howtos/espresso-howto.mdx`.

