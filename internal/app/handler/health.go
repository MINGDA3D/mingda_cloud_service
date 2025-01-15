package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/redis"
	"mingda_cloud_service/internal/pkg/rabbitmq"
	"mingda_cloud_service/internal/pkg/response"
)

type ServiceStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string          `json:"status"`
	Services  []ServiceStatus `json:"services"`
	Timestamp int64          `json:"timestamp"`
}

// HealthCheck 健康检查处理器
func HealthCheck(c *gin.Context) {
	health := HealthResponse{
		Status:    "ok",
		Services:  make([]ServiceStatus, 0),
		Timestamp: time.Now().Unix(),
	}

	// 检查数据库连接
	dbStatus := checkDatabase()
	health.Services = append(health.Services, dbStatus)
	if dbStatus.Status != "ok" {
		health.Status = "error"
	}

	// 检查Redis连接
	redisStatus := checkRedis(c.Request.Context())
	health.Services = append(health.Services, redisStatus)
	if redisStatus.Status != "ok" {
		health.Status = "error"
	}

	// 检查RabbitMQ连接
	mqStatus := checkRabbitMQ()
	health.Services = append(health.Services, mqStatus)
	if mqStatus.Status != "ok" {
		health.Status = "error"
	}

	response.Success(c, health)
}

func checkDatabase() ServiceStatus {
	status := ServiceStatus{
		Name:   "mysql",
		Status: "ok",
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		status.Status = "error"
		status.Message = "获取数据库连接失败"
		return status
	}

	if err := sqlDB.Ping(); err != nil {
		status.Status = "error"
		status.Message = "数据库连接失败"
		return status
	}

	return status
}

func checkRedis(ctx context.Context) ServiceStatus {
	status := ServiceStatus{
		Name:   "redis",
		Status: "ok",
	}

	if err := redis.Client.Ping(ctx).Err(); err != nil {
		status.Status = "error"
		status.Message = "Redis连接失败"
		return status
	}

	return status
}

func checkRabbitMQ() ServiceStatus {
	status := ServiceStatus{
		Name:   "rabbitmq",
		Status: "ok",
	}

	if !rabbitmq.IsConnected() {
		status.Status = "error"
		status.Message = "RabbitMQ连接已断开"
		return status
	}

	return status
} 