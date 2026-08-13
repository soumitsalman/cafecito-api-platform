package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

// finalizePage trims to limit and encodes a next cursor from the last returned row when more rows exist.
func finalizePage[T any](rows []T, limit int, cursor_of func(item T) *Cursor) Page[T] {
	items := rows
	has_next := false
	if len(rows) > limit {
		items = rows[:limit]
		has_next = true
	}
	var next *Cursor
	if has_next && len(items) > 0 {
		next = cursor_of(items[len(items)-1])
	}
	return Page[T]{Items: items, NextCursor: next}
}

// whereIDs appends the ID predicate.
func whereIDs(where *[]string, params pgx.NamedArgs, ids []uuid.UUID, prefix string) {
	if len(ids) > 0 {
		*where = append(*where, prefix+"id = ANY(@ids)")
		params["ids"] = ids
	}
}

// whereSourceIDs appends the source ID predicate.
func whereSourceIDs(where *[]string, params pgx.NamedArgs, source_ids []uuid.UUID, prefix string) {
	if len(source_ids) > 0 {
		*where = append(*where, prefix+"source = ANY(@source_ids)")
		params["source_ids"] = source_ids
	}
}

// whereKind appends the kind predicate.
func whereKind(where *[]string, params pgx.NamedArgs, kind string, prefix string) {
	if kind != "" {
		*where = append(*where, prefix+"kind = @kind")
		params["kind"] = kind
	}
}

// whereCreated appends the from and to time range predicate.
func whereCreated(where *[]string, params pgx.NamedArgs, from, to time.Time, prefix string) {
	if !from.IsZero() {
		*where = append(*where, prefix+"created >= @from")
		params["from"] = from
	}
	if !to.IsZero() {
		*where = append(*where, prefix+"created <= @to")
		params["to"] = to
	}
}

// whereTags appends the allowlisted any tag predicate.
func whereTags(where *[]string, params pgx.NamedArgs, tags []string, prefix string) {
	if len(tags) > 0 {
		*where = append(*where, prefix+"tags && @tags")
		params["tags"] = tags
	}
}

// appendTagsFtsWhere appends the allowlisted any tag FTS predicate.
func whereTagsFts(where *[]string, params pgx.NamedArgs, tags []string, prefix string) {
	if len(tags) > 0 {
		*where = append(*where, prefix+"tags_fts @@ plainto_tsquery('simple', @tags)")
		params["tags"] = strings.Join(tags, " & ")
	}
}

func whereDigest(where *[]string, params pgx.NamedArgs, filters *Filters, prefix string) {
	if len(filters.Companies) > 0 {
		*where = append(*where, prefix+"digest->'companies' ?| @companies")
		params["companies"] = filters.Companies
	}
	if len(filters.People) > 0 {
		*where = append(*where, prefix+"digest->'people' ?| @people")
		params["people"] = filters.People
	}
	if len(filters.Products) > 0 {
		*where = append(*where, prefix+"digest->'products' ?| @products")
		params["products"] = filters.Products
	}
	if len(filters.Regions) > 0 {
		*where = append(*where, prefix+"digest->'regions' ?| @regions")
		params["regions"] = filters.Regions
	}
	if len(filters.EventTypes) > 0 {
		*where = append(*where, prefix+"digest->'event_type' ?| @event_types")
		params["event_types"] = filters.EventTypes
	}
	if len(filters.ImpactLevels) > 0 {
		*where = append(*where, prefix+"digest->'impact_level' ?| @impact_levels")
		params["impact_levels"] = filters.ImpactLevels
	}
	if len(filters.ImpactedDomains) > 0 {
		*where = append(*where, prefix+"digest->'impacted_domains' ?| @impacted_domains")
		params["impacted_domains"] = filters.ImpactedDomains
	}
}

