package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
	datautils "github.com/soumitsalman/data-utils"
)

const (
	MIN_CANDIDATE_LIMIT = 2048
)

const (
	_BEAN_COLUMNS_BASE         = "id, url, kind, created, author, image_url, language, categories, sentiments, entities, regions, ideology, title, source_id, base_url, domain_name, site_name, cluster_id"
	_BEAN_COLUMNS_SUMMARY      = "summary"
	_BEAN_COLUMNS_CONTENT      = "CASE WHEN restricted_content THEN NULL ELSE content END AS content"
	_BEAN_COLUMNS_TREND        = "likes, comments, mentions, subscribers, related, trend_score"
	_BEAN_COLUMNS_ALL          = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY + ", " + _BEAN_COLUMNS_CONTENT + ", " + _BEAN_COLUMNS_TREND
	BEAN_COLUMNS_HEADLINES     = _BEAN_COLUMNS_BASE
	BEAN_COLUMNS_WITHOUT_TREND = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY
	BEAN_COLUMNS_WITH_TREND    = _BEAN_COLUMNS_BASE + ", " + _BEAN_COLUMNS_SUMMARY + ", " + _BEAN_COLUMNS_TREND
	BEAN_COLUMNS_PREVIEW       = "id, url, created, title, source_id, base_url, domain_name, site_name, favicon"
	// BEAN_COLUMNS_MINIMAL_WITH_TREND = BEAN_COLUMNS_PREVIEW + ", " + _BEAN_COLUMNS_TREND
)

const (
	SOURCE_COLUMNS_BASE = "id, base_url, domain_name, site_name"
	SOURCE_COLUMNS_ALL  = SOURCE_COLUMNS_BASE + ", description, favicon, rss_feed"
)

const (
	SORT_RECENT   = "created"
	SORT_TRENDING = "trend_score"
	SORT_RELEVANT = "distance"
)

var (
	ErrNonExistentID = errors.New("Item with this ID does not exist")
	DBError          = errors.New("Database or query error")
)

