package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/repository"
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
	cfg *config.Config
	db  *gorm.DB

	// repository
	clientRepository repository.ClientRepository
}

func NewClientUsecase(
	cfg *config.Config,
	db *gorm.DB,

	// repository
	clientRepository repository.ClientRepository,
) *ClientUsecaseImpl {
	return &ClientUsecaseImpl{
		cfg: cfg,
		db:  db,

		// repository
		clientRepository: clientRepository,
	}
}
