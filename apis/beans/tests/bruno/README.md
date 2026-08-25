# Beans Bruno collection (maintainer)

Executable HTTP examples for the **local Beans process** (`baseUrl` defaults to `http://localhost:8080`). Gateway clients use `https://api.cafecito.tech` and the `/beans` prefix.

This collection is **not** the public contract. Align requests with `apis/beans/router` and `config/beans.oas.json`. Portal examples live under `docs/pages/products/beans/`.

## Variables

| Variable | Typical value | Notes |
| --- | --- | --- |
| `baseUrl` | `http://localhost:8080` | No `/beans` prefix on the backend. |
| `apiKey` | Bearer token | Required for every request except `GET /health` when the process has `API_KEY` set. Public clients send `Authorization: Bearer <key>`. |

Do not document or rely on private backend headers in public copy.

## Contract reminders

- Collection envelope: `data`, `pagination.limit`, `pagination.num_results`, `pagination.next_cursor`, `meta.as_of`.
- Errors: `{ "error": { "code", "message" } }`. Empty collections are HTTP 200. Missing detail is 404.
- `GET /top-headlines` maps to public `GET /beans/articles/top-headlines`. Do not send `content_type`, `ids`, `urls`, `from`, or `to`.
- Feed routes (`/articles/latest`, `/articles/trending`) reject `ids`, `urls`, `from`, and `to`.
- Request `content_type=post` is invalid. `post` may appear on Article responses.
- `full_content=true` requests body text when available; it is not a full-text guarantee.

Operation names in this folder match public `operationId` values (for example Search Articles → `searchArticles`, List Top Headlines → `getTopHeadlines`).
