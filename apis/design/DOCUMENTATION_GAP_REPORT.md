| Field | Value |
|---|---|
| Status | **current** |
| Authority | Internal audit of documentation surfaces; **not** the public API contract |
| Audience | Maintainers planning documentation and contract cleanup |
| Last verified | 2026-08-25 |
| Owner role | Documentation steward |
| Superseded by | n/a |

**Index:** [README.md](README.md)

Settled **published** Beans collection envelope (do not treat older proposal text as live): `{ data, pagination.limit, pagination.num_results, pagination.next_cursor, meta.as_of }`. Findings below that still argue for omitting `num_results` or restricting `meta.as_of` are historical relative to gateway OpenAPI.

**Scope:** Repository-wide static audit of the Cafecito documentation system.

This report identifies information that is missing, contradictory, stale, or difficult to discover. It covers the public Zudoku portal, gateway OpenAPI files, backend router contracts, generated Swagger artifacts, service READMEs, design plans, MCP definitions, and executable examples.

This is an internal planning report. It is not itself a public API contract.

## How to use this report

Each finding has four parts:

- **Gap:** what a user, maintainer, search engine, generated client, or AI agent cannot reliably determine.
- **Missing or incorrect in:** the files where the information is absent, stale, or self-conflicting.
- **Edit or add:** the files that should change to resolve the gap.
- **Completion check:** the evidence that should exist after the work is complete.

The required contract cascade is:

```text
router bindings / handlers / response structs
  -> router annotations and doc comments
  -> generated backend Swagger
  -> gateway OpenAPI
  -> portal pages and examples
  -> MCP catalog, Markdown output, search, and AI indexes
```

Do not hand-edit generated backend Swagger. Change the router source and regenerate `apis/beans/docs/*` or `apis/espresso/docs/*`.

Priority levels:

- **P0:** Information can cause an invalid request, broken generated client, or materially misleading integration.
- **P1:** Information is needed for a complete and dependable developer experience.
- **P2:** Information improves search, AI retrieval, or long-term maintainability.

## Executive summary

The public portal has strong narrative and workflow coverage. The main weakness is that the repository contains several competing contracts:

1. Current router bindings and response structs.
2. Generated backend Swagger.
3. Manually maintained gateway OpenAPI.
4. Current portal pages.
5. Older service READMEs and design plans.

The highest-priority work is to reconcile those surfaces before adding more prose. Otherwise new pages will continue to amplify conflicting behavior.

| Category | Main problem | Priority |
|---|---|---|
| Inconsistent information | Runtime, OpenAPI, examples, READMEs, and MCP documentation disagree on parameters, pagination, formats, status codes, and tool coverage. | P0 |
| Needed for usability | Common conventions, troubleshooting, client setup, dynamic response fields, and cross-product recipes are incomplete or distributed across multiple pages. | P1 |
| Needed for discovery | SEO metadata is broad, machine-readable schemas are weak, MCP/LLM catalogs are incomplete, and repository-facing documentation contains stale paths and terminology. | P1/P2 |


## Current-state re-audit - 2026-08-25

This section supersedes the open/closed interpretation of the original findings below. The detailed findings remain as audit evidence; an item is still open only when this disposition marks it active or partial.

### Verification snapshot

- **Contract cascade: PASS.** npm run verify:api-contracts reports Beans 17 backend / 18 gateway operations and Espresso 14 backend / 15 gateway operations, with zero unexplained mismatches. Reviewed differences are in config/api-contract-exceptions.json.
- **Portal build: PASS.** npm run build --workspace docs produces 65 prerendered routes, Pagefind, sitemap, published Markdown, llms.txt, and llms-full.txt.
- **Live documentation checks: PASS.** Examples, terms, lifecycle, links, inventory, generated-output, and forbidden-public-term checks produce no live issues.
- **Complete docs verifier: PASS.** `npm run verify:docs` passes live checks, positioning parity, generated-output checks, and all negative fixtures.
- **Deployment enforcement: PARTIAL.** api-contract-docs.yml is workflow_dispatch only; deploy workflows do not depend on it. The check is reproducible but not an automatic or required gate.

### Updated intent distance

| API | Current distance | Evidence | Remaining gap |
|---|---|---|---|
| Beans | Near | Publisher-content positioning, broad content-purpose examples, and migration mappings for TheNewsAPI, World News API, GNews, finlight, NewsAPI.ai, and NewsData.io are present; the start-page Markdown table, entry pages, OpenAPI, portal metadata, and generated Markdown are checked for wording parity. | Keep the parity expectations current when wording changes. |
| Espresso | Near | Market/business intelligence positioning, GDELT and Perigon mapping, evidence/provenance, monitoring, research-brief, and early-warning workflows are present; the start-page Markdown table, OpenAPI, portal metadata, and generated Markdown are checked for wording parity. | Keep the parity expectations current when wording changes. |

### Disposition by category

#### Inconsistent information

| Finding | Status | Current state |
|---|---|---|
| I-01 | Partial / controlled | Active READMEs and indexes are current. Superseded plans still contain old paths, but lifecycle labels and generated-output exclusion control the risk. |
| I-02, I-03 | Resolved | Beans post is response-only, top-headlines rejects inapplicable filters, and current router/OAS/portal/Bruno surfaces agree. |
| I-04, I-05 | Resolved on active surfaces | Current pagination and Espresso terminology/formats agree. Old cursor, sip, relation, and text claims remain only in labeled historical design records. |
| I-06 | **Resolved** | Active docs now state that unknown or route-inapplicable query parameters are rejected with HTTP 400 and the ErrorResponse envelope across Beans and Espresso; unsupported Espresso `response_type` is also documented as 400. |
| I-07, I-08 | Resolved | Bearer, health, 401/429 boundaries, and current MCP tool inventories are represented and structurally checked. |
| I-09, I-10 | Resolved / controlled | Product status, frozen limits, and public/internal boundaries are current; generated public output contains none of the tested internal terms. |

#### Needed for usability

| Finding | Status | Current state |
|---|---|---|
| U-01 | **Resolved** | Shared conventions now document strict unknown/inapplicable query rejection with the ErrorResponse envelope. |
| U-02, U-03, U-05, U-06, U-07, U-09, U-10, U-11 | Resolved or monitor | Schemas, Espresso route matrices, reusable clients, MCP setup, dynamic-field guidance, freshness semantics, SDK/versioning/limits guidance, and Bruno links are now present. |
| U-04 | **Resolved** | Troubleshooting now identifies unsupported `response_type` as HTTP 400 and directs clients to JSON. |
| U-08 | **Resolved** | The workflow has corrected numbering plus complete Node.js and Python composite clients covering Beans search, Espresso Events/Signals, and Event evidence. |

#### Needed for discovery

| Finding | Status | Current state |
|---|---|---|
| D-01, D-03, D-04, D-06, D-08, D-09 | Resolved or monitor | Product-specific metadata, intent-to-route API overview, MCP catalog, page descriptions, Bruno links, and search vocabulary are present. |
| D-02 | Resolved for crawl baseline; scoped for enhancement | Canonical URLs, sitemap, Pagefind, Markdown, LLM indexes, `metadata.robots`, and checked-in `robots.txt` are present. A custom global Open Graph/Twitter/JSON-LD plugin is intentionally not used; product-specific `<Head>` metadata remains an optional future enhancement. |
| D-05, D-07 | Partial / controlled | Structural OAS parity and lifecycle headers are present, but semantic prose parity is not machine-checked and repository crawlers can still encounter historical plans. |
| D-10 | Resolved | Generated output contains required current slugs and no tested internal terms; `npm run verify:docs` passes its generated-output checks and all negative fixtures. |

### Current remaining edit set

No active edits remain from this re-audit. `docs/zudoku.config.tsx` uses Zudoku native metadata, canonical URL, and sitemap configuration, with `docs/public/robots.txt` for crawler directives. If product-specific sharing cards or structured data become necessary, add small local `<Head>` blocks to the relevant product overview pages rather than restoring a global plugin.

### Current RCCA status

The architectural correction is underway, not complete. CA-1, CA-2, CA-4, CA-5, CA-6, and CA-7 have concrete repository artifacts. CA-3 has a manual workflow, but it is deliberately not a pull-request, push, or required deployment gate.

The remaining root cause is narrower: Cafecito now has structural reconciliation and documentation hygiene controls, but semantic prose parity and the final verification gate remain manually enforced.

# Product intent and documentation distance

## Assessment method

The two cited design references establish product intent and the vocabulary users will bring from comparable APIs. They are comparison authorities, not a requirement for Cafecito to clone every provider route.

This section measures documentation distance from the stated product intent, not implementation quality:

- Near: the public docs name the product category, explain the main user jobs, identify the core resources, and show the next route.
- Moderate: the core resources are documented, but the category, competitive frame, boundaries, or important use cases are not obvious.
- Far: the current docs describe a different product, omit the main resources, or make the intended use case difficult to recognize.

## Intent-distance summary

| API | Stated intent | Current documentation alignment | Distance | Main remaining gap |
|---|---|---|---|---|
| Beans | A multi-source News and Blogs API comparable to TheNewsAPI, World News API, GNews, finlight, NewsAPI.ai/Event Registry, and NewsData.io, with additional publisher content such as earnings reports, financial reports, litigation/lawsuits, official statements, research papers, podcasts, and technical documents. | Strong on core retrieval: Articles, Sources, Stories, feeds, discovery, related reading, mentions, filters, and full-content projection are documented. | Near | Keep the canonical source and parity check current when wording changes. |
| Espresso | A market and Business Intelligence API comparable to GDELT and Perigon, covering event/news intelligence, normalized events, stories or evidence, entities, geography, search, source context, and explainable intelligence. | Strong on the Cafecito model: Events, Signals, Evidence, Sources, discovery, structured filters, and agent workflows are documented. | Near | Keep the canonical source and parity check current when wording changes. |

### Overall conclusion

Beans and Espresso are now near their stated intents on active documentation surfaces. The canonical positioning source, gateway OpenAPI descriptions, portal pages, SEO metadata, generated Markdown, and parity verifier keep the category and provider-boundary language aligned. Remaining work is maintenance when product wording changes.

The highest-value additions are:

1. A clear category statement for each product.
2. A provider-pattern mapping that uses the named comparison APIs.
3. Explicit statements about equivalent, different, and unavailable capabilities.
4. Use-case examples that start from market or publisher questions rather than from Cafecito resource names alone.
5. A cross-product evidence workflow: Beans publisher material to Espresso intelligence and evidence.

