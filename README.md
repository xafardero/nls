# nls

A fast, terminal-based network scanner that lists hosts in a network.


## Build
```sh
git clone <your-repo-url>
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
