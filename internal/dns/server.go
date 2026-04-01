// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package dns

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/AndreySenov/rioni/internal/relay"
	"github.com/miekg/dns"
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type server struct {
	relay        relay.Relay
	writeTimeout time.Duration
	udp          *dns.Server
	tcp          *dns.Server
	waitGroup    sync.WaitGroup
	shutdownOnce sync.Once
}

func NewServer(config cfg.Dns, relay relay.Relay) Server {
	server := &server{
		relay:        relay,
		writeTimeout: config.WriteTimeout(),
		waitGroup:    sync.WaitGroup{},
		shutdownOnce: sync.Once{},
	}

	server.udp = &dns.Server{
		Net:          "udp",
		Addr:         config.Address(),
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
		Handler:      dns.HandlerFunc(server.handle),
	}

	server.tcp = &dns.Server{
		Net:          "tcp",
		Addr:         config.Address(),
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
		Handler:      dns.HandlerFunc(server.handle),
	}

	return server
}

func (s *server) Start(ctx context.Context) error {
	if s == nil || s.udp == nil || s.tcp == nil {
		return errors.New("dns server is not initialized")
	}

	udpErr := make(chan error, 1)
	tcpErr := make(chan error, 1)
	s.waitGroup.Add(2)

	go func() {
		defer s.waitGroup.Done()
		udpErr <- s.udp.ListenAndServe()
	}()

	go func() {
		defer s.waitGroup.Done()
		tcpErr <- s.tcp.ListenAndServe()
	}()

	select {
	case err := <-udpErr:
		_ = s.Shutdown(ctx)
		return err
	case err := <-tcpErr:
		_ = s.Shutdown(ctx)
		return err
	case <-ctx.Done():
		_ = s.Shutdown(context.Background())
		return ctx.Err()
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	if s == nil || s.udp == nil || s.tcp == nil {
		return errors.New("dns server is not initialized")
	}

	var err error

	s.shutdownOnce.Do(func() {
		done := make(chan struct{})

		go func() {
			_ = s.udp.ShutdownContext(ctx)
			_ = s.tcp.ShutdownContext(ctx)
			s.waitGroup.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-done:
		}
	})

	return err
}

func (s *server) handle(w dns.ResponseWriter, r *dns.Msg) {
	if s == nil || s.relay == nil {
		slog.Error("dns server is not initialized")
		dns.HandleFailed(w, r)
		return
	}

	query, err := r.Pack()
	if err != nil {
		slog.Error("failed to pack DNS message", "error", err)
		dns.HandleFailed(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.writeTimeout)
	defer cancel()

	responseBytes, err := s.relay.Exchange(ctx, query)
	if err != nil {
		slog.Error("failed to query upstream", "error", err)
		dns.HandleFailed(w, r)
		return
	}

	response := new(dns.Msg)
	if err := response.Unpack(responseBytes); err != nil {
		slog.Error("failed to unpack DNS response", "error", err)
		dns.HandleFailed(w, r)
		return
	}

	if response.Id != r.Id {
		slog.Error("upstream answer id mismatch")
		dns.HandleFailed(w, r)
		return
	}

	if !response.Response {
		slog.Error("upstream returned query instead of answer")
		dns.HandleFailed(w, r)
		return
	}

	if err := w.WriteMsg(response); err != nil {
		slog.Error("failed to write DNS response", "error", err)
		return
	}
}
