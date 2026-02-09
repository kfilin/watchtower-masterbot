package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/kfilin/watchtower-masterbot/bot"
	"github.com/kfilin/watchtower-masterbot/config"
	"github.com/kfilin/watchtower-masterbot/health"
	"github.com/kfilin/watchtower-masterbot/web"
)

func main() {
	// 0. Load optional .env file
	_ = godotenv.Load()

	// 1. Load Config
	cfg := config.Load()

	// 2. Initialize Bot (Graceful Error Handling)
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		encryptionKey = "default-encryption-key-change-in-production"
		log.Println("⚠️  Using default encryption key - set ENCRYPTION_KEY for production")
	}

	botInstance, err := bot.NewBot(cfg.TelegramToken, cfg.AdminID, encryptionKey, cfg.WebAppURL)

	// 3. Start Health & Web Server
	log.Printf("🏥 Starting Health & Web Server on port %s...", cfg.HealthPort)

	registerWeb := func(mux *http.ServeMux) {
		if err == nil {
			webServer := web.NewServer(botInstance.GetManager(), cfg.AdminID, cfg.TelegramToken)
			webServer.RegisterHandlers(mux)
			log.Println("⚡ Retro Terminal TWA registered at /terminal")
		}
	}

	if hErr := health.StartServer(cfg.HealthPort, registerWeb); hErr != nil {
		log.Printf("❌ Health server failed to start: %v", hErr)
		os.Exit(1)
	}
	if err != nil {
		log.Println("---------------------------------------------------")
		log.Printf("❌ Telegram Bot Error: %v", err)
		log.Println("💡 The application will continue running with health endpoints only")
		log.Printf("💡 Health endpoints available at: http://localhost:%s/health", cfg.HealthPort)
		log.Println("💡 Fix the token and restart to enable Telegram features")
		log.Println("---------------------------------------------------")
		health.SetBotStatus("failed")

		// Don't exit! Keep the health server running
		// Wait for shutdown signal
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("🛑 Shutting down health server...")
		health.Shutdown()
		os.Exit(0)
	}

	// 4. Start Bot (only if initialization succeeded)
	health.SetBotStatus("running")
	go botInstance.Start()
	log.Printf("✅ Telegram bot started successfully! Health endpoints available at: http://localhost:%s/health", cfg.HealthPort)

	// 5. Keep Alive
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("🛑 Shutting down...")
	health.Shutdown()
}
