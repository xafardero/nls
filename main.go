package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/schollz/progressbar/v3"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "IP", Width: 10},
		{Title: "MAC", Width: 10},
		{Title: "Hostname", Width: 10},
		{Title: "OS", Width: 10},
	}
	rows := []table.Row{}
	sx := nmap_scaner()

	// Count the number of each OS for all hosts.
	for _, host := range sx.Hosts {
		i := 0
		ip := "none"
		mac := "none"
		hostname := "none"
		os := "none"

		if len(host.Addresses) > 0 {
			ip = host.Addresses[0].Addr
		}
		if len(host.Addresses) > 1 {
			mac = host.Addresses[1].Addr
		}
		if len(host.Hostnames) > 0 {
			hostname = host.Hostnames[0].Name
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
		})
		i++
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
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

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func nmap_scaner() *nmap.Run {
	bar := progressbar.Default(100, "Scanning network...")
	// Simulate progress while nmap runs
	ch := make(chan *nmap.Run)
	go func() {
		scanner, err := nmap.NewScanner(
			context.Background(),
			nmap.WithTargets("192.168.1.0/24"),
			nmap.WithFastMode(),
			nmap.WithOSDetection(), // Needs to run with sudo
		)
		if err != nil {
			log.Fatalf("unable to create nmap scanner: %v", err)
		}
		result, warnings, err := scanner.Run()
		if len(*warnings) > 0 {
			log.Printf("run finished with warnings: %s\n", *warnings)
		}
		if err != nil {
			log.Fatalf("nmap scan failed: %v", err)
		}
		ch <- result
	}()
	// Show progress bar until scan is done
	for i := 0; i < 100; i++ {
		select {
		case result := <-ch:
			bar.Finish()
			return result
		default:
			bar.Add(1)
		}
		time.Sleep(50 * time.Millisecond)
	}
	result := <-ch
	bar.Finish()
	return result
}
