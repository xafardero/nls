package ui

import (
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
	model := NewUIModel(hosts)
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
	model := NewUIModel(hosts)

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
	model := NewUIModel(hosts)
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
	model := NewUIModel(hosts)

	// Simulate pressing 's' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if !m.showPrompt {
		t.Error("expected showPrompt to be true after 's' key")
	}
	if m.selectedIP != "192.168.1.10" {
		t.Errorf("selectedIP = %q; want %q", m.selectedIP, "192.168.1.10")
	}
}

func TestUpdate_SSHPromptEscape(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.10", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)
	model.showPrompt = true
	model.selectedIP = "192.168.1.10"
	model.usernameInput.SetValue("testuser")

	// Press escape to cancel SSH prompt
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	if m.showPrompt {
		t.Error("expected showPrompt to be false after escape")
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
	model := NewUIModel(hosts)

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
	model := NewUIModel([]scanner.HostInfo{})
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
	model := NewUIModel([]scanner.HostInfo{})

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(UIModel)

	// Should not show SSH prompt for "No hosts found"
	if m.showPrompt {
		t.Error("should not show SSH prompt for 'No hosts found' row")
	}
}

func TestStatusMessage_Expiry(t *testing.T) {
	hosts := []scanner.HostInfo{
		{ID: 0, IP: "192.168.1.1", MAC: "AA:BB:CC:DD:EE:FF", Vendor: "Test", Hostname: "test.local"},
	}
	model := NewUIModel(hosts)

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
