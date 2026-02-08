package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (c *ClientUsecaseImpl) Register(ctx context.Context, req *model.ReqRegister) (*model.ResRegister, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	client := new(entity.Client)
	converter.ModelReqRegisterToEntityClient(req, client)

	err = client.PrepareCredentialsForStorage()
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	err = c.clientRepository.Create(ctx, c.db, client)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.ResRegister)
	err = converter.EntityClientToModelResRegister(client, res)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return res, nil
}
