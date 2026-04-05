#!/usr/bin/env bash

# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

# Runs all tests

STATUS=0

sh ./dns_udp_tests.sh
dns_udp_tests_exit_code=$?

sh ./dns_tcp_tests.sh
dns_tcp_tests_exit_code=$?

sh ./doh_get_tests.sh
doh_get_tests_exit_code=$?

sh ./doh_post_tests.sh
doh_post_tests_exit_code=$?

if [ "$dns_udp_tests_exit_code" -eq 0 ]; then
  echo "DNS (UDP) tests PASSED"
else
  STATUS=1
  echo "DNS (UDP) tests FAILED"
fi

if [ "$dns_tcp_tests_exit_code" -eq 0 ]; then
  echo "DNS (TCP) tests PASSED"
else
  STATUS=1
  echo "DNS (TCP) tests FAILED"
fi

if [ "$doh_get_tests_exit_code" -eq 0 ]; then
  echo "DNS-over-HTTPS (GET) tests PASSED"
else
  STATUS=1
  echo "DNS-over-HTTPS (GET) tests FAILED"
fi

if [ "$doh_post_tests_exit_code" -eq 0 ]; then
  echo "DNS-over-HTTPS (POST) tests PASSED"
else
  STATUS=1
  echo "DNS-over-HTTPS (POST) tests FAILED"
fi

if [ "$STATUS" -eq 0 ]; then
  echo "All tests passed"
  exit 0
else
  exit 1
fi
