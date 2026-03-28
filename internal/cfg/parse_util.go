// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"strconv"
	"strings"
)

func parseSize(s string) int64 {
	s = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(s)), " ", "")

	if strings.HasSuffix(s, "kb") {
		if val, err := strconv.ParseInt(strings.TrimSuffix(s, "kb"), 10, 64); err == nil {
			return val << 10
		}
	}

	if strings.HasSuffix(s, "mb") {
		if val, err := strconv.ParseInt(strings.TrimSuffix(s, "mb"), 10, 64); err == nil {
			return val << 20
		}
	}

	if strings.HasSuffix(s, "gb") {
		if val, err := strconv.ParseInt(strings.TrimSuffix(s, "gb"), 10, 64); err == nil {
			return val << 30
		}
	}

	if strings.HasSuffix(s, "b") {
		if val, err := strconv.ParseInt(strings.TrimSuffix(s, "b"), 10, 64); err == nil {
			return val
		}
	}

	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}

	return 0
}
