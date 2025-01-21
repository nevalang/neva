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
            echo "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    echo "${os}-${arch}"
}

# Determine latest release tag
LATEST_TAG=$(curl -s https://api.github.com/repos/nevalang/neva/releases/latest | grep "tag_name" | cut -d '"' -f 4)
echo "Latest tag is $LATEST_TAG"

# Determine platform
PLATFORM=$(detect_platform)
echo "Platform is $PLATFORM"

# Build the release URL
BIN_NAME="neva"
BIN_URL="https://github.com/nevalang/neva/releases/download/$LATEST_TAG/${BIN_NAME}-${PLATFORM}"

# Download the binary
echo "Downloading..."
curl -L $BIN_URL -o $BIN_NAME

# Make the binary executable
chmod +x $BIN_NAME

# Move the binary to a location in the user's PATH
INSTALL_DIR="/usr/local/bin"
sudo mv $BIN_NAME $INSTALL_DIR
echo "$BIN_NAME installed successfully to $INSTALL_DIR"
