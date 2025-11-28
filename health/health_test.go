package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(HealthHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
}
