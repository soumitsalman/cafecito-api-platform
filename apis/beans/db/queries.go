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
)

const (
	MIN_CANDIDATE_LIMIT = 256
)

var (
	ErrNonExistentID = errors.New("Item with this ID does not exist")
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
func buildScalarWhere(filters *Filters) ([]string, pgx.NamedArgs) {
	where := []string{}
	params := pgx.NamedArgs{}

	if len(filters.IDs) > 0 {
		where = append(where, "id = ANY(@ids)")
		params["ids"] = filters.IDs
	}
	if len(filters.URLs) > 0 {
		where = append(where, "url = ANY(@urls)")
		params["urls"] = filters.URLs
	}
	if len(filters.Sources) > 0 {
		where = append(where, "source = ANY(@sources)")
		params["sources"] = filters.Sources
	}
	if len(filters.ExcludeSources) > 0 {
		where = append(where, "source != ALL(@exclude_sources)")
		params["exclude_sources"] = filters.ExcludeSources
	}
	if len(filters.Domains) > 0 {
		where = append(where, "base_url = ANY(@domains)")
		params["domains"] = filters.Domains
	}
	if len(filters.ExcludeDomains) > 0 {
		where = append(where, "base_url != ALL(@exclude_domains)")
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
	if !filters.UpdatedFrom.IsZero() {
		where = append(where, "updated >= @updated_from")
		params["updated_from"] = filters.UpdatedFrom
	}
	if !filters.UpdatedTo.IsZero() {
		where = append(where, "updated <= @updated_to")
		params["updated_to"] = filters.UpdatedTo
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
	if filters.ClusterID != "" {
		where = append(where, "cluster_id = @cluster_id")
		params["cluster_id"] = filters.ClusterID
	}
	return where, params
}

// buildScalarOrderQuery constructs the query for latest and trending beans
// table: latest_beans_view or trending_beans_view or aggregated_beans_view.
// sort: Present in page.Cursor.Sort. page must have a cursor and a valid sort value or else the query will fail
// if filters.Embedding is present, then filter.Distance must be present. Default/not-assigned value = 0.0, returns exact match
func buildScalarOrderQuery(table string, filters *Filters, page *PageRequest, select_columns string) (string, pgx.NamedArgs) {
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
				where = append(where, fmt.Sprintf("(created, id) < (@cursor_value, @cursor_id)"))
				params["cursor_value"] = *page.Cursor.Created
				params["cursor_id"] = *page.Cursor.ID
			}
		case SORT_TRENDING:
			order_by = "trend_score DESC, id DESC"
			if page.Cursor.ID != nil {
				where = append(where, fmt.Sprintf("(trend_score, id) < (@cursor_value, @cursor_id)"))
				params["cursor_value"] = *page.Cursor.TrendScore
				params["cursor_id"] = *page.Cursor.ID
			}
		}
	}
	where_expr := ""
	if len(where) > 0 {
		where_expr = "WHERE " + strings.Join(where, " AND ")
	}
	params["limit"] = page.Limit + 1

	query := fmt.Sprintf(`
		SELECT %s -- columns
		FROM %s -- latest_beans_view or trending_beans_view or aggregated_beans_view
		%s -- WHERE
		ORDER BY %s -- order by created DESC, id DESC or trend_score DESC, id DESC
		LIMIT @limit`,
		buildSelect(select_columns, filters.FullContent),
		table,
		where_expr,
		order_by,
	)
	return query, params
}

func searchCandidateLimit(limit int) int {
	candidate_limit := limit * 4
	if candidate_limit < MIN_CANDIDATE_LIMIT {
		candidate_limit = MIN_CANDIDATE_LIMIT
	}
	return candidate_limit
}

