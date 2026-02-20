package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/zagvozdeen/ola/internal/api"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/db"
	"github.com/zagvozdeen/ola/internal/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	log := logger.New(cfg)
	pool := db.New(ctx, cfg, log)
	defer pool.Close()

	api.New(cfg, log, pool).Run(ctx)
}
