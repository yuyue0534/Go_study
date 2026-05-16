package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host   string
	Port   string
	DBPath string
}

func Load() Config {
	return Config{
		Host:   envOrDefault("STUDENT_SYS_HOST", "127.0.0.1"),
		Port:   envOrDefault("STUDENT_SYS_PORT", "8080"),
		DBPath: envOrDefault("STUDENT_SYS_DB_PATH", "./data/student.db"),
	}
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
