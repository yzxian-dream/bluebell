package routers

import (
	"github.com/gin-contrib/pprof"
	"net/http"
	"webapp/controller"
	"webapp/settings"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//r := gin.New()
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})
	v1 := r.Group("api/v1")
	//v1.Use(middlewares.RateLimitMiddleware(2*time.Second, 1))
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.PostDetailHandler)
	v1.GET("/posts", controller.GetPostListHandler)
	v1.GET("/posts2", controller.GetPostList2Handler)
	v1.GET("/ping", func(context *gin.Context) {
		//如果是登录用户，判断请求头中是否有有效的jwt
		context.String(http.StatusOK, "pong")
	})
	v1.Use(controller.JWTAuthMiddleware())
	{

		v1.POST("/post", controller.CreatePostHandler)
		v1.POST("/vote", controller.PostVote)
	}

	pprof.Register(r)

}
