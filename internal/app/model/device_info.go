package model

import (
	"time"
	"gorm.io/gorm"
)

// DeviceInfo 设备基础信息
type DeviceInfo struct {
	gorm.Model
	DeviceSN        string    `gorm:"column:device_sn;type:varchar(64);uniqueIndex;not null" json:"device_sn"`         // 设备SN码
	DeviceModel     string    `gorm:"column:device_model;type:varchar(32);not null" json:"device_model"`               // 机型
	HardwareVersion string    `gorm:"column:hardware_version;type:varchar(32)" json:"hardware_version"`                // 硬件版本号
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;not null" json:"update_time"`                    // 更新时间
}

// SoftwareVersions 软件版本信息
type SoftwareVersions struct {
	gorm.Model
	DeviceSN            string    `gorm:"column:device_sn;type:varchar(64);index;not null" json:"device_sn"`              // 设备SN码
	KlipperVersion      string    `gorm:"column:klipper_version;type:varchar(32)" json:"klipper_version"`                 // Klipper版本
	KlipperScreenVersion string    `gorm:"column:klipper_screen_version;type:varchar(32)" json:"klipper_screen_version"`   // KlipperScreen版本
	MoonrakerVersion   string    `gorm:"column:moonraker_version;type:varchar(32)" json:"moonraker_version"`             // Moonraker版本
	MainsailVersion    string    `gorm:"column:mainsail_version;type:varchar(32)" json:"mainsail_version"`               // Mainsail版本
	CrowsnestVersion   string    `gorm:"column:crowsnest_version;type:varchar(32)" json:"crowsnest_version"`             // Crowsnest版本
	MainboardFirmware  string    `gorm:"column:mainboard_firmware;type:varchar(32)" json:"mainboard_firmware"`           // 主板固件版本
	PrintheadFirmware  string    `gorm:"column:printhead_firmware;type:varchar(32)" json:"printhead_firmware"`           // 打印头板固件版本
	LevelingFirmware   string    `gorm:"column:leveling_firmware;type:varchar(32)" json:"leveling_firmware"`             // 快速调平板固件版本
	ReportTime        time.Time `gorm:"column:report_time;type:datetime;not null" json:"report_time"`                    // 上报时间
}

// TableName 指定DeviceInfo表名
func (DeviceInfo) TableName() string {
	return "md_device_info"
}

// TableName 指定SoftwareVersions表名
func (SoftwareVersions) TableName() string {
	return "md_software_versions"
} 