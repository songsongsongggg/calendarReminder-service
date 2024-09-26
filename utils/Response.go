package utils

import (
	"encoding/json"
	"net/http"
)

// Response 统一的响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功的响应
func SuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	response := Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ErrorResponse 错误的响应
func ErrorResponse(w http.ResponseWriter, code int, message string) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
