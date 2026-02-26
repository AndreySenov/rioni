#!/bin/sh
# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

set -e

NOLOGIN="/usr/sbin/nologin"
if [ ! -x "$NOLOGIN" ]; then
    NOLOGIN="/sbin/nologin"
fi
if [ ! -x "$NOLOGIN" ]; then
    NOLOGIN="/bin/false"
fi

if ! getent group rioni >/dev/null 2>&1; then
    groupadd --system rioni
fi

if ! getent passwd rioni >/dev/null 2>&1; then
    useradd \
        --system \
        --gid rioni \
        --home-dir /var/lib/rioni \
        --no-create-home \
        --shell "$NOLOGIN" \
        rioni
fi
