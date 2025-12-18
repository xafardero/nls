package main

import (
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"nls/internal/scanner"
	"nls/internal/ui"
)

func main() {
	cidr := "192.168.1.0/24"
	if len(os.Args) > 1 {
		cidr = os.Args[1]
	}

	if _, _, err := net.ParseCIDR(cidr); err != nil {
		fmt.Println("Invalid CIDR format:", cidr)
		os.Exit(1)
	}

	hosts, err := scanner.Scan(cidr)
	if err != nil {
		fmt.Println("Error running scanner:", err)
		os.Exit(1)
	}
	model := ui.NewUIModel(hosts)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running ui:", err)
		os.Exit(1)
	}
}
