package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

//go:generate moq -out=../mock/MockCache.go -pkg=mock . Cache

type Cache interface {
}

var _ Cache = &CacheImpl{}

type CacheImpl struct {
	Config *viper.Viper
	rdb    *redis.Client
}

func NewCache(cfg *viper.Viper, rdb *redis.Client) *CacheImpl {
	return &CacheImpl{
		Config: cfg,
		rdb:    rdb,
	}
}
