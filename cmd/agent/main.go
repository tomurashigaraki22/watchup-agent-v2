package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"watchup-agent/internal/agent"
	"watchup-agent/internal/auth"
	"watchup-agent/internal/client"
	"watchup-agent/internal/config"
)

const Version = "1.0.0"

func main() {
	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("WatchUp Agent v%s starting...", Version)

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("Server ID: %s", cfg.ServerID)
	log.Printf("Endpoint: %s", cfg.Endpoint)
	log.Printf("Interval: %s", cfg.Interval)

	// Validate server_id is set
	if cfg.ServerID == "" {
		log.Printf("Server ID not configured, starting setup...")
		if err := setupServerID(cfg); err != nil {
			log.Fatalf("Server ID setup failed: %v", err)
		}
	}

	// Phase 1: Authentication
	log.Printf("Initializing authentication...")
	authFlow := auth.NewAuthFlow(cfg.Endpoint, cfg.ServerID, Version, cfg.ServerID, cfg.Auth.TokenFile)
	
	token, err := authFlow.EnsureAuthenticated()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Create authenticated HTTP client
	apiClient := client.NewClient(cfg.Endpoint, token)
	
	// Validate token with server
	log.Printf("Validating token with server...")
	if err := apiClient.ValidateToken(); err != nil {
		log.Printf("Token validation failed: %v", err)
		log.Printf("Starting re-authentication...")
		
		// Invalidate stored token and retry
		authFlow.InvalidateToken()
		token, err = authFlow.EnsureAuthenticated()
		if err != nil {
			log.Fatalf("Re-authentication failed: %v", err)
		}
		
		// Update client with new token
		apiClient = client.NewClient(cfg.Endpoint, token)
		if err := apiClient.ValidateToken(); err != nil {
			log.Fatalf("Token validation still failing: %v", err)
		}
	}

	log.Printf("Authentication successful!")

	// Phase 2 & 3: Create and start the agent
	log.Printf("Initializing metrics collection...")
	agentInstance := agent.NewAgent(cfg, apiClient)

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start agent in a goroutine
	agentDone := make(chan error, 1)
	go func() {
		agentDone <- agentInstance.Start(ctx)
	}()

	fmt.Println()
	fmt.Println("🚀 WatchUp Agent is running!")
	fmt.Printf("   Server ID: %s\n", cfg.ServerID)
	fmt.Printf("   Endpoint: %s\n", cfg.Endpoint)
	fmt.Printf("   Metrics interval: %s\n", cfg.Interval)
	fmt.Printf("   Token file: %s\n", cfg.Auth.TokenFile)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the agent.")
	fmt.Println()

	// Wait for shutdown signal or agent error
	select {
	case <-sigChan:
		log.Println("Shutdown signal received, stopping agent...")
		cancel() // Cancel the context to stop the agent
		
		// Wait for agent to finish
		<-agentDone
		
	case err := <-agentDone:
		if err != nil && err != context.Canceled {
			log.Printf("Agent stopped with error: %v", err)
		}
	}

	// Print final statistics
	stats := agentInstance.GetStats()
	if stats.TotalSent > 0 || stats.TotalFailed > 0 {
		fmt.Println()
		fmt.Printf("📊 Final Statistics:\n")
		fmt.Printf("   Metrics sent: %d\n", stats.TotalSent)
		fmt.Printf("   Failed sends: %d\n", stats.TotalFailed)
		fmt.Printf("   Success rate: %.1f%%\n", stats.SuccessRate)
		if !stats.LastSentAt.IsZero() {
			fmt.Printf("   Last successful send: %s\n", stats.LastSentAt.Format("2006-01-02 15:04:05"))
		}
	}

	log.Println("Agent stopped successfully")
}
