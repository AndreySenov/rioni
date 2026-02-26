#!/usr/bin/env bash

# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

# Runs all tests

STATUS=0

sh ./dns_tests.sh
dns_tests_exit_code=$?

sh ./doh_get_tests.sh
doh_get_tests_exit_code=$?

sh ./doh_post_tests.sh
doh_post_tests_exit_code=$?

if [ "$dns_tests_exit_code" -eq 0 ]; then
  echo "DNS tests PASSED"
else
  STATUS=1
  echo "DNS tests FAILED"
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