// buildKNNSearchQuery constructs the vector (semantic) search query for both latest and trending beans
// table: latest_beans_view or trending_beans_view or aggregated_beans_view with cosine distance CTE.
// filters.Embedding, filters.Limit must be present. filters.Distance will be considered only if > 0.0
func buildKNNSearchQuery(table string, filters *Filters, page *PageRequest, columns string) (string, pgx.NamedArgs) {
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
func (b *PGSack) QueryLatestBeans(ctx context.Context, filters Filters, page PageRequest, columns string) (Page[Bean], error) {
	if page.Cursor == nil {
		page.Cursor = &Cursor{Sort: SORT_RECENT}
	}
	page.Cursor.Sort = SORT_RECENT

	query, params := buildScalarOrderQuery("latest_beans_view", &filters, &page, columns)
	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, err
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, Sort: SORT_RECENT, ID: &bean.ID, Created: &bean.Created}
	}), nil
}

// TrendingArticles returns articles ranked by trend_score descending.
// Forces sort by trend_score descending
func (b *PGSack) QueryTrendingBeans(ctx context.Context, filters Filters, page PageRequest, columns string) (Page[Bean], error) {
	if page.Cursor == nil {
		page.Cursor = &Cursor{Sort: SORT_TRENDING}
	}
	page.Cursor.Sort = SORT_TRENDING
	query, params := buildScalarOrderQuery("trending_beans_view", &filters, &page, columns)
	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, err
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, Sort: SORT_TRENDING, ID: &bean.ID, TrendScore: &bean.TrendScore.Float64}
	}), nil
}

// QueryBeans returns beans matching filters, ordered by created descending.
// Dispatches to vector (cosine distance CTE) or scalar query based on whether an embedding is present in the filters.
// Forces sort by created descending if no embedding is present.
// Forces sort by relevance if an embedding is present and filters.Distance > 0.0.
func (b *PGSack) QueryBeans(ctx context.Context, filters Filters, page PageRequest, columns string) (Page[Bean], error) {
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
		return Page[Bean]{}, err
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		if len(filters.Embedding) > 0 {
			return &Cursor{Version: _CURSOR_VERSION, ID: &bean.ID, Distance: &bean.Distance.Float64}
		} else {
			return &Cursor{Version: _CURSOR_VERSION, ID: &bean.ID, Created: &bean.Created}
		}
	}), nil
}

// // exclusiveEndOfUTCDay returns 00:00:00 UTC of the calendar day after day.
// // Use with a strict-less-than comparison so a YYYY-MM-DD `to` includes that whole UTC date.
// func exclusiveEndOfUTCDay(day time.Time) time.Time {
// 	utc := day.UTC()
// 	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
// }

