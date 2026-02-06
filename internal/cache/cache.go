package cache

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/cachekey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:generate moq -out=../mock/MockCache.go -pkg=mock . Cache

// TODO: if redis down, use in memory cache

type Cache interface {
	SetClientRequestCountIfExist(ctx context.Context, apiKey string, datetime time.Time, value int) error
	SetClientRequestCount(ctx context.Context, apiKey string, datetime time.Time, value int, ttl time.Duration) error
	IncrementTopClientRequestCountHourly(ctx context.Context, timestamp time.Time, increment int, member string) error
	GetTop3ClientRequestCount24Hour(ctx context.Context) (model.APIKeyCountList, error)
	GetCountByAPIKeyAndDate(ctx context.Context, apiKey string, datetime time.Time) (int, error)
}

var _ Cache = &CacheImpl{}

type CacheImpl struct {
	Config     *viper.Viper
	rdb        *redis.Client
	inMemCache *expirable.LRU[string, int]
}

func NewCache(cfg *viper.Viper, rdb *redis.Client, inMemCache *expirable.LRU[string, int]) *CacheImpl {
	return &CacheImpl{
		Config:     cfg,
		rdb:        rdb,
		inMemCache: inMemCache,
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
		if isRedisUnavailable(err) {
			x.Logger.WithError(err).WithContext(ctx).Warn("redis unavailable, using in memory cache")
			c.inMemCache.Add(key, value)
			return nil
		}
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

func (c *CacheImpl) SetClientRequestCount(ctx context.Context, apiKey string, datetime time.Time, value int, ttl time.Duration) error {
	key := cachekey.BuildClientRequestCountKey(apiKey, datetime)

	err := c.rdb.Set(ctx, string(key), value, ttl).Err()
	if err != nil {
		if isRedisUnavailable(err) {
			x.Logger.WithError(err).WithContext(ctx).Warn("redis unavailable, using in memory cache")
			c.inMemCache.Add(key, value)
			return nil
		}
		return errkit.AddFuncName(err)
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

func (c *CacheImpl) GetTop3ClientRequestCount24Hour(ctx context.Context) (model.APIKeyCountList, error) {
	result, err := c.rdb.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   cachekey.TopClientRequestCount24H,
		Start: 0,
		Stop:  2,
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := model.APIKeyCountList{}
	for _, z := range result {
		apiKeyCount := new(model.APIKeyCount)
		converter.RedisZToModelAPIKeyCount(&z, apiKeyCount)

		res = append(res, *apiKeyCount)
	}

	return res, nil
}

func (c *CacheImpl) GetCountByAPIKeyAndDate(ctx context.Context, apiKey string, datetime time.Time) (int, error) {
	key := cachekey.BuildClientRequestCountKey(apiKey, datetime)

	val, err := c.rdb.Get(ctx, key).Int()
	if err != nil {
		if isRedisUnavailable(err) {
			x.Logger.WithError(err).WithContext(ctx).Warn("redis unavailable, using in memory cache")
			localVal, ok := c.inMemCache.Get(key)
			if !ok {
				err := redis.Nil
				return 0, errkit.AddFuncName(err)
			}
			return localVal, nil
		}
		return 0, errkit.AddFuncName(err)
	}
	return val, nil
}
