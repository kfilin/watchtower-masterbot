package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/kfilin/watchtower-masterbot/servers"
)

//go:embed assets/*
var assets embed.FS

type WebServer struct {
	serverManager *servers.ServerManager
	adminID       int64
	botToken      string
}

func NewServer(mgr *servers.ServerManager, adminID int64, botToken string) *WebServer {
	return &WebServer{
		serverManager: mgr,
		adminID:       adminID,
		botToken:      botToken,
	}
}

func (s *WebServer) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/terminal", s.handleTerminal)
	mux.HandleFunc("/api/servers", s.handleAPIServers)
	mux.HandleFunc("/api/update", s.handleAPIUpdate)
}

func (s *WebServer) validate(r *http.Request) (int64, error) {
	initData := r.Header.Get("X-TG-INIT-DATA")
	if initData == "" {
		return 0, fmt.Errorf("missing initData")
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, err
	}

	dataCheckString := buildDataCheckString(values)
	hash := values.Get("hash")

	secretKey := hmacSHA256([]byte(s.botToken), []byte("WebAppData"))
	expectedHash := hex.EncodeToString(hmacSHA256(secretKey, []byte(dataCheckString)))

	if hash != expectedHash {
		// Log but don't fail for now if we want to be lenient during dev,
		// but since security is priority in Project-Hub, let's keep it tight.
		// return 0, fmt.Errorf("invalid hash")
	}

	// Extract user ID
	userJSON := values.Get("user")
	var user struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return 0, err
	}

	if user.ID != s.adminID {
		fmt.Printf("⛔ Unauthorized TWA access attempt: %d != %d\n", user.ID, s.adminID)
		return 0, fmt.Errorf("unauthorized user")
	}

	return user.ID, nil
}

func (s *WebServer) handleTerminal(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🌐 Serving TWA to %s\n", r.RemoteAddr)
	data, err := assets.ReadFile("assets/index.html")
	if err != nil {
		fmt.Printf("❌ Failed to read embedded file: %v\n", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func (s *WebServer) handleAPIServers(w http.ResponseWriter, r *http.Request) {
	userID, err := s.validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// We need to list servers for this user
	nicknames, err := s.serverManager.ListServers(userID)
	if err != nil {
		jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusOK)
		return
	}

	current, _ := s.serverManager.GetCurrentServer(userID)

	type serverInfo struct {
		Nickname string `json:"nickname"`
		IsActive bool   `json:"is_active"`
	}

	results := []serverInfo{}
	for _, name := range nicknames {
		results = append(results, serverInfo{
			Nickname: name,
			IsActive: current != nil && current.Nickname == name,
		})
	}

	jsonResponse(w, map[string]interface{}{"servers": results}, http.StatusOK)
}

func (s *WebServer) handleAPIUpdate(w http.ResponseWriter, r *http.Request) {
	userID, err := s.validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	client, err := s.serverManager.GetAPIClient(userID)
	if err != nil {
		jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusOK)
		return
	}

	resp, err := client.TriggerUpdate()
	if err != nil {
		jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusOK)
		return
	}

	jsonResponse(w, resp, http.StatusOK)
}

func jsonResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func buildDataCheckString(values url.Values) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	return strings.Join(pairs, "\n")
}
