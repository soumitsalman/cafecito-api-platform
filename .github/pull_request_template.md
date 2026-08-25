## Summary

<!-- What changed and why. -->

## API-change checklist

Complete every item that applies. Leave N/A with a reason when a surface is untouched.

- [ ] Runtime behavior and tests (`apis/<service>/router/`, `apis/<service>/tests/`)
- [ ] Swagger annotations updated
- [ ] Generated Swagger regenerated (`swag` from [`apis/README.md`](../apis/README.md); no hand-edits)
- [ ] Gateway OpenAPI reconciled (`config/beans.oas.json` / `config/espresso.oas.json`)
- [ ] MCP export / documented tools reconciled
- [ ] Portal pages and examples reconciled (`docs/pages/`)
- [ ] Product intent and not-for language still correct
- [ ] Frozen public policy still true: free tier **100/min** and **50,000/month**; **Bearer** except **health**; **no private backend header** in public docs
- [ ] `npm run verify:api-contracts` and `npm run verify:docs` pass locally (optional: run **API Contract & Docs** via `workflow_dispatch`; it does not gate CI)

## Definition of Done

A public API change is done only when the cascade is closed: tests, annotations, generated Swagger, gateway OAS, MCP, portal examples, and local `verify:api-contracts` / `verify:docs`.

## Test plan

- [ ]