## Beans distance from intent

### What the current documentation already does well

The Beans pages already document most of the route families expected from a News API:

- Article search and exact lookup.
- Latest, top-headlines, and trending feeds.
- Publisher and Source discovery.
- Story grouping and member Articles.
- Related publisher reading.
- External Article mentions.
- Content-type, publisher, date, label, and relevance filters.
- Compact collections followed by targeted detail calls.
- MCP access and agent-oriented pagination.

The Beans overview also has a provider-parallel route table, and the migration page lists a broad content taxonomy including earnings_report, financial_report, lawsuit, official_statement, research_paper, and technical_documentation.

### Where the current docs are distant from the intent

#### B-I1. The top-level category is too narrow

The portal calls Beans News & Blogs and Beans News API & MCP. That is understandable, but it hides the intended broader positioning as a publisher-content API that also covers financial, legal, corporate, research, technical, and other public content.

**Distance:** Low to moderate.

**Why it matters:** A user searching for an earnings reports API, lawsuit monitoring API, corporate statements API, or research-content API may not identify Beans as relevant.

**Files to edit:**

- docs/pages/start/overview.mdx
- docs/pages/api-overview.mdx
- docs/pages/products/beans/overview.mdx
- docs/pages/products/beans/migration.mdx
- docs/zudoku.config.tsx

**Change needed:**

Use a consistent phrase such as:

> Beans is a read-only publisher-content API for news, blogs, financial and earnings reports, litigation and lawsuits, official statements, research, technical documents, and related coverage context.

Keep News API in the title for search familiarity, but add publisher content and named content examples in the description and first screen.

#### B-I2. The broader content taxonomy is not discoverable as a set of user jobs

The type list exists, but it is presented mainly as a parameter enum. It does not show why a user would select earnings_report, lawsuit, official_statement, or research_paper.

**Distance:** Moderate for non-news use cases.

**Files to edit:**

- docs/pages/products/beans/overview.mdx
- docs/pages/products/beans/scenarios.mdx
- docs/pages/products/beans/migration.mdx
- config/beans.oas.json for enum descriptions and examples

**Change needed:**

Add a public content-purpose table:

| User need | Beans content types | Starting route |
|---|---|---|
| General news and blogs | news, blog | /beans/articles/search |
| Company financial updates | earnings_report, financial_report | Article search with content_type |
| Legal and regulatory monitoring | lawsuit, enforcement_action, contract | Article search with type and entity filters |
| Corporate communications | official_statement, press_release | Article search with type, Source, and date |
| Research and technical monitoring | research_paper, technical_documentation, whitepaper | Article search with type and topic filters |
| Audio and publication discovery | podcast, site | Article search with type and Source filters |

Add one short runnable example for earnings monitoring and one for litigation or official-statement monitoring.

#### B-I3. Provider migration language is not prominent enough

The cited market report provides direct comparisons to TheNewsAPI, World News API, GNews, finlight, NewsAPI.ai/Event Registry, and NewsData.io. The current Beans docs include provider parallels, but a user migrating from one of those APIs does not get a clear provider-to-Beans route and parameter map.

**Distance:** Moderate for competitive discovery and migration.

**Files to edit:**

- docs/pages/products/beans/migration.mdx
- docs/pages/products/beans/overview.mdx
- docs/pages/api-overview.mdx
- docs/pages/start/overview.mdx

**Change needed:**

Add a provider migration matrix:

| Coming from | Typical task | Beans replacement |
|---|---|---|
| TheNewsAPI | Article cards, source catalog, related articles | Article search, Sources, Similar Articles |
| World News API | Search, top news, entity and source filters | Article search, top-headlines, discovery, Source filters |
| GNews | Search, latest, ranked headlines | Article search, latest, top-headlines |
| finlight | Financial news and company-related monitoring | Article search with financial types, entities, Sources |
| NewsAPI.ai/Event Registry | News/blog selection, concepts, categories, duplicates, analysis | Article search, content types, labels, Stories, related coverage |
| NewsData.io | Search, continuation, AI tags, sentiment, full content where available | Article search, cursor pagination, normalized labels, full_content |

The matrix must describe behavioral differences rather than claiming identical fields or content rights.

#### B-I4. Content availability and full-content expectations need clearer positioning

The market report emphasizes that providers differ substantially in article-body availability, licensing, and truncation. Beans documents full_content=true, but the product positioning should make clear that this requests available content and is not a universal full-text guarantee.

**Distance:** Moderate for users migrating from full-content providers.

**Files to edit:**

- docs/pages/products/beans/overview.mdx
- docs/pages/products/beans/migration.mdx
- docs/pages/guides/api-conventions.mdx
- config/beans.oas.json

**Change needed:**

Add a clearly visible note:

> full_content=true requests body content when Beans has it available. It does not guarantee a full publisher copy for every Article. Applications should handle a missing or partial content value and use the canonical url for attribution.

Add this to the response schema description and a migration warning.

#### B-I5. The Beans boundary with Espresso is good but should be expressed in provider terms

The current docs say to use Espresso for Events, Signals, impact, forecast, and evidence relationships. That is useful. It should be expanded to explain that Beans is the publisher-material layer, while Espresso is the structured intelligence layer.

**Distance:** Low.

**Files to edit:**

- docs/pages/products/beans/overview.mdx
- docs/pages/products/espresso/overview.mdx
- Add docs/pages/guides/cross-product-workflows.mdx
- docs/pages/api-overview.mdx

**Change needed:**

Add a two-column decision table:

| If the user asks | Use |
|---|---|
| What did publishers report? | Beans |
| Which sources covered the same story? | Beans Stories and Sources |
| What external discussion exists around an Article? | Beans Mentions |
| What happened in structured intelligence terms? | Espresso Events |
| What does it mean or what is the outlook? | Espresso Signals |
| What evidence supports the conclusion? | Espresso Evidence |

## Espresso distance from intent

### What the current documentation already does well

Espresso already presents a coherent Cafecito-specific model:

- Events represent concrete developments.
- Signals represent synthesized conclusions.
- Evidence and Sources provide context and provenance.
- Discovery routes help resolve structured filter values.
- Workflows show event-to-evidence, event-to-signal, and signal-to-event traversal.
- JSON, YAML, and TOON are positioned for REST, MCP, and agent use.
- Route descriptions explain intent-first retrieval and selective enrichment.

This is a good foundation, but it is not yet enough to communicate the intended market and event intelligence category.

### Where the current docs are distant from the intent

#### E-I1. The market and event intelligence category is under-named

The current docs say business intelligence and Events and Signals, but they do not consistently use the broader terms a GDELT or Perigon user will search for:

- market intelligence API
- business intelligence API
- event intelligence API
- news intelligence API
- event monitoring
- company or entity monitoring
- impact and outlook
- evidence-backed intelligence

**Distance:** Moderate.

**Files to edit:**

- docs/pages/start/overview.mdx
- docs/pages/api-overview.mdx
- docs/pages/products/espresso/overview.mdx
- docs/pages/products/espresso/migration.mdx
- docs/zudoku.config.tsx

**Change needed:**

Use a consistent category statement:

> Espresso is a market and business intelligence API for searching concrete Events, interpreting Signals, and tracing evidence and Sources for monitoring, research, risk analysis, and AI-agent workflows.

Keep Events and Signals as the Cafecito resource names, but introduce the category terms before the resource model.

#### E-I2. The docs do not map GDELT and Perigon concepts to Espresso concepts

The cited industry reference shows the user expectations created by GDELT and Perigon:

- Event and story search.
- Event and story detail.
- Article or evidence traversal.
- Entity and geography discovery.
- Event classification and time filters.
- Source context.
- Aggregation and counts.
- Semantic or free-text search.
- Stable pagination and normalized response fields.

Espresso documents some equivalents, but it does not tell a user which concept is equivalent, which is Cafecito-specific, and which is not currently exposed.

**Distance:** Moderate to high for migration and comparison users.

**Files to edit:**

- docs/pages/products/espresso/migration.mdx
- docs/pages/products/espresso/overview.mdx
- docs/pages/products/espresso/workflows.mdx
- docs/pages/api-overview.mdx
- config/espresso.oas.json

**Change needed:**

Add an explicit concept map:

| GDELT or Perigon pattern | Espresso interpretation |
|---|---|
| Event search | GET /espresso/events |
| Event detail | GET /espresso/events/{event_id} |
| Story or article evidence | Event Evidence and Source routes |
| Signal or semantic interpretation | GET /espresso/signals |
| Supporting event traversal | GET /espresso/signals/{signal_id}/events |
| Entity, region, and event-type discovery | Espresso discovery routes |
| Source or publisher context | Espresso Sources |
| Free-text or semantic search | q on supported collection routes |
| Cursor continuation | cursor and next_cursor |
| Provider-specific story/count/vector routes | State whether no direct Espresso equivalent exists |

The last row is important. The docs must label intentionally unsupported provider concepts rather than implying that Espresso is a one-for-one GDELT or Perigon replacement.

#### E-I3. Signals are a differentiated product capability but are not explained as market intelligence

Signals are the strongest Espresso differentiator, but the current documentation describes them mainly as a resource type. It does not give enough market-facing examples of what a Signal means:

- impact or outlook for a company or sector
- policy or regulatory implication
- market movement or risk
- cross-source development
- emerging trend or change in direction

**Distance:** Moderate.

**Files to edit:**

- docs/pages/products/espresso/overview.mdx
- docs/pages/products/espresso/workflows.mdx
- docs/pages/products/espresso/migration.mdx
- docs/pages/guides/mcp-ai-agents.mdx

**Change needed:**

Add a Signals use-case table and examples:

| User question | Espresso path |
|---|---|
| What happened? | Search Events |
| What is the likely business or market implication? | Search Signals |
| What evidence supports that implication? | Retrieve Signal-linked Events and Event Evidence |
| Which companies, people, products, or regions are affected? | Use structured filters and discovery |
| Is this a one-off event or a developing pattern? | Compare related Events, Signals, Sources, and freshness metadata |

Use representative public payloads that clearly distinguish a concrete Event from a synthesized Signal.

#### E-I4. Market monitoring workflows are underrepresented

The current workflows focus on retrieval and traversal. They do not yet show the applications implied by the industry intent:

- monitoring a company, sector, region, or topic
- building an event-driven research brief
- producing an evidence-backed market update
- detecting an emerging risk or opportunity
- following a Signal back to supporting Events and Sources

