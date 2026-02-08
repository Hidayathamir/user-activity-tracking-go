package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *Config) {
	db := NewDatabase(cfg)

	sqlDB, err := db.DB()
	x.PanicIfErr(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	x.PanicIfErr(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+cfg.GetDatabaseMigrations(),
		"postgres",
		driver,
	)
	x.PanicIfErr(err)

	err = m.Up()
	x.PanicIfErr(err)
}
