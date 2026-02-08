package route

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/dependency_injection"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(ginEngine *gin.Engine, controllers *dependency_injection.Controllers, middlewares *dependency_injection.Middlewares) {
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ginEngine.POST("/api/register", controllers.ClientController.Register)
	ginEngine.POST("/api/login", controllers.ClientController.Login)
	ginEngine.GET("/api/client/me", middlewares.AuthMiddleware, controllers.ClientController.GetClientDetail)

	ginEngine.POST("/api/logs", middlewares.InternalServiceMiddleware, controllers.RequestLogController.RecordAPIHit)
	ginEngine.GET("/api/usage/daily", middlewares.AuthMiddleware, controllers.RequestLogController.GetClientDailyRequestCount)
	ginEngine.GET("/api/usage/top", middlewares.AuthMiddleware, controllers.RequestLogController.GetTop3ClientRequestCount24Hour)
}
