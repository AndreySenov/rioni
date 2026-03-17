#!/bin/sh
# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

set -e

REPO="AndreySenov/rioni"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"

ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  ARCH_TYPE="amd64" ;;
    aarch64) ARCH_TYPE="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

if command -v dpkg >/dev/null; then
    PKG_EXT="deb"
    INSTALL_CMD="sudo dpkg -i"
elif command -v rpm >/dev/null; then
    PKG_EXT="rpm"
    INSTALL_CMD="sudo rpm -i"
else
    echo "No supported package manager found (dpkg or rpm)."
    exit 1
fi

echo "Detecting latest version for $ARCH_TYPE ($PKG_EXT)..."

DOWNLOAD_URL=$(curl -s $GITHUB_API | grep "browser_download_url" | grep "$ARCH_TYPE" | grep "\.$PKG_EXT" | cut -d '"' -f 4 | head -n 1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find a suitable package for $ARCH_TYPE and $PKG_EXT."
    exit 1
fi

TEMP_PKG="/tmp/rioni_latest.$PKG_EXT"
echo "Downloading $DOWNLOAD_URL..."
curl -sL "$DOWNLOAD_URL" -o "$TEMP_PKG"

echo "Installing Rioni..."
$INSTALL_CMD "$TEMP_PKG"

rm "$TEMP_PKG"
echo "Rioni installed successfully!"
echo "To start the service, run: sudo systemctl enable --now rioni"
