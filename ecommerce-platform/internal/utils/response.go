package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// JSON 返回JSON响应
func JSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessMessage 成功消息响应
func SuccessMessage(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, Response{
		Code:    200,
		Message: message,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(w http.ResponseWriter, data interface{}, total, page, size int) {
	JSON(w, http.StatusOK, PageResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求错误
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// Unauthorized 未授权
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

// Forbidden 禁止访问
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

// NotFound 未找到
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalError 内部错误
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken 生成随机Token
func GenerateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ParseInt 解析整数
func ParseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseInt64 解析int64
func ParseInt64(s string, defaultVal int64) int64 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseFloat64 解析float64
func ParseFloat64(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetFirstImage 获取第一张图片
func GetFirstImage(images string) string {
	if images == "" {
		return ""
	}
	imgs := strings.Split(images, ",")
	if len(imgs) > 0 {
		return imgs[0]
	}
	return ""
}
