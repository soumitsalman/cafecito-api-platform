# Query-Parameter Naming Design

Status: Product design decision
Scope: Beans and Espresso public APIs
Reviewed: 2026-08-18

## Decision

Use one parameter name, `q`, for unstructured user-provided text across
collection and discovery routes in both Beans and Espresso.

`q` means:

> Search this route's public searchable content using its documented relevance
> behavior.

`q` therefore means “unstructured search text,” not “semantic search.” Semantic
behavior must be explicit in the route or operation contract.

The parameter name does not promise a particular search algorithm. Depending on
the route, the service may use full-text search, tolerant matching, prefix or
contains matching, semantic retrieval, or a combination of these techniques.
The route documentation must state the searchable fields and matching behavior.

Exact and structured filters remain separate parameters, such as `domains`,
`categories`, `entities`, `regions`, and `event_types`.

## Why `q` is uniform

Users and AI clients learn one simple rule: put flexible search text in `q`.
This avoids requiring clients to guess whether a route uses `name`, `prefix`,
`search`, `prompt`, or another provider-specific name.

Uniform naming also lets common SDKs and agent tools expose the same input field
across Articles, Events, Signals, Sources, Stories, and discovery resources.

The market references show that `q` is widely used for general text search,
including PredictHQ Events and Places and Perigon article and story routes. They
also show that there is no industry-wide parameter name or algorithm for every
kind of text retrieval. No reviewed provider defines one universal parameter name
and one universal algorithm for every resource type.

## Search behavior by operation

| Operation | Input | Meaning |
|---|---|---|
| Collection or discovery `GET` route | `q` | Free-text relevance retrieval over the fields documented for that route. The route may use lexical, fuzzy, prefix, contains, semantic, or hybrid matching. |
| Explicit semantic search route | Request-body `q` | Natural-language semantic retrieval. The operation must be named or documented as semantic. |
| Explicit semantic threshold | `score_threshold` | Optional minimum semantic-match score; valid only for semantic search. |
| Exact structured filter | Route-specific parameter | Exact or enumerated matching, such as `domains`, `categories`, or `event_types`. |

Examples:

```text
GET /events?q=earnings outlook weakened
GET /sources?q=financial times
GET /tags?q=capital market
GET /entities?q=ralph lauren
POST /events/search
{
  "q": "companies whose outlook weakened because of rising input costs",
  "score_threshold": 0.72
}
```

The semantic operation is distinguished by its route and request shape, not by
renaming the query field. A semantic endpoint may also accept structured
filters in its JSON body. If one route supports both ordinary and semantic
search, it must expose an explicit `search_mode` or use separate operations.

## Why not use `name` or `prefix` as the general parameter?

`name` implies that only a canonical name field is searched. That is too narrow
for a Source search that may reasonably include a publisher name, domain,
aliases, or description.

`prefix` promises starts-with matching. It is appropriate only for a strictly
defined autocomplete operation. It excludes contains matching, token matching,
typo tolerance, and relevance ranking.

Therefore, `name` and `prefix` may be added later as explicit, narrower filters,
but they are not the primary flexible-search parameter.

## Counterarguments considered

Uniform `q` has real disadvantages:

- The same parameter can represent semantic retrieval, full-text search,
  fuzzy metadata matching, or prefix lookup.
- The parameter name does not identify the searchable fields or the ranking
  algorithm.
- Semantic search has different latency, cost, scoring, and caching behavior
  from a Source or tag lookup.
- An AI client could send a natural-language question to a discovery route that
  expects a short canonical label.
- `prefix` is clearer when an endpoint intentionally promises starts-with
  autocomplete behavior.

These objections do not outweigh uniform naming, but they require precise
operation-level documentation, strict validation of unsupported parameters, and
separate treatment for explicit semantic search.

## User and AI-client impact

- Users learn one consistent input name across both products.
- AI agents can generate the same `q` field for different search tools.
- Route descriptions still make the search mode and searchable fields explicit.
- Exact filters remain predictable and are not overloaded into natural-language
  search.
- The API avoids falsely claiming that every `q` query is semantic.
- Route descriptions identify the search mode, searchable fields, ranking behavior,
  and whether a semantic score or threshold is available. A common parameter name
  is not a substitute for that documentation.

## Industry alignment

The market report distinguishes lexical, fuzzy, and semantic retrieval; a
free-text parameter alone does not identify the algorithm. PredictHQ uses `q`
for full-text search, while Perigon uses a separate vector-search operation
with a natural-language `prompt` and relevance `score`. GDELT uses its own
`search` parameter for event and story retrieval.

Beans and Espresso use `q` as a product-level normalization of these varied
provider conventions. Provider-specific names remain compatibility references,
not additional Beans or Espresso parameters.

## Recent practitioner research

Reddit discussions created after January 2026 did not establish a consensus that
`q`, `name`, or `prefix` is the industry-standard search parameter. They did
support three related design principles:

- Query parameters are normal for collection filtering and search.
- Complex or nested searches may be better represented in a request body.
- APIs should reject unsupported parameters clearly instead of silently ignoring
  them.

The July 2026 HTTP `QUERY` discussions focused on the limitations of putting
complex search criteria into URLs, not on choosing between `q` and `name`.
The March 2026 API-validation discussion supported clear `400` errors for
unknown parameters. These are practitioner observations, not formal standards.

## Decision assessment

Uniform `q` has the stronger overall product case for Beans and Espresso:

- It gives users one consistent text-search input.
- It lets SDKs and AI tools expose the same field across resources.
- It leaves the implementation free to improve ranking without renaming the
  public parameter.

The alternative of using `q` for semantic search, `name` for Sources, and
`prefix` for discovery is stronger only when the API intentionally promises a
strict field-limited or starts-with contract. It is not the default product
experience.

## Rejected alternative

Using `q` for semantic search, `name` for Sources, and `prefix` for discovery
would make the parameter names more narrowly descriptive, but would create
different client conventions for operations that are all user-entered text
search. It would also make AI tools select parameters based on route-specific
field assumptions rather than a stable search interface.

The rejected approach remains available for narrower future filters. For
example, `prefix` may be added to a dedicated autocomplete operation when
starts-with behavior is an explicit requirement.

