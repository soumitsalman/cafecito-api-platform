# Backend response cache design

Updated: 2026-10-16  
Status: Implemented using `github.com/gin-contrib/cache` v1.4.4.  
Scope: Beans and Espresso backend public data GET routes.

## Decision

Use the library's `CachePage` decorator and native Redis/in-memory persistence
stores. The user accepted dropping the custom implementation's disk storage,
semantic query keys, capacity controls, circuit breaker, operation budgets,
statistics, miss coalescing, independent waiter cancellation, and strict response
eligibility/header policies. These are no longer implementation requirements.

The previous bespoke cache under `apis/shared/cache/` and service-specific
parameter-preparation hooks have been removed. v1.4.4 retains compatibility with
the repository's Go 1.25 module while avoiding the broader dependency updates of
v1.4.5; v1.4.6 requires Go 1.26.

## Integration

1. `apis/shared/cache.go`:
   - Select the library store from environment settings and register the
     library response type for serialization.
   - Return a `ResponseCache` wrapper holding `Store`, `Path`, and `TTL`.
     `Wrap(handler)` delegates to `cache.CachePage` with the initialized TTL;
     a nil wrapper or nil store passes the handler through unchanged.
   - Use `persistence.NewInMemoryStore` when `RESPONSE_CACHE_PATH` is missing or
     blank; use `persistence.NewRedisCacheWithURL` for `redis://` or `rediss://`.
2. `apis/beans/router/routes.go` and `apis/espresso/router/routes.go`:
   - Wrap all existing public data GET handlers with `ResponseCache.Wrap`,
     after authentication, using one configured TTL for every cached route.
   - Keep health, Swagger, and preview routes unwrapped. Omitted/nil optional
     wrapper injection leaves existing router callers uncached.
   - Bind and validate normally on misses; cached responses skip the handler.
3. Service entrypoints:
   - Initialize one wrapper per process and inject it into the router.
   - Retain ordinary HTTP shutdown in `apis/shared/server.go`, independently of
     caching. Persistence lifecycle behavior belongs to the library.

## Configuration

| Setting | Default | Behavior |
|---|---|---|
| `RESPONSE_CACHE_ENABLED` | `true` | `false` disables initialization and route caching |
| `RESPONSE_CACHE_PATH` | blank | Redis connection string, or in-memory storage when absent/blank |
| `RESPONSE_CACHE_TTL` | `10m` | Positive Go duration applied uniformly to cached routes; blank also uses the default |

Invalid boolean settings, invalid/nonpositive TTLs, and malformed/non-Redis paths fail initialization with
credential-free messages. Connection establishment and further Redis URL handling
belong to the library. A configured Redis failure does not select memory instead.
`RESPONSE_CACHE_ENABLED=false` skips path/TTL validation and returns a nil wrapper.
Environment settings are captured at initialization; changes require restart.

Removed settings are `RESPONSE_CACHE_DIR`, `RESPONSE_CACHE_DISK_MAX_BYTES`,
`RESPONSE_CACHE_DISK_MAX_ENTRIES`, and `RESPONSE_CACHE_MAX_ENTRY_BYTES`. They are
not read. No disk directory or Fly Volume is used.

Each service uses a separate Redis instance or database number through its own
`RESPONSE_CACHE_PATH`. The native Redis store is used directly, including its
`Flush` behavior, which clears its configured database. Do not share that database
with unrelated data. Memory stores are independently allocated.

## Request identity and response behavior

`CachePage` derives its key from `RequestURI()`, including the concrete path and
raw query string. Parameter order, explicit defaults, spelling, and equivalent
filter values are intentionally not canonicalized. `CachePageWithoutQuery` is
not used: filters, cursors, projections, and response formats must remain distinct.

The library captures response status, headers, and bytes. Its writer caches
statuses below 300 and updates entries during writes; `CachePage` deletes an
entry if its handler aborts. There is no application-owned response-size limit,
private/cookie exclusion, or replay-header allowlist. JSON, YAML, and TOON retain
their serialized representations. A hit returns the stored `meta.as_of` and cursor.
Authentication still runs on every request; validation runs inside the handler on
misses. Gateway authentication, quotas, and metering remain outside this cache.

