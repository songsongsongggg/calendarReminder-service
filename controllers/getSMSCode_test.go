package controllers_test

import (
	"calendarReminder-service/config"
	"calendarReminder-service/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 测试 GetSMSCode 函数的功能测试
func TestGetSMSCode(t *testing.T) {
	// 初始化 Mock Redis 客户端
	config.InitMockRedis()

	// 模拟 HTTP 请求
	req, err := http.NewRequest("GET", "/getSMSCode?mobile=15014354723", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 创建 HTTP ResponseRecorder 用于记录响应
	rr := httptest.NewRecorder()

	// 创建 Redis 客户端 mock
	mockRedisClient := config.GetRedisClient()

	// 使用闭包传递 Redis 客户端
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetSMSCode(w, r, mockRedisClient)
	})

	// 调用处理函数
	handler.ServeHTTP(rr, req)

	// 检查状态码是否为 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("处理程序返回错误状态码: 得到 %v 想要 %v", status, http.StatusOK)
	}

	// 检查响应体内容是否正确
	expected := "验证码已发送至 15014354723"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("处理程序返回了意外的响应体: 得到 %v 想要 %v", rr.Body.String(), expected)
	}
}
