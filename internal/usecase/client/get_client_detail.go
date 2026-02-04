package client

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (c *ClientUsecaseImpl) GetClientDetail(ctx context.Context, req *model.ReqGetClientDetail) (*model.ResGetClientDetail, error) {
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

	res := new(model.ResGetClientDetail)
	err = converter.EntityClientToModelResCurrent(client, res)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return res, nil
}
