package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/cachekey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/timekit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)

	rdb := config.NewRedis(viperConfig)
	defer x.PanicIfErrForDefer(rdb.Close)

	myCron := cron.New()

	pattern := viperConfig.GetString(configkey.CronPattern)

	_, err := myCron.AddFunc(pattern, func() {
		x.Logger.Info("running unionTopClients job")
		err := unionTopClients(context.Background(), rdb)
		x.LogIfErr(err)
	})
	x.PanicIfErr(err)

	x.Logger.Info("running initial unionTopClients job")
	err = unionTopClients(context.Background(), rdb)
	x.LogIfErr(err)

	myCron.Start()
	x.Logger.WithField("pattern", pattern).Info("cron scheduler started, running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	x.Logger.Info("shutting down cron scheduler...")
	myCron.Stop()
	x.Logger.Info("cron scheduler stopped")
}

func unionTopClients(ctx context.Context, rdb *redis.Client) error {
	pipe := rdb.TxPipeline()

	last24HKeys := getLast24hKeys()

	pipe.ZUnionStore(ctx, cachekey.TopClientRequestCount24H, &redis.ZStore{
		Keys: last24HKeys,
	})

	pipe.Expire(ctx, cachekey.TopClientRequestCount24H, 1*time.Hour)

	pipe.Set(ctx, cachekey.TopClientRequestCount24HSourceKeys, strings.Join(last24HKeys, "\n"), 1*time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return err
}

func getLast24hKeys() []string {
	var keys []string
	now := time.Now().UTC()
	now = timekit.TruncateToHour(now)

	for i := range 24 {
		t := now.Add(time.Duration(-i) * time.Hour)

		key := cachekey.BuildTopClientRequestCountHourlyKey(t)

		keys = append(keys, key)
	}

	return keys
}
