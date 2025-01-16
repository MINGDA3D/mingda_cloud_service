package handler

import (
	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/service"
	"mingda_cloud_service/internal/pkg/constants"
	"mingda_cloud_service/internal/pkg/response"
)

// DeviceInfoHandler 设备信息处理器
type DeviceInfoHandler struct {
	deviceInfoService *service.DeviceInfoService
}

// NewDeviceInfoHandler 创建设备信息处理器实例
func NewDeviceInfoHandler() *DeviceInfoHandler {
	return &DeviceInfoHandler{
		deviceInfoService: service.NewDeviceInfoService(),
	}
}

// ReportDeviceInfo 处理设备信息上报请求
func (h *DeviceInfoHandler) ReportDeviceInfo(c *gin.Context) {
	var req service.DeviceInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 从上下文获取设备SN并验证
	deviceSN := c.GetString(constants.ContextDeviceSN)
	if deviceSN != req.DeviceInfo.DeviceSN {
		response.ErrorWithMsg(c, "设备SN不匹配")
		return
	}

	// 调用服务处理上报请求
	if err := h.deviceInfoService.ReportDeviceInfo(&req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"success": true})
} 