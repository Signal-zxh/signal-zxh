package handler

import (
	"net/http"
	"strconv"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/gin-gonic/gin"
)

type PostHandler struct{}

func (h *PostHandler) GetPosts(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, title FROM posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []model.Post

	for rows.Next() {
		var p model.Post
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
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req model.Post

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := db.DB.Exec("INSERT INTO posts(title) VALUES(?)", req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"title": req.Title,
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	var req model.Post

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := db.DB.Exec("UPDATE posts SET title = ? WHERE id = ?", req.Title, id)
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
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	res, err := db.DB.Exec("DELETE FROM posts WHERE id = ?", id)
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
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	row := db.DB.QueryRow("SELECT id, title FROM posts WHERE id = ?", id)

	var post model.Post

	err = row.Scan(&post.ID, &post.Title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	c.JSON(http.StatusOK, post)
}
