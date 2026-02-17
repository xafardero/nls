package ui

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"golang.org/x/term"

	"nls/internal/scanner"
)

// ColumnWeights defines the proportional width allocation for table columns.
// Values represent the percentage of available width each column should occupy.
type ColumnWeights struct {
	IP       float64
	MAC      float64
	Vendor   float64
	Hostname float64
}

// DefaultColumnWeights returns the standard column width distribution.
// IP gets 20%, while MAC, Vendor, and Hostname each get approximately 26.67%.
func DefaultColumnWeights() ColumnWeights {
	return ColumnWeights{
		IP:       0.20,
		MAC:      0.27,
		Vendor:   0.26,
		Hostname: 0.27,
	}
}

// getTerminalSize returns the current terminal width and height.
// Falls back to environment variables or default values
// if terminal size detection fails.
func getTerminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = DefaultTermWidth
		height = DefaultTermHeight
		if envW, ok := os.LookupEnv("COLUMNS"); ok {
			if val, err := strconv.Atoi(envW); err == nil {
				width = val
			}
		}
		if envH, ok := os.LookupEnv("LINES"); ok {
			if val, err := strconv.Atoi(envH); err == nil {
				height = val - DefaultTermHeightPad
			}
		}
	} else {
		width = w
		height = h - DefaultTermHeightPad
	}
	return
}

// buildColumns creates table column definitions based on terminal width.
// Columns are proportionally sized using the provided weights.
// If sortCol > 0, adds a sort indicator (↑/↓) to the sorted column's title.
func buildColumns(width int, weights ColumnWeights, sortCol int, ascending bool) []table.Column {
	remaining := width - TableIDWidth - TablePaddingWidth

	ipWidth := int(float64(remaining) * weights.IP)
	macWidth := int(float64(remaining) * weights.MAC)
	vendorWidth := int(float64(remaining) * weights.Vendor)
	hostnameWidth := int(float64(remaining) * weights.Hostname)

	// Helper to add sort indicator
	addSortIndicator := func(title string, col int) string {
		if col == sortCol {
			if ascending {
				return title + " ↑"
			}
			return title + " ↓"
		}
		return title
	}

	return []table.Column{
		{Title: "Id", Width: TableIDWidth},
		{Title: addSortIndicator("IP", 1), Width: ipWidth},
		{Title: addSortIndicator("MAC", 2), Width: macWidth},
		{Title: addSortIndicator("Vendor", 3), Width: vendorWidth},
		{Title: addSortIndicator("Hostname", 4), Width: hostnameWidth},
	}
}

// buildRows converts a slice of HostInfo into table rows.
// Returns a single "No hosts found" row if the input is empty.
func buildRows(hosts []scanner.HostInfo) []table.Row {
	if len(hosts) == 0 {
		return []table.Row{{"-", "No hosts found", "-", "-", "-"}}
	}

	rows := make([]table.Row, 0, len(hosts))
	for _, h := range hosts {
		rows = append(rows, table.Row{
			strconv.Itoa(h.ID),
			h.IP,
			h.MAC,
			h.Vendor,
			h.Hostname,
		})
	}
	return rows
}

// filterHosts returns a filtered slice of hosts matching the search query.
// The query is matched case-insensitively against IP, MAC, Vendor, and Hostname fields.
func filterHosts(hosts []scanner.HostInfo, query string) []scanner.HostInfo {
	if query == "" {
		return hosts
	}

	query = strings.ToLower(query)
	filtered := make([]scanner.HostInfo, 0)

	for _, h := range hosts {
		if strings.Contains(strings.ToLower(h.IP), query) ||
			strings.Contains(strings.ToLower(h.MAC), query) ||
			strings.Contains(strings.ToLower(h.Vendor), query) ||
			strings.Contains(strings.ToLower(h.Hostname), query) {
			filtered = append(filtered, h)
		}
	}

	return filtered
}

// sortHosts returns a sorted copy of hosts based on the specified column.
// col: 1=IP, 2=MAC, 3=Vendor, 4=Hostname
func sortHosts(hosts []scanner.HostInfo, col int, ascending bool) []scanner.HostInfo {
	if col == 0 || len(hosts) == 0 {
		return hosts
	}

	// Create a copy to avoid modifying original
	sorted := make([]scanner.HostInfo, len(hosts))
	copy(sorted, hosts)

	sort.Slice(sorted, func(i, j int) bool {
		var less bool
		switch col {
		case 1: // IP - compare numerically
			less = compareIPs(sorted[i].IP, sorted[j].IP)
		case 2: // MAC
			less = strings.Compare(sorted[i].MAC, sorted[j].MAC) < 0
		case 3: // Vendor
			less = strings.Compare(sorted[i].Vendor, sorted[j].Vendor) < 0
		case 4: // Hostname
			less = strings.Compare(sorted[i].Hostname, sorted[j].Hostname) < 0
		default:
			return false
		}

		if ascending {
			return less
		}
		return !less
	})

	return sorted
}

// compareIPs compares two IP addresses numerically.
// Returns true if ip1 < ip2.
func compareIPs(ip1, ip2 string) bool {
	// Handle "none" sentinel value
	if ip1 == "none" && ip2 == "none" {
		return false
	}
	if ip1 == "none" {
		return false // "none" sorts last
	}
	if ip2 == "none" {
		return true
	}

	// Split into octets
	parts1 := strings.Split(ip1, ".")
	parts2 := strings.Split(ip2, ".")

	// Compare each octet numerically
	for i := 0; i < 4 && i < len(parts1) && i < len(parts2); i++ {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		if err1 != nil || err2 != nil {
			// Fall back to string comparison if not valid numbers
			return strings.Compare(ip1, ip2) < 0
		}

		if num1 != num2 {
			return num1 < num2
		}
	}

	// If all octets equal, IPs are equal
	return false
}
