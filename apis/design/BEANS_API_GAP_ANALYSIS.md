# Beans API Gap Assessment

Live rechecked on 2026-08-20 from 20:40-20:43 UTC against `http://0.0.0.0:8080`. Collection requests used `limit=5` or `limit=20`.

## Overall assessment

The implementation is not ready for the proposed V1 contract. Core Article, Source, Similar, and Mention routes respond, but the public data boundary, response schema, feed windows, discovery routes, Stories, and count route have material gaps.

## Latest live re-run

- `/health`: HTTP 200.
- `/articles/search?limit=5`, `/articles/latest?limit=5`, `/articles/trending?limit=5`, and `/sources?limit=5`: HTTP 200 with five records and cursors.
- Opaque search cursor continuation: HTTP 200 with five additional records.
- `/articles/search?from=1900-01-01&to=1900-01-02&limit=5`: HTTP 200 with zero records, confirming Search applies date filters.
- `/articles/latest?from=1900-01-01&to=1900-01-02&limit=5`: HTTP 200 with five 2026 records, confirming Latest ignores those bounds.
- `/articles/trending?from=1900-01-01&to=1900-01-02&limit=5`: HTTP 200 with five records, confirming Trending ignores those bounds.
- `/articles/latest?q=google&limit=5`: HTTP 400 because `score_threshold` is required.
- `/top-headlines?limit=5`: HTTP 200 with `data: []`.
- `/regions?limit=5` and `/stories?limit=5`: no response within 45 seconds during the re-run. An earlier same-day Stories probe returned HTTP 500 `db_unavailable`.
- `/tags?limit=5`: HTTP 404.
- `/articles/count?limit=5`: HTTP 400 because `count` is parsed by the Article UUID route.
- Article detail with `full_content=true`: HTTP 200; content is returned, but the sampled Article still has `source: null` and `content_type: post`.
- Similar Articles: HTTP 200 with an empty collection. Mentions: HTTP 200 with one record. Source detail: HTTP 200.
- An impossible Source UUID filter still returned five Sources.

The refreshed 20-record profile was unchanged: Search and Latest each returned `blog`, `news`, `podcast`, and `post`; `source` was null for 16/20 records, `regions` for 8/20, and `entities` for 18/20. Trending returned trend objects for 20/20 records, but `trend_score` was present for 0/20.

## Highest-priority gaps

| Priority | Gap | Live evidence |
|---|---|---|
| P0 | Article scope is not restricted to News and Blog | `/articles/search?limit=5` returned `post`; `content_type=post` is accepted. The proposal requires public queries to return only `news` and `blog`. |
| P0 | Required payload normalization is not enforced | In 20-Article samples, `source` was null for 16/20; `regions` and `entities` were frequently null. The proposal requires nested Source and enrichment arrays as `[]`. |
| P0 | Latest/trending date windows are ignored | `latest?from=1900-01-01&to=1900-01-02&limit=5` still returned 2026 Articles. Trending behaved similarly. Date-window code is commented out in `apis/beans/router/routes.go`. |
| P0 | Stories are exposed before persistent Story identity exists | `/stories?limit=5` did not complete within 45 seconds on the latest run; an earlier same-day probe returned HTTP 500 `db_unavailable`. Live Article `story_id` values were URL-like cluster keys, not stable UUIDs. |
| P0 | Regions discovery is unusable at the requested page size | `/regions?limit=5` produced no response within 45 seconds. |
| P1 | `/tags` and `/articles/count` are unavailable | `/tags?limit=5` returned 404. `/articles/count?limit=5` was captured by `/articles/:id` and returned a UUID parsing error. Registration is absent in `apis/beans/router/routes.go`. |
| P1 | Trending payload omits `trend_score` | 20 live trending records contained no `trend.trend_score`, although the proposal requires it. `toArticleDocument` does not populate it in `apis/beans/router/responses.go`. |
| P1 | Source UUID filtering is ignored | An impossible UUID filter still returned the first five Sources. `sourceListParams.IDs` is never passed into `QuerySources`. |
| P1 | Feed semantic-query contract is stricter than proposed | `/articles/latest?q=google&limit=5` and trending returned 400 because `score_threshold` was required. The proposal makes it optional when `q` is present. |
| P1 | Collection envelope has contract drift | Live collections include `pagination.num_results` and `meta.as_of` universally. The proposal requires only `limit`/`next_cursor`, and conditional metadata. |
| P2 | Discovery values are not consistently normalized | Unfiltered categories returned title-case values; entities included malformed newline-delimited strings. |
| P2 | Top headlines is operationally empty | `/top-headlines?limit=5` returned HTTP 200 with `data: []`; current data appears stale relative to the fixed 24-hour window. |

## Route status

- B01 Search: responds, but violates Article scope, Source/enrichment requirements, and envelope rules.
- B02 Article detail: responds; full content works, but payload normalization and Story identity remain gaps.
- B03 Latest: responds, but ignores proposed date behavior.
- B04 Top headlines: responds but currently returns no records.
- B05 Trending: responds, but ignores date bounds and omits `trend_score`.
- B06 Similar: responds with valid empty collections.
- B07 Mentions: responds with valid mention data.
- B09–B11 Stories: registered, but not production-ready; live collection returns 500 and IDs are derived cluster keys.
- B12 Sources: responds; UUID filtering is ignored.
- B13 Source detail: responds.
- B14 Categories, B15 Entities, B17 Sentiments: respond, with normalization-quality issues.
- B16 Regions: timeout observed.
- B18 Tags: missing, HTTP 404.
- B19 Count: missing; route falls through to Article detail parsing.

## Recommended implementation order

1. Enforce the public `news`/`blog` boundary and normalize nullable arrays, Source objects, and label values.
2. Correct latest/trending windows and make `from`/`to` bind and apply as specified.
3. Restore exact response envelopes and conditional metadata.
4. Fix trending `trend_score` and metric naming.
5. Either complete persistent Story UUID/membership support or remove B09–B11 from publication.
6. Implement `/tags` and `/articles/count`.
7. Fix Regions query performance and Source `ids` filtering.

The proposal’s explicitly future capabilities—language, structured geography, canonical entities, numeric sentiment, push delivery, and Story history—should remain separate future gaps, not blockers for this V1 renovation.
