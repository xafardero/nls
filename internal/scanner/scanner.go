package scanner

import (
	"context"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/schollz/progressbar/v3"
)

type Scanner interface {
	Scan(target string) (*nmap.Run, error)
}

type NmapScanner struct{}

func (s *NmapScanner) Scan(target string) (*nmap.Run, error) {
	bar := progressbar.NewOptions(-1, progressbar.OptionSetDescription("Scanning network..."), progressbar.OptionSpinnerType(14))
	ch := make(chan *nmap.Run)
	chErr := make(chan error)

	go func() {
		scanner, err := nmap.NewScanner(
			context.Background(),
			nmap.WithTargets(target),
			nmap.WithFastMode(),
			nmap.WithOSDetection(),
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

	// Animate spinner until scan is done
	for {
		select {
		case result := <-ch:
			bar.Finish()
			return result, nil
		case err := <-chErr:
			bar.Finish()
			return nil, err
		default:
			bar.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
