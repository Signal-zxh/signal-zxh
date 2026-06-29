package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/Signal-zxh/signal-zxh/router"
)

type mockPostService struct{}

func (m *mockPostService) GetPostByID(id int) (model.Post, error) {
	return model.Post{}, nil
}

func (m *mockPostService) GetPosts() ([]model.Post, error) {
	return nil, nil
}

func (m *mockPostService) GetPostsByPage(page, pageSize int) ([]model.Post, int, error) {
	return []model.Post{
		{ID: 1, Title: "Test Post", Content: "Content"},
	}, 10, nil
}

func (m *mockPostService) CreatePost(title, content string, userID int) (int64, error) {
	return 0, nil
}

func (m *mockPostService) UpdatePost(id int, title, content string) error {
	return nil
}

func (m *mockPostService) DeletePost(id int) error {
	return nil
}

type mockPostServiceError struct{}

func (m *mockPostServiceError) GetPostByID(id int) (model.Post, error) {
	return model.Post{}, nil
}

func (m *mockPostServiceError) GetPosts() ([]model.Post, error) {
	return nil, nil
}

func (m *mockPostServiceError) GetPostsByPage(page, pageSize int) ([]model.Post, int, error) {
	return nil, 0, errors.New("internal error")
}

func (m *mockPostServiceError) CreatePost(title, content string, userID int) (int64, error) {
	return 0, nil
}

func (m *mockPostServiceError) UpdatePost(id int, title, content string) error {
	return nil
}

func (m *mockPostServiceError) DeletePost(id int) error {
	return nil
}

func TestLogin(t *testing.T) {
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "123456")

	postHandler := handler.NewPostHandler(&mockPostService{})
	r := router.SetupRouter(postHandler)

	body := `{"username":"admin","password":"123456"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("login failed, got %d", w.Code)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	tokenData, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中没有data字段")
	}
	token, ok := tokenData["token"].(string)
	if !ok || token == "" {
		t.Fatal("未获取token")
	}
}

func TestGetPosts_Success(t *testing.T) {
	postHandler := handler.NewPostHandler(&mockPostService{})
	r := router.SetupRouter(postHandler)

	req := httptest.NewRequest("GET", "/posts?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetPosts() got status code %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	if resp["code"] != float64(0) {
		t.Errorf("GetPosts() got code %v, want 0", resp["code"])
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中没有data字段")
	}

	posts, ok := data["data"].([]interface{})
	if !ok {
		t.Fatal("data.data 不是数组")
	}

	if len(posts) != 1 {
		t.Errorf("GetPosts() got %d posts, want 1", len(posts))
	}

	if data["total"] != float64(10) {
		t.Errorf("GetPosts() got total %v, want 10", data["total"])
	}

	if data["page"] != float64(1) {
		t.Errorf("GetPosts() got page %v, want 1", data["page"])
	}

	if data["page_size"] != float64(10) {
		t.Errorf("GetPosts() got page_size %v, want 10", data["page_size"])
	}
}

func TestGetPosts_DefaultParameters(t *testing.T) {
	postHandler := handler.NewPostHandler(&mockPostService{})
	r := router.SetupRouter(postHandler)

	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetPosts() got status code %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中没有data字段")
	}

	if data["page"] != float64(1) {
		t.Errorf("GetPosts() got page %v, want 1 (default)", data["page"])
	}

	if data["page_size"] != float64(10) {
		t.Errorf("GetPosts() got page_size %v, want 10 (default)", data["page_size"])
	}
}

func TestGetPosts_InvalidParameters(t *testing.T) {
	postHandler := handler.NewPostHandler(&mockPostService{})
	r := router.SetupRouter(postHandler)

	req := httptest.NewRequest("GET", "/posts?page=abc&page_size=xyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetPosts() got status code %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中没有data字段")
	}

	if data["page"] != float64(1) {
		t.Errorf("GetPosts() got page %v, want 1 (invalid page defaulted)", data["page"])
	}

	if data["page_size"] != float64(10) {
		t.Errorf("GetPosts() got page_size %v, want 10 (invalid page_size defaulted)", data["page_size"])
	}
}

func TestGetPosts_ServiceError(t *testing.T) {
	postHandler := handler.NewPostHandler(&mockPostServiceError{})
	r := router.SetupRouter(postHandler)

	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetPosts() got status code %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	if resp["code"] != float64(1) {
		t.Errorf("GetPosts() got code %v, want 1", resp["code"])
	}

	if resp["message"] != "internal error" {
		t.Errorf("GetPosts() got message %v, want 'internal error'", resp["message"])
	}
}
