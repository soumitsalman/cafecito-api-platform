# Cafecito API Implementations
Updated: 2026-08-14

## Coding Guideline (`apis/`)

When writing or editing Go code under `apis/`, follow:

- Local variables: `lower_snake_case`
- Constants: `UPPER_SNAKE_CASE`
- Private functions: `camelCase`
- Public functions: `PascalCase`
- **Tests**: all Go tests belong in the service's `tests/` directory (e.g. `apis/beans/tests/`, `apis/espresso/tests/`). **Never** place `*_test.go` files next to production code under packages like `router/`, `cupboard/`, `beansack/`, `nlp/`, etc.

## Runtime Config

- Required: `PG_CONNECTION_STRING`
- Optional: `EMBEDDER_BASE_URL`, `EMBEDDER_API_KEY`, `EMBEDDER_MODEL`, `PORT`, `API_KEY`
- `API_KEY` format is semicolon-separated `Header=Value`.
- If `API_KEY` is unset, backend auth is disabled.

## Running Tests

Integration tests live under each service's `tests/` directory and need a reachable database. Env is loaded from that service's `.env` (at least `PG_CONNECTION_STRING`; beans also needs embedder vars for some tests).

```bash
# Beans
cd apis/beans && go test ./tests/...

# Espresso
cd apis/espresso && go test ./tests/...
```

For testing vector search, run llama-server from the repo root and set `EMBEDDER_BASE_URL=http://localhost:10000`. It stays in the foreground; success looks like `listening on http://127.0.0.1:10000` (about 1s). This local binary is CPU-only.
```bash
apis/.tools/llama-server/llama-server \
  --model apis/.models/F2LLM-v2-80M.Q8_0.gguf \
  --embedding \
  --pooling last \
  --embd-normalize 2 \
  --verbosity 3 \
  --ctx-size 16384 \
  --parallel 32 \
  --batch-size 512 \
  --ubatch-size 512 \
  --host 127.0.0.1 \
  --n-gpu-layers all \
  --port 10000
```

## Settled Decisions

- `digest->'event_type' ?| @event_types` and `impact_level` are valid PostgreSQL JSONB scalar filters.
- Do not report these filters as a gap or replace them with `->>` / `= ANY(...)`.
- Current stored kinds are only `event` and `signal`.


## Schema Definitions
### Espresso DB / Cupboard

```postgresql
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE OR REPLACE FUNCTION immutable_tags_to_text(tags text[])
RETURNS text
LANGUAGE sql
IMMUTABLE
PARALLEL SAFE
AS $$
    SELECT array_to_string(COALESCE(tags, '{}'), ' ');
$$;

CREATE TABLE IF NOT EXISTS sips (
    id UUID PRIMARY KEY,
    kind TEXT NOT NULL,
    created TIMESTAMPTZ NOT NULL,    
    source UUID,    
    embedding vector(320) NOT NULL,
    tags TEXT[], 
    tags_fts tsvector GENERATED ALWAYS AS (
        to_tsvector('simple', immutable_tags_to_text(tags))
    ) STORED,
    digest JSONB,       
    url TEXT, -- used for deriving id 
    base_url TEXT, -- used for deriving source
    ts DATE DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS sources (
    id UUID PRIMARY KEY,
    base_url TEXT NOT NULL, -- used for deriving id
    domain_name TEXT,    
    site_name TEXT,
    description TEXT,
    favicon TEXT,
    rss_feed TEXT,
    ts DATE DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS relations (
    -- NOTE: from_id and to_id are supposed to be foreign keys to sips but ignoring it to improve performance
    from_id UUID NOT NULL,
    to_id UUID NOT NULL,
    relationship TEXT NOT NULL,
    ts DATE DEFAULT CURRENT_DATE,
    UNIQUE(from_id, to_id, relationship)
);

CREATE INDEX IF NOT EXISTS idx_sips_url ON sips(url);
CREATE INDEX IF NOT EXISTS idx_sips_base_url ON sips(base_url);
CREATE INDEX IF NOT EXISTS idx_sips_kind ON sips(kind);
CREATE INDEX IF NOT EXISTS idx_sips_created ON sips(created);
CREATE INDEX IF NOT EXISTS idx_sips_source ON sips(source);
CREATE INDEX IF NOT EXISTS idx_sips_tags ON sips(tags);
CREATE INDEX IF NOT EXISTS idx_sips_tags_fts ON sips USING gin(tags_fts);
CREATE INDEX IF NOT EXISTS idx_sips_embedding_hnsw ON sips USING hnsw (embedding vector_cosine_ops) WITH (m = 24, ef_construction = 128);

CREATE INDEX IF NOT EXISTS idx_sources_base_url ON sources(base_url);

CREATE INDEX IF NOT EXISTS idx_relations_from_id ON relations(from_id);
CREATE INDEX IF NOT EXISTS idx_relations_to_id ON relations(to_id);
CREATE INDEX IF NOT EXISTS idx_relations_relationship ON relations(relationship);

```

