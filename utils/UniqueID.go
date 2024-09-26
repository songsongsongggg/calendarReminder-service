package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
)

// IDGenerator 接口定义
type IDGenerator interface {
	GenerateUniqueID() (string, error)
}

// SimpleIDGenerator 是 ID 生成器的简单实现
type SimpleIDGenerator struct{}

// GenerateUniqueID 实现 IDGenerator 接口中的方法
func (g *SimpleIDGenerator) GenerateUniqueID() (string, error) {
	return GenerateUniqueID() // 调用全局的 GenerateUniqueID 方法
}

// GenerateUniqueID 生成一个唯一的 16 位 ID
func GenerateUniqueID() (string, error) {
	bytes := make([]byte, 8) // 8 字节 = 64 位
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// 使用十六进制表示，将字节转换为字符串，并截取前 16 位
	id := fmt.Sprintf("%x", bytes) // 转换为十六进制字符串
	return id[:16], nil            // 返回前 16 位
}

// GenerateUUID 生成唯一的 UUID
func GenerateUUID() string {
	return uuid.New().String() // 生成并返回一个 UUID 字符串
}
