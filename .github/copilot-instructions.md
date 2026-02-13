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
```
nls/
├── cmd/nls/
│   └── main.go              - Entry point (minimal, delegates to app)
├── internal/
│   ├── app/                 - Application orchestration layer
│   │   ├── app.go           - App coordination & workflow
│   │   ├── config.go        - Configuration management
│   │   └── config_test.go   - Config validation tests
│   ├── progress/            - Progress reporting abstraction
│   │   ├── reporter.go      - Reporter interface + NoOp implementation
│   │   └── spinner.go       - Spinner implementation
│   ├── scanner/             - Network scanning using nmap
│   │   ├── scanner.go       - Scanner interface
│   │   ├── nmap.go          - NmapScanner implementation
│   │   ├── types.go         - HostInfo struct definition
│   │   └── scanner_test.go  - Table-driven tests
│   └── ui/                  - Interactive TUI (Bubbletea/Bubbles)
│       ├── model.go         - UIModel & initialization
│       ├── view.go          - Rendering logic
│       ├── update.go        - Event handling (Init/Update)
│       ├── styles.go        - Lipgloss styling
│       ├── helpers.go       - Helper functions (columns, rows, terminal)
│       └── helpers_test.go  - UI tests
├── go.mod
└── README.md
```

### Data Flow
1. `main.go` creates Config, injects dependencies (Scanner, ProgressReporter)
2. `app.Run()` validates config, executes scan workflow, launches UI
3. `NmapScanner.Scan()` runs nmap with progress feedback via Reporter interface
4. Scanner extracts host info (IP, MAC, Vendor, Hostname) into `HostInfo` structs
5. UI displays results in interactive table using Charmbracelet Bubbles

### Design Patterns
- **Dependency Injection**: Scanner and ProgressReporter injected into App
- **Interface-based Design**: Scanner and Reporter are interfaces (mockable)
- **Configuration Management**: Centralized Config struct with validation
- **MVC-like Separation**: UI split into Model/View/Update/Styles/Helpers
- **Factory Pattern**: `NewNmapScanner()`, `NewUIModel()` constructors

### Key Dependencies
- `github.com/Ullaakut/nmap/v3` - Go wrapper for nmap
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - Pre-built TUI components (table, textinput)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/schollz/progressbar/v3` - Progress spinner
- `golang.org/x/term` - Terminal size detection

### App Package (`internal/app`)
- **Config**: Centralized configuration with CIDR, Timeout, ShowProgress
- **App**: Orchestrates scan workflow (validate → scan → UI)
- **Validation**: CIDR format and timeout validation before scan
- **Context Management**: Timeout applied via `context.WithTimeout`

### Progress Package (`internal/progress`)
- **Reporter Interface**: `Start()`, `Update()`, `Finish()` methods
- **Spinner**: ProgressBar-based implementation
- **NoOp**: Silent implementation for testing/non-interactive use
- **Benefit**: Scanner decoupled from progress display library

### Scanner Package (`internal/scanner`)
- **Scanner Interface**: `Scan(ctx, target) ([]HostInfo, error)` for mockability
- **NmapScanner**: Implementation using nmap library
  - Accepts `progress.Reporter` via constructor
  - Uses buffered channels to prevent goroutine leaks
  - Context-aware for cancellation support
- **extractHostInfo()**: Extracts IP (first), MAC+Vendor (second), Hostname (first)
- **HostInfo**: Struct with ID, IP, MAC, Vendor, Hostname fields
- **IDs**: Assigned sequentially starting from 0
- **Errors**: Wrapped with context using `fmt.Errorf` and `%w`

### UI Package (`internal/ui`)
- **model.go**: UIModel struct, constants, NewUIModel() constructor
- **view.go**: Rendering logic (View(), renderPromptView(), renderNormalView())
- **update.go**: Event handling (Init(), Update(), keyboard handlers)
- **styles.go**: Lipgloss styles (base, selected, prompt)
- **helpers.go**: Utility functions (buildColumns, buildRows, getTerminalSize)
  - ColumnWeights for flexible column sizing (20% IP, 27% MAC, 26% Vendor, 27% Hostname)
  - Terminal size fallback via COLUMNS/LINES env vars
- **Table Interaction**:
  - `q`/`ctrl+c`: quit
  - `esc`: toggle table focus
  - `s`: initiate SSH connection
  - `enter`: connect (when in SSH prompt)
  - `↑`/`↓` or `j`/`k`: navigate rows

### Styling Conventions
- Uses lipgloss for terminal styling
- Border color: `lipgloss.Color("240")` (dark gray)
- Selected row: yellow text (`229`) on blue background (`57`), bold + underlined
- Table height: defaults to MinTableHeight (7), adjusts to terminal
- SSH prompt: rounded border with 50-character width
- Constants: TableIDWidth (5), TablePaddingWidth (8), SSHUsernameMaxLen (32)

### Error Handling
- Config validation errors: returned before scan starts
- Scanner errors: context-wrapped, propagated to app layer
- UI errors: returned from tea.Program.Run()
- Main function uses `run()` pattern to allow deferred cleanup
- All errors use `fmt.Errorf` with `%w` for error wrapping

## Testing Conventions

### Test Structure
- Use table-driven tests for comprehensive coverage
- Test files named `*_test.go` in same package
- Helper functions marked with `t.Helper()`
- Subtests with `t.Run()` for better organization

### Test Coverage Goals
- Critical business logic (scanner, UI helpers, config): 80%+
- Focus on public exported functions
- Test behavior, not implementation details
- Use `progress.NoOp{}` for testing scanner without UI feedback
- Mock Scanner interface for testing app layer

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

## Code Organization Principles

### Package Responsibilities
- **cmd/nls**: Entry point only, minimal logic
- **internal/app**: Application orchestration, not business logic
- **internal/scanner**: Network scanning, isolated from UI/progress
- **internal/progress**: Progress feedback, abstracted from implementation
- **internal/ui**: Terminal UI, separated by concern (MVC-like)

### Dependency Rules
- App depends on: scanner (interface), progress (interface), ui
- Scanner depends on: progress (interface), nmap library
- UI depends on: scanner (types only), bubbletea, bubbles, lipgloss
- Progress depends on: nothing (interface) / progressbar (implementation)
- No circular dependencies

### File Size Guidelines
- Keep files under 150 lines when possible
- Split large files by responsibility (model/view/update)
- One primary concept per file

### Testing Philosophy
- Unit tests for business logic (extractHostInfo, buildColumns, config validation)
- Integration tests for full workflows (future: app_test.go with mock scanner)
- No tests for pure UI rendering (Bubbletea handles this)
- Table-driven tests for comprehensive coverage

## Notes
- Requires nmap installed on system
- Must run with root/sudo for nmap ping scan to work
- Tests can run without root/sudo (unit tests only)
- Releases automated via GitHub Actions on version tags
- All code follows golang-patterns skill conventions
- Documentation follows Go doc comment standards
- Avoid inline comments except for complex logic; prefer descriptive function/variable names
- Architecture follows 2026 Go best practices (interfaces, DI, separation of concerns)
