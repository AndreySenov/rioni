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
    INSTALL_CMD="sudo dnf install -y"
else
    echo "No supported package manager found (dpkg or rpm)."
    exit 1
fi

echo "Detecting latest version for $ARCH_TYPE ($PKG_EXT)..."

RELEASE_DATA=$(curl -s "$GITHUB_API")

DOWNLOAD_URL=$(echo "$RELEASE_DATA" | grep "browser_download_url" | grep "$ARCH_TYPE" | grep "\.$PKG_EXT" | cut -d '"' -f 4 | head -n 1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find a suitable package for $ARCH_TYPE and $PKG_EXT."
    exit 1
fi

CHECKSUM_URL=$(echo "$RELEASE_DATA" | grep "browser_download_url" | grep "$DOWNLOAD_URL\.sha256" | cut -d '"' -f 4 | head -n 1)

TEMP_PKG="/tmp/rioni_latest.$PKG_EXT"
TEMP_SUM="/tmp/rioni_latest.sha256"

echo "Downloading package $DOWNLOAD_URL..."
curl -sL "$DOWNLOAD_URL" -o "$TEMP_PKG"

if [ -n "$CHECKSUM_URL" ]; then
    echo "Verifying SHA256 checksum..."
    curl -sL "$CHECKSUM_URL" -o "$TEMP_SUM"

    EXPECTED_HASH=$(awk '{print $1}' "$TEMP_SUM")
    LOCAL_HASH=$(sha256sum "$TEMP_PKG" | awk '{print $1}')

    if [ "$EXPECTED_HASH" = "$LOCAL_HASH" ]; then
        echo "Checksum verified: OK"
    else
        echo "ERROR: Checksum mismatch!"
        echo "Expected: $EXPECTED_HASH"
        echo "Actual:   $LOCAL_HASH"
        rm -f "$TEMP_PKG" "$TEMP_SUM"
        exit 1
    fi
    rm -f "$TEMP_SUM"
else
    echo "Warning: No .sha256 file found for the package. Skipping verification."
fi

echo "Installing Rioni..."
$INSTALL_CMD "$TEMP_PKG"

rm -f "$TEMP_PKG"
echo "Installation completed."
