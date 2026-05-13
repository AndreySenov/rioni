// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package relay

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndreySenov/rioni/internal/cfg"
)

type Relay interface {
	Exchange(ctx context.Context, dnsMessage []byte) ([]byte, error)
}

type relay struct {
	client   Client
	upstream []string
}

func NewRelay(config cfg.Rioni) Relay {
	return &relay{
		client:   NewClient(config),
		upstream: append([]string(nil), config.Relay.Upstream...),
	}
}

func (r *relay) Exchange(ctx context.Context, dnsMessage []byte) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("relay is not initialized")
	}
	if r.client == nil {
		return nil, fmt.Errorf("relay upstream client is not initialized")
	}
	if len(dnsMessage) == 0 {
		return nil, fmt.Errorf("dns message is empty")
	}
	if len(r.upstream) == 0 {
		return nil, fmt.Errorf("no upstream configured")
	}

	type result struct {
		endpoint string
		response []byte
		err      error
	}

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan result, len(r.upstream))

	for _, u := range r.upstream {
		endpoint := u
		go func() {
			response, err := r.client.Query(reqCtx, endpoint, dnsMessage)
			results <- result{endpoint: endpoint, response: response, err: err}
		}()
	}

	var errs []error
	for i := 0; i < len(r.upstream); i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("exchange canceled: %w", ctx.Err())
		case r := <-results:
			if r.err == nil {
				cancel()
				return r.response, nil
			}
			errs = append(errs, fmt.Errorf("error querying upstream %s: %w", r.endpoint, r.err))
		}
	}

	return nil, fmt.Errorf("upstream queries failed: %w", errors.Join(errs...))
}
