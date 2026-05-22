package app

import (
	"context"
	"errors"
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

func TestApp_Run_InvalidConfig_EmptyCIDR(t *testing.T) {
	cfg := &Config{CIDR: "", Timeout: 5 * time.Minute}
	a := New(cfg, &mockScanner{})
	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for empty CIDR, got nil")
	}
	if !strings.Contains(err.Error(), "invalid configuration") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApp_Run_InvalidConfig_BadCIDR(t *testing.T) {
	cfg := &Config{CIDR: "not-a-cidr", Timeout: 5 * time.Minute}
	a := New(cfg, &mockScanner{})
	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid CIDR, got nil")
	}
	if !strings.Contains(err.Error(), "invalid configuration") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApp_Run_InvalidConfig_ZeroTimeout(t *testing.T) {
	cfg := &Config{CIDR: "192.168.1.0/24", Timeout: 0}
	a := New(cfg, &mockScanner{})
	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for zero timeout, got nil")
	}
	if !strings.Contains(err.Error(), "invalid configuration") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApp_Run_ScanError(t *testing.T) {
	cfg := &Config{CIDR: "192.168.1.0/24", Timeout: 5 * time.Minute}
	scanErr := errors.New("permission denied")
	a := New(cfg, &mockScanner{err: scanErr})
	err := a.Run(context.Background())
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}
	if !strings.Contains(err.Error(), "scan network") {
		t.Errorf("error should wrap scan error, got: %v", err)
	}
	if !errors.Is(err, scanErr) {
		t.Errorf("error chain should contain original scan error")
	}
}

func TestApp_Run_ContextCancelled(t *testing.T) {
	cfg := &Config{CIDR: "192.168.1.0/24", Timeout: 5 * time.Minute}
	// Scanner that returns context error
	a := New(cfg, &mockScanner{err: context.Canceled})
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancelled
	err := a.Run(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}
