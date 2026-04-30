package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// DeviceCodeResponse represents the response from device registration
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// StatusResponse represents the response from status polling
type StatusResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

// DeviceInfo contains information about the agent device
type DeviceInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
	ServerID string `json:"server_id"`
}

// RegisterDevice initiates the device linking flow
func RegisterDevice(baseURL, deviceName, version, serverID string) (*DeviceCodeResponse, error) {
	// Get system information
	hostname, _ := getHostname()
	deviceInfo := DeviceInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  version,
		Hostname: hostname,
		ServerID: serverID, // Include server_id for backend validation
	}

	payload := map[string]interface{}{
		"device_name": deviceName,
		"device_info": deviceInfo,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Post(baseURL+"/agents/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 409 {
		return nil, fmt.Errorf("server ID '%s' already exists for your account. Use a different server_id or deactivate the existing agent", serverID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	var result DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// PollForToken polls the status endpoint until the device is approved or expires
func PollForToken(baseURL, deviceCode string, interval time.Duration, maxDuration time.Duration) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	startTime := time.Now()
	
	for {
		// Check if we've exceeded the maximum duration
		if time.Since(startTime) > maxDuration {
			return "", fmt.Errorf("polling timeout exceeded")
		}

		// Make status request
		url := fmt.Sprintf("%s/agents/status?device_code=%s", baseURL, deviceCode)
		resp, err := client.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to check status: %w", err)
		}

		var status StatusResponse
		err = json.NewDecoder(resp.Body).Decode(&status)
		resp.Body.Close()

		if err != nil {
			return "", fmt.Errorf("failed to decode status response: %w", err)
		}

		switch status.Status {
		case "approved":
			return status.AccessToken, nil
		case "denied":
			return "", fmt.Errorf("device linking was denied by user")
		case "expired":
			return "", fmt.Errorf("device code has expired")
		case "pending":
			// Continue polling
			time.Sleep(interval)
			continue
		default:
			return "", fmt.Errorf("unknown status: %s", status.Status)
		}
	}
}

// getHostname returns the system hostname
func getHostname() (string, error) {
	// This is a simple implementation - could be enhanced
	return "agent-host", nil
}