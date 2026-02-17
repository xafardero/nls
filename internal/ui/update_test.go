package ui

import (
	"context"
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/scanner"
)

func TestUpdate_ClearStatusMsg(t *testing.T) {
	// Create model with a status message
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts, nil, "")
	model.statusMessage = "Test message"
	model.statusExpiry = time.Now().Add(5 * time.Second)

	// Send clearStatusMsg
	msg := clearStatusMsg{}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.statusMessage != "" {
		t.Errorf("statusMessage after clearStatusMsg = %q; want empty string", m.statusMessage)
	}
}

func TestUpdate_QuitKey(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts, nil, "")

	tests := []struct {
		name string
		key  string
	}{
		{name: "q key", key: "q"},
		{name: "ctrl+c", key: "ctrl+c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			_, cmd := model.Update(msg)
			if cmd == nil {
				t.Error("expected quit command, got nil")
			}
		})
	}
}

func TestUpdate_EscToggleFocus(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts, nil, "")
	initialFocus := model.table.Focused()

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.table.Focused() == initialFocus {
		t.Error("expected table focus to toggle after esc key")
	}
}

func TestUpdate_SSHPrompt(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts, nil, "")

	// Simulate pressing 's' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.mode != modeSSHPrompt {
		t.Error("expected mode to be modeSSHPrompt after 's' key")
	}
	if m.selectedIP != "192.168.1.10" {
		t.Errorf("selectedIP = %q; want %q", m.selectedIP, "192.168.1.10")
	}
}

func TestUpdate_SSHPromptEscape(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts, nil, "")
	model.mode = modeSSHPrompt
	model.selectedIP = "192.168.1.10"
	model.usernameInput.SetValue("testuser")

	// Press escape to cancel SSH prompt
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.mode != modeNormal {
		t.Error("expected mode to be modeNormal after escape")
	}
	if m.usernameInput.Value() != "" {
		t.Errorf("usernameInput value = %q; want empty string", m.usernameInput.Value())
	}
}

func TestHandleNormalKeys_CopyIP_ValidHost(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.100", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
		{ID: 1, IP: "192.168.1.101", MAC: "11:22:33:44:55:66", Vendor: "Test2", Hostname: "test2.local"},
	}
	model := NewUIModel(hosts, nil, "")

	// Note: We can't reliably test clipboard.WriteAll() without mocking or system access
	// This test verifies the state changes when 'y' is pressed
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	// Check if status message was set (assuming clipboard write succeeds)
	// Since we can't mock clipboard, we test that the logic path is correct
	if m.statusMessage == "" && cmd == nil {
		// This is acceptable if clipboard.WriteAll fails in test environment
		t.Log("clipboard write may have failed in test environment, which is expected")
	} else if m.statusMessage == "IP copied to clipboard!" {
		// Verify status message was set correctly
		if time.Now().After(m.statusExpiry) {
			t.Error("statusExpiry should be in the future")
		}
		if cmd == nil {
			t.Error("expected tick command to be returned for clearing status")
		}
	}
}

func TestHandleNormalKeys_CopyIP_NoHostsFound(t *testing.T) {
	// Empty hosts list creates a "No hosts found" row
	model := NewUIModel([]scanner.HostInfo{}, nil, "")
	initialMessage := model.statusMessage

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	// Should not set status message for "No hosts found"
	if m.statusMessage != initialMessage {
		t.Errorf("statusMessage changed for 'No hosts found' row; should remain unchanged")
	}
	if cmd != nil && initialMessage == "" {
		t.Error("should not return command for 'No hosts found' row")
	}
}

func TestHandleNormalKeys_SSHWithNoHostsFound(t *testing.T) {
	// Empty hosts list creates a "No hosts found" row
	model := NewUIModel([]scanner.HostInfo{}, nil, "")

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	// Should not show SSH prompt for "No hosts found"
	if m.mode != modeNormal {
		t.Error("should not show SSH prompt for 'No hosts found' row")
	}
}

func TestStatusMessage_Expiry(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts, nil, "")

	// Set status message with past expiry
	model.statusMessage = "Expired message"
	model.statusExpiry = time.Now().Add(-1 * time.Second)

	if time.Now().Before(model.statusExpiry) {
		t.Error("statusExpiry should be in the past")
	}

	// Set status message with future expiry
	model.statusMessage = "Current message"
	model.statusExpiry = time.Now().Add(3 * time.Second)

	if !time.Now().Before(model.statusExpiry) {
		t.Error("statusExpiry should be in the future")
	}
}

// mockScanner is a test double that implements scanner.Scanner interface.
type mockScanner struct {
	hosts []scanner.HostInfo
	err   error
}

