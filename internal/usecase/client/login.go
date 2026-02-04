package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (c *ClientUsecaseImpl) Login(ctx context.Context, req *model.ReqLogin) (*model.ResLogin, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	client := new(entity.Client)
	err = c.ClientRepository.FindByName(ctx, c.DB, client, req.Name)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	err = client.ValidatePassword(req.Password)
	if err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	token, err := c.signAccessToken(ctx, client.Name)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.ResLogin)
	res.Token = token

	return res, nil
}
