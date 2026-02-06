package requestlog

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (r *RequestLogUsecaseImpl) GetTop3ClientRequestCount24Hour(ctx context.Context, req *model.ReqGetTop3ClientRequestCount24Hour) (*model.ResGetTop3ClientRequestCount24Hour, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	apiKeyCountList, err := r.Cache.GetTop3ClientRequestCount24Hour(ctx)
	x.LogIfErrContext(ctx, err)

	if len(apiKeyCountList) == 0 {
		apiKeyCountList, err = r.ClientRequestCountRepository.GetTop3ClientRequestCount24Hour(ctx, r.DB)
		if err != nil {
			return nil, errkit.AddFuncName(err)
		}
	}

	nameCountList := model.NameCountList{}

	for _, v := range apiKeyCountList {
		client := new(entity.Client)
		err := r.ClientRepository.FindByAPIKey(ctx, r.DB, client, v.APIKey)
		if err != nil {
			return nil, errkit.AddFuncName(err)
		}

		nameCountList = append(nameCountList, model.NameCount{
			Name:  client.Name,
			Count: v.Count,
		})
	}

	res := new(model.ResGetTop3ClientRequestCount24Hour)
	res.NameCountList = nameCountList

	return res, nil
}
