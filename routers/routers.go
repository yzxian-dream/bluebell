package routers

import (
	"net/http"
	"webapp/controller"
	"webapp/settings"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	v1 := r.Group("api/v1")
	v1.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})

	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/community", controller.JWTAuthMiddleware(), controller.CommunityHandler)
	v1.GET("/community/:id", controller.JWTAuthMiddleware(), controller.CommunityDetailHandler)
	v1.POST("/post", controller.JWTAuthMiddleware(), controller.CreatePostHandler)
	v1.GET("/post/:id", controller.JWTAuthMiddleware(), controller.PostDetailHandler)
	v1.GET("/post", controller.JWTAuthMiddleware(), controller.GetPostListHandler)
	v1.POST("/vote", controller.JWTAuthMiddleware(), controller.PostVote)
	v1.GET("/post2", controller.JWTAuthMiddleware(), controller.GetPostList2Handler)
	v1.GET("/ping", controller.JWTAuthMiddleware(), func(context *gin.Context) {
		//如果是登录用户，判断请求头中是否有有效的jwt
		context.String(http.StatusOK, "pong")
	})
	v1.Use(controller.JWTAuthMiddleware())
}
