package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ValidateSign 验证签名
func ValidateSign(sn, secret string, timestamp int64, sign string) bool {
	// 签名规则：sha256(sn + secret + timestamp)
	data := fmt.Sprintf("%s%s%d", sn, secret, timestamp)
	hash := sha256.Sum256([]byte(data))
	expectedSign := hex.EncodeToString(hash[:])
	return sign == expectedSign
} 