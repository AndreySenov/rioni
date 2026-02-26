#!/bin/sh
# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

set -e

if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload >/dev/null 2>&1 || true
fi

echo
echo "Rioni has been installed."
echo
echo "Review the configuration:"
echo "  /etc/rioni/rioni.cfg.yml"
echo
echo "Then enable and start the service with:"
echo "  sudo systemctl enable --now rioni"
echo
echo "To check status:"
echo "  systemctl status rioni"
echo
