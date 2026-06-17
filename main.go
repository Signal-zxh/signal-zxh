package main

import (
	"log"
	"os"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type CreatePostRequest struct {
	Title string `json:"title"`
}

type UpdatePostRequest struct {
	Title string `json:"title"`
}

func main() {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")

	if err := db.Init(dsn); err != nil {
		log.Fatal("db connect failed:", err)
	}

	r := gin.Default()

	// 静态页面（知识图谱前端放这里）
	r.Static("/static", "./static")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	h := &handler.PostHandler{}
	r.GET("/posts", h.GetPosts)

	r.GET("/posts/:id", h.GetPostByID)

	r.POST("/posts", h.CreatePost)

	r.DELETE("/posts/:id", h.DeletePost)

	r.PUT("/posts/:id", h.UpdatePost)

	r.Run(":8080") // 监听 8080 端口
}
