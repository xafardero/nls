package ui

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/progress"
	"nls/internal/scanner"
)

// clearStatusMsg is sent after a delay to clear the status message.
type clearStatusMsg struct{}

// rescanCompleteMsg is sent when a rescan finishes successfully.
type rescanCompleteMsg struct {
	hosts []scanner.HostInfo
}

// rescanErrorMsg is sent when a rescan fails.
type rescanErrorMsg struct {
	err error
}

// doRescan performs a network rescan in a goroutine and returns the result as a message.
// Creates a new scanner with NoOp progress reporter to avoid terminal output conflicts with the TUI.
// If a non-NmapScanner is passed (e.g., mock for testing), it uses that scanner directly.
func doRescan(s scanner.Scanner, cidr string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// For NmapScanner, create a silent version to avoid progress output interfering with TUI
		// For other scanners (e.g., mocks in tests), use the provided scanner
		scannerToUse := s
		if _, ok := s.(*scanner.NmapScanner); ok {
			scannerToUse = scanner.NewNmapScanner(progress.NoOp{})
		}

		hosts, err := scannerToUse.Scan(ctx, cidr)
		if err != nil {
			return rescanErrorMsg{err: err}
		}
		return rescanCompleteMsg{hosts: hosts}
	}
}

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
	case clearStatusMsg:
		m.statusMessage = ""
		return m, nil

	case rescanCompleteMsg:
		// Update hosts with new scan results
		m.isScanning = false
		m.allHosts = msg.hosts

		// Reapply current filter if active
		if m.searchActive {
			m.filteredHosts = filterHosts(m.allHosts, m.searchQuery)
		} else {
			m.filteredHosts = m.allHosts
		}

		// Rebuild the table
		m = m.rebuildTable()

		// Show success message
		m.statusMessage = fmt.Sprintf("Rescan complete: %d host(s) found", len(m.allHosts))
		m.statusExpiry = time.Now().Add(3 * time.Second)
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})

	case rescanErrorMsg:
		// Handle scan error
		m.isScanning = false
		m.statusMessage = fmt.Sprintf("Rescan failed: %v", msg.err)
		m.statusExpiry = time.Now().Add(5 * time.Second)
		return m, tea.Tick(5*time.Second, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})

	case tea.KeyMsg:
		// Ignore keyboard input while scanning
		if m.isScanning {
			return m, nil
		}

		// Route to appropriate handler based on view mode
		switch m.mode {
		case modeHelp:
			return m.handleHelpKeys(msg)
		case modeSearch:
			return m.handleSearchKeys(msg)
		case modeSSHPrompt:
			return m.handleSSHPromptKeys(msg)
		default: // modeNormal
			return m.handleNormalKeys(msg)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// handleHelpKeys handles keyboard input when help screen is shown.
func (m UIModel) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "?":
		m.mode = modeNormal
		m.table.Focus()
		return m, nil
	}
	return m, nil
}

// handleSearchKeys handles keyboard input when search mode is active.
func (m UIModel) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel search, return to normal mode
		m.mode = modeNormal
		m.searchInput.SetValue("")
		m.searchInput.Blur()
		m.table.Focus()
		return m, nil

	case "enter":
		// Apply search filter
		query := m.searchInput.Value()
		m.searchQuery = query
		m.searchActive = query != ""

		// Filter hosts
		if m.searchActive {
			m.filteredHosts = filterHosts(m.allHosts, query)
		} else {
			m.filteredHosts = m.allHosts
		}

		// Rebuild table with filtered and sorted data
		m = m.rebuildTable()

		m.mode = modeNormal
		m.searchInput.Blur()
		m.table.Focus()

		// Show status message
		m.statusMessage = fmt.Sprintf("Found %d host(s)", len(m.filteredHosts))
		m.statusExpiry = time.Now().Add(2 * time.Second)
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})

	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}
}

