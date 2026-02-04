package x

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/spf13/viper"
)

func SetupAll(viperConfig *viper.Viper) {
	SetupLogger(viperConfig)
	SetupValidator(viperConfig)
	AESKey = viperConfig.GetString(configkey.AESKey)
}
