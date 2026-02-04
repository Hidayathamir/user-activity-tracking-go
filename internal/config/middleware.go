package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	AuthMiddleware            gin.HandlerFunc
	InternalServiceMiddleware gin.HandlerFunc
}

func SetupMiddlewares(usecases *Usecases) *Middlewares {
	AuthMiddleware := middleware.NewAuthMiddleware(usecases.ClientUsecase)
	InternalServiceMiddleware := middleware.NewInternalServiceMiddleware()

	return &Middlewares{
		AuthMiddleware:            AuthMiddleware,
		InternalServiceMiddleware: InternalServiceMiddleware,
	}
}