// handleSSHPromptKeys handles keyboard input when SSH prompt is shown.
func (m UIModel) handleSSHPromptKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		m.usernameInput.SetValue("")
		m.table.Focus()
		return m, nil

	case "enter":
		username := m.usernameInput.Value()
		if username == "" {
			return m, nil
		}
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
	case "?":
		// Show help screen
		m.mode = modeHelp
		m.table.Blur()
		return m, nil

	case "/":
		// Activate search mode
		m.mode = modeSearch
		m.searchInput.Focus()
		m.table.Blur()
		return m, nil

	case "esc":
		if m.table.Focused() {
			m.table.Blur()
		} else {
			m.table.Focus()
		}

	case "q", "ctrl+c":
		return m, tea.Quit

	case "1", "2", "3", "4":
		// Sort by column
		col := int(msg.String()[0] - '0')
		if m.sortColumn == col {
			// Toggle sort direction
			m.sortAscending = !m.sortAscending
		} else {
			// New column, default to ascending
			m.sortColumn = col
			m.sortAscending = true
		}
		m = m.rebuildTable()
		return m, nil

	case "r":
		// Trigger network rescan
		if m.isScanning {
			// Already scanning, ignore
			return m, nil
		}
		m.isScanning = true
		return m, doRescan(m.scanner, m.cidr)

	case "y":
		// Copy IP to clipboard
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 1 && selectedRow[1] != "No hosts found" {
			ip := selectedRow[1]
			if err := clipboard.WriteAll(ip); err == nil {
				m.statusMessage = "IP copied to clipboard!"
				m.statusExpiry = time.Now().Add(2 * time.Second)
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return clearStatusMsg{}
				})
			}
		}

	case "m":
		// Copy MAC to clipboard
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 2 && selectedRow[1] != "No hosts found" {
			mac := selectedRow[2]
			if err := clipboard.WriteAll(mac); err == nil {
				m.statusMessage = "MAC address copied to clipboard!"
				m.statusExpiry = time.Now().Add(2 * time.Second)
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return clearStatusMsg{}
				})
			}
		}

	case "h":
		// Copy hostname to clipboard
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 4 && selectedRow[1] != "No hosts found" {
			hostname := selectedRow[4]
			if err := clipboard.WriteAll(hostname); err == nil {
				m.statusMessage = "Hostname copied to clipboard!"
				m.statusExpiry = time.Now().Add(2 * time.Second)
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return clearStatusMsg{}
				})
			}
		}

	case "a":
		// Copy all fields to clipboard
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 4 && selectedRow[1] != "No hosts found" {
			// Format: IP\tMAC\tVendor\tHostname
			allFields := fmt.Sprintf("%s\t%s\t%s\t%s",
				selectedRow[1], selectedRow[2], selectedRow[3], selectedRow[4])
			if err := clipboard.WriteAll(allFields); err == nil {
				m.statusMessage = "All fields copied to clipboard!"
				m.statusExpiry = time.Now().Add(2 * time.Second)
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return clearStatusMsg{}
				})
			}
		}

	case "s":
		// SSH to selected host
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 1 && selectedRow[1] != "No hosts found" {
			m.selectedIP = selectedRow[1]
			m.mode = modeSSHPrompt
			m.table.Blur()
			m.usernameInput.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// rebuildTable rebuilds the table with current filter and sort settings.
func (m UIModel) rebuildTable() UIModel {
	// Apply sort to filtered hosts
	hostsToDisplay := m.filteredHosts
	if m.sortColumn > 0 {
		hostsToDisplay = sortHosts(hostsToDisplay, m.sortColumn, m.sortAscending)
	}

	// Rebuild columns with sort indicator
	width, _ := getTerminalSize()
	weights := DefaultColumnWeights()
	columns := buildColumns(width, weights, m.sortColumn, m.sortAscending)

	// Rebuild rows
	rows := buildRows(hostsToDisplay)

	// Update table
	m.table.SetColumns(columns)
	m.table.SetRows(rows)

	return m
}
