package utils

import (
	"hash/crc32"
	"short_link_generation/storage"
	"strings"
)

// 生成短链接
func GenerateShortCode(url string) string {
	//把长连接转为短链接:CRC32哈希+Base62编码,取到前6位
	//1.CRC32哈希
	checksum := crc32.ChecksumIEEE([]byte(url))

	//2.转换为base62
	base62Str := ToBase62(uint64(checksum))

	//3.补齐长度
	if len(base62Str) < 6 {
		base62Str += strings.Repeat("0", 6-len(base62Str))
	}

	//4.返回短链接
	return base62Str[:6]
}

// 冲突检测函数
func IsCollection(store *storage.MemoryStore, code string) bool {
	store.Mu.RLock()
	defer store.Mu.RUnlock()
	_, exists := store.UrlMap[code]
	return exists
}

// 验证是否是合法的短码
func IsValidCode(code string) bool {
	//长度为6
	if len(code) != 6 {
		return false
	}

	//在短码中只允许出现字母和数字
	for _, c := range code {
		if !isBase62Char(c) {
			return false
		}
	}

	return true
}

// 验证是否符合base62编码
func isBase62Char(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
