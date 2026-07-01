# Monorepo migration: GitHub secrets and Azure

After merging backend deploy workflows into this repo, copy repository secrets from the old standalone repos into **cafecito-api-manager** (Settings → Secrets and variables → Actions).

## Beans (`deploy-beans.yml`)

| Secret | Source repo |
|--------|-------------|
| `CAFECITOBEANSAPI_AZURE_CLIENT_ID` | go-beans-api |
| `CAFECITOBEANSAPI_AZURE_TENANT_ID` | go-beans-api |
| `CAFECITOBEANSAPI_AZURE_SUBSCRIPTION_ID` | go-beans-api |
| `CAFECITOBEANSAPI_REGISTRY_USERNAME` | go-beans-api |
| `CAFECITOBEANSAPI_REGISTRY_PASSWORD` | go-beans-api |

## Espresso (`deploy-espresso.yml`)

| Secret | Source repo |
|--------|-------------|
| `CAFECITOESPRESSOAPI_AZURE_CLIENT_ID` | go-espresso-api |
| `CAFECITOESPRESSOAPI_AZURE_TENANT_ID` | go-espresso-api |
| `CAFECITOESPRESSOAPI_AZURE_SUBSCRIPTION_ID` | go-espresso-api |
| `CAFECITOESPRESSOAPI_REGISTRY_USERNAME` | go-espresso-api |
| `CAFECITOESPRESSOAPI_REGISTRY_PASSWORD` | go-espresso-api |

## Azure Deployment Center

For each Container App in Azure Portal:

1. **cafecito-beans-api** — set GitHub repo to `soumitsalman/cafecito-api-manager`, branch `main`, app path `services/beans` (if the portal exposes a path setting).
2. **cafecito-espresso-api** — same repo, app path `services/espresso`.

Workflow `appSourcePath` is the authoritative Docker build context; keep Azure Deployment Center aligned to avoid regenerated workflow drift.

## Zuplo path filters

In Zuplo project → Source Control, exclude `services/**` from deploy triggers so backend-only pushes do not redeploy the gateway.

## First deploy after merge

Backend workflows will not auto-run until `services/beans/**` or `services/espresso/**` changes. Use **workflow_dispatch** on `deploy-beans.yml` and `deploy-espresso.yml` for an initial smoke deploy after secrets are copied.
