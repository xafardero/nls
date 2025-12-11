package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/scanner"
	"nls/internal/ui"
)

func main() {
	scannerService := &scanner.NmapScanner{}
	scanResult, err := scannerService.Scan("192.168.1.0/24")
	if err != nil {
		fmt.Println("Error running scanner:", err)
		os.Exit(1)
	}

	model := ui.NewUIModel(scanResult)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
