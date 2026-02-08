package config

import (
	"fmt"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func NewKafkaWriter(cfg *Config) *kafka.Writer {
	address := cfg.GetKafkaAddress()

	waitForKafka(address, 30, 2*time.Second)

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(address),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return kafkaWriter
}

func NewKafkaReader(cfg *Config, groupID string, topic string) *kafka.Reader {
	address := cfg.GetKafkaAddress()

	waitForKafka(address, 30, 2*time.Second)

	const _10MB = 10e6

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: _10MB,
	})

	return kafkaReader
}

// waitForKafka will waits for Kafka to be fully ready with retry logic.
func waitForKafka(address string, maxRetries int, retryInterval time.Duration) {
	for i := range maxRetries {
		conn, err := kafka.Dial("tcp", address)
		if err != nil {
			x.Logger.WithFields(logrus.Fields{
				"attempt": fmt.Sprintf("%d/%d", i+1, maxRetries),
				"err":     err.Error(),
				"sleep":   retryInterval,
			}).Warn("waiting for Kafka to be ready...")
			time.Sleep(retryInterval)
			continue
		}
		_, err = conn.Brokers()
		x.PanicIfErr(conn.Close())
		if err == nil {
			return
		}
		x.Logger.WithFields(logrus.Fields{
			"attempt": fmt.Sprintf("%d/%d", i+1, maxRetries),
			"err":     err.Error(),
			"sleep":   retryInterval,
		}).Warn("kafka connected but not ready...")
		time.Sleep(retryInterval)
	}
	err := fmt.Errorf("kafka not ready after %d attempts", maxRetries)
	panic(err)
}
