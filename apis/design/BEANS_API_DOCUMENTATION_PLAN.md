# Beans API Documentation Plan

Status: Documentation implementation plan

Updated: 2026-08-19

References:

- Target contract: [BEANS_API_ROUTE_PROPOSAL.md](BEANS_API_ROUTE_PROPOSAL.md)
- Parallel documentation model: [ESPRESSO_API_DOCUMENTATION_PLAN.md](ESPRESSO_API_DOCUMENTATION_PLAN.md)
- Industry comparison: [INDUSTRY_NEWS_API_ROUTE_REFERENCE.md](INDUSTRY_NEWS_API_ROUTE_REFERENCE.md)

## 1. Authority and V1 scope

The route proposal is the authority for the public Beans vocabulary. The public product is read-only publisher Article retrieval. All 15 stored content types are public `content_type` values:

`blog`, `contract`, `earnings_report`, `enforcement_action`, `financial_report`, `lawsuit`, `news`, `official_statement`, `podcast`, `post`, `press_release`, `research_paper`, `site`, `technical_documentation`, and `whitepaper`.

Omitting `content_type` includes all stored types. Provider labels such as `pr` are not Beans stored types or aliases.

| Area | V1 documentation rule |
|---|---|
| Product boundary | Beans finds and verifies what publishers published. Espresso covers business intelligence Events and Signals. |
| Article identity | Use `Article` and its UUID. Do not use storage names or imply durable Story membership. |
| Content type | Use any of the 15 stored types; omission means all stored types. |
| Pagination | `limit` is 1-100 with default 20. `cursor` is opaque and clients follow `pagination.next_cursor`. Do not document `offset` or `page`. |
| Empty results | Collections return HTTP 200 with `data: []`, pagination, and meta. |
| Dates | `from` and `to` use inclusive `YYYY-MM-DD` publication or observation bounds as described by the route. |
| Search | `q` is natural-language search. `score_threshold` requires `q`. `ids` and `urls` are exact search filters. |
| Projection | `full_content=true` requests the available Article body and may increase response size. |
| Errors | Use the standard error envelope and HTTP 400, 401, 404, 429, or 500 semantics. Do not describe success envelopes or HTTP 204 for empty collections. |
| Deferred routes | Stories remain gated until persistent Story UUID and membership behavior is available. Do not publish count or propagation routes. |

## 2. LLM-friendly operation strings

Every public operation description should use the same labeled sequence used by the Espresso plan:

1. `Answers:` state the user or agent question the operation resolves.
2. `Returns:` identify the resource, envelope, ordering, and empty-result behavior.
3. `Use when:` state the correct selection condition.
4. `Do not use when:` name nearby operations that are better matches.
5. `Search/filter behavior:` define accepted filters, exactness, OR or AND behavior, and rejected filters.
6. `Time:` define the date field, window, defaults, and input format.
7. `Sort/pagination:` define ordering and cursor behavior.
8. `Missing fields:` identify nullable or omitted data without exposing storage details.
9. `Next step:` give the next useful operation or workflow transition.

Example:

```text
Answers: Which publisher Articles match this topic or filter set?
Returns: A cursor-paginated collection of Article records. Empty results are HTTP 200 with data: [].
Use when: semantic search, exact IDs or URLs, or combined Article filters are needed.
Do not use when: newest, headline-ranked, or trend-ranked results are needed.
```

Descriptions must use the same resource names and parameter meanings across Go Swagger, gateway OpenAPI, portal pages, and MCP guidance. They must not expose database tables, embeddings, vector indexes, internal headers, or gateway rewrites.

## 3. Documentation delivery sequence

1. `apis/beans/router/routes.go`:
   - Replace stale operation summaries and descriptions with the labeled LLM contract above.
   - Use `Article`, `Source`, `Mention`, and discovery labels consistently.
   - Document canonical paths: `/articles/search`, `/articles/latest`, `/articles/top-headlines`, `/articles/trending`, `/articles/{id}`, `/articles/{id}/similar`, `/articles/{id}/mentions`, `/sources`, and `/sources/{id}`.
   - Document discovery paths `/categories`, `/entities`, `/regions`, and `/sentiments`.
   - Keep Stories and other deferred routes out of the live public route set.

