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
			echo "不支持的架构: $arch" >&2
			exit 1
			;;
	esac
	printf '%s-%s\n' "$os" "$arch"
}

latest_tag=$(curl -fsSL https://api.github.com/repos/nevalang/neva/releases/latest | grep '"tag_name"' | cut -d '"' -f 4)
if [[ -z "$latest_tag" ]]; then
	echo '无法确定最新 Neva 版本。' >&2
	exit 1
fi
echo "最新版本: $latest_tag"

platform=$(detect_platform)
echo "平台: $platform"

install_dir="${NEVA_INSTALL_DIR:-$HOME/.local/bin}"
binary_path="$(mktemp)"
trap 'rm -f "$binary_path"' EXIT

echo '下载中...'
curl -fL "https://github.moeyy.xyz/https://github.com/nevalang/neva/releases/download/$latest_tag/neva-$platform" -o "$binary_path"
chmod +x "$binary_path"
mkdir -p "$install_dir"
install -m 0755 "$binary_path" "$install_dir/neva"

echo "Neva 已成功安装到 $install_dir/neva"
case ":$PATH:" in
	*":$install_dir:"*) ;;
	*) echo "请将 $install_dir 添加到 PATH，然后重启终端或 VS Code。" ;;
esac
