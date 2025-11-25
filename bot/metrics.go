package bot

import (
    "fmt"
    "net/http"
    "time"
)

var (
    totalUpdatesTriggered int
    lastUpdateTime       time.Time
    botStartTime         = time.Now()
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    
    metrics := fmt.Sprintf(`# HELP watchtower_updates_total Total number of updates triggered
# TYPE watchtower_updates_total counter
watchtower_updates_total %d

# HELP watchtower_last_update_time_seconds Timestamp of last update
# TYPE watchtower_last_update_time_seconds gauge
watchtower_last_update_time_seconds %d

# HELP watchtower_bot_uptime_seconds Bot uptime in seconds
# TYPE watchtower_bot_uptime_seconds gauge
watchtower_bot_uptime_seconds %.0f
`,
        totalUpdatesTriggered,
        lastUpdateTime.Unix(),
        time.Since(botStartTime).Seconds(),
    )
    
    w.Write([]byte(metrics))
}

func RecordUpdateTriggered() {
    totalUpdatesTriggered++
    lastUpdateTime = time.Now()
}