**Distance:** Moderate.

**Files to edit or add:**

- docs/pages/products/espresso/workflows.mdx
- docs/pages/products/espresso/overview.mdx
- Add docs/pages/guides/market-intelligence-workflows.mdx if the workflow set becomes too large.
- Add the cross-product workflow in docs/pages/guides/cross-product-workflows.mdx.

**Change needed:**

Add at least three workflows:

1. Company or sector monitor: discover filter values, search Events, continue pagination, enrich only new or high-impact records.
2. Evidence-backed market brief: search Signals, retrieve supporting Events, inspect Evidence and Sources, present as_of.
3. Early-warning workflow: search Events for a bounded topic or region, compare Signals, retain IDs and cursors, and alert only after evidence is available.

#### E-I5. The relationship between raw publisher material, Events, and Signals is not visible enough

The internal capability analysis describes a useful evidence chain from Beans articles and propagation to Espresso intelligence. The public docs currently separate the products but do not show how an application can move between them.

**Distance:** Moderate for composite applications.

**Files to edit or add:**

- Add docs/pages/guides/cross-product-workflows.mdx.
- Update docs/pages/products/beans/overview.mdx.
- Update docs/pages/products/espresso/overview.mdx.
- Update docs/pages/guides/mcp-ai-agents.mdx.
- Update docs/pages/api-overview.mdx.

**Change needed:**

Show one canonical sequence:

    Beans Article or Story discovery
      -> select topic, entity, Source, or publisher evidence
      -> Espresso Event search
      -> Espresso Signal search or detail
      -> supporting Events and Evidence
      -> cite publisher and intelligence provenance separately

Explain that Beans answers what publisher material exists, while Espresso answers what the structured intelligence means and what supports it.

## Intent-specific changes grouped by report category

### Inconsistent information

Add intent consistency checks to the existing I findings:

- Beans naming: product titles and metadata must say News API plus broader publisher content. Do not let News & Blogs be the only category label.
- Beans content types: the same broad taxonomy must be present in the request binding, response schema, gateway OpenAPI, migration table, route examples, and top-level product description.
- Espresso naming: replace stale sips, actions, and response_type=text language in repository documentation, while using market/business/event intelligence in current public descriptions.
- Espresso comparison boundaries: add a clear statement to the migration page that GDELT and Perigon are comparison patterns, not claims of route-for-route equivalence.
- Cross-product boundary: keep Beans publisher-material language separate from Espresso structured-intelligence language. Do not describe Beans as producing Signals or Espresso as a general news archive.
- Gateway descriptions: update config/beans.oas.json and config/espresso.oas.json info descriptions, tags, and operation summaries so machine-readable positioning matches the public pages.
- Design authorities: add a dated intent-authority note to apis/design/BEANS_ESPRESSO_DOCUMENTATION_PLAN.md and apis/design/API_CAPABILITY_GAP_ANALYSIS.md, pointing to the two cited market references.

### Needed for usability

Add intent-driven entry points rather than only resource-driven explanations:

- Add a Beans What are you monitoring table for news, blogs, earnings, financial reports, litigation, lawsuits, official statements, research, technical documents, and podcasts.
- Add an Espresso What are you trying to understand table for concrete developments, market implications, outlook, evidence, source coverage, and monitoring.
- Add one provider migration matrix for Beans and one concept map for Espresso.
- Add one example for each of the three Beans non-news categories with the highest business value: earnings or financial reports, litigation or enforcement, and official statements or press releases.
- Add three Espresso market workflows: monitor, brief, and early warning.
- Add an explicit no-equivalence or not-currently-exposed note for industry routes that do not have a public Espresso equivalent.
- Add a shared Beans-to-Espresso workflow that demonstrates when to use both products.
- Expand the common conventions page with source/content availability, freshness, provenance, and provider-difference guidance.

### Needed for discovery

Update search and AI vocabulary to match both the category and the provider expectations:

- Beans metadata and page descriptions should include: news API, blog API, publisher content API, earnings reports API, financial reports API, lawsuit monitoring, litigation news, press releases, official statements, research papers, technical documents, and story clustering.
- Espresso metadata and page descriptions should include: market intelligence API, business intelligence API, event intelligence API, news intelligence API, event monitoring, company monitoring, impact, outlook, evidence, provenance, GDELT comparison, and Perigon comparison.
- The API overview should route article/news search to Beans and event/market intelligence to Espresso.
- The MCP guide should describe tool selection using the same intent vocabulary, not only tool names.
- The LLM Markdown and llms.txt output should include the provider concept maps and the explicit REST-only or not-currently-exposed boundaries.
- The portal metadata in docs/zudoku.config.tsx should remove unrelated terms and add the intent terms above.
- The migration pages should use provider names in titles or descriptions where that is useful for search discovery, while keeping factual comparisons bounded by the cited reference documents.

## Intent-distance completion criteria

### Beans is close enough to intent when

- The first screen says Beans is both a News API and a broader publisher-content API.
- A user can immediately find examples for news, blogs, earnings, financial reports, litigation, official statements, research, and technical content.
- The provider migration page maps TheNewsAPI, World News API, GNews, finlight, NewsAPI.ai/Event Registry, and NewsData.io to Beans route families.
- Full-content availability and publisher attribution limits are explicit.
- The docs clearly say Beans provides publisher material and coverage context, not synthesized Espresso intelligence.
- Search snippets and LLM output contain both News API terms and broader publisher-content terms.

### Espresso is close enough to intent when

- The first screen says Espresso is market, business, event, and news intelligence, not only generic business intelligence.
- A GDELT or Perigon user can map event, story, evidence, entity, geography, source, search, and aggregation concepts to Espresso or see that a concept is not currently exposed.
- Signals are explained through market, impact, outlook, and risk examples.
- The docs show monitoring, evidence-backed brief, and early-warning workflows.
- The docs distinguish Espresso Signals from raw provider events and Beans publisher Articles.
- Search snippets and LLM output contain market-intelligence, event-intelligence, evidence, provenance, GDELT, and Perigon vocabulary.


---

# 1. Inconsistent information that is self-conflicting or self-contradictory

## I-01. Current portal structure conflicts with repository documentation indexes

**Priority:** P0

### Gap

The current portal uses `start/`, `products/`, and `guides/`, but repository indexes and documentation plans still point to older `howtos/`, `introduction.mdx`, and `pages/pricing.mdx` paths. A maintainer, repository crawler, or AI agent following the README can be sent to files that no longer exist.

The old indexes also describe Espresso as “sips,” mention `same_as` and `derived_from`, and document `response_type=text`, while the current public contract uses Events, Signals, and JSON/YAML/TOON output.

### Missing or incorrect in

- [`docs/README.md`](docs/README.md): old portal paths, old product terminology, old pricing, and an obsolete product lineup.
- [`README.md`](README.md): repository setup and path references that need a current-docs cross-check.
- [`apis/README.md`](apis/README.md): legacy route and parameter documentation.
- [`apis/README.md`](apis/README.md): legacy “sip” model, response format, and relation-route documentation.
- [`apis/design/BEANS_ESPRESSO_DOCUMENTATION_PLAN.md`](apis/design/BEANS_ESPRESSO_DOCUMENTATION_PLAN.md): old destination paths.
- [`apis/design/BEANS_API_DOCUMENTATION_PLAN.md`](apis/design/BEANS_API_DOCUMENTATION_PLAN.md): old portal path assumptions.
- [`apis/design/ESPRESSO_API_DOCUMENTATION_PLAN.md`](apis/design/ESPRESSO_API_DOCUMENTATION_PLAN.md): old portal path assumptions.

### Edit or add

- Rewrite `docs/README.md` to index the current `docs/pages/start`, `docs/pages/products`, and `docs/pages/guides` tree.
- Update `README.md` to link to the current portal and identify backend READMEs as maintainer documentation.
- Rewrite the public-facing portions of both service READMEs around the current route contract, or clearly label them as historical/maintainer-only documents.
- Update all three documentation plans to current paths and mark superseded contract decisions.
- Keep old redirects in [`docs/zudoku.config.tsx`](docs/zudoku.config.tsx) only when they intentionally preserve external links.

### Completion check

Run a stale-path scan for `howtos/`, `introduction.mdx`, `pages/pricing.mdx`, `sips`, `same_as`, `derived_from`, and `response_type=text`. Every remaining match should be intentional and explicitly marked historical or internal.

## I-02. Beans advertises `post`, but the request binding does not accept it

**Priority:** P0

### Gap

The public Article response enum contains `post`, the gateway OpenAPI enum contains `post`, and Beans scenarios and migration docs list it. The current `ContentType` request binding in [`apis/beans/router/params.go`](apis/beans/router/params.go) omits `post`. The Beans overview also omits `post`, creating a second contradiction inside the portal.

The repository needs one explicit decision:

- `post` is a supported request filter and must be added to the runtime binding; or
- `post` is response-only/not publicly filterable and must be removed from the gateway spec and all public examples.

### Missing or incorrect in

- Runtime request binding: `apis/beans/router/params.go`.
- Runtime response enum: `apis/beans/router/responses.go`.
- Generated Beans Swagger: `apis/beans/docs/swagger.json`, `swagger.yaml`, and `docs.go`.
- Gateway contract: `config/beans.oas.json`.
- Human docs: `docs/pages/products/beans/overview.mdx`, `scenarios.mdx`, and `migration.mdx`.
- Contract intent: `apis/design/BEANS_API_ROUTE_PROPOSAL.md`.

### Edit or add

If `post` is supported:

- Add it to the binding in `apis/beans/router/params.go`.
- Add request/response contract tests in `apis/beans/tests/`.
- Update route annotations if the enum is emitted there.
- Regenerate `apis/beans/docs/*`.
- Verify `config/beans.oas.json` and all Beans pages.

If `post` is not supported as a request filter:

- Remove it from the public request parameter enum and request examples.
- Keep it only in response documentation if that is genuinely supported.
- Explain the distinction in the Beans overview.

### Completion check

The same content-type set must appear in `params.go`, response types, generated Swagger, gateway OpenAPI, route tables, examples, and MCP descriptions. A contract test should exercise `content_type=post` if it is supported.

## I-03. Beans top-headlines example sends a parameter that its documentation forbids

**Priority:** P0

### Gap

The Beans Scenario 1 example sends `content_type=news` to `/articles/top-headlines`, while the Beans overview says that route is fixed to news and does not accept `content_type`, `ids`, `urls`, `from`, or `to`.

