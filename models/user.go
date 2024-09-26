package models

// 用户实体类
type User struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	Mobile    string   `gorm:"unique;not null" json:"mobile"`
	CreatorID string   `gorm:"unique;not null" json:"creator_id"`
	CreatedAt JSONTime `json:"created_at"`
	UpdatedAt JSONTime `json:"updated_at"`
}