// buildSipWhere builds the shared structured filter predicates for Event-family and Signal queries.
// kind_expr is the kind predicate; expand_source_same_as enables the SAME_AS evidence expansion for source_ids (Event-family only).
func buildWhere(filters *Filters, alias string) ([]string, pgx.NamedArgs) {
	where := []string{}
	params := pgx.NamedArgs{}
	if alias != "" {
		alias += "."
	}

	whereIDs(&where, params, filters.IDs, alias)
	whereSourceIDs(&where, params, filters.SourceIDs, alias)
	whereKind(&where, params, filters.Kind, alias)
	whereCreated(&where, params, filters.CreatedFrom, filters.CreatedTo, alias)
	whereTagsFts(&where, params, filters.Tags, alias)
	whereDigest(&where, params, filters, alias)

	// if len(filters.SourceIDs) > 0 {
	// 	if expand_source_same_as {
	// 		where = append(where, `(
	// 			s.source = ANY(@source_ids)
	// 			OR EXISTS (
	// 				SELECT 1
	// 				FROM relations AS src_same
	// 				JOIN sips AS evidence
	// 				  ON evidence.id = CASE
	// 				      WHEN src_same.from_id = s.id THEN src_same.to_id
	// 				      ELSE src_same.from_id
	// 				  END
	// 				 AND evidence.kind LIKE 'event%'
	// 				WHERE src_same.relationship = 'SAME_AS'
	// 				  AND (src_same.from_id = s.id OR src_same.to_id = s.id)
	// 				  AND evidence.source = ANY(@source_ids)
	// 			)
	// 		)`)
	// 	} else {
	// 		where = append(where, "s.source = ANY(@source_ids)")
	// 	}
	// 	params["source_ids"] = filters.SourceIDs
	// }
	// if len(filters.EventTypes) > 0 {
	// 	where = append(where, "s.digest->>'event_type' = ANY(@event_types)")
	// 	params["event_types"] = filters.EventTypes
	// }
	// if len(filters.ImpactLevels) > 0 {
	// 	where = append(where, "s.digest->>'impact_level' = ANY(@impact_levels)")
	// 	params["impact_levels"] = filters.ImpactLevels
	// }
	// if len(filters.ImpactedDomains) > 0 {
	// 	where = append(where, "(s.digest->'impacted_domains') ?| @impacted_domains")
	// 	params["impacted_domains"] = filters.ImpactedDomains
	// }
	// for key, vals := range map[string][]string{
	// 	"companies": filters.Companies,
	// 	"people":    filters.People,
	// 	"products":  filters.Products,
	// 	"regions":   filters.Regions,
	// } {
	// 	if len(vals) > 0 {
	// 		where = append(where, fmt.Sprintf("(s.digest->'%s') ?| @%s", key, key))
	// 		params[key] = vals
	// 	}
	// }
	// appendTagWhere(&where, params, filters.Tags, filters.TagMode, "s")
	return where, params
}

func buildScalarQuery(filters *Filters, page *PageRequest) (string, pgx.NamedArgs) {
	expr_fmt := `
	SELECT id, kind, created, tags, digest
	FROM sips
	%s -- WHERE cursor AND filters
	ORDER BY created DESC, id DESC
	LIMIT @limit
	`
	where, params := buildWhere(filters, "")
	if page.Cursor != nil && page.Cursor.Created != nil && page.Cursor.ID != nil {
		where = append(where, "(created, id) < (@cursor_created, @cursor_id)")
		params["cursor_created"] = page.Cursor.Created
		params["cursor_id"] = *page.Cursor.ID
	}
	where_expr := ""
	if len(where) > 0 {
		where_expr = "WHERE " + strings.Join(where, " AND ")
	}

	params["limit"] = page.Limit + 1
	return fmt.Sprintf(expr_fmt, where_expr), params
}

