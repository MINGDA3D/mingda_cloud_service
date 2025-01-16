package handler

import (
	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/service"
	"mingda_cloud_service/internal/pkg/constants"
	"mingda_cloud_service/internal/pkg/response"
)

// PrintTaskHandler 打印任务处理器
type PrintTaskHandler struct {
	printTaskService *service.PrintTaskService
}

// NewPrintTaskHandler 创建打印任务处理器实例
func NewPrintTaskHandler() *PrintTaskHandler {
	return &PrintTaskHandler{
		printTaskService: service.NewPrintTaskService(),
	}
}

// ReportPrintStatus 处理打印任务状态上报请求
func (h *PrintTaskHandler) ReportPrintStatus(c *gin.Context) {
	// 1. 获取设备SN
	deviceSN := c.GetString(constants.ContextDeviceSN)
	if deviceSN == "" {
		response.Error(c, response.ErrUnauthorized, "未获取到设备SN")
		return
	}

	// 2. 绑定请求参数
	var req service.PrintTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, "参数绑定失败")
		return
	}

	// 3. 处理打印任务状态上报
	if err := h.printTaskService.ReportPrintStatus(deviceSN, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetDevicePrintTasks 获取设备打印任务列表
func (h *PrintTaskHandler) GetDevicePrintTasks(c *gin.Context) {
	// 1. 获取设备SN
	deviceSN := c.GetString(constants.ContextDeviceSN)
	if deviceSN == "" {
		response.Error(c, response.ErrUnauthorized, "未获取到设备SN")
		return
	}

	// 2. 获取状态过滤参数
	status := c.Query("status")

	// 3. 查询打印任务列表
	tasks, err := h.printTaskService.GetDevicePrintTasks(deviceSN, status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, tasks)
}

// GetTaskHistory 获取任务状态变更历史
func (h *PrintTaskHandler) GetTaskHistory(c *gin.Context) {
	// 1. 获取任务ID
	taskID := c.Param("task_id")
	if taskID == "" {
		response.Error(c, response.ErrInvalidParams, "未提供任务ID")
		return
	}

	// 2. 查询任务历史
	history, err := h.printTaskService.GetTaskHistory(taskID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, history)
} 