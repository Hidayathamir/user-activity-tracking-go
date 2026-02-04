package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http"
	"github.com/spf13/viper"
)

type Controllers struct {
	ClientController     *http.ClientController
	RequestLogController *http.RequestLogController
}

func SetupControllers(viperConfig *viper.Viper, usecases *Usecases) *Controllers {
	ClientController := http.NewClientController(viperConfig, usecases.ClientUsecase)
	RequestLogController := http.NewRequestLogController(viperConfig, usecases.RequestLogUsecase)

	return &Controllers{
		ClientController:     ClientController,
		RequestLogController: RequestLogController,
	}
}
