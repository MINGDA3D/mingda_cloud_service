package service

import (
	"fmt"
	"time"
	"mingda_cloud_service/internal/app/model"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/errors"
)

// PrintTaskService 打印任务服务
type PrintTaskService struct{}

// NewPrintTaskService 创建打印任务服务实例
func NewPrintTaskService() *PrintTaskService {
	return &PrintTaskService{}
}

// PrintTaskRequest 打印任务状态上报请求
type PrintTaskRequest struct {
	TaskID          string     `json:"task_id" binding:"required"`
	FileName        string     `json:"file_name" binding:"required"`
	Status          string     `json:"status" binding:"required"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Progress        float64    `json:"progress"`
	Duration        int        `json:"duration" binding:"required"`
	FilamentUsed    float64    `json:"filament_used"`
	LayersCompleted int        `json:"layers_completed"`
	ErrorCode       string     `json:"error_code"`
	CancelReason    string     `json:"cancel_reason"`
}

// ReportPrintStatus 上报打印任务状态
func (s *PrintTaskService) ReportPrintStatus(deviceSN string, req *PrintTaskRequest) error {
	// 1. 验证任务状态是否有效
	if !isValidTaskStatus(req.Status) {
		return errors.New(errors.ErrInvalidParams, "无效的任务状态")
	}

	// 2. 开启事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		return errors.New(errors.ErrDatabase, "开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. 查找或创建打印任务
	var task model.PrintTask
	result := tx.Where("task_id = ?", req.TaskID).First(&task)
	if result.Error != nil {
		// 如果任务不存在，创建新任务
		task = model.PrintTask{
			TaskID:          req.TaskID,
			DeviceSN:        deviceSN,
			FileName:        req.FileName,
			StartTime:       req.StartTime,
			Status:          req.Status,
			Progress:        req.Progress,
			Duration:        req.Duration,
			FilamentUsed:    req.FilamentUsed,
			LayersCompleted: req.LayersCompleted,
		}
		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.ErrDatabase, fmt.Sprintf("创建打印任务失败: %v", err))
		}
	} else {
		// 4. 记录状态变更历史
		if task.Status != req.Status {
			history := model.PrintTaskHistory{
				TaskID:         req.TaskID,
				DeviceSN:       deviceSN,
				PreviousStatus: task.Status,
				CurrentStatus:  req.Status,
				ChangeTime:     time.Now(),
			}
			if err := tx.Create(&history).Error; err != nil {
				tx.Rollback()
				return errors.New(errors.ErrDatabase, fmt.Sprintf("记录状态变更历史失败: %v", err))
			}
		}

		// 5. 更新任务状态
		updates := map[string]interface{}{
			"status":           req.Status,
			"progress":        req.Progress,
			"duration":        req.Duration,
			"filament_used":   req.FilamentUsed,
			"layers_completed": req.LayersCompleted,
		}

		// 根据状态设置特定字段
		switch req.Status {
		case model.TaskStatusCompleted, model.TaskStatusCancelled, model.TaskStatusError:
			updates["end_time"] = req.EndTime
			if req.ErrorCode != "" {
				updates["error_code"] = req.ErrorCode
			}
			if req.CancelReason != "" {
				updates["cancel_reason"] = req.CancelReason
			}
		}

		if err := tx.Model(&task).Updates(updates).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.ErrDatabase, fmt.Sprintf("更新打印任务失败: %v", err))
		}
	}

	// 6. 提交事务
	if err := tx.Commit().Error; err != nil {
		return errors.New(errors.ErrDatabase, fmt.Sprintf("提交事务失败: %v", err))
	}

	return nil
}

// GetDevicePrintTasks 获取设备打印任务列表
func (s *PrintTaskService) GetDevicePrintTasks(deviceSN string, status string) ([]model.PrintTask, error) {
	var tasks []model.PrintTask
	query := database.DB.Where("device_sn = ?", deviceSN)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Order("create_time DESC").Find(&tasks).Error; err != nil {
		return nil, errors.New(errors.ErrDatabase, fmt.Sprintf("查询打印任务失败: %v", err))
	}
	
	return tasks, nil
}

// GetTaskHistory 获取任务状态变更历史
func (s *PrintTaskService) GetTaskHistory(taskID string) ([]model.PrintTaskHistory, error) {
	var history []model.PrintTaskHistory
	if err := database.DB.Where("task_id = ?", taskID).
		Order("change_time DESC").
		Find(&history).Error; err != nil {
		return nil, errors.New(errors.ErrDatabase, fmt.Sprintf("查询任务历史失败: %v", err))
	}
	return history, nil
}

// isValidTaskStatus 验证任务状态是否有效
func isValidTaskStatus(status string) bool {
	validStatuses := map[string]bool{
		model.TaskStatusIdle:      true,
		model.TaskStatusPrinting:  true,
		model.TaskStatusPaused:    true,
		model.TaskStatusResumed:   true,
		model.TaskStatusCompleted: true,
		model.TaskStatusCancelled: true,
		model.TaskStatusError:     true,
	}
	return validStatuses[status]
} 