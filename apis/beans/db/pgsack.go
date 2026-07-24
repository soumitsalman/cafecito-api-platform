package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
	"github.com/soumitsalman/cafecito-api-platform/apis/internal/config"
	datautils "github.com/soumitsalman/data-utils"
)

const (
	_TIMEOUT        = 10
	_POOL_SIZE      = 32
	_CONN_LIFETIME  = 5
	_CONN_IDLE_TIME = 5
)

const (
	_TRENDING_BEANS_VIEW   = "trending_beans_view"
	_AGGREGATED_BEANS_VIEW = "aggregated_beans_view"
)

// Name of mandatory tables.
const (
	BEANS            = "beans"
	PUBLISHERS       = "publishers"
	CHATTERS         = "chatters"
	RELATED_BEANS    = "related_beans"
	BEAN_RELATIONS   = "bean_relations"
	FIXED_CATEGORIES = "fixed_categories"
	FIXED_SENTIMENTS = "fixed_sentiments"
)

const (
	CORE_BEAN_FIELDS                = "url, kind, title, summary, author, source, image_url, created, categories, sentiments, regions, entities"
	CORE_PUBLISHER_FIELDS           = "source, base_url, site_name, description, favicon"
	EXTENDED_BEAN_FIELDS            = "url, kind, title, summary, author, source, image_url, created, categories, sentiments, regions, entities, base_url, site_name, description, favicon, likes, comments, shares, related"
	UNRESTRICTED_CONTENT_CONDITIONS = "restricted_content IS NULL AND content IS NOT NULL"
	ORDER_BY_LATEST                 = "created DESC"
	ORDER_BY_TRENDING               = "trend_score DESC"
	ORDER_BY_DISTANCE               = "distance ASC"
)

var (
	ErrNotImplemented         = errors.New("method not implemented")
	ErrVectorDistanceRequired = errors.New("vector counts require a distance threshold")
)

type Condition struct {
	URLs       []string
	Kind       []string
	Created    time.Time
	Updated    time.Time
	Collected  time.Time
	Categories []string
	Regions    []string
	Entities   []string
	Tags       []string
	Sources    []string
	Embedding  []float32
	Distance   *float64
	Extra      []string // CAUTION: This is a catch-all for any additional conditions. Use with care to avoid SQL injection.
}

type Pagination struct {
	Limit  int
	Offset int
}

type Beansack interface {
	QueryBeans(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]BeanAggregate, error)
	QueryLatestBeans(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]Bean, error)
	QueryTrendingBeans(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]BeanTrend, error)
	QueryPublishers(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]Publisher, error)
	QueryChatters(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]Chatter, error)
	QueryPropagation(ctx context.Context, urls []string) ([]PropagationResult, error)

	DistinctCategories(ctx context.Context, page Pagination) ([]string, error)
	DistinctSentiments(ctx context.Context, page Pagination) ([]string, error)
	DistinctEntities(ctx context.Context, page Pagination) ([]string, error)
	DistinctRegions(ctx context.Context, page Pagination) ([]string, error)
	DistinctSources(ctx context.Context, page Pagination) ([]string, error)

	CountRows(ctx context.Context, table string, conditions Condition) (int64, error)
	Close()
}

type PGSack struct {
	db *pgxpool.Pool
}

const _SQL_HNSW_SETTINGS = `
SET hnsw.iterative_scan = strict_order;
SET hnsw.ef_search = 100;
`

func NewPGSack(ctx context.Context, connString string) *PGSack {
	config, err := pgxpool.ParseConfig(connString)
	NoError(err)

	config.MinConns = 1
	config.MaxConns = _POOL_SIZE
	config.MaxConnLifetime = time.Minute * _CONN_LIFETIME
	config.MaxConnIdleTime = time.Minute * _CONN_IDLE_TIME
	config.HealthCheckPeriod = time.Minute * _CONN_LIFETIME
	config.ConnConfig.ConnectTimeout = time.Minute * _TIMEOUT
	if config.ConnConfig.RuntimeParams == nil {
		config.ConnConfig.RuntimeParams = map[string]string{}
	}
	config.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprintf("%dmin", _TIMEOUT)
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := pgxvec.RegisterTypes(ctx, conn); err != nil {
			return err
		}
		_, err := conn.Exec(ctx, _SQL_HNSW_SETTINGS)
		return err
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	NoError(err)
	NoError(db.Ping(ctx)) // Quick health check on startup

	return &PGSack{db: db}
}

