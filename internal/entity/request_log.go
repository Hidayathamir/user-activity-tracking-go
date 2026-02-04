package entity

import (
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/table"
)

type RequestLog struct {
	ID        int       `gorm:"column:id;primaryKey"`
	APIKey    string    `gorm:"column:api_key"`
	IP        string    `gorm:"column:ip"`
	Endpoint  string    `gorm:"column:endpoint"`
	Timestamp time.Time `gorm:"column:timestamp"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Client Client `gorm:"foreignKey:api_key;references:api_key"`
}

func (c *RequestLog) TableName() string {
	return table.RequestLog
}

type RequestLogList []RequestLog
