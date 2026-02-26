// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package relay

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/AndreySenov/rioni/internal/cfg"
)

func NewTransport(config cfg.Rioni) http.RoundTripper {
	dnsServers := normalizeDnsServers(config.Relay.Client)

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			if len(dnsServers) == 0 {
				return nil, fmt.Errorf("no DNS servers configured")
			}

			type dialResult struct {
				conn net.Conn
				err  error
			}

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			results := make(chan dialResult, len(dnsServers))

			var wg sync.WaitGroup
			wg.Add(len(dnsServers))

			for _, s := range dnsServers {
				dnsServer := s
				go func() {
					defer wg.Done()

					d := net.Dialer{
						Timeout: config.Relay.Client.Timeout(),
					}

					conn, err := d.DialContext(ctx, network, dnsServer)

					select {
					case results <- dialResult{conn: conn, err: err}:
					case <-ctx.Done():
						if conn != nil {
							_ = conn.Close()
						}
					}
				}()
			}

			go func() {
				wg.Wait()
				close(results)
			}()

			var (
				firstConn net.Conn
				lastErr   error
			)

			for result := range results {
				if result.err == nil && firstConn == nil {
					firstConn = result.conn
					cancel()
					continue
				}

				if result.conn != nil {
					_ = result.conn.Close()
				}

				if result.err != nil {
					lastErr = result.err
				}
			}

			if firstConn != nil {
				return firstConn, nil
			}

			return nil, fmt.Errorf("failed to dial any DNS server (tried %d): %w", len(dnsServers), lastErr)
		},
	}

	dialer := &net.Dialer{
		Resolver: resolver,
		Timeout:  config.Relay.Client.Timeout(),
	}

	return &http.Transport{
		DialContext:     dialer.DialContext,
		IdleConnTimeout: config.Server.Http.IdleTimeout(),
	}
}

func normalizeDnsServers(config cfg.ClientConfig) []string {
	const defaultPort = "53"

	dnsServers := make([]string, 0, len(config.Dns))
	for _, dnsServer := range config.Dns {
		if dnsServer == "" {
			continue
		}

		if _, _, err := net.SplitHostPort(dnsServer); err != nil {
			dnsServer = net.JoinHostPort(dnsServer, defaultPort)
		}

		dnsServers = append(dnsServers, dnsServer)
	}

	return dnsServers
}
