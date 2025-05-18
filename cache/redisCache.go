package cache

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/redis/go-redis/v9"
	"shortenLink/config"
)

var ctx context.Context

// 多级缓存设计
type RedisCache struct {
	Redis *redis.Client
	Local *ristretto.Cache
}

var RedisCacheIns *RedisCache

func NewRedisCache(cfg config.RedisConfig) *RedisCache {
	local, _ := ristretto.NewCache(&ristretto.Config{ //具体设置待学习和商议
		NumCounters: 1e7,     // 10M
		MaxCost:     1 << 30, // 最大内存1GB
		BufferItems: 64,
	})
	rd := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
	return &RedisCache{
		Redis: rd,
		Local: local,
	}
}

func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
