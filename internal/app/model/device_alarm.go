package model

import (
	"time"
)

// DeviceAlarm 设备告警信息
type DeviceAlarm struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	DeviceSN    string    `gorm:"column:device_sn;index" json:"device_sn"`           // 设备SN
	AlarmType   int       `gorm:"column:alarm_type" json:"alarm_type"`               // 告警类型：1-存储空间不足，2-CPU温度过高，3-内存不足
	AlarmLevel  int       `gorm:"column:alarm_level" json:"alarm_level"`             // 告警级别：1-提示，2-警告，3-严重
	AlarmValue  float64   `gorm:"column:alarm_value" json:"alarm_value"`             // 告警时的具体数值
	AlarmDesc   string    `gorm:"column:alarm_desc" json:"alarm_desc"`               // 告警描述
	Status      int       `gorm:"column:status" json:"status"`                       // 状态：0-未处理，1-已处理，2-已忽略
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	ResolveTime *time.Time `gorm:"column:resolve_time" json:"resolve_time"`          // 告警解除时间
	ResolveDesc string    `gorm:"column:resolve_desc" json:"resolve_desc"`           // 告警解除说明
}

// TableName 表名
func (DeviceAlarm) TableName() string {
	return "md_device_alarm"
}

// AlarmTypeStorage 存储空间不足
const AlarmTypeStorage = 1

// AlarmTypeCPUTemp CPU温度过高
const AlarmTypeCPUTemp = 2

// AlarmTypeMemory 内存不足
const AlarmTypeMemory = 3

// AlarmLevelInfo 提示
const AlarmLevelInfo = 1

// AlarmLevelWarning 警告
const AlarmLevelWarning = 2

// AlarmLevelCritical 严重
const AlarmLevelCritical = 3

// AlarmStatusPending 未处理
const AlarmStatusPending = 0

// AlarmStatusResolved 已处理
const AlarmStatusResolved = 1

// AlarmStatusIgnored 已忽略
const AlarmStatusIgnored = 2 