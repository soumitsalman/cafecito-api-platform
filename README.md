## Cafecito API Manager

[Cafecito's](https://cafecito.tech) [API](https://api.cafecito.tech) & MCP gateways and [developer portal](https://developer.cafecito.tech) is hosted on [Zuplo](https://zuplo.com/).
This monorepo contains the Zuplo gateway, Zudoku developer portal, and backend Go services for **Project Cafecito** products:

- Beans API & MCP (`services/beans/`)
- Espresso API & MCP (`services/espresso/`)
- Cortado API & MCP (future, gateway routes only)
- Latte API & MCP (future, gateway routes only)

## Repository layout

| Path | Purpose |
|------|---------|
| `config/`, `modules/` | Zuplo API gateway (TypeScript) |
| `docs/` | Zudoku developer portal |
| `services/beans/` | Beans Go API |
| `services/espresso/` | Espresso Go API |
| `services/TEI-Dockerfile` | Shared local embedder image |
| `docker-compose.yml` | Local stack: TEI + both APIs |

## Gateway (Zuplo)

```bash
npm install
npm run dev
```

Open [http://localhost:9000](http://localhost:9000).

Edit product routes in `config/`:

- `config/beans.oas.json`
- `config/espresso.oas.json`
- `config/cortado.oas.json`
- `config/latte.oas.json`

## Developer portal (Zudoku)

```bash
cd docs
npm install
npm run dev
```

Production build: `cd docs && npm run build`

## Backend services (Go)

Native run (requires `.env` with `PG_CONNECTION_STRING`, `EMBEDDER_BASE_URL`):

```bash
cd services/beans && make run    # :8080
cd services/espresso && make run # :8080 (run one at a time natively)
```

Docker (from repo root; shared TEI on `:10000`):

```bash
docker compose up --build                  # beans :8080, espresso :8081
docker compose up --build tei beansapi     # beans only
docker compose up --build tei espressoapi  # espresso only
```

Or use Makefile shortcuts: `make docker-up`, `make docker-beans`, `make docker-espresso`.

Each service needs a `services/<name>/.env` file for Docker Compose (`env_file`).

## Learn more

- [Zuplo documentation](https://zuplo.com/docs)
- [Zuplo Discord](https://discord.zuplo.com)
