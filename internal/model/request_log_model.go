package model

import "time"

type ClientRequestLogEvent struct {
	APIKey    string    `json:"api_key"`
	IP        string    `json:"ip"`
	Endpoint  string    `json:"endpoint"`
	Timestamp time.Time `json:"timestamp"`
}

type ClientRequestLogEventList []ClientRequestLogEvent

type ReqBatchConsumeClientRequestLogEvent struct {
	EventList ClientRequestLogEventList
}

type ReqRecordAPIHit struct {
	APIKey    string    `json:"api_key"     validate:"required"`
	IP        string    `json:"ip"          validate:"required"`
	Endpoint  string    `json:"endpoint"    validate:"required"`
	Timestamp time.Time `format:"date-time" json:"timestamp"    swaggertype:"string" validate:"required"`
}

type ResRecordAPIHit struct {
	OK bool `json:"ok"`
}

type ReqGetTop3ClientRequestCount24Hour struct {
}

type ResGetTop3ClientRequestCount24Hour struct {
	NameCountList NameCountList `json:"name_count_list"`
}

type NameCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type NameCountList []NameCount

type APIKeyCount struct {
	APIKey string
	Count  int
}

type APIKeyCountList []APIKeyCount
