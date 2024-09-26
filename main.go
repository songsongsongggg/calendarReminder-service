package main

import (
	"calendarReminder-service/config"
	"calendarReminder-service/models"
	"calendarReminder-service/rabbitmq"
	"calendarReminder-service/routes"
	"calendarReminder-service/services"
	"calendarReminder-service/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// 在main函数中，修复路由注册

func main() {
	// 加载配置文件
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	// 设置默认时区为 Asia/Shanghai
	if err := config.SetDefaultTimezone(); err != nil {
		log.Fatalf("设置时区失败: %v", err)
	}

	// 初始化 Redis、MySQL 和 RabbitMQ
	config.InitRedis()
	config.InitMySQL()
	// 初始化 RabbitMQ
	config.InitRabbitMQ()

	// 初始化 RabbitMQ
	if err := rabbitmq.SetupRabbitMQ(); err != nil {
		log.Fatalf("RabbitMQ 初始化失败: %v", err)
	}

	// 启动消息消费
	go func() {
		if err := rabbitmq.ConsumeReminders(); err != nil {
			log.Fatalf("消费消息失败: %v", err)
		}
	}()

	// 自动迁移表结构
	config.DB.AutoMigrate(&models.User{}, &models.Reminder{})

	// 初始化 ID 生成器和 UserService
	idGen := &utils.SimpleIDGenerator{}
	userService := services.NewUserService(config.DB, idGen)
	reminderService := services.NewReminderService(config.DB)

	// 初始化路由
	router := mux.NewRouter()

	// 注册用户登录、登出和短信验证码的路由，传递router
	routes.PassportRoutes(router, userService)
	// 注册提醒功能的路由
	routes.ReminderRoutes(router, userService, reminderService)

	// 启动服务
	log.Println("服务启动在端口 :9900")
	if err := http.ListenAndServe(":9900", router); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