func (p *PGSack) QueryBeans(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]BeanAggregate, error) {
	items, err := fetchBeans(ctx, p, _AGGREGATED_BEANS_VIEW, conditions, nil, page, columns)
	if err != nil {
		return nil, err
	}
	return datautils.Transform(items, func(item *dataRow) BeanAggregate { return item.toBeanAggregate() }), nil
}

func (p *PGSack) QueryLatestBeans(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]Bean, error) {
	items, err := fetchBeans(ctx, p, BEANS, conditions, []string{ORDER_BY_LATEST}, page, columns)
	if err != nil {
		return nil, err
	}
	return datautils.Transform(items, func(item *dataRow) Bean { return item.toBean() }), nil
}

func (p *PGSack) QueryTrendingBeans(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]BeanTrend, error) {
	items, err := fetchBeans(ctx, p, _TRENDING_BEANS_VIEW, conditions, []string{ORDER_BY_TRENDING}, page, columns)
	if err != nil {
		return nil, err
	}
	return datautils.Transform(items, func(item *dataRow) BeanTrend { return item.toBeanTrend() }), nil
}

// TODO: how to pass in text/tag/keyword based search
func (p *PGSack) QueryPublishers(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]Publisher, error) {
	query, args := p.BuildScalarSQL(PUBLISHERS, conditions, nil, page, columns)
	items, err := fetchAll[dataRow](ctx, p.db, query, args)
	if err != nil {
		return nil, err
	}
	return datautils.Transform(items, func(item *dataRow) Publisher { return item.toPublisher() }), nil
}

func (p *PGSack) QueryChatters(ctx context.Context, conditions Condition, page Pagination, columns []string) ([]Chatter, error) {
	query, args := p.BuildScalarSQL(CHATTERS, conditions, nil, page, columns)
	items, err := fetchAll[dataRow](ctx, p.db, query, args)
	if err != nil {
		return nil, err
	}
	return datautils.Transform(items, func(item *dataRow) Chatter { return item.toChatter() }), nil
}

const _SQL_PROPAGATION_RELATED = `
SELECT
    rb.url AS seed_url,
    b.url,
    b.created,
    b.source,
    p.site_name
FROM related_beans rb
INNER JOIN beans b ON b.url = rb.related_url
LEFT JOIN publishers p ON p.source = b.source
WHERE rb.url = ANY(@urls)
ORDER BY rb.url, b.created DESC`

const _SQL_PROPAGATION_SHARES = `
SELECT
    c.url AS seed_url,
    c.chatter_url,
    c.source,
    c.forum,
    c.collected,
    c.comments,
    c.likes
FROM chatters c
WHERE c.url = ANY(@urls)
ORDER BY c.url, c.collected DESC`

type propagationRelatedRow struct {
	SeedURL  string         `db:"seed_url"`
	URL      string         `db:"url"`
	Created  sql.NullTime   `db:"created"`
	Source   string         `db:"source"`
	SiteName sql.NullString `db:"site_name"`
}

type propagationShareRow struct {
	SeedURL    string         `db:"seed_url"`
	ChatterURL string         `db:"chatter_url"`
	Source     string         `db:"source"`
	Forum      sql.NullString `db:"forum"`
	Collected  sql.NullTime   `db:"collected"`
	Comments   sql.NullInt64  `db:"comments"`
	Likes      sql.NullInt64  `db:"likes"`
}

