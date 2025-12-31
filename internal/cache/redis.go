package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/randhir/aegis-core/internal/config"
	"github.com/randhir/aegis-core/internal/logger"
)

var Client *redis.Client

func ConnectRedis() error {
	cfg := config.AppConfig.Redis

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("addr", cfg.Addr),
	)

	return nil
}

func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

