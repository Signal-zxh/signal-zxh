package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/Signal-zxh/signal-zxh/service"
	"github.com/Signal-zxh/signal-zxh/utils"
	"github.com/gin-gonic/gin"
)

type PostHandler struct{}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := service.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail("internal error"))
		return
	}

	c.JSON(http.StatusOK, model.Success(posts))
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Fail("no user"))
		return
	}

	uid := userID.(int)

	id, err := service.CreatePost(req.Title, req.Content, uid)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, model.Fail("invalid input"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"id":      id,
		"title":   req.Title,
		"content": req.Content,
	}))
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid id"))
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	err = service.UpdatePost(id, req.Title, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.Fail("post not found"))
		} else if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, model.Fail("invalid input"))
		} else {
			c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"message": "updated successfully",
		"id":      id,
		"title":   req.Title,
	}))
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid id"))
		return
	}

	err = service.DeletePost(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.Fail("post not found"))
		} else {
			c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"message": "deleted successfully",
	}))
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid id"))
		return
	}

	post, err := service.GetPostByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.Fail("post not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(post))
}

func (h *PostHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("bad request"))
		return
	}

	if req.Username != os.Getenv("ADMIN_USERNAME") || req.Password != os.Getenv("ADMIN_PASSWORD") {
		c.JSON(http.StatusUnauthorized, model.Fail("invalid credentials"))
		return
	}

	token, err := utils.GenerateToken(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail("failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"token": token,
	}))
}
