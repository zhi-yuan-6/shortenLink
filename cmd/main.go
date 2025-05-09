package main

import (
	"github.com/gin-gonic/gin"
	"short_link_generation/routers"
	"short_link_generation/storage"
)

func main() {
	store := storage.NewMemoryStore() // 创建内存存储实例

	// 创建一个默认的路由引擎
	r := gin.Default()

	routers.SetupRouters(r, store) // 设置路由

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
