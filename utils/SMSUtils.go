package utils

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置文件
func InitConfig() {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 添加当前目录（开发环境使用）
	viper.AddConfigPath(".")
	// 如果是测试环境，添加一个绝对路径，确保 viper 能找到配置文件 // 当前目录
	viper.AddConfigPath("/Users/wuzhisong/wuzhisong/GolandProjects/calendarReminder-service") // 项目根目录

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("致命错误：无法读取配置文件: %s", err))
	}
}

// SendSMS 发送短信验证码
func SendSMS(mobile string, random string) error {
	client, err := CreateClient()
	if err != nil {
		return fmt.Errorf("创建阿里云客户端时出错: %v", err)
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("迎客知识"),
		TemplateCode:  tea.String("SMS_461375482"),
		PhoneNumbers:  tea.String(mobile),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", random)),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return fmt.Errorf("发送短信时出错: %v", err)
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}

		// 打印完整的错误信息以便调试
		fmt.Printf("发送短信时出错: %v\n", error)
		fmt.Println("错误消息:", tea.StringValue(error.Message))

		if error.Data != nil {
			fmt.Println("错误详情:", tea.StringValue(error.Data))
		}

		return fmt.Errorf("短信发送失败: %v", error.Message)
	}
	return nil
}

// SendSMSReminder 发送通知短信
func SendSMSReminder(content string, mobile string) error {
	client, err := CreateClient()
	if err != nil {
		return fmt.Errorf("创建阿里云客户端时出错: %v", err)
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("迎客知识"),
		TemplateCode:  tea.String("SMS_473770239"),
		PhoneNumbers:  tea.String(mobile),
		TemplateParam: tea.String(fmt.Sprintf("{\"value\":\"%s\"}", content)),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return fmt.Errorf("发送短信时出错: %v", err)
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}

		// 打印完整的错误信息以便调试
		fmt.Printf("发送短信时出错: %v\n", error)
		fmt.Println("错误消息:", tea.StringValue(error.Message))

		if error.Data != nil {
			fmt.Println("错误详情:", tea.StringValue(error.Data))
		}

		return fmt.Errorf("短信发送失败: %v", error.Message)
	}
	return nil
}

// CreateClient 创建阿里云短信服务客户端，读取配置文件中的AccessKeyId和AccessKeySecret
func CreateClient() (*dysmsapi20170525.Client, error) {
	// 初始化viper配置
	InitConfig()

	// 从配置文件中读取AccessKeyId、AccessKeySecret和Endpoint
	accessKeyId := viper.GetString("alibabaCloud.accessKeyId")
	accessKeySecret := viper.GetString("alibabaCloud.accessKeySecret")
	endpoint := "dysmsapi.aliyuncs.com"

	// 校验配置文件是否成功读取
	if accessKeyId == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("配置文件中缺少必要的阿里云配置信息")
	}

	// 配置阿里云客户端
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	return dysmsapi20170525.NewClient(config)
}