// finalizePage trims to limit and encodes a next cursor from the last returned row when more rows exist.
// A non-positive limit returns every row and no next cursor (unbounded callers).
func finalizePage[T any](rows []T, limit int, cursor_of func(item T) *Cursor) Page[T] {
	if limit <= 0 {
		return Page[T]{Items: rows}
	}
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

func buildSelect(columns string, full_content bool) string {
	if columns == "" {
		columns = _BEAN_COLUMNS_ALL
	}
	if !strings.Contains(columns, "content") && full_content {
		columns += ", " + _BEAN_COLUMNS_CONTENT
	}
	return columns
}

// buildScalarWhere constructs the shared WHERE clause predicates for article queries.
// All filters reference columns on the beans table (or views that inherit beans.*).
func buildScalarWhere(filters *BeanFilters) ([]string, pgx.NamedArgs) {
	where := []string{}
	params := pgx.NamedArgs{}

	if len(filters.IDs) > 0 {
		where = append(where, "id = ANY(@ids)")
		params["ids"] = filters.IDs
	}
	if len(filters.ExcludeIDs) > 0 {
		where = append(where, "id != ALL(@exclude_ids)")
		params["exclude_ids"] = filters.ExcludeIDs
	}
	if len(filters.URLs) > 0 {
		where = append(where, "url = ANY(@urls)")
		params["urls"] = filters.URLs
	}
	if len(filters.Sources) > 0 {
		where = append(where, "source_id = ANY(@source_ids)")
		params["source_ids"] = filters.Sources
	}
	if len(filters.ExcludeSources) > 0 {
		where = append(where, "source_id != ALL(@exclude_source_ids)")
		params["exclude_source_ids"] = filters.ExcludeSources
	}
	if len(filters.Domains) > 0 {
		where = append(where, "domain_name = ANY(@domains)")
		params["domains"] = filters.Domains
	}
	if len(filters.ExcludeDomains) > 0 {
		where = append(where, "domain_name != ALL(@exclude_domains)")
		params["exclude_domains"] = filters.ExcludeDomains
	}
	if filters.Kind != "" {
		where = append(where, "kind = @kind")
		params["kind"] = filters.Kind
	}
	if !filters.CreatedFrom.IsZero() {
		where = append(where, "created >= @created_from")
		params["created_from"] = filters.CreatedFrom
	}
	if !filters.CreatedTo.IsZero() {
		where = append(where, "created <= @created_to")
		params["created_to"] = filters.CreatedTo
	}
	if !filters.ObservedFrom.IsZero() {
		where = append(where, "observed >= @observed_from")
		params["observed_from"] = filters.ObservedFrom
	}
	if !filters.ObservedTo.IsZero() {
		where = append(where, "observed <= @observed_to")
		params["observed_to"] = filters.ObservedTo
	}
	if len(filters.Tags) > 0 {
		where = append(where, "tags @@ plainto_tsquery('simple', @tags_query)")
		params["tags_query"] = strings.Join(filters.Tags, " & ")
	}
	if len(filters.Authors) > 0 {
		for i, author := range filters.Authors {
			param_key := fmt.Sprintf("author_%d", i)
			where = append(where, "author ILIKE '%' || @"+param_key+" || '%'")
			params[param_key] = strings.TrimSpace(author)
		}
	}
	if len(filters.Categories) > 0 {
		where = append(where, "categories && @categories")
		params["categories"] = filters.Categories
	}
	if len(filters.ExcludeCategories) > 0 {
		where = append(where, "NOT (categories && @exclude_categories)")
		params["exclude_categories"] = filters.ExcludeCategories
	}
	if len(filters.Regions) > 0 {
		where = append(where, "regions && @regions")
		params["regions"] = filters.Regions
	}
	if len(filters.Entities) > 0 {
		where = append(where, "entities && @entities")
		params["entities"] = filters.Entities
	}
	if len(filters.Sentiments) > 0 {
		where = append(where, "sentiments && @sentiments")
		params["sentiments"] = filters.Sentiments
	}
	if len(filters.Languages) > 0 {
		where = append(where, "language = ANY(@languages)")
		params["languages"] = filters.Languages
	}
	if filters.ClusterID != uuid.Nil {
		where = append(where, "cluster_id = @cluster_id")
		params["cluster_id"] = filters.ClusterID
	}
	if len(filters.Extra) > 0 {
		where = append(where, filters.Extra...)
	}
	return where, params
}

// buildScalarOrderQuery constructs the query for latest and trending beans
// table: latest_beans_view or trending_beans_view or aggregated_beans_view.
// sort: Present in page.Cursor.Sort. page must have a cursor and a valid sort value or else the query will fail
// if filters.Embedding is present, then filter.Distance must be present. Default/not-assigned value = 0.0, returns exact match
func buildScalarOrderQuery(table string, filters *BeanFilters, page *PageRequest, select_columns string) (string, pgx.NamedArgs) {
	where, params := buildScalarWhere(filters)
	if len(filters.Embedding) > 0 {
		where = append(where, "embedding <=> @embedding <= @distance")
		params["embedding"] = pgvector.NewVector(filters.Embedding)
		params["distance"] = filters.Distance
	}

	order_by := ""
	if page.Cursor != nil {
		switch page.Cursor.Sort {
		case SORT_RECENT:
			order_by = "created DESC, id DESC"
			if page.Cursor.ID != nil {
				where = append(where, "(created, id) < (@cursor_value, @cursor_id)")
				params["cursor_value"] = *page.Cursor.Created
				params["cursor_id"] = *page.Cursor.ID
			}
		case SORT_TRENDING:
			order_by = "trend_score DESC, id DESC"
			if page.Cursor.ID != nil {
				where = append(where, "(trend_score, id) < (@cursor_value, @cursor_id)")
				params["cursor_value"] = *page.Cursor.TrendScore
				params["cursor_id"] = *page.Cursor.ID
			}
		}
	}
	where_expr := ""
	if len(where) > 0 {
		where_expr = "WHERE " + strings.Join(where, " AND ")
	}
	limit_expr := ""
	if page.Limit > 0 {
		params["limit"] = page.Limit + 1
		limit_expr = "LIMIT @limit"
	}

	query := fmt.Sprintf(`
		SELECT %s -- columns
		FROM %s -- latest_beans_view or trending_beans_view or aggregated_beans_view
		%s -- WHERE
		ORDER BY %s -- order by created DESC, id DESC or trend_score DESC, id DESC
		%s`,
		buildSelect(select_columns, filters.FullContent),
		table,
		where_expr,
		order_by,
		limit_expr,
	)
	return query, params
}

func searchCandidateLimit(limit int) int {
	return max(limit*4, MIN_CANDIDATE_LIMIT)
}

// buildKNNSearchQuery constructs the vector (semantic) search query for both latest and trending beans
// table: latest_beans_view or trending_beans_view or aggregated_beans_view with cosine distance CTE.
// filters.Embedding, filters.Limit must be present. filters.Distance will be considered only if > 0.0
func buildKNNSearchQuery(table string, filters *BeanFilters, page *PageRequest, columns string) (string, pgx.NamedArgs) {
	where, params := buildScalarWhere(filters)
	if page.Cursor != nil && page.Cursor.Distance != nil && page.Cursor.ID != nil {
		where = append(where, "(embedding <=> @embedding, id) > (@cursor_distance, @cursor_id)")
		params["cursor_distance"] = *page.Cursor.Distance
		params["cursor_id"] = *page.Cursor.ID
	}
	inner_where_expr := ""
	if len(where) > 0 {
		inner_where_expr = "WHERE " + strings.Join(where, " AND ")
	}
	outer_where_expr := ""
	if filters.Distance > 0 {
		outer_where_expr = "WHERE distance <= @distance"
		params["distance"] = filters.Distance
	}
	params["embedding"] = pgvector.NewVector(filters.Embedding)
	params["candidate_limit"] = searchCandidateLimit(page.Limit)
	params["limit"] = page.Limit + 1

	query := fmt.Sprintf(`
		WITH nearest_results AS MATERIALIZED (
			SELECT *, embedding <=> @embedding AS distance
			FROM %s	-- latest_beans_view or trending_beans_view or aggregated_beans_view
			%s -- WHERE inner
			ORDER BY distance ASC
			LIMIT @candidate_limit
		)
		SELECT %s, distance FROM nearest_results
		%s
		ORDER BY distance ASC, id ASC
		LIMIT @limit`,
		table,
		inner_where_expr,
		buildSelect(columns, filters.FullContent),
		outer_where_expr,
	)
	return query, params
}

// TrendingArticles returns articles ranked by trend_score descending.
// Forces sort by created descending
func (b *PGSack) QueryLatestBeans(ctx context.Context, filters BeanFilters, page PageRequest, columns string) (Page[Bean], error) {
	if page.Cursor == nil {
		page.Cursor = &Cursor{Sort: SORT_RECENT}
	}
	page.Cursor.Sort = SORT_RECENT

	query, params := buildScalarOrderQuery("latest_beans_view", &filters, &page, columns)
	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, DBError
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, Sort: SORT_RECENT, ID: &bean.ID, Created: &bean.Created}
	}), nil
}

