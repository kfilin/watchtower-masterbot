package main

import (
	"log"
	
	"watchtower-masterbot/bot"
	"watchtower-masterbot/config"
	"watchtower-masterbot/servers"
	"watchtower-masterbot/health"  // ADD HEALTH IMPORT
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Initialize server manager with encryption
	serverManager := servers.NewManager(cfg.EncryptionKey)
	
	// START HEALTH SERVER - ADD THIS CRITICAL LINE
	health.StartHealthServer()
	
	// Create and start bot
	watchtowerBot := bot.NewBot(cfg, serverManager)
	
	log.Println("Starting WatchtowerMasterBot...")
	if err := watchtowerBot.Start(); err != nil {
		log.Fatal("Failed to start bot:", err)
	}
}
