package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Signal-zxh/signalzxh-blog/middleware"
	"github.com/Signal-zxh/signalzxh-blog/utils"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func TestAuth_NoToken(t *testing.T) {
	r := setupRouter()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Auth() got status code %d, want %d", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	if resp["message"] != "no token" {
		t.Errorf("Auth() got message %v, want 'no token'", resp["message"])
	}
}

func TestAuth_BadFormat(t *testing.T) {
	r := setupRouter()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Auth() got status code %d, want %d", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	if resp["message"] != "bad format" {
		t.Errorf("Auth() got message %v, want 'bad format'", resp["message"])
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	r := setupRouter()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Auth() got status code %d, want %d", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	if resp["message"] != "invalid token" {
		t.Errorf("Auth() got message %v, want 'invalid token'", resp["message"])
	}
}

func TestAuth_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := utils.GenerateToken(1)
	if err != nil {
		t.Fatalf("生成token失败：%v", err)
	}

	r := setupRouter()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			t.Error("user_id未设置")
		}
		if userID != 1 {
			t.Errorf("user_id got %v, want 1", userID)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Auth() got status code %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	if resp["code"] != float64(0) {
		t.Errorf("Auth() got code %v, want 0", resp["code"])
	}
}

func TestLogger_SetTraceId(t *testing.T) {
	r := setupRouter()
	r.GET("/test", middleware.Logger(), func(c *gin.Context) {
		traceId, exists := c.Get("traceId")
		if !exists {
			t.Error("traceId未设置")
		}
		if traceId == "" {
			t.Error("traceId为空")
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Logger() got status code %d, want %d", w.Code, http.StatusOK)
	}
}

func TestLogger_NextCalled(t *testing.T) {
	r := setupRouter()
	nextCalled := false
	r.GET("/test", middleware.Logger(), func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("next handler未被调用")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Logger() got status code %d, want %d", w.Code, http.StatusOK)
	}
}
