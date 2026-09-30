#!/bin/sh
set -eu

REPO="xafardero/nls"
INSTALL_DIR="${NLS_INSTALL_DIR:-/usr/local/bin}"

detect_os() {
	case "$(uname -s)" in
	Linux) echo "linux" ;;
	Darwin) echo "macos" ;;
	*)
		echo "error: unsupported OS: $(uname -s)" >&2
		exit 1
		;;
	esac
}

detect_arch() {
	case "$(uname -m)" in
	x86_64 | amd64) echo "amd64" ;;
	aarch64 | arm64) echo "arm64" ;;
	*)
		echo "error: unsupported architecture: $(uname -m)" >&2
		exit 1
		;;
	esac
}

latest_version() {
	curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
		grep '"tag_name":' |
		sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

install_binary() {
	os="$1"
	arch="$2"
	version="$3"

	if [ "$os" = "macos" ] && [ "$arch" != "arm64" ]; then
		echo "error: no macOS ${arch} build is published; only macos-arm64 is available" >&2
		exit 1
	fi

	asset="nls-${os}-${arch}"
	url="https://github.com/${REPO}/releases/download/${version}/${asset}"

	tmp="$(mktemp)"
	trap 'rm -f "$tmp"' EXIT

	echo "Downloading ${asset} ${version}..."
	curl -fsSL "$url" -o "$tmp"
	chmod +x "$tmp"

	if [ -w "$INSTALL_DIR" ]; then
		mv "$tmp" "${INSTALL_DIR}/nls"
	else
		echo "Installing to ${INSTALL_DIR} requires sudo..."
		sudo mv "$tmp" "${INSTALL_DIR}/nls"
	fi
	trap - EXIT
}

nmap_install_hint() {
	if command -v brew >/dev/null 2>&1; then
		echo "brew install nmap"
	elif command -v apt-get >/dev/null 2>&1; then
		echo "sudo apt-get install nmap"
	elif command -v dnf >/dev/null 2>&1; then
		echo "sudo dnf install nmap"
	elif command -v pacman >/dev/null 2>&1; then
		echo "sudo pacman -S nmap"
	elif command -v apk >/dev/null 2>&1; then
		echo "sudo apk add nmap"
	else
		echo "see https://nmap.org/download"
	fi
}

warn_if_nmap_missing() {
	if command -v nmap >/dev/null 2>&1; then
		return
	fi

	echo "warning: nmap is not installed; nls needs it to scan the network." >&2
	echo "To install nmap: $(nmap_install_hint)" >&2
}

main() {
	os="$(detect_os)"
	arch="$(detect_arch)"
	version="$(latest_version)"

	if [ -z "$version" ]; then
		echo "error: could not determine the latest nls version" >&2
		exit 1
	fi

	install_binary "$os" "$arch" "$version"

	echo "nls ${version} installed to ${INSTALL_DIR}/nls"
	warn_if_nmap_missing
	echo "Run 'sudo nls <CIDR>' to get started."
}

main
