// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	RioniConfigFileEnvVar            = "RIONI_CONFIG_FILE"
	RioniServerDnsAddressEnvVar      = "RIONI_SERVER_DNS_ADDRESS"
	RioniServerHttpAddressEnvVar     = "RIONI_SERVER_HTTP_ADDRESS"
	RioniServerHttpTlsCertFileEnvVar = "RIONI_SERVER_HTTP_TLS_CERT_FILE"
	RioniServerHttpTlsKeyFileEnvVar  = "RIONI_SERVER_HTTP_TLS_KEY_FILE"
)

func ResolvePath(configFlag *flag.Flag) (string, error) {
	if path, set, err := resolveEnvVar(); set {
		if err == nil {
			return path, nil
		}
		return "", err
	}

	if path, set, err := resolveFlag(configFlag); set {
		if err == nil {
			return path, nil
		}
		return "", err
	}

	return "", fmt.Errorf("config path is not provided: use --config CLI option or %s environment variable", RioniConfigFileEnvVar)
}

func resolveFlag(configFlag *flag.Flag) (string, bool, error) {
	if configFlag == nil || configFlag.Value.String() == "" {
		return "", false, nil
	}

	path, err := resolve(configFlag.Value.String())
	if err != nil {
		return "", true, err
	}

	return path, true, nil
}

func resolveEnvVar() (string, bool, error) {
	if env, ok := os.LookupEnv(RioniConfigFileEnvVar); ok {
		if env == "" {
			return "", true, fmt.Errorf("environment variable %q is set but empty", RioniConfigFileEnvVar)
		}

		path, err := resolve(env)
		if err != nil {
			return "", true, err
		}

		return path, true, nil
	}
	return "", false, nil
}

func resolve(path string) (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("configuration path %q: %w", path, err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("configuration path %q is a directory", path)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("configuration path %q: %w", path, err)
	}

	return abs, nil
}
