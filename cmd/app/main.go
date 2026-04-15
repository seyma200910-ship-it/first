package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"first/internal/cache"
	"first/internal/config"
	"first/internal/db"
	"first/internal/pkg"
	"first/internal/server"
	"first/internal/users"
)

func main() {
	ctx := context.Background()

	// config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// logger
	logger := pkg.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	slog.Info("config loaded",
		slog.String("env", cfg.AppEnv),
		slog.String("http_port", cfg.HTTPPort),
	)

	// postgres
	pool, err := db.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer pool.Close()

	slog.Info("pool connected")

	// redis
	redisClient, err := cache.NewClient(
		ctx,
		cfg.RedisAddr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("redis connection error: %v", err)
	}
	defer redisClient.Close()

	slog.Info("redis connected")

	// users module
	userRepo := users.NewPostgresRepository(pool)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	// router
	router := server.NewRouter(server.Handlers{
		UserHandler: userHandler,
	})

	// http server
	srv := server.New(router, cfg.HTTPPort)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("server run error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	slog.Info("application stopped")
}
