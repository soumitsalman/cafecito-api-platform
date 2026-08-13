# Cafecito API Platform
Updated: 2026-08-13

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
- Beans API spec/docs drift: `apis/beans/docs/swagger.yaml`, `config/beans.oas.json`, and `docs/pages/howtos/beans-howto.mdx`.
- Espresso endpoint behavior: `apis/espresso/router/routes.go`, `apis/espresso/router/types.go`, then `apis/espresso/cupboard/database.go`.
- Espresso API spec/docs drift: `apis/espresso/docs/swagger.yaml`, `config/espresso.oas.json`, and `docs/pages/howtos/espresso-howto.mdx`.
