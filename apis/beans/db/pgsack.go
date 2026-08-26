package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
)

type PGSack struct {
	db *pgxpool.Pool
}

func NewPGSack(ctx context.Context, connString string) *PGSack {
	db, err := utils.NewConnection(ctx, connString)
	utils.NoError(err)
	return &PGSack{db: db}
}

func (p *PGSack) Close() {
	if p != nil && p.db != nil {
		p.db.Close()
	}
}
