package redis

import (
	"context"
	"fmt"
	"time"

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

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Set 设置值
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
} 