func buildVectorSearchQuery(filters *Filters, page *PageRequest) (string, pgx.NamedArgs) {
	expr_fmt := `
	WITH nearest_results AS MATERIALIZED (
		SELECT 
			id, kind, created, tags, digest, 
			embedding <=> @embedding AS distance
		FROM sips
		%s -- WHERE cursor AND filters
		ORDER BY distance ASC
		LIMIT @limit
	)
	SELECT * FROM nearest_results 
	%s -- WHERE distance <= @distance
	ORDER BY distance ASC, id ASC
	LIMIT @limit`
	where, params := buildWhere(filters, "")
	if page.Cursor != nil && page.Cursor.ID != nil {
		where = append(where, "(embedding <=> @embedding, id) > (@cursor_distance, @cursor_id)")
		params["cursor_distance"] = *page.Cursor.Distance
		params["cursor_id"] = *page.Cursor.ID
	}
	inner_where_expr := ""
	if len(where) > 0 {
		inner_where_expr = "WHERE " + strings.Join(where, " AND ")
	}
	outer_where_expr := ""
	if filters.Distance != nil {
		outer_where_expr = "WHERE distance <= @distance"
		params["distance"] = *filters.Distance
	}

	params["embedding"] = pgvector.NewVector(filters.Embedding)
	params["limit"] = page.Limit + 1
	return fmt.Sprintf(expr_fmt, inner_where_expr, outer_where_expr), params
}

// QuerySips returns Sip records matching the filters.
func (p *Cupboard) QuerySips(ctx context.Context, filters Filters, page PageRequest) (Page[Sip], error) {
	// sip fields to return id, created, kind, tags, digest
	expr, params := "", pgx.NamedArgs{}
	if len(filters.Embedding) > 0 {
		expr, params = buildVectorSearchQuery(&filters, &page)
	} else {
		expr, params = buildScalarQuery(&filters, &page)
	}
	rows, err := fetchAll[Sip](ctx, p.db, expr, params)
	if err != nil {
		return Page[Sip]{}, err
	}
	return finalizePage(rows, page.Limit, func(s Sip) *Cursor {
		if len(filters.Embedding) > 0 {
			return &Cursor{Version: _CURSOR_VERSION, ID: &s.ID, Distance: &s.Distance.Float64}
		} else {
			return &Cursor{Version: _CURSOR_VERSION, ID: &s.ID, Created: &s.Created}
		}
	}), nil
}

// GetSip retrieves one Sip record by UUID and optional kind. Returns (Zero, nil) when not found.
func (p *Cupboard) GetSip(ctx context.Context, id uuid.UUID, kind string) (ExtendedSip, error) {
	// sip fields to return id, kind, created, tags, digest, source, url, base_url
	query := `SELECT 
		s.id, s.kind, s.created, s.tags, s.digest, s.source, s.url, s.base_url, 
		src.domain_name, src.site_name, src.description, src.favicon, src.rss_feed
	FROM sips s
	LEFT JOIN sources src ON s.source = src.id
	WHERE s.id = @id`
	params := pgx.NamedArgs{"id": id}
	if kind != "" {
		query += " AND s.kind = @kind"
		params["kind"] = kind
	}
	sip, err := fetchOne[ExtendedSip](ctx, p.db, query+" LIMIT 1", params)
	if errors.Is(err, pgx.ErrNoRows) {
		return ExtendedSip{}, nil
	}
	return sip, err
}

// SipExists checks if a Sip record exists with the given id and optional kind.
func (p *Cupboard) SipExists(ctx context.Context, id uuid.UUID, kind string) (bool, error) {
	query := "SELECT 1 FROM sips WHERE id = @id"
	params := pgx.NamedArgs{"id": id}
	if kind != "" {
		query += " AND kind = @kind"
		params["kind"] = kind
	}
	query = fmt.Sprintf("SELECT EXISTS(%s LIMIT 1)", query)
	return fetchOneScalar[bool](ctx, p.db, query, params)
}

