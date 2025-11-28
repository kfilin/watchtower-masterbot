package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("ADMIN_USER_ID", "12345")
	os.Setenv("HEALTH_PORT", "8080")
	
	cfg := Load()  // ✅ CORRECTED: Load() not LoadConfig()
	
	if cfg.TelegramToken != "test-token" {
		t.Errorf("Expected TelegramToken 'test-token', got '%s'", cfg.TelegramToken)
	}
	if cfg.AdminID != 12345 {
		t.Errorf("Expected AdminID 12345, got '%d'", cfg.AdminID)
	}
	if cfg.HealthPort != "8080" {
		t.Errorf("Expected HealthPort '8080', got '%s'", cfg.HealthPort)
	}
}

func TestLoadConfigMissingToken(t *testing.T) {
	// Save original and unset
	originalToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	
	cfg := Load()  // ✅ CORRECTED: Load() not LoadConfig()
	
	if cfg.TelegramToken != "" {
		t.Error("Expected empty Telegram token when not set")
	}
	
	// Restore
	if originalToken != "" {
		os.Setenv("TELEGRAM_BOT_TOKEN", originalToken)
	}
}
