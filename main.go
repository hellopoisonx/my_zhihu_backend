package main

import (
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/controller"
	"my_zhihu_backend/app/dao"
	_ "my_zhihu_backend/app/log"
	"my_zhihu_backend/app/middleware"
	"my_zhihu_backend/app/repository"
	"my_zhihu_backend/app/router"
	"my_zhihu_backend/app/service"
	util2 "my_zhihu_backend/app/util"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	redisClient := repository.NewRedisClient(config.C)
	db := repository.NewMysqlDBConn(config.C)
	repository.AutoMigrate(db)
	userDAO := dao.NewUserDAO(config.C, db)
	authDAO := dao.NewAuthDAO(redisClient, config.C)
	articleDAO := dao.NewArticleDAO(db)
	util := new(util2.Util)
	userService := service.NewUserService(userDAO, config.C, util)
	authService := service.NewAuthService(authDAO, userDAO, config.C, util)
	articleService := service.NewArticleService(articleDAO, util)
	userController := controller.NewUserController(userService, config.C)
	authController := controller.NewAuthController(authService, config.C)
	articleController := controller.NewArticleController(articleService, config.C)

	r := gin.Default()
	r.Use(middleware.HandleError(), middleware.RateLimit())
	router.InitAuthRouter(r, authController, authService)
	router.InitUsersRouter(r, userController, authService)
	router.InitArticleRouter(r, articleController, authService)
	err := r.Run(config.C().App.ListenAddr)
	if err != nil {
		return
	}
}
