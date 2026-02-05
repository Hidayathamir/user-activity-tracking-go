package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/layer"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ ClientRequestCountRepository = &ClientRequestCountRepositoryMwLogger{}

type ClientRequestCountRepositoryMwLogger struct {
	Next ClientRequestCountRepository
}

func NewClientRequestCountRepositoryMwLogger(next ClientRequestCountRepository) *ClientRequestCountRepositoryMwLogger {
	return &ClientRequestCountRepositoryMwLogger{
		Next: next,
	}
}

func (c *ClientRequestCountRepositoryMwLogger) IncrementCount(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time, count int) (int, error) {
	newCount, err := c.Next.IncrementCount(ctx, db, apiKey, datetime, count)

	fields := logrus.Fields{
		"apiKey":   apiKey,
		"datetime": datetime,
		"count":    count,
		"newCount": newCount,
	}
	x.LogMw(ctx, fields, err, layer.Repository)

	return newCount, err
}
