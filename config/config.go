package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TelegramToken string
	AdminID       int64
	HealthPort    string
	EncryptionKey string
	WebAppURL     string
}

func Load() *Config {
	return &Config{
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		AdminID:       getEnvAsInt("ADMIN_USER_ID", 0),
		HealthPort:    getEnv("HEALTH_PORT", "8080"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		WebAppURL:     getEnv("WEBAPP_URL", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSpace(value)
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int64) int64 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}
