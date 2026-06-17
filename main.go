package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CreatePostRequest struct {
	Title string `json:"title"`
}

type UpdatePostRequest struct {
	Title string `json:"title"`
}

func main() {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	r.GET("/posts", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title FROM posts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var result []Post

		for rows.Next() {
			var p Post
			err := rows.Scan(&p.ID, &p.Title)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result = append(result, p)
		}

		c.JSON(http.StatusOK, gin.H{
			"posts": result,
		})
	})

	r.GET("/posts/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid id")
			return
		}

		row := db.QueryRow("SELECT id, title FROM posts WHERE id = ?", id)

		var post Post

		err = row.Scan(&post.ID, &post.Title)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post not found",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	})

	r.POST("/posts", func(c *gin.Context) {
		var req CreatePostRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := db.Exec("INSERT INTO posts(title) VALUES(?)", req.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := res.LastInsertId()

		c.JSON(http.StatusOK, gin.H{
			"id":    id,
			"title": req.Title,
		})
	})

	r.DELETE("/posts/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid id")
			return
		}

		res, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, _ := res.RowsAffected()

		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "deleted successfully",
		})
	})

	r.PUT("/posts/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid id")
			return
		}

		var req UpdatePostRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := db.Exec("UPDATE posts SET title = ? WHERE id = ?", req.Title, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, _ := res.RowsAffected()

		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "updated successfully",
			"id":      id,
			"title":   req.Title,
		})
	})

	r.Run(":8080") // 监听 8080 端口
}
