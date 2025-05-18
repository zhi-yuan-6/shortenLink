package redirect

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortenLink/models"
	"shortenLink/utils"
	"sync"
)

var codeLocks = &sync.Map{}

func RedirectHandler(c *gin.Context) {
	//获取参数
	code := c.Param("code")
	if !utils.IsValidCode(code) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid short code"})
		return
	}

	//获取短码对应的长链接
	/*//store.Mu.RLock()
	//originalURL, exists := store.UrlMap[code]
	//store.Mu.RUnlock()*/
	//调用service
	originalURL, err := utils.GetOriginalURL(code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "short code not found"})
		return
	}

	//原子操作递增计数器(已经改为在查询时递增)
	//store.IncrementVisit(code)
	lock, _ := codeLocks.LoadOrStore(code, &sync.Mutex{})
	mu := lock.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	err = models.IncrementVisit(code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"original_url": originalURL,
	})
	//c.Redirect(http.StatusFound, originalURL)
}
