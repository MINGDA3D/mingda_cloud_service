package cache

import (
    "context"
    "fmt"
    "mingda_cloud_service/internal/pkg/config"
    "github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })
    
    // 测试连接
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    
    return client, nil
}
