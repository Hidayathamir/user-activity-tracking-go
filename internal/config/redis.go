package config

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.GetRedisAddress(),
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
	})
	return rdb
}
