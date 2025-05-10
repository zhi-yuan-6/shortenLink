package routers

import (
	"github.com/gin-gonic/gin"
	"shortenLink/handler/health"
	"shortenLink/handler/redirect"
	"shortenLink/handler/shorten"
	"shortenLink/handler/stats"
	"shortenLink/middleware"
	"shortenLink/storage"
)

func SetupRouters(r *gin.Engine, store *storage.MemoryStore) {
	r.Use(middleware.RequestLogger(), middleware.CustomRecovery(), middleware.NewConcurrencyLimiter(500).Limit())

	r.GET("/:code", redirect.RedirectHandler(store))
	api := r.Group("/api")
	{
		api.GET("/health", health.HealthHandler)
		api.POST("/shorten", shorten.ShortenHandler(store))
		api.GET("/stats/:code", stats.StatsHandler(store))
	}

}
