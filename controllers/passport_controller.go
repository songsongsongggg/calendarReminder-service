package controllers

import (
	"calendarReminder-service/config"
	"calendarReminder-service/models"
	"calendarReminder-service/services"
	"calendarReminder-service/utils"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Redis中存储短信验证码的Key前缀
const MOBILE_SMSCODE = "MOBILE_SMSCODE:"

// Redis中存储用户Token的Key前缀
const REDIS_USER_TOKEN = "REDIS_USER_TOKEN:"

// 发送短信验证码接口
func GetSMSCode(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	// 从请求中提取手机号码
	mobile := r.URL.Query().Get("mobile")
	log.Printf("收到发送验证码请求, 手机号: %s", mobile)

	// 校验手机号码的格式（可以使用正则表达式）
	if !utils.IsValidPhoneNumber(mobile) {
		log.Printf("手机号格式不正确: %s", mobile)
		utils.ErrorResponse(w, http.StatusBadRequest, "手机号格式不正确")
		return
	}

	// 根据用户的IP地址限制请求频率，防止刷短信
	userIP := utils.GetRequestIP(r)
	ipKey := fmt.Sprintf(MOBILE_SMSCODE + userIP)

	// 使用Redis的SetNX方法设置一个带过期时间的键，确保60秒内只能发送一次验证码
	_, err := redisClient.SetNX(config.Ctx, ipKey, mobile, 60*time.Second).Result()
	if err != nil {
		log.Printf("请求频繁，60秒内只能请求一次验证码, IP: %s, 手机号: %s", userIP, mobile)
		utils.ErrorResponse(w, http.StatusTooManyRequests, "60秒内只能请求一次验证码")
		return
	}

	// 生成6位随机验证码
	rand.Seed(time.Now().UnixNano())
	random := fmt.Sprintf("%06d", rand.Intn(1000000))
	log.Printf("生成验证码: %s, 手机号: %s", random, mobile)

	// 将验证码存储到Redis，有效期30分钟
	redisKey := fmt.Sprintf(MOBILE_SMSCODE + mobile)
	err = redisClient.Set(config.Ctx, redisKey, random, 30*time.Minute).Err()
	if err != nil {
		log.Printf("存储验证码失败: %s", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "存储验证码失败")
		return
	}

	// 调用第三方服务发送短信验证码
	err = utils.SendSMS(mobile, random)
	if err != nil {
		log.Printf("发送短信失败: %s, 手机号: %s", err, mobile)
		utils.ErrorResponse(w, http.StatusInternalServerError, "发送短信失败")
		return
	}

	log.Printf("验证码已发送至手机号: %s", mobile)
	utils.SuccessResponse(w, nil, "验证码已发送成功")
}

// 登录或注册处理
func Login(w http.ResponseWriter, r *http.Request, userService services.UserService, redisClient *redis.Client) {
	// 从请求中提取手机号码和验证码
	var reqBody struct {
		Mobile  string `json:"mobile"`
		SmsCode string `json:"smsCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("请求体解析失败: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "请求体解析失败")
		return
	}
	log.Printf("收到登录/注册请求, 手机号: %s, 验证码: %s", reqBody.Mobile, reqBody.SmsCode)

	// 校验验证码是否正确
	redisKey := fmt.Sprintf(MOBILE_SMSCODE + reqBody.Mobile)
	redisSmsCode, err := redisClient.Get(config.Ctx, redisKey).Result()
	if err != nil || redisSmsCode != reqBody.SmsCode {
		log.Printf("验证码错误, 手机号: %s, 输入的验证码: %s, Redis中的验证码: %s", reqBody.Mobile, reqBody.SmsCode, redisSmsCode)
		utils.ErrorResponse(w, http.StatusUnauthorized, "验证码错误")
		return
	}

	// 查询用户是否存在
	userExists, err := userService.QueryMobileIsExist(reqBody.Mobile)
	if err != nil {
		log.Printf("查询用户是否存在失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("查询用户失败: %v", err))
		return
	}

	var user *models.User // 定义用户变量
	if userExists {
		// 如果用户存在，则获取用户信息
		user, err = userService.GetUserByMobile(reqBody.Mobile)
		if err != nil {
			log.Printf("获取用户信息失败: %v", err)
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取用户信息失败: %v", err))
			return
		}
		log.Printf("用户存在，获取用户信息成功, 用户ID: %s", user.CreatorID)
	} else {
		// 用户不存在，创建新用户
		user, err = userService.CreateUser(reqBody.Mobile)
		if err != nil {
			log.Printf("注册用户失败: %v", err)
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("注册用户失败: %v", err))
			return
		}
		log.Printf("新用户注册成功, 用户ID: %s", user.CreatorID)
	}

	// 生成分布式Token，并存储在Redis中
	token := utils.GenerateUUID() // 生成唯一的 UUID token
	tokenKey := fmt.Sprintf(REDIS_USER_TOKEN + user.CreatorID)

	if err = redisClient.Set(config.Ctx, tokenKey, token, 24*time.Hour).Err(); err != nil {
		log.Printf("存储 Token 失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("生成 Token 失败: %v", err))
		return
	}
	log.Printf("Token 生成并存储成功, 用户ID: %s, Token: %s", user.CreatorID, token)

	// 设置cookies
	setCookie(w, "token", token, 30*24*time.Hour)
	setCookie(w, "creator_id", user.CreatorID, 30*24*time.Hour)
	log.Printf("Cookies 设置成功, 用户ID: %s", user.CreatorID)

	// 删除已使用的验证码
	if err := redisClient.Del(config.Ctx, redisKey).Err(); err != nil {
		log.Printf("删除验证码失败: %v", err)
	}

	// 返回成功响应
	utils.SuccessResponse(w, nil, "登录成功")
}

// 设置Cookie的帮助函数
func setCookie(w http.ResponseWriter, name, value string, maxAge time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true, // 设置为 HttpOnly 防止 JavaScript 访问
	})
	log.Printf("Cookie 设置成功: %s = %s", name, value)
}

// 用户退出登录
func Logout(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	// 从请求中提取用户ID
	var reqBody struct {
		CreatorID string `json:"creator_id"` // 将 userID 字段设为字符串类型
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("请求体解析失败: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "请求体解析失败")
		return
	}

	log.Printf("用户ID是: %s", reqBody.CreatorID)

	// 删除用户的登录Token
	tokenKey := fmt.Sprintf(REDIS_USER_TOKEN + reqBody.CreatorID)
	err := redisClient.Del(config.Ctx, tokenKey).Err()
	if err != nil {
		log.Printf("删除Token失败: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("删除Token失败: %v", err))
		return
	}

	log.Printf("用户登出成功, 用户ID: %s", reqBody.CreatorID)
	utils.SuccessResponse(w, nil, "用户登出成功")
}
