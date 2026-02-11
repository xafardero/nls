// Package ui provides a terminal user interface for displaying network scan results.
// It uses the Bubbletea framework to create an interactive table view with
// keyboard navigation and SSH connection capabilities.
package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"nls/internal/scanner"
)

const (
	TableIDWidth          = 5
	TablePaddingWidth     = 8
	MinTableHeight        = 7
	DefaultTermWidth      = 100
	DefaultTermHeight     = 20
	DefaultTermHeightPad  = 5
	SSHUsernameMaxLen     = 32
	SSHUsernameInputWidth = 40
)

type ColumnWeights struct {
	IP       float64
	MAC      float64
	Vendor   float64
	Hostname float64
}

func DefaultColumnWeights() ColumnWeights {
	return ColumnWeights{
		IP:       0.20,
		MAC:      0.27,
		Vendor:   0.26,
		Hostname: 0.27,
	}
}

type UIModel struct {
	table         table.Model
	showPrompt    bool
	usernameInput textinput.Model
	selectedIP    string
}

// getTerminalSize returns the current terminal width and height.
// Falls back to environment variables or default values
// if terminal size detection fails.
func getTerminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = DefaultTermWidth
		height = DefaultTermHeight
		if envW, ok := os.LookupEnv("COLUMNS"); ok {
			if val, err := strconv.Atoi(envW); err == nil {
				width = val
			}
		}
		if envH, ok := os.LookupEnv("LINES"); ok {
			if val, err := strconv.Atoi(envH); err == nil {
				height = val - DefaultTermHeightPad
			}
		}
	} else {
		width = w
		height = h - DefaultTermHeightPad
	}
	return
}

// buildColumns creates table column definitions based on terminal width.
// Columns are proportionally sized using the provided weights.
func buildColumns(width int, weights ColumnWeights) []table.Column {
	// Calculate available width after fixed columns and padding
	remaining := width - TableIDWidth - TablePaddingWidth

	// Assign proportional widths based on weights
	ipWidth := int(float64(remaining) * weights.IP)
	macWidth := int(float64(remaining) * weights.MAC)
	vendorWidth := int(float64(remaining) * weights.Vendor)
	hostnameWidth := int(float64(remaining) * weights.Hostname)

	return []table.Column{
		{Title: "Id", Width: TableIDWidth},
		{Title: "IP", Width: ipWidth},
		{Title: "MAC", Width: macWidth},
		{Title: "Vendor", Width: vendorWidth},
		{Title: "Hostname", Width: hostnameWidth},
	}
}

func buildRows(hosts []scanner.HostInfo) []table.Row {
	if len(hosts) == 0 {
		return []table.Row{{"-", "No hosts found", "-", "-", "-"}}
	}

	rows := make([]table.Row, 0, len(hosts))
	for _, h := range hosts {
		rows = append(rows, table.Row{
			strconv.Itoa(h.ID),
			h.IP,
			h.MAC,
			h.Vendor,
			h.Hostname,
		})
	}
	return rows
}

func getBaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
}

func getPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(50)
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
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true).
		Underline(true)
	t.SetStyles(s)

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

func (m UIModel) Init() tea.Cmd { return nil }

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showPrompt {
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
				m.usernameInput, cmd = m.usernameInput.Update(msg)
				return m, cmd
			}
		}

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
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the UI based on current state.
// Shows either the normal table view or the SSH prompt overlay.
func (m UIModel) View() string {
	baseView := getBaseStyle().Render(m.table.View())

	if m.showPrompt {
		prompt := fmt.Sprintf("SSH to %s\n\n%s\n\n[enter: connect] [esc: cancel]",
			m.selectedIP,
			m.usernameInput.View(),
		)
		promptBox := getPromptStyle().Render(prompt)

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

	footer := "[q/ctrl+c: quit] [esc: focus/blur] [s: ssh]"

	// Use strings.Builder for efficient string concatenation
	var b strings.Builder
	b.WriteString(baseView)
	b.WriteString("\n")
	b.WriteString(footer)
	return b.String()
}