func normalizeExactValues(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func (b *PGSack) getBeanURL(ctx context.Context, id uuid.UUID) (string, error) {
	seed_url, err := utils.FetchOneScalar[string](ctx, b.db,
		"SELECT url FROM beans WHERE id = @id",
		pgx.NamedArgs{"id": id},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNonExistentID
	}
	return seed_url, err
}

// QuerySimilarBeans returns known related publisher coverage for one Article UUID.
// related_beans edges are treated as undirected. Missing IDs return ErrNonExistentID.
func (b *PGSack) QuerySimilarBeans(ctx context.Context, id uuid.UUID, filters Filters, page PageRequest, columns string) (Page[Bean], error) {
	url, err := b.getBeanURL(ctx, id)
	if err != nil {
		return Page[Bean]{}, err
	}

	where, params := buildScalarWhere(&filters)
	if page.Cursor != nil && page.Cursor.Created != nil && page.Cursor.ID != nil {
		where = append(where, "(created, id) < (@cursor_created, @cursor_id)")
		params["cursor_created"] = *page.Cursor.Created
		params["cursor_id"] = *page.Cursor.ID
	}
	params["id"] = id
	params["seed_url"] = url
	params["limit"] = page.Limit + 1

	where_expr := ""
	if len(where) > 0 {
		where_expr = " AND " + strings.Join(where, " AND ")
	}
	query := fmt.Sprintf(`
		WITH target_urls AS (
			SELECT related_url AS target_url
			FROM related_beans
			WHERE url = @seed_url		
		)
		SELECT %s
		FROM target_urls
		INNER JOIN latest_beans_view ON url = target_url
		WHERE id <> @id
			%s -- additional where
		ORDER BY created DESC, id DESC
		LIMIT @limit`,
		buildSelect(columns, filters.FullContent),
		where_expr,
	)
	rows, err := utils.FetchAll[Bean](ctx, b.db, query, params)
	if err != nil {
		return Page[Bean]{}, err
	}
	return finalizePage(rows, page.Limit, func(bean Bean) *Cursor {
		return &Cursor{Version: _CURSOR_VERSION, ID: &bean.ID, Created: &bean.Created}
	}), nil
}

// QueryMentions returns the latest observed social/forum posts linking an Article URL.
// Missing IDs return ErrNonExistentID. Empty membership is an empty page.
func (b *PGSack) QueryMentions(ctx context.Context, id uuid.UUID, filters MentionFilters, page PageRequest) (Page[Mention], error) {
	seed_url, err := b.getBeanURL(ctx, id)
	if err != nil {
		return Page[Mention]{}, err
	}

	inner_where := []string{"ch.url = @url"}
	params := pgx.NamedArgs{
		"url":   seed_url,
		"limit": page.Limit + 1,
	}
	platforms := normalizeExactValues(filters.Platforms)
	if len(platforms) > 0 {
		inner_where = append(inner_where, "LOWER(ch.source) = ANY(@platforms)")
		params["platforms"] = platforms
	}
	forums := normalizeExactValues(filters.Forums)
	if len(forums) > 0 {
		inner_where = append(inner_where, "LOWER(ch.forum) = ANY(@forums)")
		params["forums"] = forums
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
				ch.chatter_url, ch.source, ch.forum, ch.collected,
				ch.likes, ch.comments, ch.subscribers
			FROM chatters ch
			WHERE %s
			ORDER BY ch.chatter_url, ch.collected DESC
		)
		SELECT chatter_url, source, forum, collected, likes, comments, subscribers
		FROM latest_chatters
		%s
		ORDER BY collected DESC, chatter_url DESC
		LIMIT @limit`,
		strings.Join(inner_where, " AND "),
		outer_where_expr,
	)
	rows, err := utils.FetchAll[Mention](ctx, b.db, query, params)
	if err != nil {
		return Page[Mention]{}, err
	}
	return finalizePage(rows, page.Limit, func(mention Mention) *Cursor {
		url := mention.URL
		observed_at := mention.ObservedAt
		return &Cursor{Version: _CURSOR_VERSION, Created: &observed_at, TextKey: &url}
	}), nil
}

// GetBean retrieves one bean record by UUID. Returns (zero, ErrInvalidID) when not found.
func (b *PGSack) GetBean(ctx context.Context, id uuid.UUID, full_content bool) (Bean, error) {
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
		return Bean{}, ErrNonExistentID
	}
	return bean, err
}

// QuerySources returns source records matching the optional query and domain filters.
func (b *PGSack) QuerySources(ctx context.Context, q string, domains []string, page PageRequest, columns string) (Page[Source], error) {
	where := []string{}
	params := pgx.NamedArgs{
		"limit": page.Limit + 1,
	}
	if q != "" {
		where = append(where, "STARTS_WITH(LOWER(site_name), @q) OR STARTS_WITH(LOWER(source), @q) OR STARTS_WITH(LOWER(base_url), @q)")
		params["q"] = strings.ToLower(strings.TrimSpace(q))
	}
	if len(domains) > 0 {
		where = append(where, "domain_name = ANY(@domains)")
		params["domains"] = domains
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
		return Page[Source]{}, err
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
		return Page[string]{}, err
	}
	return finalizePage(rows, page.Limit, func(item string) *Cursor {
		return &Cursor{
			Version: _CURSOR_VERSION,
			TextKey: &item,
		}
	}), nil
}

func coalesceStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func storyCandidateLimit(limit int) int {
	candidate_limit := limit * 16
	if candidate_limit < 128 {
		candidate_limit = 128
	}
	return candidate_limit
}

func buildStoryMemberWhere(filters *Filters) ([]string, pgx.NamedArgs) {
	where, params := buildScalarWhere(filters)
	where = append([]string{"cluster_id IS NOT NULL", "cluster_id <> ''"}, where...)
	return where, params
}

func (b *PGSack) StoryExists(ctx context.Context, story_id string) (bool, error) {
	if story_id == "" {
		return false, nil
	}
	return utils.FetchOneScalar[bool](
		ctx,
		b.db,
		`SELECT EXISTS(SELECT 1 FROM trend_aggregates WHERE cluster_id = @cluster_id)`,
		pgx.NamedArgs{"cluster_id": story_id},
	)
}

func (b *PGSack) GetStory(ctx context.Context, story_id string) (Cluster, error) {
	if story_id == "" {
		return Cluster{}, ErrNonExistentID
	}
	stories, err := b.hydrateStories(ctx, []string{story_id})
	if err != nil {
		return Cluster{}, err
	}
	if len(stories) == 0 {
		return Cluster{}, ErrNonExistentID
	}
	return stories[0], nil
}

func (b *PGSack) QueryStories(ctx context.Context, filters Filters, page PageRequest) (Page[Cluster], error) {
	var rows []clusterBase
	var err error
	if len(filters.Embedding) > 0 {
		rows, err = b.queryStoryPageByRelevance(ctx, &filters, &page)
	} else {
		rows, err = b.queryStoryPageByRecency(ctx, &filters, &page)
	}
	if err != nil {
		return Page[Cluster]{}, err
	}

	paged := finalizePage(rows, page.Limit, func(row clusterBase) *Cursor {
		id := row.ID
		if len(filters.Embedding) > 0 {
			distance := row.Distance.Float64
			return &Cursor{Version: _CURSOR_VERSION, TextKey: &id, Distance: &distance}
		}
		created := row.LastCreated
		return &Cursor{Version: _CURSOR_VERSION, TextKey: &id, Created: &created}
	})

	ids := make([]string, 0, len(paged.Items))
	for _, row := range paged.Items {
		ids = append(ids, row.ID)
	}
	stories, err := b.hydrateStories(ctx, ids)
	if err != nil {
		return Page[Cluster]{}, err
	}
	return Page[Cluster]{Items: stories, NextCursor: paged.NextCursor}, nil
}

func (b *PGSack) queryStoryPageByRecency(ctx context.Context, filters *Filters, page *PageRequest) ([]clusterBase, error) {
	where, params := buildStoryMemberWhere(filters)
	cursor_where := ""
	if page.Cursor != nil && page.Cursor.Created != nil && page.Cursor.TextKey != nil {
		cursor_where = "WHERE (last_created, id) < (@cursor_created, @cursor_id)"
		params["cursor_created"] = *page.Cursor.Created
		params["cursor_id"] = *page.Cursor.TextKey
	}
	params["min_bean_count"] = filters.MinBeanCount
	params["limit"] = page.Limit + 1

	query := fmt.Sprintf(`
		WITH matched AS (
			SELECT DISTINCT cluster_id
			FROM trending_beans_view
			WHERE %s
		),
		stats AS (
			SELECT
				b.cluster_id AS id,
				MAX(b.created) AS last_created
			FROM trending_beans_view b
			INNER JOIN matched m ON m.cluster_id = b.cluster_id
			GROUP BY b.cluster_id
			HAVING COUNT(*) >= @min_bean_count
		)
		SELECT * FROM stats
		%s
		ORDER BY last_created DESC, id DESC
		LIMIT @limit`,
		strings.Join(where, " AND "),
		cursor_where,
	)
	return utils.FetchAll[clusterBase](ctx, b.db, query, params)
}

func (b *PGSack) queryStoryPageByRelevance(ctx context.Context, filters *Filters, page *PageRequest) ([]clusterBase, error) {
	where, params := buildStoryMemberWhere(filters)
	cursor_where := ""
	if page.Cursor != nil && page.Cursor.Distance != nil && page.Cursor.TextKey != nil {
		cursor_where = "WHERE (distance, id) > (@cursor_distance, @cursor_id)"
		params["cursor_distance"] = *page.Cursor.Distance
		params["cursor_id"] = *page.Cursor.TextKey
	}
	distance_having := ""
	if filters.Distance > 0 {
		distance_having = "HAVING MIN(distance) <= @distance"
		params["distance"] = filters.Distance
	}
	params["embedding"] = pgvector.NewVector(filters.Embedding)
	params["candidate_limit"] = storyCandidateLimit(page.Limit)
	params["min_bean_count"] = filters.MinBeanCount
	params["limit"] = page.Limit + 1

	query := fmt.Sprintf(`
		WITH nearest AS MATERIALIZED (
			SELECT cluster_id, embedding <=> @embedding AS distance
			FROM trending_beans_view
			WHERE %s
			ORDER BY distance ASC
			LIMIT @candidate_limit
		),
		matched AS (
			SELECT cluster_id, MIN(distance) AS distance
			FROM nearest
			GROUP BY cluster_id
			%s
		),
		stats AS (
			SELECT
				m.cluster_id AS id,
				m.distance AS distance,
				MAX(b.created) AS last_created
			FROM matched m
			INNER JOIN trending_beans_view b ON b.cluster_id = m.cluster_id
			GROUP BY m.cluster_id, m.distance
			HAVING COUNT(*) >= @min_bean_count
		)
		SELECT * FROM stats
		%s
		ORDER BY distance ASC, id ASC
		LIMIT @limit`,
		strings.Join(where, " AND "),
		distance_having,
		cursor_where,
	)
	return utils.FetchAll[clusterBase](ctx, b.db, query, params)
}

func (b *PGSack) hydrateStories(ctx context.Context, ids []string) ([]Cluster, error) {
	if len(ids) == 0 {
		return []Cluster{}, nil
	}
	params := pgx.NamedArgs{"ids": ids}
	stats_query := `
		WITH members AS MATERIALIZED (
			SELECT tr.cluster_id, b.created, b.source, b.categories, b.regions, b.entities
			FROM beans b
			INNER JOIN trend_aggregates tr ON b.url = tr.url
			WHERE tr.cluster_id = ANY(@ids)
		),
		stats AS (
			SELECT
				cluster_id AS id,
				MIN(created) AS first_created,
				MAX(created) AS last_created,
				COUNT(*)::int AS bean_count,
				COUNT(DISTINCT source)::int AS source_count
			FROM members
			GROUP BY cluster_id
		),
		label_freq AS (
			SELECT m.cluster_id, lbl.type, lbl.val, COUNT(*) AS cnt
			FROM members m
			CROSS JOIN LATERAL (
				SELECT 'category' AS type, v AS val FROM unnest(COALESCE(m.categories, '{}')) AS v
				UNION ALL
				SELECT 'region', v FROM unnest(COALESCE(m.regions, '{}')) AS v
				UNION ALL
				SELECT 'entity', v FROM unnest(COALESCE(m.entities, '{}')) AS v
			) lbl
			GROUP BY m.cluster_id, lbl.type, lbl.val
		),
		ranked AS (
			SELECT cluster_id, type, val, cnt,
				ROW_NUMBER() OVER (PARTITION BY cluster_id, type ORDER BY cnt DESC, val ASC) AS rn
			FROM label_freq
		),
		cat_agg AS (
			SELECT cluster_id, ARRAY_AGG(val ORDER BY cnt DESC, val ASC) AS categories
			FROM ranked WHERE type = 'category' AND rn <= 10 GROUP BY cluster_id
		),
		region_agg AS (
			SELECT cluster_id, ARRAY_AGG(val ORDER BY cnt DESC, val ASC) AS regions
			FROM ranked WHERE type = 'region' AND rn <= 10 GROUP BY cluster_id
		),
		entity_agg AS (
			SELECT cluster_id, ARRAY_AGG(val ORDER BY cnt DESC, val ASC) AS entities
			FROM ranked WHERE type = 'entity' AND rn <= 10 GROUP BY cluster_id
		),
		tag_freq AS (
			SELECT cluster_id, val, SUM(cnt) AS cnt
			FROM label_freq GROUP BY cluster_id, val
		),
		tag_ranked AS (
			SELECT cluster_id, val, cnt,
				ROW_NUMBER() OVER (PARTITION BY cluster_id ORDER BY cnt DESC, val ASC) AS rn
			FROM tag_freq
		),
		tag_agg AS (
			SELECT cluster_id, ARRAY_AGG(val ORDER BY cnt DESC, val ASC) AS tags
			FROM tag_ranked WHERE rn <= 10 GROUP BY cluster_id
		)
		SELECT
			s.id,
			s.first_created,
			s.last_created,
			s.bean_count,
			s.source_count,
			COALESCE(c.categories, '{}') AS categories,
			COALESCE(r.regions, '{}') AS regions,
			COALESCE(e.entities, '{}') AS entities,
			COALESCE(t.tags, '{}') AS tags
		FROM stats s
		LEFT JOIN cat_agg c ON c.cluster_id = s.id
		LEFT JOIN region_agg r ON r.cluster_id = s.id
		LEFT JOIN entity_agg e ON e.cluster_id = s.id
		LEFT JOIN tag_agg t ON t.cluster_id = s.id`

	stats, err := utils.FetchAll[Cluster](ctx, b.db, stats_query, params)
	if err != nil {
		return nil, err
	}
	by_id := make(map[string]Cluster, len(stats))
	for _, story := range stats {
		story.Categories = coalesceStrings(story.Categories)
		story.Regions = coalesceStrings(story.Regions)
		story.Entities = coalesceStrings(story.Entities)
		story.Tags = coalesceStrings(story.Tags)
		by_id[story.ID] = story
	}

	top_query := fmt.Sprintf(`
		SELECT %s
		FROM (
			SELECT %s,
				ROW_NUMBER() OVER (
					PARTITION BY cluster_id
					ORDER BY trend_score DESC NULLS LAST, created DESC, id DESC
				) AS rn
			FROM trending_beans_view
			WHERE cluster_id = ANY(@ids)
		) ranked
		WHERE rn <= 3
		ORDER BY cluster_id, rn`,
		BEAN_COLUMNS_WITHOUT_TREND,
		BEAN_COLUMNS_WITHOUT_TREND,
	)
	articles, err := utils.FetchAll[Bean](ctx, b.db, top_query, params)
	if err != nil {
		return nil, err
	}
	grouped := make(map[string][]Bean, len(ids))
	for _, article := range articles {
		if !article.ClusterID.Valid {
			continue
		}
		grouped[article.ClusterID.String] = append(grouped[article.ClusterID.String], article)
	}

	stories := make([]Cluster, 0, len(ids))
	for _, id := range ids {
		story, ok := by_id[id]
		if !ok {
			continue
		}
		story.TopArticles = grouped[id]
		if story.TopArticles == nil {
			story.TopArticles = []Bean{}
		}
		for _, article := range story.TopArticles {
			if article.Title.Valid && article.Title.String != "" {
				story.Title = article.Title.String
				break
			}
		}
		stories = append(stories, story)
	}
	return stories, nil
}
