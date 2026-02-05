package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/layer"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ClientUsecase = &ClientUsecaseMwLogger{}

type ClientUsecaseMwLogger struct {
	Next ClientUsecase
}

func NewClientUsecaseMwLogger(next ClientUsecase) *ClientUsecaseMwLogger {
	return &ClientUsecaseMwLogger{
		Next: next,
	}
}

func (c *ClientUsecaseMwLogger) Register(ctx context.Context, req *model.ReqRegister) (*model.ResRegister, error) {
	res, err := c.Next.Register(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}

func (c *ClientUsecaseMwLogger) Login(ctx context.Context, req *model.ReqLogin) (*model.ResLogin, error) {
	res, err := c.Next.Login(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}

func (c *ClientUsecaseMwLogger) GetClientDetail(ctx context.Context, req *model.ReqGetClientDetail) (*model.ResGetClientDetail, error) {
	res, err := c.Next.GetClientDetail(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}

func (c *ClientUsecaseMwLogger) Verify(ctx context.Context, req *model.ReqVerify) (*model.ClientAuth, error) {
	res, err := c.Next.Verify(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}
