package app

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"nls/internal/scanner"
)

type mockScanner struct {
	hosts []scanner.HostInfo
	err   error
}

func (m *mockScanner) Scan(_ context.Context, _ string) ([]scanner.HostInfo, error) {
	return m.hosts, m.err
}

func TestApp_Run_RequiresRoot(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("test runs as root; skipping non-root privilege check test")
	}
	cfg := &Config{CIDR: "192.168.1.0/24", Timeout: 5 * time.Minute}
	app := New(cfg, &mockScanner{})
	err := app.Run(context.Background())
	if err == nil {
		t.Fatal("expected error when not root, got nil")
	}
	if !strings.Contains(err.Error(), "root") {
		t.Errorf("error should mention root, got: %v", err)
	}
}

func TestApp_Run_InvalidConfig(t *testing.T) {
	cfg := &Config{CIDR: "", Timeout: 5 * time.Minute}
	app := New(cfg, &mockScanner{})
	err := app.Run(context.Background())
	if err == nil {
		t.Fatal("expected validation error for empty CIDR, got nil")
	}
	if !strings.Contains(err.Error(), "invalid configuration") {
		t.Errorf("error should mention invalid configuration, got: %v", err)
	}
}