func (p *PGSack) QueryPropagation(ctx context.Context, urls []string) ([]PropagationResult, error) {
	args := pgx.NamedArgs{"urls": urls}

	relatedRows, err := fetchAll[propagationRelatedRow](ctx, p.db, _SQL_PROPAGATION_RELATED, args)
	if err != nil {
		return nil, err
	}

	shareRows, err := fetchAll[propagationShareRow](ctx, p.db, _SQL_PROPAGATION_SHARES, args)
	if err != nil {
		return nil, err
	}

	coverageByURL := map[string][]PropagationCoverage{}
	for _, r := range relatedRows {
		coverageByURL[r.SeedURL] = append(coverageByURL[r.SeedURL], PropagationCoverage{
			URL:      r.URL,
			Created:  r.Created.Time,
			Source:   r.Source,
			SiteName: r.SiteName.String,
		})
	}

	mentionsByURL := map[string][]PropagationMention{}
	for _, s := range shareRows {
		mentionsByURL[s.SeedURL] = append(mentionsByURL[s.SeedURL], PropagationMention{
			ShareURL: s.ChatterURL,
			Source:   s.Source,
			Forum:    s.Forum.String,
			Observed: s.Collected.Time,
			Comments: s.Comments.Int64,
			Likes:    s.Likes.Int64,
		})
	}

	results := make([]PropagationResult, 0, len(urls))
	for _, url := range urls {
		cov := coverageByURL[url]
		men := mentionsByURL[url]
		if cov == nil {
			cov = []PropagationCoverage{}
		}
		if men == nil {
			men = []PropagationMention{}
		}
		results = append(results, PropagationResult{URL: url, Coverage: cov, Mentions: men})
	}
	return results, nil
}

func (p *PGSack) DistinctCategories(ctx context.Context, page Pagination) ([]string, error) {
	// SELECT DISTINCT unnest(categories) AS category FROM beans WHERE categories IS NOT NULL ORDER BY category
	query, args := p.BuildScalarSQL(BEANS, Condition{Extra: []string{"categories IS NOT NULL"}}, []string{"category"}, page, []string{"DISTINCT unnest(categories) AS category"})
	return fetchAllScalar[string](ctx, p.db, query, args)
}

func (p *PGSack) DistinctSentiments(ctx context.Context, page Pagination) ([]string, error) {
	query, args := p.BuildScalarSQL(BEANS, Condition{Extra: []string{"sentiments IS NOT NULL"}}, []string{"sentiment"}, page, []string{"DISTINCT unnest(sentiments) AS sentiment"})
	return fetchAllScalar[string](ctx, p.db, query, args)
}

func (p *PGSack) DistinctEntities(ctx context.Context, page Pagination) ([]string, error) {
	query, args := p.BuildScalarSQL(BEANS, Condition{Extra: []string{"entities IS NOT NULL"}}, []string{"entity"}, page, []string{"DISTINCT unnest(entities) AS entity"})
	return fetchAllScalar[string](ctx, p.db, query, args)
}

func (p *PGSack) DistinctRegions(ctx context.Context, page Pagination) ([]string, error) {
	query, args := p.BuildScalarSQL(BEANS, Condition{Extra: []string{"regions IS NOT NULL"}}, []string{"region"}, page, []string{"DISTINCT unnest(regions) AS region"})
	return fetchAllScalar[string](ctx, p.db, query, args)
}

func (p *PGSack) DistinctSources(ctx context.Context, page Pagination) ([]string, error) {
	query, args := p.BuildScalarSQL(PUBLISHERS, Condition{}, []string{"source"}, page, []string{"source"})
	return fetchAllScalar[string](ctx, p.db, query, args)
}

func (p *PGSack) CountRows(ctx context.Context, table string, conditions Condition) (int64, error) {
	if len(conditions.Embedding) > 0 {
		if conditions.Distance == nil {
			return 0, ErrVectorDistanceRequired
		}
		query, args := p.BuildVectorCountSQL(table, conditions)
		return fetchOneScalar[int64](ctx, p.db, query, args)
	}
	query, args := p.BuildScalarSQL(table, conditions, nil, Pagination{}, []string{"count(*)"})
	return fetchOneScalar[int64](ctx, p.db, query, args)
}

func (p *PGSack) Close() {
	if p != nil && p.db != nil {
		p.db.Close()
	}
}

