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

## Architecture

### Project Structure
- `cmd/nls/main.go` - Entry point: parses CIDR, runs scanner, launches UI
- `internal/scanner/` - Network scanning logic using nmap
- `internal/ui/` - Interactive terminal table using Bubbletea/Bubbles

### Data Flow
1. Main validates CIDR argument
2. `scanner.Scan()` runs nmap ping scan with progress spinner
3. Scanner extracts host info (IP, MAC, Vendor, Hostname) into `HostInfo` structs
4. UI displays results in interactive table using Charmbracelet Bubbles

### Key Dependencies
- `github.com/Ullaakut/nmap/v3` - Go wrapper for nmap
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - Pre-built TUI components (table)
- `github.com/charmbracelet/lipgloss` - Terminal styling

## Conventions

### Scanner Package
- `scanner.Scan()` is synchronous with spinner feedback
- Uses goroutines internally with channels for async nmap execution
- Always extracts: IP (first address), MAC+Vendor (second address), Hostname (first hostname)
- Returns `[]HostInfo` with ID assigned sequentially

### UI Package
- Follows Bubbletea's Model-View-Update pattern
- `NewUIModel()` calculates responsive column widths based on terminal size
- Table is focused by default
- Keyboard shortcuts:
  - `q`/`ctrl+c`: quit
  - `esc`: toggle table focus
  - `enter`: select row (prints to console on exit)

### Styling Conventions
- Uses lipgloss for terminal styling
- Border color: `lipgloss.Color("240")` (dark gray)
- Selected row: yellow text (`229`) on blue background (`57`), bold + underlined
- Table height: terminal height minus 5 lines for borders/footer

### Error Handling
- Invalid CIDR: prints error and exits with code 1
- Scanner errors: logged to console, returned to main
- Terminal size fallback: defaults to 100x20, checks COLUMNS/LINES env vars

## Notes
- Requires nmap installed on system
- Must run with root/sudo for nmap ping scan to work
- No tests currently in codebase
- Releases automated via GitHub Actions on version tags
