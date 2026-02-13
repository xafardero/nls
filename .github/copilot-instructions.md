# Copilot Instructions for nls

## Project Overview
nls is a terminal-based network scanner that lists hosts in a network using nmap's ping scan. It combines a Go-based network scanning backend with a Bubbletea TUI (terminal user interface) for interactive display.

## Build and Run Commands

### Build
```bash
go build -o nls ./cmd/nls
```

### Run
Requires root privileges for nmap ping scan:
```bash
sudo ./nls [CIDR]
# Example: sudo ./nls 10.0.0.0/24
# Default: 192.168.1.0/24 if no CIDR provided
```

### Cross-platform Builds
The project supports Linux (amd64/arm64) and macOS (arm64):
```bash
GOOS=linux GOARCH=amd64 go build -o nls-linux-amd64 ./cmd/nls
GOOS=linux GOARCH=arm64 go build -o nls-linux-arm64 ./cmd/nls
GOOS=darwin GOARCH=arm64 go build -o nls-macos-arm64 ./cmd/nls
```

### Testing
The project follows TDD methodology with table-driven tests:
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...
```

## Architecture

### Project Structure
- `cmd/nls/main.go` - Entry point: parses CIDR, runs scanner, launches UI
- `internal/scanner/` - Network scanning logic using nmap
  - `scanner.go` - Scan function with context support and progress display
  - `types.go` - HostInfo struct definition
  - `scanner_test.go` - Table-driven tests for host extraction
- `internal/ui/` - Interactive terminal table using Bubbletea/Bubbles
  - `ui.go` - TUI model, view, and update logic
  - `ui_test.go` - Tests for UI helper functions

### Data Flow
1. Main validates CIDR argument
2. `scanner.Scan()` runs nmap ping scan with progress spinner
3. Scanner extracts host info (IP, MAC, Vendor, Hostname) into `HostInfo` structs
4. UI displays results in interactive table using Charmbracelet Bubbles

### Key Dependencies
- `github.com/Ullaakut/nmap/v3` - Go wrapper for nmap
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - Pre-built TUI components (table)
- `github.com/charmaccepts context for cancellation support
- Uses buffered channels to prevent goroutine leaks
- `extractHostInfo()` is the main testable function
- Always extracts: IP (first address), MAC+Vendor (second address), Hostname (first hostname)
- Returns `[]HostInfo` with ID assigned sequentially
- Errors are wrapped with context using `fmt.Errorf` and `%w`
- Slices are preallocated when size is known
### Scanner Package
- `scanner.Scan()` is synchronous with spinner feedback
- Uses goroutines internally with channels for async nmap execution
- Always extracts: IP (first address), MAC+Vendor (second address), Hostname (first hostname)
- Returns `[]HostInfo` with ID assigned sequentially

##Styles accessed via `getBaseStyle()` and `getPromptStyle()` functions
- Table is focused by default
- Keyboard shortcuts:
  - `q`/`ctrl+c`: quit
  - `esc`: toggle table focus
  - `s`: initiate SSH connection
  - `enter`: connect (when in SSH prompt)
- Helper functions (`buildColumns`, `buildRows`) are well-tested
  - `q`/`ctrl+c`: quit
  - `esc`: toggle table focus
  - `enter`: select row (prints to console on exit)

### Styling Conventions
- Uses lipgloss for terminal styling
- Border color: `lipgloss.Color("240")` (dark gray)
- Selected row: yellow text (`229`) on blue background (`57`), bold + underlined
- Table height: wrapped error returned from `run()` function
- Scanner errors: context-wrapped, returned to main
- Terminal size fallback: defaults to 100x20, checks COLUMNS/LINES env vars
- Main function uses `run()` pattern to allow deferred cleanup

## Testing Conventions

### Test Structure
- Use table-driven tests for comprehensive coverage
- Test files named `*_test.go` in same package
- Helper functions marked with `t.Helper()`
- Subtests with `t.Run()` for better organization

### Test Coverage Goals
- Critical business logic (scanner, UI helpers): 80%+
- Focus on public exported functions
- Test behavior, not implementation details

### Running Tests
```bash
go test ./...                    # All tests
go test -v ./...                 # Verbose
go test -cover ./...             # With coverage
go test -race ./...              # Race detection
go test -run TestName ./...      # Specific test
```
- Invalid CIDR: prints error and exits with code 1
- Scanner errors: logged to console, returned to main
- Terminal size fallback: defaults to 100x20, checks COLUMNS/LINES env vars

## Notes
- Requires nmap installed on system
- Must run with root/sudo for nmap ping scan to work
- Tests can run without root/sudo (unit tests only)
- Releases automated via GitHub Actions on version tags
- All code follows golang-patterns skill conventions
- Documentation follows Go doc comment standards
