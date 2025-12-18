# nls

A fast, terminal-based network scanner that lists hosts in a network.

## Download
You can download the latest release for Linux (amd64/arm64) or macOS (amd64/arm64) from the [Releases page](https://github.com/xafardero/nls/releases).

Example (Linux amd64):
```sh
curl -L https://github.com/xafardero/nls/releases/download/v0.1.1/nls-linux-amd64 -o nls
chmod +x nls
sudo ./nls [CIDR]
```

Example (macOS arm64):
```sh
curl -L https://github.com/xafardero/nls/releases/download/v0.1.1/nls-macos-arm64 -o nls
chmod +x nls
sudo ./nls [CIDR]
```

## Build from source
```sh
git clone https://github.com/xafardero/nls.git
cd nls
go build -o nls ./cmd/nls
```

## Usage
Run as root (required for nmap ping scan):

```sh
sudo ./nls [CIDR]
```
- `[CIDR]` is optional. If omitted, defaults to `192.168.1.0/24`.
- Example: `sudo ./nls 10.0.0.0/24`

## Features
- Scans a given CIDR subnet for live hosts using nmap's ping scan
- Responsive, interactive terminal UI with keyboard navigation
- Displays IP, MAC, Vendor, and Hostname for each discovered host
- Customizable scan range via command-line argument

## Keyboard Shortcuts
- `q` or `ctrl+c`: Quit
- `esc`: Focus/blur table
- `enter`: Select row

## Description
**nls** is a network scanner that quickly discovers live hosts in a subnet using nmap's ping scan. Results are shown in a modern, interactive terminal table UI, making it easy to browse and analyze your local network.

---
MIT License
