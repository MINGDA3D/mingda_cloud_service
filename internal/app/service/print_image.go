package service

import (
    "fmt"
    "path/filepath"
    "time"
    "mime/multipart"
    "os"
    "mingda_cloud_service/internal/app/model"
    "mingda_cloud_service/internal/pkg/database"
    "mingda_cloud_service/internal/pkg/errors"
)

// PrintImageService 打印图片服务
type PrintImageService struct {
    uploadDir string // 图片上传目录
    baseURL   string // 图片访问基础URL
}

// NewPrintImageService 创建打印图片服务实例
func NewPrintImageService(uploadDir, baseURL string) *PrintImageService {
    // 确保上传目录存在
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        panic(fmt.Sprintf("创建上传目录失败: %v", err))
    }
    
    return &PrintImageService{
        uploadDir: uploadDir,
        baseURL:   baseURL,
    }
}

// UploadPrintImage 上传打印图片
func (s *PrintImageService) UploadPrintImage(deviceSN, taskID string, file *multipart.FileHeader) (*model.PrintImage, error) {
    // 生成文件名
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("%s_%s_%d%s", deviceSN, taskID, time.Now().Unix(), ext)
    
    // 构建文件路径
    relativePath := filepath.Join(deviceSN, taskID, filename)
    fullPath := filepath.Join(s.uploadDir, relativePath)
    
    // 确保目录存在
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return nil, errors.NewWithError(errors.ErrSystem, fmt.Errorf("创建目录失败: %v", err))
    }
    
    // 保存文件
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return nil, errors.NewWithError(errors.ErrSystem, fmt.Errorf("创建目录失败: %v", err))
    }
    
    if err := saveUploadedFile(file, fullPath); err != nil {
        return nil, errors.NewWithError(errors.ErrSystem, fmt.Errorf("保存文件失败: %v", err))
    }
    
    // 构建图片URL
    imageURL := fmt.Sprintf("%s/%s", s.baseURL, relativePath)
    
    // 创建数据库记录
    image := &model.PrintImage{
        DeviceSN:   deviceSN,
        TaskID:     taskID,
        ImagePath:  relativePath,
        ImageURL:   imageURL,
        Status:     0, // 未检测
        HasDefect:  false,
        CreateTime: time.Now(),
        UpdateTime: time.Now(),
    }
    
    if err := database.DB.Create(image).Error; err != nil {
        // 删除已上传的文件
        os.Remove(fullPath)
        return nil, errors.NewWithError(errors.ErrDatabase, err)
    }
    
    return image, nil
}

// GetPrintImages 获取打印图片列表
func (s *PrintImageService) GetPrintImages(deviceSN, taskID string) ([]model.PrintImage, error) {
    var images []model.PrintImage
    query := database.DB.Where("device_sn = ?", deviceSN)
    
    if taskID != "" {
        query = query.Where("task_id = ?", taskID)
    }
    
    if err := query.Order("create_time DESC").Find(&images).Error; err != nil {
        return nil, errors.NewWithError(errors.ErrDatabase, err)
    }
    
    return images, nil
}

// 保存上传的文件
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
    
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()
    
    // 创建文件
    if err = os.WriteFile(dst, nil, 0644); err != nil {
        return err
    }
    
    // 打开文件准备写入
    out, err = os.OpenFile(dst, os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer out.Close()
    
    // 从源文件拷贝到目标文件
    buf := make([]byte, 1024*1024) // 1MB buffer
    for {
        n, err := src.Read(buf)
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return err
        }
        if n == 0 {
            break
        }
        
        if _, err := out.Write(buf[:n]); err != nil {
            return err
        }
    }
    
    return nil
} 