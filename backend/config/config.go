package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 集中保存服务配置，从 .env / 环境变量读取（敏感值不入 Git）。
type Config struct {
	ServerPort    string // Gin 监听端口，默认 8080
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string // 管理员 Token 签名密钥
	UploadDir     string // V1 图片本地存储目录
	StorageDriver string // V1: local；二期: oss

	AdminUser     string // 首次启动自动创建的超级管理员
	AdminPassword string // 其初始密码（首次登录后建议修改）
}

// Load 读取 .env（不存在不报错）并填充 Config。
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] 未找到 .env，将使用系统环境变量（.env.example 可作模板）")
	}
	return &Config{
		ServerPort:    getenv("SERVER_PORT", "8080"),
		DBHost:        getenv("DB_HOST", "127.0.0.1"),
		DBPort:        getenv("DB_PORT", "3306"),
		DBUser:        getenv("DB_USER", "root"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        getenv("DB_NAME", "project_showcase"),
		JWTSecret:     getenv("JWT_SECRET", "dev-only-change-me"),
		UploadDir:     getenv("UPLOAD_DIR", "./uploads"),
		StorageDriver: getenv("STORAGE_DRIVER", "local"),
		AdminUser:     getenv("ADMIN_USER", "admin"),
		AdminPassword: getenv("ADMIN_INIT_PASSWORD", "admin123456"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
