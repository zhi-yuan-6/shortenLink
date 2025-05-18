package services

import (
	"shortenLink/models"
	"sync"
	"time"
)

func StatsService(code string) (int64, time.Time, error) {
	//获取原始URL用于存在性校验
	/*store.Mu.RLock()
	_, exists := store.UrlMap[code]
	store.Mu.RUnlock()*/
	_, err := models.GetOriginalURL(code)

	if err != nil {
		return 0, time.Time{}, err
	}

	//原始读取统计数据
	/*record, ok := store.Visits.Load(code)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"visits": 0,
		})
		return
	}

	//获取拜访次数和时间
	visits := atomic.LoadInt64(&record.(*memory.VisitRecord).Count)
	lastVisit := time.Unix(0, atomic.LoadInt64(&record.(*memory.VisitRecord).LastVisit))*/
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()
	visits, lastVisit, err := models.GetVisitCount(code)
	if err != nil {
		return 0, time.Time{}, err
	}
	return visits, lastVisit, nil
}
