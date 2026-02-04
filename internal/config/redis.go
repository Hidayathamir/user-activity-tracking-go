package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viperConfig *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: viperConfig.GetString(configkey.RedisAddress),
	})
	return rdb
}
