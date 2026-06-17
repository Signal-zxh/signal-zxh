package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/gin-gonic/gin"
)

type PostHandler struct{}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := db.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req model.Post

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := db.CreatePost(req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	err = db.UpdatePost(id, req.Title)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = db.DeletePost(id)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	post, err := db.GetPostByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post not found",
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, post)
}
