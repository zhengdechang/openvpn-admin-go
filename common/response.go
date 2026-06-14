package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// OK 返回成功响应（带数据）
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

// OKMsg 返回成功响应（带消息，无数据）
func OKMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{Success: true, Message: message})
}

// Fail 返回失败响应
func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, Response{Success: false, Error: message})
}

// BadRequest 400
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, message)
}

// Unauthorized 401
func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, message)
}

// Forbidden 403
func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, message)
}

// NotFound 404
func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, message)
}

// InternalError 500
func InternalError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, message)
}
