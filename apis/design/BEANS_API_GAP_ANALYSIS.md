# Beans API Gap Assessment

| Field | Value |
|---|---|
| Status | **historical** |
| Authority | Snapshot of gaps versus a **superseded** route proposal |
| Audience | Maintainers reading 2026-08-24 test evidence |
| Last verified | 2026-08-25 |
| Owner role | API maintainers (archival) |
| Superseded by | [`config/beans.oas.json`](../../config/beans.oas.json); do not treat the proposal as live contract |

The P0 “collection envelope” row recorded that runtime already returned `pagination.num_results` and universal `meta.as_of`. That shape **matches the published API**. It was a gap only against the old proposal, which permitted only `limit` and `next_cursor` and restricted `meta`. Every collection probe in that campaign used `limit=5` or higher.

## Live re-run status

The requested direct live re-run could not be performed. At 2026-08-24 19:36 UTC, `http://localhost:8080/health` failed with connection refused (`HTTP_STATUS=000`), and `ss` confirmed that no process was listening on TCP port 8080. No current claim below is presented as evidence from that unavailable endpoint.

The prior 2026-08-20 live observations are historical only and are not used to assess the changed implementation. The current findings come from the fresh router-contract tests, which start `router.NewRouter` with the configured test database. They validate the same in-repo handler and response code, but are not a substitute for a direct live re-run.

## Test execution

- `go test -count=1 ./tests/...`: rerun after the new test criteria and router coverage; the contract suite is failing.
- `go test -count=1 -v ./tests -run ^TestRouter`: fails at the refreshed discovery assertions.
- Focused failures were reproduced with `-count=1`:
  - `TestRouterSearchRejectsUnsupportedParameters`
  - `TestRouterSourcesFilterByIDs`
  - `TestRouterScoreThresholdRequiresQ`
  - `TestRouterTrendingHonorsObservationWindow`
  - `TestRouterArticlesCount`

## Verified gaps

| Priority | Gap | Fresh router-test evidence |
| --- | --- | --- |
| P0 | Collection envelope and metadata are outside the proposal | Collections return `pagination.num_results` and universal `meta.as_of`. The proposal permits only `limit` and `next_cursor`, with metadata restricted to B05 and B11. |
| P0 | Public Article scope and normalization are not enforced | `content_type=post` returns HTTP 200 with five `post` Articles. Sampled records have `source: null`, nullable enrichment arrays, and in at least one case omit `story_id`; the proposal requires `news`/`blog`, a compact Source, arrays, and a nullable `story_id` key. |
| P0 | Unsupported request parameters are silently ignored | B01 returns HTTP 200 for `page`, `offset`, and `sort`, rather than the required invalid-request error. |
| P0 | Trending ignores its observation window | `/articles/trending?from=2000-01-01&to=2000-01-02&limit=5` returns five records, including Articles from 2024 and 2026, rather than an empty result. |
| P1 | Source ID filtering is ignored | `/sources?ids=<existing-id>&limit=5` returns five unrelated Sources, and a random UUID also returns five Sources. |
| P1 | Discovery payloads do not match the direct-discovery contract | Category results include `{value,type}` and title-case labels, rather than normalized `{value}` entries. They also carry the universal envelope metadata above. |
| P1 | Score-threshold validation is incomplete | B01, B03, B04, and B05 return HTTP 200 for `score_threshold=0.6` without `q`; B09 correctly returns HTTP 400. |
| P1 | Headline and Trend response shapes drift | The tested top-headlines response includes `trend`, although B04 is canonical Article only. Trend records now include `trend_score`, but expose singular `audience` instead of required `audiences`. |
| P1 | B18 Tags and B19 Count are still unavailable | `/tags` remains commented out of route registration. `/articles/count` is captured by `/articles/:id`; the B19 test receives `400 invalid UUID length: 5` instead of a count response. |

## Route status from refreshed evidence

| Route group | Status |
| --- | --- |
| B01 Search | Failing contract tests: scope, required Article normalization, envelope, unknown-parameter rejection, and score-threshold validation. |
| B03 Latest | Failing score-threshold validation and shared collection envelope contract. |
| B04 Top headlines | Now returns records in the router test environment, but fails score-threshold validation and adds a prohibited trend payload/metadata. |
| B05 Trending | Failing observation-window behavior, public Article scope/normalization, envelope, and `audience` field name. `trend_score` is now present. |
| B09–B11 Stories | Only the score-threshold negative case was freshly confirmed (passes). Positive Story/list/member behavior remains unverified against the unavailable direct live service. |
| B12 Sources | Failing Source ID filter and shared envelope/metadata contract. |
| B14–B18 Discovery | Discovery contract fails for Categories; Tags remains unregistered. The other discovery routes need direct-live confirmation once the service is restored. |
| B19 Count | Failing; route is not registered before the Article-ID route. |

## Recommended implementation order

1. Make the collection envelope route-aware and normalize public Article fields: `news`/`blog` only, required compact Source, non-null arrays, and an explicit nullable `story_id`.
2. Reject unsupported query keys and enforce the `score_threshold`/`q` rule consistently across B01, B03, B04, B05, and B09.
3. Bind and apply the B05 observation window; reject B03/B04 date bounds.
4. Pass Source `ids` into the source query, then register and implement B18 Tags and B19 Count.
5. Emit normalized direct-discovery entries, remove B04 trend data, and rename `audience` to `audiences`.
6. Start or restore the service at `localhost:8080`, then re-run the same `limit=5` probes before treating any test-environment finding as a live conclusion.
