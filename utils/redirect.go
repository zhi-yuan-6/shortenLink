package utils

import (
	"context"
	"log"
	"shortenLink/cache"
	"shortenLink/models"
	"time"
)

var ctx = context.Background()

// 获取原始url
func GetOriginalURL(code string) (string, error) {
	//1.查本地缓存
	if url, ok := cache.RedisCacheIns.Local.Get(code); ok {
		log.Println("Cache hit in local cache")
		return url.(string), nil
	}
	//2.查redis
	url, err := cache.RedisCacheIns.Redis.Get(ctx, code).Result()
	if err == nil {
		cache.RedisCacheIns.Local.Set(code, url, 0)
		log.Println("Cache hit in local redis")
		return url, nil
	}
	//3.查数据库
	originalURL, err := models.GetOriginalURL(code)
	if err != nil {
		return "", err
	}

	//4.回填缓存
	go func() {
		cache.RedisCacheIns.Redis.Set(ctx, code, originalURL, 24*time.Hour)
		cache.RedisCacheIns.Local.Set(code, originalURL, 0)
	}()

	return originalURL, nil
}

/*// 删除url
func DeleteURL(code string) error {
	// 1. 先删缓存
	cache.RedisCacheIns.Local.Del(code)

	// 2. 再删redis
	if err := cache.RedisCacheIns.Redis.Del(ctx, code).Err(); err != nil {
		return err
	}

	// 3. 更新数据库
	if err := models.DeleteShortenURL(code); err != nil {
		return err
	}

	return nil
}
*/
