package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Database DatabaseConfig `yaml:"database" json:"database"`
	Log      LogConfig      `yaml:"log" json:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
	Mode string `yaml:"mode" json:"mode"` // debug 或 release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `yaml:"type" json:"type"`         // sqlite, mysql, postgres
	DSN      string `yaml:"dsn" json:"dsn"`           // 连接字符串
	Host     string `yaml:"host" json:"host"`         // 主机
	Port     int    `yaml:"port" json:"port"`         // 端口
	User     string `yaml:"user" json:"user"`         // 用户名
	Password string `yaml:"password" json:"password"` // 密码
	DBName   string `yaml:"dbname" json:"dbname"`     // 数据库名
	SSLMode  string `yaml:"sslmode" json:"sslmode"`   // SSL模式 (postgres)
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" json:"level"`           // debug, info, warn, error
	Format     string `yaml:"format" json:"format"`         // 日志格式
	Output     string `yaml:"output" json:"output"`         // 输出方式: file, stdout
	FilePath   string `yaml:"filepath" json:"filepath"`     // 日志文件路径
	MaxSize    int    `yaml:"maxsize" json:"maxsize"`       // 单个文件最大大小(MB)
	MaxBackups int    `yaml:"maxbackups" json:"maxbackups"` // 最大备份数
	MaxAge     int    `yaml:"maxage" json:"maxage"`         // 最大保留天数
	Compress   bool   `yaml:"compress" json:"compress"`     // 是否压缩备份
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Type:    "sqlite",
			DSN:     "./data/sqli_lab.db",
			Host:    "localhost",
			Port:    3306,
			User:    "root",
			DBName:  "security",
			SSLMode: "disable",
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "[%timestamp%] [%level%] [GID:%goroutine%] [%module%] [%function%] [%file%:%line%] - %message%",
			Output:     "file",
			FilePath:   "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
	}
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	// 如果配置文件存在，读取并解析
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 环境变量覆盖
	applyEnvOverrides(cfg)

	return cfg, nil
}

// applyEnvOverrides 应用环境变量覆盖
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SQLI_LAB_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SQLI_LAB_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}
	if v := os.Getenv("SQLI_LAB_MODE"); v != "" {
		cfg.Server.Mode = v
	}
	if v := os.Getenv("SQLI_LAB_DB_TYPE"); v != "" {
		cfg.Database.Type = v
	}
	if v := os.Getenv("SQLI_LAB_DB_DSN"); v != "" {
		cfg.Database.DSN = v
	}
	if v := os.Getenv("SQLI_LAB_DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("SQLI_LAB_DB_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Database.Port)
	}
	if v := os.Getenv("SQLI_LAB_DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("SQLI_LAB_DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("SQLI_LAB_DB_NAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("SQLI_LAB_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("SQLI_LAB_LOG_FILE"); v != "" {
		cfg.Log.FilePath = v
	}
}
