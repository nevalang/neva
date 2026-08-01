#!/usr/bin/env bash

set -euo pipefail

detect_platform() {
	local os arch
	os=$(uname -s | tr '[:upper:]' '[:lower:]')
	arch=$(uname -m)
	case "$arch" in
		x86_64) arch="amd64" ;;
		arm64|aarch64) arch="arm64" ;;
		loong64|loongarch64) arch="loong64" ;;
		*)
			echo "Unsupported architecture: $arch" >&2
			exit 1
			;;
	esac
	printf '%s-%s\n' "$os" "$arch"
}

latest_tag=$(curl -fsSL https://api.github.com/repos/nevalang/neva/releases/latest | grep '"tag_name"' | cut -d '"' -f 4)
if [[ -z "$latest_tag" ]]; then
	echo 'Could not determine the latest Neva release tag.' >&2
	exit 1
fi
echo "Latest tag is $latest_tag"

platform=$(detect_platform)
echo "Platform is $platform"

install_dir="${NEVA_INSTALL_DIR:-$HOME/.local/bin}"
binary_path="$(mktemp)"
trap 'rm -f "$binary_path"' EXIT

echo 'Downloading...'
curl -fL "https://github.com/nevalang/neva/releases/download/$latest_tag/neva-$platform" -o "$binary_path"
chmod +x "$binary_path"
mkdir -p "$install_dir"
install -m 0755 "$binary_path" "$install_dir/neva"

echo "Neva installed successfully to $install_dir/neva"
case ":$PATH:" in
	*":$install_dir:"*) ;;
	*) echo "Add $install_dir to PATH, then restart your terminal or VS Code." ;;
esac