The implementation currently shares a broad parameter struct across feed routes. That makes the public route contract unclear: a parameter may be accepted by binding but intentionally absent from the advertised route contract.

### Missing or incorrect in

- Example: `docs/pages/products/beans/scenarios.mdx`.
- Route explanation: `docs/pages/products/beans/overview.mdx`.
- Shared route binding: `apis/beans/router/params.go`.
- Route annotation and handler: `apis/beans/router/routes.go`.
- Generated and gateway specs: `apis/beans/docs/*` and `config/beans.oas.json`.

### Edit or add

- Remove `content_type` from the top-headlines examples immediately.
- Define a route-specific `topHeadlinesParams` type if the API must reject unsupported filters rather than silently tolerate them.
- Update the handler annotation, generated Swagger, gateway OpenAPI, route matrix, and tests together.
- State whether unknown query parameters are ignored or rejected.

### Completion check

Every example in `beans/scenarios.mdx` must validate against the exact parameter list for its route. Add a contract test for unsupported top-headlines parameters if rejection is intended.

## I-04. Espresso pagination disagrees across runtime, Swagger, gateway OpenAPI, and examples

**Priority:** P0

### Gap

The current Espresso `Pagination` response struct serializes `limit`, `num_results`, and `next_cursor`, but no `cursor`. The generated Swagger and gateway OpenAPI require both `cursor` and `next_cursor`, and the public portal examples show `cursor` in JSON but stale `page` fields in YAML and TOON.

**2026-08-25 note:** Published Beans and Espresso collection pagination is `limit`, `num_results` (this page only), and `next_cursor` (no `pagination.cursor` in the response). `meta.as_of` is required on Beans collections. Prefer gateway OpenAPI over this finding if they disagree. The older target-contract preference for `pagination.cursor` in the response is **superseded**.

### Missing or incorrect in

- Runtime response shape: `apis/espresso/router/responses.go`.
- Runtime response construction: `apis/espresso/router/routes.go`.
- Generated Swagger: `apis/espresso/docs/swagger.json`, `swagger.yaml`, and `docs.go`.
- Gateway contract: `config/espresso.oas.json`.
- Human examples: `docs/pages/products/espresso/overview.mdx` and `migration.mdx`.
- Shared guidance: `docs/pages/guides/api-conventions.mdx`.

### Edit or add

If `cursor` is part of the public response:

- Add and populate the serialized field in `apis/espresso/router/responses.go` and `routes.go`.
- Update annotations and regenerate `apis/espresso/docs/*`.
- Verify the gateway schema and examples.
- Replace `page: null` in YAML/TOON examples with the agreed cursor representation.

If `cursor` is request-only:

- Remove it from generated Swagger, gateway OpenAPI, and human response examples.
- Keep only `next_cursor` in the response contract.

### Completion check

JSON, YAML, and TOON examples must be projections of the same payload. The next-page instruction must say to send `next_cursor` unchanged as the next request’s `cursor`, without decoding or constructing tokens.

## I-05. Espresso resource terminology and response formats are stale in the service README

**Priority:** P0

### Gap

The service README describes a `sip`/`action`/`event`/`signal` model and `response_type=text`. The current public product describes Events and Signals, exposes no public Actions route, and documents JSON, YAML, and TOON.

### Missing or incorrect in

- `apis/README.md`, especially its data-model and response-format sections.
- `apis/design/ESPRESSO_API_DOCUMENTATION_PLAN.md` where old terminology remains.
- Any generated or gateway descriptions that still mention old formats or Actions.

### Edit or add

- Rewrite the service README around the current public Event/Signal route contract.
- Remove or clearly isolate storage-oriented “sip” terminology.
- Replace `response_type=text` with the current supported enum after verifying `apis/espresso/router/params.go`.
- Document what is stable versus dynamic in Event and Signal payloads.
- Remove public references to Actions unless an Action route is intentionally restored.

### Completion check

Searching the repository for `response_type=text`, `sip`, and public `action` route claims should return only intentionally historical/internal references.

## I-06. Error, empty-result, and status-code behavior conflicts across documents

**Priority:** P0

### Gap

Older service READMEs describe `204` for empty results and a flat `{ "error": "..." }` body. Current public descriptions use `200` with an empty collection and structured error envelopes. A new user cannot know which behavior to implement, and generated clients may produce the wrong success handling.

### Missing or incorrect in

- `apis/README.md`.
- `apis/README.md`.
- `apis/beans/router/responses.go` and `apis/espresso/router/responses.go`.
- Router annotations and generated Swagger in both services.
- `config/beans.oas.json` and `config/espresso.oas.json`.
- `docs/pages/guides/api-conventions.mdx`.
- Product overviews where error examples are incomplete or inconsistent.

### Edit or add

- Define one error envelope with exact serialized fields and examples.
- Define empty collection behavior separately from missing detail behavior.
- Align status-code tables in both gateway specs and both product overviews.
- Add examples for 400, 401, 404, 429, and 500 responses.
- Update local READMEs and regenerate backend Swagger after annotation changes.

### Completion check

A user should be able to determine from one conventions page whether an empty collection is `200`, whether a missing item is `404`, and how to parse every documented error response.

## I-07. Authentication metadata is incomplete or invalid in the public OpenAPI files

**Priority:** P0

### Gap

Beans declares a root security requirement referencing `ApiKeyAuth`, but the corresponding `securitySchemes` definition is absent. Espresso lacks a declared server URL and root security requirement even though the portal requires a Bearer API key.

This breaks or weakens “Try it,” SDK generation, automated API clients, and AI tools that rely on OpenAPI rather than prose.

### Missing or incorrect in

- `config/beans.oas.json`.
- `config/espresso.oas.json`.
- `docs/pages/guides/api-conventions.mdx` and `docs/pages/start/api-keys.mdx` as the human authentication authority.
- `config/policies.json` as the gateway behavior to verify against, not necessarily to change.

### Edit or add

- Declare the correct public server URL in both specs.
- Declare the exact security scheme and header format actually accepted by the gateway.
- Apply security requirements consistently to REST and MCP operations.
- Make the gateway spec, portal examples, and policy behavior agree on `Authorization: Bearer` versus any backend-only header.

### Completion check

`jq` validation should show a declared server and security scheme in each public spec. A generated client should be able to infer the base URL and authentication mechanism without reading portal prose.

## I-08. MCP tool lists conflict with the REST product documentation

**Priority:** P0

### Gap

Beans product docs describe Story workflows and refer to `listStories`, `getStory`, and `listStoryArticles`, but the Beans MCP operation list does not expose those tools. Espresso’s MCP guide lists only a subset of the discovery and source tools present in the gateway spec.

An AI agent may plan a tool call that does not exist or fail to discover a tool that does.

### Missing or incorrect in

- Beans MCP definitions in `config/beans.oas.json`.
- Espresso MCP definitions in `config/espresso.oas.json`.
- `docs/pages/guides/mcp-ai-agents.mdx`.
- `docs/pages/products/beans/overview.mdx`.
- `docs/pages/products/espresso/workflows.mdx`.
- Actual MCP route registration/implementation if missing tools are meant to be added.

### Edit or add

Choose and document one of these states for every capability:

- REST and MCP supported.
- REST-only.
- MCP-only.
- Planned/not yet available.

Then update the MCP catalog, tool descriptions, product workflows, and actual route exposure. Add a complete table of tool name, purpose, required inputs, REST equivalent, and expected output.

### Completion check

The MCP guide, gateway `x-zuplo-route` definitions, and actual MCP tool list must have the same operation inventory. Stories must be explicitly labeled REST-only if they remain absent from MCP.

## I-09. Pricing and product-status information is time-stale

**Priority:** P0

### Gap

The pricing page and `docs/README.md` say the free launch preview ends June 30, 2026. The current date is August 25, 2026. The repository README also lists MediCafe even though the current portal navigation lists Beans, Espresso, and Cortado.

### Missing or incorrect in

- `docs/pages/guides/pricing-limits.mdx`.
- `docs/README.md`.
- `docs/pages/start/overview.mdx` if product availability or limits have changed.
- `docs/pages/api-overview.mdx` if it contains product availability claims.

### Edit or add

- Replace expired launch language with current plan, limit, and availability information.
- Add an explicit “last updated” date and owner for pricing/limits.
- State whether limits apply per key, user, account, or product.
- Remove or label future products that do not have a public API.

### Completion check

No public page should contain an expired date or an unavailable product presented as live.

## I-10. Public/internal documentation boundaries are unclear

**Priority:** P1

### Gap

The service READMEs contain database schemas, table names, vector dimensions, indexes, relation storage, and internal pipeline details. They are not clearly labeled as internal maintainer documents. Repository crawlers and AI agents can treat them as public API documentation and extract obsolete implementation concepts.

The public pricing page also uses “vector search,” while the public documentation boundary calls for user-facing semantic-search language without exposing retrieval implementation.

### Missing or incorrect in

- `apis/README.md`.
- `apis/README.md`.
- `docs/README.md`.
- `docs/pages/guides/pricing-limits.mdx`.
- Any public page or generated description containing persistence, embedding, relation-storage, or infrastructure details.

### Edit or add

- Add an explicit maintainer-only heading to service READMEs, or split internal implementation notes from public API usage.
- Remove obsolete implementation details from any public/crawlable surface.
- Use “semantic search” in public docs unless provider behavior requires more specific terminology.
- Add a short public-boundary rule to the documentation contribution guidance.

### Completion check

An AI crawler reading only public portal content should learn product behavior, not database tables, vector dimensions, relation direction, or backend topology.

---

# 2. Information needed for usability

These items are not necessarily contradictions. They are missing instructions or missing contract detail that make it harder for a new API user to complete a task successfully.

## U-01. Add one canonical common-conventions page

**Priority:** P1

### Gap

Authentication, base URLs, pagination, response formats, error envelopes, empty results, date semantics, rate limits, and retries are distributed between product pages and short guides. A user must read several pages before implementing shared client behavior.

### Missing from

- `docs/pages/guides/api-conventions.mdx`.
- `docs/pages/start/overview.mdx`.
- `docs/pages/start/first-api-call.mdx`.
- Product-specific overview pages, which each explain only part of the shared behavior.

### Edit or add

Expand `docs/pages/guides/api-conventions.mdx` with:

- Base URL and authentication table.
- Exact request headers and content types.
- Collection and detail envelopes.
- Empty collection and missing detail behavior.
- Opaque cursor rules.
- `num_results` semantics.
- `as_of` and freshness semantics.
- Date filter meaning and timezone behavior.
- Response-format selection and content types.
- Status-code and error-envelope table.
- Rate-limit, retry, and backoff guidance.

