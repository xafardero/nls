package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	scanner := &NmapScanner{}
	scanResult, err := scanner.Scan("192.168.1.0/24")
	if err != nil {
		fmt.Println("Error running scanner:", err)
		os.Exit(1)
	}

	model := NewUIModel(scanResult)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
