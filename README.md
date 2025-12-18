# nls

A fast, terminal-based network scanner that lists hosts in a network.

## Download
You can download the latest release for Linux (amd64/arm64) or macOS (amd64/arm64) from the [Releases page](https://github.com/xafardero/nls/releases).

Example (Linux amd64):
```sh
curl -L https://github.com/xafardero/nls/releases/latest/download/nls-linux-amd64 -o nls
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
