package requestlog

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/cache"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/gateway/messaging"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseRequestLog.go -pkg=mock . RequestLogUsecase

type RequestLogUsecase interface {
	RecordAPIHit(ctx context.Context, req *model.ReqRecordAPIHit) (*model.ResRecordAPIHit, error)
	BatchConsumeClientRequestLogEvent(ctx context.Context, req *model.ReqBatchConsumeClientRequestLogEvent) error
	GetClientDailyRequestCount(ctx context.Context, req *model.ReqGetClientDailyRequestCount) (*model.ResGetClientDailyRequestCount, error)
	GetTop3ClientRequestCount24Hour(ctx context.Context, req *model.ReqGetTop3ClientRequestCount24Hour) (*model.ResGetTop3ClientRequestCount24Hour, error)
}

var _ RequestLogUsecase = &RequestLogUsecaseImpl{}

type RequestLogUsecaseImpl struct {
	Config *viper.Viper
	DB     *gorm.DB

	// cache
	Cache cache.Cache

	// repository
	RequestLogRepository         repository.RequestLogRepository
	ClientRequestCountRepository repository.ClientRequestCountRepository
	ClientRepository             repository.ClientRepository

	// Producer
	Producer messaging.Producer
}

func NewRequestLogUsecase(
	Config *viper.Viper,
	DB *gorm.DB,

	// cache
	Cache cache.Cache,

	// repository
	RequestLogRepository repository.RequestLogRepository,
	ClientRequestCountRepository repository.ClientRequestCountRepository,
	ClientRepository repository.ClientRepository,

	// Producer
	Producer messaging.Producer,
) *RequestLogUsecaseImpl {
	return &RequestLogUsecaseImpl{
		Config: Config,
		DB:     DB,

		// cache
		Cache: Cache,

		// repository
		RequestLogRepository:         RequestLogRepository,
		ClientRequestCountRepository: ClientRequestCountRepository,
		ClientRepository:             ClientRepository,

		// Producer
		Producer: Producer,
	}
}
