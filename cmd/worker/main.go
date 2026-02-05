package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/config"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/messaging/route"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)

	db := config.NewDatabase(viperConfig)

	rdb := config.NewRedis(viperConfig)
	defer x.PanicIfErrForDefer(rdb.Close)

	var kafkaWriter *kafka.Writer = nil // we dont need kafka writer

	usecases := config.SetupUsecases(viperConfig, db, rdb, kafkaWriter)

	consumers := config.SetupConsumers(viperConfig, usecases)

	topicHandler := route.Setup(consumers)

	runWorker(viperConfig, topicHandler)
}

func runWorker(viperConfig *viper.Viper, topicHandler route.TopicHandler) {
	groupID := viperConfig.GetString(configkey.KafkaGroupID)

	localLogger := x.Logger.WithFields(logrus.Fields{
		"groupID": groupID,
	})

	localLogger.Debug("setup context for cancelation when receive signal notify")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		localLogger.Debug("listening signal notify")
		<-sigChan
		localLogger.Debug("receive signal notify, shutting down worker application")
		cancel()
	}()

	var wg sync.WaitGroup

	localLogger.Debug("start all single msg handler")
	for topicName, handler := range topicHandler.SingleMsgHandler {
		wg.Add(1)

		localLogger.WithFields(logrus.Fields{
			"topic": topicName,
		}).Debug("spawn goroutine for single message handler")

		go func(t string, h route.FuncSingleMsgHandler) {
			defer wg.Done()
			startConsumer(ctx, viperConfig, groupID, t, h)
		}(topicName, handler)
	}

	localLogger.Debug("start all batch msg handler")
	for topicName, batchHandlerConfig := range topicHandler.BatchHandler {
		wg.Add(1)

		localLogger.WithFields(logrus.Fields{
			"topic":             topicName,
			"batchSize":         batchHandlerConfig.BatchSize,
			"flushInterval (s)": batchHandlerConfig.FlushInterval.Seconds(),
		}).Debug("spawn goroutine for batch message handler")

		go func(t string, bhc route.BatchHandlerConfig) {
			defer wg.Done()
			startBatchConsumer(ctx, viperConfig, groupID, t, bhc)
		}(topicName, batchHandlerConfig)
	}

	localLogger.Debug("all workers started, waiting for messages...")
	localLogger.Debug("blocking runtime until all consumers stop")
	wg.Wait()
	localLogger.Debug("application stopped")
}

// startConsumer encapsulates the logic for a single topic consumer (processes messages one by one)
func startConsumer(ctx context.Context, viperConfig *viper.Viper, groupID, topicName string, handler route.FuncSingleMsgHandler) {
	localLogger := x.Logger.WithFields(logrus.Fields{
		"groupID": groupID,
		"topic":   topicName,
		"type":    "single consume",
	})

	kafkaReader := config.NewKafkaReader(viperConfig, groupID, topicName)
	defer func() {
		localLogger.Debug("closing consumer")
		err := kafkaReader.Close()
		if err != nil {
			localLogger.WithError(err).Debug("error close kafka reader")
		}
	}()

	localLogger.Debug("consumer started")

	for {
		if ctx.Err() != nil {
			localLogger.Debug("context cancelled, will stop loop listening")
			return
		}

		message, err := kafkaReader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				localLogger.WithError(err).Debug("error fetch message because context cancelled, stop loop listening")
				return
			}
			localLogger.WithError(err).Debug("error fetch message")
			time.Sleep(1 * time.Second)
			continue
		}

		err = handler(ctx, message)
		if err != nil {
			localLogger.WithError(err).Debug("handler failed")
		}

		err = kafkaReader.CommitMessages(ctx, message)
		if err != nil {
			localLogger.WithError(err).Debug("failed to commit message")
		}
	}
}

// startBatchConsumer encapsulates the logic for a batch topic consumer
// It collects messages and processes them in batches based on batch size or flush interval
func startBatchConsumer(ctx context.Context, viperConfig *viper.Viper, groupID, topicName string, bhc route.BatchHandlerConfig) {
	localLogger := x.Logger.WithFields(logrus.Fields{
		"groupID":           groupID,
		"topic":             topicName,
		"batchSize":         bhc.BatchSize,
		"flushInterval (s)": bhc.FlushInterval.Seconds(),
		"type":              "batch consume",
	})

	kafkaReader := config.NewKafkaReader(viperConfig, groupID, topicName)
	defer func() {
		localLogger.Debug("closing consumer")
		err := kafkaReader.Close()
		if err != nil {
			localLogger.WithError(err).Debug("error close kafka reader")
		}
	}()

	localLogger.Debug("consumer started")

	batch := make([]kafka.Message, 0, bhc.BatchSize)
	ticker := time.NewTicker(bhc.FlushInterval)
	defer ticker.Stop()

	// Channel to receive fetched messages
	messageChan := make(chan kafka.Message)
	errChan := make(chan error)

	// Goroutine to fetch messages
	go func() {
		for {
			if ctx.Err() != nil {
				localLogger.Debug("context cancelled, will stop loop listening")
				return
			}

			message, err := kafkaReader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					localLogger.WithError(err).Debug("error fetch message because context cancelled, stop loop listening")
					return
				}
				errChan <- err
				continue
			}
			messageChan <- message
		}
	}()

	// processBatch handles the batch processing and committing
	processBatch := func() {
		if len(batch) == 0 {
			return
		}

		localLogger.WithField("batchLen", len(batch)).Debug("processing batch")

		err := bhc.Handler(ctx, batch)
		if err != nil {
			localLogger.WithError(err).Debug("handler failed")
		}

		err = kafkaReader.CommitMessages(ctx, batch...)
		if err != nil {
			localLogger.WithError(err).Debug("failed to commit batch messages")
		}

		batch = batch[:0]
	}

	for {
		select {
		case <-ctx.Done():
			localLogger.Debug("context cancelled, will stop loop listening")
			return

		case <-ticker.C:
			// Flush interval reached, process current batch
			if len(batch) > 0 {
				localLogger.WithField("trigger", "interval").Debug("flush interval reached")
				processBatch()
			}

		case message := <-messageChan:
			batch = append(batch, message)
			localLogger.WithField("batchLen", len(batch)).Debug("message added to batch")

			// Check if batch is full
			if len(batch) >= bhc.BatchSize {
				localLogger.WithField("trigger", "size").Debug("batch size reached")
				processBatch()
			}

		case err := <-errChan:
			localLogger.WithError(err).Debug("error fetch message")
			time.Sleep(1 * time.Second)
		}
	}
}
