#!/bin/bash

# Linutils Rakesh One-Line Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/rak626/linutils-rakesh/main/install.sh | bash

set -e

REPO="rak626/linutils-rakesh"
BINARY_NAME="linutils-rakesh"
INSTALL_DIR="$HOME/.local/bin"

echo "--- Linutils Rakesh Installer ---"

# 1. Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# 2. Detect Architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# 3. Fetch latest release from GitHub API
echo "Fetching latest version info..."
LATEST_RELEASE=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep "tag_name" | cut -d '"' -f 4)

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not find latest release. Falling back to manual build check..."
    # Fallback logic if release doesn't exist yet: tell user to build from source
    echo "No pre-compiled releases found yet. Please build from source:"
    echo "git clone https://github.com/$REPO.git && cd linutils-rakesh && go build -o linutils-rakesh main.go"
    exit 1
fi

echo "Latest version: $LATEST_RELEASE"

# 4. Download binary
URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}-linux-${ARCH}"
echo "Downloading from $URL..."

curl -L "$URL" -o "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# 5. Ensure bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Adding $INSTALL_DIR to PATH in .bashrc and .zshrc..."
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    [ -f "$HOME/.zshrc" ] && echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc"
    echo "Please restart your terminal or run 'source ~/.bashrc' to use the command."
fi

echo "--- Installation Complete! ---"
echo "You can now run the tool by typing: $BINARY_NAME"
