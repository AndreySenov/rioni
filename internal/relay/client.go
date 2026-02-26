// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package relay

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/AndreySenov/rioni/internal/httpx"
)

type Client interface {
	Query(ctx context.Context, endpoint string, dnsMessage []byte) ([]byte, error)
}

type client struct {
	readLimit int64
	client    *http.Client
}

func NewClient(config cfg.Rioni) Client {
	return &client{
		readLimit: config.Relay.Client.ReadLimit(),
		client: &http.Client{
			Transport: NewTransport(config),
			Timeout:   config.Relay.Client.Timeout(),
		},
	}
}

func (u *client) Query(ctx context.Context, endpoint string, dnsMessage []byte) ([]byte, error) {
	if u == nil || u.client == nil {
		return nil, fmt.Errorf("upstream client is not initialized")
	}
	if len(dnsMessage) == 0 {
		return nil, fmt.Errorf("dns message is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(dnsMessage))
	if err != nil {
		return nil, fmt.Errorf("error creating upstream request: %w", err)
	}

	req.Header.Set(httpx.HeaderContentType, httpx.ContentTypeApplicationDnsMessage)
	req.Header.Set(httpx.HeaderAccept, httpx.ContentTypeApplicationDnsMessage)

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending upstream request: %w", err)
	}
	defer resp.Body.Close()

	if httpx.IsStatusNotOk(resp) {
		return nil, fmt.Errorf("unexpected upstream response status code: %d", resp.StatusCode)
	}

	if contentType := resp.Header.Get(httpx.HeaderContentType); contentType != httpx.ContentTypeApplicationDnsMessage {
		return nil, fmt.Errorf("unexpected upstream response content type: %s", contentType)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, u.readLimit+1))
	if err != nil {
		return nil, fmt.Errorf("error reading upstream response body: %w", err)
	}

	if int64(len(body)) > u.readLimit {
		return nil, fmt.Errorf("upstream response too large (>%d bytes)", u.readLimit)
	}

	return body, nil
}
