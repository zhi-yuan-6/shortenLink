package stats

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_generation/storage"
	"short_link_generation/utils"
	"sync/atomic"
	"time"
)

func StatsHandler(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取参数
		code := c.Param("code")
		if !utils.IsValidCode(code) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid short code",
			})
		}

		//获取原始URL用于存在性校验
		store.Mu.RLock()
		_, exists := store.UrlMap[code]
		store.Mu.RUnlock()

		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		//原始读取统计数据
		record, ok := store.Visits.Load(code)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"visits": 0,
			})
			return
		}

		//获取拜访次数和时间
		visits := atomic.LoadInt64(&record.(*storage.VisitRecord).Count)
		lastVisit := time.Unix(0, atomic.LoadInt64(&record.(*storage.VisitRecord).LastVisit))

		c.JSON(http.StatusOK, gin.H{
			"visits":     visits,
			"last_visit": lastVisit.Format(time.RFC3339),
		})
	}
}
