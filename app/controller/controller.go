package controller

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func do[T any](c *gin.Context, ttl time.Duration, f func(ctx context.Context, req *T) (*response.Response, app_error.AppError)) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ttl)
	defer cancel()
	var req T
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(app_error.NewInputError("invalid parameters", app_error.ErrCodeInvalidParameters, err))
		return
	}
	if resp, err := f(timeout, &req); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func doQuery[T any](c *gin.Context, ttl time.Duration, f func(ctx context.Context, req *T) (*response.Response, app_error.AppError)) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ttl)
	defer cancel()
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(app_error.NewInputError("invalid parameters", app_error.ErrCodeInvalidParameters, err))
		return
	}
	if resp, err := f(timeout, &req); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}