2. `apis/beans/router/params.go` and `apis/beans/router/responses.go`:
   - Enforce and describe the public `content_type` enum as all 15 stored types listed in the proposal.
   - Keep cursor bounds at default 20 and maximum 100.
   - Keep route-specific filters aligned with the proposal: Article filters, date bounds, exact IDs or URLs, Source filters, and mention platform or forum filters.
   - Keep collection and detail envelopes consistent with `data`, `pagination`, and `meta`.

3. `apis/beans/docs/swagger.yaml`, `apis/beans/docs/swagger.json`, and `apis/beans/docs/docs.go`:
   - Regenerate from the revised annotations with Swaggo.
   - Verify operation IDs, paths, response schemas, public enums, and error responses.
   - Never hand-edit generated files.

   Regeneration command:

   ```bash
   cd apis/beans
   go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g main.go -o docs
   ```

4. `config/beans.oas.json`:
   - Mirror the generated public operations under `/beans`.
   - Keep gateway operation IDs, parameters, descriptions, schemas, and error behavior parallel with backend Swagger.
   - Set the `content_type` parameter and Article response enums to all 15 stored types; do not add provider-only `pr`.
   - Keep `/beans/mcp` and its tool mappings aligned with the live operation inventory.
   - Preserve existing authentication and gateway policy behavior.

5. `docs/pages/products/beans/overview.mdx`:
   - Make this the primary human and AI-agent route guide.
   - Explain the Beans and Espresso product boundary, Article and Source concepts, route selection, filters, cursor pagination, response envelopes, and missing-field behavior.
   - Include one canonical JSON collection example and one detail workflow.
   - State the complete stored-type `content_type` rule in the main content-type section.

6. `docs/pages/products/beans/migration.mdx`:
   - Map common provider concepts to Beans routes and parameters.
   - Compare search, latest, headlines, trending, similar, mentions, Sources, and discovery.
   - Mark provider concepts without a Beans equivalent as non-equivalents instead of inventing aliases.

7. `docs/pages/guides/mcp-ai-agents.mdx`:
   - List Beans MCP tools using the same operation IDs as the gateway contract.
   - Give an agent workflow: discover filters, search or select an Article, retrieve detail, then follow similar or mentions.
   - Explain that `content_type` accepts every stored type listed in the contract and that omission includes all stored types.

8. `docs/pages/guides/api-conventions.mdx`, `docs/pages/start/first-api-call.mdx`, and `docs/pages/start/api-keys.mdx`:
   - Align authentication, cursor pagination, empty collections, error handling, and first-call examples with the Beans and Espresso contracts.
   - Do not use legacy offset, page, propagation, or stored-kind examples.

9. `docs/pages/start/overview.mdx` and `docs/zudoku.config.tsx`:
   - Keep product positioning current and add the Beans migration page to product navigation.
   - Do not expose internal route or implementation details.

## 4. Provider comparison matrix

| User need | Beans route | Provider comparison guidance |
|---|---|---|
| Broad article search | `GET /beans/articles/search` | Translate provider full-text or semantic search to `q`; map exact IDs and URLs to `ids` and `urls`. |
| Recent publication feed | `GET /beans/articles/latest` | Map freshness or recency feeds to `from`, `to`, and cursor pagination. |
| Headline attention | `GET /beans/articles/top-headlines` | Use the fixed recent headline window; do not add provider time parameters that the route rejects. |
| Attention trend | `GET /beans/articles/trending` | Map ranking or engagement feeds to the trend route and its date window. |
| Related reading | `GET /beans/articles/{id}/similar` | Treat results as related reading, not guaranteed Story membership. |
| External observations | `GET /beans/articles/{id}/mentions` | Treat results as observations on external platforms, not publisher Articles. |
| Publisher metadata | `GET /beans/sources` and `GET /beans/sources/{id}` | Map publisher or source lookup to Source metadata and UUID filters. |
| Filter discovery | `GET /beans/categories`, `/entities`, `/regions`, `/sentiments` | Discover normalized values before applying Article filters. |

## 5. Cross-surface validation

- Compare backend Swagger paths and operation IDs with `config/beans.oas.json` after every annotation change.
- Verify that every public `content_type` enum contains exactly the 15 stored types from the proposal.
- Search examples and prose for missing stored kinds, provider `pr` aliases, `offset`, `page`, propagation, and internal implementation vocabulary.
- Validate every portal example against the gateway path, query parameter names, response envelope, and status behavior.
- Confirm MCP tool names and route descriptions match the gateway operation IDs.
- Keep deferred Stories, counts, propagation, and internal data concepts absent from public reference and workflows.
