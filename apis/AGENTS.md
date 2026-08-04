## Coding Guideline (`apis/`)

When writing or editing Go code under `apis/`, follow:

- Local variables: `lower_snake_case`
- Constants: `UPPER_SNAKE_CASE`
- Private functions: `camelCase`
- Public functions: `PascalCase`
- **Tests**: all Go tests belong in the service's `tests/` directory (e.g. `apis/beans/tests/`, `apis/espresso/tests/`). **Never** place `*_test.go` files next to production code under packages like `router/`, `cupboard/`, `beansack/`, `nlp/`, etc.