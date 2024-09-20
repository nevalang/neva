#!/bin/bash

# Function to detect platform and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    case $arch in
        x86_64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        loong64|loongarch64)
            arch="loong64"
            ;;
        *)
            echo "不支持的架构: $arch"
            exit 1
            ;;
    esac
    echo "${os}-${arch}"
}

# Determine latest release tag
LATEST_TAG=$(curl -s https://api.github.com/repos/nevalang/neva/releases/latest | grep "tag_name" | cut -d '"' -f 4)
echo "最新版本: $LATEST_TAG"

# Determine platform
PLATFORM=$(detect_platform)
echo "平台: $PLATFORM"

# Build the release URL
BIN_NAME="neva"
BIN_URL="https://github.moeyy.xyz/https://github.com/nevalang/neva/releases/download/$LATEST_TAG/${BIN_NAME}-${PLATFORM}"

# Download the binary
echo "下载中..."
curl -L $BIN_URL -o $BIN_NAME

# Make the binary executable
chmod +x $BIN_NAME

# Move the binary to a location in the user's PATH
INSTALL_DIR="/usr/local/bin"
sudo mv $BIN_NAME $INSTALL_DIR
echo "已经将 $BIN_NAME 成功安装到 $INSTALL_DIR"
