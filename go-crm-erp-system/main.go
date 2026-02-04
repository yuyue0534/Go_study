package main

import (
	"crm-erp-system/config"
	"crm-erp-system/database"
	"crm-erp-system/router"
	"log"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 初始化路由
	r := router.SetupRouter()

	// 启动服务
	port := config.AppConfig.Port
	log.Printf("服务启动在端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
