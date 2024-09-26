package controllers_test

import (
	"bytes"
	"calendarReminder-service/controllers"
	"calendarReminder-service/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockReminderService 是一个模拟的提醒服务，用于测试
type MockReminderService struct {
	CreateReminderFunc          func(reminder *models.Reminder) error
	GetRemindersByCreatorIDFunc func(creatorID string) ([]models.Reminder, error)
	DeleteReminderFunc          func(id, creatorID string) error
	UpdateReminderFunc          func(id string, reminder *models.Reminder, creatorID string) error
}

// 实现 ReminderService 接口
func (m *MockReminderService) CreateReminder(reminder *models.Reminder) error {
	return m.CreateReminderFunc(reminder)
}

func (m *MockReminderService) GetRemindersByCreatorID(creatorID string) ([]models.Reminder, error) {
	return m.GetRemindersByCreatorIDFunc(creatorID)
}

func (m *MockReminderService) DeleteReminder(id, creatorID string) error {
	return m.DeleteReminderFunc(id, creatorID)
}

func (m *MockReminderService) UpdateReminder(id string, reminder *models.Reminder, creatorID string) error {
	return m.UpdateReminderFunc(id, reminder, creatorID)
}

// 测试创建提醒接口
func TestCreateReminder(t *testing.T) {
	service := &MockReminderService{
		CreateReminderFunc: func(reminder *models.Reminder) error {
			// 模拟成功的创建行为
			return nil
		},
	}

	// 创建一个 HTTP 请求
	reminder := models.Reminder{Content: "测试提醒"}
	body, _ := json.Marshal(reminder)
	req := httptest.NewRequest(http.MethodPost, "/reminders", bytes.NewBuffer(body))
	req.AddCookie(&http.Cookie{Name: "creator_id", Value: "test_user"})

	// 创建响应记录器
	rr := httptest.NewRecorder()

	// 调用 CreateReminder 控制器
	controllers.CreateReminder(rr, req, service)

	// 验证响应状态码
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusCreated, status)
	}
}

// 测试获取提醒接口
func TestGetReminders(t *testing.T) {
	service := &MockReminderService{
		GetRemindersByCreatorIDFunc: func(creatorID string) ([]models.Reminder, error) {
			return []models.Reminder{{CreatorID: creatorID, Content: "测试提醒"}}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/reminders", nil)
	req.AddCookie(&http.Cookie{Name: "creator_id", Value: "test_user"})
	rr := httptest.NewRecorder()

	controllers.GetReminders(rr, req, service)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusOK, status)
	}
}

// 测试删除提醒接口
func TestDeleteReminder(t *testing.T) {
	service := &MockReminderService{
		DeleteReminderFunc: func(id, creatorID string) error {
			return nil // 模拟成功删除
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/reminders?id=1", nil)
	req.AddCookie(&http.Cookie{Name: "creator_id", Value: "test_user"})
	rr := httptest.NewRecorder()

	controllers.DeleteReminder(rr, req, service)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusNoContent, status)
	}
}

// 测试更新提醒接口
func TestUpdateReminder(t *testing.T) {
	service := &MockReminderService{
		UpdateReminderFunc: func(id string, reminder *models.Reminder, creatorID string) error {
			return nil // 模拟成功更新
		},
	}

	reminder := models.Reminder{Content: "更新的测试提醒"}
	body, _ := json.Marshal(reminder)
	req := httptest.NewRequest(http.MethodPut, "/reminders?id=1", bytes.NewBuffer(body))
	req.AddCookie(&http.Cookie{Name: "creator_id", Value: "test_user"})
	rr := httptest.NewRecorder()

	controllers.UpdateReminder(rr, req, service)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusOK, status)
	}
}
