package controller

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/request"
	"my_zhihu_backend/app/response"
	"my_zhihu_backend/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *service.AuthService
	cfg     config.ReadConfigFunc
}

func NewAuthController(service *service.AuthService, cfg config.ReadConfigFunc) *AuthController {
	return &AuthController{service: service, cfg: cfg}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	doWithBody(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, req *request.AuthLoginRequest) (*response.Response, app_error.AppError) {
		at, rt, aExp, rExp, err := ctrl.service.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return &response.Response{
			Code:          0,
			Message:       "login",
			Ok:            true,
			InternalError: false,
			Body: response.AuthLoginResponse{
				AccessToken: response.TokenResponse{
					Token:    at,
					ExpireAt: aExp,
				},
				RefreshToken: response.TokenResponse{
					Token:    rt,
					ExpireAt: rExp,
				},
			},
		}, nil
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ctrl.cfg().Service.Timeout)
	defer cancel()
	id := c.MustGet("id").(model.UserId)
	if err := ctrl.service.Logout(timeout, id); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Code:          0,
		Message:       "logout",
		Ok:            true,
		InternalError: false,
	})
}

func (ctrl *AuthController) Renew(c *gin.Context) {
	doWithBody(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, req *request.AuthRenewAccessTokenRequest) (*response.Response, app_error.AppError) {
		token, exp, err := ctrl.service.RenewAccessToken(ctx, req)
		if err != nil {
			return nil, err
		}
		return &response.Response{
			Code:          0,
			Message:       "renewed",
			Ok:            true,
			InternalError: false,
			Body: response.TokenResponse{
				Token:    token,
				ExpireAt: exp,
			},
		}, nil
	})
}
