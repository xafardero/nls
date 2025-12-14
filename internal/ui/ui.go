package ui

import (
	"os"
	"strconv"

	"github.com/Ullaakut/nmap/v3"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
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

func NewUIModel(scanResult *nmap.Run) UIModel {
	width, height := getTerminalSize()
	if height < 7 {
		height = 7
	}

	columns := []table.Column{
		{Title: "Id", Width: width / 25},
		{Title: "IP", Width: width / 10},
		{Title: "MAC", Width: width / 5},
		{Title: "Vendor", Width: width / 5},
		{Title: "Hostname", Width: width / 10},
	}
	rows := []table.Row{}
	for i, host := range scanResult.Hosts {
		ip := "none"
		mac := "none"
		vendor := "none"
		hostname := "none"
		if len(host.Addresses) > 0 {
			ip = host.Addresses[0].Addr
		}
		if len(host.Addresses) > 1 {
			mac = host.Addresses[1].Addr
			vendor = host.Addresses[1].Vendor
		}
		if len(host.Hostnames) > 0 {
			hostname = host.Hostnames[0].Name
		}
		rows = append(rows, table.Row{
			strconv.Itoa(i),
			ip,
			mac,
			vendor,
			hostname,
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
