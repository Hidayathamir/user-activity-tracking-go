package main

import (
	_ "github.com/Hidayathamir/user-activity-tracking-go/api" // need import for swagger
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/route"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

// General API Info
// for swag (https://github.com/swaggo/swag)

//	@title	User Activity Tracking Go

//	@securityDefinitions.apikey	ApiKeyJWTAuth
//	@in							header
//	@name						Authorization
//	@description				JWT authorization

//	@securityDefinitions.apikey	ApiKeyXInternalSecret
//	@in							header
//	@name						X-Internal-Secret
//	@description				X Internal Secret

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)

	db := config.NewDatabase(viperConfig)

	rdb := config.NewRedis(viperConfig)
	defer x.PanicIfErrForDefer(rdb.Close)

	kafkaWriter := config.NewKafkaProducer(viperConfig)
	defer x.PanicIfErrForDefer(kafkaWriter.Close)

	usecases := config.SetupUsecases(viperConfig, db, rdb, kafkaWriter)

	controllers := config.SetupControllers(viperConfig, usecases)

	middlewares := config.SetupMiddlewares(usecases)

	ginEngine := config.NewGinEngine(viperConfig)

	route.Setup(ginEngine, controllers, middlewares)

	err := ginEngine.Run("0.0.0.0:9090")
	x.PanicIfErr(err)
}
