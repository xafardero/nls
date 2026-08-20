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
	echo "Run 'sudo nls <CIDR>' to get started."
}

main
