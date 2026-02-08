package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (c *ClientUsecaseImpl) Verify(ctx context.Context, req *model.ReqVerify) (*model.ClientAuth, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	clientName, err := c.parseAccessToken(ctx, req.Token)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	client := new(entity.Client)
	err = c.clientRepository.FindByName(ctx, c.db, client, clientName)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	clientAuth := new(model.ClientAuth)
	err = converter.EntityClientToModelClientAuth(client, clientAuth)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return clientAuth, nil
}