- Events and signals are stored in `sips` table
- Relationship (SAME_AS, DERVIED_FROM) between events and signals in `relations` table
- Sources in `sources` table.
- Each `source` column in `sips` match `id` column in `sources`
- `source` column in `sips` IS NULLABLE for all signals and some events that are computed internally rather than sourced from external publishers


### Beans DB / Beansack

```postgresql
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE OR REPLACE FUNCTION immutable_tags_to_text(
    a varchar[],
    b varchar[],
    c varchar[]
)
RETURNS text
LANGUAGE sql
IMMUTABLE
PARALLEL SAFE
AS $$
    SELECT array_to_string(
        (
            SELECT array_agg(elem)
            FROM unnest(
                COALESCE(a, '{}') ||
                COALESCE(b, '{}') ||
                COALESCE(c, '{}')
            ) AS elem
            WHERE elem IS NOT NULL
        ),
        ' '
    );
$$;

-- CONTENT TABLES
CREATE TABLE IF NOT EXISTS beans (
    -- CORE FIELDS
    id UUID,
    url VARCHAR NOT NULL PRIMARY KEY,
    kind VARCHAR,
    title VARCHAR,
    author VARCHAR,
    source VARCHAR,
    image_url VARCHAR,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    collected TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- TEXT HEAVY FIELDS
    summary TEXT,
    content TEXT,
    restricted_content BOOLEAN,

    -- CLASSIFICATION FIELDS
    embedding vector(320), -- vector length is not easily mutable once set, so hardcoding it for now
    categories VARCHAR[],
    sentiments VARCHAR[],

    -- COMPRESSED EXTRACTION FIELDS
    regions VARCHAR[],
    entities VARCHAR[],

    -- TEXT SEARCH FIELD
    tags TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('simple', immutable_tags_to_text(regions, entities, categories))
    ) STORED
);

CREATE TABLE IF NOT EXISTS publishers (
    id UUID,
    source VARCHAR NOT NULL PRIMARY KEY,
    base_url VARCHAR NOT NULL,
    site_name VARCHAR,
    description TEXT,
    favicon VARCHAR,
    rss_feed VARCHAR,
    collected TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chatters (
    chatter_url VARCHAR NOT NULL,
    -- this is a foreign key to beans.url but not enforced due to insertion sequence
    url VARCHAR NOT NULL,
    source VARCHAR,
    forum VARCHAR,
    collected TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    subscribers INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS related_beans (
    url VARCHAR NOT NULL,
    related_url VARCHAR NOT NULL,
    collected TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (url, related_url)
);


CREATE MATERIALIZED VIEW IF NOT EXISTS trend_aggregates AS
WITH
    max_chatters AS (
        SELECT
            chatter_url,
            MAX(likes) as likes,
            MAX(comments) as comments
        FROM chatters
        GROUP BY chatter_url
    ),
    first_seen_max_chatters AS (
        SELECT
            fs.chatter_url,
            MIN(fs.collected) as collected
        FROM chatters fs
        LEFT JOIN max_chatters mx ON fs.chatter_url = mx.chatter_url
        WHERE fs.likes = mx.likes AND fs.comments = mx.comments
        GROUP BY fs.chatter_url
    ),
    chatter_stats AS (
        SELECT
            url,
            DATE(MAX(collected)) as updated,
            SUM(likes) as likes,
            SUM(comments) as comments,
            SUM(subscribers) as subscribers,
            COUNT(chatter_url) as shares
        FROM (
            SELECT ch.* FROM chatters ch
            LEFT JOIN first_seen_max_chatters fs ON fs.chatter_url = ch.chatter_url
            WHERE fs.collected = ch.collected
        )
        GROUP BY url
    ),
    related_stats AS (
        SELECT url, COUNT(*) AS related
        FROM related_beans
        GROUP BY url
    ),
    related_freq AS (
        SELECT related_url AS cand, COUNT(*)::int AS cnt
        FROM related_beans
        GROUP BY related_url
    ),
    cluster_candidates AS (
        SELECT url AS bean_url, url AS cand FROM related_beans
        UNION
        SELECT url, related_url FROM related_beans
    ),
    cluster_ids AS (
        SELECT DISTINCT ON (cc.bean_url)
            cc.bean_url AS url,
            cc.cand AS cluster_id
        FROM cluster_candidates cc
        LEFT JOIN related_freq rf ON rf.cand = cc.cand
        ORDER BY cc.bean_url, COALESCE(rf.cnt, 0) DESC, cc.cand
    ),
    active AS (
        SELECT url FROM chatter_stats
        UNION
        SELECT url FROM related_stats
    ),
    trend_stats AS (
        SELECT
            a.url,
            COALESCE(cg.likes, 0) as likes,
            COALESCE(cg.comments, 0) as comments,
            COALESCE(cg.subscribers, 0) as subscribers,
            COALESCE(cg.shares, 0) as shares,
            COALESCE(rg.related, 0) as related,
            GREATEST(DATE(b.created), COALESCE(cg.updated, DATE(b.created))) as updated,
            ci.cluster_id
        FROM active a
        INNER JOIN beans b ON b.url = a.url
        LEFT JOIN chatter_stats cg ON a.url = cg.url
        LEFT JOIN related_stats rg ON a.url = rg.url
        LEFT JOIN cluster_ids ci ON ci.url = a.url
    )
SELECT
    *,
    ((100*related + 50*comments + 10*shares + likes) / (CURRENT_DATE + 2 - updated))::float AS trend_score
FROM trend_stats
WHERE GREATEST(likes, comments, shares, related) > 0;


CREATE VIEW IF NOT EXISTS beans_sources_view AS
SELECT
    b.*,
    p.id AS source_id, p.base_url, p.site_name, p.description, p.favicon, p.rss_feed
FROM beans b
LEFT JOIN publishers p ON b.source = p.source;

-- PRIMARY DIFF: between latest vs trending
-- trending requires some chatter or related items. Hence INNER JOIN trend_aggregates
-- latest does not require chatter or related items. Hence LEFT JOIN trend_aggregates
CREATE VIEW IF NOT EXISTS latest_beans_view AS
SELECT
    b.*,
    tr.updated, tr.comments, tr.shares, tr.likes, tr.subscribers, tr.related, tr.trend_score, tr.cluster_id
FROM beans_sources_view b
LEFT JOIN trend_aggregates tr ON b.url = tr.url;


CREATE VIEW IF NOT EXISTS trending_beans_view AS
SELECT
    b.*,
    tr.updated, tr.comments, tr.shares, tr.likes, tr.subscribers, tr.related, tr.trend_score, tr.cluster_id
FROM beans_sources_view b
INNER JOIN trend_aggregates tr ON b.url = tr.url;


CREATE VIEW IF NOT EXISTS aggregated_beans_view AS
WITH related_groups AS (
    SELECT url, ARRAY_AGG(related_url) AS related_urls
    FROM related_beans
    GROUP BY url
)
SELECT
    b.*,
    tr.updated, tr.comments, tr.shares, tr.likes, tr.subscribers, tr.related, tr.trend_score, tr.cluster_id,
    rel.related_urls
FROM beans_sources_view b
LEFT JOIN trend_aggregates tr ON b.url = tr.url
LEFT JOIN related_groups rel ON b.url = rel.url;

-- INDEXES --
-- beans
CREATE INDEX IF NOT EXISTS idx_beans_kind ON beans(kind);
CREATE INDEX IF NOT EXISTS idx_beans_created ON beans(created DESC);
CREATE INDEX IF NOT EXISTS idx_beans_source ON beans(source);
CREATE INDEX IF NOT EXISTS idx_beans_categories ON beans USING gin(categories);
CREATE INDEX IF NOT EXISTS idx_beans_entities ON beans USING gin(entities);
CREATE INDEX IF NOT EXISTS idx_beans_regions ON beans USING gin(regions);
-- tags search
CREATE INDEX IF NOT EXISTS idx_beans_tags ON beans USING gin(tags);
-- vector search
CREATE INDEX IF NOT EXISTS idx_beans_embedding_hnsw_cosine ON beans USING hnsw (embedding vector_cosine_ops)
    WITH (m = 24, ef_construction = 128);

-- publishers
CREATE INDEX IF NOT EXISTS idx_publishers_source ON publishers(source);

-- chatters
CREATE INDEX IF NOT EXISTS idx_chatters_url ON chatters(url);
CREATE INDEX IF NOT EXISTS idx_chatters_collected ON chatters(collected DESC);

-- related_beans
CREATE INDEX IF NOT EXISTS idx_related_beans_related_url ON related_beans(related_url);
CREATE INDEX IF NOT EXISTS idx_related_beans_collected ON related_beans(collected DESC);
CREATE INDEX IF NOT EXISTS idx_chatters_chatter_url ON chatters(chatter_url);

CREATE UNIQUE INDEX IF NOT EXISTS idx_trend_agg_url ON trend_aggregates(url);
```

