// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadEnv(t *testing.T) {
	t.Run("empty env", func(t *testing.T) {
		cfg, err := ReadEnv()
		require.NoError(t, err)

		relay := cfg.Rioni.Relay
		require.Nil(t, relay.Upstream)
		require.Equal(t, "5s", relay.Client.TimeoutStr)
		require.Equal(t, 5*time.Second, relay.Client.Timeout())
		require.Equal(t, "1mb", relay.Client.ReadLimitStr)
		require.Equal(t, int64(1024*1024), relay.Client.ReadLimit())
		require.Nil(t, relay.Client.Dns)

		http := cfg.Rioni.Server.Http
		require.Equal(t, "true", http.EnableStr)
		require.True(t, http.IsEnabled())
		require.Equal(t, ":443", http.AddressStr)
		require.Equal(t, ":443", http.Address())
		require.Equal(t, "5s", http.ReadHeaderTimeoutStr)
		require.Equal(t, 5*time.Second, http.ReadHeaderTimeout())
		require.Equal(t, "10s", http.ReadTimeoutStr)
		require.Equal(t, 10*time.Second, http.ReadTimeout())
		require.Equal(t, "64kb", http.ReadLimitStr)
		require.Equal(t, int64(64*1024), http.ReadLimit())
		require.Equal(t, "10s", http.WriteTimeoutStr)
		require.Equal(t, 10*time.Second, http.WriteTimeout())
		require.Equal(t, "30s", http.IdleTimeoutStr)
		require.Equal(t, 30*time.Second, http.IdleTimeout())

		tls := http.Tls
		require.Zero(t, tls.CertFileStr)
		require.Zero(t, tls.CertFile())
		require.Zero(t, tls.KeyFileStr)
		require.Zero(t, tls.KeyFile())
		require.Equal(t, "true", tls.BuildSelfSignedStr)
		require.True(t, tls.IsBuildSelfSigned())

		dns := cfg.Rioni.Server.Dns
		require.Equal(t, "true", dns.EnableStr)
		require.True(t, dns.IsEnabled())
		require.Equal(t, ":53", dns.AddressStr)
		require.Equal(t, ":53", dns.Address())
		require.Equal(t, "2s", dns.ReadTimeoutStr)
		require.Equal(t, 2*time.Second, dns.ReadTimeout())
		require.Equal(t, "2s", dns.WriteTimeoutStr)
		require.Equal(t, 2*time.Second, dns.WriteTimeout())

		log := cfg.Rioni.Log
		require.Equal(t, "json", log.FormatStr)
		require.Equal(t, "json", log.Format())
	})

	t.Run("non-empty env", func(t *testing.T) {
		t.Setenv(EnvRioniRelayUpstream, "https://upstream-1.example,https://upstream-2.example")
		t.Setenv(EnvRioniRelayClientDns, "1.1.1.1:53,8.8.8.8:53")
		t.Setenv(EnvRioniRelayClientTimeout, "7s")
		t.Setenv(EnvRioniRelayClientReadLimit, "2mb")

		t.Setenv(EnvRioniServerHttpEnable, "false")
		t.Setenv(EnvRioniServerHttpAddress, "127.0.0.1:8443")
		t.Setenv(EnvRioniServerHttpReadHeaderTimeout, "6s")
		t.Setenv(EnvRioniServerHttpReadTimeout, "11s")
		t.Setenv(EnvRioniServerHttpReadLimit, "96kb")
		t.Setenv(EnvRioniServerHttpWriteTimeout, "12s")
		t.Setenv(EnvRioniServerHttpIdleTimeout, "33s")
		t.Setenv(EnvRioniServerHttpTlsCertFile, "/tmp/cert.pem")
		t.Setenv(EnvRioniServerHttpTlsKeyFile, "/tmp/key.pem")
		t.Setenv(EnvRioniServerHttpTlsBuildSelfSigned, "false")

		t.Setenv(EnvRioniServerDnsEnable, "false")
		t.Setenv(EnvRioniServerDnsAddress, "127.0.0.1:1053")
		t.Setenv(EnvRioniServerDnsReadTimeout, "3s")
		t.Setenv(EnvRioniServerDnsWriteTimeout, "4s")
		t.Setenv(EnvRioniLogFormat, "text")

		cfg, err := ReadEnv()
		require.NoError(t, err)

		require.Equal(t, []string{"https://upstream-1.example", "https://upstream-2.example"}, cfg.Rioni.Relay.Upstream)
		require.Equal(t, []string{"1.1.1.1:53", "8.8.8.8:53"}, cfg.Rioni.Relay.Client.Dns)
		require.Equal(t, "7s", cfg.Rioni.Relay.Client.TimeoutStr)
		require.Equal(t, 7*time.Second, cfg.Rioni.Relay.Client.Timeout())
		require.Equal(t, "2mb", cfg.Rioni.Relay.Client.ReadLimitStr)
		require.Equal(t, int64(2*1024*1024), cfg.Rioni.Relay.Client.ReadLimit())

		http := cfg.Rioni.Server.Http
		require.Equal(t, "false", http.EnableStr)
		require.False(t, http.IsEnabled())
		require.Equal(t, "127.0.0.1:8443", http.AddressStr)
		require.Equal(t, "127.0.0.1:8443", http.Address())
		require.Equal(t, "6s", http.ReadHeaderTimeoutStr)
		require.Equal(t, 6*time.Second, http.ReadHeaderTimeout())
		require.Equal(t, "11s", http.ReadTimeoutStr)
		require.Equal(t, 11*time.Second, http.ReadTimeout())
		require.Equal(t, "96kb", http.ReadLimitStr)
		require.Equal(t, int64(96*1024), http.ReadLimit())
		require.Equal(t, "12s", http.WriteTimeoutStr)
		require.Equal(t, 12*time.Second, http.WriteTimeout())
		require.Equal(t, "33s", http.IdleTimeoutStr)
		require.Equal(t, 33*time.Second, http.IdleTimeout())
		require.Equal(t, "/tmp/cert.pem", http.Tls.CertFileStr)
		require.Equal(t, "/tmp/cert.pem", http.Tls.CertFile())
		require.Equal(t, "/tmp/key.pem", http.Tls.KeyFileStr)
		require.Equal(t, "/tmp/key.pem", http.Tls.KeyFile())
		require.Equal(t, "false", http.Tls.BuildSelfSignedStr)
		require.False(t, http.Tls.IsBuildSelfSigned())

		dns := cfg.Rioni.Server.Dns
		require.Equal(t, "false", dns.EnableStr)
		require.False(t, dns.IsEnabled())
		require.Equal(t, "127.0.0.1:1053", dns.AddressStr)
		require.Equal(t, "127.0.0.1:1053", dns.Address())
		require.Equal(t, "3s", dns.ReadTimeoutStr)
		require.Equal(t, 3*time.Second, dns.ReadTimeout())
		require.Equal(t, "4s", dns.WriteTimeoutStr)
		require.Equal(t, 4*time.Second, dns.WriteTimeout())

		log := cfg.Rioni.Log
		require.Equal(t, "text", log.FormatStr)
		require.Equal(t, "text", log.Format())
	})
}

