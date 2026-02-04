package main

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)

	config.Migrate(viperConfig)
}
