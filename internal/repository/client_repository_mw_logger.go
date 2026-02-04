package repository

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ ClientRepository = &ClientRepositoryMwLogger{}

type ClientRepositoryMwLogger struct {
	Next ClientRepository
}

func NewClientRepositoryMwLogger(next ClientRepository) *ClientRepositoryMwLogger {
	return &ClientRepositoryMwLogger{
		Next: next,
	}
}

func (c *ClientRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, client *entity.Client) error {
	err := c.Next.Create(ctx, db, client)

	fields := logrus.Fields{
		"client": client,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (c *ClientRepositoryMwLogger) FindByName(ctx context.Context, db *gorm.DB, client *entity.Client, name string) error {
	err := c.Next.FindByName(ctx, db, client, name)

	fields := logrus.Fields{
		"client": client,
	}
	x.LogMw(ctx, fields, err)

	return err
}
