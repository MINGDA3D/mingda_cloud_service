package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/pkg/config"
	"mingda_cloud_service/internal/pkg/database"
	"mingda_cloud_service/internal/pkg/logger"
	"mingda_cloud_service/internal/pkg/redis"
	"mingda_cloud_service/internal/pkg/rabbitmq"
	"mingda_cloud_service/internal/app/handler"
	"mingda_cloud_service/internal/pkg/middleware"
)

type App struct {
	config *config.Config
	engine *gin.Engine
}

func NewApp(cfg *config.Config) (*App, error) {
	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		return nil, fmt.Errorf("init logger error: %v", err)
	}

	// 初始化数据库连接
	if err := database.Init(cfg.Database); err != nil {
		return nil, fmt.Errorf("init database error: %v", err)
	}

	// 初始化Redis连接
	if err := redis.Init(cfg.Redis); err != nil {
		return nil, fmt.Errorf("init redis error: %v", err)
	}

	// 初始化RabbitMQ连接
	if err := rabbitmq.Init(cfg.RabbitMQ); err != nil {
		return nil, fmt.Errorf("init rabbitmq error: %v", err)
	}

	// 设置gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建gin引擎
	engine := gin.New()

	// 使用日志和恢复中间件
	engine.Use(gin.Logger(), gin.Recovery())

	return &App{
		config: cfg,
		engine: engine,
	}, nil
}

func (a *App) Run() error {
	// 注册路由
	a.registerRoutes()

	// 启动HTTP服务
	addr := fmt.Sprintf(":%d", a.config.Server.Port)
	return http.ListenAndServe(addr, a.engine)
}

func (a *App) registerRoutes() {
	// 使用中间件
	a.engine.Use(middleware.Logger())
	a.engine.Use(gin.Recovery())

	// 创建处理器
	authHandler := handler.NewAuthHandler(a.config.Server.JWTSecret, a.config.Server.AESKey)
	deviceInfoHandler := handler.NewDeviceInfoHandler()
	deviceStatusHandler := handler.NewDeviceStatusHandler()
	deviceAlarmHandler := handler.NewDeviceAlarmHandler()
	printTaskHandler := handler.NewPrintTaskHandler()

	// API v1 路由组
	v1 := a.engine.Group("/api/v1")
	{
		// 公开接口
		v1.GET("/health", handler.HealthCheck)
		
		// 设备认证接口
		v1.POST("/devices/register", authHandler.Register)
		v1.POST("/devices/auth", authHandler.Authenticate)
		v1.POST("/devices/refresh", authHandler.RefreshToken)

		// 需要认证的接口
		auth := v1.Group("/", middleware.AuthRequired())
		{
			// 设备信息相关路由
			deviceGroup := auth.Group("/device")
			{
				deviceGroup.POST("/info", deviceInfoHandler.ReportDeviceInfo)
				deviceGroup.POST("/status", deviceStatusHandler.ReportDeviceStatus)
				// 设备告警相关路由
				deviceGroup.POST("/alarm", deviceAlarmHandler.ReportDeviceAlarm)
				deviceGroup.GET("/alarms", deviceAlarmHandler.GetDeviceAlarms)
				deviceGroup.POST("/alarm/:id/resolve", deviceAlarmHandler.ResolveAlarm)
				deviceGroup.POST("/alarm/:id/ignore", deviceAlarmHandler.IgnoreAlarm)
				// 打印任务相关路由
				deviceGroup.POST("/print/status", printTaskHandler.ReportPrintStatus)
				deviceGroup.GET("/print/tasks", printTaskHandler.GetDevicePrintTasks)
				deviceGroup.GET("/print/task/:task_id/history", printTaskHandler.GetTaskHistory)
			}
		}
	}
} 