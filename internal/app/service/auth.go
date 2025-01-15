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
	// 创建JWT Claims
	claims := jwt.MapClaims{
		"device_id": device.ID,
		"sn":        device.SN,
		"exp":       time.Now().Add(24 * time.Hour).Unix(), // 24小时有效期
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("generate token error: %v", err)
	}

	// 保存token记录
	deviceToken := mdmodel.DeviceToken{
		DeviceID: device.ID,
		Token:    tokenString,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}
	if err := database.DB.Create(&deviceToken).Error; err != nil {
		return "", fmt.Errorf("save token error: %v", err)
	}

	return tokenString, nil
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