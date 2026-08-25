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


--
-- Name: beans beans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.beans
    ADD CONSTRAINT beans_pkey PRIMARY KEY (id);


--
-- Name: publishers publishers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.publishers
    ADD CONSTRAINT publishers_pkey PRIMARY KEY (id);


--
-- Name: related_beans related_beans_v2_bean_id_related_bean_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.related_beans
    ADD CONSTRAINT related_beans_v2_bean_id_related_bean_id_key UNIQUE (bean_id, related_bean_id);


--
-- Name: idx_beans_categories; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_categories ON public.beans USING gin (categories);


--
-- Name: idx_beans_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_created ON public.beans USING btree (created DESC);


--
-- Name: idx_beans_embedding_hnsw_cosine; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_embedding_hnsw_cosine ON public.beans USING hnsw (embedding public.vector_cosine_ops) WITH (m='24', ef_construction='128');


--
-- Name: idx_beans_entities; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_entities ON public.beans USING gin (entities);


--
-- Name: idx_beans_kind; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_kind ON public.beans USING btree (kind);


--
-- Name: idx_beans_regions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_regions ON public.beans USING gin (regions);


--
-- Name: idx_beans_source_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_source_id ON public.beans USING btree (source_id);


--
-- Name: idx_beans_tags; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_tags ON public.beans USING gin (tags);


--
-- Name: idx_beans_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_beans_url ON public.beans USING btree (url);


--
-- Name: idx_chatters_bean_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chatters_bean_id ON public.chatters USING btree (bean_id);


--
-- Name: idx_chatters_chatter_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chatters_chatter_url ON public.chatters USING btree (chatter_url);


--
-- Name: idx_chatters_collected; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chatters_collected ON public.chatters USING btree (collected DESC);


--
-- Name: idx_chatters_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chatters_url ON public.chatters USING btree (url);


--
-- Name: idx_publishers_base_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_publishers_base_url ON public.publishers USING btree (base_url);


--
-- Name: idx_publishers_source; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_publishers_source ON public.publishers USING btree (domain_name);


--
-- Name: idx_trend_aggregates_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_trend_aggregates_id ON public.trend_aggregates USING btree (id);

