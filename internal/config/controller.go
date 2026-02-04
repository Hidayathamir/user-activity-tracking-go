package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http"
	"github.com/spf13/viper"
)

type Controllers struct {
	ClientController *http.ClientController
}

func SetupControllers(viperConfig *viper.Viper, usecases *Usecases) *Controllers {
	ClientController := http.NewClientController(viperConfig, usecases.ClientUsecase)

	return &Controllers{
		ClientController: ClientController,
	}
}
