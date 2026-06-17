package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CreatePostRequest struct {
	Title string `json:"title"`
}

func main() {
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("db connect failed:", err)
	}

	log.Println("MySQL connected successfully")

	r := gin.Default()

	// 静态页面（知识图谱前端放这里）
	r.Static("/static", "./static")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	posts := []Post{
		{ID: 1, Title: "第一篇文章"},
		{ID: 2, Title: "第二篇文章"},
	}
	r.GET("/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})
	})

	r.GET("/posts/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		for _, post := range posts {
			if post.ID == id {
				c.JSON(http.StatusOK, post)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Post not found",
		})
	})

	r.POST("/posts", func(c *gin.Context) {
		var req CreatePostRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid json",
			})
			return
		}

		newID := len(posts) + 1
		newPost := Post{
			ID:    newID,
			Title: req.Title,
		}
		posts = append(posts, newPost)
		c.JSON(http.StatusOK, newPost)
	})

	r.Run(":8080") // 监听 8080 端口
}
