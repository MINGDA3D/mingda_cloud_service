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
		return errors.New(errors.ErrDatabase, "更新设备信息失败")
	}

	// 2. 记录软件版本信息
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

	if err := tx.Create(softwareVersions).Error; err != nil {
		tx.Rollback()
		return errors.New(errors.ErrDatabase, "记录软件版本信息失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return errors.New(errors.ErrDatabase, "提交事务失败")
	}

	return nil
} 