Link this page prominently from both product overviews and both first-call examples.

### Completion check

A new user should be able to implement a generic pagination/error wrapper using only this page and the relevant route reference.

## U-02. Add field-level descriptions and examples to the machine-readable contracts

**Priority:** P1

### Gap

Many public schemas have properties and enums but no field descriptions. Espresso Event and Signal payloads are largely dynamic objects. Human pages explain some fields, but generated clients and AI agents cannot reliably infer field meaning, nullability, or stable versus variable fields.

### Missing from

- `apis/beans/router/responses.go` and `params.go` doc comments/Swagger annotations.
- `apis/espresso/router/responses.go` and `params.go` doc comments/Swagger annotations.
- Generated `apis/beans/docs/*` and `apis/espresso/docs/*`.
- `config/beans.oas.json` and `config/espresso.oas.json`.
- Espresso human overview where dynamic fields are not fully explained.

### Edit or add

- Add public field descriptions to response structs and route parameter annotations.
- Document required, nullable, optional, and conditionally present fields.
- Define stable common fields for Events and Signals, with an explicit extension-field policy.
- Add one representative request and response example per route family, including empty and error responses.
- Regenerate backend Swagger and reconcile gateway schemas manually.

### Completion check

An agent using only the OpenAPI document should know what `as_of`, `num_results`, `next_cursor`, `content_type`, `impact_level`, `event_types`, and dynamic digest fields mean.

## U-03. Add a route-by-route request and response matrix for Espresso

**Priority:** P1

### Gap

Beans has a reasonably detailed route matrix. Espresso has good narrative route selection but less centralized route-by-route documentation of accepted filters, required path values, output shape, and follow-up routes.

### Missing from

- `docs/pages/products/espresso/overview.mdx`.
- `docs/pages/products/espresso/workflows.mdx`.
- `config/espresso.oas.json` where parameter descriptions/examples are thin.

### Edit or add

Add a table with one row per public route containing:

- User intent.
- Required path/query parameters.
- Supported filters.
- Default and maximum limits.
- Response envelope.
- Empty-result behavior.
- Common follow-up routes.
- REST/MCP availability.

Use the router bindings and tests as the input source rather than copying the older README.

### Completion check

Users should not need to infer route-specific parameters by comparing multiple examples.

## U-04. Add troubleshooting and failure-recovery guidance

**Priority:** P1

### Gap

The portal explains successful calls but does not provide a single troubleshooting path for invalid parameters, missing credentials, rate limits, empty results, stale cursors, unsupported formats, or temporary service failures.

### Missing from

- `docs/pages/start/` has no dedicated troubleshooting page.
- `docs/pages/guides/api-conventions.mdx` has no complete recovery table.
- `docs/pages/guides/mcp-ai-agents.mdx` has no MCP-specific troubleshooting.

### Edit or add

Add `docs/pages/start/troubleshooting.mdx` and link it from `docs/zudoku.config.tsx`, `start/overview.mdx`, and `start/first-api-call.mdx`.

Include:

- 401/403 authentication diagnosis.
- 400 parameter validation examples.
- 404 detail-resource behavior.
- 429 retry guidance.
- 5xx fallback behavior.
- Empty result interpretation.
- Cursor expiration or invalid-cursor handling.
- Response-format fallback.
- MCP connection and tool-discovery failures.

### Completion check

Each documented error should answer: what happened, whether to retry, what to change, and which page or route to use next.

## U-05. Add complete reusable REST client examples

**Priority:** P1

### Gap

The portal has JavaScript, Python, and cURL snippets, but most are isolated calls. It lacks reusable examples showing authentication, query construction, pagination, error handling, and selective detail enrichment in one place.

### Missing from

- `docs/pages/start/first-api-call.mdx`.
- `docs/pages/products/beans/scenarios.mdx`.
- `docs/pages/products/espresso/workflows.mdx`.
- No dedicated client-patterns page.

### Edit or add

Add `docs/pages/guides/client-patterns.mdx` and link it from the start pages. Provide:

- TypeScript/Node client wrapper.
- Python client wrapper.
- cURL shell pattern.
- Cursor continuation helper.
- Structured error handling.
- `response_type` selection.
- Environment-variable and secret-handling guidance.

Keep product workflows focused on product decisions and link to this page for reusable transport code.

### Completion check

A user should be able to copy one complete client pattern and then change only the route and parameters.

## U-06. Add operational MCP client setup, not just MCP concepts

**Priority:** P1

### Gap

The MCP guide explains why to use MCP and names endpoints/tools, but does not provide copyable configuration for common MCP clients, transport expectations, authentication placement, or connection troubleshooting.

### Missing from

- `docs/pages/guides/mcp-ai-agents.mdx`.
- `docs/pages/start/overview.mdx`, which links to MCP but does not show setup.
- Gateway MCP operation metadata, which uses generic request/response objects rather than tool-specific schemas.

### Edit or add

Expand `mcp-ai-agents.mdx` with:

- Generic MCP client configuration.
- Common desktop/editor configuration examples where supported.
- Endpoint and authentication fields.
- Tool discovery/handshake expectations.
- One complete agent task from tool selection through result enrichment.
- Troubleshooting for authentication, missing tools, empty results, and format selection.

If client-specific configuration cannot be guaranteed, add a generic configuration and explicitly state the supported transport.

### Completion check

A user should be able to configure an MCP client without reverse-engineering the gateway OpenAPI extension.

## U-07. Explain dynamic Event and Signal payloads as a usable schema

**Priority:** P1

### Gap

Espresso deliberately preserves variable upstream fields, but the documentation does not give users enough guidance to safely consume those records. The public contract needs to distinguish stable fields from optional or provider-specific fields.

### Missing from

- `apis/espresso/router/responses.go` and generated schemas.
- `config/espresso.oas.json` Event/Signal schemas.
- `docs/pages/products/espresso/overview.mdx`.
- `docs/pages/products/espresso/workflows.mdx`.

### Edit or add

Document:

- Stable identifiers and timestamps.
- Stable classification/filter fields.
- Common evidence/source fields.
- Optional fields that may be absent or null.
- Extension fields that clients must ignore safely.
- Recommended parsing strategy for JSON, YAML, and TOON.

Add a stable schema envelope even if the inner digest remains extensible.

### Completion check

An SDK or agent can safely extract identifiers, timestamps, summaries, sources, and follow-up IDs without assuming every record has the same digest fields.

## U-08. Add a canonical Beans-to-Espresso composite workflow

**Priority:** P1

### Gap

The individual product workflows are strong, but there is no first-class scenario showing how to combine the services into one application. The product boundary is explained, but the handoff is not demonstrated.

### Missing from

- `docs/pages/products/beans/scenarios.mdx`.
- `docs/pages/products/espresso/workflows.mdx`.
- `docs/pages/guides/mcp-ai-agents.mdx`.
- `docs/pages/api-overview.mdx`, which is too brief to orient a cross-product workflow.

### Edit or add

Add `docs/pages/guides/cross-product-workflows.mdx` and add it to `docs/zudoku.config.tsx`.

Include at least:

1. Discover a topic, publisher, or article in Beans.
2. Select a small set of relevant identifiers or entities.
3. Search Espresso Events or Signals using those terms.
4. Retrieve evidence or supporting Events selectively.
5. Present provenance and freshness to the end user.

Provide cURL, Python, or TypeScript examples and explain when the products should not be combined.

### Completion check

The workflow should show concrete data handoff, not just link from one product page to another.

## U-09. Make freshness and temporal semantics explicit

**Priority:** P1

### Gap

The APIs expose `as_of` and date filters, but users need a single explanation of whether dates mean record creation, publication, occurrence, forecast, or collection time. Espresso’s route annotations say `from`/`to` apply to `created_at`, but this rule is not centralized across products.

### Missing from

- `docs/pages/guides/api-conventions.mdx`.
- `docs/pages/products/beans/overview.mdx`.
- `docs/pages/products/espresso/overview.mdx`.
- `docs/pages/products/beans/migration.mdx` and `espresso/migration.mdx`.

### Edit or add

- Add a shared date/filter semantics section.
- Define timezone and date-only interpretation.
- Explain `as_of` as response freshness metadata.
- Label captured examples with their capture date and explain that values may change.
- Show how applications should display freshness to users.

### Completion check

An application developer can choose date filters without guessing which timestamp is being bounded.

## U-10. State SDK, versioning, support, and limits expectations

**Priority:** P1

### Gap

The start pages do not clearly state whether official SDKs exist, how API compatibility is versioned, where breaking changes are announced, or which limits apply to REST versus MCP.

### Missing from

- `docs/pages/start/overview.mdx`.
- `docs/pages/start/first-api-call.mdx`.
- `docs/pages/guides/pricing-limits.mdx`.
- `docs/pages/api-overview.mdx`.

### Edit or add

Add a short “Using the API in production” section covering:

- Official SDK availability or explicit REST/MCP-only status.
- API versioning policy.
- Backward-compatibility expectations.
- Deprecation notices.
- REST/MCP rate and quota accounting.
- Support and issue-reporting channels.

### Completion check

New users know what integration surface to select and how to plan for contract changes.

## U-11. Surface the executable Bruno examples as maintained examples

**Priority:** P1

### Gap

Current Bruno collections under `apis/beans/tests/bruno` and `apis/espresso/tests/bruno` are valuable executable examples, but the public portal does not point users to them and there is no stated relationship between those examples and the public contract.

### Missing from

- Beans and Espresso product pages.
- `docs/pages/guides/api-conventions.mdx`.
- Repository README indexes.

### Edit or add

- Add a maintained-examples section to each product overview.
- Link to the repository collection or publish a sanitized downloadable collection.
- Add a README in each Bruno directory explaining variables, authentication, and expected service availability.
- Add a CI check that validates example paths and required parameters against the contract tests.

### Completion check

Users can find a runnable collection, understand its variables, and know whether it is tested against the current contract.

---

# 3. Information needed for discovery through search engines and AI agents

## D-01. Replace broad and partly irrelevant site metadata

**Priority:** P1

### Gap

The portal metadata includes useful terms such as “news api,” “news mcp,” and “market intelligence,” but also unrelated or premature terms such as “billing software,” “PR management,” and “tech startup.” This can dilute search relevance and cause an AI index to associate the portal with products that are not currently documented.

### Missing or incorrect in

