package auth

import (
	"fmt"
	"os"
	"strings"
)

// SaveToken saves the authentication token to a file with secure permissions
func SaveToken(token, filepath string) error {
	// Ensure directory exists
	dir := filepath[:strings.LastIndex(filepath, "/")]
	if dir != "" {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create token directory: %w", err)
		}
	}

	// Write token to file with restricted permissions (owner read/write only)
	if err := os.WriteFile(filepath, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// LoadToken loads the authentication token from a file
func LoadToken(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("token file does not exist")
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("token file is empty")
	}

	return token, nil
}

// TokenExists checks if a token file exists and is readable
func TokenExists(filepath string) bool {
	if _, err := os.Stat(filepath); err != nil {
		return false
	}

	// Try to read the token to ensure it's valid
	token, err := LoadToken(filepath)
	return err == nil && token != ""
}

// DeleteToken removes the token file
func DeleteToken(filepath string) error {
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}
	return nil
}

// ValidateTokenFormat performs basic validation on token format
func ValidateTokenFormat(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	// Basic JWT format check (should have 3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("token does not appear to be a valid JWT")
	}

	return nil
}