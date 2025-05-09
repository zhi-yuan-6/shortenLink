package utils

import "strings"

const Base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ToBase62 将 uint64 数字转换为 Base62 字符串
func ToBase62(num uint64) string {
	if num == 0 {
		return "0"
	}
	//var result []byte
	var result strings.Builder
	for num > 0 {
		res := num % 62
		result.WriteByte(Base62Chars[res])
		//result = append(result, base62Chars[res])
		num /= 62
	}
	//反转字节 排序
	/*for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)*/
	return reverseString(result.String())
}

// reverseString 将字符串反转
func reverseString(s string) string {
	//将字符串转换为字符切片
	r := []rune(s)
	//反转字符切片
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	//将字符切片转换为字符串
	return string(r)
}
