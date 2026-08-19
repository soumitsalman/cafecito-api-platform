package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

const (
	_TIMEOUT        = 10
	_POOL_SIZE      = 32
	_CONN_LIFETIME  = 5
	_CONN_IDLE_TIME = 5
)

const _SQL_HNSW_SETTINGS = `
SET hnsw.iterative_scan = strict_order;
SET hnsw.ef_search = 100;
`

func NewConnection(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// PGX fetch helpers
func FetchOne[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) (T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		var zero T
		return zero, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[T])
}

func FetchOneScalar[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) (T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		var zero T
		return zero, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowTo[T])
}

func FetchAll[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) ([]T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[T])
}

func FetchAllScalar[T any](ctx context.Context, db *pgxpool.Pool, query string, args pgx.NamedArgs) ([]T, error) {
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowTo[T])
}
