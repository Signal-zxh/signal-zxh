package handler_test

import (
	"encoding/json"
	"github.com/Signal-zxh/signal-zxh/router"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	// 登录获取token
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "123")

	r := router.SetupRouter()
	body := `{"username":"admin","password":"123"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// 创建响应记录器
	w := httptest.NewRecorder()
	// 模拟访问
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("login failed, got %d", w.Code)
	}
	// 解析响应体
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}
	// 提取token并验证非空
	tokenData, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中没有data字段")
	}
	token, ok := tokenData["token"].(string)
	if !ok || token == "" {
		t.Fatal("未获取token")
	}
}

// TODO: 完成剩下单元测试
