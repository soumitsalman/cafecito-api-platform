## Cafecito API Manager

[Cafecito's](https://cafecito.tech) [API](https://api.cafecito.tech) & MCP gateways and [developer portal](https://developer.cafecito.tech) is hosted on [Zuplo](https://zuplo.com/).
This monorepo contains the Zuplo gateway, Zudoku developer portal, and backend Go apis for **Project Cafecito** products:

- Beans API & MCP (`apis/beans/`)
- Espresso API & MCP (`apis/espresso/`)
- Cortado API & MCP (future, gateway routes only)
- Latte API & MCP (future, gateway routes only)

## Repository layout

- `config/`, `modules/` — Zuplo gateway routes, policies, and handlers.
- `docs/` — Zudoku developer portal and product documentation.
- `apis/beans/` — Beans Go API.
- `apis/espresso/` — Espresso Go API.
- `apis/internal/` — shared Go embedding client and utilities.
- `apis/Dockerfile` and `apis/entrypoint.sh` — production container that runs one API with a co-located `llama-server`.
- `docker-compose.yml` — local Compose file for both API containers.

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

## Backend apis (Go)

Native run requires PostgreSQL plus an embedding endpoint. Start a local `llama-server` with embedding enabled, or point `EMBEDDER_BASE_URL` at any OpenAI-compatible embedding URL supported by the API. Then set `PG_CONNECTION_STRING` and `EMBEDDER_BASE_URL` in the service `.env` file and run:

```bash
cd apis/beans && go run .    # :8080
cd apis/espresso && go run . # :8080 (run one at a time natively)
```

The API expects the embedding endpoint to provide embeddings; a chat/completions-only URL is not sufficient.

### Docker Compose

Compose builds the shared `apis/Dockerfile` twice, once for each API. Each container also starts its own `llama-server`, so no separate embedder or `tei` service is needed:

```bash
docker compose up --build beansapi       # Beans on :8080
docker compose up --build espressoapi    # Espresso on :8081
docker compose up --build                # Both APIs
```

### Build and deploy the API image

Build from the `apis/` directory, selecting the API with `SERVICE`:

```bash
docker build --build-arg SERVICE=beans -t cafecito-beans:latest ./apis
docker build --build-arg SERVICE=espresso -t cafecito-espresso:latest ./apis
```

The Dockerfile downloads and hard-codes the `F2LLM-v2-80M.Q8_0.gguf` model, exposes the API on port `8080`, and starts the model server on the container loopback interface. Changing `MODEL_URL` at build time is supported, but runtime environment variables do not switch the baked-in model. Push the resulting image to your registry and deploy it with the platform’s container service configuration.

Each service needs an `apis/<name>/.env` file for Compose (`env_file`), including its database connection and backend API key settings.

## Learn more

- [Zuplo documentation](https://zuplo.com/docs)
- [Zuplo Discord](https://discord.zuplo.com)
