package db

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/logger"
)

//go:embed migrations
var fs embed.FS

func New(ctx context.Context, cfg *config.Config, log *logger.Logger) *pgxpool.Pool {
	pool, err := connect(ctx, cfg, log)
	if err != nil {
		log.Error("Fatal error: failed to connect to db", err)
		os.Exit(1)
	}
	return pool
}

func connect(ctx context.Context, cfg *config.Config, log *logger.Logger) (*pgxpool.Pool, error) {
	connCfg, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	))
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, connCfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	err = migrate(ctx, log, pool)
	if err != nil {
		return nil, fmt.Errorf("run postgres migrations: %w", err)
	}
	return pool, nil
}

func migrate(ctx context.Context, log goose.Logger, pool *pgxpool.Pool) (err error) {
	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if closeErr := db.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	goose.SetBaseFS(fs)
	goose.SetLogger(log)

	err = goose.SetDialect("pgx")
	if err != nil {
		return err
	}

	//err = goose.DownContext(ctx, db, "migrations")
	//if err != nil {
	//	return fmt.Errorf("failed to down migrations: %w", err)
	//}

	err = goose.UpContext(ctx, db, "migrations")
	if err != nil {
		return err
	}

	return err
}
