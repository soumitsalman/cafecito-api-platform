# Cafecito API Implementations
Updated: 2026-08-25

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

The manual **API Contract & Docs** workflow (`.github/workflows/api-contract-docs.yml`, `workflow_dispatch` only) starts `pgvector/pgvector:pg16`, applies `.github/ci/schema-beans.sql` or `schema-espresso.sql`, seeds Espresso from `.github/ci/seed-espresso.sql`, and relies on Beans `tests/fixtures_test.go` for deterministic rows. It does not run on PRs or gate deploys. Router vector-search tests that call a live embedder are skipped in that workflow (`-skip 'TestRouterVectorSearch'`). Stress tests skip when no process is listening on `:8080`.

```bash
# Beans (seeds CI fixtures; fake embedder)
cd apis/beans && go test ./tests/...

# Espresso (hermetic contract tests always; DB tests need PG + fixtures)
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

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;
COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';
CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;
COMMENT ON EXTENSION vector IS 'vector data type and ivfflat and hnsw access methods';


CREATE FUNCTION public.immutable_tags_to_text(tags text[]) RETURNS text
    LANGUAGE sql IMMUTABLE PARALLEL SAFE
    AS $$
    SELECT array_to_string(COALESCE(tags, '{}'), ' ');
$$;


--
-- Name: show_db_tree(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.show_db_tree() RETURNS TABLE(tree_structure text)
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- First show all databases
    RETURN QUERY
    SELECT ':file_folder: ' || datname || ' (DATABASE)'
    FROM pg_database 
    WHERE datistemplate = false;

    -- Then show current database structure
    RETURN QUERY
    WITH RECURSIVE 
    -- Get schemas
    schemas AS (
        SELECT 
            n.nspname AS object_name,
            1 AS level,
            n.nspname AS path,
            'SCHEMA' AS object_type
        FROM pg_namespace n
        WHERE n.nspname NOT LIKE 'pg_%' 
        AND n.nspname != 'information_schema'
    ),

    -- Get all objects (tables, views, functions, etc.)
    objects AS (
        SELECT 
            c.relname AS object_name,
            2 AS level,
            s.path || ' → ' || c.relname AS path,
            CASE c.relkind
                WHEN 'r' THEN 'TABLE'
                WHEN 'v' THEN 'VIEW'
                WHEN 'm' THEN 'MATERIALIZED VIEW'
                WHEN 'i' THEN 'INDEX'
                WHEN 'S' THEN 'SEQUENCE'
                WHEN 'f' THEN 'FOREIGN TABLE'
            END AS object_type
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        JOIN schemas s ON n.nspname = s.object_name
        WHERE c.relkind IN ('r','v','m','i','S','f')

        UNION ALL

        SELECT 
            p.proname AS object_name,
            2 AS level,
            s.path || ' → ' || p.proname AS path,
            'FUNCTION' AS object_type
        FROM pg_proc p
        JOIN pg_namespace n ON n.oid = p.pronamespace
        JOIN schemas s ON n.nspname = s.object_name
    ),

    -- Combine schemas and objects
    combined AS (
        SELECT * FROM schemas
        UNION ALL
        SELECT * FROM objects
    )

    -- Final output with tree-like formatting
    SELECT 
        REPEAT('    ', level) || 
        CASE 
            WHEN level = 1 THEN '└── :open_file_folder: '
            ELSE '    └── ' || 
                CASE object_type
                    WHEN 'TABLE' THEN ':bar_chart: '
                    WHEN 'VIEW' THEN ':eye: '
                    WHEN 'MATERIALIZED VIEW' THEN ':newspaper: '
                    WHEN 'FUNCTION' THEN ':zap: '
                    WHEN 'INDEX' THEN ':mag: '
                    WHEN 'SEQUENCE' THEN ':1234: '
                    WHEN 'FOREIGN TABLE' THEN ':globe_with_meridians: '
                    ELSE ''
                END
        END || object_name || ' (' || object_type || ')'
    FROM combined
    ORDER BY path;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: relations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.relations (
    from_id uuid NOT NULL,
    to_id uuid NOT NULL,
    relationship text NOT NULL,
    ts date DEFAULT CURRENT_DATE
);


--
-- Name: sips; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sips (
    id uuid NOT NULL,
    kind text NOT NULL,
    created timestamp with time zone NOT NULL,
    source uuid,
    embedding public.vector(320),
    tags text[],
    tags_fts tsvector GENERATED ALWAYS AS (to_tsvector('simple'::regconfig, public.immutable_tags_to_text(tags))) STORED,
    digest jsonb,
    url text,
    base_url text,
    ts date DEFAULT CURRENT_DATE
);


--
-- Name: sources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sources (
    id uuid NOT NULL,
    base_url text NOT NULL,
    domain_name text,
    site_name text,
    description text,
    favicon text,
    rss_feed text,
    ts date DEFAULT CURRENT_DATE
);


--
-- Name: relations relations_from_id_to_id_relationship_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.relations
    ADD CONSTRAINT relations_from_id_to_id_relationship_key UNIQUE (from_id, to_id, relationship);

ALTER TABLE ONLY public.sips
    ADD CONSTRAINT sips_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.sources
    ADD CONSTRAINT sources_pkey PRIMARY KEY (id);


CREATE INDEX idx_relations_from_id ON public.relations USING btree (from_id);

CREATE INDEX idx_relations_relationship ON public.relations USING btree (relationship);

CREATE INDEX idx_relations_to_id ON public.relations USING btree (to_id);

CREATE INDEX idx_sips_base_url ON public.sips USING btree (base_url);

CREATE INDEX idx_sips_created ON public.sips USING btree (created);

CREATE INDEX idx_sips_embedding_hnsw ON public.sips USING hnsw (embedding public.vector_cosine_ops) WITH (m='24', ef_construction='128');

CREATE INDEX idx_sips_kind ON public.sips USING btree (kind);

CREATE INDEX idx_sips_source ON public.sips USING btree (source);

CREATE INDEX idx_sips_tags ON public.sips USING btree (tags);

CREATE INDEX idx_sips_tags_fts ON public.sips USING gin (tags_fts);

CREATE INDEX idx_sips_url ON public.sips USING btree (url);

CREATE INDEX idx_sources_base_url ON public.sources USING btree (base_url);
```

