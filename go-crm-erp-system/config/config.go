package config

import (
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		DBPath:    getEnv("DB_PATH", "./crm_erp.db"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
