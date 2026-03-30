// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := config()
		err := Check(config)
		require.NoError(t, err)
	})

	t.Run("missing upstream", func(t *testing.T) {
		config := config()
		config.Rioni.Relay.Upstream = nil
		err := Check(config)
		require.EqualError(t, err, "relay configuration error: at least one upstream DNS-over-HTTPS server must be configured")
	})

	t.Run("missing client dns", func(t *testing.T) {
		config := config()
		config.Rioni.Relay.Client.Dns = nil
		err := Check(config)
		require.EqualError(t, err, "relay client configuration error: at least one DNS server must be configured to resolve upstream domains")
	})

	t.Run("http enabled missing address", func(t *testing.T) {
		config := config()
		config.Rioni.Server.Http.AddressStr = ""
		err := Check(config)
		require.EqualError(t, err, "DNS-over-HTTPS server configuration error: listen address must be specified")
	})

	t.Run("http enabled missing cert file", func(t *testing.T) {
		config := config()
		config.Rioni.Server.Http.Tls.CertFileStr = ""
		err := Check(config)
		require.EqualError(t, err, "DNS-over-HTTPS server TLS configuration error: certificate file path must be specified")
	})

	t.Run("http enabled missing key file", func(t *testing.T) {
		config := config()
		config.Rioni.Server.Http.Tls.KeyFileStr = ""
		err := Check(config)
		require.EqualError(t, err, "DNS-over-HTTPS server TLS configuration error: private key file path must be specified")
	})

	t.Run("dns enabled missing address", func(t *testing.T) {
		config := config()
		config.Rioni.Server.Dns.AddressStr = ""
		err := Check(config)
		require.EqualError(t, err, "DNS server configuration error: listen address must be specified")
	})

	t.Run("http disabled no validation", func(t *testing.T) {
		config := config()
		config.Rioni.Server.Http.EnableStr = "false"
		config.Rioni.Server.Http.AddressStr = ""
		err := Check(config)
		require.NoError(t, err)
	})

	t.Run("dns disabled no validation", func(t *testing.T) {
		config := config()
		config.Rioni.Server.Dns.EnableStr = "false"
		config.Rioni.Server.Dns.AddressStr = ""
		err := Check(config)
		require.NoError(t, err)
	})

}

func config() Config {
	return Config{
		Rioni: Rioni{
			Relay: Relay{
				Upstream: []string{"https://dns.google/dns-query"},
				Client: ClientConfig{
					Dns: []string{"8.8.8.8:53"},
				},
			},
			Server: Server{
				Http: HttpConfig{
					EnableStr:  "true",
					AddressStr: ":443",
					Tls: TlsConfig{
						CertFileStr: "/tmp/cert.pem",
						KeyFileStr:  "/tmp/key.pem",
					},
				},
				Dns: DnsConfig{
					EnableStr:  "true",
					AddressStr: ":53",
				},
			},
		},
	}
}
