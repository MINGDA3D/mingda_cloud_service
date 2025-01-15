package main

import (
	"flag"
	"log"

	"mingda_cloud_service/internal/app"
	"mingda_cloud_service/internal/pkg/config"
)

var configFile = flag.String("config", "configs/config.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化应用
	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// 启动服务
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
