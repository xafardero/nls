package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/scanner"
)

func TestHandleNormalKeys_HelpScreen(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)

	// Press ? to show help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.mode != modeHelp {
		t.Error("expected mode to be modeHelp after '?' key")
	}
}

func TestHandleHelpKeys_Exit(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)
	model.mode = modeHelp

	tests := []struct {
		name string
		key  string
	}{
		{name: "esc to exit", key: "esc"},
		{name: "q to exit", key: "q"},
		{name: "? to exit", key: "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model
			var msg tea.KeyMsg
			if tt.key == "esc" {
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}

			updatedModel, _ := m.Update(msg)
			result := updatedModel.(UIModel)

			if result.mode != modeNormal {
				t.Errorf("expected mode to be modeNormal after %s, got %v", tt.key, result.mode)
			}
		})
	}
}

func TestHandleNormalKeys_SearchMode(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)

	// Press / to activate search
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.mode != modeSearch {
		t.Error("expected mode to be modeSearch after '/' key")
	}
}

func TestHandleSearchKeys_Cancel(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)
	model.mode = modeSearch
	model.searchInput.SetValue("test query")

	// Press esc to cancel
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.mode != modeNormal {
		t.Error("expected mode to be modeNormal after esc in search mode")
	}
	if m.searchInput.Value() != "" {
		t.Errorf("search input should be cleared, got %q", m.searchInput.Value())
	}
}

func TestHandleSearchKeys_ApplyFilter(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Apple", Hostname: "test.local"},
		{ID: 1, IP: "192.168.1.2", MAC: "11:22:33:44:55:66", Vendor: "Samsung", Hostname: "phone.local"},
	}
	model := NewUIModel(hosts)
	model.mode = modeSearch
	model.searchInput.SetValue("apple")

	// Press enter to apply filter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.mode != modeNormal {
		t.Error("expected mode to be modeNormal after enter in search mode")
	}
	if !m.searchActive {
		t.Error("expected searchActive to be true after applying filter")
	}
	if m.searchQuery != "apple" {
		t.Errorf("searchQuery = %q; want %q", m.searchQuery, "apple")
	}
	if len(m.filteredHosts) != 1 {
		t.Errorf("filteredHosts length = %d; want 1", len(m.filteredHosts))
	}
}

func TestHandleNormalKeys_SortColumns(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Vendor A", Hostname: "host1"},
		{ID: 1, IP: "192.168.1.5", MAC: "11:22:33:44:55:66", Vendor: "Vendor B", Hostname: "host2"},
	}
	model := NewUIModel(hosts)

	tests := []struct {
		name            string
		key             string
		expectedSortCol int
		expectedAsc     bool
	}{
		{name: "sort by IP (1)", key: "1", expectedSortCol: 1, expectedAsc: true},
		{name: "sort by MAC (2)", key: "2", expectedSortCol: 2, expectedAsc: true},
		{name: "sort by Vendor (3)", key: "3", expectedSortCol: 3, expectedAsc: true},
		{name: "sort by Hostname (4)", key: "4", expectedSortCol: 4, expectedAsc: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			updatedModel, _ := m.Update(msg)
			result := updatedModel.(UIModel)

			if result.sortColumn != tt.expectedSortCol {
				t.Errorf("sortColumn = %d; want %d", result.sortColumn, tt.expectedSortCol)
			}
			if result.sortAscending != tt.expectedAsc {
				t.Errorf("sortAscending = %v; want %v", result.sortAscending, tt.expectedAsc)
			}
		})
	}
}

func TestHandleNormalKeys_SortToggle(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Vendor A", Hostname: "host1"},
	}
	model := NewUIModel(hosts)
	model.sortColumn = 1
	model.sortAscending = true

	// Press 1 again to toggle sort direction
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.sortColumn != 1 {
		t.Errorf("sortColumn = %d; want 1", m.sortColumn)
	}
	if m.sortAscending {
		t.Error("expected sortAscending to be false after toggle")
	}
}

func TestHandleNormalKeys_CopyMAC(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)

	// Press m to copy MAC
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	if cmd == nil {
		t.Error("expected command for status message tick")
	}
	if m.statusMessage == "" {
		t.Error("expected status message after copying MAC")
	}
}

func TestHandleNormalKeys_CopyHostname(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)

	// Press h to copy hostname
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	if cmd == nil {
		t.Error("expected command for status message tick")
	}
	if m.statusMessage == "" {
		t.Error("expected status message after copying hostname")
	}
}

func TestHandleNormalKeys_CopyAll(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)

	// Press a to copy all fields
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	if cmd == nil {
		t.Error("expected command for status message tick")
	}
	if m.statusMessage == "" {
		t.Error("expected status message after copying all fields")
	}
}

func TestView_HelpMode(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)
	model.mode = modeHelp

	view := model.View()

	if view == "" {
		t.Error("View() returned empty string in help mode")
	}
	// Help text should be visible
	if len(view) < 100 {
		t.Error("Help view seems too short")
	}
}

func TestView_SearchMode(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)
	model.mode = modeSearch

	view := model.View()

	if view == "" {
		t.Error("View() returned empty string in search mode")
	}
}

func TestView_NormalModeWithFilter(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Apple", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)
	model.searchActive = true
	model.searchQuery = "apple"

	view := model.View()

	if view == "" {
		t.Error("View() returned empty string")
	}
	// Should show filter indicator in footer
	if len(view) < 10 {
		t.Error("View with filter seems too short")
	}
}
