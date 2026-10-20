# Contributing to nls

Thank you for contributing to nls. Contributions to code, documentation, tests, and bug reports are welcome.

## Before You Start

- Search existing issues and pull requests to avoid duplicate work.
- For substantial changes, open an issue first to agree on the approach.
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development Setup

1. Install [Go](https://go.dev/dl/) using the version declared in `go.mod`.
2. Install [nmap](https://nmap.org/download) if you need to run a real scan.
3. Clone the repository and download dependencies:

   ```sh
   git clone https://github.com/xafardero/nls.git
   cd nls
   go mod download
   ```

## Make a Change

1. Create a branch from `main` named `fix/<description>` or `feat/<description>`.
2. Keep changes focused and add or update tests for changed behavior.
3. Format and verify the change:

   ```sh
   gofmt -w .
   go test -race ./...
   go build ./cmd/nls
   ```

Unit tests do not require root access or nmap. Running an actual scan does require nmap and usually root privileges:

```sh
sudo go run ./cmd/nls 192.168.1.0/24
```

## Pull Requests

- Use a clear title that describes the user-visible change.
- Explain what changed and how you tested it.
- Keep the pull request limited to one concern.
- Ensure CI passes before requesting review.

## Reporting Security Vulnerabilities

Do not open a public issue for a suspected security vulnerability. Use the repository's private security reporting channel instead.
