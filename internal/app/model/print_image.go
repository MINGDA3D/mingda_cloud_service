package model

import (
    "time"
)

// PrintImage 打印图片记录
type PrintImage struct {
    ID          int64     `gorm:"primaryKey;column:id" json:"id"`
    DeviceSN    string    `gorm:"column:device_sn;type:varchar(64);index" json:"device_sn"`         // 设备SN
    TaskID      string    `gorm:"column:task_id;type:varchar(64);index" json:"task_id"`             // 打印任务ID
    ImagePath   string    `gorm:"column:image_path;type:varchar(255)" json:"image_path"`            // 图片存储路径
    ImageURL    string    `gorm:"column:image_url;type:varchar(255)" json:"image_url"`              // 图片访问URL
    Status      int       `gorm:"column:status;type:tinyint" json:"status"`                         // 检测状态：0-未检测 1-检测中 2-检测完成
    HasDefect   bool      `gorm:"column:has_defect;type:tinyint" json:"has_defect"`                // 是否存在缺陷
    DefectType  string    `gorm:"column:defect_type;type:varchar(64)" json:"defect_type"`          // 缺陷类型
    Confidence  float64   `gorm:"column:confidence;type:decimal(5,2)" json:"confidence"`            // 检测置信度
    CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`             // 创建时间
    UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`             // 更新时间
}

// TableName 表名
func (PrintImage) TableName() string {
    return "md_print_images"
} 