package config

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var RedisClient *redis.Client
var Ctx = context.Background()
var DB *gorm.DB

// LoadConfig 加载配置文件
func LoadConfig() error {
	viper.SetConfigName("config") // 配置文件名称 (不带扩展名)
	viper.SetConfigType("yaml")   // 配置文件格式
	viper.AddConfigPath(".")      // 当前目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 确保配置已成功加载
	log.Println("配置文件加载成功")

	return nil
}

// 初始化 Redis 连接
func InitRedis() {
	// 从配置文件中获取 Redis 配置
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	// 检查是否正确加载了 Redis 配置
	if host == "" || port == 0 {
		log.Fatalf("Redis 配置加载失败，host: %s, port: %d", host, port)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("正在连接 Redis: %s", addr)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // Redis 密码
		DB:       db,       // 使用的 Redis 数据库
	})

	// 尝试 ping Redis 以确保连接成功
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("连接 Redis 失败: %v", err))
	}

	fmt.Println("Redis 连接成功")
}

// 获取 Redis 客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}

// 初始化 MySQL 连接
func InitMySQL() {
	// 从配置文件中获取 MySQL 配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
		viper.GetString("mysql.charset"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接 MySQL 数据库失败: %v", err))
	}

	fmt.Println("MySQL 连接成功")
}

// RabbitMQ 配置结构体
type RabbitMQConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	VirtualHost string `mapstructure:"virtual-host"`
}

// 全局 RabbitMQ 连接
var RabbitMQConn *amqp.Connection

// 初始化 RabbitMQ 连接
func InitRabbitMQ() {
	var rabbitConfig RabbitMQConfig

	// 从配置文件中加载 RabbitMQ 配置
	if err := viper.UnmarshalKey("rabbitmq", &rabbitConfig); err != nil {
		log.Fatalf("读取 RabbitMQ 配置失败: %v", err)
	}

	// 拼接 RabbitMQ 的 DSN
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		rabbitConfig.Username,
		rabbitConfig.Password,
		rabbitConfig.Host,
		rabbitConfig.Port,
		rabbitConfig.VirtualHost,
	)

	// 输出 DSN 进行调试
	fmt.Printf("RabbitMQ DSN: %s\n", dsn)

	// 尝试连接 RabbitMQ
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}

	// 如果连接成功，赋值给全局变量
	RabbitMQConn = conn

	// 输出成功的日志
	log.Println("RabbitMQ 连接成功")
}

// 关闭 RabbitMQ 连接
func CloseRabbitMQ() {
	if RabbitMQConn != nil {
		if err := RabbitMQConn.Close(); err != nil {
			log.Printf("关闭 RabbitMQ 连接失败: %v", err)
		} else {
			log.Println("RabbitMQ 连接已关闭")
		}
	}
}

// SetDefaultTimezone 设置全局时区为 Asia/Shanghai
func SetDefaultTimezone() error {
	loc, err := time.LoadLocation("Asia/Shanghai") // 设置为上海时区
	if err != nil {
		return err
	}
	time.Local = loc // 设置全局时区
	log.Printf("时区已设置为: %s", loc.String())
	return nil
}
