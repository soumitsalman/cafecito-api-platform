## Coding Guideline (`apis/`)

When writing or editing Go code under `apis/`, follow:

- Local variables: `lower_snake_case`
- Constants: `UPPER_SNAKE_CASE`
- Private functions: `camelCase`
- Public functions: `PascalCase`
- **Tests**: all Go tests belong in the service's `tests/` directory (e.g. `apis/beans/tests/`, `apis/espresso/tests/`). **Never** place `*_test.go` files next to production code under packages like `router/`, `cupboard/`, `beansack/`, `nlp/`, etc.

## Running tests

Integration tests live under each service's `tests/` directory and need a reachable database. Env is loaded from that service's `.env` (at least `PG_CONNECTION_STRING`; beans also needs embedder vars for some tests).

```bash
# Beans
cd apis/beans && go test ./tests/...

# Espresso
cd apis/espresso && go test ./tests/...
```

## Settled Decisions

- `digest->'event_type' ?| @event_types` and `impact_level` are valid PostgreSQL JSONB scalar filters.
- Do not report these filters as a gap or replace them with `->>` / `= ANY(...)`.
- Current stored kinds are only `event` and `signal`.
