// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package logx

import (
	"log/slog"
	"os"
)

func NewLog(format string, app string, appVersion string) *slog.Logger {
	switch format {
	case "json":
		return NewJsonLog(app, appVersion)
	default:
		return NewTextLog(app, appVersion)
	}
}

func NewJsonLog(app string, appVersion string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	return newLog(handler, app, appVersion)
}

func NewTextLog(app string, appVersion string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	return newLog(handler, app, appVersion)
}

func newLog(handler slog.Handler, app string, appVersion string) *slog.Logger {
	return slog.New(handler).With(AppKey, app, AppVersionKey, appVersion)
}
