package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10, // 🔥 аналог MaxConns
		MinIdleConns: 2,  // 🔥 аналог MinConns
	})

	// 🔁 проверка подключения
	if err := pingWithRetry(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

func pingWithRetry(ctx context.Context, client *redis.Client) error {
	var err error

	for range 10 {
		err = client.Ping(ctx).Err()
		if err == nil {
			return nil
		}

		time.Sleep(time.Second)
	}

	return fmt.Errorf("redis not available: %w", err)
}
