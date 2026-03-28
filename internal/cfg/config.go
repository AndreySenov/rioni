// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"time"
)

type Config struct {
	Rioni Rioni `yaml:"rioni" envPrefix:"RIONI_"`
}

type Rioni struct {
	Relay  Relay  `yaml:"relay" envPrefix:"RELAY_"`
	Server Server `yaml:"server" envPrefix:"SERVER_"`
}

type Relay struct {
	Upstream []string     `yaml:"upstream" env:"UPSTREAM" envSeparator:","`
	Client   ClientConfig `yaml:"client" envPrefix:"CLIENT_"`
}

type ClientConfig struct {
	TimeoutStr   string   `yaml:"timeout" env:"TIMEOUT" envDefault:"5s"`
	ReadLimitStr string   `yaml:"read-limit" env:"READ_LIMIT" envDefault:"1mb"`
	Dns          []string `yaml:"dns" env:"DNS" envSeparator:","`
}

func (c ClientConfig) Timeout() time.Duration {
	d, _ := time.ParseDuration(c.TimeoutStr)
	return d
}

func (c ClientConfig) ReadLimit() int64 {
	return parseSize(c.ReadLimitStr)
}

type Server struct {
	Http HttpConfig `yaml:"http" envPrefix:"HTTP_"`
	Dns  DnsConfig  `yaml:"dns" envPrefix:"DNS_"`
}

type HttpConfig struct {
	EnableStr            string    `yaml:"enable" env:"ENABLE" envDefault:"true"`
	AddressStr           string    `yaml:"address" env:"ADDRESS" envDefault:":443"`
	ReadHeaderTimeoutStr string    `yaml:"read-header-timeout" env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeoutStr       string    `yaml:"read-timeout" env:"READ_TIMEOUT" envDefault:"10s"`
	ReadLimitStr         string    `yaml:"read-limit" env:"READ_LIMIT" envDefault:"64kb"`
	WriteTimeoutStr      string    `yaml:"write-timeout" env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeoutStr       string    `yaml:"idle-timeout" env:"IDLE_TIMEOUT" envDefault:"30s"`
	Tls                  TlsConfig `yaml:"tls" envPrefix:"TLS_"`
}

func (h HttpConfig) IsEnabled() bool {
	return parseBool(h.EnableStr)
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
	CertFileStr        string `yaml:"cert-file" env:"CERT_FILE"`
	KeyFileStr         string `yaml:"key-file" env:"KEY_FILE"`
	BuildSelfSignedStr string `yaml:"build-self-signed" env:"BUILD_SELF_SIGNED" envDefault:"true"`
}

func (t TlsConfig) CertFile() string {
	return t.CertFileStr
}

func (t TlsConfig) KeyFile() string {
	return t.KeyFileStr
}

func (t TlsConfig) IsBuildSelfSigned() bool {
	return parseBool(t.BuildSelfSignedStr)
}

type DnsConfig struct {
	EnableStr       string `yaml:"enable" env:"ENABLE" envDefault:"true"`
	AddressStr      string `yaml:"address" env:"ADDRESS" envDefault:":53"`
	ReadTimeoutStr  string `yaml:"read-timeout" env:"READ_TIMEOUT" envDefault:"2s"`
	WriteTimeoutStr string `yaml:"write-timeout" env:"WRITE_TIMEOUT" envDefault:"2s"`
}

func (d DnsConfig) IsEnabled() bool {
	return parseBool(d.EnableStr)
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
