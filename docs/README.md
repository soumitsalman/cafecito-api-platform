# Cafecito Developer Portal

Zudoku-powered developer portal for [Project Cafecito](https://cafecito.tech) APIs and MCP servers. Source lives in `docs/pages/`; navigation and OpenAPI mounting are configured in [`zudoku.config.tsx`](zudoku.config.tsx).

Hosted at [developer.cafecito.tech](https://developer.cafecito.tech). Built on [Zudoku](https://zudoku.dev) and published alongside the Zuplo gateway.

## Product status (canonical)

| Product | Status | Public surface |
| --- | --- | --- |
| **Beans** | Live | Beans is a publisher-content API for news, blogs, financial and earnings reports, litigation and lawsuits, official statements, research, technical documents, and related coverage context. |
| **Espresso** | Live | Espresso is a market, business, event, and news intelligence API for searching concrete Events, interpreting Signals, and tracing evidence and Sources. GDELT and Perigon are comparison patterns, not claims of route-for-route equivalence. |
| **Cortado** | Future | No public API |
| **Latte** | Future | No public API |

One API key authenticates every live product. There is no official SDK. Gateway: `https://api.cafecito.tech`. Limits: ongoing free tier, 100 requests/min and 50,000/month per user (see [Pricing and limits](pages/guides/pricing-limits.mdx)).

Approved one-liners are on [Getting started](pages/start/overview.mdx).

## Portal content summary

### Getting Started — `pages/start/`

| Page | File | Summary |
|------|------|---------|
| Overview | `start/overview.mdx` | Product chooser: Beans and Espresso live; Cortado and Latte future. |
| API Keys | `start/api-keys.mdx` | Create keys in the portal; one key for live REST and MCP. Bearer auth. Health is public. |
| First API call | `start/first-api-call.mdx` | First Beans and Espresso requests. |

### Products — `pages/products/`

| Page | File |
|------|------|
| Beans | `products/beans/overview.mdx`, `scenarios.mdx`, `migration.mdx` |
| Espresso | `products/espresso/overview.mdx`, `workflows.mdx`, `migration.mdx` |
| Cortado | `products/cortado/overview.mdx` (coming soon) |

### Guides — `pages/guides/`

| Page | File |
|------|------|
| MCP & AI agents | `guides/mcp-ai-agents.mdx` |
| Cross-product workflow | `guides/cross-product-workflow.mdx` |
| API conventions | `guides/api-conventions.mdx` |
| Reusable clients | `guides/client-patterns.mdx` |
| Troubleshooting | `guides/troubleshooting.mdx` |
| Pricing and limits | `guides/pricing-limits.mdx` |

### API Reference — `pages/api-overview.mdx`

Interactive OpenAPI reference for **Beans** (`/api/beans`) and **Espresso** (`/api/espresso`), mounted from `../config/*.oas.json`.

### Contact — `pages/contact.mdx`

Bug reports and feature requests via GitHub issue templates.

### Company & Policies

| Page | File |
|------|------|
| About Us | `company/about-us.md` |
| Privacy Policy | `company/privacy-policy.md` |
| Terms of Use | `company/terms-of-use.md` |

## Local development

From the repo root:

```bash
cd docs
npm install
npm run dev
```

Or from the gateway root: `npm run docs` (Zuplo docs integration).

Production build: `cd docs && npm run build`

## Clerk auth (API key creation)

If API key `createKey` is not invoked, the user is typically not authenticated in the portal.

Set these environment variables before starting docs:

- `ZUDOKU_PUBLIC_CLERK_PUBLISHABLE_KEY` (preferred)
- `ZUDOKU_PUBLIC_CLERK_JWT_TEMPLATE_NAME` (preferred, default: `dev-portal`)
- `CLERK_PUBLISHABLE_KEY` and `CLERK_JWT_TEMPLATE_NAME` are also supported as fallback names.
- `ZUDOKU_FAIL_ON_DEMO_CLERK_KEY` (default: `true`) fails fast in non-production if the demo key is still in use.

In Clerk, allow the Developer Portal callback URL:

- Local: `http://localhost:3000/oauth/callback`
- Hosted: `https://<your-docs-domain>/oauth/callback`

Ensure the JWT template name in Clerk matches `CLERK_JWT_TEMPLATE_NAME`.

## Learn more

- Backend Go services (run, tests, Swagger): [`../apis/README.md`](../apis/README.md)
- [Zuplo Developer Portal docs](https://zuplo.com/docs/dev-portal/introduction)
