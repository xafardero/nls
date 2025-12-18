package ui

import (
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"nls/internal/scanner"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type UIModel struct {
	table table.Model
}

func getTerminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 100
		height = 20
		if envW, ok := os.LookupEnv("COLUMNS"); ok {
			if val, err := strconv.Atoi(envW); err == nil {
				width = val
			}
		}
		if envH, ok := os.LookupEnv("LINES"); ok {
			if val, err := strconv.Atoi(envH); err == nil {
				height = val - 5
			}
		}
	} else {
		width = w
		height = h - 5
	}
	return
}

func buildColumns(width int) []table.Column {
	// Calculate available width after fixed columns and padding
	idWidth := 5
	padding := 8 // for borders and spacing
	remaining := width - idWidth - padding
	// Assign proportional widths
	ipWidth := remaining / 5
	macWidth := remaining / 5
	vendorWidth := remaining / 5
	hostnameWidth := remaining / 5
	return []table.Column{
		{Title: "Id", Width: idWidth},
		{Title: "IP", Width: ipWidth},
		{Title: "MAC", Width: macWidth},
		{Title: "Vendor", Width: vendorWidth},
		{Title: "Hostname", Width: hostnameWidth},
	}
}



func buildRows(hosts []scanner.HostInfo) []table.Row {
	rows := []table.Row{}
	for _, h := range hosts {
		rows = append(rows, table.Row{
			strconv.Itoa(h.ID),
			h.IP,
			h.MAC,
			h.Vendor,
			h.Hostname,
		})
	}
	if len(rows) == 0 {
		rows = append(rows, table.Row{"-", "No hosts found", "-", "-", "-"})
	}
	return rows
}

func NewUIModel(hosts []scanner.HostInfo) UIModel {
	width, height := getTerminalSize()
	if height < 7 {
		height = 7
	}

	columns := buildColumns(width)
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
	return UIModel{t}
}

func (m UIModel) Init() tea.Cmd { return nil }

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m UIModel) View() string {
	footer := "[q/ctrl+c: quit] [esc: focus/blur] [enter: select]"
	return baseStyle.Render(m.table.View()) + "\n" + footer
}