// TrendingArticles returns articles ranked by trend_score descending.
// Forces sort by trend_score descending
func (b *PGSack) QueryTrendingBeans(ctx context.Context, filters BeanFilters, page PageRequest, columns string) (Page[Bean], error) {
	if page.Cursor == nil {
		page.Cursor = &Cursor{Sort: SORT_TRENDING}
	}
	page.Cursor.Sort = SORT_TRENDING
	query, params := buildScalarOrderQuery("trending_beans_view", &filters, &page, columns)
	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, DBError
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, Sort: SORT_TRENDING, ID: &bean.ID, TrendScore: &bean.TrendScore.Float64}
	}), nil
}

// QueryBeans returns beans matching filters, ordered by created descending.
// Dispatches to vector (cosine distance CTE) or scalar query based on whether an embedding is present in the filters.
// Forces sort by created descending if no embedding is present.
// Forces sort by relevance if an embedding is present and filters.Distance > 0.0.
func (b *PGSack) QueryBeans(ctx context.Context, filters BeanFilters, page PageRequest, columns string) (Page[Bean], error) {
	if page.Cursor == nil {
		page.Cursor = &Cursor{}
	}
	var query string
	var params pgx.NamedArgs
	if len(filters.Embedding) > 0 {
		page.Cursor.Sort = SORT_RELEVANT
		query, params = buildKNNSearchQuery("latest_beans_view", &filters, &page, columns)
	} else {
		page.Cursor.Sort = SORT_RECENT
		query, params = buildScalarOrderQuery("latest_beans_view", &filters, &page, columns)
	}

	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, DBError
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		if len(filters.Embedding) > 0 {
			return &Cursor{Version: _CURSOR_VERSION, ID: &bean.ID, Distance: &bean.Distance.Float64}
		} else {
			return &Cursor{Version: _CURSOR_VERSION, ID: &bean.ID, Created: &bean.Created}
		}
	}), nil
}

func (b *PGSack) beanExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return utils.FetchOneScalar[bool](ctx, b.db,
		"SELECT EXISTS (SELECT 1 FROM beans WHERE id = @id LIMIT 1)",
		pgx.NamedArgs{"id": id},
	)
}

