package router

import (
	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuth(r *gin.Engine, h *handler.PostHandler) {
	auth := r.Group("/")
	auth.Use(middleware.Auth())

	auth.POST("/posts", h.CreatePost)
	auth.DELETE("/posts/:id", h.DeletePost)
	auth.PUT("/posts/:id", h.UpdatePost)
}
