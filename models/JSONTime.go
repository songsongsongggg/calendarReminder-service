package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONTime 自定义时间类型，处理时间解析与格式化
type JSONTime struct {
	time.Time
}

// 时间格式常量
const timeFormat = "2006-01-02 15:04:05"

// MarshalJSON 实现 json.Marshal 接口，确保返回 "YYYY-MM-DD HH:MM:SS" 格式
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	formatted := jt.Time.Format(timeFormat) // 格式化为 "YYYY-MM-DD HH:MM:SS"
	return json.Marshal(formatted)
}

// UnmarshalJSON 实现 json.Unmarshal 接口，解析任何传入格式并转换为 "YYYY-MM-DD HH:MM:SS"
func (jt *JSONTime) UnmarshalJSON(data []byte) error {
	var t string
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	// 解析时允许不同格式，但最终统一转换为 "YYYY-MM-DD HH:MM:SS"
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		// 如果解析 RFC3339 失败，尝试使用自定义格式
		parsedTime, err = time.Parse(timeFormat, t)
		if err != nil {
			return err
		}
	}
	jt.Time = parsedTime
	return nil
}

// Scan 实现 sql.Scanner 接口，用于从数据库中读取 time.Time 数据
func (jt *JSONTime) Scan(value interface{}) error {
	if value == nil {
		jt.Time = time.Time{}
		return nil
	}
	v, ok := value.(time.Time)
	if !ok {
		return errors.New("无法转换为 time.Time")
	}
	jt.Time = v
	return nil
}

// Value 实现 driver.Valuer 接口，将 JSONTime 转换为数据库支持的格式
func (jt JSONTime) Value() (driver.Value, error) {
	return jt.Time, nil
}
