package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/column"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryClient.go -pkg=mock . ClientRepository

type ClientRepository interface {
	Create(ctx context.Context, db *gorm.DB, client *entity.Client) error
	FindByName(ctx context.Context, db *gorm.DB, client *entity.Client, name string) error
	FindByAPIKey(ctx context.Context, db *gorm.DB, client *entity.Client, apiKey string) error
}

var _ ClientRepository = &ClientRepositoryImpl{}

type ClientRepositoryImpl struct {
	Config *viper.Viper
}

func NewClientRepository(cfg *viper.Viper) *ClientRepositoryImpl {
	return &ClientRepositoryImpl{
		Config: cfg,
	}
}

func (c *ClientRepositoryImpl) Create(ctx context.Context, db *gorm.DB, client *entity.Client) error {
	err := db.WithContext(ctx).Create(client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = errkit.SetHTTPError(err, http.StatusConflict)
		}
		return errkit.AddFuncName(err)
	}
	return nil
}

func (c *ClientRepositoryImpl) FindByName(ctx context.Context, db *gorm.DB, client *entity.Client, name string) error {
	err := db.WithContext(ctx).Where(column.Name.Eq(name)).Take(client).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}

func (c *ClientRepositoryImpl) FindByAPIKey(ctx context.Context, db *gorm.DB, client *entity.Client, apiKey string) error {
	err := db.WithContext(ctx).Where(column.APIKey.Eq(apiKey)).Take(client).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}
	return nil
}
