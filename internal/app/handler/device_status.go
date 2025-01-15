package handler

import (
	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/service"
	"mingda_cloud_service/internal/pkg/response"
)

// DeviceStatusHandler 设备状态处理器
type DeviceStatusHandler struct {
	statusService *service.DeviceStatusService
}

// NewDeviceStatusHandler 创建设备状态处理器实例
func NewDeviceStatusHandler() *DeviceStatusHandler {
	return &DeviceStatusHandler{
		statusService: service.NewDeviceStatusService(),
	}
}

// ReportDeviceStatus 处理设备状态上报请求
func (h *DeviceStatusHandler) ReportDeviceStatus(c *gin.Context) {
	var req service.DeviceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 从上下文获取设备SN
	deviceSN := c.GetString("device_sn")
	if deviceSN == "" {
		response.Error(c, errors.New(errors.ErrUnauthorized, "未授权的访问"))
		return
	}
	req.DeviceSN = deviceSN

	// 处理设备状态上报
	if err := h.statusService.ReportDeviceStatus(&req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
} 