- Events and signals are stored in `sips` table
- Relationship (SAME_AS, DERVIED_FROM) between events and signals in `relations` table
- Sources in `sources` table.
- Each `source` column in `sips` match `id` column in `sources`
- `source` column in `sips` IS NULLABLE for all signals and some events that are computed internally rather than sourced from external publishers


### Beans DB / Beansack

```sql
--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- Name: vector; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;


--
-- Name: EXTENSION vector; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION vector IS 'vector data type and ivfflat and hnsw access methods';


--
-- Name: immutable_tags_to_text(character varying[], character varying[], character varying[]); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.immutable_tags_to_text(a character varying[], b character varying[], c character varying[]) RETURNS text
    LANGUAGE sql IMMUTABLE PARALLEL SAFE
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


--
-- Name: show_db_tree(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.show_db_tree() RETURNS TABLE(tree_structure text)
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- First show all databases
    RETURN QUERY
    SELECT ':file_folder: ' || datname || ' (DATABASE)'
    FROM pg_database 
    WHERE datistemplate = false;

    -- Then show current database structure
    RETURN QUERY
    WITH RECURSIVE 
    -- Get schemas
    schemas AS (
        SELECT 
            n.nspname AS object_name,
            1 AS level,
            n.nspname AS path,
            'SCHEMA' AS object_type
        FROM pg_namespace n
        WHERE n.nspname NOT LIKE 'pg_%' 
        AND n.nspname != 'information_schema'
    ),

    -- Get all objects (tables, views, functions, etc.)
    objects AS (
        SELECT 
            c.relname AS object_name,
            2 AS level,
            s.path || ' → ' || c.relname AS path,
            CASE c.relkind
                WHEN 'r' THEN 'TABLE'
                WHEN 'v' THEN 'VIEW'
                WHEN 'm' THEN 'MATERIALIZED VIEW'
                WHEN 'i' THEN 'INDEX'
                WHEN 'S' THEN 'SEQUENCE'
                WHEN 'f' THEN 'FOREIGN TABLE'
            END AS object_type
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        JOIN schemas s ON n.nspname = s.object_name
        WHERE c.relkind IN ('r','v','m','i','S','f')

        UNION ALL

        SELECT 
            p.proname AS object_name,
            2 AS level,
            s.path || ' → ' || p.proname AS path,
            'FUNCTION' AS object_type
        FROM pg_proc p
        JOIN pg_namespace n ON n.oid = p.pronamespace
        JOIN schemas s ON n.nspname = s.object_name
    ),

    -- Combine schemas and objects
    combined AS (
        SELECT * FROM schemas
        UNION ALL
        SELECT * FROM objects
    )

    -- Final output with tree-like formatting
    SELECT 
        REPEAT('    ', level) || 
        CASE 
            WHEN level = 1 THEN '└── :open_file_folder: '
            ELSE '    └── ' || 
                CASE object_type
                    WHEN 'TABLE' THEN ':bar_chart: '
                    WHEN 'VIEW' THEN ':eye: '
                    WHEN 'MATERIALIZED VIEW' THEN ':newspaper: '
                    WHEN 'FUNCTION' THEN ':zap: '
                    WHEN 'INDEX' THEN ':mag: '
                    WHEN 'SEQUENCE' THEN ':1234: '
                    WHEN 'FOREIGN TABLE' THEN ':globe_with_meridians: '
                    ELSE ''
                END
        END || object_name || ' (' || object_type || ')'
    FROM combined
    ORDER BY path;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: beans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.beans (
    id uuid NOT NULL,
    url character varying NOT NULL,
    kind character varying,
    title character varying,
    author character varying,
    image_url character varying,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    collected timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    summary text,
    content text,
    restricted_content boolean,
    embedding public.vector(320),
    categories character varying[],
    sentiments character varying[],
    regions character varying[],
    entities character varying[],
    tags tsvector GENERATED ALWAYS AS (to_tsvector('simple'::regconfig, public.immutable_tags_to_text(regions, entities, categories))) STORED,
    source_id uuid,
    base_url character varying
);


--
-- Name: publishers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.publishers (
    id uuid NOT NULL,
    domain_name character varying CONSTRAINT publishers_source_not_null NOT NULL,
    base_url character varying NOT NULL,
    site_name character varying,
    description text,
    favicon character varying,
    rss_feed character varying,
    collected timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: beans_sources_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.beans_sources_view AS
 SELECT b.id,
    b.url,
    b.kind,
    b.title,
    b.author,
    b.image_url,
    b.created,
    b.collected,
    b.summary,
    b.content,
    b.restricted_content,
    b.embedding,
    b.categories,
    b.sentiments,
    b.regions,
    b.entities,
    b.tags,
    b.source_id,
    b.base_url,
    p.domain_name,
    p.site_name,
    p.description,
    p.favicon,
    p.rss_feed
   FROM (public.beans b
     LEFT JOIN public.publishers p ON ((b.source_id = p.id)));


--
-- Name: chatters; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chatters (
    chatter_url character varying NOT NULL,
    url character varying NOT NULL,
    source character varying,
    forum character varying,
    collected timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    likes integer DEFAULT 0,
    comments integer DEFAULT 0,
    subscribers integer DEFAULT 0,
    shares integer DEFAULT 0,
    bean_id uuid,
    platform character varying
);


--
-- Name: related_beans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.related_beans (
    bean_id uuid CONSTRAINT related_beans_v2_bean_id_not_null NOT NULL,
    related_bean_id uuid CONSTRAINT related_beans_v2_related_bean_id_not_null NOT NULL,
    collected timestamp without time zone DEFAULT CURRENT_TIMESTAMP CONSTRAINT related_beans_v2_collected_not_null NOT NULL
);


--
-- Name: trend_aggregates; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.trend_aggregates AS
 WITH RECURSIVE best_chatters AS (
         SELECT DISTINCT ON (chatters.chatter_url) chatters.chatter_url,
            chatters.bean_id,
            chatters.likes,
            chatters.comments,
            chatters.subscribers,
            chatters.collected
           FROM public.chatters
          ORDER BY chatters.chatter_url, chatters.comments DESC, chatters.likes DESC, chatters.collected
        ), chatter_stats AS (
         SELECT best_chatters.bean_id,
            date(max(best_chatters.collected)) AS first_collected,
            sum(best_chatters.likes) AS likes,
            sum(best_chatters.comments) AS comments,
            sum(best_chatters.subscribers) AS subscribers,
            count(best_chatters.chatter_url) AS mentions
           FROM best_chatters
          GROUP BY best_chatters.bean_id
        ), related_stats AS (
         SELECT edges.bean_id,
            count(DISTINCT edges.rel) AS related,
            date(min(edges.collected)) AS first_collected
           FROM ( SELECT related_beans.bean_id,
                    related_beans.related_bean_id AS rel,
                    related_beans.collected
                   FROM public.related_beans
                UNION ALL
                 SELECT related_beans.related_bean_id,
                    related_beans.bean_id,
                    related_beans.collected
                   FROM public.related_beans) edges
          WHERE (edges.bean_id <> edges.rel)
          GROUP BY edges.bean_id
        ), cluster_candidates AS (
         SELECT related_beans.bean_id,
            related_beans.bean_id AS cand,
            related_beans.collected
           FROM public.related_beans
        UNION ALL
         SELECT related_beans.bean_id,
            related_beans.related_bean_id,
            related_beans.collected
           FROM public.related_beans
        UNION ALL
         SELECT related_beans.related_bean_id,
            related_beans.bean_id,
            related_beans.collected
           FROM public.related_beans
        UNION ALL
         SELECT related_beans.related_bean_id,
            related_beans.related_bean_id,
            related_beans.collected
           FROM public.related_beans
        ), first_seen_related AS (
         SELECT cluster_candidates.cand,
            min(cluster_candidates.collected) AS first_seen
           FROM cluster_candidates
          GROUP BY cluster_candidates.cand
        ), cluster_ids AS (
         SELECT DISTINCT ON (cc.bean_id) cc.bean_id,
            cc.cand AS cluster_id
           FROM (cluster_candidates cc
             JOIN first_seen_related fs ON ((fs.cand = cc.cand)))
          ORDER BY cc.bean_id, cc.collected, fs.first_seen, cc.cand
        ), cluster_walk AS (
         SELECT cluster_ids.bean_id,
            cluster_ids.cluster_id,
            1 AS depth
           FROM cluster_ids
        UNION ALL
         SELECT w.bean_id,
            c.cluster_id,
            (w.depth + 1)
           FROM (cluster_walk w
             JOIN cluster_ids c ON ((c.bean_id = w.cluster_id)))
          WHERE ((c.cluster_id <> w.cluster_id) AND (w.depth < 32))
        ), cluster_roots AS (
         SELECT DISTINCT ON (cluster_walk.bean_id) cluster_walk.bean_id,
            cluster_walk.cluster_id
           FROM cluster_walk
          ORDER BY cluster_walk.bean_id, cluster_walk.depth DESC
        ), active AS (
         SELECT chatter_stats.bean_id
           FROM chatter_stats
        UNION
         SELECT related_stats.bean_id
           FROM related_stats
        ), trend_stats AS (
         SELECT a.bean_id AS id,
            COALESCE(cs.likes, (0)::bigint) AS likes,
            COALESCE(cs.comments, (0)::bigint) AS comments,
            COALESCE(cs.subscribers, (0)::bigint) AS subscribers,
            COALESCE(cs.mentions, (0)::bigint) AS mentions,
            COALESCE(rs.related, (0)::bigint) AS related,
            GREATEST(rs.first_collected, cs.first_collected) AS observed,
            cr.cluster_id
           FROM (((active a
             LEFT JOIN chatter_stats cs ON ((a.bean_id = cs.bean_id)))
             LEFT JOIN related_stats rs ON ((a.bean_id = rs.bean_id)))
             LEFT JOIN cluster_roots cr ON ((a.bean_id = cr.bean_id)))
        )
 SELECT id,
    likes,
    comments,
    subscribers,
    mentions,
    related,
    observed,
    cluster_id,
    ((((((100 * related) + (50 * comments)) + (10 * mentions)) + likes) / ((CURRENT_DATE + 2) - observed)))::double precision AS trend_score
   FROM trend_stats
  WHERE (GREATEST(likes, comments, mentions, related) > 0)
  WITH NO DATA;


--
-- Name: latest_beans_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.latest_beans_view AS
 SELECT b.id,
    b.url,
    b.kind,
    b.title,
    b.author,
    b.image_url,
    b.created,
    b.collected,
    b.summary,
    b.content,
    b.restricted_content,
    b.embedding,
    b.categories,
    b.sentiments,
    b.regions,
    b.entities,
    b.tags,
    b.source_id,
    b.base_url,
    b.domain_name,
    b.site_name,
    b.description,
    b.favicon,
    b.rss_feed,
    tr.observed,
    tr.comments,
    tr.mentions,
    tr.likes,
    tr.subscribers,
    tr.related,
    tr.trend_score,
    tr.cluster_id
   FROM (public.beans_sources_view b
     LEFT JOIN public.trend_aggregates tr ON ((b.id = tr.id)));


--
-- Name: trending_beans_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.trending_beans_view AS
 SELECT b.id,
    b.url,
    b.kind,
    b.title,
    b.author,
    b.image_url,
    b.created,
    b.collected,
    b.summary,
    b.content,
    b.restricted_content,
    b.embedding,
    b.categories,
    b.sentiments,
    b.regions,
    b.entities,
    b.tags,
    b.source_id,
    b.base_url,
    b.domain_name,
    b.site_name,
    b.description,
    b.favicon,
    b.rss_feed,
    tr.observed,
    tr.comments,
    tr.mentions,
    tr.likes,
    tr.subscribers,
    tr.related,
    tr.trend_score,
    tr.cluster_id
   FROM (public.beans_sources_view b
     JOIN public.trend_aggregates tr ON ((b.id = tr.id)));

ALTER TABLE ONLY public.beans
    ADD CONSTRAINT beans_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.publishers
    ADD CONSTRAINT publishers_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.related_beans
    ADD CONSTRAINT related_beans_v2_bean_id_related_bean_id_key UNIQUE (bean_id, related_bean_id);

CREATE INDEX idx_beans_categories ON public.beans USING gin (categories);

CREATE INDEX idx_beans_created ON public.beans USING btree (created DESC);

CREATE INDEX idx_beans_embedding_hnsw_cosine ON public.beans USING hnsw (embedding public.vector_cosine_ops) WITH (m='24', ef_construction='128');

CREATE INDEX idx_beans_entities ON public.beans USING gin (entities);

CREATE INDEX idx_beans_kind ON public.beans USING btree (kind);

CREATE INDEX idx_beans_regions ON public.beans USING gin (regions);

CREATE INDEX idx_beans_source_id ON public.beans USING btree (source_id);

CREATE INDEX idx_beans_tags ON public.beans USING gin (tags);

CREATE INDEX idx_beans_url ON public.beans USING btree (url);

CREATE INDEX idx_chatters_bean_id ON public.chatters USING btree (bean_id);

CREATE INDEX idx_chatters_chatter_url ON public.chatters USING btree (chatter_url);

CREATE INDEX idx_chatters_collected ON public.chatters USING btree (collected DESC);

CREATE INDEX idx_chatters_url ON public.chatters USING btree (url);

CREATE INDEX idx_publishers_base_url ON public.publishers USING btree (base_url);

CREATE INDEX idx_publishers_source ON public.publishers USING btree (domain_name);

CREATE UNIQUE INDEX idx_trend_aggregates_id ON public.trend_aggregates USING btree (id);

```

