package main

import (
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/controller"
	_ "my_zhihu_backend/app/log"
	"my_zhihu_backend/app/middleware"
	"my_zhihu_backend/app/repository"
	"my_zhihu_backend/app/router"
	"my_zhihu_backend/app/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	redisClient := repository.NewRedisClient(config.C)
	db := repository.NewMysqlDBConn(config.C)
	repository.AutoMigrate(db)
	userService := service.NewUserService(db, redisClient)
	authService := service.NewAuthService(db, redisClient)
	articleService := service.NewArticleService(db)
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)
	articleController := controller.NewArticleController(articleService)

	r := gin.Default()
	r.Use(
		cors.New(cors.Config{ // TODO: 跨域
			AllowAllOrigins:        true,
			AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:           []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:          []string{"Content-Length"},
			AllowCredentials:       true,
			AllowBrowserExtensions: true,
			AllowWebSockets:        true,
			AllowFiles:             true,
			MaxAge:                 12 * time.Hour,
		}),
		middleware.HandleError(), middleware.RateLimit(),
	)
	router.InitAuthRouter(r, authController, authService)
	router.InitUsersRouter(r, userController, authService, redisClient)
	router.InitArticleRouter(r, articleController, authService)
	err := r.Run(config.C().App.ListenAddr)
	if err != nil {
		return
	}
}
