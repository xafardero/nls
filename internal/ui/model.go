// Package ui provides a terminal user interface for displaying network scan results.
// It uses the Bubbletea framework to create an interactive table view with
// keyboard navigation and SSH connection capabilities.
package ui

import (
	"time"

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
	HelpBoxWidth          = 70
	HelpBoxPadding        = 2
	SearchInputWidth      = 50
)

// viewMode represents the current view/screen mode
type viewMode int

const (
	modeNormal viewMode = iota
	modeHelp
	modeSearch
	modeSSHPrompt
)

// Help screen content
const helpText = `Keyboard Shortcuts:

  Navigation:
    ↑/k          Move up
    ↓/j          Move down
    esc          Toggle table focus

  Actions:
    s            SSH to selected host
    y            Copy IP to clipboard
    m            Copy MAC to clipboard
    h            Copy hostname to clipboard
    a            Copy all fields to clipboard
    r            Rescan network

  Search & Sort:
    /            Search/filter hosts
    1            Sort by IP
    2            Sort by MAC
    3            Sort by Vendor
    4            Sort by Hostname

  Other:
    ?            Show this help
    q/ctrl+c     Quit

Press esc or q to close this help screen.`

// UIModel represents the state of the terminal UI.
// It manages the display table, multiple view modes, search/filter, sorting,
// and user input. UIModel requires initialization via NewUIModel and cannot
// be used with its zero value due to dependencies on Bubbletea components.
type UIModel struct {
	// Display components
	table         table.Model
	usernameInput textinput.Model
	searchInput   textinput.Model

	// Data storage
	allHosts      []scanner.HostInfo // Original host data
	filteredHosts []scanner.HostInfo // After applying search filter

	// View state
	mode          viewMode
	statusMessage string
	statusExpiry  time.Time

	// SSH state
	selectedIP string

	// Search/Filter state
	searchActive bool
	searchQuery  string

	// Sort state
	sortColumn    int // 0=none, 1=IP, 2=MAC, 3=Vendor, 4=Hostname
	sortAscending bool
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
	columns := buildColumns(width, weights, 0, false) // No initial sort
	rows := buildRows(hosts)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	t.SetStyles(tableStyles())

	// SSH username input
	ti := textinput.New()
	ti.Placeholder = "username"
	ti.Focus()
	ti.CharLimit = SSHUsernameMaxLen
	ti.Width = SSHUsernameInputWidth

	// Search input
	si := textinput.New()
	si.Placeholder = "Search (IP, MAC, Vendor, Hostname)..."
	si.CharLimit = 50
	si.Width = SearchInputWidth

	return UIModel{
		table:         t,
		allHosts:      hosts,
		filteredHosts: hosts, // Initially, no filter applied
		usernameInput: ti,
		searchInput:   si,
		mode:          modeNormal,
		searchActive:  false,
		sortColumn:    0,
		sortAscending: true,
	}
}
