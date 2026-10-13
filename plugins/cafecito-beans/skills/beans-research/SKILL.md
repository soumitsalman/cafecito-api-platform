---
name: beans-research
description: Find what publishers have reported with Cafecito Beans MCP. Use for article research, recent coverage, headlines, trends, related reading, and external mentions across news, blogs, reports, research, and official statements.
---

# Research publisher material with Beans

Beans is a public-information discovery service. Its Articles include news, blogs, financial and earnings reports, litigation, official statements, research, technical documents, and other publisher material. Coverage and update timing vary by source. Use the `cafecito-beans` MCP tools for publisher reading; Beans does not provide Espresso Events, Signals, impact, or forecasts.

## Pick the first tool by intent

| User question | MCP tool | Meaning |
| --- | --- | --- |
| What did publishers report about a topic? | `searchArticles` | Relevance and documented filters. |
| What was published most recently? | `getLatestArticles` | Publication chronology. |
| What news is prominent in the recent headline window? | `getTopHeadlines` | Recent, attention-ranked news. |
| Which articles are gaining attention? | `getTrendingArticles` | Attention-ranked coverage. |
| What did one selected article say? | `getArticle` | Detail by Article UUID; request `full_content=true` only if the body is needed. |
| What other publisher reading is related? | `getSimilarArticles` | Related Articles for a selected Article UUID. |
| Where was this Article URL observed externally? | `getArticleMentions` | External mentions, not more publisher Articles. |

Use `content_type=news` on `getLatestArticles` or `getTrendingArticles` when the user wants news only. `getTopHeadlines` is already news only; do not send `content_type` to it. For the complete set of 14 accepted content types and document-focused workflows, use the `beans-source-monitoring` skill. Use `listCategories`, `listEntities`, `listRegions`, or `listSentiments` only when an accepted filter value is unknown. Reuse a returned `value` in an applicable Article filter.

## Research sequence

1. Start with one collection tool and the smallest useful filter set. Use `q` for topic relevance; send `score_threshold` only with `q`. Date filters `from` and `to` apply to Article publication time where the tool accepts them. Do not send dates to latest, trending, or top-headlines feeds.
2. Select Article UUIDs from the returned `data`. Open only the Articles needed to answer. `full_content=true` may still return no body; use the Article's canonical `url` for attribution.
3. Call `getSimilarArticles` only for related publisher reading and `getArticleMentions` only for external observations. Neither establishes Story membership; use the `beans-story-coverage` skill for that.
4. Report the source, publication date, and canonical Article URL when present. Explain ranking accurately: latest is chronological; headlines and trending are attention-ranked. Show `meta.as_of` when freshness matters. Attribute claims to publishers and separate your inference.

Collections return `{data, pagination, meta}`. `pagination.num_results` counts this page, not all matches. An empty `data: []` is a valid result. Continue only when needed: send `pagination.next_cursor` unchanged as `cursor` to the same tool with the same filters; stop when it is `null`. Do not invent UUIDs, cursor values, parameters, or tool names. A returned Article body does not grant redistribution or other downstream rights.

The 16 exported Beans tools and current request schemas are listed in the [MCP guide](https://developer.cafecito.tech/guides/mcp-ai-agents) and [Beans API reference](https://developer.cafecito.tech/api/beans).
