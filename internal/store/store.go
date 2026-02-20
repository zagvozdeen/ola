package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zagvozdeen/ola/internal/logger"
)

type Storage interface {
}

type Store struct {
	log  *logger.Logger
	pool *pgxpool.Pool
}

var _ Storage = (*Store)(nil)

func New(log *logger.Logger, pool *pgxpool.Pool) *Store {
	return &Store{log: log, pool: pool}
}
