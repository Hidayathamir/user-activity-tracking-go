package requestlog

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/counter"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/timekit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func (r *RequestLogUsecaseImpl) BatchConsumeClientRequestLogEvent(ctx context.Context, req *model.ReqBatchConsumeClientRequestLogEvent) error {
	requestLogList := new(entity.RequestLogList)
	converter.ModelClientRequestLogEventListToEntityRequestLogList(&req.EventList, requestLogList)

	err := r.requestLogRepository.CreateAll(ctx, r.db, requestLogList)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	r.incrementClientRequestCount(ctx, req)

	return nil
}

func (r *RequestLogUsecaseImpl) incrementClientRequestCount(ctx context.Context, req *model.ReqBatchConsumeClientRequestLogEvent) {
	requestCounter := counter.RequestCounter{}

	for _, event := range req.EventList {
		var key counter.KeyRequestCounter
		key.Set(event.APIKey, timekit.TruncateToHour(event.Timestamp))
		requestCounter[key]++
	}

	for key, count := range requestCounter {
		apiKey, timestamp, err := key.Parse() // timestamp here is already truncated to hour
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn()
			continue
		}

		newCount, err := r.clientRequestCountRepository.IncrementCount(ctx, r.db, apiKey, timestamp, count)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn()
			continue
		}

		err = r.cache.SetClientRequestCountIfExist(ctx, apiKey, timestamp, newCount)
		x.LogIfErrContext(ctx, err)

		err = r.cache.IncrementTopClientRequestCountHourly(ctx, timestamp, count, apiKey)
		x.LogIfErrContext(ctx, err)
	}
}
