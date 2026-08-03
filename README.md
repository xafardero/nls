# nls

[![Test](https://github.com/xafardero/nls/actions/workflows/test.yml/badge.svg)](https://github.com/xafardero/nls/actions/workflows/test.yml)
[![Lint](https://github.com/xafardero/nls/actions/workflows/lint.yml/badge.svg)](https://github.com/xafardero/nls/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/xafardero/nls)](https://goreportcard.com/report/github.com/xafardero/nls)
[![Release](https://img.shields.io/github/v/release/xafardero/nls)](https://github.com/xafardero/nls/releases)
[![License: MIT](https://img.shields.io/github/license/xafardero/nls)](LICENSE)

A terminal-based network scanner that lists hosts in a network using nmap's ping scan. Results are displayed in an interactive terminal UI for easy browsing.

Unlike raw `nmap`/`arp-scan` output, `nls` gives you a live, sortable, filterable table you can act on directly — search by IP/MAC/vendor/hostname, SSH into a host, or copy its IP, all without leaving the terminal or re-running commands with different flags.

![Demo](img/demo.gif)

## Download
Download the latest release for Linux (amd64/arm64) or macOS (arm64) from the [Releases page](https://github.com/xafardero/nls/releases).

```sh
# Replace {OS}-{ARCH} with your platform (e.g., linux-amd64, macos-arm64)
curl -L https://github.com/xafardero/nls/releases/download/v0.2.0/nls-linux-amd64 -o nls
chmod +x nls
sudo mv nls /usr/local/bin/
```

Now you can run `nls` from anywhere.

## Build from source
```sh
git clone https://github.com/xafardero/nls.git
cd nls
go build -o nls ./cmd/nls
```

## Usage
Run as root (required for nmap ping scan):

```sh
sudo nls <CIDR>
```
- CIDR is required (e.g., `sudo nls 192.168.1.0/24`)
- Example: `sudo nls 10.0.0.0/24`
- Check the installed version: `nls --version` (or `nls -v`)

**Keyboard Shortcuts:**

**Navigation:**
- `↑`/`↓` or `j`/`k`: Navigate table
- `esc`: Toggle table focus

**Actions:**
- `s`: SSH to selected host
- `c`: Copy IP to clipboard
- `r`: Rescan network (refreshes host list)

**Search & Sort:**
- `/`: Search/filter hosts (matches IP, MAC, Vendor, or Hostname)
- `1`: Sort by IP address
- `2`: Sort by MAC address
- `3`: Sort by Vendor
- `4`: Sort by Hostname
- Press the same number again to toggle ascending/descending

**Help & Exit:**
- `?`: Show help screen with all shortcuts
- `q` or `ctrl+c`: Quit

## Features
- Fast network scanning using nmap's ping scan
- Displays IP, MAC address, vendor, and hostname for each host
- SSH directly to any host from the UI
- Live search/filter, column sorting, clipboard copy, and rescan — all without leaving the terminal

## License
[MIT](LICENSE)
