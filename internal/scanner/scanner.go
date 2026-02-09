package scanner

import (
	"context"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/schollz/progressbar/v3"
)

func Scan(target string) ([]HostInfo, error) {
	bar := progressbar.NewOptions(-1, progressbar.OptionSetDescription("Scanning network..."), progressbar.OptionSpinnerType(14))
	defer func() {
		if err := bar.Finish(); err != nil {
			log.Printf("progressbar finish error: %v", err)
		}
	}()

	ch := make(chan *nmap.Run)
	chErr := make(chan error)

	go func() {
		scanner, err := nmap.NewScanner(
			context.Background(),
			nmap.WithTargets(target),
			nmap.WithPingScan(),
		)
		if err != nil {
			chErr <- err
			return
		}
		result, warnings, err := scanner.Run()
		if len(*warnings) > 0 {
			log.Printf("run finished with warnings: %s\n", *warnings)
		}
		if err != nil {
			chErr <- err
			return
		}
		ch <- result
	}()

	for {
		select {
		case result := <-ch:
			return extractHostInfo(result), nil
		case err := <-chErr:
			return nil, err
		default:
			if err := bar.Add(1); err != nil {
				return nil, err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func extractHostInfo(scanResult *nmap.Run) []HostInfo {
	hosts := []HostInfo{}
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
		hosts = append(hosts, HostInfo{
			ID:       i,
			IP:       ip,
			MAC:      mac,
			Vendor:   vendor,
			Hostname: hostname,
		})
	}
	return hosts
}
