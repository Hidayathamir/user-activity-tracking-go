package main

import (
	_ "github.com/Hidayathamir/user-activity-tracking-go/api" // need import for swagger
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/route"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/dependency_injection"
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
	cfg := config.NewConfig()
	x.SetupAll(cfg.GetLogLevel(), cfg.GetAESKey())

	db := config.NewDatabase(cfg)

	rdb := config.NewRedis(cfg)
	defer x.PanicIfErrForDefer(rdb.Close)

	kafkaWriter := config.NewKafkaWriter(cfg)
	defer x.PanicIfErrForDefer(kafkaWriter.Close)

	usecases := dependency_injection.SetupUsecases(cfg, db, rdb, kafkaWriter)

	controllers := dependency_injection.SetupControllers(cfg, usecases)

	middlewares := dependency_injection.SetupMiddlewares(usecases)

	ginEngine := config.NewGinEngine(cfg)

	route.Setup(ginEngine, controllers, middlewares)

	x.Logger.Debug("open swagger: http://localhost:9090/swagger/index.html")
	err := ginEngine.Run("0.0.0.0:9090")
	x.PanicIfErr(err)
}
