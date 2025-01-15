package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	mdmodel "mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/errors"
	"mingda_cloud_service/internal/pkg/utils"
	"mingda_cloud_service/internal/pkg/validator"
	"github.com/gin-gonic/gin"
	"context"
)

type AuthService struct {
	jwtSecret string
	aesKey    []byte
}

func NewAuthService(jwtSecret, aesKey string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		aesKey:    []byte(aesKey),
	}
}

// RegisterDevice 注册设备
func (s *AuthService) RegisterDevice(sn, model string) (*mdmodel.Device, error) {
	// 验证SN码格式
	if err := validator.ValidateDeviceSN(sn); err != nil {
		return nil, errors.New(errors.ErrInvalidSN, err.Error())
	}

	// 检查设备是否已存在
	var device mdmodel.Device
	if err := database.DB.Where("sn = ?", sn).First(&device).Error; err == nil {
		return nil, errors.New(errors.ErrDeviceTypeInvalid, "设备已注册")
	}

	// 在生产环境下，应该从生产管理系统获取预置的密钥
	var secret string
	if gin.Mode() == gin.DebugMode {
		// 开发环境：生成随机密钥
		secret = utils.GenerateRandomString(32)
	} else {
		// 生产环境：从生产管理系统获取预置密钥
		// TODO: 实现从生产管理系统获取密钥的逻辑
		// secret = productionSystem.GetDeviceSecret(sn)
		// if secret == "" {
		//     return nil, errors.New(errors.ErrDeviceNotFound, "设备未授权")
		// }
		secret = utils.GenerateRandomString(32) // 临时使用，实际生产环境需要替换
	}

	// 创建设备记录
	device = mdmodel.Device{
		SN:          sn,
		DeviceModel: model,
		Secret:      secret,
		Status:      0, // 初始状态：未激活
		LastOnline:  time.Now(), // 设置初始在线时间
	}

	if err := database.DB.Create(&device).Error; err != nil {
		return nil, fmt.Errorf("create device error: %v", err)
	}

	return &device, nil
}

// AuthenticateDevice 设备认证
func (s *AuthService) AuthenticateDevice(sn, sign string, timestamp int64) (*mdmodel.Device, error) {
	// 检查时间戳
	if time.Now().Unix()-timestamp > 300 { // 5分钟内有效
		return nil, errors.New(errors.ErrExpired, "请求已过期")
	}

	// 获取设备信息
	var device mdmodel.Device
	if err := database.DB.Where("sn = ?", sn).First(&device).Error; err != nil {
		return nil, errors.New(errors.ErrDeviceNotFound, "设备未注册")
	}

	// 验证签名
	if !utils.ValidateSign(sn, device.Secret, timestamp, sign) {
		return nil, errors.New(errors.ErrInvalidSign, "签名验证失败")
	}

	return &device, nil
}

// GenerateToken 生成访问令牌
func (s *AuthService) GenerateToken(device *mdmodel.Device) (string, error) {
	// 生成JWT token
	token, err := utils.GenerateToken(device, s.jwtSecret, 24*time.Hour) // token有效期24小时
	if err != nil {
		return "", err
	}

	// 保存token记录
	deviceToken := &mdmodel.DeviceToken{
		DeviceID: device.ID,
		Token:    token,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	if err := database.DB.Create(deviceToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken 验证访问令牌
func (s *AuthService) ValidateToken(tokenString string) (*mdmodel.Device, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.New(errors.ErrUnauthorized, "无效的访问令牌")
	}

	// 验证claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		deviceID := uint(claims["device_id"].(float64))
		
		// 获取设备信息
		var device mdmodel.Device
		if err := database.DB.First(&device, deviceID).Error; err != nil {
			return nil, errors.New(errors.ErrDeviceNotFound, "设备不存在")
		}

		return &device, nil
	}

	return nil, errors.New(errors.ErrUnauthorized, "无效的访问令牌")
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(oldToken string) (string, error) {
	// 解析旧token
	claims, err := utils.ParseToken(oldToken, s.jwtSecret)
	if err != nil {
		return "", errors.New(errors.ErrUnauthorized, "invalid token")
	}

	// 检查token是否在黑名单中
	if s.isTokenBlacklisted(oldToken) {
		return "", errors.New(errors.ErrUnauthorized, "token has been revoked")
	}

	// 获取设备信息
	var device mdmodel.Device
	if err := database.DB.First(&device, claims.DeviceID).Error; err != nil {
		return "", errors.New(errors.ErrDeviceNotFound, "device not found")
	}

	// 检查设备状态
	if device.Status != 1 {
		return "", errors.New(errors.ErrUnauthorized, "device is disabled")
	}

	// 生成新token
	newToken, err := utils.GenerateToken(&device, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return "", err
	}

	// 将旧token加入黑名单
	if err := s.addToBlacklist(oldToken, claims.ExpiresAt); err != nil {
		return "", err
	}

	// 保存新token记录
	deviceToken := &mdmodel.DeviceToken{
		DeviceID: device.ID,
		Token:    newToken,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	if err := database.DB.Create(deviceToken).Error; err != nil {
		return "", err
	}

	return newToken, nil
}

// isTokenBlacklisted 检查token是否在黑名单中
func (s *AuthService) isTokenBlacklisted(token string) bool {
	var exists bool
	err := database.Redis.Get(context.Background(), fmt.Sprintf("token_blacklist:%s", token)).Err()
	return err == nil && exists
}

// addToBlacklist 将token加入黑名单
func (s *AuthService) addToBlacklist(token string, expireAt int64) error {
	ctx := context.Background()
	duration := time.Until(time.Unix(expireAt, 0))
	return database.Redis.Set(ctx, fmt.Sprintf("token_blacklist:%s", token), true, duration).Err()
} 