package redis

import (
	"context"
	"time"
)

// Lock 获取分布式锁
func Lock(ctx context.Context, key string, expiration time.Duration) bool {
	return Client.SetNX(ctx, "lock:"+key, 1, expiration).Val()
}

// Unlock 释放分布式锁
func Unlock(ctx context.Context, key string) {
	Client.Del(ctx, "lock:"+key)
} 