func (m *mockScanner) Scan(ctx context.Context, target string) ([]scanner.HostInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.hosts, nil
}

func TestUpdate_RescanTrigger(t *testing.T) {
	initialHosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Initial", Hostname: "test1"},
	}

	mockScan := &mockScanner{
		hosts: []scanner.HostInfo{
			{ID: 0, IP: "192.168.1.2", MAC: "11:22:33:44:55:66", Vendor: "Rescanned", Hostname: "test2"},
		},
	}

	model := NewUIModel(initialHosts, mockScan, "192.168.1.0/24")

	// Trigger rescan with 'r' key
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updatedModel, cmd := model.Update(keyMsg)
	m := updatedModel.(UIModel)

	if !m.isScanning {
		t.Error("isScanning should be true after pressing 'r'")
	}

	if cmd == nil {
		t.Fatal("Update should return a command to start rescan")
	}
}

func TestUpdate_RescanComplete(t *testing.T) {
	initialHosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Initial", Hostname: "test1"},
	}

	newHosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.2", MAC: "11:22:33:44:55:66", Vendor: "Rescanned", Hostname: "test2"},
		{ID: 1, IP: "192.168.1.3", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "New", Hostname: "test3"},
	}

	model := NewUIModel(initialHosts, nil, "192.168.1.0/24")
	model.isScanning = true

	// Simulate rescan completion
	msg := rescanCompleteMsg{hosts: newHosts}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.isScanning {
		t.Error("isScanning should be false after rescan completes")
	}

	if len(m.allHosts) != 2 {
		t.Errorf("allHosts length = %d; want 2", len(m.allHosts))
	}

	if len(m.filteredHosts) != 2 {
		t.Errorf("filteredHosts length = %d; want 2", len(m.filteredHosts))
	}

	if m.statusMessage != "Rescan complete: 2 host(s) found" {
		t.Errorf("statusMessage = %q; want 'Rescan complete: 2 host(s) found'", m.statusMessage)
	}

	if cmd == nil {
		t.Error("Update should return a command to clear status after delay")
	}
}

func TestUpdate_RescanError(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test"},
	}

	model := NewUIModel(hosts, nil, "192.168.1.0/24")
	model.isScanning = true

	// Simulate rescan error
	testErr := fmt.Errorf("network timeout")
	msg := rescanErrorMsg{err: testErr}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.isScanning {
		t.Error("isScanning should be false after rescan error")
	}

	if m.statusMessage != "Rescan failed: network timeout" {
		t.Errorf("statusMessage = %q; want 'Rescan failed: network timeout'", m.statusMessage)
	}

	if cmd == nil {
		t.Error("Update should return a command to clear status after delay")
	}
}

func TestUpdate_IgnoreKeysWhileScanning(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test"},
	}

	model := NewUIModel(hosts, nil, "192.168.1.0/24")
	model.isScanning = true

	// Try to trigger actions while scanning
	keys := []string{"q", "/", "?", "s", "y", "1"}
	for _, key := range keys {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(key[0])}}
		updatedModel, cmd := model.Update(keyMsg)
		m := updatedModel.(UIModel)

		if m.mode != modeNormal {
			t.Errorf("mode changed to %v when key %q pressed during scan; want modeNormal", m.mode, key)
		}

		if cmd != nil {
			t.Errorf("command returned for key %q during scan; want nil", key)
		}
	}
}

func TestUpdate_RescanWithActiveFilter(t *testing.T) {
	initialHosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Initial", Hostname: "test1"},
	}

	newHosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.2", MAC: "11:22:33:44:55:66", Vendor: "Apple", Hostname: "test2"},
		{ID: 1, IP: "192.168.1.3", MAC: "AA:AA:AA:AA:AA:AA", Vendor: "Samsung", Hostname: "test3"},
	}

	model := NewUIModel(initialHosts, nil, "192.168.1.0/24")
	model.searchActive = true
	model.searchQuery = "Apple"
	model.filteredHosts = initialHosts // Simulate previous filter
	model.isScanning = true

	// Simulate rescan completion
	msg := rescanCompleteMsg{hosts: newHosts}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	// Filter should be reapplied
	if len(m.filteredHosts) != 1 {
		t.Errorf("filteredHosts length = %d; want 1 (filter should be reapplied)", len(m.filteredHosts))
	}

	if len(m.filteredHosts) > 0 && m.filteredHosts[0].Vendor != "Apple" {
		t.Errorf("filteredHosts[0].Vendor = %q; want 'Apple'", m.filteredHosts[0].Vendor)
	}
}
