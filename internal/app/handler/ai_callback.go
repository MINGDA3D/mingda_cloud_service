package handler

import (
    "mingda_cloud_service/internal/app/model"
    "mingda_cloud_service/internal/pkg/errors"
    "mingda_cloud_service/internal/pkg/response"
    "github.com/gin-gonic/gin"
)

// AICallbackHandler AI回调处理器
type AICallbackHandler struct {
    db *gorm.DB
}

// NewAICallbackHandler 创建AI回调处理器
func NewAICallbackHandler(db *gorm.DB) *AICallbackHandler {
    return &AICallbackHandler{
        db: db,
    }
}

// CallbackRequest AI回调请求
type CallbackRequest struct {
    TaskID  string         `json:"task_id" binding:"required"`
    Status  string         `json:"status" binding:"required"`
    Result  *PredictResult `json:"result,omitempty"`
}

// PredictResult AI预测结果
type PredictResult struct {
    HasDefect  bool    `json:"has_defect"`
    DefectType string  `json:"defect_type"`
    Confidence float64 `json:"confidence"`
}

// HandleCallback 处理AI回调
func (h *AICallbackHandler) HandleCallback(c *gin.Context) {
    var req CallbackRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, errors.New(errors.ErrInvalidParams, "参数错误"))
        return
    }

    // 开启事务
    tx := h.db.Begin()
    if tx.Error != nil {
        response.Error(c, errors.New(errors.ErrDatabase, "开启事务失败"))
        return
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 查找对应的图片记录
    var image model.PrintImage
    if err := tx.Where("task_id = ?", req.TaskID).First(&image).Error; err != nil {
        tx.Rollback()
        response.Error(c, errors.New(errors.ErrDatabase, "查询图片记录失败"))
        return
    }

    // 更新检测结果
    updates := map[string]interface{}{
        "status": model.StatusChecked,
    }

    if req.Status == "success" && req.Result != nil {
        hasDefect := req.Result.HasDefect
        updates["has_defect"] = &hasDefect
        updates["defect_type"] = req.Result.DefectType
        updates["confidence"] = req.Result.Confidence
    }

    if err := tx.Model(&image).Updates(updates).Error; err != nil {
        tx.Rollback()
        response.Error(c, errors.New(errors.ErrDatabase, "更新检测结果失败"))
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        response.Error(c, errors.New(errors.ErrDatabase, "提交事务失败"))
        return
    }

    response.Success(c, nil)
} 