package router

import (
	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/middleware"
	_ "github.com/Signal-zxh/signalzxh-blog/docs"  // swagger docs
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(postHandler *handler.PostHandler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	// swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterPage(r)
	RegisterAuth(r, postHandler)
	RegisterAPI(r, postHandler)

	return r
}
