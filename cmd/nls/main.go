package main

import (
	"context"
	"fmt"
	"os"

	"nls/internal/app"
	"nls/internal/progress"
	"nls/internal/scanner"
)

func run() error {
	config := app.DefaultConfig()

	if len(os.Args) > 1 {
		config.CIDR = os.Args[1]
	}

	var progressReporter progress.Reporter
	if config.ShowProgress {
		progressReporter = progress.NewSpinner()
	} else {
		progressReporter = progress.NoOp{}
	}
	nmapScanner := scanner.NewNmapScanner(progressReporter)

	application := app.New(config, nmapScanner)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	return application.Run(ctx)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
