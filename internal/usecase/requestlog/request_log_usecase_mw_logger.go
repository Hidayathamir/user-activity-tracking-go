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

func (r *RequestLogUsecaseMwLogger) GetClientDailyRequestCount(ctx context.Context, req *model.ReqGetClientDailyRequestCount) (*model.ResGetClientDailyRequestCount, error) {
	res, err := r.Next.GetClientDailyRequestCount(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}

func (r *RequestLogUsecaseMwLogger) GetTop3ClientRequestCount24Hour(ctx context.Context, req *model.ReqGetTop3ClientRequestCount24Hour) (*model.ResGetTop3ClientRequestCount24Hour, error) {
	res, err := r.Next.GetTop3ClientRequestCount24Hour(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err, layer.Usecase)

	return res, err
}
