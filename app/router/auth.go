package router

import (
	"my_zhihu_backend/app/controller"
	"my_zhihu_backend/app/middleware"
	"my_zhihu_backend/app/service"

	"github.com/gin-gonic/gin"
)

func InitAuthRouter(r *gin.Engine, authController *controller.AuthController, service *service.AuthService) {
	auth := r.Group("/auth")
	{
		auth.POST("", authController.Login)
		auth.DELETE("", middleware.Auth(service), authController.Logout)
		auth.PATCH("", authController.Renew)
	}
}
