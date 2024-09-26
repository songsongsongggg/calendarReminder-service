package services

import (
	"calendarReminder-service/models"
	"calendarReminder-service/utils"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

// UserService 接口定义
type UserService interface {
	QueryMobileIsExist(mobile string) (bool, error)
	GetUserByMobile(mobile string) (*models.User, error)
	GetUserByCreatorID(CreatorID string) (*models.User, error)
	CreateUser(mobile string) (*models.User, error)
}

// UserServiceImpl 结构体实现 UserService 接口
type UserServiceImpl struct {
	db          *gorm.DB
	idGenerator utils.IDGenerator // 添加 ID 生成器接口
}

// NewUserService 创建 UserService 实例
func NewUserService(db *gorm.DB, idGen utils.IDGenerator) UserService {
	return &UserServiceImpl{db: db, idGenerator: idGen} // 使用传入的生成器
}

func (s *UserServiceImpl) CreateUser(mobile string) (*models.User, error) {
	// 1. 验证手机号的格式（您可以添加更复杂的验证逻辑）
	if len(mobile) == 0 {
		return nil, errors.New("手机号码不能为空")
	}

	// 2. 生成 CreatorID
	creatorID, err := utils.GenerateUniqueID()
	if err != nil {
		return nil, err
	}
	now := time.Now().Truncate(time.Second) // 获取当前时间并截断到秒

	// 3. 创建新用户
	user := &models.User{
		Mobile:    mobile,
		CreatorID: creatorID, // 根据需求生成 CreatorID
		CreatedAt: models.JSONTime{Time: now},
		UpdatedAt: models.JSONTime{Time: now},
	}

	// 4. 保存用户到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// 5. 返回新创建的用户信息
	return user, nil
}

// QueryMobileIsExist 检查手机号码是否已存在
func (s *UserServiceImpl) QueryMobileIsExist(mobile string) (bool, error) {
	var user models.User
	// 查询数据库，查找手机号是否存在
	result := s.db.Where("mobile = ?", mobile).First(&user)

	// 检查是否出错
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil // 用户不存在
		}
		return false, result.Error // 其他查询错误
	}

	// 用户存在
	return true, nil
}

// GetUserByMobile 根据手机号查询用户信息
func (s *UserServiceImpl) GetUserByMobile(mobile string) (*models.User, error) {
	var user models.User

	// 使用 GORM 根据手机号查询用户
	result := s.db.Where("mobile = ?", mobile).First(&user)

	// 检查查询是否出错
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 如果未找到用户，返回 nil
		}
		return nil, result.Error // 返回其他查询错误
	}

	// 返回用户信息
	return &user, nil
}

func (s *UserServiceImpl) GetUserByCreatorID(CreatorID string) (*models.User, error) {
	var user models.User

	// 使用 GORM 根据 CreatorID 查询用户
	result := s.db.Where("creator_id = ?", CreatorID).First(&user)

	// 检查查询是否出错
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 如果未找到用户，返回 nil
		}
		return nil, result.Error // 返回其他查询错误
	}

	// 输出调试信息
	log.Printf("查询到的用户信息: %+v", user)

	// 返回用户信息的指针
	return &user, nil
}
