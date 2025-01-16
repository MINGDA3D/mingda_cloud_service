package service

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "time"
    
    "mingda_cloud_service/internal/app/model"
    "mingda_cloud_service/internal/pkg/ai"
    "mingda_cloud_service/internal/pkg/config"
    "mingda_cloud_service/internal/pkg/errors"
    "gorm.io/gorm"
)

// PrintImageService 打印图片服务
type PrintImageService struct {
    db       *gorm.DB
    config   *config.Config
    aiClient *ai.Client
}

// NewPrintImageService 创建打印图片服务
func NewPrintImageService(db *gorm.DB, cfg *config.Config) *PrintImageService {
    aiClient := ai.NewClient(
        cfg.AI.BaseURL,
        fmt.Sprintf("%s/api/v1/ai/callback", cfg.Server.BaseURL),
    )
    
    return &PrintImageService{
        db:       db,
        config:   cfg,
        aiClient: aiClient,
    }
}

// UploadPrintImage 上传打印图片
func (s *PrintImageService) UploadPrintImage(file *multipart.FileHeader, deviceSN, taskID string) error {
    // 创建上传目录
    uploadDir := filepath.Join("uploads", "images", time.Now().Format("20060102"))
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return errors.New(errors.ErrSystem, fmt.Sprintf("创建目录失败: %v", err))
    }

    // 生成文件名
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("%s_%s%s", taskID, time.Now().Format("150405"), ext)
    filepath := filepath.Join(uploadDir, filename)

    // 保存文件
    src, err := file.Open()
    if err != nil {
        return errors.New(errors.ErrSystem, fmt.Sprintf("打开上传文件失败: %v", err))
    }
    defer src.Close()

    dst, err := os.Create(filepath)
    if err != nil {
        return errors.New(errors.ErrSystem, fmt.Sprintf("创建目标文件失败: %v", err))
    }
    defer dst.Close()

    if _, err = io.Copy(dst, src); err != nil {
        return errors.New(errors.ErrSystem, fmt.Sprintf("保存文件失败: %v", err))
    }

    // 构建图片URL
    imageURL := fmt.Sprintf("%s/images/%s/%s", s.config.Server.BaseURL, time.Now().Format("20060102"), filename)

    // 开启事务
    tx := s.db.Begin()
    if tx.Error != nil {
        return errors.New(errors.ErrDatabase, "开启事务失败")
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 创建图片记录
    image := &model.PrintImage{
        TaskID:    taskID,
        DeviceSN:  deviceSN,
        ImagePath: filepath,
        ImageURL:  imageURL,
        Status:    model.StatusChecking, // 设置为检测中状态
    }

    if err := tx.Create(image).Error; err != nil {
        tx.Rollback()
        return errors.New(errors.ErrDatabase, fmt.Sprintf("保存图片记录失败: %v", err))
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return errors.New(errors.ErrDatabase, "提交事务失败")
    }

    // 触发AI检测
    go func() {
        if err := s.aiClient.RequestPredict(imageURL, taskID); err != nil {
            // 更新状态为未检测，等待重试
            s.db.Model(&model.PrintImage{}).
                Where("task_id = ?", taskID).
                Update("status", model.StatusPending)
        }
    }()

    return nil
}

// GetPrintImages 获取打印图片列表
func (s *PrintImageService) GetPrintImages(deviceSN string, taskID string) ([]model.PrintImage, error) {
    var images []model.PrintImage
    query := s.db.Where("device_sn = ?", deviceSN)
    
    if taskID != "" {
        query = query.Where("task_id = ?", taskID)
    }
    
    if err := query.Find(&images).Error; err != nil {
        return nil, errors.New(errors.ErrDatabase, fmt.Sprintf("查询图片列表失败: %v", err))
    }
    
    return images, nil
} 