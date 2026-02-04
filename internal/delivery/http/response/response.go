package response

import (
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/gin-gonic/gin"
)

type WebResponse[T any] struct {
	Data         T             `json:"data"`
	Paging       *PageMetadata `json:"paging"`
	ErrorMessage string        `json:"error_message"`
	ErrorDetail  []string      `json:"error_detail"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

func Data(c *gin.Context, status int, data any) {
	res := WebResponse[any]{}
	res.Data = data
	c.JSON(status, res)
}

func DataPaging(c *gin.Context, status int, data any, paging *PageMetadata) {
	res := WebResponse[any]{}
	res.Data = data
	res.Paging = paging
	c.JSON(status, res)
}

func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}

	httpErr := errkit.GetHTTPError(err)

	res := WebResponse[any]{}
	res.ErrorMessage = httpErr.Message
	res.ErrorDetail = errkit.Split(err)

	c.JSON(httpErr.HTTPCode, res)
}