// QuerySameSips returns sips that are SAME_AS the given id.
func (p *Cupboard) QuerySameSips(ctx context.Context, id uuid.UUID, filters Filters, page PageRequest) (Page[Sip], error) {
	expr_fmt := `
	WITH same_scope AS (
		SELECT @id AS anchor_id
		UNION
		SELECT CASE WHEN from_id = @id THEN to_id ELSE from_id END
		FROM relations
		WHERE relationship = 'SAME_AS'
		  AND (from_id = @id OR to_id = @id)
	)
	SELECT id, created, kind, source, url, base_url
	FROM sips
	INNER JOIN same_scope ON id = anchor_id
	%s -- cursor AND filters
	ORDER BY created DESC, id DESC
	LIMIT @limit`

	where, params := buildWhere(&filters, "")
	if page.Cursor != nil {
		where = append(where, "((created, id) < (@cursor_created, @cursor_id))")
		params["cursor_created"] = page.Cursor.Created
		params["cursor_id"] = page.Cursor.ID
	}
	where_expr := ""
	if len(where) > 0 {
		where_expr = "WHERE " + strings.Join(where, " AND ")
	}

	params["id"] = id
	params["limit"] = page.Limit + 1
	rows, err := fetchAll[Sip](ctx, p.db, fmt.Sprintf(expr_fmt, where_expr), params)
	if err != nil {
		return Page[Sip]{}, err
	}
	return finalizePage(rows, page.Limit, func(s Sip) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, ID: &s.ID, Created: &s.Created}
	}), nil
}

// QueryDerivedSips returns Sips derived from the Sip.
func (p *Cupboard) QueryDerivedSips(ctx context.Context, id uuid.UUID, filters Filters, page PageRequest) (Page[Sip], error) {
	anchor_expr := `
	SELECT from_id AS anchor_id 
	FROM relations 
	WHERE relationship = 'DERIVED_FROM' AND to_id = @anchor_id
	`
	return p.queryDerivedRelations(ctx, anchor_expr, &id, &filters, &page)
}

// QueryDerivedSips returns Sips which were used for deriving the Sip.
func (p *Cupboard) QuerySupportingSips(ctx context.Context, id uuid.UUID, filters Filters, page PageRequest) (Page[Sip], error) {
	anchor_expr := `
	SELECT to_id AS anchor_id 
	FROM relations 
	WHERE relationship = 'DERIVED_FROM' AND from_id = @anchor_id
	`
	return p.queryDerivedRelations(ctx, anchor_expr, &id, &filters, &page)
}

func (p *Cupboard) queryDerivedRelations(ctx context.Context, anchor_expr string, anchor_id *uuid.UUID, filters *Filters, page *PageRequest) (Page[Sip], error) {
	expr_fmt := `
	WITH anchor_scope AS (%s)
	SELECT id, created, kind, tags, digest 
	FROM sips
	INNER JOIN anchor_scope ON id = anchor_id
	%s -- cursor AND filters
	ORDER BY created DESC, id DESC
	LIMIT @limit`

	where, params := buildWhere(filters, "")
	if page.Cursor != nil {
		where = append(where, "(created, id) < (@cursor_created, @cursor_id)")
		params["cursor_created"] = page.Cursor.Created
		params["cursor_id"] = page.Cursor.ID
	}
	where_expr := ""
	if len(where) > 0 {
		where_expr = "WHERE " + strings.Join(where, " AND ")
	}

	params["anchor_id"] = anchor_id
	params["limit"] = page.Limit + 1
	rows, err := fetchAll[Sip](ctx, p.db, fmt.Sprintf(expr_fmt, anchor_expr, where_expr), params)
	if err != nil {
		return Page[Sip]{}, err
	}
	return finalizePage(rows, page.Limit, func(s Sip) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, ID: &s.ID, Created: &s.Created}
	}), nil
}

// GetEventRelationCounts returns the explicit relationship metadata used by R02.
func (p *Cupboard) CountRelations(ctx context.Context, id uuid.UUID) (RelationCounts, error) {
	query := `
	SELECT
    COUNT(*) FILTER (
        WHERE relationship = 'SAME_AS'
          AND (from_id = @anchor_id OR to_id = @anchor_id)
    ) AS same_as_count,

    COUNT(*) FILTER (
        WHERE relationship = 'DERIVED_FROM'
          AND from_id = @anchor_id
    ) AS derived_from_count,

    COUNT(*) FILTER (
        WHERE relationship = 'DERIVED_FROM'
          AND to_id = @anchor_id
    ) AS derived_to_count
	FROM relations;
	`
	return fetchOne[RelationCounts](ctx, p.db, query, pgx.NamedArgs{"anchor_id": id})
}

