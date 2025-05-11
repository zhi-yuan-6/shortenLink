package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	//Redis    RedisConfig    `mapstructure:"redis"`
}

type PostgresConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"dbname"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxIdleTime  time.Duration `mapstructure:"max_idle_time"`
	SSLMode      string        `mapstructure:"sslmode"`
	AutoMigrate  bool          `mapstructure:"auto_migrate"` // 新增自动迁移配置
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		_, filename, _, _ := runtime.Caller(0) // 获取当前文件的路径
		configPath = filepath.Join(filepath.Dir(filename), "../config.yaml")
	}
	viper.SetConfigFile(configPath) // 设置配置文件名

	//设置环境变量前缀并自动绑定
	viper.SetEnvPrefix("App")
	viper.AutomaticEnv()

	//优先加载环境变量覆盖配置
	viper.BindEnv("database.postgres.password")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func NewPostgres(p PostgresConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Password, p.DBName)

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
