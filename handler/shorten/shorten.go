package shorten

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortenLink/services"
	"shortenLink/storage"
)

// ShortLinkHandler 处理短链接请求,POST请求
type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func ShortenHandler(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求参数url
		//url := c.Query("url")
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid format:" + err.Error(),
			})
			return
		}

		//幂等性检查
		store.Mu.RLock()
		if code, exists := store.ReverseMap[req.URL]; exists {
			store.Mu.RUnlock()
			c.JSON(http.StatusOK, gin.H{
				"short_link": buildShortURL(c, code),
			})
			return
		}
		store.Mu.RUnlock()

		//生成短码（带冲突检测）
		code := services.Shorten(req.URL, store)
		// 返回短链接
		c.JSON(http.StatusOK, gin.H{
			"short_link": buildShortURL(c, code),
		})
	}
}

func buildShortURL(c *gin.Context, code string) string {
	return c.Request.Host + "/" + code
}
