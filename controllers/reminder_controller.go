package controllers

import (
	"calendarReminder-service/models"
	"calendarReminder-service/rabbitmq"
	"calendarReminder-service/services"
	"calendarReminder-service/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// 创建提醒
func CreateReminder(w http.ResponseWriter, r *http.Request, userService services.UserService, reminderService services.ReminderService) {
	log.Println("开始处理创建提醒的请求")

	var reminder models.Reminder
	// 解析请求体
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		log.Printf("解析请求体失败: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "请求体解析失败")
		return
	}

	// 获取当前时间
	now := time.Now().Add(time.Hour).Truncate(time.Second) // 获取当前时间并截断到秒

	// 设置提醒的创建和更新时间
	reminder.CreatedAt = models.JSONTime{Time: now}
	reminder.UpdatedAt = models.JSONTime{Time: now}

	err := reminderService.CreateReminder(&reminder)
	if err != nil {
		log.Printf("创建提醒失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "创建提醒失败")
		return
	}

	// 获取时间字符串（假设没有时区）
	remindAtStr := reminder.RemindAt.String() // 这里获取的是一个带时区的字符串
	remindAtStrWithoutTZ := remindAtStr[:19]  // 去掉时区信息

	// 假设提醒时间通过请求体传递，解析前端传入的时间字符串
	remindAt, err := time.Parse("2006-01-02 15:04:05", remindAtStrWithoutTZ) // 使用自定义的 JSONTime 类型
	if err != nil {
		log.Printf("时间解析失败: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "时间格式不正确")
		return
	}
	now = time.Now().Add(8 * time.Hour).Truncate(time.Second) // 获取当前时间并截断到秒

	// 检查提醒时间是否在未来
	if remindAt.Before(now) {
		log.Printf("提醒时间无效: %v", remindAt)
		utils.ErrorResponse(w, http.StatusBadRequest, "提醒时间必须是未来的时间")
		return
	}

	// 检查提醒时间是否至少距离当前时间一小时
	//if remindAt.Sub(now) < time.Hour {
	//	log.Printf("提醒时间距离当前时间少于一小时: %v", remindAt)
	//	utils.ErrorResponse(w, "提醒时间必须至少在一个小时之后", http.StatusBadRequest)
	//	return
	//}

	// 计算提醒时间与当前时间的延迟
	delay := remindAt.Sub(now).Milliseconds()

	// 获取用户信息
	user, err := userService.GetUserByCreatorID(reminder.CreatorID)
	if err != nil || user == nil {
		log.Printf("获取用户信息失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 发布提醒消息到 RabbitMQ 延迟队列
	err = rabbitmq.PublishReminderToQueue(reminder.Content, delay, user.Mobile)
	if err != nil {
		log.Printf("发布消息到队列失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "创建提醒成功，但短信提醒无法发送")
		return
	}

	log.Println("提醒创建成功，消息已发布到队列")
	utils.SuccessResponse(w, nil, "提醒创建成功")
}

// 获取用户的提醒列表
func GetReminders(w http.ResponseWriter, r *http.Request, reminderService services.ReminderService) {
	// 日志记录：开始处理获取提醒列表的请求
	log.Println("开始处理获取提醒列表的请求")

	// 从请求的 cookie 中获取 creator_id
	creatorID, err := GetCreatorIDFromRequest(r)
	if err != nil {
		log.Printf("获取 creator_id 失败: %v", err)
		utils.ErrorResponse(w, http.StatusUnauthorized, "未授权的请求")
		return
	}
	// 日志记录：获取提醒的创建者ID
	log.Printf("获取创建者ID: %s 的提醒列表", creatorID)

	// 调用服务层获取提醒列表
	reminders, err := reminderService.GetRemindersByCreatorID(creatorID)
	if err != nil {
		log.Printf("获取提醒失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "获取提醒失败")
		return
	}

	// 日志记录：成功获取提醒列表
	log.Printf("提醒列表获取成功: %+v", reminders)

	// 返回提醒列表
	utils.SuccessResponse(w, reminders, "获取提醒列表成功")
}

// 删除提醒
func DeleteReminder(w http.ResponseWriter, r *http.Request, reminderService services.ReminderService) {
	// 日志记录：开始处理删除提醒的请求
	log.Println("开始处理删除提醒的请求")

	// 从路径参数中获取要删除的提醒ID
	id, err := getIDFromRequest(r)
	if err != nil {
		log.Println("提醒ID不存在")
		utils.ErrorResponse(w, http.StatusBadRequest, "提醒ID不存在")
		return
	}

	// 从请求的 cookie 中获取 creator_id
	creatorID, err := GetCreatorIDFromRequest(r)
	if err != nil {
		log.Printf("获取 creator_id 失败: %v", err)
		utils.ErrorResponse(w, http.StatusUnauthorized, "未授权的请求")
		return
	}

	// 日志记录：尝试删除提醒
	log.Printf("尝试删除提醒, ID: %s, 创建者ID: %s", id, creatorID)

	// 调用服务层删除提醒
	if err := reminderService.DeleteReminder(id, creatorID); err != nil {
		log.Printf("删除提醒失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "删除提醒失败")
		return
	}

	// 删除成功，返回无内容状态
	utils.SuccessResponse(w, nil, "删除提醒成功")
}

// 更新提醒
func UpdateReminder(w http.ResponseWriter, r *http.Request, reminderService services.ReminderService) {
	// 日志记录：开始处理更新提醒的请求
	log.Println("开始处理更新提醒的请求")

	// 从路径参数中获取要更新的提醒ID
	id, err := getIDFromRequest(r)
	if err != nil {
		log.Println("提醒ID不存在")
		utils.ErrorResponse(w, http.StatusBadRequest, "提醒ID不存在")
		return
	}

	// 定义提醒对象，解析请求体
	var reminder models.Reminder
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		log.Printf("解析请求体失败: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "请求体解析失败")
		return
	}

	// 从请求的 cookie 中获取 creator_id
	creatorID, err := GetCreatorIDFromRequest(r)
	if err != nil {
		log.Printf("获取 creator_id 失败: %v", err)
		utils.ErrorResponse(w, http.StatusUnauthorized, "未授权的请求")
		return
	}

	// 日志记录：尝试更新提醒
	log.Printf("尝试更新提醒, ID: %s, 创建者ID: %s, 更新内容: %+v", id, creatorID, reminder)

	// 调用服务层更新提醒
	if err := reminderService.UpdateReminder(id, &reminder, creatorID); err != nil {
		log.Printf("更新提醒失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "更新提醒失败")
		return
	}

	// 日志记录：成功更新提醒
	log.Println("提醒更新成功")

	// 返回更新后的提醒信息
	utils.SuccessResponse(w, nil, "提醒更新成功")
}

// 从请求中获取 creator_id
func GetCreatorIDFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("creator_id") // 获取存储 creator_id 的 cookie
	if err != nil {
		return "", err // 如果没有找到 cookie，返回错误
	}
	return cookie.Value, nil // 返回 cookie 的值
}

// 从请求中获取提醒ID的工具函数
func getIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return "", fmt.Errorf("提醒ID不存在")
	}
	return id, nil
}
