package middleware

import (
	"fmt"
	"strings"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/response"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/client"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/ctx/ctxclientauth"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(clientUsecase client.ClientUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAuth := strings.TrimSpace(c.GetHeader("Authorization"))
		if headerAuth == "" {
			err := fmt.Errorf("header auth empty")
			err = errkit.Unauthorized(err)
			err = errkit.AddFuncName(err)
			response.Error(c, err)
			return
		}

		var token string
		parts := strings.Fields(headerAuth)
		switch {
		case len(parts) == 1:
			token = parts[0]
		case len(parts) == 2 && strings.EqualFold(parts[0], "Bearer"):
			token = parts[1]
		default:
			err := fmt.Errorf("authorization header format invalid")
			err = errkit.Unauthorized(err)
			err = errkit.AddFuncName(err)
			response.Error(c, err)
			return
		}

		req := &model.ReqVerify{
			Token: token,
		}

		clientAuth, err := clientUsecase.Verify(c.Request.Context(), req)
		if err != nil {
			err = errkit.Unauthorized(err)
			err = errkit.AddFuncName(err)
			response.Error(c, err)
			return
		}

		c.Request = c.Request.WithContext(ctxclientauth.Set(c.Request.Context(), clientAuth))

		c.Next()
	}
}
