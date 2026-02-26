#!/usr/bin/env bash

# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

# DNS-over-HTTPS Server Test Script
# Tests DNS requests over HTTPS using dig utility

SERVER="localhost"
PORT="8443"
DOMAIN="example.com"

PASSED=0
FAILED=0

echo "==================================="
echo "DNS-over-HTTPS (POST) Server Test Script"
echo "Testing server: ${SERVER}:${PORT}"
echo "==================================="
echo ""

run_post_test() {
  local test_name=$1
  local record_type=$2

  echo "Test: POST ${test_name}"
  echo "Query: ${DOMAIN} (${record_type})"
  echo "-----------------------------------"

  dig @localhost -p ${PORT} +https ${DOMAIN} "${record_type}" +dnssec
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

run_post_test "A Record Query" "A"
run_post_test "AAAA Record Query" "AAAA"
run_post_test "MX Record Query" "MX"
run_post_test "TXT Record Query" "TXT"
run_post_test "NS Record Query" "NS"

echo ""
echo "==================================="
echo "DoH POST tests completed"
echo "Passed: ${PASSED}, Failed: ${FAILED}"
echo "==================================="
