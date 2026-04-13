package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"first/internal/cache"
	"first/internal/config"
	"first/internal/db"
	"first/internal/pkg"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	logger := pkg.NewLogger(cfg.LogLevel)

	logger.Info("config loaded",
		slog.String("env", cfg.AppEnv),
		slog.String("http_port", cfg.HTTPPort),
	)

	pool, err := db.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	slog.Info("pool connected")

	rdb, err := cache.NewClient(ctx, cfg.RedisAddr(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	slog.Info("redis connected")

	err = rdb.Set(ctx, "test", "123", time.Minute).Err()
	if err != nil {
		log.Fatal(err)
	}

	val, err := rdb.Get(ctx, "test").Result()
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("redis working", slog.String("value", val))

	// сюда потом:
	// db pool
	// redis client
	// router
	// server.Run()

}
