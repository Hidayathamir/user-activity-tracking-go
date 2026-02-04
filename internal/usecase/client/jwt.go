package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/configkey"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/golang-jwt/jwt/v5"
)

func (c *ClientUsecaseImpl) signAccessToken(ctx context.Context, clientName string) (string, error) {
	secret := c.Config.GetString(configkey.AuthJWTSecret)
	if secret == "" {
		err := fmt.Errorf("jwt secret is not configured")
		return "", errkit.AddFuncName(err)
	}

	expireSeconds := c.Config.GetInt(configkey.AuthJWTExpireSeconds)
	if expireSeconds <= 0 {
		err := fmt.Errorf("jwt expire seconds must be greater than zero")
		return "", errkit.AddFuncName(err)
	}

	issuer := c.Config.GetString(configkey.AuthJWTIssuer)
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   clientName,
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	return tokenString, nil
}

func (c *ClientUsecaseImpl) parseAccessToken(ctx context.Context, tokenString string) (clientName string, err error) {
	if tokenString == "" {
		err := fmt.Errorf("token is empty")
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName(err)
	}

	secret := c.Config.GetString(configkey.AuthJWTSecret)
	if secret == "" {
		err := fmt.Errorf("jwt secret is not configured")
		return "", errkit.AddFuncName(err)
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName(err)
	}

	if !token.Valid {
		err := fmt.Errorf("token is invalid")
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName(err)
	}

	if claims.Subject == "" {
		err := fmt.Errorf("token subject is empty")
		err = errkit.Unauthorized(err)
		return "", errkit.AddFuncName(err)
	}

	clientName = claims.Subject

	return clientName, nil
}
