# nls Agent Instructions

## Workflow Rules

- **Always create a git branch before starting any fix or feature.** Branch off `main` using the pattern `fix/<short-description>` for bugs or `feat/<short-description>` for new features. Do this before touching any files.

## Project Overview
nls is a terminal-based network scanner that lists hosts in a network using nmap's ping scan. It combines a Go-based network scanning backend with a Bubbletea TUI (terminal user interface) for interactive display.

## Build and Run Commands

See [README.md](README.md) for build, run, and usage commands. See `.github/workflows/release.yml` for the cross-platform build matrix and `.github/workflows/test.yml` for the CI test invocation.

## Architecture

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

For per-package internals (project tree, dependencies, styling constants, error-handling details), see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Testing Conventions

### Test Structure
- Use table-driven tests for comprehensive coverage
- Test files named `*_test.go` in same package
- Tests may be grouped by behavior or concern; they do not need a 1:1 filename match with production files
- Helper functions marked with `t.Helper()`
- Subtests with `t.Run()` for better organization

### Test Coverage Goals
- Critical business logic (scanner, UI helpers, config): 80%+
- Focus on behavior and critical paths, including unexported helpers when they contain core logic
- Test behavior, not implementation details
- Use `progress.NoOp{}` for testing scanner without UI feedback
- Mock Scanner interface for testing app layer

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
- Keep files focused and split them when doing so clearly improves readability or responsibility boundaries
- Split large files by responsibility (model/view/update) when the boundary is natural
- One primary concept per file

### Testing Philosophy
- Unit tests for business logic (extractHostInfo, buildColumns, config validation)
- Integration tests for full workflows (future: app_test.go with mock scanner)
- Avoid snapshot-heavy UI rendering tests; keep UI view tests lightweight and behavior-oriented
- UI tests should focus on keyboard handling, filtering, sorting, and update-state behavior
- Table-driven tests for comprehensive coverage

## Code Style

- Follow [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
  and the formatting/style section of
  [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style).
- State your assumptions.
- Interface contracts: when ownership or lifetime semantics (e.g. buffer reuse) are important,
  document it at the interface definition, not just in the implementation.
- All exposed objects must have a doc comment.
- All comments must start with a capital letter and end with a full stop.
- Use `//nolint:linter1[,linter2,...]` sparingly; prefer fixing the code.

## Notes
- Requires nmap installed on system
- Must run with root/sudo for nmap ping scan to work
- Tests can run without root/sudo (unit tests only)
- Releases automated via GitHub Actions on version tags
- Invalid CIDR: prints error and exits with code 1
- Terminal size fallback: defaults to 100x20, checks COLUMNS/LINES env vars