- `docs/zudoku.config.tsx` metadata block.
- Frontmatter in product overview, scenario, workflow, migration, guide, and start pages.
- `docs/README.md` product summaries.

### Edit or add

- Replace broad keywords with task-oriented terms: news API, article search API, publisher/source discovery, story aggregation, business intelligence API, event intelligence, signal search, evidence retrieval, MCP tools, and AI research workflows.
- Give every high-value page a specific title and description.
- Include product, audience, capability, and outcome in descriptions.
- Remove unavailable-product terms from current metadata.

### Completion check

Search snippets for Beans, Espresso, MCP, and API onboarding should be distinct rather than sharing generic portal metadata.

## D-02. Add or verify technical SEO assets

**Priority:** P2

### Gap

The repository configures Pagefind and LLM Markdown output, but no repository-owned `robots.txt`, `sitemap.xml`, JSON-LD, canonical URL, or structured social metadata was found. The portal may generate some of these at build time, but that output was not established by the source tree.

### Missing from

- `docs/zudoku.config.tsx` lacks an evident complete technical-SEO configuration.
- No checked-in `docs/public/robots.txt` or `docs/public/sitemap.xml` was found.
- Page frontmatter does not consistently provide specific titles and descriptions.

### Edit or add

- First inspect the built production output to avoid duplicating generated assets.
- If absent, configure or add `robots.txt`, sitemap generation, canonical site URL, Open Graph/Twitter metadata, and JSON-LD for the developer portal and API products.
- Ensure API reference pages are indexable while internal design documents are not presented as public product pages.

### Completion check

Production output has a valid sitemap, crawl policy, canonical URLs, and product-specific structured metadata.

## D-03. Strengthen the API overview as a search and AI landing page

**Priority:** P1

### Gap

`docs/pages/api-overview.mdx` is too short to capture the vocabulary users search for. It links to references but does not provide a concise inventory of product capabilities, route families, formats, MCP endpoints, or common tasks.

### Missing from

- `docs/pages/api-overview.mdx`.
- `docs/pages/start/overview.mdx` for cross-product capability keywords.
- Product overview pages if important route-family terms are only present in tables or code.

### Edit or add

Expand the API overview with:

- Beans capability and route-family summary.
- Espresso capability and route-family summary.
- REST and MCP endpoint links.
- Search terms mapped to product and route families.
- Common tasks and the best starting page for each.
- Links to authentication, conventions, workflows, and migration pages.

### Completion check

An AI agent or search user who asks “which Cafecito API finds articles,” “which API finds events and signals,” or “how do I use Cafecito with MCP” lands on a useful routing page.

## D-04. Create a canonical AI-agent and MCP capability catalog

**Priority:** P1

### Gap

`mcp-ai-agents.mdx` is useful but conceptual and incomplete. It does not provide a complete tool inventory, capability-to-tool mapping, REST equivalents, or explicit REST-only/planned status.

### Missing from

- `docs/pages/guides/mcp-ai-agents.mdx`.
- `config/beans.oas.json` and `config/espresso.oas.json` tool descriptions where schemas are generic.
- `docs/pages/products/beans/overview.mdx` and `espresso/workflows.mdx` where MCP references are selective.

### Edit or add

- Add complete Beans and Espresso MCP tool tables.
- Include tool name, user intent, required inputs, optional inputs, response shape, REST equivalent, and next recommended tool.
- State which tools are discovery-only and which perform enrichment.
- Mark Stories as REST-only if that remains true.
- Add a machine-readable tool manifest only if the supported MCP tooling can consume it reliably; do not create a second conflicting catalog.

### Completion check

An agent can select a tool from the catalog without reading implementation code or guessing operation names.

## D-05. Improve OpenAPI discoverability for generated clients and AI retrieval

**Priority:** P1

### Gap

The public specs contain many paths and schemas, but weak field descriptions, incomplete security metadata, sparse request examples, and dynamic Espresso objects. This makes the API technically present but semantically difficult to retrieve.

### Missing from

- `config/beans.oas.json`.
- `config/espresso.oas.json`.
- Source annotations and response comments in both routers.
- Generated Swagger artifacts.

### Edit or add

- Add stable `summary`, `description`, `tags`, parameter descriptions, enum meanings, and examples.
- Add explicit operation-level examples for common tasks.
- Add security schemes and servers.
- Document pagination and error semantics in reusable components.
- Give dynamic Event and Signal schemas a stable outer shape and documented extension policy.

### Completion check

An OpenAPI-only client or AI agent can identify the right route, construct a valid request, authenticate, parse the response, and continue pagination.

## D-06. Make page titles and descriptions specific to search intent

**Priority:** P2

### Gap

Some product pages use generic titles such as “Overview,” while other pages rely mainly on an H1 or description. Generic titles reduce useful search snippets and make LLM indexes less discriminative.

### Missing from

- `docs/pages/products/beans/overview.mdx`.
- `docs/pages/products/espresso/overview.mdx`.
- Beans and Espresso scenario/workflow/migration pages.
- Start and guide pages with generic or absent title metadata.

### Edit or add

Use titles such as:

- “Beans News and Article Search API”
- “Beans API Workflows for Feeds, Sources, and Stories”
- “Espresso Events and Signals API”
- “Espresso Workflows for Evidence and Business Intelligence”
- “Cafecito MCP Tools for AI Agents”

Descriptions should state audience, capability, and outcome rather than only repeat the page heading.

### Completion check

Every indexable page has a unique title and a concise description containing the terms users are likely to search.

## D-07. Improve repository-crawler hygiene

**Priority:** P1

### Gap

Repository crawlers and AI agents see stale READMEs and design plans before or alongside the current portal. Even if the public portal is correct, repository-level discovery can return obsolete contracts.

### Missing or incorrect in

- `README.md`.
- `docs/README.md`.
- `apis/README.md`.
- `apis/README.md`.
- `apis/design/*.md` documentation plans and gap reports.

### Edit or add

- Add a current-documentation pointer at the top of repository and service READMEs.
- Mark design plans as current, superseded, or historical with dates.
- Add a “public contract source” section pointing to router code, generated Swagger, gateway OAS, and portal docs.
- Move internal schema notes to clearly named maintainer documents if they are still useful.
- Remove dead links and obsolete route examples.

### Completion check

An AI crawler starting at the repository root reaches current documentation first and can distinguish design intent from live behavior.

## D-08. Link maintained executable examples into the discoverable documentation graph

**Priority:** P2

### Gap

The Bruno collections are discoverable only by browsing service tests. Search engines and AI agents cannot infer that they are maintained usage examples.

### Missing from

- Product overview pages.
- `docs/pages/guides/api-conventions.mdx`.
- Root and service README indexes.
- The Bruno directories lack a public-facing index explaining their purpose.

### Edit or add

- Add links from Beans and Espresso overview pages to the corresponding Bruno collection or a sanitized public export.
- Add a short `README.md` in each Bruno product directory describing setup and variables.
- Add stable example names matching public operation IDs.
- Include the collection in the documentation validation checklist.

### Completion check

Searching for a route or operation ID finds both its reference documentation and a runnable example.

## D-09. Add explicit search vocabulary and cross-links for common user intents

**Priority:** P2

### Gap

The docs explain capabilities, but common intent phrases are not consistently mapped to routes. Users and agents may search for “publisher monitoring,” “topic feed,” “related articles,” “event evidence,” “signal provenance,” or “market outlook” without finding a direct route recommendation.

### Missing from

- `docs/pages/api-overview.mdx`.
- Beans and Espresso overview route selectors.
- Scenario/workflow introductions.
- `docs/pages/guides/mcp-ai-agents.mdx`.

### Edit or add

Add intent-to-route tables and cross-links:

- News/topic feed -> Beans latest/top-headlines/search.
- Publisher monitoring -> Beans sources and source article routes.
- Related coverage/story grouping -> Beans Stories or article follow-up routes.
- What happened -> Espresso Events.
- What does it mean/outlook -> Espresso Signals.
- Why should I trust this -> Espresso evidence and supporting Events.
- Which filter value exists -> discovery routes.

### Completion check

Natural-language task phrases resolve to a product, route, example, and next call.

## D-10. Ensure generated Markdown and LLM indexes exclude stale/internal material

**Priority:** P1

### Gap

The portal enables published Markdown and `llms.txt` generation, but the source tree contains stale documentation indexes and internal design material. The build needs an explicit inclusion policy so generated AI indexes prioritize current public pages and do not expose internal implementation notes.

### Missing from

- `docs/zudoku.config.tsx` inclusion/navigation policy.
- `docs/README.md` current-content inventory.
- `apis/design/*.md` status labeling.
- Public-boundary guidance for generated Markdown/LLM output.

### Edit or add

- Verify the actual generated `llms.txt` and full Markdown output.
- Exclude internal design and maintainer-only documents from public output.
- Ensure current product, workflow, guide, and API-reference pages are included.
- Add a build check for stale paths, old terms, and missing product pages.

### Completion check

The generated LLM index contains one current description of each product and route family, with no obsolete README contract or internal schema material.

---

# 7-Why Architectural RCCA

## RCCA problem statement

**Problem:** A public Cafecito API capability can be described differently by runtime code, generated backend Swagger, gateway OpenAPI, portal prose, MCP metadata, examples, service READMEs, and design documents. Product intent can also be present in market-reference documents without becoming visible in public entry points.

**Impact:**

- A developer can construct a request that the runtime rejects.
- Generated clients and AI agents can infer the wrong server, authentication method, parameter enum, response envelope, or MCP tool.
- Search users can misunderstand Beans as only News & Blogs or fail to recognize Espresso as market and event intelligence.
- Maintainers can update one surface successfully while leaving downstream surfaces stale.
- Deployment can proceed without proving that the public contract and documentation still match the service.

**Scope:** Beans and Espresso router contracts, generated Swagger, gateway OpenAPI, Zudoku pages, MCP exports, executable examples, repository READMEs, design-intent documents, and CI workflows.

**Observed architecture:**

    Product intent and market references
      -> route proposal and implementation decisions
      -> request bindings, handlers, response structs, and tests
      -> generated backend Swagger
      -> manually maintained gateway OpenAPI and MCP exports
      -> hand-written portal pages and examples
      -> Pagefind, Markdown publishing, and LLM indexes

The flow is documented, but it is not enforced as one atomic change system.

## Seven Whys

### Why 1: Why can users and agents not fully trust the documentation?

Because the same public behavior is described differently in multiple places.

