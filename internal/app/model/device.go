package model

import (
	"time"
	"gorm.io/gorm"
)

// Device 设备模型
type Device struct {
	gorm.Model
	SN          string    `gorm:"type:varchar(32);uniqueIndex;not null" json:"sn"`           // 设备序列号
	DeviceModel string    `gorm:"type:varchar(32);not null" json:"model"`                    // 设备型号
	Name        string    `gorm:"type:varchar(64)" json:"name"`                              // 设备名称
	Secret      string    `gorm:"type:varchar(64);not null" json:"secret"`                        // 设备密钥
	Status      int       `gorm:"type:tinyint;default:0" json:"status"`                      // 设备状态
	LastOnline  time.Time `gorm:"type:datetime;not null" json:"last_online"`                 // 最后在线时间
	FirmwareVer string    `gorm:"type:varchar(32)" json:"firmware_ver"`                      // 固件版本
	IP          string    `gorm:"type:varchar(64)" json:"ip"`                                // IP地址
	MAC         string    `gorm:"type:varchar(32)" json:"mac"`                               // MAC地址
}

// DeviceToken 设备令牌模型
type DeviceToken struct {
	gorm.Model
	DeviceID uint      `gorm:"index;not null" json:"device_id"`                // 设备ID
	Token    string    `gorm:"type:varchar(256);not null" json:"token"`        // 访问令牌
	ExpireAt time.Time `gorm:"not null" json:"expire_at"`                      // 过期时间
}

// TableName 指定表名
func (Device) TableName() string {
	return "md_devices"
}

func (DeviceToken) TableName() string {
	return "md_device_tokens"
} 