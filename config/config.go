package config

import (
	"os"
	"strconv"
	
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	AdminUserID   int64
	WebhookURL    string
	WebhookPort   string
	EncryptionKey string
}

func Load() *Config {
	// Load .env file - this will work even if .env doesn't exist
	godotenv.Load()
	
	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		AdminUserID:   getEnvInt64("ADMIN_USER_ID", 304528450),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		WebhookPort:   getEnv("PORT", "8443"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", "default-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}
