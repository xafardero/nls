package ui

import (
"os"
"strconv"

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
func buildColumns(width int, weights ColumnWeights) []table.Column {
	remaining := width - TableIDWidth - TablePaddingWidth

	ipWidth := int(float64(remaining) * weights.IP)
	macWidth := int(float64(remaining) * weights.MAC)
	vendorWidth := int(float64(remaining) * weights.Vendor)
	hostnameWidth := int(float64(remaining) * weights.Hostname)

	return []table.Column{
		{Title: "Id", Width: TableIDWidth},
		{Title: "IP", Width: ipWidth},
		{Title: "MAC", Width: macWidth},
		{Title: "Vendor", Width: vendorWidth},
		{Title: "Hostname", Width: hostnameWidth},
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