// SQL query string builder utilities
// BuildScalarSQL creates a non-vector query from scalar filters, ordering, and pagination.
func (p *PGSack) BuildScalarSQL(table string, conditions Condition, orders []string, page Pagination, columns []string) (string, pgx.NamedArgs) {
	// where clause first - because we may need it before select
	where_expr, where_params := p.buildWhereExpr(conditions)

	// select fields
	fields := buildSQLFields(columns)

	base_expr := fmt.Sprintf("SELECT %s FROM %s %s", fields, table, where_expr)
	builder := strings.Builder{}
	builder.WriteString(base_expr)

	// orders
	if len(orders) > 0 {
		builder.WriteString(" ")
		builder.WriteString(buildOrderBy(orders...))
	}
	// pagination
	page_expr, page_params := buildPaginationExpr(page)
	if page_expr != "" {
		builder.WriteString(" ")
		builder.WriteString(page_expr)
	}
	// LogQuery(builder.String(), mergeParams(where_params, page_params))
	return builder.String(), mergeParams(where_params, page_params)
}

func buildSQLFields(columns []string) string {
	if len(columns) == 0 {
		return "*"
	}
	return strings.Join(columns, ", ")
}

// BuildVectorSQL creates an HNSW-compatible nearest-neighbor query with an optional distance cutoff.
func (p *PGSack) BuildVectorSQL(table string, conditions Condition, orders []string, page Pagination, fields string) (string, pgx.NamedArgs) {
	where_expr, where_params := p.buildWhereExpr(conditions)

	distance_expr := "embedding <=> @embedding::vector"
	where_params["embedding"] = pgvector.NewVector(conditions.Embedding)
	where_params["candidate_limit"] = config.VECTOR_QUERY_DEFAULT_CANDIDATE_LIMIT + ((page.Offset + page.Limit) * config.VECTOR_QUERY_CANDIDATE_LIMIT_MULTIPLIER)

	if orders == nil {
		orders = []string{ORDER_BY_DISTANCE}
	} else {
		orders = append(orders, ORDER_BY_DISTANCE)
	}

	builder := strings.Builder{}
	builder.WriteString("WITH nearest_results AS MATERIALIZED (\n")
	builder.WriteString("SELECT *, (")
	builder.WriteString(distance_expr)
	builder.WriteString(") AS distance\nFROM ")
	builder.WriteString(table)
	if where_expr != "" {
		builder.WriteString("\n")
		builder.WriteString(where_expr)
	}
	builder.WriteString("\nORDER BY ")
	builder.WriteString(distance_expr)
	builder.WriteString(" ASC\nLIMIT @candidate_limit\n)\nSELECT ")
	builder.WriteString(fields)
	builder.WriteString("\nFROM nearest_results")
	if conditions.Distance != nil {
		builder.WriteString("\nWHERE distance <= @distance")
		where_params["distance"] = *conditions.Distance
	}
	builder.WriteString("\n")
	builder.WriteString(buildOrderBy(orders...))

	page_expr, page_params := buildPaginationExpr(page)
	if page_expr != "" {
		builder.WriteString(" ")
		builder.WriteString(page_expr)
	}
	// LogQuery(builder.String(), mergeParams(candidate_params, where_params, page_params))
	return builder.String(), mergeParams(where_params, page_params)
}

// BuildVectorCountSQL creates an exact vector-distance count query without ANN ordering.
func (p *PGSack) BuildVectorCountSQL(table string, conditions Condition) (string, pgx.NamedArgs) {
	where_expr, where_params := p.buildWhereExpr(conditions)
	if conditions.Distance != nil {
		if where_expr == "" {
			where_expr = "WHERE (embedding <=> @embedding::vector) <= @distance"
		} else {
			where_expr += " AND (embedding <=> @embedding::vector) <= @distance"
		}
		where_params = mergeParams(where_params, pgx.NamedArgs{
			"embedding": pgvector.NewVector(conditions.Embedding),
			"distance":  *conditions.Distance,
		})
	}
	return fmt.Sprintf("SELECT count(*) FROM %s %s", table, where_expr), where_params
}

