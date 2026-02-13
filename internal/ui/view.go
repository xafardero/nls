package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI based on current state.
// Shows either the normal table view or the SSH prompt overlay.
func (m UIModel) View() string {
	baseView := baseStyle.Render(m.table.View())

	if m.showPrompt {
		return m.renderPromptView(baseView)
	}

	return m.renderNormalView(baseView)
}

// renderPromptView renders the SSH prompt overlay.
func (m UIModel) renderPromptView(baseView string) string {
	prompt := fmt.Sprintf("SSH to %s\n\n%s\n\n[enter: connect] [esc: cancel]",
		m.selectedIP,
		m.usernameInput.View(),
	)
	promptBox := promptStyle.Render(prompt)

	width, height := getTerminalSize()
	overlay := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		promptBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
	return overlay
}

// renderNormalView renders the standard table view with footer.
func (m UIModel) renderNormalView(baseView string) string {
	footer := "[q/ctrl+c: quit] [esc: focus/blur] [s: ssh]"

	var b strings.Builder
	b.WriteString(baseView)
	b.WriteString("\n")
	b.WriteString(footer)
	return b.String()
}
