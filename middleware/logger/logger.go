package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.记录请求开始时间
		start := time.Now()
		//2.处理请求前
		c.Next()

		//3.处理 请求后
		latency := time.Since(start)
		status := c.Writer.Status()

		//4.结构化日志输出
		logFields := map[string]interface{}{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  status,
			"latency": latency,
			"client":  c.ClientIP(),
		}

		if status >= 500 {
			//gin.DefaultErrorWriter 是一个全局的 io.Writer，通常用于写入错误日志。
			//Write 方法接受一个字节切片作为参数，并将其内容写入到 gin.DefaultErrorWriter 指向的输出流中。
			gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("[ERROR]%+v\n", logFields)))
		} else {
			fmt.Printf("[INFO]%+v\n", logFields)
		}
	}
}
