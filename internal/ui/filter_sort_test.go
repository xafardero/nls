package ui

import (
	"reflect"
	"testing"

	"nls/internal/scanner"
)

func TestFilterHosts(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Apple Inc.", Hostname: "macbook.local"},
		{ID: 1, IP: "192.168.1.10", MAC: "11:22:33:44:55:66", Vendor: "Samsung", Hostname: "phone.local"},
		{ID: 2, IP: "10.0.0.1", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Router Co", Hostname: "router"},
		{ID: 3, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Apple Inc.", Hostname: "ipad"},
	}

	tests := []struct {
		name     string
		query    string
		expected []scanner.HostInfo
	}{
		{
			name:     "empty query returns all",
			query:    "",
			expected: hosts,
		},
		{
			name:  "filter by IP",
			query: "10.0.0",
			expected: []scanner.HostInfo{
				{ID: 2, IP: "10.0.0.1", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Router Co", Hostname: "router"},
			},
		},
		{
			name:  "filter by MAC",
			query: "aa:bb:cc",
			expected: []scanner.HostInfo{
				{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Apple Inc.", Hostname: "macbook.local"},
			},
		},
		{
			name:  "filter by Vendor",
			query: "apple",
			expected: []scanner.HostInfo{
				{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Apple Inc.", Hostname: "macbook.local"},
				{ID: 3, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Apple Inc.", Hostname: "ipad"},
			},
		},
		{
			name:  "filter by Hostname",
			query: "router",
			expected: []scanner.HostInfo{
				{ID: 2, IP: "10.0.0.1", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Router Co", Hostname: "router"},
			},
		},
		{
			name:  "case insensitive search",
			query: "APPLE",
			expected: []scanner.HostInfo{
				{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Apple Inc.", Hostname: "macbook.local"},
				{ID: 3, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Apple Inc.", Hostname: "ipad"},
			},
		},
		{
			name:     "no matches",
			query:    "nonexistent",
			expected: []scanner.HostInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterHosts(hosts, tt.query)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("filterHosts() mismatch:\ngot:  %+v\nwant: %+v", result, tt.expected)
			}
		})
	}
}

func TestSortHosts(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"},
		{ID: 1, IP: "192.168.1.5", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Apple", Hostname: "device1"},
		{ID: 2, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Samsung", Hostname: "device2"},
	}

	tests := []struct {
		name      string
		col       int
		ascending bool
		expected  []scanner.HostInfo
	}{
		{
			name:      "no sort (col=0)",
			col:       0,
			ascending: true,
			expected:  hosts, // unchanged
		},
		{
			name:      "sort by IP ascending",
			col:       1,
			ascending: true,
			expected: []scanner.HostInfo{
				{ID: 1, IP: "192.168.1.5", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Apple", Hostname: "device1"},
				{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"},
				{ID: 2, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Samsung", Hostname: "device2"},
			},
		},
		{
			name:      "sort by IP descending",
			col:       1,
			ascending: false,
			expected: []scanner.HostInfo{
				{ID: 2, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Samsung", Hostname: "device2"},
				{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"},
				{ID: 1, IP: "192.168.1.5", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Apple", Hostname: "device1"},
			},
		},
		{
			name:      "sort by MAC ascending",
			col:       2,
			ascending: true,
			expected: []scanner.HostInfo{
				{ID: 1, IP: "192.168.1.5", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Apple", Hostname: "device1"},
				{ID: 2, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Samsung", Hostname: "device2"},
				{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"},
			},
		},
		{
			name:      "sort by Vendor ascending",
			col:       3,
			ascending: true,
			expected: []scanner.HostInfo{
				{ID: 1, IP: "192.168.1.5", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Apple", Hostname: "device1"},
				{ID: 2, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Samsung", Hostname: "device2"},
				{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"},
			},
		},
		{
			name:      "sort by Hostname descending",
			col:       4,
			ascending: false,
			expected: []scanner.HostInfo{
				{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"},
				{ID: 2, IP: "192.168.1.20", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Samsung", Hostname: "device2"},
				{ID: 1, IP: "192.168.1.5", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Apple", Hostname: "device1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortHosts(hosts, tt.col, tt.ascending)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("sortHosts() mismatch:\ngot:  %+v\nwant: %+v", result, tt.expected)
			}

			// Verify original is unchanged
			if !reflect.DeepEqual(hosts[0], scanner.HostInfo{ID: 0, IP: "192.168.1.10", MAC: "CC:CC:CC:CC:CC:CC", Vendor: "Zebra", Hostname: "device3"}) {
				t.Error("sortHosts modified original slice")
			}
		})
	}
}

func TestCompareIPs(t *testing.T) {
	tests := []struct {
		name     string
		ip1      string
		ip2      string
		expected bool // true if ip1 < ip2
	}{
		{
			name:     "simple less than",
			ip1:      "192.168.1.1",
			ip2:      "192.168.1.2",
			expected: true,
		},
		{
			name:     "simple greater than",
			ip1:      "192.168.1.10",
			ip2:      "192.168.1.5",
			expected: false,
		},
		{
			name:     "equal IPs",
			ip1:      "192.168.1.1",
			ip2:      "192.168.1.1",
			expected: false,
		},
		{
			name:     "different octets",
			ip1:      "192.168.1.255",
			ip2:      "192.168.2.1",
			expected: true,
		},
		{
			name:     "string vs numeric comparison",
			ip1:      "192.168.1.9",
			ip2:      "192.168.1.10",
			expected: true, // 9 < 10 numerically (but "9" > "10" as strings)
		},
		{
			name:     "none value sorts last",
			ip1:      "none",
			ip2:      "192.168.1.1",
			expected: false,
		},
		{
			name:     "IP beats none",
			ip1:      "192.168.1.1",
			ip2:      "none",
			expected: true,
		},
		{
			name:     "both none",
			ip1:      "none",
			ip2:      "none",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareIPs(tt.ip1, tt.ip2)

			if result != tt.expected {
				t.Errorf("compareIPs(%q, %q) = %v; want %v", tt.ip1, tt.ip2, result, tt.expected)
			}
		})
	}
}

func TestBuildColumns_WithSortIndicator(t *testing.T) {
	weights := DefaultColumnWeights()
	width := 100

	tests := []struct {
		name      string
		sortCol   int
		ascending bool
		wantTitle string
		colIndex  int
	}{
		{
			name:      "sort by IP ascending",
			sortCol:   1,
			ascending: true,
			wantTitle: "IP ↑",
			colIndex:  0,
		},
		{
			name:      "sort by IP descending",
			sortCol:   1,
			ascending: false,
			wantTitle: "IP ↓",
			colIndex:  0,
		},
		{
			name:      "sort by MAC ascending",
			sortCol:   2,
			ascending: true,
			wantTitle: "MAC ↑",
			colIndex:  1,
		},
		{
			name:      "no sort indicator on other columns",
			sortCol:   1,
			ascending: true,
			wantTitle: "Vendor", // Should not have indicator
			colIndex:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns := buildColumns(width, weights, tt.sortCol, tt.ascending)

			if columns[tt.colIndex].Title != tt.wantTitle {
				t.Errorf("column[%d].Title = %q; want %q", tt.colIndex, columns[tt.colIndex].Title, tt.wantTitle)
			}
		})
	}
}

func TestUpdateModel_RebuildTable(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Vendor A", Hostname: "host1"},
		{ID: 1, IP: "192.168.1.5", MAC: "BB:BB:BB:BB:BB:BB", Vendor: "Vendor B", Hostname: "host2"},
	}

	model := NewUIModel(hosts, nil, "")
	model.sortColumn = 1
	model.sortAscending = true

	// Rebuild should apply sort
	model = model.rebuildTable()

	// Check that table rows are sorted (first row should be 192.168.1.5 after sort)
	rows := model.table.Rows()
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}

	// After sort by IP ascending, first row should be 192.168.1.5
	if rows[0][0] != "192.168.1.5" {
		t.Errorf("first row IP = %q; want %q after sort", rows[0][0], "192.168.1.5")
	}
}
