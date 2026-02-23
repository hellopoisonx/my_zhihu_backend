package router

import (
	"my_zhihu_backend/app/controller"
	"my_zhihu_backend/app/middleware"
	"my_zhihu_backend/app/service"

	"github.com/gin-gonic/gin"
)

func InitArticleRouter(r *gin.Engine, articleController *controller.ArticleController, authService *service.AuthService) {
	q := r.Group("/questions")
	q.Use(middleware.Auth(authService))
	{
		q.POST("", articleController.PostQuestion)
		q.DELETE("/:id", articleController.DeleteQuestion)
		q.PATCH("/:id", articleController.UpdateQuestion)
		q.GET("/:id", articleController.GetQuestion)
		q.GET("", articleController.ListQuestions)
	}
}
