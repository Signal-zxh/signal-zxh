package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Signal-zxh/signalzxh-blog/agent"
	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/Signal-zxh/signalzxh-blog/utils"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

type ToolHandler struct{}

// @Summary 获取文章列表
// @Description 分页获取文章列表
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=model.PageResponse{data=[]model.Post}}
// @Failure 500 {object} model.Response
// @Router /posts [get]
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

// @Summary 创建文章
// @Description 创建新文章（需要认证）
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{title=string,content=string} true "文章内容"
// @Success 200 {object} model.Response{data=object{id=int,title=string,content=string}}
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /posts [post]
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

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.Fail("invalid user_id"))
		return
	}

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

// @Summary 更新文章
// @Description 更新指定文章（需要认证）
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body object{title=string,content=string} true "更新内容"
// @Success 200 {object} model.Response{data=object{message=string,id=int,title=string}}
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /posts/{id} [put]
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

// @Summary 删除文章
// @Description 删除指定文章（需要认证）
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} model.Response{data=object{message=string}}
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /posts/{id} [delete]
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

// @Summary 获取文章详情
// @Description 根据ID获取单篇文章
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.Response{data=model.Post}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /posts/{id} [get]
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

// @Summary 用户登录
// @Description 使用用户名密码登录获取JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{username=string,password=string} true "登录信息"
// @Success 200 {object} model.Response{data=object{token=string}}
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /login [post]
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

// @Summary HTTP探测工具
// @Description 发送HTTP请求并返回响应状态和时间
// @Tags tools
// @Accept json
// @Produce json
// @Param request body object{url=string} true "探测URL"
// @Success 200 {object} model.Response{data=object{status=int,time_ms=int}}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /probe [post]
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

// @Summary AI Agent工具
// @Description 调用AI Agent处理查询请求
// @Tags tools
// @Accept json
// @Produce json
// @Param request body object{query=string} true "查询内容"
// @Success 200 {object} model.Response{data=object{result=string}}
// @Failure 400 {object} model.Response
// @Router /agent [post]
func (t *ToolHandler) Agent(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid request"))
		return
	}

	result := agent.RouteTool(req.Query)

	c.JSON(http.StatusOK, model.Success(gin.H{
		"result": result,
	}))
}
