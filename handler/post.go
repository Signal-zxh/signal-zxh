package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Signal-zxh/signal-zxh/agent"
	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/Signal-zxh/signal-zxh/service"
	"github.com/Signal-zxh/signal-zxh/utils"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

type ToolHandler struct{}

func (h *PostHandler) GetPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	posts, total, err := h.postService.GetPostsByPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail("internal error"))
		return
	}

	c.JSON(http.StatusOK, model.Success(model.PageResponse{
		Data:     posts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
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

	id, err := h.postService.CreatePost(req.Title, req.Content, uid)
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

	err = h.postService.UpdatePost(id, req.Title, req.Content)
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

	err = h.postService.DeletePost(id)
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

	post, err := h.postService.GetPostByID(id)
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

func (t *ToolHandler) HttpProbe(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}
	start := time.Now()
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, model.Fail("url is empty"))
		return
	}
	resp, err := http.Get(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		return
	}
	defer resp.Body.Close()
	cost := time.Since(start)
	c.JSON(http.StatusOK, model.Success(gin.H{
		"status":  resp.StatusCode,
		"time_ms": cost.Milliseconds(),
	}))
}

func (t *ToolHandler) Agent(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}

	c.BindJSON(&req)

	result := agent.RouteTool(req.Query)

	c.JSON(http.StatusOK, model.Success(gin.H{
		"result": result,
	}))
}
