package main

import (
	//"analysis/api/routes"
	"analysis/configs"
	//"analysis/pkg/database"
	//"analysis/pkg/logging"
)

func main() {
	// 加载配置
	config := configs.LoadConfig()

	// 初始化日志
	logging.Init(config.Log)

	// 初始化数据库
	db, err := database.Init(config.DB)
	if err != nil {
		logging.Fatal("Failed to connect database", err)
	}

	// 初始化路由并启动服务
	router := routes.InitRouter(db, config)
	router.Run(":" + config.Server.Port)
}
