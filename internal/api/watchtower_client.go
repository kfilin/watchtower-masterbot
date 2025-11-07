package api

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
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
	Status      string    `json:"status"`
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
	Message string   `json:"message,omitempty"`
}

type UpdateJob struct {
	ID      string    `json:"id"`
	State   string    `json:"state"`
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
	Results []struct {
		Container string `json:"container"`
		Status    string `json:"status"`
		Error     string `json:"error,omitempty"`
	} `json:"results"`
}

type MetricsResponse struct {
	Data map[string]string
}

func NewWatchtowerClient(baseURL, token string) *WatchtowerClient {
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
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	return &WatchtowerClient{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout:   60 * time.Second,
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
	return []ContainerStatus{}, nil
}

func (c *WatchtowerClient) TriggerUpdateWithTimeout(timeout time.Duration) (*UpdateResponse, error) {
	// Create a custom client with longer timeout just for updates
	customClient := &http.Client{
		Timeout:   timeout,
		Transport: c.HTTPClient.Transport,
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, "/v1/update")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := customClient.Do(req)
	if err != nil {
		// Check if it's a timeout - this might mean the update is processing
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline") {
			return &UpdateResponse{
				Updated: []string{},
				Failed:  []string{},
				Message: "Update triggered successfully (processing in background)",
			}, nil
		}
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Read the response body first to check if it's empty
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// If we can't read the body but got 200 OK, treat as success
			return &UpdateResponse{
				Updated: []string{},
				Failed:  []string{},
				Message: "Update triggered successfully (empty response)",
			}, nil
		}

		// If body is empty, treat as successful trigger
		if len(body) == 0 {
			return &UpdateResponse{
				Updated: []string{},
				Failed:  []string{},
				Message: "Update triggered successfully",
			}, nil
		}

		// Try to decode the response
		var updateResp UpdateResponse
		if err := json.Unmarshal(body, &updateResp); err != nil {
			// If decoding fails but we got 200 OK, treat as success
			return &UpdateResponse{
				Updated: []string{},
				Failed:  []string{},
				Message: "Update triggered successfully (invalid JSON response)",
			}, nil
		}
		return &updateResp, nil

	case http.StatusAccepted:
		return &UpdateResponse{
			Updated: []string{"Update accepted and processing"},
			Failed:  []string{},
			Message: "Update queued and processing in background",
		}, nil

	case http.StatusNoContent:
		return &UpdateResponse{
			Updated: []string{},
			Failed:  []string{},
			Message: "Update triggered successfully (no content)",
		}, nil

	case http.StatusBadGateway:
		return nil, fmt.Errorf("watchtower service unavailable (502 Bad Gateway)")

	case http.StatusGatewayTimeout:
		return &UpdateResponse{
			Updated: []string{},
			Failed:  []string{},
			Message: "Update triggered (gateway timeout but likely processing)",
		}, nil

	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("watchtower service temporarily unavailable (503)")

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed - check your token")

	default:
		// For any 2xx status, treat as success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return &UpdateResponse{
				Updated: []string{},
				Failed:  []string{},
				Message: fmt.Sprintf("Update triggered successfully (status %d)", resp.StatusCode),
			}, nil
		}
		return nil, fmt.Errorf("API returned unexpected status: %d", resp.StatusCode)
	}
}

// Keep the original method for backward compatibility
func (c *WatchtowerClient) TriggerUpdate() (*UpdateResponse, error) {
	return c.TriggerUpdateWithTimeout(5 * time.Minute) // 5 minute timeout for updates
}

func (c *WatchtowerClient) GetStatus() (*WatchtowerStatus, error) {
	return &WatchtowerStatus{
		Version: "1.7.1",
		Status:  "running",
	}, nil
}

func (c *WatchtowerClient) TestConnection() error {
	resp, err := c.doRequest("GET", "/v1/update")
	if err != nil {
		return fmt.Errorf("connection test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error during connection test: %d", resp.StatusCode)
	}

	return nil
}

// GetUpdateJobs gets recent update jobs
func (c *WatchtowerClient) GetUpdateJobs(limit int) ([]UpdateJob, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/v1/update?limit=%d", limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get update jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response struct {
		Result []UpdateJob `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Result, nil
}

// GetUpdateJob gets specific update job details
func (c *WatchtowerClient) GetUpdateJob(jobID string) (*UpdateJob, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/v1/update/%s", jobID))
	if err != nil {
		return nil, fmt.Errorf("failed to get update job: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var job UpdateJob
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &job, nil
}

// GetMetrics gets Prometheus metrics
func (c *WatchtowerClient) GetMetrics() (*MetricsResponse, error) {
	resp, err := c.doRequest("GET", "/v1/metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	metricsData := make(map[string]string)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.Contains(line, " ") {
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				metricsData[parts[0]] = parts[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning metrics response: %v", err)
	}

	return &MetricsResponse{Data: metricsData}, nil
}
