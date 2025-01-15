package migrations

import (
	"mingda_cloud_service/internal/app/model"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Device{},
		&model.DeviceToken{},
		&model.DeviceInfo{},
		&model.SoftwareVersions{},
	)
} 