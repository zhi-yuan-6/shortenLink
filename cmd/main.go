package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"shortenLink/config"
	"shortenLink/routers"
	"shortenLink/storage"
)

func main() {
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//连接数据库并迁移
	store := storage.NewMemoryStore() // 创建内存存储实例
	sqlDB := connectSql(cfg)
	defer sqlDB.Close()

	// 创建一个默认的路由引擎
	r := gin.Default()

	routers.SetupRouters(r, store) // 设置路由

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

func connectSql(cfg *config.Config) *sql.DB {
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
	gormDB, err := config.NewGORM(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 条件执行自动迁移
	if cfg.Postgres.AutoMigrate {
		if err := config.AutoMigrate(gormDB); err != nil {
			log.Fatalf("Auto migration failed: %v", err)
		}
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB
}
