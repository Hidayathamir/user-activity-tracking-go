package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func Migrate(viperConfig *viper.Viper) {
	db := NewDatabase(viperConfig)

	sqlDB, err := db.DB()
	x.PanicIfErr(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	x.PanicIfErr(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+viperConfig.GetString(configkey.DatabaseMigrations),
		"postgres",
		driver,
	)
	x.PanicIfErr(err)

	err = m.Up()
	x.PanicIfErr(err)
}
