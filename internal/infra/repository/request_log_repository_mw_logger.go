package repository

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/layer"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ RequestLogRepository = &RequestLogRepositoryMwLogger{}

type RequestLogRepositoryMwLogger struct {
	Next RequestLogRepository
}

func NewRequestLogRepositoryMwLogger(next RequestLogRepository) *RequestLogRepositoryMwLogger {
	return &RequestLogRepositoryMwLogger{
		Next: next,
	}
}

func (r *RequestLogRepositoryMwLogger) CreateAll(ctx context.Context, db *gorm.DB, requestLogList *entity.RequestLogList) error {
	err := r.Next.CreateAll(ctx, db, requestLogList)

	fields := logrus.Fields{
		"requestLogList": requestLogList,
	}
	x.LogMw(ctx, fields, err, layer.Repository)

	return err
}