Examples already identified in this report include:

- Beans content_type values differ between request binding, response enum, gateway OpenAPI, and portal pages.
- The Beans top-headlines scenario sends a parameter that its own route documentation says is unsupported.
- Espresso pagination differs between the runtime response, generated Swagger, gateway OpenAPI, and JSON/YAML/TOON examples.
- Authentication metadata is incomplete in the public OpenAPI files.
- MCP tool catalogs do not match the capabilities described in product pages.
- Local READMEs retain obsolete routes, response formats, status behavior, and terminology.

**Immediate correction:** Resolve the P0 contradictions in I-02 through I-09 and verify every affected surface in one work packet.

### Why 2: Why is the same behavior described differently?

Because the public contract is replicated into independently edited representations:

1. Go request bindings and handlers.
2. Go response structs and serialized field names.
3. Swaggo annotations.
4. Generated backend Swagger.
5. Manually maintained gateway OpenAPI.
6. Gateway MCP operation lists and tool descriptions.
7. Zudoku product pages, workflows, and examples.
8. Service and repository READMEs.
9. Design plans and market-intent references.

Some duplication is unavoidable because each surface serves a different audience. The failure is that duplicated facts such as paths, parameters, enums, envelopes, authentication, operation IDs, and tool exposure have no automated parity check.

**Architectural correction:** Define an authority for each fact and generate or verify every duplicate representation against that authority.

### Why 3: Why are duplicated facts not generated or reconciled automatically?

Because the repository has generation commands for backend Swagger but no contract-cascade pipeline.

Current behavior:

- Swaggo generation is documented in service READMEs and AGENTS guidance.
- Generated Swagger is committed.
- Gateway OpenAPI is maintained separately.
- Portal pages and examples are authored manually.
- MCP exports are encoded in gateway OpenAPI extensions.
- No repository script compares backend paths and parameters with public gateway paths after removing the product prefix.
- No repository script compares operation IDs with MCP exports or examples.

**Architectural correction:** Add a deterministic verifier that produces a normalized contract inventory from existing sources. The inventory must be generated evidence, not another hand-maintained source of truth.

### Why 4: Why is contract drift not blocked by CI?

Because CI follows deployment topology instead of public-contract topology.

Evidence:

