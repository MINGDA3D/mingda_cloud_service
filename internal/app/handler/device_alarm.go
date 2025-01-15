package handler

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/service"
	"mingda_cloud_service/internal/pkg/middleware"
	"mingda_cloud_service/internal/pkg/response"
	"mingda_cloud_service/internal/pkg/errors"
)

// DeviceAlarmHandler 设备告警处理器
type DeviceAlarmHandler struct {
	alarmService *service.DeviceAlarmService
}

// NewDeviceAlarmHandler 创建设备告警处理器实例
func NewDeviceAlarmHandler() *DeviceAlarmHandler {
	return &DeviceAlarmHandler{
		alarmService: service.NewDeviceAlarmService(),
	}
}

// ReportDeviceAlarm 上报设备告警
func (h *DeviceAlarmHandler) ReportDeviceAlarm(c *gin.Context) {
	// 从上下文获取设备SN
	deviceSN := c.GetString(middleware.ContextDeviceSN)
	if deviceSN == "" {
		response.New(c).Error(errors.New(errors.ErrUnauthorized, "未授权的访问"))
		return
	}

	// 绑定请求参数
	var req service.DeviceAlarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.New(c).Error(errors.New(errors.ErrInvalidParams, "参数绑定错误"))
		return
	}

	// 处理请求
	if err := h.alarmService.ReportDeviceAlarm(deviceSN, &req); err != nil {
		response.New(c).Error(err)
		return
	}

	response.New(c).Success(nil)
}

// ResolveAlarm 处理告警
func (h *DeviceAlarmHandler) ResolveAlarm(c *gin.Context) {
	// 获取告警ID
	alarmID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.New(c).Error(errors.New(errors.ErrInvalidParams, "无效的告警ID"))
		return
	}

	// 绑定请求参数
	var req service.ResolveAlarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.New(c).Error(errors.New(errors.ErrInvalidParams, "参数绑定错误"))
		return
	}

	// 处理请求
	if err := h.alarmService.ResolveAlarm(alarmID, &req); err != nil {
		response.New(c).Error(err)
		return
	}

	response.New(c).Success(nil)
}

// IgnoreAlarm 忽略告警
func (h *DeviceAlarmHandler) IgnoreAlarm(c *gin.Context) {
	// 获取告警ID
	alarmID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.New(c).Error(errors.New(errors.ErrInvalidParams, "无效的告警ID"))
		return
	}

	// 处理请求
	if err := h.alarmService.IgnoreAlarm(alarmID); err != nil {
		response.New(c).Error(err)
		return
	}

	response.New(c).Success(nil)
}

// GetDeviceAlarms 获取设备告警列表
func (h *DeviceAlarmHandler) GetDeviceAlarms(c *gin.Context) {
	// 从上下文获取设备SN
	deviceSN := c.GetString(middleware.ContextDeviceSN)
	if deviceSN == "" {
		response.New(c).Error(errors.New(errors.ErrUnauthorized, "未授权的访问"))
		return
	}

	// 获取状态过滤参数
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err != nil {
			response.New(c).Error(errors.New(errors.ErrInvalidParams, "无效的状态参数"))
			return
		}
		status = &statusVal
	}

	// 获取告警列表
	alarms, err := h.alarmService.GetDeviceAlarms(deviceSN, status)
	if err != nil {
		response.New(c).Error(err)
		return
	}

	response.New(c).Success(alarms)
} 