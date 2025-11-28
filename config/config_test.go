package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("ADMIN_USER_ID", "12345")
	os.Setenv("HEALTH_PORT", "8080")
	
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if cfg.TelegramBotToken != "test-token" {
		t.Errorf("Expected TelegramBotToken 'test-token', got '%s'", cfg.TelegramBotToken)
	}
	
	if cfg.AdminUserID != 12345 {
		t.Errorf("Expected AdminUserID 12345, got '%d'", cfg.AdminUserID)
	}
	
	if cfg.HealthPort != "8080" {
		t.Errorf("Expected HealthPort '8080', got '%s'", cfg.HealthPort)
	}
}

func TestLoadConfigMissingToken(t *testing.T) {
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error for missing Telegram token, got none")
	}
	// Restore for other tests
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
}
