package router

import (
	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/Signal-zxh/signal-zxh/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())
	// 静态页面
	r.Static("/static", "./static")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// handler
	h := &handler.PostHandler{}

	r.GET("/posts", h.GetPosts)
	r.GET("/posts/:id", h.GetPostByID)
	r.POST("/posts", h.CreatePost)
	r.DELETE("/posts/:id", h.DeletePost)
	r.PUT("/posts/:id", h.UpdatePost)
	r.POST("/login", h.Login)

	return r
}
