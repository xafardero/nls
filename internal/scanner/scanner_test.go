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
					IP:       "192.168.1.1",
					MAC:      "00:11:22:33:44:55",
					Vendor:   "Router Co",
					Hostname: "router.local",
				},
				{
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

func TestExtractHostInfo_SlicePreallocation(t *testing.T) {
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

	if results[0].IP != "192.168.1.1" {
		t.Errorf("first host IP = %q; want %q", results[0].IP, "192.168.1.1")
	}
	if results[999].IP != "192.168.1.1" {
		t.Errorf("last host IP = %q; want %q", results[999].IP, "192.168.1.1")
	}
}
