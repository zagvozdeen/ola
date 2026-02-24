package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zagvozdeen/ola/internal/logger"
)

type Storage interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context)
	Rollback(ctx context.Context)
}

type Store struct {
	log  *logger.Logger
	pool *pgxpool.Pool
}

var _ Storage = (*Store)(nil)

func New(log *logger.Logger, pool *pgxpool.Pool) *Store {
	return &Store{log: log, pool: pool}
}
