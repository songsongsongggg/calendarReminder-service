package controllers_test

import (
	"calendarReminder-service/config"
	"calendarReminder-service/controllers"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试 Logout 函数的单元测试
func TestLogout(t *testing.T) {
	log.Println("启动 Logout 函数测试")

	// 初始化 Mock Redis（可以使用实际 Redis 连接或者 mock）
	config.InitRedis()

	// 构造 HTTP GET 请求，模拟登出请求
	req, err := http.NewRequest("GET", "/logout?creator_id=12345", nil) // 模拟带有 userID 的登出请求
	if err != nil {
		t.Fatalf("创建 HTTP 请求失败: %v", err)
	}

	// 创建 ResponseRecorder 用于记录 HTTP 响应结果
	rr := httptest.NewRecorder()

	// 构建 HTTP 处理器，调用 Logout 函数进行测试
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Logout(w, r, config.RedisClient)
	})

	// 发送 HTTP 请求，调用被测试的 Logout 函数
	handler.ServeHTTP(rr, req)

	// 检查 HTTP 响应状态码是否为 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Logout 处理程序返回了错误状态码: 得到 %v, 预期 %v", status, http.StatusOK)
	} else {
		log.Printf("Logout 请求成功，状态码: %v", rr.Code)
	}

	// 检查响应体内容是否为“登出成功”
	expected := "登出成功\n"
	if rr.Body.String() != expected {
		t.Errorf("Logout 处理程序返回了意外的响应体: 得到 %v, 预期 %v", rr.Body.String(), expected)
	} else {
		log.Printf("Logout 请求返回的响应体正确: %v", rr.Body.String())
	}

	log.Println("Logout 函数测试完成")
}
