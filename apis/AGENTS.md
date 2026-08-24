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
    id UUID PRIMARY KEY,
    url VARCHAR NOT NULL,
    kind VARCHAR,    
    author VARCHAR,
    source_id UUID, -- this refers to the publishers.id if present in publishers table
    base_url VARCHAR,
    image_url VARCHAR,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    collected TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- TEXT HEAVY FIELDS
    title VARCHAR,
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
    id UUID PRIMARY KEY,
    domain_name VARCHAR NOT NULL,
    base_url VARCHAR NOT NULL,
    site_name VARCHAR,
    description TEXT,
    favicon VARCHAR,
    rss_feed VARCHAR,
    collected TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chatters (
    chatter_url VARCHAR NOT NULL,    
    url VARCHAR NOT NULL,
    bean_id UUID NOT NULL, -- this refers to an item in beans table
    platform VARCHAR, -- ex: reddit, hackernews, ycombinator, linkedin
    forum VARCHAR,
    collected TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    subscribers INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS related_beans (
  bean_id uuid NOT NULL,
  related_bean_id uuid NOT NULL,
  collected timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (bean_id, related_bean_id)
);


CREATE OR REPLACE VIEW beans_sources_view AS
SELECT
    b.*,
    p.domain_name, p.site_name, p.description, p.favicon, p.rss_feed
FROM beans b
LEFT JOIN publishers p ON b.source_id = p.id;


CREATE MATERIALIZED VIEW IF NOT EXISTS trend_aggregates AS
WITH RECURSIVE
    -- per chatter_url, the peak-engagement row at the earliest time it hit that peak;
    -- lexicographic ranking (comments, then likes) keeps the chatter even when
    -- the likes and comments maxima occur on different rows
    best_chatters AS (
        SELECT DISTINCT ON (chatter_url)
            chatter_url, bean_id, likes, comments, subscribers, collected
        FROM chatters
        ORDER BY chatter_url, comments DESC, likes DESC, collected ASC
    ),
    chatter_stats AS (
        SELECT
            bean_id,
            DATE(MAX(collected)) AS first_collected,
            SUM(likes) AS likes,
            SUM(comments) AS comments,
            SUM(subscribers) AS subscribers,
            COUNT(chatter_url) AS mentions
        FROM best_chatters
        GROUP BY bean_id
    ),
    related_stats AS (
        SELECT
            bean_id,
            COUNT(DISTINCT rel) AS related,
            DATE(MIN(collected)) AS first_collected
        FROM (
            SELECT bean_id, related_bean_id AS rel, collected FROM related_beans
            UNION ALL
            SELECT related_bean_id, bean_id, collected FROM related_beans
        ) edges
        WHERE bean_id <> rel
        GROUP BY bean_id
    ),
    -- relations are logically bidirectional but stored unidirectionally;
    -- include both directions (plus self) so every bean gets a cluster_id
    cluster_candidates AS (
        SELECT bean_id, bean_id AS cand, collected FROM related_beans
        UNION ALL
        SELECT bean_id, related_bean_id, collected FROM related_beans
        UNION ALL
        SELECT related_bean_id, bean_id, collected FROM related_beans
        UNION ALL
        SELECT related_bean_id, related_bean_id, collected FROM related_beans
    ),
    -- earliest appearance of each candidate anywhere in related_beans;
    -- frozen once set (new rows always carry a later collected)
    first_seen_related AS (
        SELECT cand, MIN(collected) AS first_seen
        FROM cluster_candidates
        GROUP BY cand
    ),
    -- winner comes from the bean's earliest (immutable) relation batch,
    -- preferring the earliest-seen candidate (the cluster seed), so the
    -- pointer is stable across refreshes and late joiners inherit the seed
    cluster_ids AS (
        SELECT DISTINCT ON (cc.bean_id)
            cc.bean_id,
            cc.cand AS cluster_id
        FROM cluster_candidates cc
        JOIN first_seen_related fs ON fs.cand = cc.cand
        ORDER BY cc.bean_id, cc.collected ASC, fs.first_seen ASC, cc.cand ASC
    ),
    -- chase each bean's pointer to its root (union-find): pointers strictly
    -- decrease by (first_seen, uuid) so chains are acyclic and end at a
    -- self-pointing seed; frozen pointers make the root equally stable
    cluster_walk AS (
        SELECT bean_id, cluster_id, 1 AS depth
        FROM cluster_ids
        UNION ALL
        SELECT w.bean_id, c.cluster_id, w.depth + 1
        FROM cluster_walk w
        JOIN cluster_ids c ON c.bean_id = w.cluster_id
        WHERE c.cluster_id <> w.cluster_id
          AND w.depth < 32    -- safety cap; chains are provably finite
    ),
    cluster_roots AS (
        SELECT DISTINCT ON (bean_id)
            bean_id,
            cluster_id
        FROM cluster_walk
        ORDER BY bean_id, depth DESC    -- deepest hop = root
    ),
    active AS (
        SELECT bean_id FROM chatter_stats
        UNION
        SELECT bean_id FROM related_stats
    ),
    trend_stats AS (
        SELECT
            a.bean_id as id,
            COALESCE(cs.likes, 0) AS likes,
            COALESCE(cs.comments, 0) AS comments,
            COALESCE(cs.subscribers, 0) AS subscribers,
            COALESCE(cs.mentions, 0) AS mentions,
            COALESCE(rs.related, 0) AS related,
            GREATEST(rs.first_collected, cs.first_collected) AS observed,
            cr.cluster_id
        FROM active a
        LEFT JOIN chatter_stats cs ON a.bean_id = cs.bean_id
        LEFT JOIN related_stats rs ON a.bean_id = rs.bean_id
        LEFT JOIN cluster_roots cr ON a.bean_id = cr.bean_id
    )
SELECT
    *,
    ((100*related + 50*comments + 10*mentions + likes) / (CURRENT_DATE + 2 - observed))::float AS trend_score
FROM trend_stats
WHERE GREATEST(likes, comments, mentions, related) > 0;


CREATE OR REPLACE VIEW latest_beans_view AS
SELECT
    b.*,
    tr.observed, tr.comments, tr.mentions, tr.likes, tr.subscribers, tr.related, tr.trend_score, tr.cluster_id
FROM beans_sources_view b
LEFT JOIN trend_aggregates tr ON b.id = tr.id;


CREATE OR REPLACE VIEW trending_beans_view AS
SELECT
    b.*,
    tr.observed, tr.comments, tr.mentions, tr.likes, tr.subscribers, tr.related, tr.trend_score, tr.cluster_id
FROM beans_sources_view b
INNER JOIN trend_aggregates tr ON b.id = tr.id;


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

CREATE UNIQUE INDEX IF NOT EXISTS idx_trend_aggregates_id ON trend_aggregates (id);
```

## Documentation dependency map

Swagger annotations, gateway OpenAPI, portal pages are separate artifacts. They need to be updated separately
For any update in public routes, params and responses
1. Update Swagger annotations in `apis/<product>/router/` as the service-local contract. After annotation changes, regenerate and commit the service's Swagger outputs; never hand-edit generated `docs/docs.go`, `docs/swagger.json`, or `docs/swagger.yaml`. Always include request, response and error type definiton for each route.
2. Update api gateway definitions `../config/<product>.oas.json`'
3. Update developer portal docs under `../docs/pages` e.g. ', `../docs/pages/products/<product>/`, and their effect on shared documents like `../docs/pages/start`, `../docs/pages/guides`. Always include sample params and responses
4. Update corresponding Bruno definitions in `<product>/tests/bruno/`

### Public documentation boundary

The public surfaces are gateway OpenAPI files in `config/`, backend Swagger artifacts, and `docs/pages/`. They describe stable client behavior, not implementation. Remove any internal design detail found there; use user-facing behavior or generic availability wording instead.

Never expose:

- database internals: Espresso's `cupboard`, `sips`, `sources`, and `relations` storage; internal relationship values such as `same_as` and `derived_from`; digest payloads and their structure; Beans' `beansack` storage; table, schema, column, foreign-key, index, view, or migration details;
- retrieval internals: embeddings, vector dimensions/indexes, HNSW, embedder/model settings, caches, query implementation, or relation direction/relationship storage;
- private implementation and operations: Go package/type/handler names, internal identifiers, backend environment variables and headers, gateway rewrites or policies, infrastructure/vendor topology, credentials, or rate-limit/quota implementation.

Events, Signals, Sources, evidence, filters, pagination, response formats, and API-key requirements are public concepts. Describe them without exposing their storage or implementation.

## Renovation Plan
Read the files in [design/](apis/design) for V1 and future renovation plan