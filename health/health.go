package health

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

func StartHealthServer() {
	startTime := time.Now()
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(startTime)
		
		health := HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   "1.1.0",
			Uptime:    uptime.String(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	})
	
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Use a different port to avoid conflict with massage-bot
	go http.ListenAndServe(":8081", nil)
}
