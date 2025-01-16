package model

import (
    "time"
)

// PrintImage 打印图片模型
type PrintImage struct {
    ID          uint      `gorm:"primarykey"`
    TaskID      string    `gorm:"type:varchar(32);not null;comment:打印任务ID"`
    DeviceSN    string    `gorm:"type:varchar(32);not null;comment:设备序列号"`
    ImagePath   string    `gorm:"type:varchar(255);not null;comment:图片存储路径"`
    ImageURL    string    `gorm:"type:varchar(255);not null;comment:图片访问URL"`
    Status      int       `gorm:"type:tinyint;not null;default:0;comment:检测状态(0-未检测,1-检测中,2-检测完成)"`
    HasDefect   *bool     `gorm:"comment:是否存在缺陷"`
    DefectType  string    `gorm:"type:varchar(32);comment:缺陷类型"`
    Confidence  float64   `gorm:"type:decimal(5,4);comment:检测置信度"`
    CreateTime  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
    UpdateTime  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 表名
func (p *PrintImage) TableName() string {
    return "md_print_images"
}

// 检测状态常量
const (
    StatusPending  = 0 // 未检测
    StatusChecking = 1 // 检测中
    StatusChecked  = 2 // 检测完成
) 