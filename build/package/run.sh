#!/usr/bin/env bash
# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

# Rioni DNS Proxy Server - Launch Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BINARY="${SCRIPT_DIR}/rioni"
CONFIG="${SCRIPT_DIR}/configs/rioni.cfg.yml"

if [ ! -f "${BINARY}" ]; then
    echo "Error: Binary not found at ${BINARY}" >&2
    exit 1
fi

if [ ! -x "${BINARY}" ]; then
    echo "Error: Binary at ${BINARY} is not executable" >&2
    exit 1
fi

if [ ! -f "${CONFIG}" ]; then
    echo "Error: Config file not found at ${CONFIG}" >&2
    exit 1
fi

cd "${SCRIPT_DIR}"

exec "${BINARY}" --config "${CONFIG}"
