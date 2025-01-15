package router

import (
	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/handler"
	"mingda_cloud_service/internal/pkg/middleware"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {
	// 创建处理器实例
	authHandler := handler.NewAuthHandler()
	deviceInfoHandler := handler.NewDeviceInfoHandler()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		v1.POST("/devices/register", authHandler.RegisterDevice)
		v1.POST("/devices/auth", authHandler.AuthenticateDevice)
		v1.POST("/devices/refresh", middleware.JWTAuth(), authHandler.RefreshToken)

		// 设备信息相关路由（需要认证）
		deviceGroup := v1.Group("/device", middleware.JWTAuth())
		{
			deviceGroup.POST("/info", deviceInfoHandler.ReportDeviceInfo)
		}
	}
} 