package shorten

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortenLink/services"
)

// ShortLinkHandler 处理短链接请求,POST请求
type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func ShortenHandler(c *gin.Context) {
	// 获取请求参数url
	//url := c.Query("url")
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid format:" + err.Error(),
		})
		return
	}

	//生成短码（带冲突检测）
	code, err := services.Shorten(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成短码失败:" + err.Error(),
		})
		return
	}
	// 返回短链接
	c.JSON(http.StatusOK, gin.H{
		"short_link": buildShortURL(c, code),
	})
}

func buildShortURL(c *gin.Context, code string) string {
	return c.Request.Host + "/" + code
}
