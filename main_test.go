package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set test environment variables
	os.Setenv("HEALTH_PORT", "18080")
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("ADMIN_USER_ID", "12345")
	
	code := m.Run()
	os.Exit(code)
}

func TestConfigLoad(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if cfg.TelegramBotToken != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", cfg.TelegramBotToken)
	}
	
	if cfg.HealthPort != "18080" {
		t.Errorf("Expected health port '18080', got '%s'", cfg.HealthPort)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test missing token
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	_, err := loadConfig()
	if err == nil {
		t.Error("Expected error for missing token, got none")
	}
	
	// Restore
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
}
