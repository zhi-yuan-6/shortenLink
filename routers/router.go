package routers

import (
	"github.com/gin-gonic/gin"
	"shortenLink/handler/health"
	"shortenLink/handler/redirect"
	"shortenLink/handler/shorten"
	"shortenLink/handler/stats"
	"shortenLink/middleware"
	"shortenLink/middleware/logger"
	"shortenLink/middleware/recovery"
)

func SetupRouters(r *gin.Engine) {
	r.Use(logger.RequestLogger(), recovery.CustomRecovery(), middleware.NewConcurrencyLimiter(500).Limit())

	r.GET("/:code", redirect.RedirectHandler)
	api := r.Group("/api")
	{
		api.GET("/health", health.HealthHandler)
		api.POST("/shorten", shorten.ShortenHandler)
		api.GET("/stats/:code", stats.StatsHandler)
	}

}
