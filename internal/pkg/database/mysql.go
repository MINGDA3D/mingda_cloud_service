package database

import (
    "fmt"
    "mingda_cloud_service/internal/pkg/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func NewMySQLConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
    )
    
    return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
