package utils

import "regexp"

// 国内手机号码的正则表达式
const chinaPhoneNumberRegex = `^1[3-9]\d{9}$`

// 校验手机号是否符合国内的标准
func IsValidPhoneNumber(mobile string) bool {
	// 如果手机号为空，则返回false
	if mobile == "" {
		return false
	}

	// 使用正则表达式进行校验
	pattern := regexp.MustCompile(chinaPhoneNumberRegex)
	return pattern.MatchString(mobile)
}
