package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/topic"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/MockProducer.go -pkg=mock . Producer

type Producer interface {
	SendClientRequestLogEvent(ctx context.Context, event *model.ClientRequestLogEvent) error
}

var _ Producer = &ProducerImpl{}

type ProducerImpl struct {
	Config      *viper.Viper
	KafkaWriter *kafka.Writer
}

func NewProducer(cfg *viper.Viper, kafkaWriter *kafka.Writer) *ProducerImpl {
	return &ProducerImpl{
		Config:      cfg,
		KafkaWriter: kafkaWriter,
	}
}

func (p *ProducerImpl) SendClientRequestLogEvent(ctx context.Context, event *model.ClientRequestLogEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	err = p.KafkaWriter.WriteMessages(ctx, kafka.Message{
		Topic: topic.ClientRequestLog,
		Value: value,
	})
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
