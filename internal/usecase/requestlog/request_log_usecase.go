package requestlog

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/infra/cache"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/infra/messaging"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/infra/repository"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
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
	cfg *config.Config
	db  *gorm.DB

	// cache
	cache cache.Cache

	// repository
	requestLogRepository         repository.RequestLogRepository
	clientRequestCountRepository repository.ClientRequestCountRepository
	clientRepository             repository.ClientRepository

	// producer
	producer messaging.Producer
}

func NewRequestLogUsecase(
	cfg *config.Config,
	db *gorm.DB,

	// cache
	cache cache.Cache,

	// repository
	requestLogRepository repository.RequestLogRepository,
	clientRequestCountRepository repository.ClientRequestCountRepository,
	clientRepository repository.ClientRepository,

	// Producer
	producer messaging.Producer,
) *RequestLogUsecaseImpl {
	return &RequestLogUsecaseImpl{
		cfg: cfg,
		db:  db,

		// cache
		cache: cache,

		// repository
		requestLogRepository:         requestLogRepository,
		clientRequestCountRepository: clientRequestCountRepository,
		clientRepository:             clientRepository,

		// Producer
		producer: producer,
	}
}
