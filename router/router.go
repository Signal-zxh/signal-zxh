package router

import (
	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(postHandler *handler.PostHandler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	RegisterPage(r)
	RegisterAuth(r, postHandler)
	RegisterAPI(r, postHandler)

	return r
}
