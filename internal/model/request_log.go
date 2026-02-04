package model

import "time"

type ClientRequestLogEvent struct {
	APIKey    string    `json:"api_key"`
	IP        string    `json:"ip"`
	Endpoint  string    `json:"endpoint"`
	Timestamp time.Time `json:"timestamp"`
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
