package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
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

func (c *ClientRequestCountRepositoryMwLogger) GetTop3ClientRequestCount24Hour(ctx context.Context, db *gorm.DB) (model.APIKeyCountList, error) {
	res, err := c.Next.GetTop3ClientRequestCount24Hour(ctx, db)

	fields := logrus.Fields{
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Repository)

	return res, err
}

func (c *ClientRequestCountRepositoryMwLogger) GetCountByAPIKeyAndDate(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time) (int, error) {
	res, err := c.Next.GetCountByAPIKeyAndDate(ctx, db, apiKey, datetime)

	fields := logrus.Fields{
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Repository)

	return res, err
}
