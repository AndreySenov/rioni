// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cfg

import "errors"

func Check(config Config) error {
	relayConfig := config.Rioni.Relay
	if len(relayConfig.Upstream) == 0 {
		return errors.New("relay configuration error: at least one upstream DNS-over-HTTPS server must be configured")
	}
	if len(relayConfig.Client.Dns) == 0 {
		return errors.New("relay client configuration error: at least one DNS server must be configured to resolve upstream domains")
	}

	httpConfig := config.Rioni.Server.Http
	if httpConfig.IsEnabled() {
		if httpConfig.Address() == "" {
			return errors.New("DNS-over-HTTPS server configuration error: listen address must be specified")
		}
		if httpConfig.Tls.CertFile() == "" {
			return errors.New("DNS-over-HTTPS server TLS configuration error: certificate file path must be specified")
		}
		if httpConfig.Tls.KeyFile() == "" {
			return errors.New("DNS-over-HTTPS server TLS configuration error: private key file path must be specified")
		}
	}

	dnsConfig := config.Rioni.Server.Dns
	if dnsConfig.IsEnabled() {
		if dnsConfig.Address() == "" {
			return errors.New("DNS server configuration error: listen address must be specified")
		}
	}

	return nil
}
