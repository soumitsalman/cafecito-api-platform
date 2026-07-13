<proposed_plan>
# Cafecito Repository Restructuring Strategies

**Summary**

The repository is already a valid monorepo, but its identity and Go module paths still reflect the former standalone backend repositories. The safest long-term direction is to rename it around the Cafecito platform while keeping Zuplo’s gateway files at the repository root.

Zuplo treats the connected Git repository as the gateway project and expects gateway configuration such as `config/`, `modules/`, `zuplo.jsonc`, and optionally `docs/` in that project structure. Git pushes can trigger gateway deployments and branch environments. [Zuplo GitOps guidance](https://zuplo.com/blog/manage-your-apis-with-gitops), [Zuplo project structure example](https://zuplo.com/examples/basic-api-gateway)

## Strategy 1: Platform Monorepo, Minimal Migration Recommended

Repository name:

```text
cafecito-platform
```

Structure:

```text
cafecito-platform/
├── config/                 # Zuplo OpenAPI and policies
├── modules/                # Zuplo custom policies and handlers
├── docs/                   # Zudoku developer portal
├── services/
│   ├── beans/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── router/
│   │   ├── beansack/
│   │   ├── nlp/
│   │   ├── go.mod
│   │   └── Dockerfile
│   └── espresso/
│       ├── cmd/
│       ├── internal/
│       ├── router/
│       ├── cupboard/
│       ├── nlp/
│       ├── go.mod
│       └── Dockerfile
├── deploy/
├── docker-compose.yml
└── zuplo.jsonc
```

Go module paths:

```text
github.com/soumitsalman/cafecito-platform/services/beans
github.com/soumitsalman/cafecito-platform/services/espresso
```

Changes:

- Rename the GitHub repository to `cafecito-platform`.
- Update both `go.mod` files and all internal imports.
- Rename Go binaries from `beansapi` and `espressoapi` to `beans-api` and `espresso-api`.
- Optionally move service entrypoints into `services/<name>/cmd/<name>/main.go`.
- Keep `config/`, `modules/`, `docs/`, and `zuplo.jsonc` at the root.
- Keep `services/` as the backend boundary to minimize Azure, Docker Compose, GitHub Actions, and Zuplo filter changes.
- Update workflows, Dockerfiles, Makefiles, READMEs, generated Swagger metadata, and docs links.

This preserves the existing Zuplo project root and public API routes while making the repository identity and Go imports accurate.

## Strategy 2: Platform Monorepo With Explicit Gateway and API Naming

Use the same repository name, but rename the backend boundary:

```text
cafecito-platform/
├── config/
├── modules/
├── docs/
├── apis/
│   ├── beans/
│   └── espresso/
├── deploy/
├── docker-compose.yml
└── zuplo.jsonc
```

Go module paths:

```text
github.com/soumitsalman/cafecito-platform/apis/beans
github.com/soumitsalman/cafecito-platform/apis/espresso
```

Benefits:

- `apis/` communicates that these are independently deployed public API products.
- The repository becomes easier to extend with future APIs.
- The distinction between Zuplo gateway code and backend API code is clearer.

Costs:

- Every `services/...` reference must change in GitHub Actions, Docker Compose, Azure deployment contexts, documentation, Makefiles, ignore files, and Zuplo path filters.
- Azure `appSourcePath` values must be updated.
- Any configured Zuplo source-control path filtering must be reviewed and tested.
- This produces a larger diff without changing the fundamental Zuplo repository boundary.

This is a good choice only if the team strongly prefers `apis/` over `services/`.

## Strategy 3: Separate Gateway and Backend Repositories

Repositories:

```text
cafecito-gateway
cafecito-beans-api
cafecito-espresso-api
```

`cafecito-gateway` contains:

```text
config/
modules/
docs/
package.json
zuplo.jsonc
```

Each backend repository contains only one Go API, its Docker build, tests, and deployment workflow.

Benefits:

- Zuplo receives a repository containing only gateway code.
- Backend deployments become fully independent.
- Each Go module can use a clean canonical path such as:

```text
github.com/soumitsalman/beans-api
github.com/soumitsalman/espresso-api
```

Costs and risks:

- The Zuplo GitHub connection must be moved or recreated against the new gateway repository.
- Zudoku’s relative imports to `../config/*.oas.json` must be revised if docs move.
- Backend deployment secrets, Azure repository configuration, branch settings, and workflows must be migrated.
- Shared local development requires a separate orchestration repository or documented multi-repository checkout.
- Cross-repository API specification synchronization becomes a new maintenance problem.
- Zuplo branch environments and deployment checks must be revalidated after reconnection.

This is architecturally cleanest, but it is unnecessary unless independent team ownership, release cadence, or access control is important.

## Recommended Migration Sequence

1. Create a migration branch and rename the repository to `cafecito-platform`.
2. Keep the Zuplo project files at the repository root.
3. First update only Go module paths and internal imports.
4. Run backend tests and builds without moving directories.
5. Introduce `cmd/<service>/main.go` and `internal/` incrementally if desired.
6. Rename binaries and Docker image build outputs.
7. Update Azure workflows, Docker Compose, Makefiles, service documentation, and generated Swagger references.
8. Add CI checks for:
   - `go test ./...` in each service.
   - `go vet ./...`.
   - OpenAPI validation.
   - Consistency between `config/*.oas.json` and generated backend Swagger files.
   - Zuplo lint/test.
9. Confirm Zuplo’s connected repository, production branch, branch environments, and source-control path settings before merging.
10. Deploy a gateway-only change and a backend-only change separately:
    - Gateway change must deploy Zuplo.
    - Backend-only change must deploy Azure without changing gateway behavior.
11. Verify public `/beans` and `/espresso` routes, backend health checks, authentication headers, MCP routes, and developer portal rendering.
12. Rename or archive the old repository only after the new repository URL and redirects have been verified.

## Zuplo-Specific Constraints

- Keep `config/`, `modules/`, `docs/`, `package.json`, and `zuplo.jsonc` at the Zuplo-connected repository root.
- Treat relative imports such as `./modules/...` and `../config/*.oas.json` as migration-sensitive.
- Route metadata references to OpenAPI files must be updated together with file moves.
- Do not rely solely on the repository README’s `services/**` path-filter assumption. Verify the actual configured Zuplo source-control behavior in the Zuplo portal, because the official documentation confirms repository-level GitOps but does not document a general nested project-root mechanism.
- Keep public gateway paths unchanged during the migration.
- Validate that backend-only commits do not trigger unnecessary gateway deployments, if Zuplo’s configured source-control settings support that behavior.

## Assumptions

- The preferred repository identity is `cafecito-platform`.
- The repository should remain a single repository.
- Public API paths and domains remain unchanged.
- The first migration should preserve the existing Zuplo project rather than creating a new one.
- `services/beans` and `services/espresso` remain the initial backend locations unless there is a strong reason to accept the broader `apis/` path migration.
</proposed_plan>
