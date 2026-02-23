package middleware

import (
	"errors"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/log"
	"my_zhihu_backend/app/response"
	"my_zhihu_backend/app/service"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

var ErrInvalidAuthorizationHeader = app_error.NewInputError("invalid authorization header", app_error.ErrCodeInvalidAuthorizationHeader, nil)

func Auth(service *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		seq := strings.Split(h, " ")
		if len(seq) != 2 {
			_ = c.Error(ErrInvalidAuthorizationHeader)
			c.Abort()
			return
		}
		if seq[1] == "" {
			_ = c.Error(ErrInvalidAuthorizationHeader)
			c.Abort()
			return
		}
		id, err := service.ValidateAccessToken(seq[1])
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Set("id", id)
		c.Next()
	}
}

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				var appErr app_error.AppError
				if err, ok := e.(error); ok {
					appErr = app_error.NewInternalError(app_error.ErrCodeUnknown, err)
				} else {
					appErr = app_error.NewInternalError(app_error.ErrCodeUnknown, errors.New(e.(string)))
				}
				log.L().Error(appErr.Msg(), appErr.ErrorField()...)
				c.JSON(http.StatusOK, &response.Response{
					Code:          int(appErr.Code()),
					Message:       appErr.Msg(),
					Ok:            false,
					InternalError: true,
					Body:          nil,
				})
			}
		}()
		c.Next()
		if err := c.Errors.Last(); err != nil {
			appErr, ok := errors.AsType[app_error.AppError](err)
			if !ok {
				appErr = app_error.NewInternalError(app_error.ErrCodeUnknown, err)
			}
			log.L().Error(appErr.Msg(), appErr.ErrorField()...)
			c.JSON(http.StatusOK, &response.Response{
				Code:          int(appErr.Code()),
				Message:       appErr.Msg(),
				Ok:            false,
				InternalError: appErr.Type() == app_error.ErrTypeInternal,
				Body:          nil,
			})
		}
	}
}

func RateLimit() gin.HandlerFunc {
	limiters := new(sync.Map)
	return func(c *gin.Context) {
		addr := c.Request.RemoteAddr

		log.L().Info("request", zap.String("addr", addr))

		if _, exists := limiters.Load(addr); !exists {
			limiters.Store(addr, rate.NewLimiter(5, 1))
		}
		a, _ := limiters.Load(addr)
		limiter := a.(*rate.Limiter)
		if limiter.Allow() {
			c.Next()
			return
		}
		_ = c.Error(app_error.ErrTooManyRequests)
		c.Abort()
	}
}
