package health

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	// Set bot status to healthy first
	SetBotStatus("running")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check for expected JSON content using contains
	if !strings.Contains(rr.Body.String(), `"status":"healthy"`) {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"bot_status":"running"`) {
		t.Errorf("handler missing bot status: got %v", rr.Body.String())
	}
}

func TestHealthHandlerDegraded(t *testing.T) {
	// Set bot status to failed to test degraded status
	SetBotStatus("failed")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Should show degraded status when bot is failed
	if !strings.Contains(rr.Body.String(), `"status":"degraded"`) {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"bot_status":"failed"`) {
		t.Errorf("handler missing failed bot status: got %v", rr.Body.String())
	}
}

func TestReadyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLiveHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/live", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(liveHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "ALIVE"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSetBotStatus(t *testing.T) {
	// Test that SetBotStatus actually changes the status
	SetBotStatus("testing")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	// Should contain the bot status we set
	if !strings.Contains(rr.Body.String(), `"bot_status":"testing"`) {
		t.Errorf("handler returned unexpected bot status: got %v", rr.Body.String())
	}
}

func TestHealthServerLifecycle(t *testing.T) {
	// Test that we can start and shutdown the server
	err := StartServer("18081", nil) // Use a different port to avoid conflicts
	if err != nil {
		t.Errorf("Failed to start health server: %v", err)
	}

	// Give it a moment to start
	time.Sleep(200 * time.Millisecond)

	// Test that the server is responding
	resp, err := http.Get("http://localhost:18081/health")
	if err != nil {
		t.Errorf("Health server not responding: %v", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %v", resp.StatusCode)
		}
	}

	// Shutdown the server
	Shutdown()
}
