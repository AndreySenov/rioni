// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package relay

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/miekg/dns"
)

type Relay interface {
	Exchange(ctx context.Context, query *dns.Msg) (*dns.Msg, error)
}

type relay struct {
	client   Client
	upstream []string
}

type queryResult struct {
	endpoint string
	response *dns.Msg
	err      error
}

func NewRelay(config cfg.Rioni) Relay {
	return &relay{
		client:   NewClient(config),
		upstream: append([]string(nil), config.Relay.Upstream...),
	}
}

func (r *relay) Exchange(ctx context.Context, query *dns.Msg) (*dns.Msg, error) {
	if r == nil {
		return nil, errors.New("relay is not initialized")
	}
	if r.client == nil {
		return nil, errors.New("relay upstream client is not initialized")
	}
	if len(r.upstream) == 0 {
		return nil, errors.New("no upstream configured")
	}
	if query == nil {
		return nil, errors.New("query is nil")
	}
	if query.Response {
		return nil, errors.New("query is a response")
	}

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results, err := r.query(reqCtx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
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

func (r *relay) query(ctx context.Context, query *dns.Msg) (chan queryResult, error) {
	bytes, packErr := query.Pack()
	if packErr != nil {
		return nil, packErr
	}

	results := make(chan queryResult, len(r.upstream))

	for _, u := range r.upstream {
		endpoint := u
		go func() {
			res, err := r.client.Query(ctx, endpoint, bytes)
			if err != nil {
				results <- queryResult{endpoint: endpoint, err: err}
				return
			}

			response := new(dns.Msg)
			err = response.Unpack(res)
			if err != nil {
				results <- queryResult{endpoint: endpoint, err: err}
				return
			}

			if !response.Response {
				results <- queryResult{endpoint: endpoint, err: errors.New("upstream returned query instead of answer")}
				return
			}

			if response.Id != query.Id {
				results <- queryResult{endpoint: endpoint, err: errors.New("upstream answer id mismatch")}
				return
			}

			results <- queryResult{endpoint: endpoint, response: response, err: nil}
		}()
	}

	return results, nil
}
