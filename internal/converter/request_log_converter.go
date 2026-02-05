package converter

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/segmentio/kafka-go"
)

func ModelReqRecordAPIHitToModelClientRequestLogEvent(req *model.ReqRecordAPIHit, event *model.ClientRequestLogEvent) {
	event.APIKey = req.APIKey
	event.IP = req.IP
	event.Endpoint = req.Endpoint
	event.Timestamp = req.Timestamp
}

func ModelClientRequestLogEventListToEntityRequestLogList(eventList *model.ClientRequestLogEventList, requestLogList *entity.RequestLogList) {
	for _, event := range *eventList {
		requestLog := entity.RequestLog{}
		ModelClientRequestLogEventToEntityRequestLog(&event, &requestLog)
		(*requestLogList) = append((*requestLogList), requestLog)
	}
}

func ModelClientRequestLogEventToEntityRequestLog(event *model.ClientRequestLogEvent, requestLog *entity.RequestLog) {
	requestLog.APIKey = event.APIKey
	requestLog.IP = event.IP
	requestLog.Endpoint = event.Endpoint
	requestLog.Timestamp = event.Timestamp
}

func KafkaMessageListToModelReqBatchConsumeClientRequestLogEvent(ctx context.Context, messages []kafka.Message, req *model.ReqBatchConsumeClientRequestLogEvent) {
	req.EventList = model.ClientRequestLogEventList{}

	for _, message := range messages {
		event := model.ClientRequestLogEvent{}

		err := json.Unmarshal(message.Value, &event)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn()
			continue
		}

		req.EventList = append(req.EventList, event)
	}
}
