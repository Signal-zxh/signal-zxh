package router

import (
	"net/http"

	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/gin-gonic/gin"
)

func RegisterAPI(r *gin.Engine) {
	h := &handler.PostHandler{}
	t := &handler.ToolHandler{}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Success(gin.H{
			"message": "pong",
		}))
	})

	r.POST("/login", h.Login)
	r.GET("/posts", h.GetPosts)
	r.GET("/posts/:id", h.GetPostByID)

	api := r.Group("/api/tools")
	api.POST("/http", t.HttpProbe)
	api.POST("/agent", t.Agent)
}
