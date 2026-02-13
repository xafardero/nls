package ui

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"nls/internal/scanner"
)

func TestBuildColumns(t *testing.T) {
	weights := DefaultColumnWeights()

	tests := []struct {
		name  string
		width int
		want  []table.Column
	}{
		{
			name:  "standard terminal width",
			width: 100,
			want: []table.Column{
				{Title: "Id", Width: 5},
				{Title: "IP", Width: 17},
				{Title: "MAC", Width: 23},
				{Title: "Vendor", Width: 22},
				{Title: "Hostname", Width: 23},
			},
		},
		{
			name:  "narrow terminal",
			width: 50,
			want: []table.Column{
				{Title: "Id", Width: 5},
				{Title: "IP", Width: 7},
				{Title: "MAC", Width: 9},
				{Title: "Vendor", Width: 9},
				{Title: "Hostname", Width: 9},
			},
		},
		{
			name:  "wide terminal",
			width: 200,
			want: []table.Column{
				{Title: "Id", Width: 5},
				{Title: "IP", Width: 37},
				{Title: "MAC", Width: 50},
				{Title: "Vendor", Width: 48},
				{Title: "Hostname", Width: 50},
			},
		},
		{
			name:  "minimum width",
			width: 20,
			want: []table.Column{
				{Title: "Id", Width: 5},
				{Title: "IP", Width: 1},
				{Title: "MAC", Width: 1},
				{Title: "Vendor", Width: 1},
				{Title: "Hostname", Width: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildColumns(tt.width, weights)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildColumns(%d, weights) mismatch:\ngot:  %+v\nwant: %+v", tt.width, got, tt.want)
			}
		})
	}
}

func TestDefaultColumnWeights(t *testing.T) {
	weights := DefaultColumnWeights()

	// Verify weights are set correctly
	if weights.IP != 0.20 {
		t.Errorf("IP weight = %f; want 0.20", weights.IP)
	}
	if weights.MAC != 0.27 {
		t.Errorf("MAC weight = %f; want 0.27", weights.MAC)
	}
	if weights.Vendor != 0.26 {
		t.Errorf("Vendor weight = %f; want 0.26", weights.Vendor)
	}
	if weights.Hostname != 0.27 {
		t.Errorf("Hostname weight = %f; want 0.27", weights.Hostname)
	}

	// Verify weights sum to approximately 1.0 (100%)
	sum := weights.IP + weights.MAC + weights.Vendor + weights.Hostname
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("weights sum = %f; want approximately 1.0", sum)
	}
}

func TestBuildColumns_CustomWeights(t *testing.T) {
	customWeights := ColumnWeights{
		IP:       0.30,
		MAC:      0.30,
		Vendor:   0.20,
		Hostname: 0.20,
	}

	width := 100
	columns := buildColumns(width, customWeights)

	// Verify custom weights are applied
	// remaining = 100 - 5 - 8 = 87
	// IP: 87 * 0.30 = 26.1 → 26
	// MAC: 87 * 0.30 = 26.1 → 26
	// Vendor: 87 * 0.20 = 17.4 → 17
	// Hostname: 87 * 0.20 = 17.4 → 17
	expected := []table.Column{
		{Title: "Id", Width: TableIDWidth},
		{Title: "IP", Width: 26},
		{Title: "MAC", Width: 26},
		{Title: "Vendor", Width: 17},
		{Title: "Hostname", Width: 17},
	}

	if !reflect.DeepEqual(columns, expected) {
		t.Errorf("buildColumns with custom weights mismatch:\ngot:  %+v\nwant: %+v", columns, expected)
	}
}

