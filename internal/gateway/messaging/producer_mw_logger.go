package messaging

import (
	"context"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ Producer = &ProducerMwLogger{}

type ProducerMwLogger struct {
	Next Producer
}

func NewProducerMwLogger(next Producer) *ProducerMwLogger {
	return &ProducerMwLogger{
		Next: next,
	}
}

func (p *ProducerMwLogger) SendClientRequestLogEvent(ctx context.Context, event *model.ClientRequestLogEvent) error {
	err := p.Next.SendClientRequestLogEvent(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
