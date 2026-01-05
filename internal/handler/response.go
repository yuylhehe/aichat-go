package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 标准成功响应
type SuccessResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse 标准错误响应
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// SendSuccess 发送成功响应
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Code: 0,
		Data: data,
	})
}

// SendCreated 发送创建成功响应
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Code: 0,
		Data: data,
	})
}

// SendError 发送错误响应
func SendError(c *gin.Context, httpCode int, errCode int, message string) {
	c.JSON(httpCode, ErrorResponse{
		Code:  errCode,
		Error: message,
	})
}

// SendBadRequest 发送 400 错误
func SendBadRequest(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, 400, message)
}

// SendInternalError 发送 500 错误
func SendInternalError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, 500, message)
}

// SendNotFound 发送 404 错误
func SendNotFound(c *gin.Context, message string) {
	SendError(c, http.StatusNotFound, 404, message)
}

// SendUnauthorized 发送 401 错误
func SendUnauthorized(c *gin.Context, message string) {
	SendError(c, http.StatusUnauthorized, 401, message)
}
