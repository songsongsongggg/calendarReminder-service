package tests__test

import (
	"calendarReminder-service/models"
	"calendarReminder-service/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

// initDB 初始化测试数据库
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	// 运行自动迁移
	err = db.AutoMigrate(&models.Reminder{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

// 测试创建提醒功能
func TestCreateReminder(t *testing.T) {
	db := initDB() // 初始化 GORM 数据库
	service := services.NewReminderService(db)

	reminder := models.Reminder{
		CreatorID: "test_user",
		Content:   "Test Reminder",
		RemindAt:  models.JSONTime{Time: time.Now().Add(time.Hour).Truncate(time.Second)},
	}

	// 执行创建提醒操作
	err := service.CreateReminder(&reminder)
	assert.NoError(t, err, "创建提醒失败")
	log.Printf("成功创建提醒: %+v\n", reminder)
}
