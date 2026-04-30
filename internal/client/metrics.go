package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"watchup-agent/internal/metrics"
)

// SendMetricsWithRetry sends metrics with exponential backoff retry logic
func (c *Client) SendMetricsWithRetry(payload *metrics.MetricsPayload, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := c.SendMetrics(payload)
		if err == nil {
			// Success
			if attempt > 0 {
				log.Printf("Metrics sent successfully after %d retries", attempt)
			}
			return nil
		}

		lastErr = err

		// Check if it's an auth error - don't retry these
		if IsAuthError(err) {
			return fmt.Errorf("authentication error (not retrying): %w", err)
		}

		// Don't sleep on the last attempt
		if attempt < maxRetries {
			// Exponential backoff: 1s, 2s, 4s, 8s, 16s
			backoffDuration := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			log.Printf("Metrics send failed (attempt %d/%d): %v. Retrying in %v...", 
				attempt+1, maxRetries+1, err, backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	return fmt.Errorf("failed to send metrics after %d attempts: %w", maxRetries+1, lastErr)
}

// SendMetrics sends a metrics payload to the server (updated from client.go)
func (c *Client) SendMetrics(payload *metrics.MetricsPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	resp, err := c.makeRequest("POST", "/metrics", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("authentication failed - token may be expired")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("metrics submission failed with status %d", resp.StatusCode)
	}

	return nil
}

// MetricsStats tracks metrics sending statistics
type MetricsStats struct {
	TotalSent     int64     `json:"total_sent"`
	TotalFailed   int64     `json:"total_failed"`
	LastSentAt    time.Time `json:"last_sent_at"`
	LastFailedAt  time.Time `json:"last_failed_at"`
	LastError     string    `json:"last_error"`
	SuccessRate   float64   `json:"success_rate"`
}

// MetricsTracker tracks metrics sending statistics
type MetricsTracker struct {
	stats MetricsStats
}

// NewMetricsTracker creates a new metrics tracker
func NewMetricsTracker() *MetricsTracker {
	return &MetricsTracker{
		stats: MetricsStats{},
	}
}

// RecordSuccess records a successful metrics send
func (mt *MetricsTracker) RecordSuccess() {
	mt.stats.TotalSent++
	mt.stats.LastSentAt = time.Now()
	mt.updateSuccessRate()
}

// RecordFailure records a failed metrics send
func (mt *MetricsTracker) RecordFailure(err error) {
	mt.stats.TotalFailed++
	mt.stats.LastFailedAt = time.Now()
	mt.stats.LastError = err.Error()
	mt.updateSuccessRate()
}

// GetStats returns current statistics
func (mt *MetricsTracker) GetStats() MetricsStats {
	return mt.stats
}

// updateSuccessRate calculates the success rate
func (mt *MetricsTracker) updateSuccessRate() {
	total := mt.stats.TotalSent + mt.stats.TotalFailed
	if total > 0 {
		mt.stats.SuccessRate = float64(mt.stats.TotalSent) / float64(total) * 100
	}
}