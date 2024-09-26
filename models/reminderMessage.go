package models

// 定义一个结构体来封装提醒内容和手机号
type ReminderMessage struct {
	Content string `json:"content"`
	Mobile  string `json:"mobile"`
}
