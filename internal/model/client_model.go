package model

import "time"

type ReqRegister struct {
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResRegister struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReqLogin struct {
	Name     string `json:"name"     validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResLogin struct {
	Token string `json:"token"`
}

type ReqGetClientDetail struct {
	Name string `json:"name" validate:"required"`
}

type ResGetClientDetail struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReqVerify struct {
	Token string `validate:"required"`
}

type ClientAuth struct {
	ID              int       `gorm:"column:id"`
	Name            string    `gorm:"column:name"`
	Email           string    `gorm:"column:email"`
	EmailDecrypted  string    `gorm:"column:email_decrypted"`
	APIKey          string    `gorm:"column:api_key"`
	APIKeyDecrypted string    `gorm:"column:api_key_decrypted"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}
