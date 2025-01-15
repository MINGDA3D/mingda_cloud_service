package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/utils"
)

func TestAuthService_RefreshToken(t *testing.T) {
	// 初始化测试环境
	setupTestEnv(t)

	// 创建测试用例
	tests := []struct {
		name      string
		setupFunc func() (string, *model.Device) // 准备测试数据
		wantErr   bool
		errMsg    string
	}{
		{
			name: "正常刷新token",
			setupFunc: func() (string, *model.Device) {
				// 创建测试设备
				device := &model.Device{
					SN:          "M1A2401A0100001",
					DeviceModel: "MD-400D",
					Status:      1, // 正常状态
					LastOnline:  time.Now(),
				}
				database.DB.Create(device)

				// 为设备生成token
				authService := NewAuthService("test_secret", "test_key")
				token, _ := authService.GenerateToken(device)
				return token, device
			},
			wantErr: false,
		},
		{
			name: "使用无效token",
			setupFunc: func() (string, *model.Device) {
				return "invalid_token", nil
			},
			wantErr: true,
			errMsg:  "invalid token",
		},
		{
			name: "设备已禁用",
			setupFunc: func() (string, *model.Device) {
				// 创建已禁用的测试设备
				device := &model.Device{
					SN:          "M1A2401A0100002",
					DeviceModel: "MD-400D",
					Status:      0, // 禁用状态
					LastOnline:  time.Now(),
				}
				database.DB.Create(device)

				// 为设备生成token
				authService := NewAuthService("test_secret", "test_key")
				token, _ := authService.GenerateToken(device)
				return token, device
			},
			wantErr: true,
			errMsg:  "device is disabled",
		},
		{
			name: "token已在黑名单中",
			setupFunc: func() (string, *model.Device) {
				// 创建测试设备
				device := &model.Device{
					SN:          "M1A2401A0100003",
					DeviceModel: "MD-400D",
					Status:      1,
					LastOnline:  time.Now(),
				}
				database.DB.Create(device)

				// 为设备生成token
				authService := NewAuthService("test_secret", "test_key")
				token, _ := authService.GenerateToken(device)
				
				// 将token加入黑名单
				authService.addToBlacklist(token, time.Now().Add(time.Hour).Unix())
				
				return token, device
			},
			wantErr: true,
			errMsg:  "token has been revoked",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := NewAuthService("test_secret", "test_key")
			oldToken, _ := tt.setupFunc()

			// 执行刷新
			newToken, err := authService.RefreshToken(oldToken)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, newToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newToken)
				assert.NotEqual(t, oldToken, newToken)

				// 验证旧token是否已加入黑名单
				assert.True(t, authService.isTokenBlacklisted(oldToken))

				// 验证新token是否可用
				claims, err := authService.ValidateToken(newToken)
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestAuthService_TokenFlow(t *testing.T) {
	// 初始化测试环境
	setupTestEnv(t)

	// 创建测试设备
	device := &model.Device{
		SN:          "M1A2401A0100001",
		DeviceModel: "MD-400D",
		Secret:      "test_device_secret",
		Status:      1,
		LastOnline:  time.Now(),
	}
	database.DB.Create(device)

	t.Run("完整token流程测试", func(t *testing.T) {
		authService := NewAuthService("test_jwt_secret", "test_aes_key")

		// 1. 设备认证
		timestamp := time.Now().Unix()
		sign := utils.GenerateSign(device.SN, device.Secret, timestamp)
		authedDevice, err := authService.AuthenticateDevice(device.SN, sign, timestamp)
		assert.NoError(t, err)
		assert.NotNil(t, authedDevice)
		assert.Equal(t, device.SN, authedDevice.SN)

		// 2. 生成token
		token, err := authService.GenerateToken(authedDevice)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 3. 验证token
		claims, err := utils.ParseToken(token, "test_jwt_secret")
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, device.SN, claims.DeviceSN)

		// 4. 使用token访问受保护的资源
		validatedDevice, err := authService.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, validatedDevice)
		assert.Equal(t, device.ID, validatedDevice.ID)

		// 5. 检查设备token记录
		var deviceToken model.DeviceToken
		err = database.DB.Where("device_id = ?", device.ID).First(&deviceToken).Error
		assert.NoError(t, err)
		assert.Equal(t, token, deviceToken.Token)
		assert.True(t, deviceToken.ExpireAt.After(time.Now()))
	})

	t.Run("异常场景测试", func(t *testing.T) {
		authService := NewAuthService("test_jwt_secret", "test_aes_key")

		// 1. 使用错误的签名
		timestamp := time.Now().Unix()
		wrongSign := "wrong_sign"
		_, err := authService.AuthenticateDevice(device.SN, wrongSign, timestamp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "签名验证失败")

		// 2. 使用过期的时间戳
		expiredTimestamp := time.Now().Add(-10 * time.Minute).Unix()
		sign := utils.GenerateSign(device.SN, device.Secret, expiredTimestamp)
		_, err = authService.AuthenticateDevice(device.SN, sign, expiredTimestamp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "请求已过期")

		// 3. 使用无效的token
		invalidToken := "invalid.token.string"
		_, err = authService.ValidateToken(invalidToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的访问令牌")

		// 4. 设备被禁用后使用token
		// 先获取有效token
		validToken, _ := authService.GenerateToken(device)
		// 禁用设备
		database.DB.Model(device).Update("status", 0)
		_, err = authService.ValidateToken(validToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "设备不可用")
	})
}

// 添加辅助函数用于生成签名
func TestAuthService_GenerateSign(t *testing.T) {
	sn := "M1A2401A0100001"
	secret := "test_secret"
	timestamp := time.Now().Unix()

	sign := utils.GenerateSign(sn, secret, timestamp)
	assert.NotEmpty(t, sign)

	// 验证相同输入生成相同签名
	sign2 := utils.GenerateSign(sn, secret, timestamp)
	assert.Equal(t, sign, sign2)

	// 验证不同输入生成不同签名
	sign3 := utils.GenerateSign(sn, "different_secret", timestamp)
	assert.NotEqual(t, sign, sign3)
}

// setupTestEnv 初始化测试环境
func setupTestEnv(t *testing.T) {
	// 初始化测试数据库连接
	err := database.Init(database.Config{
		Driver:   "sqlite3",
		Database: ":memory:", // 使用内存数据库进行测试
	})
	assert.NoError(t, err)

	// 初始化测试Redis连接
	err = database.InitRedis(database.RedisConfig{
		Addr: "localhost:6379",
		DB:   1, // 使用单独的数据库进行测试
	})
	assert.NoError(t, err)

	// 清理测试数据
	database.DB.Exec("DELETE FROM md_devices")
	database.DB.Exec("DELETE FROM md_device_tokens")
	database.Redis.FlushDB()
} 