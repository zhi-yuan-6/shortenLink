package main

import "C"
import (
	"github.com/gin-gonic/gin"
	"log"
	"shortenLink/cache"
	"shortenLink/config"
	"shortenLink/dto"
	"shortenLink/models"
	"shortenLink/routers"
)

// 配置client
func main() {
	err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//连接数据库并迁移
	/*store := memory.NewMemoryStore() // 创建内存存储实例*/
	InitSql()

	//初始化多级缓存
	cache.RedisCacheIns = cache.NewRedisCache(config.Cfg.Redis)

	// 创建一个默认的路由引擎
	r := gin.Default()

	routers.SetupRouters(r) // 设置路由

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

func InitSql() {
	/*db, err := config.NewPostgres(config.PostgresConfig{
		Host:         cfg.Postgres.Host,
		Port:         cfg.Postgres.Port,
		User:         cfg.Postgres.User,
		Password:     cfg.Postgres.Password,
		DBName:       cfg.Postgres.DBName,
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
		MaxIdleTime:  cfg.Postgres.MaxIdleTime,
		SSLMode:      cfg.Postgres.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()*/
	err := dto.NewGORM()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 条件执行自动迁移
	if config.Cfg.Postgres.AutoMigrate {
		err := dto.DB.AutoMigrate(&models.ShortUrl{}, &models.VisitStats{})
		if err != nil {
			log.Fatalf("Auto migration failed: %v", err)
		}
	}
	/*sqlDB, err := storage.DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB*/
}
