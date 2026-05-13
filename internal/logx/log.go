// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package logx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AndreySenov/rioni/internal/httpx"
	"github.com/miekg/dns"
)

func WarnForHttpRequest(logger *slog.Logger, httpRequest *http.Request, message string, args ...any) {
	logForHttpRequest(logger, httpRequest, slog.LevelWarn, message, args...)
}

func ErrorForHttpRequest(logger *slog.Logger, httpRequest *http.Request, message string, args ...any) {
	logForHttpRequest(logger, httpRequest, slog.LevelError, message, args...)
}

func logForHttpRequest(logger *slog.Logger, httpRequest *http.Request, level slog.Level, message string, args ...any) {
	if logger == nil {
		logger = slog.Default()
	}

	ctx := context.Background()
	if httpRequest != nil {
		ctx = httpRequest.Context()

		requestAttrs := httpx.GetRequestAttrs(httpRequest)
		if requestAttrs != nil && len(requestAttrs) > 0 {
			r := slog.Group("http_request", requestAttrs...)
			args = append(args, r)
		}
	}

	logger.Log(ctx, level, message, args...)
}

func ErrorForDnsRequest(logger *slog.Logger, w dns.ResponseWriter, message string, args ...any) {
	logForDnsRequest(logger, w, slog.LevelError, message, args...)
}

func logForDnsRequest(logger *slog.Logger, w dns.ResponseWriter, level slog.Level, message string, args ...any) {
	if logger == nil {
		logger = slog.Default()
	}

	if w != nil {
		requestAttrs := []any{
			"remote_addr", w.RemoteAddr(),
		}
		args = append(args, requestAttrs...)
	}

	logger.Log(context.Background(), level, message, args...)
}
