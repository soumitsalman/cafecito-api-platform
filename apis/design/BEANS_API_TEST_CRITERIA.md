# Beans API test criteria

| Field | Value |
|---|---|
| Status | **superseded** |
| Authority | None for live envelope; tests in `apis/beans/tests/` follow shipped behavior |
| Audience | Historical test-planning readers |
| Last verified | 2026-08-25 |
| Owner role | API maintainers (archival) |
| Superseded by | [`config/beans.oas.json`](../../config/beans.oas.json), `apis/beans/tests/` |

**Published envelope (use this, not C01/C02 below):** collections return `data`, `pagination.limit`, `pagination.num_results`, `pagination.next_cursor`, and `meta.as_of`.

The rows below were written against [BEANS_API_ROUTE_PROPOSAL.md](BEANS_API_ROUTE_PROPOSAL.md) and are **historical**. Each collection request in that campaign used `limit=5` or higher.

## Cross-route contract

| ID | Positive criteria | Negative criteria |
| --- | --- | --- |
| C01 | Collection responses are `200` envelopes with `data`, `pagination.limit`, and `pagination.next_cursor`; cursors page without duplicate IDs. Empty selections are `200` with `data: []`. | Invalid `limit`, cursor, UUID, timestamp, or structured query input is a `400` error envelope (`error.code`, `error.message`). |
| C02 | B01, B03, B04, B12, and B14–B18 omit collection `meta`; B05 and B11 include only their documented metadata. B02, B10, and B13 omit detail `meta`. | `pagination` contains no `num_results`, `page`, or `total`; routes do not accept `page`, `offset`, or `sort`. |
| C03 | Canonical Articles always include the documented fields: required title, compact Source, `story_id` (nullable UUID), array fields as arrays, and only `news` or `blog` content types. `content` is present only with `full_content=true`. | Unsupported public content types (for example `post`) and invalid filter encodings are rejected. |
| C04 | Trend-capable Article responses include nullable `likes`, `comments`, `mentions`, `audiences`, `related`, and `trend_score`; label values are normalized `snake_case`. | Non-trend routes do not add `trend`. |
| C05 | Source lists expose exactly `id`, `domain`, `name`, and `url`; source detail adds nullable metadata. Story IDs are UUIDs, list previews are canonical Articles, and story detail contains `articles`, `links`, and `stats`. | Malformed or unknown article, source, and story IDs return the documented error behavior (malformed `400`, missing `404`). |

## Route criteria

| Routes | Positive criteria | Negative criteria |
| --- | --- | --- |
| B01 `GET /articles/search` | Unfiltered browsing; text query with or without `score_threshold`; exact `ids`/`urls`; all article filters; date range; cursor; and full-content projection. | `score_threshold` without `q`; invalid date/cursor/limit; unsupported content type; `page`, `offset`, and `sort`. |
| B02 `GET /articles/{id}` | Canonical Article by UUID; `full_content=true` projection. | Malformed ID and unknown ID. |
| B03 `GET /articles/latest` | Prior-seven-day feed with permitted article filters, text query, optional score threshold, and full-content projection. | `ids`, `urls`, `from`, and `to` are rejected; score threshold requires `q`. |
| B04 `GET /top-headlines` | Prior-24-hour, news-only canonical Articles with permitted filters and no trend/meta. | `ids`, `urls`, `from`, and `to` are rejected; score threshold requires `q`. |
| B05 `GET /articles/trending` | Trend Articles with `meta.as_of`; `from`/`to` constrain the attention-observation window; text query may omit score threshold. | `ids` and `urls` are rejected; score threshold requires `q`. |
| B06–B07 Article subresources | Similar Articles honor allowed Article filters; mentions honor platforms, forums, and dates. Both return canonical envelopes. | Malformed/missing parent IDs and invalid filter values. |
| B09–B11 Stories | Browse/search stories with optional score threshold; UUID story detail; member Articles and B11 pagination/meta (`story_id`, `as_of`). | Score threshold without `q`; malformed/missing story IDs; invalid pagination/cursor. |
| B12–B13 Sources | Browse/search Sources, exact `ids` and `domains`, cursor pagination, and rich detail. | Malformed/missing source ID and invalid pagination. |
| B14–B18 Discovery | Categories, entities, regions, sentiments, and tags return normalized direct-discovery values with browse/search and cursor pagination. | Invalid pagination and unsupported pagination aliases. |
| B19 `GET /articles/count` | `from` and `to` are required; permitted scalar/CSV filters and documented `group_by` values produce count rows. | Missing/invalid bounds; invalid `group_by`; `q`, `full_content`, `score_threshold`, `cursor`, and `limit` are rejected. |
