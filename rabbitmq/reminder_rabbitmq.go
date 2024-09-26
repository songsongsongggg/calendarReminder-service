package rabbitmq

import (
	"calendarReminder-service/config"
	"calendarReminder-service/models"
	"calendarReminder-service/utils"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

var (
	conn *amqp.Connection
	once sync.Once
)

// RabbitMQ 配置
const (
	exchangeName = "reminder_exchange"    // 交换机名称
	queueName    = "reminder_queue"       // 队列名称
	routingKey   = "reminder.routing_key" // 路由键
)

// setupRabbitMQ 设置交换机和队列，并进行绑定
func SetupRabbitMQ() error {
	log.Println("开始设置 RabbitMQ")

	// 创建一个新的信道
	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		log.Printf("RabbitMQ 信道创建失败: %v", err)
		return err
	}
	defer ch.Close()
	log.Println("信道创建成功")

	// 声明延迟消息交换机
	err = ch.ExchangeDeclare(
		exchangeName,        // 交换机名称
		"x-delayed-message", // 交换机类型
		true,                // 持久化
		false,               // 自动删除
		false,               // 独占
		false,               // 阻塞
		amqp.Table{
			"x-delayed-type": "topic", // 使用主题交换机作为基础
		},
	)
	if err != nil {
		log.Printf("交换机声明失败: %v", err)
		return err
	}
	log.Printf("交换机 [%s] 声明成功", exchangeName)

	// 声明队列
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("队列声明失败: %v", err)
		return err
	}
	log.Printf("队列 [%s] 声明成功", queueName)

	// 将队列绑定到交换机
	err = ch.QueueBind(
		queueName,    // 队列名称
		"reminder.*", // 路由键模式
		exchangeName, // 交换机名称
		false,
		nil,
	)
	if err != nil {
		log.Printf("队列绑定失败: %v", err)
		return err
	}
	log.Printf("队列 [%s] 与交换机 [%s] 绑定成功，路由键模式 [%s]", queueName, exchangeName, "reminder.*")

	return nil
}

// PublishReminderToQueue 发布消息到队列
func PublishReminderToQueue(content string, delay int64, mobile string) error {
	log.Printf("开始发布消息: 内容=%s, 手机号=%s, 延迟=%d 毫秒", content, mobile, delay)

	// 继续使用北京时间
	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		log.Printf("RabbitMQ 信道创建失败: %v", err)
		return err
	}
	defer ch.Close()

	reminderMsg := models.ReminderMessage{
		Content: content,
		Mobile:  mobile,
	}

	body, err := json.Marshal(reminderMsg)
	if err != nil {
		log.Printf("消息体序列化失败: %v", err)
		return err
	}
	log.Printf("时间间隔: %v", delay)

	err = ch.Publish(
		exchangeName, // 交换机名称
		routingKey,   // 路由键
		false,        // 是否强制发布
		false,        // 是否强制发送给消费者
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers: amqp.Table{
				"x-delay": delay, // 延迟时间为基于北京时间的毫秒数
			},
		})
	if err != nil {
		log.Printf("消息发布失败: %v", err)
		return err
	}

	log.Printf("消息发布成功，内容: %s，手机号: %s，延迟 %d 毫秒发送", content, mobile, delay)
	return nil
}

// ConsumeReminders 消费队列中的消息
func ConsumeReminders() error {
	log.Println("准备消费消息")

	// 使用全局的 RabbitMQ 连接
	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		log.Printf("无法创建 RabbitMQ 信道: %v", err)
		return err
	}
	defer ch.Close()
	log.Println("信道创建成功")

	// 开始消费队列中的消息
	msgs, err := ch.Consume(
		queueName, // 队列名称
		"",        // 消费者标识符（空则自动生成）
		true,      // 是否自动应答
		false,     // 是否独占
		false,     // 是否阻塞
		false,     // 额外属性
		nil,
	)
	if err != nil {
		log.Printf("无法消费消息: %v", err)
		return err
	}

	log.Println("开始消费消息")

	// 处理消息
	for msg := range msgs {
		log.Printf("接收到消息: %s", string(msg.Body))

		var reminderMsg models.ReminderMessage
		// 反序列化消息体
		if err := json.Unmarshal(msg.Body, &reminderMsg); err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}

		// 日志输出解析的消息内容
		log.Printf("解析成功，内容: %s，手机号: %s", reminderMsg.Content, reminderMsg.Mobile)

		// 调用 utils.SendSMSReminder 发送短信提醒
		err := utils.SendSMSReminder(reminderMsg.Content, reminderMsg.Mobile)
		if err != nil {
			log.Printf("短信发送失败: %v", err)
		} else {
			log.Printf("短信发送成功，手机号: %s", reminderMsg.Mobile)
		}
	}

	return nil
}
