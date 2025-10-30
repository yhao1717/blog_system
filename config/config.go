package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBDriver   string
	DBPath     string
	JWTSecret  string
	ServerPort string
	LogLevel   string
}

func LoadConfig() *Config {
	config := &Config{
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBPath:     getEnv("DB_PATH", "blog.db"),
		JWTSecret:  getEnv("JWT_SECRET", "your_default_secret_key_change_in_production"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}

	// 验证必要的配置
	if config.JWTSecret == "your_default_secret_key_change_in_production" {
		log.Println("WARNING: Using default JWT secret key. Change this in production!")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
