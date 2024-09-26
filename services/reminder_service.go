package services

import (
	"calendarReminder-service/models"
	"gorm.io/gorm"
)

// ReminderService 提醒服务接口
type ReminderService interface {
	CreateReminder(reminder *models.Reminder) error
	GetRemindersByCreatorID(creatorID string) ([]models.Reminder, error)
	DeleteReminder(id string, creatorID string) error
	UpdateReminder(id string, reminder *models.Reminder, creatorID string) error
}

// ReminderServiceImpl 提醒服务实现
type ReminderServiceImpl struct {
	db *gorm.DB
}

// NewReminderService 创建 ReminderService 实现
func NewReminderService(db *gorm.DB) ReminderService {
	return &ReminderServiceImpl{db: db}
}

// CreateReminder 创建提醒
func (s *ReminderServiceImpl) CreateReminder(reminder *models.Reminder) error {
	return s.db.Create(reminder).Error
}

// GetRemindersByCreatorID 获取指定用户的提醒列表
func (s *ReminderServiceImpl) GetRemindersByCreatorID(creatorID string) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := s.db.Where("creator_id = ?", creatorID).Find(&reminders).Error
	return reminders, err
}

// DeleteReminder 删除提醒
func (s *ReminderServiceImpl) DeleteReminder(id string, creatorID string) error {
	return s.db.Where("id = ? AND creator_id = ?", id, creatorID).Delete(&models.Reminder{}).Error
}

// UpdateReminder 更新提醒
func (s *ReminderServiceImpl) UpdateReminder(id string, reminder *models.Reminder, creatorID string) error {
	return s.db.Model(&models.Reminder{}).Where("id = ? AND creator_id = ?", id, creatorID).Updates(reminder).Error
}
