package entity

import (
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/table"
)

type ClientRequestCount struct {
	ID       int       `gorm:"column:id;primaryKey"`
	APIKey   string    `gorm:"column:api_key"`
	Datetime time.Time `gorm:"column:datetime"`
	Count    int       `gorm:"column:count"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Client Client `gorm:"foreignKey:api_key;references:api_key"`
}

func (c *ClientRequestCount) TableName() string {
	return table.ClientRequestCount
}

type ClientRequestCountList []ClientRequestCount

type Top3ClientRequestCount struct {
	APIKey   string `gorm:"column:api_key"`
	TotalSum int    `gorm:"column:total_sum"`
}

type Top3ClientRequestCountList []Top3ClientRequestCount
