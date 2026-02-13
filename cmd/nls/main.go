package main

import (
	"context"
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/scanner"
	"nls/internal/ui"
)

func run() error {
	cidr := "192.168.1.0/24"
	if len(os.Args) > 1 {
		cidr = os.Args[1]
	}

	if _, _, err := net.ParseCIDR(cidr); err != nil {
		return fmt.Errorf("invalid CIDR format %s: %w", cidr, err)
	}

	hosts, err := scanner.Scan(context.Background(), cidr)
	if err != nil {
		return fmt.Errorf("scan network: %w", err)
	}

	model := ui.NewUIModel(hosts)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		return fmt.Errorf("run ui: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
