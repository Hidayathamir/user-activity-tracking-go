package cache

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/layer"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ Cache = &CacheMwLogger{}

type CacheMwLogger struct {
	Next Cache
}

func NewCacheMwLogger(next Cache) *CacheMwLogger {
	return &CacheMwLogger{
		Next: next,
	}
}

func (c *CacheMwLogger) SetClientRequestCountIfExist(ctx context.Context, apiKey string, datetime time.Time, value int) error {
	err := c.Next.SetClientRequestCountIfExist(ctx, apiKey, datetime, value)

	fields := logrus.Fields{
		"apiKey":   apiKey,
		"datetime": datetime,
		"value":    value,
	}
	x.LogMw(ctx, fields, err, layer.Cache)

	return err
}

func (c *CacheMwLogger) SetClientRequestCount(ctx context.Context, apiKey string, datetime time.Time, value int, ttl time.Duration) error {
	err := c.Next.SetClientRequestCount(ctx, apiKey, datetime, value, ttl)

	fields := logrus.Fields{
		"apiKey":   apiKey,
		"datetime": datetime,
		"value":    value,
		"ttl":      ttl,
	}
	x.LogMw(ctx, fields, err, layer.Cache)

	return err
}

func (c *CacheMwLogger) IncrementTopClientRequestCountHourly(ctx context.Context, timestamp time.Time, increment int, member string) error {
	err := c.Next.IncrementTopClientRequestCountHourly(ctx, timestamp, increment, member)

	fields := logrus.Fields{
		"timestamp": timestamp,
		"increment": increment,
		"member":    member,
	}
	x.LogMw(ctx, fields, err, layer.Cache)

	return err
}

func (c *CacheMwLogger) GetTop3ClientRequestCount24Hour(ctx context.Context) (model.APIKeyCountList, error) {
	res, err := c.Next.GetTop3ClientRequestCount24Hour(ctx)

	fields := logrus.Fields{
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Cache)

	return res, err
}

func (c *CacheMwLogger) GetCountByAPIKeyAndDate(ctx context.Context, apiKey string, datetime time.Time) (int, error) {
	res, err := c.Next.GetCountByAPIKeyAndDate(ctx, apiKey, datetime)

	fields := logrus.Fields{
		"apiKey":   apiKey,
		"datetime": datetime,
		"res":      res,
	}
	x.LogMw(ctx, fields, err, layer.Cache)

	return res, err
}
