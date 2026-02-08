package requestlog

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (r *RequestLogUsecaseImpl) RecordAPIHit(ctx context.Context, req *model.ReqRecordAPIHit) (*model.ResRecordAPIHit, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	event := new(model.ClientRequestLogEvent)
	converter.ModelReqRecordAPIHitToModelClientRequestLogEvent(req, event)

	err = r.producer.SendClientRequestLogEvent(ctx, event)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	res := &model.ResRecordAPIHit{OK: true}

	return res, nil
}
