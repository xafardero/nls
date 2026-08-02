# Architecture Reference

Detailed reference for `nls`'s internal structure. See [AGENTS.md](../AGENTS.md) for the high-level data flow and design patterns; this file covers per-package internals.

## Project Structure
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
│       ├── filter_sort_test.go - Filter and sort behavior tests
│       ├── helpers_test.go  - UI helper tests
│       ├── keyboard_test.go - Keyboard interaction tests
│       └── update_test.go   - Update loop and state transition tests
├── go.mod
└── README.md
```

## Key Dependencies
- `github.com/Ullaakut/nmap/v3` - Go wrapper for nmap
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - Pre-built TUI components (table, textinput)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/schollz/progressbar/v3` - Progress spinner
- `golang.org/x/term` - Terminal size detection

## App Package (`internal/app`)
- **Config**: Centralized configuration with CIDR, Timeout, ShowProgress
- **App**: Orchestrates scan workflow (validate → scan → UI)
- **Validation**: CIDR format and timeout validation before scan
- **Context Management**: Timeout applied via `context.WithTimeout`

## Progress Package (`internal/progress`)
- **Reporter Interface**: `Start()`, `Update()`, `Finish()` methods
- **Spinner**: ProgressBar-based implementation
- **NoOp**: Silent implementation for testing/non-interactive use
- **Benefit**: Scanner decoupled from progress display library

## Scanner Package (`internal/scanner`)
- **Scanner Interface**: `Scan(ctx, target) ([]HostInfo, error)` for mockability
- **NmapScanner**: Implementation using nmap library
  - Accepts `progress.Reporter` via constructor
  - Uses buffered channels to prevent goroutine leaks
  - Context-aware for cancellation support
- **extractHostInfo()**: Extracts IP (first), MAC+Vendor (second), Hostname (first)
- **HostInfo**: Struct with ID, IP, MAC, Vendor, Hostname fields
- **IDs**: Assigned sequentially starting from 0
- **Errors**: Wrapped with context using `fmt.Errorf` and `%w`

## UI Package (`internal/ui`)
- **model.go**: UIModel struct, constants, NewUIModel() constructor
- **view.go**: Rendering logic (View(), renderHelpView(), renderSearchView(), renderSSHPromptView(), renderNormalView())
- **update.go**: Event handling (Init(), Update(), keyboard handlers, rescan workflow)
- **styles.go**: Lipgloss styles (base, selected, prompt)
- **helpers.go**: Utility functions (buildColumns, buildRows, getTerminalSize, filtering, sorting)
  - ColumnWeights for flexible column sizing (20% IP, 27% MAC, 26% Vendor, 27% Hostname)
  - Terminal size fallback via COLUMNS/LINES env vars
- **Table Interaction**:
  - `q`/`ctrl+c`: quit
  - `esc`: toggle table focus
  - `?`: show or close help screen
  - `/`: open host search/filter
  - `c`: copy selected host IP to clipboard
  - `r`: rescan the current CIDR
  - `s`: initiate SSH connection
  - `enter`: connect (when in SSH prompt)
  - `1`-`4`: sort by IP, MAC, Vendor, or Hostname
  - `↑`/`↓` or `j`/`k`: navigate rows

### Styling Conventions
- Uses lipgloss for terminal styling
- Border color: `lipgloss.Color("240")` (dark gray)
- Selected row: yellow text (`229`) on blue background (`57`), bold + underlined
- Table height: defaults to MinTableHeight (7), adjusts to terminal
- SSH prompt: rounded border with 50-character width
- Key constants: TablePaddingWidth (8), SSHUsernameMaxLen (32), HelpBoxWidth (70), SearchInputWidth (50)

## Error Handling
- Config validation errors: returned before scan starts
- Scanner errors: context-wrapped, propagated to app layer
- UI errors: returned from tea.Program.Run()
- Main function uses `run()` pattern to allow deferred cleanup
- All errors use `fmt.Errorf` with `%w` for error wrapping