// QuerySimilarBeans returns known related publisher coverage for one Article UUID.
// related_beans edges are treated as undirected. Missing IDs return ErrNonExistentID.
func (b *PGSack) QuerySimilarBeans(ctx context.Context, id uuid.UUID, filters BeanFilters, page PageRequest, columns string) (Page[Bean], error) {
	where, params := buildScalarWhere(&filters)
	if page.Cursor != nil && page.Cursor.Created != nil && page.Cursor.ID != nil {
		where = append(where, "(created, id) < (@cursor_created, @cursor_id)")
		params["cursor_created"] = *page.Cursor.Created
		params["cursor_id"] = *page.Cursor.ID
	}
	params["id"] = id
	params["limit"] = page.Limit + 1

	where_expr := ""
	if len(where) > 0 {
		where_expr = " AND " + strings.Join(where, " AND ")
	}
	query := fmt.Sprintf(`
		WITH target_ids AS (
			SELECT related_bean_id AS target_id FROM related_beans WHERE bean_id = @id
			UNION
			SELECT bean_id AS target_id FROM related_beans WHERE related_bean_id = @id
		)
		SELECT %s
		FROM target_ids
		INNER JOIN trending_beans_view ON id = target_id
		WHERE id <> @id
			AND EXISTS (SELECT 1 FROM beans WHERE id = @id)
			%s -- additional where
		ORDER BY created DESC, id DESC
		LIMIT @limit`,
		buildSelect(columns, filters.FullContent),
		where_expr,
	)
	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, DBError
	}
	if len(rows) == 0 {
		exists, err := b.beanExists(ctx, id)
		if err != nil {
			return Page[Bean]{}, DBError
		}
		if !exists {
			return Page[Bean]{}, ErrNonExistentID
		}
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, ID: &bean.ID, Created: &bean.Created}
	}), nil
}

// QueryMentions returns the latest observed social/forum posts linking an Article URL.
// Missing IDs return ErrNonExistentID. Empty membership is an empty page.
func (b *PGSack) QueryMentions(ctx context.Context, id uuid.UUID, filters MentionFilters, page PageRequest) (Page[Mention], error) {
	inner_where := []string{
		"ch.bean_id = @bean_id",
		"EXISTS (SELECT 1 FROM beans WHERE id = @bean_id)",
	}
	params := pgx.NamedArgs{
		"bean_id": id,
		"limit":   page.Limit + 1,
	}
	if len(filters.Platforms) > 0 {
		inner_where = append(inner_where, "LOWER(ch.platform) = ANY(@platforms)")
		params["platforms"] = filters.Platforms
	}
	if len(filters.Forums) > 0 {
		inner_where = append(inner_where, "LOWER(ch.forum) = ANY(@forums)")
		params["forums"] = filters.Forums
	}
	if !filters.ObservedFrom.IsZero() {
		inner_where = append(inner_where, "ch.collected >= @observed_from")
		params["observed_from"] = filters.ObservedFrom
	}
	if !filters.ObservedTo.IsZero() {
		inner_where = append(inner_where, "ch.collected <= @observed_to")
		params["observed_to"] = filters.ObservedTo
	}

	outer_where_expr := ""
	if page.Cursor != nil && page.Cursor.Created != nil && page.Cursor.TextKey != nil {
		outer_where_expr = "WHERE (collected, chatter_url) < (@cursor_collected, @cursor_url)"
		params["cursor_collected"] = *page.Cursor.Created
		params["cursor_url"] = *page.Cursor.TextKey
	}

	query := fmt.Sprintf(`
		WITH latest_chatters AS (
			SELECT DISTINCT ON (ch.chatter_url)
				ch.chatter_url, ch.platform, ch.forum, ch.collected,
				ch.likes, ch.comments, ch.subscribers
			FROM chatters ch
			WHERE %s
			ORDER BY ch.chatter_url, ch.collected DESC
		)
		SELECT 
			chatter_url, platform, forum, 
			collected as observed,
			likes, comments, subscribers
		FROM latest_chatters
		%s
		ORDER BY observed DESC, chatter_url DESC
		LIMIT @limit`,
		strings.Join(inner_where, " AND "),
		outer_where_expr,
	)
	rows, err := utils.FetchAll[Mention](ctx, b.db, query, params)
	if err != nil {
		return Page[Mention]{}, DBError
	}
	if len(rows) == 0 {
		exists, err := b.beanExists(ctx, id)
		if err != nil {
			return Page[Mention]{}, DBError
		}
		if !exists {
			return Page[Mention]{}, ErrNonExistentID
		}
	}
	return finalizePage(rows, page.Limit, func(mention Mention) *Cursor {
		url := mention.URL
		observed_at := mention.Observed
		return &Cursor{Version: _CURSOR_VERSION, Created: &observed_at, TextKey: &url}
	}), nil
}

