package routers

import (
	"github.com/gin-gonic/gin"
	"short_link_generation/controllers/health"
	"short_link_generation/controllers/redirect"
	"short_link_generation/controllers/shortLink"
	"short_link_generation/controllers/stats"
	"short_link_generation/middleware"
	"short_link_generation/storage"
)

func SetupRouters(r *gin.Engine, store *storage.MemoryStore) {
	r.Use(middleware.RequestLogger(), middleware.CustomRecovery(), middleware.NewConcurrencyLimiter(500).Limit())

	r.GET("/:code", redirect.RedirectHandler(store))
	api := r.Group("/api")
	{
		api.GET("/health", health.HealthHandler)
		api.POST("/shorten", shortLink.ShortenHandler(store))
		api.GET("/stats/:code", stats.StatsHandler(store))
	}

}
