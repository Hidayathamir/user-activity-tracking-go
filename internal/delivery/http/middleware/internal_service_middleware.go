package middleware

import (
	"fmt"
	"strings"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/response"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/gin-gonic/gin"
)

func NewInternalServiceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const expectedSecret = "internal-secret"
		headerSecret := strings.TrimSpace(c.GetHeader("X-Internal-Secret"))

		if headerSecret == "" {
			err := fmt.Errorf("header secret empty")
			err = errkit.Unauthorized(err)
			err = errkit.AddFuncName(err)
			response.Error(c, err)
			c.Abort()
			return
		}

		if headerSecret != expectedSecret {
			err := fmt.Errorf("header secret invalid")
			err = errkit.Unauthorized(err)
			err = errkit.AddFuncName(err)
			response.Error(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
