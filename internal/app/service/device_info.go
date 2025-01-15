package service

import (
	"time"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/errors"
)

// DeviceInfoService 设备信息服务
type DeviceInfoService struct{}

// NewDeviceInfoService 创建设备信息服务实例
func NewDeviceInfoService() *DeviceInfoService {
	return &DeviceInfoService{}
}

// DeviceInfoRequest 设备信息上报请求
type DeviceInfoRequest struct {
	DeviceInfo struct {
		DeviceSN        string `json:"device_sn" binding:"required"`
		DeviceModel     string `json:"device_model" binding:"required"`
		HardwareVersion string `json:"hardware_version"`
	} `json:"device_info" binding:"required"`
	SoftwareVersions struct {
		Klipper        string `json:"klipper" binding:"required"`
		KlipperScreen  string `json:"klipper_screen" binding:"required"`
		Moonraker      string `json:"moonraker"`
		Mainsail       string `json:"mainsail"`
		Crowsnest      string `json:"crowsnest"`
		Firmware       struct {
			Mainboard string `json:"mainboard" binding:"required"`
			Printhead string `json:"printhead" binding:"required"`
			Leveling  string `json:"leveling"`
		} `json:"firmware" binding:"required"`
	} `json:"software_versions" binding:"required"`
}

// ReportDeviceInfo 上报设备信息
func (s *DeviceInfoService) ReportDeviceInfo(req *DeviceInfoRequest) error {
	// 开启事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新设备基础信息
	deviceInfo := &model.DeviceInfo{
		DeviceSN:        req.DeviceInfo.DeviceSN,
		DeviceModel:     req.DeviceInfo.DeviceModel,
		HardwareVersion: req.DeviceInfo.HardwareVersion,
		UpdateTime:      time.Now(),
	}

	// 使用 Upsert 操作
	if err := tx.Where("device_sn = ?", deviceInfo.DeviceSN).
		Assign(deviceInfo).
		FirstOrCreate(deviceInfo).Error; err != nil {
		tx.Rollback()
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	// 2. 检查并更新软件版本信息
	softwareVersions := &model.SoftwareVersions{
		DeviceSN:            req.DeviceInfo.DeviceSN,
		KlipperVersion:      req.SoftwareVersions.Klipper,
		KlipperScreenVersion: req.SoftwareVersions.KlipperScreen,
		MoonrakerVersion:    req.SoftwareVersions.Moonraker,
		MainsailVersion:     req.SoftwareVersions.Mainsail,
		CrowsnestVersion:    req.SoftwareVersions.Crowsnest,
		MainboardFirmware:   req.SoftwareVersions.Firmware.Mainboard,
		PrintheadFirmware:   req.SoftwareVersions.Firmware.Printhead,
		LevelingFirmware:    req.SoftwareVersions.Firmware.Leveling,
		ReportTime:         time.Now(),
	}

	// 查询是否存在记录
	var existingVersions model.SoftwareVersions
	err := tx.Where("device_sn = ?", softwareVersions.DeviceSN).First(&existingVersions).Error
	if err != nil {
		if database.DB.Error.Error() == "record not found" {
			// 不存在记录，直接插入
			if err := tx.Create(softwareVersions).Error; err != nil {
				tx.Rollback()
				return errors.NewWithError(errors.ErrDatabase, err)
			}
		} else {
			// 其他数据库错误
			tx.Rollback()
			return errors.NewWithError(errors.ErrDatabase, err)
		}
	} else {
		// 存在记录，检查版本是否有变化
		if existingVersions.KlipperVersion != softwareVersions.KlipperVersion ||
			existingVersions.KlipperScreenVersion != softwareVersions.KlipperScreenVersion ||
			existingVersions.MoonrakerVersion != softwareVersions.MoonrakerVersion ||
			existingVersions.MainsailVersion != softwareVersions.MainsailVersion ||
			existingVersions.CrowsnestVersion != softwareVersions.CrowsnestVersion ||
			existingVersions.MainboardFirmware != softwareVersions.MainboardFirmware ||
			existingVersions.PrintheadFirmware != softwareVersions.PrintheadFirmware ||
			existingVersions.LevelingFirmware != softwareVersions.LevelingFirmware {
			// 版本有变化，更新记录
			if err := tx.Model(&model.SoftwareVersions{}).
				Where("device_sn = ?", softwareVersions.DeviceSN).
				Updates(map[string]interface{}{
					"klipper_version":       softwareVersions.KlipperVersion,
					"klipper_screen_version": softwareVersions.KlipperScreenVersion,
					"moonraker_version":     softwareVersions.MoonrakerVersion,
					"mainsail_version":      softwareVersions.MainsailVersion,
					"crowsnest_version":     softwareVersions.CrowsnestVersion,
					"mainboard_firmware":    softwareVersions.MainboardFirmware,
					"printhead_firmware":    softwareVersions.PrintheadFirmware,
					"leveling_firmware":     softwareVersions.LevelingFirmware,
					"report_time":          softwareVersions.ReportTime,
				}).Error; err != nil {
				tx.Rollback()
				return errors.NewWithError(errors.ErrDatabase, err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	return nil
} 