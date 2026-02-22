package router

import (
	"my_zhihu_backend/app/controller"
	"my_zhihu_backend/app/middleware"
	"my_zhihu_backend/app/service"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(r *gin.Engine, ctrl *controller.UserController, service *service.AuthService) {
	users := r.Group("/users")
	{
		users.POST("", ctrl.CreateNewUser)                                          // 创建用户
		users.GET("", middleware.Auth(service), ctrl.SearchUserByUsername)          // 搜索用户
		users.GET("/:id", middleware.Auth(service), ctrl.GetUser)                   // 获取用户信息
		users.DELETE("/me", middleware.Auth(service), ctrl.DeleteUser)              // 删除用户
		users.PATCH("/me", middleware.Auth(service), ctrl.UpdateUser)               // 更新用户信息
		users.POST("/follow/:id", middleware.Auth(service), ctrl.AddFollowing)      // 关注用户
		users.DELETE("/follow/:id", middleware.Auth(service), ctrl.RemoveFollowing) // 取消关注用户
		users.GET("/followers/:id", middleware.Auth(service), ctrl.GetFollowers)    // 获取粉丝列表
		users.GET("/followings/:id", middleware.Auth(service), ctrl.GetFollowings)  // 获取关注列表
	}
}
