package scanner

import (
	"reflect"
	"testing"

	"github.com/Ullaakut/nmap/v3"
)

func TestExtractHostInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    *nmap.Run
		expected []HostInfo
	}{
		{
			name: "single host with all fields",
			input: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{Addr: "192.168.1.10", AddrType: "ipv4"},
							{Addr: "AA:BB:CC:DD:EE:FF", Vendor: "Apple Inc.", AddrType: "mac"},
						},
						Hostnames: []nmap.Hostname{
							{Name: "macbook.local"},
						},
					},
				},
			},
			expected: []HostInfo{
				{
					ID:       0,
					IP:       "192.168.1.10",
					MAC:      "AA:BB:CC:DD:EE:FF",
					Vendor:   "Apple Inc.",
					Hostname: "macbook.local",
				},
			},
		},
		{
			name: "multiple hosts",
			input: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{Addr: "192.168.1.1", AddrType: "ipv4"},
							{Addr: "00:11:22:33:44:55", Vendor: "Router Co", AddrType: "mac"},
						},
						Hostnames: []nmap.Hostname{
							{Name: "router.local"},
						},
					},
					{
						Addresses: []nmap.Address{
							{Addr: "192.168.1.2", AddrType: "ipv4"},
							{Addr: "AA:BB:CC:DD:EE:00", Vendor: "Device Inc", AddrType: "mac"},
						},
						Hostnames: []nmap.Hostname{
							{Name: "device.local"},
						},
					},
				},
			},
			expected: []HostInfo{
				{
					ID:       0,
					IP:       "192.168.1.1",
					MAC:      "00:11:22:33:44:55",
					Vendor:   "Router Co",
					Hostname: "router.local",
				},
				{
					ID:       1,
					IP:       "192.168.1.2",
					MAC:      "AA:BB:CC:DD:EE:00",
					Vendor:   "Device Inc",
					Hostname: "device.local",
				},
			},
		},
		{
			name: "host with IP only",
			input: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{Addr: "192.168.1.100", AddrType: "ipv4"},
						},
					},
				},
			},
			expected: []HostInfo{
				{
					ID:       0,
					IP:       "192.168.1.100",
					MAC:      "none",
					Vendor:   "none",
					Hostname: "none",
				},
			},
		},
		{
			name: "host with no addresses",
			input: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{},
					},
				},
			},
			expected: []HostInfo{
				{
					ID:       0,
					IP:       "none",
					MAC:      "none",
					Vendor:   "none",
					Hostname: "none",
				},
			},
		},
		{
			name: "empty scan result",
			input: &nmap.Run{
				Hosts: []nmap.Host{},
			},
			expected: []HostInfo{},
		},
		{
			name: "host with multiple hostnames - takes first",
			input: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{Addr: "192.168.1.50", AddrType: "ipv4"},
						},
						Hostnames: []nmap.Hostname{
							{Name: "primary.local"},
							{Name: "secondary.local"},
						},
					},
				},
			},
			expected: []HostInfo{
				{
					ID:       0,
					IP:       "192.168.1.50",
					MAC:      "none",
					Vendor:   "none",
					Hostname: "primary.local",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractHostInfo(tt.input)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("extractHostInfo() mismatch:\ngot:  %+v\nwant: %+v", got, tt.expected)
			}
		})
	}
}

func TestExtractHostInfo_IDSequential(t *testing.T) {
	// Verify that IDs are assigned sequentially starting from 0
	scanResult := &nmap.Run{
		Hosts: []nmap.Host{
			{Addresses: []nmap.Address{{Addr: "192.168.1.1"}}},
			{Addresses: []nmap.Address{{Addr: "192.168.1.2"}}},
			{Addresses: []nmap.Address{{Addr: "192.168.1.3"}}},
			{Addresses: []nmap.Address{{Addr: "192.168.1.4"}}},
			{Addresses: []nmap.Address{{Addr: "192.168.1.5"}}},
		},
	}

	results := extractHostInfo(scanResult)

	if len(results) != 5 {
		t.Fatalf("expected 5 hosts, got %d", len(results))
	}

	for i, host := range results {
		if host.ID != i {
			t.Errorf("host %d has ID %d; want %d", i, host.ID, i)
		}
	}
}

func TestExtractHostInfo_SlicePreallocation(t *testing.T) {
	// Test that the function handles large result sets efficiently
	// This is a behavioral test - we can't directly test preallocation
	// but we verify it handles 1000 hosts without issue
	hosts := make([]nmap.Host, 1000)
	for i := range hosts {
		hosts[i] = nmap.Host{
			Addresses: []nmap.Address{
				{Addr: "192.168.1.1", AddrType: "ipv4"},
			},
		}
	}

	scanResult := &nmap.Run{Hosts: hosts}
	results := extractHostInfo(scanResult)

	if len(results) != 1000 {
		t.Errorf("expected 1000 hosts, got %d", len(results))
	}

	// Verify all have sequential IDs
	for i, result := range results {
		if result.ID != i {
			t.Errorf("host at index %d has ID %d", i, result.ID)
			break
		}
	}
}
