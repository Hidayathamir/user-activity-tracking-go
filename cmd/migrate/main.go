package main

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func main() {
	cfg := config.NewConfig()
	x.SetupAll(cfg.GetLogLevel(), cfg.GetAESKey())

	config.Migrate(cfg)
}
