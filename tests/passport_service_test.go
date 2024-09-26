package tests__test

import (
	"calendarReminder-service/models"
	"calendarReminder-service/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// MockIDGenerator 用于模拟 ID 生成器
type MockIDGenerator struct{}

func (m *MockIDGenerator) GenerateUniqueID() (string, error) {
	return "mock_creator_id", nil // 返回一个固定的 ID
}

// TestUserService 测试 UserService
func TestUserService(t *testing.T) {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接到内存数据库: %v", err)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("无法自动迁移模型: %v", err)
	}

	// 创建 UserService 实例，使用 MockIDGenerator
	userService := services.NewUserService(db, &MockIDGenerator{})

	//// 测试创建用户
	mobile := "1234567890"
	user, err := userService.CreateUser(mobile)
	if err != nil {
		t.Fatalf("创建用户时出错: %v", err)
	}

	t.Logf("创建的用户信息: %+v", user) // 打印用户信息

	// 验证用户信息
	if user.Mobile != mobile {
		t.Errorf("期望手机号码为 '1234567890'，但实际为 %v", user.Mobile) // 错误信息使用中文
	}

	//// 测试手机号码是否存在
	exists, err := userService.QueryMobileIsExist(mobile)
	if err != nil {
		t.Fatalf("查询手机号码是否存在时出错: %v", err)
	}
	if !exists {
		t.Errorf("期望手机号码 '%s' 存在，但实际不存在", mobile)
	}

	//// 测试获取用户信息
	fetchedUser, err := userService.GetUserByMobile(mobile)
	if err != nil {
		t.Fatalf("根据手机号码获取用户信息时出错: %v", err)
	}

	t.Logf("测试获取用户信息: %+v", fetchedUser) // 打印用户信息

	// 验证获取的用户信息
	if fetchedUser == nil {
		t.Errorf("期望获取到用户信息，但实际返回 nil")
	} else if fetchedUser.Mobile != mobile {
		t.Errorf("期望获取的手机号码为 '%s'，但实际为 '%s'", mobile, fetchedUser.Mobile)
	}
}
