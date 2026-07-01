# Cafecito Developer Portal

Zudoku-powered developer portal for [Project Cafecito](https://cafecito.tech) APIs and MCP servers. Source lives in `docs/pages/`; navigation and OpenAPI mounting are configured in [`zudoku.config.tsx`](zudoku.config.tsx).

Hosted at [developer.cafecito.tech](https://developer.cafecito.tech). Built on [Zudoku](https://zudoku.dev) and published alongside the Zuplo gateway.

## Portal content summary

### Welcome — `pages/introduction.mdx`

Landing page for the portal. Introduces live products (**Beans** news/blog intelligence, **Espresso** market intelligence) and **Cortado** (coming soon). Quick start: create an API key, call `https://api.cafecito.tech/<product>/<path>` with `Authorization: Bearer YOUR-API-KEY`.

### Getting Started

| Page | File | Summary |
|------|------|---------|
| API Keys | `howtos/api-keys.mdx` | Create keys in the portal; one key works across all products and MCPs. Bearer auth on every request. |
| Beans | `howtos/beans-howto.mdx` | News/blog aggregator (7,000+ sources). Semantic search, trends, propagation, tag filters. MCP tool order and REST examples. |
| Espresso | `howtos/espresso-howto.mdx` | Business intelligence **sips** (events, signals). UUID IDs, tag filtering, relationships (`same_as`, `derived_from`), `response_type=text` for agents. |
| Cortado | `howtos/cortado-howto.mdx` | Social media automation — placeholder; scheduling, posting, analytics planned. |

### MCP — `pages/howtos/mcp-howto.mdx`

Each product exposes an MCP server at `https://api.cafecito.tech/<product>/mcp` with the same API key. Documents Beans and Espresso tools mapped to REST paths.

### API Reference — `pages/api-overview.mdx`

Interactive OpenAPI reference for **Beans** (`/api/beans`) and **Espresso** (`/api/espresso`), mounted from `../config/*.oas.json`.

### Pricing — `pages/pricing.mdx`

Free launch preview until June 30, 2026. Beans: 100 req/min, 50k req/month. Metering is per API call (MCP counts the same). Unlimited keys share one meter.

### Contact — `pages/contact.mdx`

Bug reports and feature requests via GitHub issue templates on `cafecito-api-manager`.

### Company & Policies

| Page | File | Summary |
|------|------|---------|
| About Us | `company/about-us.md` | Founders, mission, product lineup (Beans live, Espresso, Cortado, MediCafe). |
| Privacy Policy | `company/privacy-policy.md` | Data collection, usage, cookies, retention (effective Sep 16, 2024). |
| Terms of Use | `company/terms-of-use.md` | Acceptable use, service terms for Beans/Espresso/Cortado (effective Sep 16, 2024). |

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

- [Zuplo Developer Portal docs](https://zuplo.com/docs/dev-portal/introduction)
- [Zudoku GitHub repository](https://github.com/zuplo/zudoku)