func mergeParams(maps ...pgx.NamedArgs) pgx.NamedArgs {
	merged := pgx.NamedArgs{}
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func (p *PGSack) buildWhereExpr(conditions Condition) (string, pgx.NamedArgs) {
	parts := make([]string, 0, 10) // preallocate for expected conditions
	args := pgx.NamedArgs{}

	if len(conditions.URLs) > 0 {
		parts = append(parts, "url = ANY(@urls)")
		args["urls"] = conditions.URLs // pgx handles []string as array automatically
	}

	if len(conditions.Kind) > 0 {
		parts = append(parts, "kind = ANY(@kind)")
		args["kind"] = conditions.Kind
	}

	if !conditions.Created.IsZero() {
		parts = append(parts, "created >= @created_from")
		args["created_from"] = conditions.Created
	}

	if !conditions.Collected.IsZero() {
		parts = append(parts, "collected >= @collected_from")
		args["collected_from"] = conditions.Collected
	}

	if !conditions.Updated.IsZero() {
		parts = append(parts, "updated >= @updated_from")
		args["updated_from"] = conditions.Updated
	}

	if len(conditions.Categories) > 0 {
		parts = append(parts, "categories && @categories")
		args["categories"] = conditions.Categories
	}

	if len(conditions.Regions) > 0 {
		parts = append(parts, "regions && @regions")
		args["regions"] = conditions.Regions
	}

	if len(conditions.Entities) > 0 {
		parts = append(parts, "entities && @entities")
		args["entities"] = conditions.Entities
	}

	if len(conditions.Tags) > 0 {
		parts = append(parts, "tags @@ plainto_tsquery('simple', @tags_query)")
		args["tags_query"] = strings.Join(conditions.Tags, " & ") // "tag1 & tag2 & tag3"
	}

	if len(conditions.Sources) > 0 {
		parts = append(parts, "source = ANY(@sources)")
		args["sources"] = conditions.Sources
	}

	if len(conditions.Extra) > 0 {
		parts = append(parts, conditions.Extra...)
	}

	if len(parts) == 0 {
		return "", nil
	}
	return fmt.Sprintf("WHERE %s", strings.Join(parts, " AND ")), args
}

func buildOrderBy(order ...string) string {
	if len(order) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(order, ", ")
}

func buildPaginationExpr(page Pagination) (string, pgx.NamedArgs) {
	parts := make([]string, 0, 2)
	args := pgx.NamedArgs{}

	if page.Limit > 0 {
		parts = append(parts, "LIMIT @limit")
		args["limit"] = page.Limit
	}
	if page.Offset > 0 {
		parts = append(parts, "OFFSET @offset")
		args["offset"] = page.Offset
	}

	if len(parts) == 0 {
		return "", nil
	}
	return strings.Join(parts, " "), args
}

// PGX fetch helpers
func fetchOne[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) (T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		var zero T
		return zero, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[T])
}

func fetchOneScalar[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) (T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		var zero T
		return zero, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowTo[T])
}

func fetchAll[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) ([]T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[T])
}

func fetchAllScalar[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) ([]T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowTo[T])
}

func fetchBeans(ctx context.Context, p *PGSack, table string, conditions Condition, orders []string, page Pagination, columns []string) ([]dataRow, error) {
	var query string
	var args pgx.NamedArgs
	if len(conditions.Embedding) > 0 {

		query, args = p.BuildVectorSQL(table, conditions, orders, page, buildSQLFields(columns))
	} else {
		query, args = p.BuildScalarSQL(table, conditions, orders, page, columns)
	}
	return fetchAll[dataRow](ctx, p.db, query, args)
}

// PGX marshalling and unmarshalling for custom types

