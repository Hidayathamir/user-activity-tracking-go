package repository

import "github.com/spf13/viper"

//go:generate moq -out=../mock/MockRepositoryClientRequestCount.go -pkg=mock . ClientRequestCountRepository

type ClientRequestCountRepository interface {
}

var _ ClientRequestCountRepository = &ClientRequestCountRepositoryImpl{}

type ClientRequestCountRepositoryImpl struct {
	Config *viper.Viper
}

func NewClientRequestCountRepository(cfg *viper.Viper) *ClientRequestCountRepositoryImpl {
	return &ClientRequestCountRepositoryImpl{
		Config: cfg,
	}
}
