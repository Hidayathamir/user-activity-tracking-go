package repository

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryRequestLog.go -pkg=mock . RequestLogRepository

type RequestLogRepository interface {
	CreateAll(ctx context.Context, db *gorm.DB, requestLogList *entity.RequestLogList) error
}

var _ RequestLogRepository = &RequestLogRepositoryImpl{}

type RequestLogRepositoryImpl struct {
	Config *viper.Viper
}

func NewRequestLogRepository(cfg *viper.Viper) *RequestLogRepositoryImpl {
	return &RequestLogRepositoryImpl{
		Config: cfg,
	}
}

func (r *RequestLogRepositoryImpl) CreateAll(ctx context.Context, db *gorm.DB, requestLogList *entity.RequestLogList) error {
	err := db.WithContext(ctx).Create(requestLogList).Error
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
