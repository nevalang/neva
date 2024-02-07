#!/bin/bash

BIN_URL="https://github.com/nevalang/neva/releases/download/v0.3.0/neva"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="neva"

echo "Downloading $BIN_NAME..."
curl -L $BIN_URL -o $BIN_NAME
chmod +x $BIN_NAME
mv $BIN_NAME $INSTALL_DIR
echo "$BIN_NAME installed successfully to $INSTALL_DIR"

# example usage:
# curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/install.sh | bash
