package service

import (
	"time"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/errors"
	"gorm.io/gorm"
)

// DeviceStatusService 设备状态服务
type DeviceStatusService struct {
	offlineCheckTicker *time.Ticker
	stopChan          chan struct{}
}

// NewDeviceStatusService 创建设备状态服务实例
func NewDeviceStatusService() *DeviceStatusService {
	service := &DeviceStatusService{
		offlineCheckTicker: time.NewTicker(1 * time.Minute), // 每分钟检查一次
		stopChan:          make(chan struct{}),
	}
	
	// 启动离线检测任务
	go service.startOfflineCheck()
	
	return service
}

// startOfflineCheck 启动离线检测任务
func (s *DeviceStatusService) startOfflineCheck() {
	for {
		select {
		case <-s.offlineCheckTicker.C:
			if err := s.CheckDeviceOnlineStatus(); err != nil {
				// TODO: 添加日志记录
				continue
			}
		case <-s.stopChan:
			s.offlineCheckTicker.Stop()
			return
		}
	}
}

// Stop 停止服务
func (s *DeviceStatusService) Stop() {
	close(s.stopChan)
}

// DeviceStatusRequest 设备状态上报请求
type DeviceStatusRequest struct {
	DeviceSN       string  `json:"device_sn"`                              // 设备SN，以token中的值为准
	StorageTotal   int64   `json:"storage_total" binding:"required"`
	StorageUsed    int64   `json:"storage_used" binding:"required"`
	StorageFree    int64   `json:"storage_free" binding:"required"`
	CPUUsage       float64 `json:"cpu_usage" binding:"required"`
	CPUTemperature float64 `json:"cpu_temperature" binding:"required"`
	MemoryTotal    int64   `json:"memory_total" binding:"required"`
	MemoryUsed     int64   `json:"memory_used" binding:"required"`
	MemoryFree     int64   `json:"memory_free" binding:"required"`
}

// ReportDeviceStatus 上报设备状态
func (s *DeviceStatusService) ReportDeviceStatus(tokenDeviceSN string, req *DeviceStatusRequest) error {
	// 使用token中的设备SN
	req.DeviceSN = tokenDeviceSN

	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 记录设备状态
	status := &model.DeviceStatus{
		DeviceSN:       req.DeviceSN,
		StorageTotal:   req.StorageTotal,
		StorageUsed:    req.StorageUsed,
		StorageFree:    req.StorageFree,
		CPUUsage:       req.CPUUsage,
		CPUTemperature: req.CPUTemperature,
		MemoryTotal:    req.MemoryTotal,
		MemoryUsed:     req.MemoryUsed,
		MemoryFree:     req.MemoryFree,
		ReportTime:     time.Now(),
	}

	if err := tx.Create(status).Error; err != nil {
		tx.Rollback()
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	// 2. 更新设备在线状态
	var online model.DeviceOnline
	err := tx.Where("device_sn = ?", req.DeviceSN).First(&online).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在记录，创建新记录
			now := time.Now()
			online = model.DeviceOnline{
				DeviceSN:       req.DeviceSN,
				IsOnline:       true,
				LastReportTime: now,
				CreateTime:     now,
				UpdateTime:     now,
			}
			if err := tx.Create(&online).Error; err != nil {
				tx.Rollback()
				return errors.NewWithError(errors.ErrDatabase, err)
			}
		} else {
			// 其他数据库错误
			tx.Rollback()
			return errors.NewWithError(errors.ErrDatabase, err)
		}
	} else {
		// 更新在线状态
		now := time.Now()
		if err := tx.Model(&online).Updates(map[string]interface{}{
			"is_online":        true,
			"last_report_time": now,
			"offline_time":     gorm.Expr("NULL"), // 设备上线时，清除离线时间
			"update_time":      now,
		}).Error; err != nil {
			tx.Rollback()
			return errors.NewWithError(errors.ErrDatabase, err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	return nil
}

// CheckDeviceOnlineStatus 检查设备在线状态
func (s *DeviceStatusService) CheckDeviceOnlineStatus() error {
	// 获取当前时间
	now := time.Now()
	offlineThreshold := now.Add(-10 * time.Minute) // 10分钟未上报则判定为离线

	// 先查询是否有需要更新的设备
	var count int64
	if err := database.DB.Model(&model.DeviceOnline{}).
		Where("is_online = ? AND last_report_time < ?", true, offlineThreshold).
		Count(&count).Error; err != nil {
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	// 只有在有需要更新的设备时才执行更新操作
	if count > 0 {
		if err := database.DB.Model(&model.DeviceOnline{}).
			Where("is_online = ? AND last_report_time < ?", true, offlineThreshold).
			Updates(map[string]interface{}{
				"is_online":    false,
				"offline_time": now,
				"update_time": now,
			}).Error; err != nil {
			return errors.NewWithError(errors.ErrDatabase, err)
		}
	}

	return nil
} 