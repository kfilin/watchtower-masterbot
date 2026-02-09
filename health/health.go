package health

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
	BotStatus string    `json:"bot_status,omitempty"`
}

var (
	startTime   time.Time
	server      *http.Server
	serverWg    sync.WaitGroup
	botStatus   string = "initializing"
	healthMutex sync.RWMutex
)

func init() {
	startTime = time.Now()
}

// StartServer starts the health check server on the specified port
func StartServer(port string, registerExtra func(*http.ServeMux)) error {
	if port == "" {
		port = "8080" // Default port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler)
	mux.HandleFunc("/live", liveHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("⚠️ 404 Unhandled: %s", r.URL.Path)
		http.NotFound(w, r)
	})

	if registerExtra != nil {
		registerExtra(mux)
	}

	// Logging middleware
	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🔍 HTTP Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		mux.ServeHTTP(w, r)
	})

	server = &http.Server{
		Addr:    ":" + port,
		Handler: loggingHandler,
	}

	log.Printf("🏥 Health server starting on port %s", port)

	// Start server in background
	serverWg.Add(1)
	go func() {
		defer serverWg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ Health server error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	return nil
}

// SetBotStatus updates the bot status for health checks
func SetBotStatus(status string) {
	healthMutex.Lock()
	defer healthMutex.Unlock()
	botStatus = status
}

// Shutdown gracefully shuts down the health server
func Shutdown() {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		serverWg.Wait()
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	healthMutex.RLock()
	currentBotStatus := botStatus
	healthMutex.RUnlock()

	// Determine overall status based on bot status
	overallStatus := "healthy"
	if currentBotStatus == "failed" {
		overallStatus = "degraded"
	}

	response := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.1.0",
		Uptime:    time.Since(startTime).String(),
		BotStatus: currentBotStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode health response", http.StatusInternalServerError)
	}
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ALIVE"))
}
