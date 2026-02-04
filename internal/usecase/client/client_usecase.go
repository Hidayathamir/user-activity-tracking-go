package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseClient.go -pkg=mock . ClientUsecase

type ClientUsecase interface {
	Register(ctx context.Context, req *model.ReqRegister) (*model.ResRegister, error)
	Login(ctx context.Context, req *model.ReqLogin) (*model.ResLogin, error)
	GetClientDetail(ctx context.Context, req *model.ReqGetClientDetail) (*model.ResGetClientDetail, error)
	Verify(ctx context.Context, req *model.ReqVerify) (*model.ClientAuth, error)
}

var _ ClientUsecase = &ClientUsecaseImpl{}

type ClientUsecaseImpl struct {
	Config *viper.Viper
	DB     *gorm.DB

	// repository
	ClientRepository repository.ClientRepository
}

func NewClientUsecase(
	Config *viper.Viper,
	DB *gorm.DB,

	// repository
	ClientRepository repository.ClientRepository,
) *ClientUsecaseImpl {
	return &ClientUsecaseImpl{
		Config: Config,
		DB:     DB,

		// repository
		ClientRepository: ClientRepository,
	}
}
