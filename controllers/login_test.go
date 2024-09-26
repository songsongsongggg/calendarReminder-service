package controllers_test

import (
	"bytes"
	"calendarReminder-service/config"
	"calendarReminder-service/controllers"
	"calendarReminder-service/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"log"
)

// MockUserService 用于模拟 UserService 的行为
type MockUserService struct {
	mock.Mock
}

// 模拟 QueryMobileIsExist 方法，返回是否存在该手机号
func (m *MockUserService) QueryMobileIsExist(mobile string) (bool, error) {
	args := m.Called(mobile)
	return args.Bool(0), args.Error(1)
}

// 模拟 GetUserByMobile 方法，返回用户信息
func (m *MockUserService) GetUserByMobile(mobile string) (*models.User, error) {
	args := m.Called(mobile)
	return args.Get(0).(*models.User), args.Error(1)
}

// 模拟 CreateUser 方法，返回新创建的用户
func (m *MockUserService) CreateUser(mobile string) (*models.User, error) {
	args := m.Called(mobile)
	return args.Get(0).(*models.User), args.Error(1)
}

// 测试 Login 函数的单元测试
func TestLogin(t *testing.T) {
	log.Println("启动 Login 函数测试")

	// 初始化 Mock Redis（可以使用实际 Redis 连接或者 mock）
	config.InitMockRedis()

	// 预设 Redis 中的验证码，确保验证码校验通过
	mobile := "15014354723"
	redisKey := fmt.Sprintf("%s:%s", controllers.MOBILE_SMSCODE, mobile)
	expectedSmsCode := "123456"
	err := config.RedisClient.Set(config.Ctx, redisKey, expectedSmsCode, 30*time.Minute).Err()
	if err != nil {
		t.Fatalf("设置 Redis 验证码失败: %v", err)
	}

	// 创建 MockUserService 实例
	mockUserService := new(MockUserService)

	// 假设手机号为 15014354723，模拟用户已存在的情况
	mockUserService.On("QueryMobileIsExist", "15014354723").Return(true, nil)
	mockUserService.On("GetUserByMobile", "15014354723").Return(&models.User{CreatorID: "12345"}, nil)

	// 构造请求体，模拟登录请求
	reqBody := map[string]string{
		"mobile":  "15014354723", // 模拟的手机号码
		"smsCode": "123456",      // 模拟的验证码
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	// 创建 HTTP POST 请求，用于模拟登录请求
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("创建 HTTP 请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json") // 设置请求头为 JSON 格式

	// 创建 ResponseRecorder 记录 HTTP 响应结果
	rr := httptest.NewRecorder()

	// 构建 HTTP 处理器，调用 Login 函数进行测试
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Login(w, r, mockUserService, config.RedisClient)
	})

	// 发送 HTTP 请求，调用被测试的 Login 函数
	handler.ServeHTTP(rr, req)

	// 检查 HTTP 响应状态码是否为 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Login 处理程序返回了错误状态码: 得到 %v, 预期 %v", status, http.StatusOK)
	} else {
		log.Printf("Login 请求成功，状态码: %v", rr.Code)
	}

	// 检查响应体内容是否为“登录成功”
	expected := "登录成功"
	if rr.Body.String() != expected {
		t.Errorf("Login 处理程序返回了意外的响应体: 得到 %v, 预期 %v", rr.Body.String(), expected)
	} else {
		log.Printf("Login 请求返回的响应体正确: %v", rr.Body.String())
	}

	log.Println("Login 函数测试完成")
}
