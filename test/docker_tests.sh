#!/usr/bin/env bash

# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

HTTP_PORT="8443"
DNS_PORT="8053"

docker run -d \
  --name rioni_docker_test \
  -p ${HTTP_PORT}:443/tcp \
  -p ${DNS_PORT}:53/tcp \
  -p ${DNS_PORT}:53/udp \
  ghcr.io/andreysenov/rioni:latest-arm64 || exit 1

sh ./all_tests.sh

docker stop rioni_docker_test
docker rm rioni_docker_test
