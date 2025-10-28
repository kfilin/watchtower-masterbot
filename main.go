package main

import (
	"log"
	
	"watchtower-masterbot/bot"
	"watchtower-masterbot/config"
	"watchtower-masterbot/servers"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Initialize server manager with encryption
	serverManager := servers.NewManager(cfg.EncryptionKey)
	
	// Create and start bot
	watchtowerBot := bot.NewBot(cfg, serverManager)
	
	log.Println("Starting WatchtowerMasterBot...")
	if err := watchtowerBot.Start(); err != nil {
		log.Fatal("Failed to start bot:", err)
	}
}
