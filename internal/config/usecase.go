package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/cache"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/gateway/messaging"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/client"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/requestlog"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	ClientUsecase     client.ClientUsecase
	RequestLogUsecase requestlog.RequestLogUsecase
}

func SetupUsecases(
	viperConfig *viper.Viper,
	db *gorm.DB,
	rdb *redis.Client,
	kafkaWriter *kafka.Writer,
) *Usecases {
	// setup caches
	var Cache cache.Cache
	Cache = cache.NewCache(viperConfig, rdb)
	Cache = cache.NewCacheMwLogger(Cache)

	// setup producer
	var Producer messaging.Producer
	Producer = messaging.NewProducer(viperConfig, kafkaWriter)
	Producer = messaging.NewProducerMwLogger(Producer)

	// setup repositories
	var ClientRepository repository.ClientRepository
	ClientRepository = repository.NewClientRepository(viperConfig)
	ClientRepository = repository.NewClientRepositoryMwLogger(ClientRepository)

	var RequestLogRepository repository.RequestLogRepository
	RequestLogRepository = repository.NewRequestLogRepository(viperConfig)
	RequestLogRepository = repository.NewRequestLogRepositoryMwLogger(RequestLogRepository)

	var ClientRequestCountRepository repository.ClientRequestCountRepository
	ClientRequestCountRepository = repository.NewClientRequestCountRepository(viperConfig)
	ClientRequestCountRepository = repository.NewClientRequestCountRepositoryMwLogger(ClientRequestCountRepository)

	// setup usecases
	var ClientUsecase client.ClientUsecase
	ClientUsecase = client.NewClientUsecase(viperConfig, db, ClientRepository)
	ClientUsecase = client.NewClientUsecaseMwLogger(ClientUsecase)

	var RequestLogUsecase requestlog.RequestLogUsecase
	RequestLogUsecase = requestlog.NewRequestLogUsecase(viperConfig, db, Cache, RequestLogRepository, ClientRequestCountRepository, Producer)
	RequestLogUsecase = requestlog.NewRequestLogUsecaseMwLogger(RequestLogUsecase)

	// returning

	return &Usecases{
		ClientUsecase:     ClientUsecase,
		RequestLogUsecase: RequestLogUsecase,
	}
}
