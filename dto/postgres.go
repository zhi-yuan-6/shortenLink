package dto

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"shortenLink/config"
)

var DB *gorm.DB

func NewPostgres(p config.PostgresConfig) (*sqlx.DB, error) {
	connStr := config.Cfg.Postgres.DSN()

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	//连接池配置
	db.SetMaxOpenConns(p.MaxOpenConns)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetConnMaxIdleTime(p.MaxIdleTime)

	//测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}
	return db, nil
}

func NewGORM() error {
	dsn := config.Cfg.Postgres.DSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.Cfg.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Cfg.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(config.Cfg.Postgres.MaxIdleTime)

	return nil
}
