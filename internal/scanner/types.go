// Package scanner provides network scanning functionality using nmap.
// It performs ping scans to discover active hosts on a network and
// extracts information about discovered devices including IP addresses,
// MAC addresses, vendor information, and hostnames.
package scanner

// HostInfo represents information about a discovered network host.
// All string fields use "none" as a sentinel value when information
// is not available.
type HostInfo struct {
	ID       int
	IP       string
	MAC      string
	Vendor   string
	Hostname string
}
