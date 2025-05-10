package services

import (
	"shortenLink/storage"
	"shortenLink/utils"
	"time"
)

func Shorten(url string, store *storage.MemoryStore) string {
	var code string
	for i := 0; i < 3; i++ { //最多尝试三次
		code = utils.GenerateShortCode(url)
		if !utils.IsCollection(store, code) {
			break
		}
		//处理冲突：追加随机字符 time.Now().UnixNano():获取当前时间的纳秒级时间戳，返回值是一个整数，使用计算出的索引从 utils.Base62Chars 字符串中获取一个字符。
		code += string(utils.Base62Chars[time.Now().UnixNano()%62])
		if !utils.IsCollection(store, code) {
			break
		}
	}

	//存储短链接
	store.Mu.Lock()
	store.UrlMap[code] = url
	store.ReverseMap[url] = code
	store.CreatedTime[code] = time.Now()
	store.Mu.Unlock()

	return code
}
