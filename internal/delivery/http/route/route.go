package route

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(ginEngine *gin.Engine, controllers *config.Controllers, middlewares *config.Middlewares) {
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ginEngine.POST("/api/register", controllers.ClientController.Register)
	ginEngine.POST("/api/login", controllers.ClientController.Login)
	ginEngine.GET("/api/client/me", middlewares.AuthMiddleware, controllers.ClientController.GetClientDetail)
}
