// Package app provides application-level orchestration and configuration.
// It coordinates the scanner, UI, and other components to execute the
// network scanning workflow.
package app

import (
	"fmt"
	"net"
	"time"
)

// Config holds the application configuration settings.
// It centralizes all configurable parameters for the network scanner.
type Config struct {
	// CIDR is the network range to scan (e.g., "192.168.1.0/24")
	CIDR string

	// Timeout is the maximum duration for the scan operation
	Timeout time.Duration

	// ShowProgress determines whether to display a progress spinner
	ShowProgress bool
}

// DefaultConfig returns a Config with sensible default values.
// CIDR defaults to the common home network range 192.168.1.0/24.
func DefaultConfig() *Config {
	return &Config{
		CIDR:         "192.168.1.0/24",
		Timeout:      5 * time.Minute,
		ShowProgress: true,
	}
}

// Validate checks if the configuration is valid.
// Returns an error if CIDR is invalid or timeout is non-positive.
func (c *Config) Validate() error {
	if _, _, err := net.ParseCIDR(c.CIDR); err != nil {
		return fmt.Errorf("invalid CIDR format %s: %w", c.CIDR, err)
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got %v", c.Timeout)
	}

	return nil
}
