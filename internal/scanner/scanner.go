package scanner

import (
	"context"
)

// Scanner defines the interface for network scanning operations.
// Implementations can use different scanning tools (nmap, custom, etc.)
// or provide mock implementations for testing.
type Scanner interface {
	// Scan performs a network scan on the specified CIDR target.
	// Returns a list of discovered hosts or an error if the scan fails.
	// The context can be used to cancel the scan operation.
	Scan(ctx context.Context, target string) ([]HostInfo, error)
}
