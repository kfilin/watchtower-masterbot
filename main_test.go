package main

import (
	"os"
	"testing"

	"github.com/kfilin/watchtower-masterbot/config"
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
	cfg := config.Load()
	
	if cfg.TelegramToken != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", cfg.TelegramToken)
	}
	
	if cfg.HealthPort != "18080" {
		t.Errorf("Expected health port '18080', got '%s'", cfg.HealthPort)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test missing token - should not error but return empty
	originalToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	
	cfg := config.Load()
	if cfg.TelegramToken != "" {
		t.Error("Expected empty token when not set")
	}
	
	// Restore
	os.Setenv("TELEGRAM_BOT_TOKEN", originalToken)
}
