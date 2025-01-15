package service

import (
	"time"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/errors"
)

// DeviceAlarmService 设备告警服务
type DeviceAlarmService struct{}

// NewDeviceAlarmService 创建设备告警服务实例
func NewDeviceAlarmService() *DeviceAlarmService {
	return &DeviceAlarmService{}
}

// DeviceAlarmRequest 设备告警上报请求
type DeviceAlarmRequest struct {
	AlarmType  int     `json:"alarm_type" binding:"required"`  // 告警类型
	AlarmLevel int     `json:"alarm_level" binding:"required"` // 告警级别
	AlarmValue float64 `json:"alarm_value" binding:"required"` // 告警值
	AlarmDesc  string  `json:"alarm_desc" binding:"required"`  // 告警描述
}

// ReportDeviceAlarm 上报设备告警
func (s *DeviceAlarmService) ReportDeviceAlarm(deviceSN string, req *DeviceAlarmRequest) error {
	// 参数校验
	if !isValidAlarmType(req.AlarmType) {
		return errors.New(errors.ErrInvalidParams, "无效的告警类型")
	}
	if !isValidAlarmLevel(req.AlarmLevel) {
		return errors.New(errors.ErrInvalidParams, "无效的告警级别")
	}

	// 创建告警记录
	alarm := &model.DeviceAlarm{
		DeviceSN:   deviceSN,
		AlarmType:  req.AlarmType,
		AlarmLevel: req.AlarmLevel,
		AlarmValue: req.AlarmValue,
		AlarmDesc:  req.AlarmDesc,
		Status:     model.AlarmStatusPending,
	}

	if err := database.DB.Create(alarm).Error; err != nil {
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	return nil
}

// ResolveAlarmRequest 处理告警请求
type ResolveAlarmRequest struct {
	ResolveDesc string `json:"resolve_desc" binding:"required"` // 处理说明
}

// ResolveAlarm 处理告警
func (s *DeviceAlarmService) ResolveAlarm(alarmID int64, req *ResolveAlarmRequest) error {
	now := time.Now()
	
	// 更新告警状态
	if err := database.DB.Model(&model.DeviceAlarm{}).
		Where("id = ? AND status = ?", alarmID, model.AlarmStatusPending).
		Updates(map[string]interface{}{
			"status":       model.AlarmStatusResolved,
			"resolve_time": now,
			"resolve_desc": req.ResolveDesc,
			"update_time":  now,
		}).Error; err != nil {
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	return nil
}

// IgnoreAlarm 忽略告警
func (s *DeviceAlarmService) IgnoreAlarm(alarmID int64) error {
	// 更新告警状态为已忽略
	if err := database.DB.Model(&model.DeviceAlarm{}).
		Where("id = ? AND status = ?", alarmID, model.AlarmStatusPending).
		Updates(map[string]interface{}{
			"status":      model.AlarmStatusIgnored,
			"update_time": time.Now(),
		}).Error; err != nil {
		return errors.NewWithError(errors.ErrDatabase, err)
	}

	return nil
}

// GetDeviceAlarms 获取设备告警列表
func (s *DeviceAlarmService) GetDeviceAlarms(deviceSN string, status *int) ([]model.DeviceAlarm, error) {
	var alarms []model.DeviceAlarm
	query := database.DB.Where("device_sn = ?", deviceSN)
	
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	
	if err := query.Order("create_time DESC").Find(&alarms).Error; err != nil {
		return nil, errors.NewWithError(errors.ErrDatabase, err)
	}
	
	return alarms, nil
}

// 内部辅助函数

func isValidAlarmType(alarmType int) bool {
	return alarmType == model.AlarmTypeStorage ||
		alarmType == model.AlarmTypeCPUTemp ||
		alarmType == model.AlarmTypeMemory
}

func isValidAlarmLevel(alarmLevel int) bool {
	return alarmLevel == model.AlarmLevelInfo ||
		alarmLevel == model.AlarmLevelWarning ||
		alarmLevel == model.AlarmLevelCritical
} 