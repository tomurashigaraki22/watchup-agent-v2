package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"watchup-agent/internal/client"
	"watchup-agent/internal/config"
	"watchup-agent/internal/metrics"
)

// Agent represents the main agent with its components
type Agent struct {
	config         *config.Config
	metricsClient  *client.Client
	collector      *metrics.Collector
	tracker        *client.MetricsTracker
	isRunning      bool
}

// NewAgent creates a new agent instance
func NewAgent(cfg *config.Config, apiClient *client.Client) *Agent {
	return &Agent{
		config:        cfg,
		metricsClient: apiClient,
		collector:     metrics.NewCollector(cfg),
		tracker:       client.NewMetricsTracker(),
		isRunning:     false,
	}
}

// Start begins the main agent loop
func (a *Agent) Start(ctx context.Context) error {
	if a.isRunning {
		return fmt.Errorf("agent is already running")
	}

	a.isRunning = true
	log.Printf("Starting metrics collection loop (interval: %v)", a.config.Interval)

	// Create ticker for metrics collection
	ticker := time.NewTicker(a.config.Interval)
	defer ticker.Stop()

	// Send initial metrics immediately
	a.collectAndSend()

	// Main loop
	for {
		select {
		case <-ctx.Done():
			log.Printf("Agent loop stopping due to context cancellation")
			a.isRunning = false
			return ctx.Err()

		case <-ticker.C:
			a.collectAndSend()
		}
	}
}

// collectAndSend collects metrics and sends them to the server
func (a *Agent) collectAndSend() {
	log.Printf("Collecting metrics...")

	// Collect metrics
	payload, err := a.collector.CollectAll()
	if err != nil {
		log.Printf("Failed to collect metrics: %v", err)
		a.tracker.RecordFailure(err)
		return
	}

	log.Printf("Collected metrics for %d categories", len(payload.Metrics))

	// Send metrics with retry
	err = a.metricsClient.SendMetricsWithRetry(payload, 5) // Max 5 retries
	if err != nil {
		log.Printf("Failed to send metrics: %v", err)
		a.tracker.RecordFailure(err)

		// Check if it's an auth error
		if client.IsAuthError(err) {
			log.Printf("Authentication error detected - agent may need re-authentication")
			// In a production system, this could trigger re-authentication
		}
		return
	}

	log.Printf("Metrics sent successfully")
	a.tracker.RecordSuccess()

	// Log statistics periodically
	stats := a.tracker.GetStats()
	if (stats.TotalSent+stats.TotalFailed)%10 == 0 && stats.TotalSent+stats.TotalFailed > 0 {
		log.Printf("Metrics stats: %d sent, %d failed, %.1f%% success rate", 
			stats.TotalSent, stats.TotalFailed, stats.SuccessRate)
	}
}

// Stop gracefully stops the agent
func (a *Agent) Stop() {
	log.Printf("Stopping agent...")
	a.isRunning = false
}

// IsRunning returns whether the agent is currently running
func (a *Agent) IsRunning() bool {
	return a.isRunning
}

// GetStats returns current metrics statistics
func (a *Agent) GetStats() client.MetricsStats {
	return a.tracker.GetStats()
}