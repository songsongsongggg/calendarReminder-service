package utils

import (
	"net/http"
	"strings"
)

// GetRequestIP 获取用户请求的 IP 地址
// 处理多级代理，尝试从 X-Forwarded-For 等头部信息获取客户端的真实 IP
func GetRequestIP(r *http.Request) string {
	// 从 X-Forwarded-For 头中获取 IP
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" || strings.ToLower(ip) == "unknown" {
		// 从 Proxy-Client-IP 头中获取 IP
		ip = r.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || strings.ToLower(ip) == "unknown" {
		// 从 WL-Proxy-Client-IP 头中获取 IP
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || strings.ToLower(ip) == "unknown" {
		// 从 HTTP_CLIENT_IP 头中获取 IP
		ip = r.Header.Get("HTTP_CLIENT_IP")
	}
	if ip == "" || strings.ToLower(ip) == "unknown" {
		// 从 HTTP_X_FORWARDED_FOR 头中获取 IP
		ip = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if ip == "" || strings.ToLower(ip) == "unknown" {
		// 最后从 RemoteAddr 获取
		ip = r.RemoteAddr
	}

	// 如果 IP 是通过 X-Forwarded-For 获取的，可能是逗号分隔的多个 IP 地址
	// 取第一个有效 IP
	if strings.Contains(ip, ",") {
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	return ip
}
