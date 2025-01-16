package handler

import (
    "mingda_cloud_service/internal/app/service"
    "mingda_cloud_service/internal/pkg/errors"
    "mingda_cloud_service/internal/pkg/response"
    "mingda_cloud_service/internal/pkg/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "mingda_cloud_service/internal/pkg/config"
)

// PrintImageHandler 打印图片处理器
type PrintImageHandler struct {
    imageService *service.PrintImageService
}

// NewPrintImageHandler 创建打印图片处理器
func NewPrintImageHandler(db *gorm.DB, cfg *config.Config) *PrintImageHandler {
    return &PrintImageHandler{
        imageService: service.NewPrintImageService(db, cfg),
    }
}

// UploadPrintImage 上传打印图片
func (h *PrintImageHandler) UploadPrintImage(c *gin.Context) {
    // 从上下文获取设备SN
    deviceSN := c.GetString(middleware.ContextDeviceSN)
    if deviceSN == "" {
        response.Error(c, errors.New(errors.ErrInvalidParams, "设备SN不能为空"))
        return
    }

    // 获取任务ID
    taskID := c.PostForm("task_id")
    if taskID == "" {
        response.Error(c, errors.New(errors.ErrInvalidParams, "任务ID不能为空"))
        return
    }

    // 获取上传的文件
    file, err := c.FormFile("file")
    if err != nil {
        response.Error(c, errors.New(errors.ErrInvalidParams, "获取上传文件失败"))
        return
    }

    // 上传图片
    if err := h.imageService.UploadPrintImage(file, deviceSN, taskID); err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil)
}

// GetPrintImages 获取打印图片列表
func (h *PrintImageHandler) GetPrintImages(c *gin.Context) {
    // 从上下文获取设备SN
    deviceSN := c.GetString(middleware.ContextDeviceSN)
    if deviceSN == "" {
        response.Error(c, errors.New(errors.ErrInvalidParams, "设备SN不能为空"))
        return
    }

    // 获取任务ID（可选）
    taskID := c.Query("task_id")

    // 获取图片列表
    images, err := h.imageService.GetPrintImages(deviceSN, taskID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, images)
} 