// QuerySources returns source records matching the optional query and domain filters.
func (p *Cupboard) QuerySources(ctx context.Context, q string, domains []string, page PageRequest) (Page[Source], error) {
	expr_fmt := `
	SELECT id, base_url, domain_name, site_name, description, favicon, rss_feed
	FROM sources
	WHERE base_url > @c_base_url -- base_url is 1:1 with id
		%s -- q
		%s -- domains
	ORDER BY base_url ASC 
	LIMIT @limit
	`
	q_expr, domains_expr := "", ""
	c_base_url := ""
	if page.Cursor != nil && page.Cursor.TextKey != nil {
		c_base_url = *page.Cursor.TextKey
	}
	params := pgx.NamedArgs{
		"c_base_url": c_base_url,
		"limit":      page.Limit + 1,
	}
	if len(q) > 0 {
		q_expr = "AND (domain_name ILIKE '%' || @q || '%' OR site_name ILIKE '%' || @q || '%' OR base_url ILIKE '%' || @q || '%')"
		params["q"] = strings.TrimSpace(q)
	}
	if len(domains) > 0 {
		domains_expr = "AND domain_name = ANY(@domains)"
		params["domains"] = domains
	}

	rows, err := fetchAll[Source](ctx, p.db, fmt.Sprintf(expr_fmt, q_expr, domains_expr), params)
	if err != nil {
		return Page[Source]{}, err
	}
	return finalizePage(rows, page.Limit, func(source Source) *Cursor {
		key := source.BaseURL
		return &Cursor{Version: _CURSOR_VERSION, TextKey: &key}
	}), nil
}

// GetSource retrieves one source record by UUID. Returns (zero, nil) when not found.
func (p *Cupboard) GetSource(ctx context.Context, id uuid.UUID) (Source, error) {
	query := "SELECT id, base_url, domain_name, site_name, description, favicon, rss_feed FROM sources WHERE id = @id LIMIT 1"
	source, err := fetchOne[Source](ctx, p.db, query, pgx.NamedArgs{"id": id})
	if errors.Is(err, pgx.ErrNoRows) {
		return Source{}, nil
	}
	return source, err
}

// QueryTags returns the distinct tag strings, optionally filtered by substring and kind scope.
func (p *Cupboard) QueryTags(ctx context.Context, q string, kinds []string, page PageRequest) (Page[string], error) {
	expr_fmt := `
	WITH values AS (
		SELECT DISTINCT unnest(tags) AS tag FROM sips
		WHERE tags IS NOT NULL
		%s -- kinds
		%s -- tags_fts
	)
	SELECT tag FROM values
	WHERE tag > @c_tag -- cursor
	%s -- tag_ilike
	ORDER BY tag ASC
	LIMIT @limit
	`
	kinds_expr, tags_fts_expr, tag_ilike_expr := "", "", ""
	c_tag := ""
	if page.Cursor != nil && page.Cursor.TextKey != nil {
		c_tag = *page.Cursor.TextKey
	}
	params := pgx.NamedArgs{
		"c_tag": c_tag,
		"limit": page.Limit + 1,
	}

	if len(q) > 0 {
		tags_fts_expr = "AND tags_fts @@ plainto_tsquery('simple', @q)"
		tag_ilike_expr = "AND tag ILIKE '%' || @q || '%'"
		params["q"] = strings.TrimSpace(q)
	}
	if len(kinds) > 0 {
		kinds_expr = "AND kind = ANY(@kinds)"
		params["kinds"] = kinds
	}

	rows, err := fetchAllScalar[string](ctx, p.db, fmt.Sprintf(expr_fmt, kinds_expr, tags_fts_expr, tag_ilike_expr), params)
	if err != nil {
		return Page[string]{}, err
	}
	return finalizePage(rows, page.Limit, func(item string) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, TextKey: &item}
	}), nil
}

