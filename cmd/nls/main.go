package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"nls/internal/app"
	"nls/internal/progress"
	"nls/internal/scanner"
)

var version = "dev"

func parseArgs(arguments []string) (showVersion bool, cidr string) {
	fs := flag.NewFlagSet("nls", flag.ContinueOnError)
	versionFlag := fs.Bool("version", false, "print version and exit")
	vFlag := fs.Bool("v", false, "print version and exit")
	_ = fs.Parse(arguments)
	if fs.NArg() > 0 {
		cidr = fs.Arg(0)
	}
	showVersion = *versionFlag || *vFlag
	return
}

func run() error {
	showVersion, cidr := parseArgs(os.Args[1:])

	if showVersion {
		fmt.Printf("nls %s\n", version)
		return nil
	}

	config := app.DefaultConfig()
	config.CIDR = cidr

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
