package routes

import (
	"calendarReminder-service/config"
	"calendarReminder-service/controllers"
	"calendarReminder-service/services"
	"net/http"

	"github.com/gorilla/mux"
)

func PassportRoutes(r *mux.Router, userService services.UserService) {

	// 发送短信验证码接口
	r.HandleFunc("/getSMSCode", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetSMSCode(w, r, config.RedisClient) // 直接使用 config.RedisClient
	}).Methods("GET")

	// 登录接口
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.Login(w, r, userService, config.RedisClient) // 直接使用 config.RedisClient
	}).Methods("POST")

	// 登出接口
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		controllers.Logout(w, r, config.RedisClient) // 直接使用 config.RedisClient
	}).Methods("POST")

}

func ReminderRoutes(r *mux.Router, userService services.UserService, reminderService services.ReminderService) {
	// POST 和 GET 请求的路由处理
	r.HandleFunc("/reminders", func(w http.ResponseWriter, r *http.Request) {
		// POST: 创建提醒
		if r.Method == http.MethodPost {
			controllers.CreateReminder(w, r, userService, reminderService)
		}
		// GET: 获取提醒列表
		if r.Method == http.MethodGet {
			controllers.GetReminders(w, r, reminderService)
		}
	}).Methods(http.MethodPost, http.MethodGet)

	// DELETE 和 PUT 请求的路由处理
	r.HandleFunc("/reminders/{id}", func(w http.ResponseWriter, r *http.Request) {
		// DELETE: 删除提醒
		if r.Method == http.MethodDelete {
			controllers.DeleteReminder(w, r, reminderService)
		}
		// PUT: 更新提醒
		if r.Method == http.MethodPut {
			controllers.UpdateReminder(w, r, reminderService)
		}
	}).Methods(http.MethodDelete, http.MethodPut)
}
