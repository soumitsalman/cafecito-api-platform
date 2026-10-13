---
name: beans-story-coverage
description: Compare publisher coverage of the same story with Cafecito Beans MCP. Use when a user asks who covered an event, which articles belong to a story, how many sources covered it, or how coverage changed over time.
---

# Trace a story across publishers

Beans Stories group related publisher Articles. Use Story membership to answer cross-publisher coverage questions. A `getSimilarArticles` result is related reading, not guaranteed membership in a Story; `getArticleMentions` is external discussion, not publisher coverage.

1. Call `listStories` with the topic or relevant documented filters. If the user asks for a multi-article group, `min_article_count` can narrow the collection. Select a returned Story UUID; do not invent one.
2. Call `getStory` for its title, publication bounds, `article_count`, `source_count`, and available `top_articles` previews. Those previews are not the complete member list.
3. Call `listStoryArticles` with the Story UUID when the user needs its publisher Articles. Continue with its own `pagination.next_cursor` and unchanged filters when more members are needed. Selectively call `getArticle` for a member that needs detail or available body text.
4. Present which publishers reported the Story with Article titles, publication dates, and canonical URLs when available. Use the returned counts and bounds as stated; do not infer a total from `pagination.num_results` or claim complete coverage of all publishers.

If a selected Article has a `story_id`, use that ID with `getStory` and `listStoryArticles` to inspect its assigned group. A Source record describes a publisher; an Article, Story, Source, and Mention are different resources. Empty Story collections are valid results, and coverage availability changes over time.

See the [Beans workflows](https://developer.cafecito.tech/products/beans/scenarios) and [MCP tool catalog](https://developer.cafecito.tech/guides/mcp-ai-agents).
