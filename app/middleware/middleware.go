package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/cache"
	"my_zhihu_backend/app/log"
	"my_zhihu_backend/app/response"
	"my_zhihu_backend/app/service"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

var l = log.L().With(zap.String("module", "middleware"))

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
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		clientIP := c.ClientIP()

		log.L().Info("request", zap.String("ip", clientIP))

		limiterInterface, _ := limiters.LoadOrStore(clientIP, rate.NewLimiter(5, 1))
		limiter := limiterInterface.(*rate.Limiter)

		timeout, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()

		if err := limiter.Wait(timeout); err != nil {
			_ = c.Error(app_error.ErrTooManyRequests.WithError(err))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

type queryResponseHijack struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (q *queryResponseHijack) Write(b []byte) (int, error) {
	q.body.Write(b)
	return q.ResponseWriter.Write(b)
}

func CacheQuery(client *redis.Client, prefix, filterName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cacher := cache.NewJsonCacher(client, 15*time.Minute, prefix, func(ctx context.Context, args ...any) (*response.Response, app_error.AppError) {
			return nil, app_error.ErrRedisCacheKeyNotExists
		}, cache.NewBloomFilter(filterName, client))
		timeout, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()
		resp, err := cacher.Get(timeout, c.Request.URL.RequestURI(), c)
		if err == nil && resp != nil { // 缓存命中 直接返回
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		hijack := &queryResponseHijack{c.Writer, bytes.NewBufferString("")} // 缓存未命中 开始劫持ResponseWriter
		c.Writer = hijack
		c.Next()
		if c.Writer.Status() == http.StatusOK && c.Errors.Last() == nil { // controller 正常响应
			go func(rawJson []byte, key string) {
				var resp response.Response
				if err := json.Unmarshal(rawJson, &resp); err != nil {
					l.Error("failed to unmarshal json to response.Response", app_error.ErrInvalidJsonBody.WithError(err).ErrorField()...)
					return
				}
				if err := cacher.Put(context.Background(), key, resp); err != nil {
					l.Error("failed to cache response", err.ErrorField()...)
				}
			}(hijack.body.Bytes(), c.Request.URL.RequestURI())
		}
	}
}
