package config

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./../../")
	v.AddConfigPath("./../")
	v.AddConfigPath("./")

	config := &Config{Viper: v}

	err := config.ReadConfig()

	x.PanicIfErr(err)

	return config
}

func (c *Config) ReadConfig() error {
	return c.ReadInConfig()
}

// --- AES ---

func (c *Config) GetAESKey() string {
	return c.GetString(configkey.AESKey)
}

// --- App ---

func (c *Config) GetAppName() string {
	return c.GetString(configkey.AppName)
}

// --- Auth JWT ---

func (c *Config) GetAuthJWTSecret() string {
	return c.GetString(configkey.AuthJWTSecret)
}

func (c *Config) GetAuthJWTIssuer() string {
	return c.GetString(configkey.AuthJWTIssuer)
}

func (c *Config) GetAuthJWTExpireSeconds() int {
	return c.GetInt(configkey.AuthJWTExpireSeconds)
}

// --- Cron ---

func (c *Config) GetCronPattern() string {
	return c.GetString(configkey.CronPattern)
}

// --- Database ---

func (c *Config) GetDatabaseMigrations() string {
	return c.GetString(configkey.DatabaseMigrations)
}

func (c *Config) GetDatabaseUsername() string {
	return c.GetString(configkey.DatabaseUsername)
}

func (c *Config) GetDatabasePassword() string {
	return c.GetString(configkey.DatabasePassword)
}

func (c *Config) GetDatabaseHost() string {
	return c.GetString(configkey.DatabaseHost)
}

func (c *Config) GetDatabasePort() int {
	return c.GetInt(configkey.DatabasePort)
}

func (c *Config) GetDatabaseName() string {
	return c.GetString(configkey.DatabaseName)
}

func (c *Config) GetDatabasePoolIdle() int {
	return c.GetInt(configkey.DatabasePoolIdle)
}

func (c *Config) GetDatabasePoolMax() int {
	return c.GetInt(configkey.DatabasePoolMax)
}

func (c *Config) GetDatabasePoolLifetime() int {
	return c.GetInt(configkey.DatabasePoolLifetime)
}

// --- Kafka ---

func (c *Config) GetKafkaAddress() string {
	return c.GetString(configkey.KafkaAddress)
}

func (c *Config) GetKafkaGroupID() string {
	return c.GetString(configkey.KafkaGroupID)
}

// --- Log ---

func (c *Config) GetLogLevel() string {
	return c.GetString(configkey.LogLevel)
}

// --- Redis ---

func (c *Config) GetRedisAddress() string {
	return c.GetString(configkey.RedisAddress)
}
