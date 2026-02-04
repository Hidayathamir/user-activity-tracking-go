package http

import (
	"net/http"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/delivery/http/response"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/usecase/client"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/ctx/ctxclientauth"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ClientController struct {
	Config  *viper.Viper
	Usecase client.ClientUsecase
}

func NewClientController(cfg *viper.Viper, useCase client.ClientUsecase) *ClientController {
	return &ClientController{
		Config:  cfg,
		Usecase: useCase,
	}
}

// Register godoc
//
//	@Summary		Register a new client
//	@Description	Register a new client with name, email, and password
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.ReqRegister	true	"Client Registration Request"
//	@Success		201		{object}	response.WebResponse[model.ResRegister]
//	@Router			/api/register [post]
func (_c *ClientController) Register(c *gin.Context) {
	req := new(model.ReqRegister)
	err := c.ShouldBindJSON(req)
	if err != nil {
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	res, err := _c.Usecase.Register(c.Request.Context(), req)
	if err != nil {
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, res)
}

// Login godoc
//
//	@Summary		Login a client
//	@Description	Login a client with name and password
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.ReqLogin	true	"Client Login Request"
//	@Success		200		{object}	response.WebResponse[model.ResLogin]
//	@Router			/api/login [post]
func (_c *ClientController) Login(c *gin.Context) {
	req := new(model.ReqLogin)
	err := c.ShouldBindJSON(req)
	if err != nil {
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	res, err := _c.Usecase.Login(c.Request.Context(), req)
	if err != nil {
		err = errkit.AddFuncName(err)
		x.Logger.WithContext(c.Request.Context()).WithError(err).Error()
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, res)
}

// GetClientDetail godoc
//
//	@Summary		Get client detail
//	@Description	Get client detail
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.WebResponse[model.ResGetClientDetail]
//	@Router			/api/client/me [get]
//	@Security		ApiKeyJWTAuth
func (_c *ClientController) GetClientDetail(c *gin.Context) {
	clientAuth := ctxclientauth.Get(c.Request.Context())

	res := new(model.ResGetClientDetail)
	converter.ModelClientAuthToModelResGetClientDetail(clientAuth, res)

	response.Data(c, http.StatusOK, res)
}
