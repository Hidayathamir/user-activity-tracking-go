package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/cache"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/client"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	ClientUsecase client.ClientUsecase
}

func SetupUsecases(
	viperConfig *viper.Viper,
	db *gorm.DB,
	rdb *redis.Client,
) *Usecases {
	// setup caches
	var Cache cache.Cache
	Cache = cache.NewCache(viperConfig, rdb)
	Cache = cache.NewCacheMwLogger(Cache)
	_ = Cache

	// setup repositories
	var ClientRepository repository.ClientRepository
	ClientRepository = repository.NewClientRepository(viperConfig)
	ClientRepository = repository.NewClientRepositoryMwLogger(ClientRepository)

	// setup usecases
	var ClientUsecase client.ClientUsecase
	ClientUsecase = client.NewClientUsecase(viperConfig, db, ClientRepository)
	ClientUsecase = client.NewClientUsecaseMwLogger(ClientUsecase)

	return &Usecases{
		ClientUsecase: ClientUsecase,
	}
}
