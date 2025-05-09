package redirect

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_generation/storage"
	"short_link_generation/utils"
)

func RedirectHandler(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取参数
		code := c.Param("code")
		if !utils.IsValidCode(code) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid short code"})
			return
		}

		//获取短码对应的长链接
		store.Mu.RLock()
		originalURL, exists := store.UrlMap[code]
		store.Mu.RUnlock()

		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		//原子操作递增计数器
		store.IncrementVisit(code)

		c.JSON(http.StatusFound, gin.H{
			"original_url": originalURL,
		})
		//c.Redirect(http.StatusFound, originalURL)
	}
}