- .github/workflows/zpl-deploy-gateway.yml ignores apis/** changes and currently runs npm ci and npm run lint.
- .github/workflows/fly-deploy-beans.yml deploys after checkout and Fly setup without a contract-test or documentation-parity job.
- .github/workflows/fly-deploy-espresso.yml has the same direct-deploy pattern.
- The root package.json has lint and gateway test commands but no cross-surface contract verification command.
- docs/package.json has a Zudoku build command, but the gateway workflow does not establish a documentation-contract gate.

A backend route change and its public documentation are one logical product change even though they live in separate deployment units.

**Architectural correction:** Add an independent API-contract-and-docs workflow triggered by backend, gateway, portal, MCP, example, and design-authority changes. Deployment workflows should depend on or duplicate the required verification gate.

### Why 5: Why can incorrect examples and semantic claims survive review?

Because examples and prose are reviewed as content rather than executable contract fixtures.

Current examples are distributed across:

- Product overview pages.
- Beans scenarios.
- Espresso workflows.
- Migration guides.
- MCP guidance.
- Bruno collections.
- OpenAPI examples.

There is no shared example registry or automated check that confirms:

- The route exists.
- Every query parameter is accepted on that route.
- Enum values match runtime binding.
- The expected envelope matches the response schema.
- Pagination examples use the current fields.
- MCP tool names exist in the exported catalog.

**Architectural correction:** Promote representative examples into testable fixtures and derive or validate public snippets against them.

### Why 6: Why do stale READMEs, plans, and intent documents remain discoverable as current truth?

Because documents have no enforced lifecycle or audience metadata.

The repository contains live contracts, target proposals, gap analyses, market references, generated outputs, maintainer notes, and historical READMEs. They often lack a consistent declaration of:

- Status: current, target, superseded, historical, or generated.
- Authority: runtime, public contract, product intent, or implementation note.
- Audience: public user, maintainer, operator, or internal design.
- Last verified date.
- Owning role.
- Replacement document when superseded.

Search engines, repository users, and AI agents cannot reliably infer those distinctions from location alone.

**Architectural correction:** Require lifecycle metadata and current-authority links at the top of every design document and service README. Exclude internal and historical documents from public Markdown and LLM indexes.

### Why 7: Why are contract cascade, lifecycle, and ownership controls not enforced?

Because the API change process has guidance but no enforceable Definition of Done and no single accountable contract owner per product.

AGENTS guidance describes the documentation dependency chain, but completion still depends on a contributor remembering every surface. There is no required owner sign-off or machine gate for:

- Runtime-to-Swagger parity.
- Swagger-to-gateway parity.
- Gateway-to-MCP parity.
- Gateway-to-portal route and parameter parity.
- Intent-to-positioning parity.
- Example executability.
- Document lifecycle and expiry.

**Architectural correction:** Establish product-contract ownership, CODEOWNERS coverage, a required API-change checklist, and CI checks that make an incomplete cascade unmergeable.

## Root-cause conclusion

The primary root cause is not insufficient documentation effort. It is an architectural control failure:

> Cafecito represents one public API contract in multiple independently maintained artifacts without a generated reconciliation layer, while CI and ownership are organized around deployment units rather than the end-to-end public contract.

Three root causes must be addressed together:

| Root cause | Classification | Consequence |
|---|---|---|
| Contract facts are manually replicated across code, specs, portal pages, MCP metadata, and examples. | Architecture | Drift is expected whenever one surface changes. |
| CI validates deployment units but not the cross-surface public contract. | Control architecture | Contradictory documentation can pass CI and deploy. |
| Documents and contract surfaces lack explicit lifecycle state and accountable ownership. | Governance architecture | Stale or internal material remains discoverable and no role owns closure. |

Contributing factors:

- Generated Swagger and manual gateway OpenAPI have different production purposes.
- Product prefixes exist at the gateway but not on backend routes.
- MCP tools are a selected subset rather than an automatic mirror of REST.
- Espresso has intentionally extensible records that need both stable-core and extension-field documentation.
- Portal prose contains task guidance that cannot be fully generated from OpenAPI.
- Provider comparison and product intent change on a different cadence from runtime behavior.

These factors justify multiple representations. They do not justify unverified divergence.

## Corrective-action architecture

### CA-0. Contain current misinformation

**Purpose:** Stop known contradictions from spreading while structural controls are built.

Actions:

1. Resolve I-02 through I-09.
2. Mark stale service READMEs as maintainer-only until rewritten.
3. Add current, target, and historical status labels to design documents.
4. Remove or qualify claims not confirmed by runtime bindings and response structs.
5. Keep public one-line product intent consistent across start pages, API overview, product pages, gateway descriptions, and MCP guidance.

Primary files:

- apis/beans/router/params.go
- apis/beans/router/responses.go
- apis/beans/router/routes.go
- apis/espresso/router/params.go
- apis/espresso/router/responses.go
- apis/espresso/router/routes.go
- config/beans.oas.json
- config/espresso.oas.json
- docs/pages/**
- apis/README.md
- apis/README.md
- docs/README.md

Exit condition: no known P0 contradiction remains open without an explicit planned/not-available label.

### CA-1. Define authority by fact type

**Purpose:** Avoid creating a new monolithic source of truth that cannot represent runtime, gateway, product, and human-guidance concerns correctly.

| Fact | Authority | Verified projections |
|---|---|---|
| Route existence and backend path | Router registration and contract tests | Generated Swagger, gateway path after prefix mapping, portal route table |
| Request parameters, defaults, limits, and enums | Binding structs plus validation tests | Annotations, Swagger, gateway OpenAPI, examples |
| Serialized response and error fields | Response structs plus handler tests | Swagger schemas, gateway schemas, portal examples |
| Public server, authentication, policy, and product prefix | Gateway OpenAPI and gateway policy configuration | Start guides, API conventions, generated clients |
| MCP tool exposure | Gateway MCP operation configuration | MCP guide and product workflows |
| Product intent and category | Dated intent reference approved by product owner | Product one-liners, metadata, migration pages, OpenAPI descriptions |
| Human workflow and route-selection guidance | Product documentation, validated against the contract inventory | Pagefind, Markdown, llms.txt |
| Generated backend Swagger | Router annotations generated with the pinned command | Committed generated artifacts |

Edit:

- AGENTS.md
- apis/AGENTS.md
- apis/design/BEANS_ESPRESSO_DOCUMENTATION_PLAN.md
- Add a short documentation-contribution page under docs or repository guidance.

Exit condition: every duplicated public fact has a named authority and a parity rule.

### CA-2. Add a generated contract-cascade verifier

**Purpose:** Detect cross-surface drift without introducing another manually edited contract.

Add scripts/verify-api-contract-cascade.mjs or an equivalent implementation that:

1. Loads apis/beans/docs/swagger.json and apis/espresso/docs/swagger.json.
2. Loads config/beans.oas.json and config/espresso.oas.json.
3. Applies the explicit backend-to-gateway prefix mapping.
4. Compares route methods, operation IDs, parameters, defaults, enums, required flags, response status codes, and envelope schemas.
5. Verifies servers and security schemes.
6. Extracts MCP operation inventories and checks documented tool names.
7. Emits a normalized, generated report for review.
8. Uses a small reviewed exception file only for intentional gateway-only, backend-only, health, docs, and MCP differences.

Recommended files:

- Add scripts/verify-api-contract-cascade.mjs
- Add config/api-contract-exceptions.json
- Update package.json with verify:api-contracts
- Update AGENTS.md with the verification command

The exception file must include a reason, owner, and expiry/review date. It must not become a permanent list of unexplained drift.

Exit condition: introducing an undocumented route, parameter, enum, security, response, or MCP mismatch fails locally and in CI.

### CA-3. Add a documentation-contract CI workflow

**Purpose:** Align CI with the public contract rather than only deployment units.

Add .github/workflows/api-contract-docs.yml triggered by changes to:

- apis/beans/router/**
- apis/espresso/router/**
- apis/beans/docs/**
- apis/espresso/docs/**
- config/beans.oas.json
- config/espresso.oas.json
- docs/**
- scripts/verify-*.mjs
- API documentation plans and current intent authorities

Required jobs:

1. Run Beans router contract tests.
2. Run Espresso router contract tests.
3. Regenerate Swagger with pinned commands and fail if committed generated artifacts differ.
4. Validate both gateway OpenAPI files.
5. Run the contract-cascade verifier.
6. Run documentation example and stale-term checks.
7. Build the Zudoku portal.
8. Inspect generated Markdown and llms.txt for required pages and forbidden internal terms.

Modify:

- .github/workflows/zpl-deploy-gateway.yml
- .github/workflows/fly-deploy-beans.yml
- .github/workflows/fly-deploy-espresso.yml

The gateway workflow can continue to ignore backend-only changes for deployment. The independent contract workflow must not ignore them.

Exit condition: no Beans or Espresso deployment can proceed from a commit that fails the relevant contract-and-docs gate.

### CA-4. Make examples executable contract assets

**Purpose:** Prevent invalid snippets and stale envelopes.

Add scripts/verify-doc-examples.mjs or equivalent behavior that:

- Maps named examples to operation IDs.
- Verifies paths and accepted query parameters.
- Validates JSON examples against public schemas where schemas are stable.
- Checks pagination field names across JSON, YAML, and TOON examples.
- Checks MCP tool names against exported tools.
- Reports illustrative examples that cannot be executed.

Use the Bruno collections as the executable source for representative requests where practical:

- apis/beans/tests/bruno/**
- apis/espresso/tests/bruno/**

Edit:

- docs/pages/products/beans/overview.mdx
- docs/pages/products/beans/scenarios.mdx
- docs/pages/products/espresso/overview.mdx
- docs/pages/products/espresso/workflows.mdx
- docs/pages/guides/mcp-ai-agents.mdx
- package.json

Exit condition: every public route family has at least one validated request, success response, empty response where applicable, and error example.

### CA-5. Add document lifecycle and audience controls

**Purpose:** Prevent stale design and maintainer material from appearing current.

Required metadata for apis/design/*.md and service READMEs:

- Status: current, target, superseded, historical, or generated.
- Authority: product intent, live contract, implementation note, or external reference.
- Audience.
- Last verified date.
- Owner role.
- Superseded-by link when applicable.

Edit:

- apis/design/*.md
- apis/README.md
- apis/README.md
- docs/README.md
- docs/zudoku.config.tsx

Add a verification script that fails on missing metadata for designated authority files and flags expired review dates.

Exit condition: a repository user or AI agent can distinguish live behavior, intended future behavior, external comparison, internal implementation, and historical material from the first screen.

### CA-6. Centralize product positioning without centralizing the whole contract

**Purpose:** Keep Beans and Espresso intent consistent across human and machine discovery.

Add a small shared portal product-definition module, for example docs/product-definitions.ts, containing:

- Product name and status.
- Approved one-line intent.
- Positive use cases.
- Not-for boundary.
- Core public resource names.
- Search/AI metadata terms.

Use or validate it from:

- docs/pages/start/overview.mdx
- docs/pages/api-overview.mdx
- Beans and Espresso overview pages
- docs/zudoku.config.tsx
- docs/pages/guides/mcp-ai-agents.mdx

Gateway OpenAPI descriptions should be checked against the approved wording, not generated blindly from portal code.

Exit condition: product entry points cannot silently diverge on what Beans and Espresso are for.

### CA-7. Establish ownership and change governance

**Purpose:** Make contract integrity an accountable product responsibility.

Add or update:

- .github/CODEOWNERS
- Pull-request template or API-change checklist
- AGENTS.md
- apis/AGENTS.md

| Role | Accountability |
|---|---|
| Beans service owner | Runtime, bindings, response structs, annotations, tests |
| Espresso service owner | Runtime, bindings, response structs, annotations, tests |
| Gateway contract owner | Public prefixes, auth, security, rate policy, gateway OpenAPI, MCP exports |
| Documentation owner | Portal information architecture, examples, search metadata, Markdown/LLM output |
| Product owner | Intent, product boundary, provider comparison, capability availability |
| CI owner | Contract verifier, generated-artifact checks, workflow health |

Required Definition of Done for a public API change:

1. Runtime behavior and tests are complete.
2. Annotations are updated.
3. Generated Swagger is regenerated.
4. Gateway OpenAPI is reconciled.
5. MCP exposure is reconciled.
6. Portal pages and examples are reconciled.
7. Intent and boundary language remain correct.
8. Contract, docs build, examples, and stale-term checks pass.

Exit condition: each failed gate has one role responsible for resolution, and API changes cannot be marked complete before the cascade closes.

## Corrective-action matrix

| Action | Why addressed | Priority | Files | Owner role | Verification |
|---|---|---|---|---|---|
| Resolve known contract contradictions | Why 1 | P0 | Router files, generated Swagger, gateway OAS, affected pages | Service plus gateway owner | Contract tests and parity report |
| Define authority by fact type | Why 2 | P0 | AGENTS files and documentation plan | Product plus architecture owner | Authority table reviewed |
| Add contract-cascade verifier | Why 3 | P0 | New verifier, package.json, exception file | CI plus gateway owner | Intentional mutation causes CI failure |
| Add cross-surface CI workflow | Why 4 | P0 | New workflow and deploy workflow dependencies | CI owner | Required check on API/doc changes |
| Validate examples and MCP tools | Why 5 | P1 | Example verifier, Bruno, portal pages | Documentation plus service owner | Invalid parameter/tool fixture fails |
| Add lifecycle metadata | Why 6 | P1 | Design docs, READMEs, portal inclusion rules | Documentation owner | Metadata/staleness check passes |
| Add CODEOWNERS and Definition of Done | Why 7 | P1 | CODEOWNERS, PR template, AGENTS files | Product and engineering leads | Required reviewers and checklist |
| Centralize product definitions | Why 2 and 6 | P1 | Shared portal module and entry pages | Product plus documentation owner | One-line intent equality check |
| Verify generated search/AI output | Why 4 through 6 | P1 | Zudoku build and CI checks | Documentation owner | Required pages present; forbidden terms absent |

## Preventive controls

### Per-change controls

- Backend route changes trigger service tests, Swagger regeneration checks, gateway parity, and docs build.
- Gateway OAS changes trigger backend-path comparison, security checks, MCP inventory checks, and portal build.
- Portal example changes trigger operation/parameter/schema validation.
- Product-intent changes trigger review of entry pages, metadata, OpenAPI descriptions, migration pages, and MCP guidance.

### Scheduled controls

- Monthly stale-document scan for review dates, expired pricing, dead paths, and superseded terminology.
- Monthly generated-contract inventory comparison.
- Quarterly provider-intent review against the cited market references.
- Quarterly production-output inspection for Pagefind, sitemap, Markdown, and llms.txt content.

### Health indicators

Track:

- Number of unexplained contract exceptions.
- Number and age of stale documents.
- Count of backend routes missing from gateway OAS.
- Count of gateway operations missing from portal route tables.
- Count of documented MCP tools missing from exports.
- Percentage of public route families with validated examples.
- Age of pricing, product status, and provider-comparison verification.

Targets:

- Zero unexplained P0 mismatches.
- Zero expired exceptions.
- Zero undocumented public operations.
- One validated example packet per route family.
- All current-authority documents reviewed within their declared interval.

## RCCA closure criteria

The architectural RCCA is closed only when:

1. All P0 contradictions in this report are resolved or explicitly gated as unavailable.
2. The authority-by-fact model is documented.
3. The contract-cascade verifier runs locally and in CI.
4. Backend, gateway, portal, and MCP changes participate in one required verification workflow.
5. Generated Swagger drift is detected automatically.
6. Public examples and MCP tool names are validated.
7. Current, target, historical, generated, and maintainer documents are visibly distinguishable.
8. Beans and Espresso product intent is consistent across entry pages, OpenAPI descriptions, metadata, and AI-agent guidance.
9. Named owner roles and a Definition of Done are active.
10. A seeded mismatch proves each critical gate fails before deployment.

## Non-goals and residual risk

- The RCCA does not require replacing human product guidance with generated prose.
- It does not require gateway OpenAPI to become a byte-for-byte copy of backend Swagger.
- It does not require every REST route to be exposed as an MCP tool.
- It does not make external provider comparisons timeless; those remain dated research.
- It does not eliminate dynamic Espresso fields; it requires a stable documented core and explicit extension behavior.

Intentional differences remain valid when they are explicit, owned, tested, documented, and time-bounded.


# Recommended implementation order

## Wave 0: Establish the contract source of truth

Resolve I-02 through I-08 first:

1. Decide the supported Beans content-type set, including `post`.
2. Decide the Espresso pagination response shape.
3. Decide exact error, empty-result, status-code, and authentication behavior.
4. Decide which capabilities are REST-only versus MCP-enabled.
5. Update router bindings, handlers, response structs, annotations, and contract tests.
6. Regenerate backend Swagger.
7. Update gateway OpenAPI.
8. Update portal examples and route matrices.

## Wave 1: Make successful use dependable

Address U-01 through U-11:

1. Expand shared conventions.
2. Add field descriptions and examples to machine-readable contracts.
3. Add troubleshooting and complete client patterns.
4. Add operational MCP setup.
5. Explain dynamic Espresso fields and freshness semantics.
6. Add a cross-product Beans-to-Espresso workflow.
7. Link maintained Bruno examples.
8. Document versioning, support, and limits.

## Wave 2: Improve search and AI discovery

Address D-01 through D-10:

1. Correct metadata and page titles.
2. Expand the API overview and intent-to-route vocabulary.
3. Complete MCP and OpenAPI catalogs.
4. Verify or add technical SEO assets.
5. Clean repository crawler entry points.
6. Verify generated Markdown and `llms.txt` output.

# Verification checklist

After each contract cascade, run:

- `go test ./router` in `apis/beans` and `apis/espresso`.
- The relevant router contract tests, including pagination, invalid parameters, empty collections, and error responses.
- The documented Swaggo generation command for each service.
- `jq empty config/beans.oas.json config/espresso.oas.json`.
- A script comparing route registrations, annotations, generated Swagger, gateway paths, and portal route tables.
- A stale-term scan for `response_type=text`, `offset`, `/related/`, `same_as`, `derived_from`, `propagation`, old `howtos/` paths, and expired pricing dates.
- A link checker for current `docs/pages` paths and redirects.
- A built-portal check for Pagefind, sitemap/robots/canonical output, Markdown publishing, and `llms.txt` contents.
- A manual review of every cURL/JavaScript/Python example in the changed pages.

The report should be considered complete only when the runtime contract, generated contract, gateway contract, human documentation, MCP catalog, and executable examples agree.
