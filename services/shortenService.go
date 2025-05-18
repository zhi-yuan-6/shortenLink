package services

import (
	"fmt"
	"gorm.io/gorm"
	"shortenLink/dto"
	"shortenLink/models"
	"shortenLink/utils"
	"sync"
	"time"
)

var urlLocks = &sync.Map{}

func Shorten(url string) (string, error) {
	//获取当前URL对应的锁
	lock, _ := urlLocks.LoadOrStore(url, &sync.Mutex{})
	mu := lock.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	//幂等性检查
	/*store.Mu.RLock()
	if code, exists := store.ReverseMap[req.URL]; exists {
		store.Mu.RUnlock()
		c.JSON(http.StatusOK, gin.H{
			"short_link": buildShortURL(c, code),
		})
		return
	}
	store.Mu.RUnlock()*/
	code, err := models.GetShortenCode(url) //需修改为本地，redis，数据库三级查询
	if err == nil {
		return code, nil
	}
	if err != gorm.ErrRecordNotFound {
		return "", err
	}

	for i := 0; i < 3; i++ {
		//最多尝试三次
		salt := time.Now().UnixNano()
		code = utils.GenerateShortCode(fmt.Sprintf("%s,%d", url, salt))
		conflict := utils.IsCollection(code)
		if !conflict { //没冲突
			break
		}
	}

	//存储短链接
	/*store.Mu.Lock()
	store.UrlMap[code] = url
	store.ReverseMap[url] = code
	store.CreatedTime[code] = time.Now()
	store.Mu.Unlock()*/
	//开启事务
	tx := dto.DB.Begin()
	defer func() { //回滚
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	expiresAt := time.Now().Add(24 * time.Hour)
	shortCode := models.ShortUrl{
		ShortCode:   code,
		OriginalUrl: url,
		ExpiresAt:   &expiresAt,
	}
	if err := tx.Create(&shortCode).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Create(&models.VisitStats{ShortCode: code}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		return "", err
	}
	/*err = shortCode.CreateShortenUrl()
	if err != nil {
		return "", err
	}
	//创建visit记录
	err = models.CreateVisitStats(code)
	if err != nil {
		return "", err
	}*/
	return code, nil
}
