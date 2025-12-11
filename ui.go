package main

import (
	"os"
	"strconv"

	"github.com/Ullaakut/nmap/v3"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Exported UIModel type handles table UI logic (Single Responsibility)
type UIModel struct {
	table table.Model
}

// Exported NewUIModel constructor
func NewUIModel(scanResult *nmap.Run) UIModel {
	width := 100
	height := 20
	if w, ok := os.LookupEnv("COLUMNS"); ok {
		if val, err := strconv.Atoi(w); err == nil {
			width = val
		}
	}
	if h, ok := os.LookupEnv("LINES"); ok {
		if val, err := strconv.Atoi(h); err == nil {
			height = val - 5 // leave some space for borders
		}
	}
	columns := []table.Column{
		{Title: "Id", Width: width / 12},
		{Title: "IP", Width: width / 8},
		{Title: "MAC", Width: width / 8},
		{Title: "Hostname", Width: width / 6},
		{Title: "OS", Width: width / 4},
		{Title: "Vendor", Width: width / 4},
	}
	rows := []table.Row{}
	for i, host := range scanResult.Hosts {
		ip := "none"
		mac := "none"
		hostname := "none"
		os := "none"
		vendor := "none"
		if len(host.Addresses) > 0 {
			ip = host.Addresses[0].Addr
		}
		if len(host.Addresses) > 1 {
			mac = host.Addresses[1].Addr
		}
		if len(host.Hostnames) > 0 {
			hostname = host.Hostnames[0].Name
		}
		if len(host.Addresses) > 1 {
			vendor = host.Addresses[1].Vendor
		}
		if len(host.OS.Matches) > 0 {
			os = host.OS.Matches[0].Name
		}
		rows = append(rows, table.Row{
			strconv.Itoa(i),
			ip,
			mac,
			hostname,
			os,
			vendor,
		})
	}
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
		Bold(false)
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
	return baseStyle.Render(m.table.View()) + "\n"
}
