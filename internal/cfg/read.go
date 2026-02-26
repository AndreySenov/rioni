// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"os"

	"github.com/a8m/envsubst"
	"gopkg.in/yaml.v3"
)

func Read(path string) (Config, error) {
	return read[Config](path)
}

func read[T any](path string) (T, error) {
	var zero T

	b, err := os.ReadFile(path)
	if err != nil {
		return zero, err
	}

	substituted, err := envsubst.Bytes(b)
	if err != nil {
		return zero, err
	}

	var result T
	if err := yaml.Unmarshal(substituted, &result); err != nil {
		return zero, err
	}

	return result, nil
}
