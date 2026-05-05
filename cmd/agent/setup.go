package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"watchup-agent/internal/config"
)

// setupServerID prompts the user to set their server ID if it's empty
func setupServerID(cfg *config.Config) error {
	if cfg.ServerID != "" {
		return nil // Already configured
	}

	// Check if running in non-interactive mode (e.g., systemd service)
	fileInfo, _ := os.Stdin.Stat()
	isInteractive := (fileInfo.Mode() & os.ModeCharDevice) != 0

	if !isInteractive {
		// Running as a service or in non-interactive mode
		fmt.Println()
		fmt.Println("❌ Server ID Configuration Required")
		fmt.Println()
		fmt.Println("The agent cannot start because 'server_id' is not configured.")
		fmt.Println()
		fmt.Println("To configure the server ID:")
		fmt.Println("   1. Edit the config file:")
		fmt.Println("      sudo nano /etc/watchup-agent/config.yaml")
		fmt.Println()
		fmt.Println("   2. Set the 'server_id' field to a unique identifier:")
		fmt.Println("      server_id: \"web-prod-01\"")
		fmt.Println()
		fmt.Println("   3. Restart the agent:")
		fmt.Println("      sudo systemctl restart watchup-agent")
		fmt.Println()
		fmt.Println("📋 Server ID Guidelines:")
		fmt.Println("   • Must be unique within your WatchUp account")
		fmt.Println("   • Use descriptive names like: web-prod-01, db-server-main, api-gateway-1")
		fmt.Println("   • 3-50 characters, alphanumeric and hyphens only")
		fmt.Println()
		return fmt.Errorf("server_id not configured - please edit config.yaml")
	}

	// Interactive mode - prompt for input
	fmt.Println()
	fmt.Println("🔧 Server ID Setup Required")
	fmt.Println("Your agent needs a unique server identifier.")
	fmt.Println()
	fmt.Println("📋 Server ID Guidelines:")
	fmt.Println("   • Must be unique within your WatchUp account")
	fmt.Println("   • Use descriptive names like: web-prod-01, db-server-main, api-gateway-1")
	fmt.Println("   • 3-50 characters, alphanumeric and hyphens only")
	fmt.Println("   • Cannot be changed after linking (choose carefully)")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("Enter server ID: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		serverID := strings.TrimSpace(input)
		if err := validateServerID(serverID); err != nil {
			fmt.Printf("❌ %v Please try again.\n", err)
			continue
		}

		// Update config
		cfg.ServerID = serverID
		
		// Save updated config
		if err := config.Save(cfg, "config.yaml"); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Server ID set to: %s\n", serverID)
		fmt.Println("Configuration saved to config.yaml")
		fmt.Println()
		fmt.Println("🔗 Next: This agent will be linked to your WatchUp account")
		fmt.Println("   • You'll get a web link and code to approve this agent")
		fmt.Println("   • Once approved, metrics will be sent automatically")
		break
	}

	return nil
}

// validateServerID validates the server ID format and requirements
func validateServerID(serverID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	if len(serverID) < 3 {
		return fmt.Errorf("server ID must be at least 3 characters long")
	}

	if len(serverID) > 50 {
		return fmt.Errorf("server ID must be 50 characters or less")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range serverID {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return fmt.Errorf("server ID can only contain letters, numbers, hyphens, and underscores")
		}
	}

	// Cannot start or end with hyphen/underscore
	if serverID[0] == '-' || serverID[0] == '_' || 
	   serverID[len(serverID)-1] == '-' || serverID[len(serverID)-1] == '_' {
		return fmt.Errorf("server ID cannot start or end with hyphen or underscore")
	}

	return nil
}