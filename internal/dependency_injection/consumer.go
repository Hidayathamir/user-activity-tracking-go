package dependency_injection

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/messaging"
)

type Consumers struct {
	RequestLogConsumer *messaging.RequestLogConsumer
}

func SetupConsumers(cfg *config.Config, usecases *Usecases) *Consumers {
	RequestLogConsumer := messaging.NewRequestLogConsumer(cfg, usecases.RequestLogUsecase)

	return &Consumers{
		RequestLogConsumer: RequestLogConsumer,
	}
}
