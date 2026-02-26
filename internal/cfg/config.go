// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Rioni Rioni `yaml:"rioni"`
}

type Rioni struct {
	Relay  Relay  `yaml:"relay"`
	Server Server `yaml:"server"`
}

type Relay struct {
	Upstream []string     `yaml:"upstream"`
	Client   ClientConfig `yaml:"client"`
}

type ClientConfig struct {
	TimeoutStr   string   `yaml:"timeout"`
	ReadLimitStr string   `yaml:"read-limit"`
	Dns          []string `yaml:"dns"`
}

func (c ClientConfig) Timeout() time.Duration {
	d, _ := time.ParseDuration(c.TimeoutStr)
	return d
}

func (c ClientConfig) ReadLimit() int64 {
	return parseSize(c.ReadLimitStr)
}

type Server struct {
	Http HttpConfig `yaml:"http"`
	Dns  DnsConfig  `yaml:"dns"`
}

type HttpConfig struct {
	AddressStr           string    `yaml:"address"`
	ReadHeaderTimeoutStr string    `yaml:"read-header-timeout"`
	ReadTimeoutStr       string    `yaml:"read-timeout"`
	ReadLimitStr         string    `yaml:"read-limit"`
	WriteTimeoutStr      string    `yaml:"write-timeout"`
	IdleTimeoutStr       string    `yaml:"idle-timeout"`
	Tls                  TlsConfig `yaml:"tls"`
}

func (h HttpConfig) Address() string {
	return h.AddressStr
}

func (h HttpConfig) ReadHeaderTimeout() time.Duration {
	d, _ := time.ParseDuration(h.ReadHeaderTimeoutStr)
	return d
}

func (h HttpConfig) ReadTimeout() time.Duration {
	d, _ := time.ParseDuration(h.ReadTimeoutStr)
	return d
}

func (h HttpConfig) ReadLimit() int64 {
	return parseSize(h.ReadLimitStr)
}

func (h HttpConfig) WriteTimeout() time.Duration {
	d, _ := time.ParseDuration(h.WriteTimeoutStr)
	return d
}

func (h HttpConfig) IdleTimeout() time.Duration {
	d, _ := time.ParseDuration(h.IdleTimeoutStr)
	return d
}

type TlsConfig struct {
	CertFileStr     string `yaml:"cert-file"`
	KeyFileStr      string `yaml:"key-file"`
	BuildSelfSigned bool   `yaml:"build-self-signed"`
}

func (t TlsConfig) CertFile() string {
	return t.CertFileStr
}

func (t TlsConfig) KeyFile() string {
	return t.KeyFileStr
}

func (t TlsConfig) IsBuildSelfSigned() bool {
	return t.BuildSelfSigned
}

type DnsConfig struct {
	AddressStr      string `yaml:"address"`
	ReadTimeoutStr  string `yaml:"read-timeout"`
	WriteTimeoutStr string `yaml:"write-timeout"`
}

func (d DnsConfig) Address() string {
	return d.AddressStr
}

func (d DnsConfig) ReadTimeout() time.Duration {
	dur, _ := time.ParseDuration(d.ReadTimeoutStr)
	return dur
}

func (d DnsConfig) WriteTimeout() time.Duration {
	dur, _ := time.ParseDuration(d.WriteTimeoutStr)
	return dur
}

func parseSize(s string) int64 {
	s = strings.TrimSpace(strings.ToLower(s))

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

	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}

	return 0
}
