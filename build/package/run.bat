@echo off
REM Copyright (c) 2026 Andrey Senov
REM SPDX-License-Identifier: Apache-2.0

REM Rioni DNS Proxy Server - Launch Script

setlocal

set "SCRIPT_DIR=%~dp0"

set "BINARY=%SCRIPT_DIR%rioni.exe"
set "CONFIG=%SCRIPT_DIR%configs\rioni.cfg.yml"

if not exist "%BINARY%" (
    echo Error: Binary not found at %BINARY% >&2
    exit /b 1
)

if not exist "%CONFIG%" (
    echo Error: Config file not found at %CONFIG% >&2
    exit /b 1
)

cd /d "%SCRIPT_DIR%"

"%BINARY%" --config "%CONFIG%"
