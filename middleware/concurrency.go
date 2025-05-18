package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConcurrencyLimiter struct {
	sem chan struct{}
}

// NewConcurrencyLimiter 创建一个新的 ConcurrencyLimiter 实例。 注：带学习
func NewConcurrencyLimiter(max int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		sem: make(chan struct{}, max),
	}
}

func (cl *ConcurrencyLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case cl.sem <- struct{}{}:
			defer func() { <-cl.sem }()
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		}

	}
}
