package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length/2)
	rand.Read(b)
	return hex.EncodeToString(b)
} 