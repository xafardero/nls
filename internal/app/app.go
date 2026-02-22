package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/scanner"
	"nls/internal/ui"
)

// App represents the main application orchestrator.
// It coordinates the scanning and UI components.
type App struct {
	config  *Config
	scanner scanner.Scanner
}

// New creates a new App instance with the provided configuration and scanner.
// The scanner parameter allows for dependency injection and testing with mock implementations.
func New(config *Config, s scanner.Scanner) *App {
	return &App{
		config:  config,
		scanner: s,
	}
}

// Run executes the main application workflow:
// 1. Validates configuration
// 2. Performs network scan
// 3. Launches interactive UI with results
//
// Returns an error if validation, scanning, or UI execution fails.
func (a *App) Run(ctx context.Context) error {
	if err := a.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	hosts, err := a.scanner.Scan(ctx, a.config.CIDR)
	if err != nil {
		return fmt.Errorf("scan network: %w", err)
	}

	model := ui.NewUIModel(hosts, a.scanner, a.config.CIDR)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("run ui: %w", err)
	}

	return nil
}
