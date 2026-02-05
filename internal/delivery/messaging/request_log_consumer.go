package messaging

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/requestlog"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type RequestLogConsumer struct {
	Config  *viper.Viper
	Usecase requestlog.RequestLogUsecase
}

func NewRequestLogConsumer(cfg *viper.Viper, useCase requestlog.RequestLogUsecase) *RequestLogConsumer {
	return &RequestLogConsumer{
		Config:  cfg,
		Usecase: useCase,
	}
}

func (r *RequestLogConsumer) ConsumeClientRequestLogEvent(ctx context.Context, messages []kafka.Message) error {
	req := new(model.ReqBatchConsumeClientRequestLogEvent)
	converter.KafkaMessageListToModelReqBatchConsumeClientRequestLogEvent(ctx, messages, req)

	err := r.Usecase.BatchConsumeClientRequestLogEvent(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
