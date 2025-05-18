package routers

import (
	"github.com/gin-gonic/gin"
	"shortenLink/handler/health"
	"shortenLink/handler/redirect"
	"shortenLink/handler/shorten"
	"shortenLink/handler/stats"
	"shortenLink/middleware"
)

func SetupRouters(r *gin.Engine) {
	r.Use(middleware.RequestLogger(), middleware.CustomRecovery(), middleware.NewConcurrencyLimiter(500).Limit())

	r.GET("/:code", redirect.RedirectHandler)
	api := r.Group("/api")
	{
		api.GET("/health", health.HealthHandler)
		api.POST("/shorten", shorten.ShortenHandler)
		api.GET("/stats/:code", stats.StatsHandler)
	}

}
