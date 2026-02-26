#!/usr/bin/env bash

# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

# DNS Server Test Script
# Tests DNS requests using dig utility

SERVER="127.0.0.1"
PORT="8053"
DOMAIN="example.com"

PASSED=0
FAILED=0

echo "==================================="
echo "DNS Server Test Script"
echo "Testing server: ${SERVER}:${PORT}"
echo "==================================="
echo ""

run_test() {
  local test_name=$1
  local record_type=$2

  echo "Test: ${test_name}"
  echo "Query: ${DOMAIN} (${record_type})"
  echo "-----------------------------------"

  dig @localhost -p ${PORT} ${DOMAIN} "${record_type}" +dnssec
  exit_code=$?

  if [ "$exit_code" -eq 0 ]; then
    PASSED=$((PASSED + 1))
    echo "✓ Test passed"
  else
    FAILED=$((FAILED + 1))
    echo "✗ Test failed"
  fi
  echo ""
}

run_test "A Record Query" "A"
run_test "AAAA Record Query" "AAAA"
run_test "MX Record Query" "MX"
run_test "TXT Record Query" "TXT"
run_test "NS Record Query" "NS"

echo ""
echo "==================================="
echo "Tests completed"
echo "Passed: ${PASSED}, Failed: ${FAILED}"
echo "==================================="

if [ "$FAILED" -gt 0 ]; then
  exit 1
fi

exit 0
