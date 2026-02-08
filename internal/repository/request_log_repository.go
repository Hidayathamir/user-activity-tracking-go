package repository

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/dbretry"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryRequestLog.go -pkg=mock . RequestLogRepository

type RequestLogRepository interface {
	CreateAll(ctx context.Context, db *gorm.DB, requestLogList *entity.RequestLogList) error
}

var _ RequestLogRepository = &RequestLogRepositoryImpl{}

type RequestLogRepositoryImpl struct {
	cfg *config.Config
}

func NewRequestLogRepository(cfg *config.Config) *RequestLogRepositoryImpl {
	return &RequestLogRepositoryImpl{
		cfg: cfg,
	}
}

func (r *RequestLogRepositoryImpl) CreateAll(ctx context.Context, db *gorm.DB, requestLogList *entity.RequestLogList) error {
	err := dbretry.Do(func() error {
		return db.WithContext(ctx).Create(requestLogList).Error
	})
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
