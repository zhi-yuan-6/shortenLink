package utils

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	url := "https://www.example.com"
	code1 := GenerateShortCode(url)
	code2 := GenerateShortCode(url)

	if code1 != code2 {
		t.Errorf("相同URL生成的短码不一致: %s != %s", code1, code2)
	}
	if len(code1) != 6 {
		t.Errorf("短码长度不符合要求，期望6位，实际得到%d位", len(code1))
	}
}
