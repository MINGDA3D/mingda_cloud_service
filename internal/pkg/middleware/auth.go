package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/errors"
	"mingda_cloud_service/internal/pkg/utils"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, errors.New(errors.ErrInvalidToken, "缺少认证头"))
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, errors.New(errors.ErrInvalidToken, "认证格式错误"))
			return
		}

		// 验证token
		claims, err := validateToken(parts[1])
		if err != nil {
			if err.Error() == "token has expired" {
				c.AbortWithStatusJSON(401, errors.New(errors.ErrTokenExpired, "访问令牌已过期"))
			} else {
				c.AbortWithStatusJSON(401, errors.New(errors.ErrInvalidToken, err.Error()))
			}
			return
		}

		// 验证设备状态
		var device model.Device
		if err := database.DB.First(&device, claims.DeviceID).Error; err != nil {
			c.AbortWithStatusJSON(401, errors.New(errors.ErrDeviceNotFound, "设备不存在"))
			return
		}

		// 检查设备状态
		if device.Status != 1 { // 假设1表示正常状态
			c.AbortWithStatusJSON(401, errors.New(errors.ErrDeviceDisabled, "设备已禁用"))
			return
		}

		// 更新最后在线时间
		database.DB.Model(&device).Update("last_online", time.Now())

		// 将设备信息存储到上下文
		c.Set("device_id", claims.DeviceID)
		c.Set("device_sn", claims.DeviceSN)
		c.Set("device", device)

		c.Next()
	}
}

func validateToken(tokenString string) (*utils.Claims, error) {
	// 从配置获取JWT密钥
	jwtSecret := "mingda3D250113PrintingCloudService2024" // 建议从配置中获取

	// 解析token
	claims, err := utils.ParseToken(tokenString, jwtSecret)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