type dataRow struct {
	URL               sql.NullString  `db:"url"`
	Kind              sql.NullString  `db:"kind"`
	Title             sql.NullString  `db:"title"`
	Summary           sql.NullString  `db:"summary"`
	Content           sql.NullString  `db:"content"`
	RestrictedContent sql.NullBool    `db:"restricted_content"`
	Author            sql.NullString  `db:"author"`
	Source            sql.NullString  `db:"source"`
	ImageUrl          sql.NullString  `db:"image_url"`
	Created           sql.NullTime    `db:"created"`
	Embedding         pgvector.Vector `db:"embedding"`
	Categories        []string        `db:"categories"`
	Sentiments        []string        `db:"sentiments"`
	Regions           []string        `db:"regions"`
	Entities          []string        `db:"entities"`
	Tags              []byte          `db:"tags"`
	Related           sql.NullInt64   `db:"related"`
	ClusterSize       sql.NullInt64   `db:"cluster_size"`
	Updated           sql.NullTime    `db:"updated"`
	Likes             sql.NullInt64   `db:"likes"`
	Comments          sql.NullInt64   `db:"comments"`
	Subscribers       sql.NullInt64   `db:"subscribers"`
	Shares            sql.NullInt64   `db:"shares"`
	Distance          float64         `db:"distance"`
	TrendScore        float64         `db:"trend_score"`
	ChatterURL        sql.NullString  `db:"chatter_url"`
	Forum             sql.NullString  `db:"forum"`
	Collected         sql.NullTime    `db:"collected"`
	BaseURL           sql.NullString  `db:"base_url"`
	SiteName          sql.NullString  `db:"site_name"`
	Description       sql.NullString  `db:"description"`
	Favicon           sql.NullString  `db:"favicon"`
	RSSFeed           sql.NullString  `db:"rss_feed"`
}

// Conversion methods from dataRow to public types
func (r *dataRow) toBean() Bean {
	return Bean{
		URL:        r.URL.String,
		Kind:       r.Kind.String,
		Title:      r.Title.String,
		Summary:    r.Summary.String,
		Content:    r.Content.String,
		Author:     r.Author.String,
		Source:     r.Source.String,
		ImageUrl:   r.ImageUrl.String,
		Created:    r.Created.Time,
		Embedding:  r.Embedding.Slice(),
		Categories: r.Categories,
		Sentiments: r.Sentiments,
		Regions:    r.Regions,
		Entities:   r.Entities,
		MergedTags: ConcatArray[string](r.Categories, r.Regions, r.Entities),
	}
}

func (r *dataRow) toBeanTrend() BeanTrend {
	return BeanTrend{
		Bean:        r.toBean(),
		Likes:       r.Likes.Int64,
		Comments:    r.Comments.Int64,
		Subscribers: r.Subscribers.Int64,
		Shares:      r.Shares.Int64,
		Related:     r.Related.Int64,
		Updated:     r.Updated.Time,
		TrendScore:  r.TrendScore,
	}
}

func (r *dataRow) toBeanAggregate() BeanAggregate {
	return BeanAggregate{
		BeanTrend:   r.toBeanTrend(),
		BaseURL:     r.BaseURL.String,
		SiteName:    r.SiteName.String,
		Description: r.Description.String,
		Favicon:     r.Favicon.String,
	}
}

func (r *dataRow) toPublisher() Publisher {
	return Publisher{
		Source:      r.Source.String,
		BaseURL:     r.BaseURL.String,
		SiteName:    r.SiteName.String,
		Description: r.Description.String,
		Favicon:     r.Favicon.String,
		RSSFeed:     r.RSSFeed.String,
		Collected:   r.Collected.Time,
	}
}

func (r *dataRow) toChatter() Chatter {
	return Chatter{
		ChatterURL:  r.ChatterURL.String,
		URL:         r.URL.String,
		Source:      r.Source.String,
		Forum:       r.Forum.String,
		Collected:   r.Collected.Time,
		Likes:       r.Likes.Int64,
		Comments:    r.Comments.Int64,
		Subscribers: r.Subscribers.Int64,
	}
}

func (r *dataRow) toChatterAggregate() ChatterAggregate {
	return ChatterAggregate{
		URL:         r.URL.String,
		Collected:   r.Collected.Time,
		Likes:       r.Likes.Int64,
		Comments:    r.Comments.Int64,
		Subscribers: r.Subscribers.Int64,
		Shares:      r.Shares.Int64,
	}
}
