// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"os"

	"github.com/a8m/envsubst"
	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

func ReadEnv() (Config, error) {
	return env.ParseAs[Config]()
}

func ReadEnvToTarget(target Config) (Config, error) {
	err := env.Parse(&target)
	if err != nil {
		var zero Config
		return zero, err
	}

	return target, nil
}

func ReadFile(path string) (Config, error) {
	return readYamlFile[Config](path)
}

func readYamlFile[T any](path string) (T, error) {
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
