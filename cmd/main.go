package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/zagvozdeen/ola/internal/api"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/db"
	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/store"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New(*cfgPath)
	log := logger.New(cfg)
	defer log.Close()
	pool := db.New(ctx, cfg, log)
	defer pool.Close()
	storage := store.New(log, pool)

	api.New(cfg, log, storage).Run(ctx)
}
