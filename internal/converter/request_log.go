package converter

import "github.com/Hidayathamir/user-activity-tracking-go/internal/model"

func ModelReqRecordAPIHitToModelClientRequestLogEvent(req *model.ReqRecordAPIHit, event *model.ClientRequestLogEvent) {
	event.APIKey = req.APIKey
	event.IP = req.IP
	event.Endpoint = req.Endpoint
	event.Timestamp = req.Timestamp
}
