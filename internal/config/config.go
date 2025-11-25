package config

import (
	"os"
)

// Config 应用配置结构
type Config struct {
	Port             string
	AuthToken        string
	DBURL            string
	BoxOfficeURL     string
	BoxOfficeAPIKey  string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		AuthToken:        getEnv("AUTH_TOKEN", ""),
		DBURL:            getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/movies?sslmode=disable"),
		BoxOfficeURL:     getEnv("BOXOFFICE_URL", ""),
		BoxOfficeAPIKey:  getEnv("BOXOFFICE_API_KEY", ""),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
