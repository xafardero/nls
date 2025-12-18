package main

import (
	"fmt"
	"os"

	"github.com/Ullaakut/nmap/v3"
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
	hosts := extractHostInfo(scanResult)
	model := ui.NewUIModel(hosts)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func extractHostInfo(scanResult *nmap.Run) []scanner.HostInfo {
	hosts := []scanner.HostInfo{}
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
		hosts = append(hosts, scanner.HostInfo{
			ID:       i,
			IP:       ip,
			MAC:      mac,
			Vendor:   vendor,
			Hostname: hostname,
		})
	}
	return hosts
}