func TestBuildRows(t *testing.T) {
	tests := []struct {
		name  string
		hosts []scanner.HostInfo
		want  []table.Row
	}{
		{
			name: "single host",
			hosts: []scanner.HostInfo{
				{
					ID:       0,
					IP:       "192.168.1.10",
					MAC:      "AA:BB:CC:DD:EE:FF",
					Vendor:   "Apple Inc.",
					Hostname: "macbook.local",
				},
			},
			want: []table.Row{
				{"0", "192.168.1.10", "AA:BB:CC:DD:EE:FF", "Apple Inc.", "macbook.local"},
			},
		},
		{
			name: "multiple hosts",
			hosts: []scanner.HostInfo{
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
			want: []table.Row{
				{"0", "192.168.1.1", "00:11:22:33:44:55", "Router Co", "router.local"},
				{"1", "192.168.1.2", "AA:BB:CC:DD:EE:00", "Device Inc", "device.local"},
			},
		},
		{
			name:  "empty host list",
			hosts: []scanner.HostInfo{},
			want: []table.Row{
				{"-", "No hosts found", "-", "-", "-"},
			},
		},
		{
			name: "hosts with 'none' values",
			hosts: []scanner.HostInfo{
				{
					ID:       0,
					IP:       "192.168.1.100",
					MAC:      "none",
					Vendor:   "none",
					Hostname: "none",
				},
			},
			want: []table.Row{
				{"0", "192.168.1.100", "none", "none", "none"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRows(tt.hosts)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildRows() mismatch:\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}

func TestBuildRows_Preallocation(t *testing.T) {
	// Test efficient handling of large host lists
	hosts := make([]scanner.HostInfo, 1000)
	for i := range hosts {
		hosts[i] = scanner.HostInfo{
			ID:       i,
			IP:       "192.168.1.1",
			MAC:      "AA:BB:CC:DD:EE:FF",
			Vendor:   "Test Vendor",
			Hostname: "test.local",
		}
	}

	rows := buildRows(hosts)

	if len(rows) != 1000 {
		t.Errorf("expected 1000 rows, got %d", len(rows))
	}

	// Spot check a few rows
	if rows[0][0] != "0" {
		t.Errorf("first row ID = %s; want %s", rows[0][0], "0")
	}
	if rows[999][0] != "999" {
		t.Errorf("last row ID = %s; want %s", rows[999][0], "999")
	}
}

func TestGetBaseStyle(t *testing.T) {
	style := getBaseStyle()

	// Verify style has expected properties
	if style.GetBorderStyle() != lipgloss.NormalBorder() {
		t.Error("expected NormalBorder style")
	}

	// Verify it's a valid lipgloss.Style
	rendered := style.Render("test")
	if rendered == "" {
		t.Error("style.Render() returned empty string")
	}
}

func TestGetPromptStyle(t *testing.T) {
	style := getPromptStyle()

	// Verify style has expected properties
	if style.GetBorderStyle() != lipgloss.RoundedBorder() {
		t.Error("expected RoundedBorder style")
	}

	// Verify width is set correctly
	if style.GetWidth() != SSHPromptWidth {
		t.Errorf("style width = %d; want %d", style.GetWidth(), SSHPromptWidth)
	}

	// Verify it's a valid lipgloss.Style
	rendered := style.Render("test")
	if rendered == "" {
		t.Error("style.Render() returned empty string")
	}
}

func TestNewUIModel(t *testing.T) {
	tests := []struct {
		name  string
		hosts []scanner.HostInfo
	}{
		{
			name: "with hosts",
			hosts: []scanner.HostInfo{
				{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
			},
		},
		{
			name:  "empty hosts",
			hosts: []scanner.HostInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewUIModel(tt.hosts)

			// Verify model is initialized
			if model.table.Cursor() < 0 {
				t.Error("table cursor not initialized")
			}

			// Verify showPrompt starts as false
			if model.showPrompt {
				t.Error("showPrompt should be false initially")
			}

			// Verify selectedIP is empty
			if model.selectedIP != "" {
				t.Errorf("selectedIP = %q; want empty string", model.selectedIP)
			}

			// Verify username input is configured
			if model.usernameInput.Placeholder != "username" {
				t.Errorf("username placeholder = %q; want %q", model.usernameInput.Placeholder, "username")
			}

			if model.usernameInput.CharLimit != SSHUsernameMaxLen {
				t.Errorf("username CharLimit = %d; want %d", model.usernameInput.CharLimit, SSHUsernameMaxLen)
			}
		})
	}
}

func TestUIModel_Init(t *testing.T) {
	model := NewUIModel([]scanner.HostInfo{})
	cmd := model.Init()

	// Init should return nil (no initial commands)
	if cmd != nil {
		t.Errorf("Init() returned non-nil command: %v", cmd)
	}
}

func TestUIModel_View(t *testing.T) {
	tests := []struct {
		name       string
		showPrompt bool
		selectedIP string
	}{
		{
			name:       "normal view",
			showPrompt: false,
		},
		{
			name:       "prompt view",
			showPrompt: true,
			selectedIP: "192.168.1.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewUIModel([]scanner.HostInfo{
				{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test"},
			})
			model.showPrompt = tt.showPrompt
			model.selectedIP = tt.selectedIP

			view := model.View()

			if view == "" {
				t.Error("View() returned empty string")
			}

			if tt.showPrompt {
				// Prompt view should contain the IP
				if model.selectedIP != "" && len(view) < len(model.selectedIP) {
					t.Error("prompt view should contain selected IP")
				}
			} else {
				// Normal view should contain footer
				if len(view) == 0 {
					t.Error("normal view should not be empty")
				}
			}
		})
	}
}
