package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../../")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	x.PanicIfErr(err)

	return config
}
