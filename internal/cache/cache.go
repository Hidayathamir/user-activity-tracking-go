package cache

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/cachekey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:generate moq -out=../mock/MockCache.go -pkg=mock . Cache

type Cache interface {
	SetClientRequestCountIfExist(ctx context.Context, apiKey string, datetime time.Time, value int) error
	IncrementTopClientRequestCountHourly(ctx context.Context, timestamp time.Time, increment int, member string) error
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

func (c *CacheImpl) SetClientRequestCountIfExist(ctx context.Context, apiKey string, datetime time.Time, value int) error {
	key := cachekey.BuildClientRequestCountKey(apiKey, datetime)

	const setOnlyIfKeyExist = "XX"
	ok, err := c.rdb.SetArgs(ctx, string(key), value, redis.SetArgs{
		Mode:    setOnlyIfKeyExist,
		KeepTTL: true,
	}).Result()
	if err != nil {
		return errkit.AddFuncName(err)
	}

	if ok == "" {
		x.Logger.WithContext(ctx).WithFields(logrus.Fields{
			"key":   key,
			"value": value,
		}).Debug("key did not exist")
	}

	return nil
}

func (c *CacheImpl) IncrementTopClientRequestCountHourly(ctx context.Context, timestamp time.Time, increment int, member string) error {
	key := cachekey.BuildTopClientRequestCountHourlyKey(timestamp)

	_, err := incrWithTTL.Run(ctx, c.rdb, []string{key}, increment, member).Result()
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
