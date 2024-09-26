package models

// 提醒实体类
type Reminder struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	CreatorID string   `gorm:"not null" json:"creator_id"`
	Content   string   `gorm:"not null" json:"content"`
	RemindAt  JSONTime `json:"remind_at"`  // 使用自定义时间类型
	CreatedAt JSONTime `json:"created_at"` // 使用自定义时间类型
	UpdatedAt JSONTime `json:"updated_at"` // 使用自定义时间类型
}
