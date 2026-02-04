package entity

import (
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/table"
)

type Client struct {
	ID     int    `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name"`
	Email  string `gorm:"column:email"`
	APIKey string `gorm:"column:api_key"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (c *Client) TableName() string {
	return table.Client
}

type ClientList []Client
