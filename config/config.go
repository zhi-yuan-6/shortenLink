package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

var Cfg *Config

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

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func LoadConfig(configPath string) error {
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
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func (pg *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode)
}
