package route

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/topic"
	"github.com/segmentio/kafka-go"
)

type (
	FuncSingleMsgHandler func(ctx context.Context, message kafka.Message) error
	FuncBatchMsgHandler  func(ctx context.Context, messages []kafka.Message) error
)

type BatchHandlerConfig struct {
	Handler       FuncBatchMsgHandler
	BatchSize     int
	FlushInterval time.Duration
}

type (
	SingleMsgHandler map[string]FuncSingleMsgHandler
	BatchHandler     map[string]BatchHandlerConfig
)

// TopicHandler holds all topic handlers
type TopicHandler struct {
	SingleMsgHandler SingleMsgHandler
	BatchHandler     BatchHandler
}

func Setup(consumers *config.Consumers) TopicHandler {
	return TopicHandler{
		SingleMsgHandler: SingleMsgHandler{},
		BatchHandler: BatchHandler{
			topic.ClientRequestLog: {
				Handler:       consumers.RequestLogConsumer.ConsumeClientRequestLogEvent,
				BatchSize:     100,
				FlushInterval: 30 * time.Second,
			},
		},
	}
}
