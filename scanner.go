package main

import (
	"context"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/schollz/progressbar/v3"
)

// Scanner interface for dependency inversion
// Can be mocked for tests
// Open/Closed principle: extend by implementing Scanner

type Scanner interface {
	Scan(target string) (*nmap.Run, error)
}

// Exported NmapScanner struct
type NmapScanner struct{}

// Exported Scan method
func (s *NmapScanner) Scan(target string) (*nmap.Run, error) {
	bar := progressbar.Default(100, "Scanning network...")
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

	for i := 0; i < 100; i++ {
		select {
		case result := <-ch:
			bar.Finish()
			return result, nil
		case err := <-chErr:
			bar.Finish()
			return nil, err
		default:
			bar.Add(1)
		}
		time.Sleep(50 * time.Millisecond)
	}
	result := <-ch
	bar.Finish()
	return result, nil
}
