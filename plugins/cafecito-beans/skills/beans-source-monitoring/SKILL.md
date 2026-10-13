---
name: beans-source-monitoring
description: Track a publisher or document type with Cafecito Beans MCP. Use for source-specific coverage and searches for earnings reports, financial reports, lawsuits, official statements, research, or technical documentation.
---

# Find coverage by publisher or document type

Use Beans for publisher material and Source metadata. These calls retrieve the current available data; they do not create alerts or scheduled monitoring.

## Publisher workflow

1. Call `listSources` with a publisher name or domain using documented `q` or `domains` filters. Select the matching Source UUID from `data[].id`.
2. Call `getSource` only if the user needs that Source's metadata.
3. Call `searchArticles` with `sources` set to the Source UUID and other relevant filters. A domain string belongs in a documented `domains` filter, not in `sources`.

## Document workflow

Use one of these 14 request values when the user needs a particular kind of publisher material:

| User need | Accepted `content_type` values |
| --- | --- |
| General news and blogs | `news`, `blog` |
| Company financial updates | `earnings_report`, `financial_report` |
| Legal and regulatory material | `lawsuit`, `enforcement_action`, `contract` |
| Corporate communications | `official_statement`, `press_release` |
| Research and technical material | `research_paper`, `technical_documentation`, `whitepaper` |
| Audio and publication discovery | `podcast`, `site` |

Call `searchArticles` with one accepted `content_type` value per request. Make separate calls for distinct requested types, keeping each result set and cursor separate. Omit the filter for broad public-information coverage. `post` may appear in responses but is not an accepted request value.

The filter is also available on `getLatestArticles`, `getTrendingArticles`, `getSimilarArticles`, `listStories`, and `listStoryArticles` when that tool's schema allows it. `getTopHeadlines` is news only and does not accept `content_type`.

Use `from` and `to` only on tools that accept them; they bound Article publication dates in inclusive UTC `YYYY-MM-DD` format. Do not add these parameters to latest, trending, or top-headlines feeds. Send `score_threshold` only together with `q`. Unknown or route-inapplicable parameters are rejected rather than ignored.

When continuing a collection, send its exact `pagination.next_cursor` as `cursor` with the same filters. Summarize the current results using publisher names, publication dates, Article titles, and canonical URLs when available. Use `meta.as_of` for snapshot freshness and avoid implying the Source catalog is exhaustive or that every publisher document is available.

See the [Beans overview](https://developer.cafecito.tech/products/beans), [workflows](https://developer.cafecito.tech/products/beans/scenarios), and [MCP guide](https://developer.cafecito.tech/guides/mcp-ai-agents).
