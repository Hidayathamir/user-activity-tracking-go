package cache

import (
	"context"
	"time"

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
