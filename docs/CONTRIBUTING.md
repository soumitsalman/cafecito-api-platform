# Contributing to Cafecito developer documentation

**Audience:** maintainers and contributors editing portal copy, OpenAPI, or related indexes  
**Last verified:** 2026-08-25  
**Owner role:** documentation steward

This file covers **how to change documentation**. Public pages themselves live under `docs/pages/` and must stay free of internal implementation detail.

## Product status to keep consistent

- **Live:** Beans and Espresso REST plus MCP.
- **Future:** Cortado and Latte are reserved; they are not live APIs.
- **SDK:** there is **no official Cafecito SDK**. Document REST and MCP only.

## Where each fact is authoritative

| Fact | Edit here | Do not treat as authority |
|---|---|---|
| Public Beans contract | `config/beans.oas.json` (after aligning backend annotations) | `apis/design/` route proposals |
| Public Espresso contract | `config/espresso.oas.json` | `apis/design/` route proposals |
| Portal narrative, workflows, migration | `docs/pages/` | design plans |
| Navigation, mounts, redirects | `docs/zudoku.config.tsx` | this file |
| Service-local Swagger | regenerate from `apis/<service>/router/` using [`apis/README.md`](../apis/README.md) | hand-edits to `docs.go` / swagger JSON |

Cascade: router annotations → generated Swagger → gateway OAS → portal pages. Artifacts do not update each other automatically. See root [`AGENTS.md`](../AGENTS.md).

## Public documentation boundary

Public surfaces (`docs/pages/`, `config/*.oas.json`, backend Swagger served to clients) describe **stable client behavior** only.

Never publish:

- persistence architecture (table names, schemas, indexes, migrations, internal relationship storage names);
- retrieval internals (embeddings, vector dimensions, HNSW, embedder configuration, caches);
- private operations (Go package names, backend env vars, gateway rewrites/policies, credentials, quota implementation).

You may describe Events, Signals, Sources, evidence, filters, pagination, response formats, and API-key requirements without explaining how they are stored.

## Beans collection envelope (do not regress)

Collections use `{ data, pagination, meta }`:

- `pagination.limit` — effective page size
- `pagination.num_results` — records **in this page**, not a total match count
- `pagination.next_cursor` — opaque token or `null`; send it back as request `cursor`
- `meta.as_of` — UTC time this collection snapshot was produced

Empty collections: HTTP 200 with `data: []`. Older design docs that omit `num_results` or restrict `as_of` to a few routes are **historical**.

## Design records

Internal notes live in [`apis/design/`](../apis/design/). Read [`apis/design/README.md`](../apis/design/README.md) for current vs historical labeling. Do not copy superseded proposal paths or envelopes into the portal.

## Portal layout

Source: `docs/pages/`. Config: `zudoku.config.tsx`. Index of pages: [`README.md`](README.md).
