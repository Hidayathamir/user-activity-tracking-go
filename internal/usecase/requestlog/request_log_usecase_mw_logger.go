package requestlog

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/layer"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ RequestLogUsecase = &RequestLogUsecaseMwLogger{}

type RequestLogUsecaseMwLogger struct {
	Next RequestLogUsecase
}

func NewRequestLogUsecaseMwLogger(next RequestLogUsecase) *RequestLogUsecaseMwLogger {
	return &RequestLogUsecaseMwLogger{
		Next: next,
	}
}

func (r *RequestLogUsecaseMwLogger) RecordAPIHit(ctx context.Context, req *model.ReqRecordAPIHit) (*model.ResRecordAPIHit, error) {
	res, err := r.Next.RecordAPIHit(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}

func (r *RequestLogUsecaseMwLogger) BatchConsumeClientRequestLogEvent(ctx context.Context, req *model.ReqBatchConsumeClientRequestLogEvent) error {
	err := r.Next.BatchConsumeClientRequestLogEvent(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return err
}
