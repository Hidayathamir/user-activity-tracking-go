package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/client"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	ClientUsecase client.ClientUsecase
}

func SetupUsecases(
	viperConfig *viper.Viper,
	db *gorm.DB,
) *Usecases {
	// setup repositories
	var ClientRepository repository.ClientRepository
	ClientRepository = repository.NewClientRepository(viperConfig)
	ClientRepository = repository.NewClientRepositoryMwLogger(ClientRepository)

	// setup use cases
	var ClientUsecase client.ClientUsecase
	ClientUsecase = client.NewClientUsecase(viperConfig, db, ClientRepository)
	ClientUsecase = client.NewClientUsecaseMwLogger(ClientUsecase)

	return &Usecases{
		ClientUsecase: ClientUsecase,
	}
}
