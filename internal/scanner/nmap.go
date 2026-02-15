package scanner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"

	"nls/internal/progress"
)

// NmapScanner implements the Scanner interface using nmap for network discovery.
// It performs ping scans to detect active hosts and extract their information.
type NmapScanner struct {
	progress progress.Reporter
	logger   *log.Logger
}

// NewNmapScanner creates a new NmapScanner with the provided progress reporter.
// If progress reporter is nil, a no-op reporter is used.
func NewNmapScanner(p progress.Reporter) *NmapScanner {
	if p == nil {
		p = progress.NoOp{}
	}
	return &NmapScanner{
		progress: p,
		logger:   log.Default(),
	}
}

// Scan performs an nmap ping scan on the specified CIDR target and returns
// a list of discovered hosts. The scan respects the provided context for
// cancellation.
//
// The function displays progress feedback during the scan and extracts
// IP addresses, MAC addresses, vendor information, and hostnames from the results.
//
// Returns an error if the scanner cannot be created or if the scan fails.
func (s *NmapScanner) Scan(ctx context.Context, target string) ([]HostInfo, error) {
	s.progress.Start("Scanning network...")

	// Buffered channels prevent goroutine leaks on context cancellation
	resultCh := make(chan *nmap.Run, 1)
	errCh := make(chan error, 1)

	go func() {
		scanner, err := nmap.NewScanner(
			ctx,
			nmap.WithTargets(target),
			nmap.WithPingScan(),
		)
		if err != nil {
			errCh <- fmt.Errorf("create scanner: %w", err)
			return
		}

		result, warnings, err := scanner.Run()
		if len(*warnings) > 0 {
			s.logger.Printf("run finished with warnings: %s\n", *warnings)
		}
		if err != nil {
			errCh <- fmt.Errorf("run scan: %w", err)
			return
		}
		resultCh <- result
	}()

	for {
		select {
		case result := <-resultCh:
			s.progress.Finish()
			return extractHostInfo(result), nil
		case err := <-errCh:
			s.progress.Finish()
			return nil, err
		case <-ctx.Done():
			s.progress.Finish()
			return nil, ctx.Err()
		default:
			s.progress.Update()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// extractHostInfo converts nmap scan results into a slice of HostInfo structs.
// It extracts the first IP address, MAC address with vendor, and hostname
// from each host. Missing fields are set to "none".
// IDs are assigned sequentially starting from 0.
func extractHostInfo(scanResult *nmap.Run) []HostInfo {
	hosts := make([]HostInfo, 0, len(scanResult.Hosts))
	for i, host := range scanResult.Hosts {
		ip := "none"
		mac := "none"
		vendor := "none"
		hostname := "none"

		if len(host.Addresses) > 0 {
			ip = host.Addresses[0].Addr
		}
		if len(host.Addresses) > 1 {
			mac = host.Addresses[1].Addr
			vendor = host.Addresses[1].Vendor
		}
		if len(host.Hostnames) > 0 {
			hostname = host.Hostnames[0].Name
		}

		hosts = append(hosts, HostInfo{
			ID:       i,
			IP:       ip,
			MAC:      mac,
			Vendor:   vendor,
			Hostname: hostname,
		})
	}
	return hosts
}
