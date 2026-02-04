package repository

import "github.com/spf13/viper"

//go:generate moq -out=../mock/MockRepositoryRequestLog.go -pkg=mock . RequestLogRepository

type RequestLogRepository interface {
}

var _ RequestLogRepository = &RequestLogRepositoryImpl{}

type RequestLogRepositoryImpl struct {
	Config *viper.Viper
}

func NewRequestLogRepository(cfg *viper.Viper) *RequestLogRepositoryImpl {
	return &RequestLogRepositoryImpl{
		Config: cfg,
	}
}
