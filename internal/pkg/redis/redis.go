package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"mingda_cloud_service/internal/pkg/config"
)

var Client *redis.Client

// Init 初始化Redis连接
func Init(cfg config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("connect redis failed: %v", err)
	}

	return nil
} 