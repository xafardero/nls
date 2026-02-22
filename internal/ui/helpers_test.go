package ui

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/table"

	"nls/internal/scanner"
)

func TestBuildColumns_DefaultWeights(t *testing.T) {
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
				{Title: "IP", Width: 18},
				{Title: "MAC", Width: 24},
				{Title: "Vendor", Width: 23},
				{Title: "Hostname", Width: 24},
			},
		},
		{
			name:  "narrow terminal",
			width: 50,
			want: []table.Column{
				{Title: "IP", Width: 8},
				{Title: "MAC", Width: 11},
				{Title: "Vendor", Width: 10},
				{Title: "Hostname", Width: 11},
			},
		},
		{
			name:  "wide terminal",
			width: 200,
			want: []table.Column{
				{Title: "IP", Width: 38},
				{Title: "MAC", Width: 51},
				{Title: "Vendor", Width: 49},
				{Title: "Hostname", Width: 51},
			},
		},
		{
			name:  "minimum width",
			width: 20,
			want: []table.Column{
				{Title: "IP", Width: 2},
				{Title: "MAC", Width: 3},
				{Title: "Vendor", Width: 3},
				{Title: "Hostname", Width: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildColumns(tt.width, weights, 0, false)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildColumns(%d, weights) mismatch:\ngot:  %+v\nwant: %+v", tt.width, got, tt.want)
			}
		})
	}
}

func TestDefaultColumnWeights(t *testing.T) {
	weights := DefaultColumnWeights()

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
	columns := buildColumns(width, customWeights, 0, false)

	// remaining = 100 - 8 = 92
	// IP: 92 * 0.30 = 27.6 → 27
	// MAC: 92 * 0.30 = 27.6 → 27
	// Vendor: 92 * 0.20 = 18.4 → 18
	// Hostname: 92 * 0.20 = 18.4 → 18
	expected := []table.Column{
		{Title: "IP", Width: 27},
		{Title: "MAC", Width: 27},
		{Title: "Vendor", Width: 18},
		{Title: "Hostname", Width: 18},
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
				{"192.168.1.10", "AA:BB:CC:DD:EE:FF", "Apple Inc.", "macbook.local"},
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
				{"192.168.1.1", "00:11:22:33:44:55", "Router Co", "router.local"},
				{"192.168.1.2", "AA:BB:CC:DD:EE:00", "Device Inc", "device.local"},
			},
		},
		{
			name:  "empty host list",
			hosts: []scanner.HostInfo{},
			want: []table.Row{
				{"No hosts found", "-", "-", "-"},
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
				{"192.168.1.100", "none", "none", "none"},
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

	if rows[0][0] != "192.168.1.1" {
		t.Errorf("first row IP = %s; want %s", rows[0][0], "192.168.1.1")
	}
	if rows[999][0] != "192.168.1.1" {
		t.Errorf("last row IP = %s; want %s", rows[999][0], "192.168.1.1")
	}
}

func TestBaseStyle(t *testing.T) {
	rendered := baseStyle.Render("test")
	if rendered == "" {
		t.Error("style.Render() returned empty string")
	}
}

func TestPromptStyle(t *testing.T) {
	if promptStyle.GetWidth() != SSHPromptWidth {
		t.Errorf("style width = %d; want %d", promptStyle.GetWidth(), SSHPromptWidth)
	}

	rendered := promptStyle.Render("test")
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
			model := NewUIModel(tt.hosts, nil, "")

			if model.table.Cursor() < 0 {
				t.Error("table cursor not initialized")
			}

			if model.mode != modeNormal {
				t.Error("mode should be modeNormal initially")
			}

			if model.selectedIP != "" {
				t.Errorf("selectedIP = %q; want empty string", model.selectedIP)
			}

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
	model := NewUIModel([]scanner.HostInfo{}, nil, "")
	cmd := model.Init()

	if cmd == nil {
		t.Error("Init() should return a WindowSize command, got nil")
	}
}

func TestUIModel_View(t *testing.T) {
	tests := []struct {
		name       string
		mode       viewMode
		selectedIP string
	}{
		{
			name: "normal view",
			mode: modeNormal,
		},
		{
			name:       "ssh prompt view",
			mode:       modeSSHPrompt,
			selectedIP: "192.168.1.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewUIModel([]scanner.HostInfo{
				{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test"},
			}, nil, "")
			model.mode = tt.mode
			model.selectedIP = tt.selectedIP

			view := model.View()

			if view == "" {
				t.Error("View() returned empty string")
			}

			if tt.mode == modeSSHPrompt {
				if model.selectedIP != "" && len(view) < len(model.selectedIP) {
					t.Error("prompt view should contain selected IP")
				}
			} else {
				if len(view) == 0 {
					t.Error("normal view should not be empty")
				}
			}
		})
	}
}
