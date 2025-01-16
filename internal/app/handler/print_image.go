package handler

import (
    "path/filepath"
    "strings"
    "github.com/gin-gonic/gin"
    "mingda_cloud_service/internal/app/service"
    "mingda_cloud_service/internal/pkg/constants"
    "mingda_cloud_service/internal/pkg/response"
    "mingda_cloud_service/internal/pkg/errors"
)

// PrintImageHandler 打印图片处理器
type PrintImageHandler struct {
    imageService *service.PrintImageService
}

// NewPrintImageHandler 创建打印图片处理器实例
func NewPrintImageHandler(uploadDir, baseURL string) *PrintImageHandler {
    return &PrintImageHandler{
        imageService: service.NewPrintImageService(uploadDir, baseURL),
    }
}

// UploadPrintImage 处理打印图片上传请求
func (h *PrintImageHandler) UploadPrintImage(c *gin.Context) {
    // 从上下文获取设备SN
    deviceSN := c.GetString(constants.ContextDeviceSN)
    if deviceSN == "" {
        response.Error(c, errors.New(errors.ErrUnauthorized, "未授权的访问"))
        return
    }
    
    // 获取任务ID
    taskID := c.PostForm("task_id")
    if taskID == "" {
        response.Error(c, errors.New(errors.ErrInvalidParams, "缺少任务ID"))
        return
    }
    
    // 获取上传的文件
    file, err := c.FormFile("image")
    if err != nil {
        response.Error(c, errors.New(errors.ErrInvalidParams, "未找到上传的图片"))
        return
    }
    
    // 检查文件类型
    if !isValidImageType(file.Filename) {
        response.Error(c, errors.New(errors.ErrInvalidParams, "不支持的图片格式，仅支持jpg、jpeg、png格式"))
        return
    }
    
    // 处理图片上传
    image, err := h.imageService.UploadPrintImage(deviceSN, taskID, file)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, image)
}

// GetPrintImages 获取打印图片列表
func (h *PrintImageHandler) GetPrintImages(c *gin.Context) {
    // 从上下文获取设备SN
    deviceSN := c.GetString(constants.ContextDeviceSN)
    if deviceSN == "" {
        response.Error(c, errors.New(errors.ErrUnauthorized, "未授权的访问"))
        return
    }
    
    // 获取任务ID参数（可选）
    taskID := c.Query("task_id")
    
    // 获取图片列表
    images, err := h.imageService.GetPrintImages(deviceSN, taskID)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, images)
}

// isValidImageType 检查是否为有效的图片类型
func isValidImageType(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
} 