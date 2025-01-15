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
	"mingda_cloud_service/internal/pkg/redis"
	"sync"
)

// DeviceLock 设备锁管理
type DeviceLock struct {
	locks sync.Map
}

func NewDeviceLock() *DeviceLock {
	return &DeviceLock{
		locks: sync.Map{},
	}
}

// Lock 获取设备锁
func (dl *DeviceLock) Lock(deviceID uint) bool {
	mutex := &sync.Mutex{}
	if _, loaded := dl.locks.LoadOrStore(deviceID, mutex); !loaded {
		mutex.Lock()
		return true
	}
	return false
}

// Unlock 释放设备锁
func (dl *DeviceLock) Unlock(deviceID uint) {
	if mutex, ok := dl.locks.Load(deviceID); ok {
		mutex.(*sync.Mutex).Unlock()
		dl.locks.Delete(deviceID)
	}
}

type AuthService struct {
	jwtSecret string
	aesKey    []byte
	deviceLock *DeviceLock
}

func NewAuthService(jwtSecret, aesKey string) *AuthService {
	return &AuthService{
		jwtSecret:  jwtSecret,
		aesKey:     []byte(aesKey),
		deviceLock: NewDeviceLock(),
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

// AuthenticateDevice 设备认证（使用分布式锁）
func (s *AuthService) AuthenticateDevice(sn, sign string, timestamp int64) (*mdmodel.Device, error) {
	ctx := context.Background()
	lockKey := fmt.Sprintf("device_lock:%s", sn)
	
	// 尝试获取分布式锁
	if !redis.Lock(ctx, lockKey, 10*time.Second) {
		return nil, errors.New(errors.ErrTooManyReq, "设备正忙")
	}
	defer redis.Unlock(ctx, lockKey)

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取设备信息（加锁）
	var device mdmodel.Device
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("sn = ?", sn).First(&device).Error; err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrDeviceNotFound, "设备未注册")
	}

	// 获取设备锁
	if !s.deviceLock.Lock(device.ID) {
		tx.Rollback()
		return nil, errors.New(errors.ErrTooManyReq, "设备正忙")
	}
	defer s.deviceLock.Unlock(device.ID)

	// 检查时间戳
	if time.Now().Unix()-timestamp > 300 { // 5分钟内有效
		tx.Rollback()
		return nil, errors.New(errors.ErrExpired, "请求已过期")
	}

	// 验证签名
	if !utils.ValidateSign(sn, device.Secret, timestamp, sign) {
		tx.Rollback()
		return nil, errors.New(errors.ErrInvalidSign, "签名验证失败")
	}

	// 更新设备状态
	if device.Status == 0 {
		if err := tx.Model(&device).Update("status", 1).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("activate device error: %v", err)
		}
		device.Status = 1
	}

	// 更新最后在线时间
	now := time.Now()
	if err := tx.Model(&device).Update("last_online", now).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update last_online error: %v", err)
	}
	device.LastOnline = now

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit transaction error: %v", err)
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

		// 检查设备状态
		if err := s.checkDeviceStatus(&device); err != nil {
			return nil, err
		}

		return &device, nil
	}

	return nil, errors.New(errors.ErrUnauthorized, "无效的访问令牌")
}

// checkDeviceStatus 检查设备状态
func (s *AuthService) checkDeviceStatus(device *mdmodel.Device) error {
	if device.Status != 1 {
		return errors.New(errors.ErrUnauthorized, "device is not activated")
	}
	return nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(oldToken string) (string, error) {
	// 解析旧token
	claims, err := utils.ParseToken(oldToken, s.jwtSecret)
	if err != nil {
		return "", errors.New(errors.ErrUnauthorized, "invalid token")
	}

	// 获取设备锁
	if !s.deviceLock.Lock(claims.DeviceID) {
		return "", errors.New(errors.ErrTooManyReq, "设备正忙")
	}
	defer s.deviceLock.Unlock(claims.DeviceID)

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查token是否在黑名单中
	if s.isTokenBlacklisted(oldToken) {
		tx.Rollback()
		return "", errors.New(errors.ErrUnauthorized, "token has been revoked")
	}

	// 获取设备信息
	var device mdmodel.Device
	if err := tx.First(&device, claims.DeviceID).Error; err != nil {
		tx.Rollback()
		return "", errors.New(errors.ErrDeviceNotFound, "device not found")
	}

	// 检查设备状态
	if err := s.checkDeviceStatus(&device); err != nil {
		tx.Rollback()
		return "", err
	}

	// 生成新token
	newToken, err := utils.GenerateToken(&device, s.jwtSecret, 24*time.Hour)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// 将旧token加入黑名单
	if err := s.addToBlacklist(oldToken, claims.ExpiresAt); err != nil {
		tx.Rollback()
		return "", err
	}

	// 保存新token记录
	deviceToken := &mdmodel.DeviceToken{
		DeviceID: device.ID,
		Token:    newToken,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	if err := tx.Create(deviceToken).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("commit transaction error: %v", err)
	}

	return newToken, nil
}

// isTokenBlacklisted 检查token是否在黑名单中
func (s *AuthService) isTokenBlacklisted(token string) bool {
	exists, err := redis.Exists(context.Background(), fmt.Sprintf("token_blacklist:%s", token))
	return err == nil && exists
}

// addToBlacklist 将token加入黑名单
func (s *AuthService) addToBlacklist(token string, expireAt int64) error {
	ctx := context.Background()
	duration := time.Until(time.Unix(expireAt, 0))
	return redis.Set(ctx, fmt.Sprintf("token_blacklist:%s", token), true, duration)
} 