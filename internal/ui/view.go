package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI based on current state.
// Shows different views based on the current mode.
func (m UIModel) View() string {
	switch m.mode {
	case modeHelp:
		return m.renderHelpView()
	case modeSearch:
		return m.renderSearchView()
	case modeSSHPrompt:
		return m.renderSSHPromptView()
	default: // modeNormal
		return m.renderNormalView()
	}
}

// renderHelpView renders the help screen.
func (m UIModel) renderHelpView() string {
	helpBox := helpStyle.Render(helpText)

	overlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		helpBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
	return overlay
}

// renderSearchView renders the search input overlay.
func (m UIModel) renderSearchView() string {
	prompt := fmt.Sprintf("Search/Filter Hosts\n\n%s\n\n[enter: apply filter] [esc: cancel]",
		m.searchInput.View(),
	)
	searchBox := promptStyle.Render(prompt)

	overlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		searchBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
	return overlay
}

// renderSSHPromptView renders the SSH prompt overlay.
func (m UIModel) renderSSHPromptView() string {
	prompt := fmt.Sprintf("SSH to %s\n\n%s\n\n[enter: connect] [esc: cancel]",
		m.selectedIP,
		m.usernameInput.View(),
	)
	promptBox := promptStyle.Render(prompt)

	overlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		promptBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
	return overlay
}

// renderNormalView renders the standard table view with footer.
func (m UIModel) renderNormalView() string {
	baseView := baseStyle.Render(m.table.View())

	// Build footer with all shortcuts
	footer := "[?: help] [/: search] [1-4: sort] [r: rescan] [y/m/h/a: copy] [s: ssh] [q: quit]"

	// Show scanning indicator if in progress
	if m.isScanning {
		footer = "⏳ Scanning network... " + footer
	}

	// Show active filter indicator
	if m.searchActive {
		footer = fmt.Sprintf("[Filter: %s] ", m.searchQuery) + footer
	}

	// Show status message if active
	if m.statusMessage != "" && time.Now().Before(m.statusExpiry) {
		footer = m.statusMessage + "  " + footer
	}

	var b strings.Builder
	b.WriteString(baseView)
	b.WriteString("\n")
	b.WriteString(footer)
	return b.String()
}
