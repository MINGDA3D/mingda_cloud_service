package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"mingda_cloud_service/internal/app/model"
)

// Claims 自定义JWT声明
type Claims struct {
	DeviceID uint   `json:"device_id"`
	DeviceSN string `json:"device_sn"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT token
func GenerateToken(device *model.Device, secret string, expireDuration time.Duration) (string, error) {
	// 设置claims
	claims := Claims{
		DeviceID: device.ID,
		DeviceSN: device.SN,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "mingda-cloud",
		},
	}

	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
} 