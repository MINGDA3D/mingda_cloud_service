package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateSign 生成签名
func GenerateSign(sn, secret string, timestamp int64) string {
	data := fmt.Sprintf("%s%s%d", sn, secret, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ValidateSign 验证签名
func ValidateSign(sn, secret string, timestamp int64, sign string) bool {
	expectedSign := GenerateSign(sn, secret, timestamp)
	return sign == expectedSign
} 