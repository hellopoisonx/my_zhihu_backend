package controller

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrInvalidParameters = app_error.NewInputError("invalid parameters", app_error.ErrCodeInvalidParameters, nil)

func doWithBody[T any](c *gin.Context, ttl time.Duration, f func(ctx context.Context, req *T) (*response.Response, app_error.AppError)) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ttl)
	defer cancel()
	var req T
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(ErrInvalidParameters.WithError(err))
		return
	}
	if resp, err := f(timeout, &req); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func doWithQuery[T any](c *gin.Context, ttl time.Duration, f func(ctx context.Context, req *T) (*response.Response, app_error.AppError)) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ttl)
	defer cancel()
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(ErrInvalidParameters.WithError(err))
		return
	}
	if resp, err := f(timeout, &req); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func doWithUserId[T any](c *gin.Context, ttl time.Duration, f func(ctx context.Context, userId model.UserId, req *T) (*response.Response, app_error.AppError)) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ttl)
	defer cancel()
	id := getCurrentUserID(c)
	var req T
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(ErrInvalidParameters.WithError(err))
		return
	}
	if resp, err := f(timeout, model.UserId(id), &req); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func doOnlyWithUserId(c *gin.Context, ttl time.Duration, f func(ctx context.Context, userId model.UserId) (*response.Response, app_error.AppError)) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ttl)
	defer cancel()
	id := getCurrentUserID(c)
	if resp, err := f(timeout, model.UserId(id)); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func getIdFromParams(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func getCurrentUserID(c *gin.Context) int64 {
	userID := c.MustGet("id").(model.UserId)

	return int64(userID)
}
