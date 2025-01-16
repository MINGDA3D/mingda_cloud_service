package model

import (
	"time"
	"gorm.io/gorm"
)

// PrintTask 打印任务模型
type PrintTask struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID          string     `gorm:"column:task_id;type:varchar(64);uniqueIndex" json:"task_id"`
	DeviceSN        string     `gorm:"column:device_sn;type:varchar(64);index" json:"device_sn"`
	FileName        string     `gorm:"column:file_name;type:varchar(255)" json:"file_name"`
	StartTime       *time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime         *time.Time `gorm:"column:end_time" json:"end_time"`
	Status          string     `gorm:"column:status;type:varchar(32);index" json:"status"`
	Progress        float64    `gorm:"column:progress;type:decimal(5,2);default:0" json:"progress"`
	Duration        int        `gorm:"column:duration;default:0" json:"duration"`
	FilamentUsed    float64    `gorm:"column:filament_used;type:decimal(10,2);default:0" json:"filament_used"`
	LayersCompleted int        `gorm:"column:layers_completed;default:0" json:"layers_completed"`
	ErrorCode       string     `gorm:"column:error_code;type:varchar(32)" json:"error_code"`
	CancelReason    string     `gorm:"column:cancel_reason;type:varchar(255)" json:"cancel_reason"`
	CreateTime      time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time  `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (PrintTask) TableName() string {
	return "md_print_tasks"
}

// PrintTaskHistory 打印任务状态变更记录模型
type PrintTaskHistory struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID         string    `gorm:"column:task_id;type:varchar(64);index" json:"task_id"`
	DeviceSN       string    `gorm:"column:device_sn;type:varchar(64);index" json:"device_sn"`
	PreviousStatus string    `gorm:"column:previous_status;type:varchar(32)" json:"previous_status"`
	CurrentStatus  string    `gorm:"column:current_status;type:varchar(32)" json:"current_status"`
	ChangeTime     time.Time `gorm:"column:change_time" json:"change_time"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (PrintTaskHistory) TableName() string {
	return "md_print_task_history"
}

// 打印任务状态常量
const (
	TaskStatusIdle      = "idle"      // 空闲状态
	TaskStatusPrinting  = "printing"  // 打印中
	TaskStatusPaused    = "paused"    // 已暂停
	TaskStatusResumed   = "resumed"   // 已恢复
	TaskStatusCompleted = "completed" // 已完成
	TaskStatusCancelled = "cancelled" // 已取消
	TaskStatusError     = "error"     // 错误中断
) 