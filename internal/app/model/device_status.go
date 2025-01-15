package model

import (
	"time"
)

// DeviceStatus 设备状态信息
type DeviceStatus struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"id"`
	DeviceSN        string    `gorm:"column:device_sn" json:"device_sn"`
	StorageTotal    int64     `gorm:"column:storage_total" json:"storage_total"`
	StorageUsed     int64     `gorm:"column:storage_used" json:"storage_used"`
	StorageFree     int64     `gorm:"column:storage_free" json:"storage_free"`
	CPUUsage        float64   `gorm:"column:cpu_usage" json:"cpu_usage"`
	CPUTemperature  float64   `gorm:"column:cpu_temperature" json:"cpu_temperature"`
	MemoryTotal     int64     `gorm:"column:memory_total" json:"memory_total"`
	MemoryUsed      int64     `gorm:"column:memory_used" json:"memory_used"`
	MemoryFree      int64     `gorm:"column:memory_free" json:"memory_free"`
	ReportTime      time.Time `gorm:"column:report_time" json:"report_time"`
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 表名
func (DeviceStatus) TableName() string {
	return "md_device_status"
}

// DeviceOnline 设备在线状态
type DeviceOnline struct {
	ID             int64      `gorm:"primaryKey;column:id" json:"id"`
	DeviceSN       string     `gorm:"column:device_sn" json:"device_sn"`
	IsOnline       bool       `gorm:"column:is_online" json:"is_online"`
	LastReportTime time.Time  `gorm:"column:last_report_time" json:"last_report_time"`
	OfflineTime    *time.Time `gorm:"column:offline_time" json:"offline_time"`
	CreateTime     time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime     time.Time  `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 表名
func (DeviceOnline) TableName() string {
	return "md_device_online"
} 