package auth

import (
	"fmt"
	"log"
	"time"
)

// AuthFlow manages the complete authentication flow
type AuthFlow struct {
	BaseURL     string
	DeviceName  string
	Version     string
	ServerID    string
	TokenFile   string
}

// NewAuthFlow creates a new authentication flow manager
func NewAuthFlow(baseURL, deviceName, version, serverID, tokenFile string) *AuthFlow {
	return &AuthFlow{
		BaseURL:    baseURL,
		DeviceName: deviceName,
		Version:    version,
		ServerID:   serverID,
		TokenFile:  tokenFile,
	}
}

// EnsureAuthenticated ensures the agent is authenticated, performing the flow if needed
func (af *AuthFlow) EnsureAuthenticated() (string, error) {
	// Check if we already have a valid token
	if TokenExists(af.TokenFile) {
		token, err := LoadToken(af.TokenFile)
		if err == nil && ValidateTokenFormat(token) == nil {
			log.Printf("Found existing token, validating...")
			return token, nil
		}
		log.Printf("Existing token is invalid, starting new authentication flow")
	}

	// Start device linking flow
	log.Printf("Starting device linking authentication...")
	return af.performDeviceLinking()
}

// performDeviceLinking executes the complete device linking flow
func (af *AuthFlow) performDeviceLinking() (string, error) {
	// Step 1: Register device
	log.Printf("Registering device with server...")
	deviceResp, err := RegisterDevice(af.BaseURL, af.DeviceName, af.Version, af.ServerID)
	if err != nil {
		return "", fmt.Errorf("device registration failed: %w", err)
	}

	// Step 2: Display instructions to user
	af.displayLinkingInstructions(deviceResp)

	// Step 3: Poll for approval
	log.Printf("Waiting for device approval...")
	pollInterval := time.Duration(deviceResp.Interval) * time.Second
	maxDuration := time.Duration(deviceResp.ExpiresIn) * time.Second
	
	token, err := PollForToken(af.BaseURL, deviceResp.DeviceCode, pollInterval, maxDuration)
	if err != nil {
		return "", fmt.Errorf("token polling failed: %w", err)
	}

	// Step 4: Save token
	log.Printf("Device approved! Saving authentication token...")
	if err := SaveToken(token, af.TokenFile); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	log.Printf("Authentication completed successfully!")
	return token, nil
}

// displayLinkingInstructions shows the user how to link the device
func (af *AuthFlow) displayLinkingInstructions(resp *DeviceCodeResponse) {
	// Override verification URL to use the correct domain
	verificationURL := "https://v2.watchup.site/agent-link"
	
	fmt.Println()
	fmt.Println("🔗 Link this agent to your Watchup account:")
	fmt.Println()
	fmt.Printf("   Visit: %s?code=%s\n", verificationURL, resp.UserCode)
	fmt.Println()
	fmt.Printf("   Code expires in %d minutes\n", resp.ExpiresIn/60)
	fmt.Println()
	fmt.Println("Waiting for authorization...")
	fmt.Println("(Press Ctrl+C to cancel)")
	fmt.Println()
}

// InvalidateToken removes the stored token (for logout/reset)
func (af *AuthFlow) InvalidateToken() error {
	log.Printf("Invalidating stored token...")
	return DeleteToken(af.TokenFile)
}