## Cache TTL

Every cached Beans and Espresso route uses the process's `RESPONSE_CACHE_TTL`,
defaulting to ten minutes. There are no route-specific TTLs. Each service can
configure its own duration. Both the native store default and `CachePage` use
the initialized duration. Existing Redis entries retain their expiry until replaced.

The library applies expiry. No jitter, stale-on-error serving, or ingestion-driven
invalidation is added. Entries from the previous implementation are not reused;
Redis entries expire naturally and old disk files may be removed separately.

## Accepted tradeoffs

- In-memory entries are per process/Machine, disappear on restart, and have no
  application-enforced capacity limit. Unique-query traffic can increase memory
  consumption until expiration and cleanup.
- Equivalent URLs can occupy multiple entries, reducing hit rates.
- Simultaneous misses can repeat backend work. `CachePageAtomic` is not selected
  because route-wide locking is unnecessary under the simplified requirements.
- Redis behavior uses the library's default pool/timeouts; no custom breaker,
  operation deadline, or cache statistics are retained.
- Incremental capture and concurrent misses inherit the library's response-writer
  behavior, including its append-to-existing-entry logic. No custom buffering or
  concurrency fix is layered over it.
- Header replay inherits the library's behavior, so future personalized handlers
  or request-specific response headers require deliberate cache exposure review.
- Cached data may lag upstream changes until expiry; pages are not a snapshot.

## Verification

Focused tests cover environment selection/disablement, default/explicit/invalid
TTL configuration, nil-wrapper pass-through, raw URL key distinctions,
TTL expiry, native status/header/abort behavior, Redis failure pass-through,
Redis round trips and flush behavior, all 31 public data routes,
format/projection separation, authentication, and avoided backend/embedding work.
The dropped custom requirements no longer have tests asserting their behavior.

Use `RESPONSE_CACHE_TEST_REDIS` only with a disposable Redis database. Service
integration tests use their normal fixture databases. Live-service stress tests
must be excluded when another local service is running. Commands and current
validation results are maintained in [`apis/README.md`](../README.md#backend-response-cache).

## References

- [Pinned middleware implementation](https://github.com/gin-contrib/cache/blob/v1.4.4/cache.go)
- [Pinned memory store](https://github.com/gin-contrib/cache/blob/v1.4.4/persistence/inmemory.go)
- [Pinned Redis store](https://github.com/gin-contrib/cache/blob/v1.4.4/persistence/redis.go)

## Validation record — 2026-10-16

- Centralized TTL follow-up: both service binaries build; focused `go vet` and
  formatting checks pass. Race-enabled shared cache tests pass with disposable
  Redis, including default/explicit/invalid TTLs, disabled/nil-wrapper behavior,
  configuration captured at initialization, and wrapper-driven expiry.
  Race-enabled Beans and Espresso cache tests pass across all 31 data routes.
  Broader suites were not rerun for this follow-up; the results below describe
  the earlier library replacement.
- Both service binaries build and focused `go vet` checks pass.
- Race-enabled shared tests pass with native memory and a disposable Redis 7
  instance, including expiry, URL identity, failure pass-through, serialization,
  and native flush behavior. Race-enabled service cache tests pass across
  all 18 Beans and 13 Espresso public data routes.
- Broader service suites retain 10 Beans and four Espresso top-level failures;
  the exact same failure sets reproduce in the pre-replacement source snapshot
  against the same isolated fixtures. No new suite failures were introduced.
- Beans verification uses a temporary database-only `ideology` compatibility
  overlay because its CI schema lacks a column already selected by runtime
  queries. The repository schema is unchanged. Live-service stress tests and the
  Espresso live-embedder router test are excluded.
- Public route annotations are unchanged. No deployment, provisioning, production
  performance claim, or Git action is part of this replacement.
