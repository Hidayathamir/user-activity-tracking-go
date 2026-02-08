package messaging

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/requestlog"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/segmentio/kafka-go"
)

type RequestLogConsumer struct {
	cfg     *config.Config
	usecase requestlog.RequestLogUsecase
}

func NewRequestLogConsumer(cfg *config.Config, usecase requestlog.RequestLogUsecase) *RequestLogConsumer {
	return &RequestLogConsumer{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (r *RequestLogConsumer) ConsumeClientRequestLogEvent(ctx context.Context, messages []kafka.Message) error {
	req := new(model.ReqBatchConsumeClientRequestLogEvent)
	converter.KafkaMessageListToModelReqBatchConsumeClientRequestLogEvent(ctx, messages, req)

	err := r.usecase.BatchConsumeClientRequestLogEvent(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
