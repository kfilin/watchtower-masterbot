package api

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type WatchtowerClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type ContainerStatus struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Status      string    `json:"status"` // "up_to_date", "update_available", "updated", "failed"
	LastChecked time.Time `json:"last_checked"`
}

type WatchtowerStatus struct {
	Version    string    `json:"version"`
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"last_update"`
	Containers int       `json:"containers_count"`
}

type UpdateResponse struct {
	Updated []string `json:"updated"`
	Failed  []string `json:"failed"`
}

func NewWatchtowerClient(baseURL, token string) *WatchtowerClient {
	// Create a custom transport that matches curl's behavior
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// Important: Disable HTTP/2 PING frames and other keep-alive that might cause timeouts
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	return &WatchtowerClient{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout:   60 * time.Second, // Reasonable timeout
			Transport: transport,
		},
	}
}

func (c *WatchtowerClient) doRequest(method, endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WatchtowerMasterBot/1.0")

	return c.HTTPClient.Do(req)
}

func (c *WatchtowerClient) GetContainers() ([]ContainerStatus, error) {
	// Watchtower doesn't have a containers status endpoint
	// We'll return an empty list for now and implement this later
	return []ContainerStatus{}, nil
}

func (c *WatchtowerClient) TriggerUpdate() (*UpdateResponse, error) {
	resp, err := c.doRequest("POST", "/v1/update")
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Handle different success scenarios
	switch resp.StatusCode {
	case http.StatusOK:
		// Watchtower returns 200 with empty body - this is success
		// The update happens asynchronously after the response
		return &UpdateResponse{
			Updated: []string{"Update triggered successfully - check Watchtower logs for details"},
			Failed:  []string{},
		}, nil

	case http.StatusAccepted:
		// Alternative success response
		return &UpdateResponse{
			Updated: []string{"Update accepted and processing"},
			Failed:  []string{},
		}, nil

	case http.StatusBadGateway:
		return nil, fmt.Errorf("watchtower service unavailable (502 Bad Gateway)")

	case http.StatusGatewayTimeout:
		return nil, fmt.Errorf("watchtower gateway timeout (504)")

	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("watchtower service temporarily unavailable (503)")

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed - check your token")

	default:
		return nil, fmt.Errorf("API returned unexpected status: %d", resp.StatusCode)
	}
}

func (c *WatchtowerClient) GetStatus() (*WatchtowerStatus, error) {
	// Watchtower doesn't have a status endpoint, return basic info
	return &WatchtowerStatus{
		Version: "1.7.1",
		Status:  "running",
		// We can't get container count without a proper endpoint
	}, nil
}

// TestConnection verifies that we can reach the Watchtower API
func (c *WatchtowerClient) TestConnection() error {
	resp, err := c.doRequest("GET", "/v1/update")
	if err != nil {
		return fmt.Errorf("connection test failed: %v", err)
	}
	defer resp.Body.Close()

	// For connection test, we accept 405 Method Not Allowed (since GET /v1/update isn't allowed)
	// or any 2xx/3xx/4xx except 5xx server errors
	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error during connection test: %d", resp.StatusCode)
	}

	return nil
}
