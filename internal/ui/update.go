package ui

import (
"fmt"
"os/exec"

tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI model.
// Returns nil as no initial commands are needed.
func (m UIModel) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and updates the model state.
// Implements the Bubbletea Update interface for event handling.
func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showPrompt {
			return m.handlePromptKeys(msg)
		}
		return m.handleNormalKeys(msg)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// handlePromptKeys handles keyboard input when SSH prompt is shown.
func (m UIModel) handlePromptKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.showPrompt = false
		m.usernameInput.SetValue("")
		m.table.Focus()
		return m, nil

	case "enter":
		username := m.usernameInput.Value()
		if username == "" {
			return m, nil
		}
		m.showPrompt = false
		m.usernameInput.SetValue("")

		sshCmd := exec.Command("ssh", fmt.Sprintf("%s@%s", username, m.selectedIP))
		return m, tea.ExecProcess(sshCmd, func(err error) tea.Msg {
return nil
})

	default:
		var cmd tea.Cmd
		m.usernameInput, cmd = m.usernameInput.Update(msg)
		return m, cmd
	}
}

// handleNormalKeys handles keyboard input in normal table view mode.
func (m UIModel) handleNormalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.table.Focused() {
			m.table.Blur()
		} else {
			m.table.Focus()
		}

	case "q", "ctrl+c":
		return m, tea.Quit

	case "s":
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 1 && selectedRow[1] != "No hosts found" {
			m.selectedIP = selectedRow[1]
			m.showPrompt = true
			m.table.Blur()
			m.usernameInput.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
