package entity

import (
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/aes"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/table"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Client struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	APIKey   string `gorm:"column:api_key"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (c *Client) TableName() string {
	return table.Client
}

type ClientList []Client

func (c *Client) PrepareCredentialsForStorage() error {
	err := c.EncryptAndSetEmail(c.Email)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	err = c.HashAndSetPassword(c.Password)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	err = c.GenerateAndEncryptAndSetAPIKey()
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}

// HashAndSetPassword will hash password and set to Password field.
func (c *Client) HashAndSetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errkit.AddFuncName(err)
	}
	c.Password = string(hashedPassword)
	return nil
}

func (c *Client) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

// GenerateAndEncryptAndSetAPIKey will generate new uuid, encrypt it, and set to APIKey field.
func (c *Client) GenerateAndEncryptAndSetAPIKey() error {
	apiKeyUUID, err := uuid.NewV7()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	apiKeyEncrypted, err := aes.Encrypt(apiKeyUUID.String())
	if err != nil {
		return errkit.AddFuncName(err)
	}
	c.APIKey = apiKeyEncrypted
	return nil
}

// DecryptAPIKey return string of decrypted APIKey field.
func (c *Client) DecryptAPIKey() (string, error) {
	apiKey, err := aes.Decrypt(c.APIKey)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}
	return apiKey, nil
}

// EncryptAndSetEmail will encrypt email and set to Email field.
func (c *Client) EncryptAndSetEmail(email string) error {
	emailEncrypted, err := aes.Encrypt(email)
	if err != nil {
		return errkit.AddFuncName(err)
	}
	c.Email = emailEncrypted
	return nil
}

// DecryptEmail return string of decrypted Email field.
func (c *Client) DecryptEmail() (string, error) {
	email, err := aes.Decrypt(c.Email)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}
	return email, nil
}
