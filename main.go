package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/router"
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/Signal-zxh/signalzxh-blog/service/cache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// @title Signalzxh Blog API
// @version 1.0
// @description 博客系统的 RESTful API 接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load env failed:", err)
	}

	if err := db.InitDB(); err != nil {
		log.Fatal("db connect failed:", err)
	}
	defer db.DB.Close()

	if err := db.InitRedis(); err != nil {
		log.Fatal("redis connect failed:", err)
	}
	defer db.RDB.Close()

	postService := service.NewPostService(db.PostRepoImpl, cache.PostCacheImpl)
	postHandler := handler.NewPostHandler(postService)

	r := router.SetupRouter(postHandler)

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