// GetBean retrieves one bean record by UUID. Returns (zero, ErrInvalidID) when not found.
func (b *PGSack) GetBean(ctx context.Context, id uuid.UUID, full_content bool) (Bean, error) {
	if id == uuid.Nil {
		return Bean{}, ErrNonExistentID
	}
	columns := BEAN_COLUMNS_WITH_TREND
	if full_content {
		columns = _BEAN_COLUMNS_ALL
	}
	query := fmt.Sprintf(`
		SELECT %s -- columns
		FROM latest_beans_view
		WHERE id = @id
		LIMIT 1`,
		columns,
	)
	bean, err := utils.FetchOne[Bean](ctx, b.db, query, pgx.NamedArgs{"id": id})
	if errors.Is(err, pgx.ErrNoRows) {
		return bean, ErrNonExistentID
	}
	if err != nil {
		return bean, DBError
	}
	return bean, err
}

// QuerySources returns source records matching the optional query and domain filters.
func (b *PGSack) QuerySources(ctx context.Context, filters SourceFilters, page PageRequest, columns string) (Page[Source], error) {
	where := []string{}
	params := pgx.NamedArgs{
		"limit": page.Limit + 1,
	}
	if filters.Q != "" {
		where = append(where, "STARTS_WITH(LOWER(site_name), @q) OR STARTS_WITH(LOWER(base_url), @q)")
		params["q"] = filters.Q
	}
	if len(filters.IDs) > 0 {
		where = append(where, "id = ANY(@ids)")
		params["ids"] = filters.IDs
	}
	if len(filters.Domains) > 0 {
		where = append(where, "domain_name = ANY(@domains)")
		params["domains"] = filters.Domains
	}
	if page.Cursor != nil && page.Cursor.TextKey != nil {
		where = append(where, "base_url > @cursor_value")
		params["cursor_value"] = *page.Cursor.TextKey
	}
	where_expr := ""
	if len(where) > 0 {
		where_expr = "WHERE " + strings.Join(where, " AND ")
	}
	if columns == "" {
		columns = SOURCE_COLUMNS_ALL
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM publishers
		%s
		ORDER BY base_url ASC
		LIMIT @limit`,
		columns,
		where_expr,
	)
	rows, err := utils.FetchAll[Source](ctx, b.db, query, params)
	if err != nil {
		return Page[Source]{}, DBError
	}
	return finalizePage(rows, page.Limit, func(source Source) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, TextKey: &source.BaseURL}
	}), nil
}

// GetSource retrieves one source record by UUID. Returns (zero, ErrNonExistentID) when not found.
func (b *PGSack) GetSource(ctx context.Context, id uuid.UUID) (Source, error) {
	query := "SELECT " + SOURCE_COLUMNS_ALL + " FROM publishers WHERE id = @id LIMIT 1"
	source, err := utils.FetchOne[Source](ctx, b.db, query, pgx.NamedArgs{"id": id})
	if errors.Is(err, pgx.ErrNoRows) {
		return Source{}, ErrNonExistentID
	}
	return source, err
}

// QueryTags returns distinct tag strings matching the optional query, scoped to a tag type column.
func (b *PGSack) QueryTags(ctx context.Context, q string, tag_type string, page PageRequest) (Page[string], error) {
	params := pgx.NamedArgs{
		"q":            q,
		"cursor_value": "",
		"limit":        page.Limit + 1,
	}
	if page.Cursor != nil && page.Cursor.TextKey != nil {
		params["cursor_value"] = *page.Cursor.TextKey
	}
	query := fmt.Sprintf(`
		WITH tag_values AS MATERIALIZED (
			SELECT DISTINCT UNNEST(%s) AS value -- tag_type
			FROM beans
			WHERE %s IS NOT NULL
		)
		SELECT value FROM tag_values
		WHERE value > @cursor_value
			AND STARTS_WITH(value, @q)
		ORDER BY value ASC
		LIMIT @limit`,
		tag_type,
		tag_type,
	)
	rows, err := utils.FetchAllScalar[string](ctx, b.db, query, params)
	if err != nil {
		return Page[string]{}, DBError
	}
	return finalizePage(rows, page.Limit, func(item string) *Cursor {
		return &Cursor{
			Version: _CURSOR_VERSION,
			TextKey: &item,
		}
	}), nil
}

func (b *PGSack) QueryClusters(ctx context.Context, filters ClusterFilters, page PageRequest) (Page[Cluster], error) {
	paged, err := b.QueryClusterPreviews(ctx, filters, page)
	if err != nil {
		return Page[Cluster]{}, err
	}
	hydrated, err := b.hydrateClusters(ctx, paged.Items, &filters.BeanFilters)
	if err != nil {
		return Page[Cluster]{}, DBError
	}
	paged.Items = hydrated
	return paged, nil
}

const _CLUSTER_QUERY_TEMPLATE = `
WITH
	candidate_beans AS (%s),
	candidate_clusters AS (%s),
	stats AS (
		SELECT
			b.cluster_id,
			MIN(b.created) AS first_created,
			MAX(b.created) AS last_created,
			COUNT(*) AS beans_count,
			COUNT(DISTINCT b.source_id) AS sources_count
		FROM trending_beans_view AS b
		WHERE EXISTS (
			SELECT 1
			FROM candidate_clusters AS c
			WHERE c.cluster_id = b.cluster_id
		)
		GROUP BY b.cluster_id
		HAVING COUNT(*) >= @min_bean_count
	)
	SELECT 
		c.cluster_id AS id, c.title, c.summary, c.image_url, %s, -- column to order by
		c.categories, c.entities, c.regions,
		s.first_created, s.last_created, s.beans_count, s.sources_count
	FROM candidate_clusters c
	INNER JOIN stats s ON s.cluster_id = c.cluster_id
	%s -- cursor
	%s -- order by
	LIMIT @limit
`

// scalar query for clusters - always sort by latest
func buildScalarQueryForClusters(filters *ClusterFilters, page *PageRequest) (string, pgx.NamedArgs) {
	candidate_beans_template := `
	SELECT * 
	FROM trending_beans_view 
	WHERE cluster_id IS NOT NULL
	%s` // %s is the placeholder for scalar filters
	candidate_clusters_template := `
	SELECT DISTINCT ON (cluster_id) * 
	FROM candidate_beans 
	ORDER BY cluster_id, created DESC`
	filter_template := ""

	where, params := buildScalarWhere(&filters.BeanFilters)
	if len(where) > 0 {
		filter_template = fmt.Sprintf("AND %s", strings.Join(where, " AND "))
	}
	candidate_beans_template = fmt.Sprintf(candidate_beans_template, filter_template)
	params["min_bean_count"] = filters.MinBeanCount
	params["limit"] = page.Limit + 1

	cursor := ""
	if page.Cursor != nil {
		cursor = "WHERE (c.created, c.cluster_id) < (@cursor_value, @cursor_id)"
		params["cursor_value"] = *page.Cursor.Created
		params["cursor_id"] = *page.Cursor.ID
	}

	return fmt.Sprintf(
		_CLUSTER_QUERY_TEMPLATE,
		candidate_beans_template,
		candidate_clusters_template,
		"c.created",
		cursor,
		"ORDER BY c.created DESC, c.cluster_id DESC",
	), params
}

// vector query for clusters - always sort by distance increasing distance
func buildVectorQueryForClusters(filters *ClusterFilters, page *PageRequest) (string, pgx.NamedArgs) {
	candidate_beans_template := `
	SELECT *, embedding <=> @embedding AS distance
	FROM trending_beans_view
	WHERE cluster_id IS NOT NULL 
	%s
	ORDER BY distance ASC
	LIMIT @candidate_limit` // %s is the placeholder for scalar filters
	candidate_clusters_template := `
	SELECT DISTINCT ON (cluster_id) *
	FROM candidate_beans
	WHERE distance <= @distance
	ORDER BY cluster_id, distance ASC`
	filter_template := ""

	where, params := buildScalarWhere(&filters.BeanFilters)
	if len(where) > 0 {
		filter_template = fmt.Sprintf("AND %s", strings.Join(where, " AND "))
	}
	candidate_beans_template = fmt.Sprintf(candidate_beans_template, filter_template)
	params["min_bean_count"] = filters.MinBeanCount
	params["limit"] = page.Limit + 1
	params["embedding"] = pgvector.NewVector(filters.Embedding)
	params["candidate_limit"] = searchCandidateLimit(page.Limit)
	params["distance"] = filters.Distance

	cursor := ""
	if page.Cursor != nil {
		cursor = "WHERE (c.distance, c.cluster_id) > (@cursor_value, @cursor_id)"
		params["cursor_value"] = *page.Cursor.Distance
		params["cursor_id"] = *page.Cursor.ID
	}

	return fmt.Sprintf(
		_CLUSTER_QUERY_TEMPLATE,
		candidate_beans_template,
		candidate_clusters_template,
		"c.distance",
		cursor,
		"ORDER BY c.distance ASC, c.cluster_id ASC",
	), params
}

// Returns a page of cluster previews.
// Preview fields: id, representative title/summary/image, first created, last created, bean count, source count.
// Categories, entities, and regions are filled later by hydrateClusters.
// Selection:
// - filter existing beans by given filters --> this contributes to the overall set to build the clusters
// - for each cluster, use most relevant (for vector search) or most recent (for scalar search) bean as representative bean
// - for each cluster, get numerical stats its global set of beans
// Pagination: page on cluster_id (NOT the filtered bean id)
func (b *PGSack) QueryClusterPreviews(ctx context.Context, filters ClusterFilters, page PageRequest) (Page[Cluster], error) {
	query, params := "", pgx.NamedArgs{}
	if len(filters.Embedding) > 0 {
		query, params = buildVectorQueryForClusters(&filters, &page)
	} else {
		query, params = buildScalarQueryForClusters(&filters, &page)
	}

	rows, err := utils.FetchAll[Cluster](ctx, b.db, query, params)
	if err != nil {
		return Page[Cluster]{}, err
	}

	return finalizePage(rows, page.Limit, func(cluster Cluster) *Cursor {
		if len(filters.Embedding) > 0 {
			return &Cursor{Version: _CURSOR_VERSION, Sort: SORT_RELEVANT, ID: &cluster.ID, Distance: &cluster.Distance.Float64}
		}
		return &Cursor{Version: _CURSOR_VERSION, Sort: SORT_RECENT, ID: &cluster.ID, Created: &cluster.Created}
	}), nil
}

func (b *PGSack) clusterExists(ctx context.Context, story_id uuid.UUID) (bool, error) {
	if story_id == uuid.Nil {
		return false, nil
	}
	return utils.FetchOneScalar[bool](
		ctx,
		b.db,
		`SELECT EXISTS(SELECT 1 FROM trend_aggregates WHERE cluster_id = @cluster_id)`,
		pgx.NamedArgs{"cluster_id": story_id},
	)
}

func (b *PGSack) GetCluster(ctx context.Context, id uuid.UUID, filters BeanFilters) (Cluster, error) {
	if id == uuid.Nil {
		return Cluster{}, ErrNonExistentID
	}
	candidate_beans_template := `
	SELECT * FROM trending_beans_view 
	WHERE cluster_id = @cluster_id %s`
	candidate_cluster_template := `
	SELECT * FROM candidate_beans 
	ORDER BY created DESC LIMIT 1`
	filter_template := ""

	where, params := buildScalarWhere(&filters)
	if len(where) > 0 {
		filter_template = fmt.Sprintf("AND %s", strings.Join(where, " AND "))
	}
	candidate_beans_template = fmt.Sprintf(candidate_beans_template, filter_template)
	params["cluster_id"] = id
	params["limit"] = 1
	params["min_bean_count"] = 1 // A direct lookup has no list threshold. NULL makes HAVING COUNT(*) >= NULL drop every group.

	query := fmt.Sprintf(
		_CLUSTER_QUERY_TEMPLATE,
		candidate_beans_template,
		candidate_cluster_template,
		"c.created",
		"",
		"",
	)
	// utils.LogQuery(query, params)
	cluster, err := utils.FetchOne[Cluster](ctx, b.db, query, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return cluster, ErrNonExistentID
	}
	if err != nil {
		// utils.LogError(err, "[DB ERROR]")
		return cluster, DBError
	}
	hydrated, err := b.hydrateClusters(ctx, []Cluster{cluster}, &filters)
	if err != nil {
		return cluster, DBError
	}
	if len(hydrated) > 0 {
		cluster = hydrated[0]
	}
	return cluster, nil
}

// hydrateClusters loads up to 3 recent member articles (one per source when possible)
// and any 10 distinct categories, entities, and regions from the filtered members.
func (b *PGSack) hydrateClusters(ctx context.Context, clusters []Cluster, filters *BeanFilters) ([]Cluster, error) {
	if len(clusters) == 0 {
		return []Cluster{}, nil
	}

	where, params := buildScalarWhere(filters)
	params["ids"] = datautils.Transform(clusters, func(cluster *Cluster) uuid.UUID { return cluster.ID })
	filter_template := ""
	if len(where) > 0 {
		filter_template = fmt.Sprintf("AND %s", strings.Join(where, " AND "))
	}

	for i := range clusters {
		clusters[i].Categories = []string{}
		clusters[i].Entities = []string{}
		clusters[i].Regions = []string{}
		clusters[i].Tags = []string{}
	}
	clusters_by_id := make(map[uuid.UUID]*Cluster, len(clusters))
	for i := range clusters {
		clusters_by_id[clusters[i].ID] = &clusters[i]
	}

	label_query := fmt.Sprintf(`
		WITH members AS (
			SELECT cluster_id, categories, entities, regions
			FROM trending_beans_view
			WHERE cluster_id = ANY(@ids)
				%s
		),
		distinct_labels AS (
			SELECT DISTINCT m.cluster_id, src.kind, value
			FROM members m
			CROSS JOIN LATERAL (
				VALUES
					('category', m.categories),
					('entity', m.entities),
					('region', m.regions)
			) AS src(kind, arr)
			CROSS JOIN LATERAL unnest(COALESCE(src.arr, '{}')) AS value
			WHERE value IS NOT NULL AND value <> ''
		),
		capped AS (
			SELECT cluster_id, kind, value
			FROM (
				SELECT
					cluster_id,
					kind,
					value,
					ROW_NUMBER() OVER (PARTITION BY cluster_id, kind) AS rn
				FROM distinct_labels
			) numbered
			WHERE rn <= 10
		)
		SELECT
			cluster_id AS id,
			COALESCE(array_agg(value) FILTER (WHERE kind = 'category'), '{}') AS categories,
			COALESCE(array_agg(value) FILTER (WHERE kind = 'entity'), '{}') AS entities,
			COALESCE(array_agg(value) FILTER (WHERE kind = 'region'), '{}') AS regions
		FROM capped
		GROUP BY cluster_id`,
		filter_template,
	)
	labels, err := utils.FetchAll[Cluster](ctx, b.db, label_query, params)
	if err != nil {
		return nil, err
	}
	for _, label := range labels {
		cluster, ok := clusters_by_id[label.ID]
		if !ok {
			continue
		}
		cluster.Categories = label.Categories
		cluster.Entities = label.Entities
		cluster.Regions = label.Regions
		cluster.Tags = ConcatArray(label.Categories, label.Entities, label.Regions)
		if cluster.Tags == nil {
			cluster.Tags = []string{}
		}
	}

	top_query := fmt.Sprintf(`
		SELECT %s
		FROM (
			SELECT *,
				ROW_NUMBER() OVER (
					PARTITION BY cluster_id
					ORDER BY source_rn ASC, created DESC, id DESC
				) AS rn
			FROM (
				SELECT *,
					ROW_NUMBER() OVER (
						PARTITION BY cluster_id, source_id
						ORDER BY created DESC, id DESC
					) AS source_rn
				FROM trending_beans_view
				WHERE cluster_id = ANY(@ids)
					%s -- filters
			) per_source
		) ranked
		WHERE rn <= 3
		ORDER BY cluster_id, rn`,
		BEAN_COLUMNS_PREVIEW,
		filter_template,
	)
	beans, err := utils.FetchAll[Bean](ctx, b.db, top_query, params)
	if err != nil {
		return nil, DBError
	}
	for _, bean := range beans {
		if cluster, ok := clusters_by_id[bean.ClusterID]; ok {
			cluster.TopBeans = append(cluster.TopBeans, bean)
		}
	}
	return clusters, nil
}

func (b *PGSack) QueryClusterMembers(ctx context.Context, cluster_id uuid.UUID, filters BeanFilters, page PageRequest, columns string) (Page[Bean], error) {
	filters.ClusterID = cluster_id
	page_out, err := b.QueryBeans(ctx, filters, page, columns)
	if err == nil && len(page_out.Items) == 0 {
		exists, err := b.clusterExists(ctx, cluster_id)
		if err != nil {
			return page_out, DBError
		}
		if !exists {
			return page_out, ErrNonExistentID
		}
	}
	return page_out, nil
}