func TestReadEnvToTarget(t *testing.T) {
	cfg, err := ReadFile("../../configs/rioni.cfg.yml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotZero(t, cfg)

	t.Setenv(EnvRioniRelayUpstream, "https://upstream-env.example")
	t.Setenv(EnvRioniRelayClientTimeout, "8s")
	t.Setenv(EnvRioniRelayClientReadLimit, "3mb")
	t.Setenv(EnvRioniServerHttpEnable, "false")
	t.Setenv(EnvRioniServerHttpAddress, "127.0.0.1:9443")
	t.Setenv(EnvRioniServerHttpReadHeaderTimeout, "7s")
	t.Setenv(EnvRioniServerHttpReadTimeout, "15s")
	t.Setenv(EnvRioniServerHttpReadLimit, "128kb")
	t.Setenv(EnvRioniServerHttpWriteTimeout, "13s")
	t.Setenv(EnvRioniServerHttpIdleTimeout, "35s")
	t.Setenv(EnvRioniServerHttpTlsBuildSelfSigned, "false")
	t.Setenv(EnvRioniServerDnsEnable, "false")
	t.Setenv(EnvRioniServerDnsAddress, "127.0.0.1:2053")
	t.Setenv(EnvRioniServerDnsReadTimeout, "5s")
	t.Setenv(EnvRioniServerDnsWriteTimeout, "6s")
	t.Setenv(EnvRioniLogFormat, "text")

	cfg, err = ReadEnvToTarget(cfg)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	relay := cfg.Rioni.Relay
	require.Equal(t, []string{"https://upstream-env.example"}, relay.Upstream)
	require.Equal(t, "8s", relay.Client.TimeoutStr)
	require.Equal(t, 8*time.Second, relay.Client.Timeout())
	require.Equal(t, "3mb", relay.Client.ReadLimitStr)
	require.Equal(t, int64(3*1024*1024), relay.Client.ReadLimit())
	require.Equal(t, []string{"8.8.8.8", "1.1.1.1"}, relay.Client.Dns)

	http := cfg.Rioni.Server.Http
	require.Equal(t, "false", http.EnableStr)
	require.False(t, http.IsEnabled())
	require.Equal(t, "127.0.0.1:9443", http.AddressStr)
	require.Equal(t, "127.0.0.1:9443", http.Address())
	require.Equal(t, "7s", http.ReadHeaderTimeoutStr)
	require.Equal(t, 7*time.Second, http.ReadHeaderTimeout())
	require.Equal(t, "15s", http.ReadTimeoutStr)
	require.Equal(t, 15*time.Second, http.ReadTimeout())
	require.Equal(t, "128kb", http.ReadLimitStr)
	require.Equal(t, int64(128*1024), http.ReadLimit())
	require.Equal(t, "13s", http.WriteTimeoutStr)
	require.Equal(t, 13*time.Second, http.WriteTimeout())
	require.Equal(t, "35s", http.IdleTimeoutStr)
	require.Equal(t, 35*time.Second, http.IdleTimeout())
	require.Equal(t, "tls/rioni.crt", http.Tls.CertFileStr)
	require.Equal(t, "tls/rioni.crt", http.Tls.CertFile())
	require.Equal(t, "tls/rioni.key", http.Tls.KeyFileStr)
	require.Equal(t, "tls/rioni.key", http.Tls.KeyFile())
	require.Equal(t, "false", http.Tls.BuildSelfSignedStr)
	require.False(t, http.Tls.IsBuildSelfSigned())

	dns := cfg.Rioni.Server.Dns
	require.Equal(t, "false", dns.EnableStr)
	require.False(t, dns.IsEnabled())
	require.Equal(t, "127.0.0.1:2053", dns.AddressStr)
	require.Equal(t, "127.0.0.1:2053", dns.Address())
	require.Equal(t, "5s", dns.ReadTimeoutStr)
	require.Equal(t, 5*time.Second, dns.ReadTimeout())
	require.Equal(t, "6s", dns.WriteTimeoutStr)
	require.Equal(t, 6*time.Second, dns.WriteTimeout())

	log := cfg.Rioni.Log
	require.Equal(t, "text", log.FormatStr)
	require.Equal(t, "text", log.Format())
}

func TestReadFile(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		cfg, err := ReadFile("../../configs/rioni.cfg.yml")
		require.NoError(t, err)
		require.NotNil(t, cfg)

		relay := cfg.Rioni.Relay
		require.Equal(t, []string{"https://dns.google/dns-query", "https://cloudflare-dns.com/dns-query"}, relay.Upstream)
		require.Zero(t, relay.Client.TimeoutStr)
		require.Zero(t, relay.Client.Timeout())
		require.Zero(t, relay.Client.ReadLimitStr)
		require.Zero(t, relay.Client.ReadLimit())
		require.Equal(t, []string{"8.8.8.8", "1.1.1.1"}, relay.Client.Dns)

		http := cfg.Rioni.Server.Http
		require.Zero(t, http.EnableStr)
		require.Zero(t, http.IsEnabled())
		require.Equal(t, ":443", http.AddressStr)
		require.Equal(t, ":443", http.Address())
		require.Zero(t, http.ReadHeaderTimeoutStr)
		require.Zero(t, http.ReadHeaderTimeout())
		require.Zero(t, http.ReadTimeoutStr)
		require.Zero(t, http.ReadTimeout())
		require.Zero(t, http.ReadLimitStr)
		require.Zero(t, http.ReadLimit())
		require.Zero(t, http.WriteTimeoutStr)
		require.Zero(t, http.WriteTimeout())
		require.Zero(t, http.IdleTimeoutStr)
		require.Zero(t, http.IdleTimeout())

		tls := http.Tls
		require.Equal(t, "tls/rioni.crt", tls.CertFileStr)
		require.Equal(t, "tls/rioni.crt", tls.CertFile())
		require.Equal(t, "tls/rioni.key", tls.KeyFileStr)
		require.Equal(t, "tls/rioni.key", tls.KeyFile())
		require.Zero(t, tls.BuildSelfSignedStr)
		require.Zero(t, tls.IsBuildSelfSigned())

		dns := cfg.Rioni.Server.Dns
		require.Zero(t, dns.EnableStr)
		require.Zero(t, dns.IsEnabled())
		require.Equal(t, ":53", dns.AddressStr)
		require.Equal(t, ":53", dns.Address())
		require.Zero(t, dns.ReadTimeoutStr)
		require.Zero(t, dns.ReadTimeout())
		require.Zero(t, dns.WriteTimeoutStr)
		require.Zero(t, dns.WriteTimeout())

		log := cfg.Rioni.Log
		require.Empty(t, log.FormatStr)
		require.Empty(t, log.Format())
	})

	t.Run("file does not exist", func(t *testing.T) {
		cfg, err := ReadFile("non-existent.yaml")
		require.Error(t, err)
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Zero(t, cfg)
	})
}
