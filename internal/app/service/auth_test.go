package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
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