// eventTagQueries defines the SQL source for each supported Event discovery type.
// The query fragments are static and selected only from the allowlisted constants.
var eventTagQueries = map[string]string{
	EVENT_TAG_TYPE_REGION: `
		SELECT jsonb_array_elements_text(COALESCE(digest->'regions', '[]'::jsonb)) AS value,
		       'region' AS type
		FROM sips
		WHERE kind = 'event'
	`,
	EVENT_TAG_TYPE_PEOPLE: `
		SELECT jsonb_array_elements_text(COALESCE(digest->'people', '[]'::jsonb)) AS value,
		       'people' AS type
		FROM sips
		WHERE kind = 'event'
	`,
	EVENT_TAG_TYPE_PRODUCT: `
		SELECT jsonb_array_elements_text(COALESCE(digest->'products', '[]'::jsonb)) AS value,
		       'product' AS type
		FROM sips
		WHERE kind = 'event'
	`,
	EVENT_TAG_TYPE_COMPANY: `
		SELECT jsonb_array_elements_text(COALESCE(digest->'companies', '[]'::jsonb)) AS value,
		       'company' AS type
		FROM sips
		WHERE kind = 'event'
	`,
	EVENT_TAG_TYPE_STOCK_TICKER: `
		SELECT jsonb_array_elements_text(COALESCE(digest->'stock_tickers', '[]'::jsonb)) AS value,
		       'stock_ticker' AS type
		FROM sips
		WHERE kind = 'event'
	`,
	EVENT_TAG_TYPE_EVENT_TYPE: `
		SELECT digest->>'event_type' AS value,
		       'event_type' AS type
		FROM sips
		WHERE kind = 'event'
		  AND digest->>'event_type' IS NOT NULL
		  AND digest->>'event_type' <> ''
	`,
}

var allEventTagTypes = []string{
	EVENT_TAG_TYPE_REGION,
	EVENT_TAG_TYPE_PEOPLE,
	EVENT_TAG_TYPE_PRODUCT,
	EVENT_TAG_TYPE_COMPANY,
	EVENT_TAG_TYPE_STOCK_TICKER,
	EVENT_TAG_TYPE_EVENT_TYPE,
}

// QueryEventTags returns one common discovery payload for any requested Event
// digest types. Results are ordered by value first and type second.
func (p *Cupboard) QueryEventTags(ctx context.Context, q string, types []string, page PageRequest) (Page[TagValue], error) {
	if len(types) == 0 {
		types = allEventTagTypes
	}

	values_queries := make([]string, 0, len(types))
	seen_types := make(map[string]struct{}, len(types))
	for _, tag_type := range types {
		if _, seen := seen_types[tag_type]; seen {
			continue
		}

		values_query, ok := eventTagQueries[tag_type]
		if !ok {
			return Page[TagValue]{}, fmt.Errorf("unsupported event tag type: %s", tag_type)
		}

		seen_types[tag_type] = struct{}{}
		values_queries = append(values_queries, values_query)
	}

	params := pgx.NamedArgs{
		"cursor_value": "",
		"cursor_type":  "",
		"limit":        page.Limit + 1,
	}
	where := []string{"(value, type) > (@cursor_value, @cursor_type)"}
	if page.Cursor != nil && page.Cursor.EventTag != nil {
		params["cursor_value"] = page.Cursor.EventTag.Value
		params["cursor_type"] = page.Cursor.EventTag.Type
	}
	if q = strings.TrimSpace(q); q != "" {
		where = append(where, "value ILIKE '%' || @q || '%'")
		params["q"] = q
	}

	query := fmt.Sprintf(
		"WITH values AS (%s) SELECT DISTINCT value, type FROM values WHERE %s ORDER BY value ASC, type ASC LIMIT @limit",
		strings.Join(values_queries, " UNION ALL "),
		strings.Join(where, " AND "),
	)
	rows, err := fetchAll[TagValue](ctx, p.db, query, params)
	if err != nil {
		return Page[TagValue]{}, err
	}

	return finalizePage(rows, page.Limit, func(item TagValue) *Cursor {
		return &Cursor{
			Version:  _CURSOR_VERSION,
			EventTag: &item,
		}
	}), nil
}
