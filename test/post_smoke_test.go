package test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

const baseURL = "http://localhost:8080"

func TestLogin(t *testing.T) {
	resp, _ := http.Post(baseURL+"/login", "application/json",
		strings.NewReader(`{"username":"admin","password":"123"}`))

	defer resp.Body.Close()

	// 解析响应体
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})

	token := data["token"].(string)
	if token == "" {
		t.Fatal("token不存在或为空")
	}
	t.Logf("登录成功，token: %s", token)
}
