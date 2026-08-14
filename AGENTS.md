# Cafecito API Platform
Updated: 2026-08-14

## System map

Project Cafecito is a monorepo: Zuplo gateway, Zudoku developer portal, and backend Go apis live in `cafecito-api-platform`.

- **Root** (`config/`, `modules/`, `docs/`): Zuplo gateway and Zudoku developer portal. Owns public API paths, API-key creation, Clerk integration, rate limits, quotas, OpenAPI specs, and docs.
- **`apis/beans/`**: Beans service — read-only news/blog aggregation API with article search, trends, source metadata, and propagation tracking.
- **`apis/espresso/`**: Espresso service — read-only business intelligence API over "sips" with events, signals, tags, related records, and token-efficient text responses for MCP/agent workflows.

Gateway paths use product prefixes such as `/beans/...` and `/espresso/...`. Backend Go apis expose their routes without those prefixes locally, typically on `:8080`.

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

## Documentation dependency map

The backend Swagger, gateway OpenAPI, and portal pages are separate artifacts. A change does not propagate between them automatically.

- Backend route behavior, request/response types, or Swagger annotations in `apis/<service>/router/` are the service-local contract. Regenerate that service's committed Swagger artifacts after changing annotations; never hand-edit generated `docs/docs.go`, `docs/swagger.json`, or `docs/swagger.yaml`.
  - Espresso: changing annotations in `apis/espresso/router/routes.go` (or referenced router types) requires, from `apis/espresso/`, `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal`.
  - Beans: changing annotations requires the documented `swag` command in `apis/beans/README.md`.
- A public Espresso contract change must also be reflected manually in `config/espresso.oas.json`. It is the gateway contract under `/espresso`, powers the portal reference at `/api/espresso` through `docs/zudoku.config.tsx`, and declares the `/espresso/mcp` server and its exported `operationId` tools. Keep gateway paths, request/response schemas, descriptions, error semantics, and MCP tool mappings consistent with the backend behavior.
- Apply the same review to `config/beans.oas.json` for Beans. `docs/zudoku.config.tsx` mounts it at `/api/beans` and mounts Espresso at `/api/espresso`.
- A public Espresso behavior change must be reflected in the relevant `docs/pages/` content:
  - product behavior, filters, response envelopes, and route examples: `docs/pages/products/espresso/overview.mdx`;
  - multi-call or agent route sequences: `docs/pages/products/espresso/workflows.mdx`;
  - product positioning and provider mappings: `docs/pages/products/espresso/migration.mdx`;
  - MCP endpoint, tool list, and agent workflow: `docs/pages/guides/mcp-ai-agents.mdx`;
  - shared authentication, pagination, or response-format behavior: `docs/pages/guides/api-conventions.mdx` and, when the quickstart changes, `docs/pages/start/first-api-call.mdx`.
- Update `docs/zudoku.config.tsx` whenever an API reference mount, documentation page, navigation item, or redirect changes. Update `docs/pages/api-overview.mdx` or `docs/pages/start/overview.mdx` only when product availability or high-level positioning changes.

### Public documentation boundary

Public surfaces are gateway OpenAPI files in `config/`, generated Swagger files that are served by a backend, and `docs/pages/`. Describe public behavior and stable contracts only. If internal implementation detail appears in one of these surfaces, remove it and replace it with user-facing behavior or a generic message such as "Service unavailable; retry."

Never expose:

- persistence architecture: Espresso's `cupboard`, `sips`, `sources`, or `relations` tables; relationship values such as `same_as` and `derived_from`; digest payloads and their structure; Beans' `beansack` and database tables; SQL, schemas, columns, foreign-key choices, indexes, materialized views, or migration details;
- retrieval implementation: embeddings, vector dimensions or indexes, HNSW, embedder/model configuration, caches, query internals, or relation direction/relationship storage;
- private code and operations: Go packages/types/handlers, internal identifiers, backend environment variables or headers, gateway rewrites/policies, infrastructure/vendor details, credentials, quotas implementation, and operational topology.

Public docs may describe Events, Signals, Sources, evidence, filters, pagination, response formats, and API-key requirements, but not how those concepts are stored or implemented.


CI/deploy:

- `.github/workflows/deploy-gateway.yml`: lint/test; `paths-ignore: apis/**`
- Zuplo deploys via GitHub integration — configure path filters to exclude `apis/**`


## API Implementation

Code: [apis/](apis/)
Read: [AGENTS.md](apis/AGENTS.md) for instructions


## First files to open by task

- Gateway route/auth/rate-limit issue: `config/policies.json`, then `modules/gate-auth.ts`, `modules/tiered-rate-limit.ts`, and the relevant `config/*.oas.json`.
- Developer portal/API-key issue: `docs/zudoku.config.tsx`, `modules/create-api-keys.ts`, `modules/clerk-webhook.ts`, `modules/consumer-ops.ts`.
- Beans endpoint behavior: `apis/beans/router/routes.go`, then `apis/beans/beansack/pgsack.go` and `apis/beans/beansack/types.go`.
- Beans API spec/docs drift: `apis/beans/router/`, `apis/beans/docs/swagger.yaml`, `config/beans.oas.json`, and `docs/pages/products/beans/overview.mdx`.
- Espresso endpoint behavior: `apis/espresso/router/routes.go`, `apis/espresso/router/types.go`, then `apis/espresso/cupboard/database.go`.
- Espresso API spec/docs drift: `apis/espresso/router/`, `apis/espresso/docs/swagger.yaml`, `config/espresso.oas.json`, `docs/pages/products/espresso/`, and the relevant `docs/pages/guides/` page.
