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
-- Name: immutable_tags_to_text(text[]); Type: FUNCTION; Schema: public; Owner: -
--

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


--
-- Name: sips sips_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sips
    ADD CONSTRAINT sips_pkey PRIMARY KEY (id);


--
-- Name: sources sources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sources
    ADD CONSTRAINT sources_pkey PRIMARY KEY (id);


--
-- Name: idx_relations_from_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_relations_from_id ON public.relations USING btree (from_id);


--
-- Name: idx_relations_relationship; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_relations_relationship ON public.relations USING btree (relationship);


--
-- Name: idx_relations_to_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_relations_to_id ON public.relations USING btree (to_id);


--
-- Name: idx_sips_base_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_base_url ON public.sips USING btree (base_url);


--
-- Name: idx_sips_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_created ON public.sips USING btree (created);


--
-- Name: idx_sips_embedding_hnsw; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_embedding_hnsw ON public.sips USING hnsw (embedding public.vector_cosine_ops) WITH (m='24', ef_construction='128');


--
-- Name: idx_sips_kind; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_kind ON public.sips USING btree (kind);


--
-- Name: idx_sips_source; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_source ON public.sips USING btree (source);


--
-- Name: idx_sips_tags; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_tags ON public.sips USING btree (tags);


--
-- Name: idx_sips_tags_fts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_tags_fts ON public.sips USING gin (tags_fts);


--
-- Name: idx_sips_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sips_url ON public.sips USING btree (url);


--
-- Name: idx_sources_base_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sources_base_url ON public.sources USING btree (base_url);


