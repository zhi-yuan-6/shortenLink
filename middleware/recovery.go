package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// CustomRecovery 是一个自定义的中间件，用于处理 panic 错误，并返回 500 错误响应。   注：带学习
func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"trace": string(debug.Stack()), //debug.Stack() 获取当前的栈跟踪信息，返回一个字节切片。
		})
	})
}