## Renovation Plan
Read the files in [design/](apis/design) for V1 and future renovation plan

## Documentation dependency map

Swagger annotations, gateway OpenAPI, and portal pages are separate artifacts. Do not assume that generating Swagger updates the public gateway contract or portal documentation.

- Treat route behavior, request/response types, and Swagger annotations in `apis/<service>/router/` as the service-local contract. After annotation changes, regenerate and commit the service's Swagger outputs; never hand-edit generated `docs/docs.go`, `docs/swagger.json`, or `docs/swagger.yaml`.
  - Espresso: from `apis/espresso/`, run `go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g router/routes.go -o docs --parseDependency --parseInternal` after changing `router/routes.go` docstrings or referenced router types.
  - Beans: use the documented `swag` command in `apis/beans/README.md` after its router annotation changes.
- For any public Espresso behavior change, then manually update `../config/espresso.oas.json`. It is the gateway contract under `/espresso`, powers `/api/espresso` via `../docs/zudoku.config.tsx`, and supplies the `/espresso/mcp` server plus its exported `operationId` tools. Keep public paths, parameters, schemas, descriptions, errors, and MCP tool mappings in sync with the backend.
- Reflect that behavior in the relevant portal pages: `../docs/pages/products/espresso/overview.mdx`, `../docs/pages/products/espresso/workflows.mdx`, and `../docs/pages/products/espresso/migration.mdx`; `../docs/pages/guides/mcp-ai-agents.mdx` for MCP changes; and `../docs/pages/guides/api-conventions.mdx` or `../docs/pages/start/first-api-call.mdx` for shared or quickstart behavior. Update `../docs/zudoku.config.tsx` only when API mounting, navigation, redirects, or page structure changes.
- Apply the same generated-Swagger → `../config/beans.oas.json` → `../docs/pages/products/beans/overview.mdx` review for Beans.

### Public documentation boundary

The public surfaces are gateway OpenAPI files in `config/`, backend Swagger artifacts, and `docs/pages/`. They describe stable client behavior, not implementation. Remove any internal design detail found there; use user-facing behavior or generic availability wording instead.

Never expose:

- database internals: Espresso's `cupboard`, `sips`, `sources`, and `relations` storage; internal relationship values such as `same_as` and `derived_from`; digest payloads and their structure; Beans' `beansack` storage; table, schema, column, foreign-key, index, view, or migration details;
- retrieval internals: embeddings, vector dimensions/indexes, HNSW, embedder/model settings, caches, query implementation, or relation direction/relationship storage;
- private implementation and operations: Go package/type/handler names, internal identifiers, backend environment variables and headers, gateway rewrites or policies, infrastructure/vendor topology, credentials, or rate-limit/quota implementation.

Events, Signals, Sources, evidence, filters, pagination, response formats, and API-key requirements are public concepts. Describe them without exposing their storage or implementation.