## Documentation dependency map

Swagger annotations, gateway OpenAPI, portal pages are separate artifacts. They need to be updated separately
For any update in public routes, params and responses
1. Update Swagger annotations in `apis/<product>/router/` as the service-local contract. After annotation changes, regenerate and commit the service's Swagger outputs; never hand-edit generated `docs/docs.go`, `docs/swagger.json`, or `docs/swagger.yaml`. Always include request, response and error type definiton for each route.
2. Update api gateway definitions `../config/<product>.oas.json`'
3. Update developer portal docs under `../docs/pages` e.g. ', `../docs/pages/products/<product>/`, and their effect on shared documents like `../docs/pages/start`, `../docs/pages/guides`. Always include sample params and responses
4. Update corresponding Bruno definitions in `<product>/bruno/`
5. Close the Definition of Done in root `AGENTS.md` and `.github/pull_request_template.md`. Maintainer run/test/swag commands: [`README.md`](README.md). Frozen public policy: **100/min**, **50,000/month**, **Bearer** except health, no private backend headers in public docs. Contract/docs GitHub workflow is manual only and does not gate CI.

### Public documentation boundary

The public surfaces are gateway OpenAPI files in `config/`, backend Swagger artifacts, and `docs/pages/`. They describe stable client behavior, not implementation. Remove any internal design detail found there; use user-facing behavior or generic availability wording instead.

Never expose:

- database internals: Espresso's `cupboard`, `sips`, `sources`, and `relations` storage; internal relationship values such as `same_as` and `derived_from`; digest payloads and their structure; Beans' `beansack` storage; table, schema, column, foreign-key, index, view, or migration details;
- retrieval internals: embeddings, vector dimensions/indexes, HNSW, embedder/model settings, caches, query implementation, or relation direction/relationship storage;
- private implementation and operations: Go package/type/handler names, internal identifiers, backend environment variables and headers, gateway rewrites or policies, infrastructure/vendor topology, credentials, or rate-limit/quota implementation.

Events, Signals, Sources, evidence, filters, pagination, response formats, and API-key requirements are public concepts. Describe them without exposing their storage or implementation.

## Renovation Plan
Read the files in [design/](apis/design) for V1 and future renovation plan