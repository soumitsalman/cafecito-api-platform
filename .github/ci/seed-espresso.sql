-- Deterministic Espresso CI fixtures (events, signals, sources, relations).
-- Embeddings are a constant 320-d vector so pgvector filters are exercisable
-- without a live embedder. Router vector-search tests still need HTTP embedder.

INSERT INTO sources (id, base_url, domain_name, site_name, description, favicon, rss_feed)
VALUES (
    'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1',
    'https://ci.example.com',
    'ci.example.com',
    'CI Source',
    'Deterministic CI publisher',
    'https://ci.example.com/favicon.ico',
    'https://ci.example.com/rss'
)
ON CONFLICT (id) DO UPDATE SET site_name = EXCLUDED.site_name;

INSERT INTO sips (id, kind, created, source, embedding, tags, digest, url, base_url)
SELECT
    id,
    kind,
    created,
    'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1'::uuid,
    array_fill(0.01::real, ARRAY[320])::vector,
    tags,
    digest,
    url,
    'https://ci.example.com'
FROM (
    VALUES
        (
            'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb01'::uuid,
            'event',
            NOW() - INTERVAL '1 hour',
            ARRAY['us', 'japan', 'china', 'ai', 'academic']::text[],
            '{"event_type":"policy","regions":["us","japan","china"],"people":["Ada Lovelace"],"companies":["OpenAI"],"products":["Espresso"],"summary":"CI fixture event 1","impact_level":"high"}'::jsonb,
            'https://ci.example.com/events/1'
        ),
        (
            'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb02'::uuid,
            'event',
            NOW() - INTERVAL '2 hours',
            ARRAY['us', 'ai']::text[],
            '{"event_type":"market","regions":["us"],"people":["Grace Hopper"],"companies":["Google"],"products":["Beans"],"summary":"CI fixture event 2","impact_level":"medium"}'::jsonb,
            'https://ci.example.com/events/2'
        ),
        (
            'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb03'::uuid,
            'event',
            NOW() - INTERVAL '3 hours',
            ARRAY['japan', 'china']::text[],
            '{"event_type":"policy","regions":["japan","china"],"people":["Alan Turing"],"companies":["OpenAI"],"products":["Latte"],"summary":"CI fixture event 3","impact_level":"low"}'::jsonb,
            'https://ci.example.com/events/3'
        ),
        (
            'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb04'::uuid,
            'event',
            NOW() - INTERVAL '4 hours',
            ARRAY['us']::text[],
            '{"event_type":"litigation","regions":["us"],"people":["Katherine Johnson"],"companies":["Acme"],"products":["Cortado"],"summary":"CI fixture event 4","impact_level":"high"}'::jsonb,
            'https://ci.example.com/events/4'
        ),
        (
            'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb05'::uuid,
            'event',
            NOW() - INTERVAL '5 hours',
            ARRAY['ai', 'academic']::text[],
            '{"event_type":"research","regions":["us"],"people":["Dorothy Vaughan"],"companies":["Google"],"products":["Espresso"],"summary":"CI fixture event 5","impact_level":"medium"}'::jsonb,
            'https://ci.example.com/events/5'
        ),
        (
            'cccccccc-cccc-4ccc-8ccc-cccccccccc01'::uuid,
            'signal',
            NOW() - INTERVAL '90 minutes',
            ARRAY['us', 'ai', 'academic']::text[],
            '{"summary":"CI fixture signal","outlook":"stable"}'::jsonb,
            'https://ci.example.com/signals/1'
        )
) AS rows(id, kind, created, tags, digest, url)
ON CONFLICT (id) DO UPDATE SET
    tags = EXCLUDED.tags,
    digest = EXCLUDED.digest,
    created = EXCLUDED.created,
    embedding = EXCLUDED.embedding;

INSERT INTO relations (from_id, to_id, relationship)
VALUES
    (
        'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb01',
        'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb02',
        'SAME_AS'
    ),
    (
        'cccccccc-cccc-4ccc-8ccc-cccccccccc01',
        'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbb01',
        'DERIVED_FROM'
    )
ON CONFLICT DO NOTHING;
