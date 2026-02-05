package http

import (
	"net/http"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/response"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/requestlog"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type RequestLogController struct {
	Config  *viper.Viper
	Usecase requestlog.RequestLogUsecase
}

func NewRequestLogController(cfg *viper.Viper, useCase requestlog.RequestLogUsecase) *RequestLogController {
	return &RequestLogController{
		Config:  cfg,
		Usecase: useCase,
	}
}

// RecordAPIHit godoc
//
//	@Summary		Record api hit
//	@Description	Record api hit
//	@Tags			RequestLog
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.ReqRecordAPIHit	true	"Record API Hit Request"
//	@Success		201		{object}	response.WebResponse[model.ResRecordAPIHit]
//	@Router			/api/logs [post]
//	@Security		ApiKeyXInternalSecret
func (r *RequestLogController) RecordAPIHit(c *gin.Context) {
	req := new(model.ReqRecordAPIHit)
	err := c.ShouldBindJSON(req)
	if err != nil {
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	res, err := r.Usecase.RecordAPIHit(c.Request.Context(), req)
	if err != nil {
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, res)
}

// GetTop3ClientRequestCount24Hour godoc
//
//	@Summary		Get top 3 client request count in last 24 hours
//	@Description	Get top 3 client request count in last 24 hours
//	@Tags			RequestLog
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.WebResponse[model.ResGetTop3ClientRequestCount24Hour]
//	@Router			/api/usage/top [get]
//	@Security		ApiKeyJWTAuth
func (r *RequestLogController) GetTop3ClientRequestCount24Hour(c *gin.Context) {
	req := new(model.ReqGetTop3ClientRequestCount24Hour)

	res, err := r.Usecase.GetTop3ClientRequestCount24Hour(c.Request.Context(), req)
	if err != nil {
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, res)
}
