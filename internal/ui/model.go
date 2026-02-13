// Package ui provides a terminal user interface for displaying network scan results.
// It uses the Bubbletea framework to create an interactive table view with
// keyboard navigation and SSH connection capabilities.
package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"

	"nls/internal/scanner"
)

// UI layout constants
const (
	TableIDWidth          = 5
	TablePaddingWidth     = 8
	MinTableHeight        = 7
	DefaultTermWidth      = 100
	DefaultTermHeight     = 20
	DefaultTermHeightPad  = 5
	SSHUsernameMaxLen     = 32
	SSHUsernameInputWidth = 40
	SSHPromptWidth        = 50
	SSHPromptPadding      = 1
)

// UIModel represents the state of the terminal UI.
// It manages the display table, SSH prompt dialog, and user input.
// UIModel requires initialization via NewUIModel and cannot be used
// with its zero value due to dependencies on Bubbletea components.
type UIModel struct {
	table         table.Model
	showPrompt    bool
	usernameInput textinput.Model
	selectedIP    string
}

// NewUIModel creates a new UI model. UIModel requires initialization
// and cannot be used with its zero value due to dependencies on
// the Bubbletea table component.
func NewUIModel(hosts []scanner.HostInfo) UIModel {
	width, height := getTerminalSize()
	if height < MinTableHeight {
		height = MinTableHeight
	}

	weights := DefaultColumnWeights()
	columns := buildColumns(width, weights)
	rows := buildRows(hosts)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	t.SetStyles(tableStyles())

	ti := textinput.New()
	ti.Placeholder = "username"
	ti.Focus()
	ti.CharLimit = SSHUsernameMaxLen
	ti.Width = SSHUsernameInputWidth

	return UIModel{
		table:         t,
		showPrompt:    false,
		usernameInput: ti,
	}
}
