package main

import (
	"log"
	"os"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")

	if err := db.Init(dsn); err != nil {
		log.Fatal("db connect failed:", err)
	}
	defer db.DB.Close()

	r := router.SetupRouter()

	r.Run(":8080") // 监听 8080 端口
}
