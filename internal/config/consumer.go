package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/messaging"
	"github.com/spf13/viper"
)

type Consumers struct {
	RequestLogConsumer *messaging.RequestLogConsumer
}

func SetupConsumers(viperConfig *viper.Viper, usecases *Usecases) *Consumers {
	RequestLogConsumer := messaging.NewRequestLogConsumer(viperConfig, usecases.RequestLogUsecase)

	return &Consumers{
		RequestLogConsumer: RequestLogConsumer,
	}
}
