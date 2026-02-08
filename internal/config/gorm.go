package config

import (
	"context"
	"fmt"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *Config) *gorm.DB {
	username := cfg.GetDatabaseUsername()
	password := cfg.GetDatabasePassword()
	host := cfg.GetDatabaseHost()
	port := cfg.GetDatabasePort()
	database := cfg.GetDatabaseName()
	idleConnection := cfg.GetDatabasePoolIdle()
	maxConnection := cfg.GetDatabasePoolMax()
	maxLifeTimeConnection := cfg.GetDatabasePoolLifetime()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC", host, port, username, password, database)

	const maxAttempts = 5

	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         &gormLogger{},
			TranslateError: true,
		})
		if err == nil {
			break
		}
		x.Logger.Warnf("database connection attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(1 * time.Second)
	}
	x.PanicIfErr(err)

	connection, err := db.DB()
	x.PanicIfErr(err)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if pingErr := connection.Ping(); pingErr == nil {
			break
		} else if attempt == maxAttempts {
			x.Logger.Panicf("failed to connect database: %v", pingErr)
		} else {
			x.Logger.Warnf("database ping attempt %d/%d failed: %v", attempt, maxAttempts, pingErr)
			time.Sleep(1 * time.Second)
		}
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type gormLogger struct{}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *gormLogger) Info(ctx context.Context, msg string, data ...any) {
	x.Logger.WithContext(ctx).Infof(msg, data...)
}

func (g *gormLogger) Warn(ctx context.Context, msg string, data ...any) {
	x.Logger.WithContext(ctx).Warnf(msg, data...)
}

func (g *gormLogger) Error(ctx context.Context, msg string, data ...any) {
	x.Logger.WithContext(ctx).Errorf(msg, data...)
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	entry := x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"elapsed": time.Since(begin),
		"rows":    rows,
		"sql":     sql,
	})

	if err != nil {
		entry.WithError(err).Error()
	} else {
		// with error nil so that in log we know is the query success or not
		entry.WithError(nil).Info()
	}
}
