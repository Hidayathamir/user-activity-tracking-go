package dependency_injection

import (
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/cache"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/gateway/messaging"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/client"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/requestlog"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Usecases struct {
	ClientUsecase     client.ClientUsecase
	RequestLogUsecase requestlog.RequestLogUsecase
}

func SetupUsecases(
	cfg *config.Config,
	db *gorm.DB,
	rdb *redis.Client,
	kafkaWriter *kafka.Writer,
) *Usecases {
	const maxKey = 1000
	const ttl = 60 * time.Second
	inMemCache := expirable.NewLRU[string, int](maxKey, nil, ttl)

	// setup caches
	var Cache cache.Cache
	Cache = cache.NewCache(cfg, rdb, inMemCache)
	Cache = cache.NewCacheMwLogger(Cache)

	// setup producer
	var Producer messaging.Producer
	Producer = messaging.NewProducer(cfg, kafkaWriter)
	Producer = messaging.NewProducerMwLogger(Producer)

	// setup repositories
	var ClientRepository repository.ClientRepository
	ClientRepository = repository.NewClientRepository(cfg)
	ClientRepository = repository.NewClientRepositoryMwLogger(ClientRepository)

	var RequestLogRepository repository.RequestLogRepository
	RequestLogRepository = repository.NewRequestLogRepository(cfg)
	RequestLogRepository = repository.NewRequestLogRepositoryMwLogger(RequestLogRepository)

	var ClientRequestCountRepository repository.ClientRequestCountRepository
	ClientRequestCountRepository = repository.NewClientRequestCountRepository(cfg)
	ClientRequestCountRepository = repository.NewClientRequestCountRepositoryMwLogger(ClientRequestCountRepository)

	// setup usecases
	var ClientUsecase client.ClientUsecase
	ClientUsecase = client.NewClientUsecase(cfg, db, ClientRepository)
	ClientUsecase = client.NewClientUsecaseMwLogger(ClientUsecase)

	var RequestLogUsecase requestlog.RequestLogUsecase
	RequestLogUsecase = requestlog.NewRequestLogUsecase(cfg, db, Cache, RequestLogRepository, ClientRequestCountRepository, ClientRepository, Producer)
	RequestLogUsecase = requestlog.NewRequestLogUsecaseMwLogger(RequestLogUsecase)

	// returning

	return &Usecases{
		ClientUsecase:     ClientUsecase,
		RequestLogUsecase: RequestLogUsecase,
	}
}
