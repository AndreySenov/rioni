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
	"github.com/AndreySenov/rioni/internal/logx"
	"github.com/AndreySenov/rioni/internal/relay"
	"github.com/miekg/dns"
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type server struct {
	logger       *slog.Logger
	relay        relay.Relay
	writeTimeout time.Duration
	udp          *dns.Server
	tcp          *dns.Server
	waitGroup    sync.WaitGroup
	shutdownOnce sync.Once
}

const (
	componentName = "dns-server"
)

func NewServer(baseLogger *slog.Logger, config cfg.Dns, relay relay.Relay) Server {
	server := &server{
		logger:       baseLogger.With(logx.ComponentKey, componentName),
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
	if s == nil {
		slog.Error("dns server is not initialized")
		handleFailed(w, r)
		return
	}

	if s.relay == nil {
		logx.ErrorForDnsRequest(s.logger, w, "DNS relay is not initialized")
		handleFailed(w, r)
		return
	}

	if r == nil {
		logx.ErrorForDnsRequest(s.logger, w, "DNS request is nil")
		handleFailed(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.writeTimeout)
	defer cancel()

	response, err := s.relay.Exchange(ctx, r)
	if err != nil {
		logx.ErrorForDnsRequest(s.logger, w, "failed to query upstream",
			"error", err,
		)
		handleFailed(w, r)
		return
	}

	if err = w.WriteMsg(response); err != nil {
		logx.ErrorForDnsRequest(s.logger, w, "failed to write DNS response",
			"error", err,
		)
		return
	}
}

func handleFailed(w dns.ResponseWriter, r *dns.Msg) {
	if w == nil || r == nil {
		return
	}
	m := new(dns.Msg)
	m.SetRcode(r, dns.RcodeServerFailure)
	_ = w.WriteMsg(m)
}
