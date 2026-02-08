package dependency_injection

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http"
)

type Controllers struct {
	ClientController     *http.ClientController
	RequestLogController *http.RequestLogController
}

func SetupControllers(cfg *config.Config, usecases *Usecases) *Controllers {
	ClientController := http.NewClientController(cfg, usecases.ClientUsecase)
	RequestLogController := http.NewRequestLogController(cfg, usecases.RequestLogUsecase)

	return &Controllers{
		ClientController:     ClientController,
		RequestLogController: RequestLogController,
	}
}
