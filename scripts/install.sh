#!/bin/bash

LATEST_TAG=$(curl -s https://api.github.com/repos/nevalang/neva/releases/latest | grep "tag_name" | cut -d '"' -f 4)
BIN_URL="https://github.com/nevalang/neva/releases/download/$LATEST_TAG/neva"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="neva"

echo "Downloading $BIN_NAME..."
curl -L $BIN_URL -o $BIN_NAME
chmod +x $BIN_NAME
sudo mv $BIN_NAME $INSTALL_DIR
echo "$BIN_NAME installed successfully to $INSTALL_DIR"

# example usage:
# curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash
