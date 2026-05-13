// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package doh

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/AndreySenov/rioni/internal/httpx"
	"github.com/AndreySenov/rioni/internal/logx"
	"github.com/AndreySenov/rioni/internal/relay"
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type server struct {
	logger       *slog.Logger
	readLimit    int64
	certFile     string
	keyFile      string
	relay        relay.Relay
	http         *http.Server
	waitGroup    sync.WaitGroup
	shutdownOnce sync.Once
}

const (
	componentName = "doh-server"
	uri           = "/dns-query"
)

func NewServer(baseLogger *slog.Logger, config cfg.Http, relay relay.Relay) Server {
	server := &server{
		logger:    baseLogger.With(logx.ComponentKey, componentName),
		readLimit: config.ReadLimit(),
		certFile:  config.Tls.CertFile(),
		keyFile:   config.Tls.KeyFile(),
		relay:     relay,
		waitGroup: sync.WaitGroup{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(uri, server.handle)

	server.http = &http.Server{
		Addr:    config.Address(),
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
		ReadHeaderTimeout: config.ReadHeaderTimeout(),
		ReadTimeout:       config.ReadTimeout(),
		WriteTimeout:      config.WriteTimeout(),
		IdleTimeout:       config.IdleTimeout(),
	}

	return server
}

func (s *server) Start(ctx context.Context) error {
	if s == nil || s.http == nil {
		return errors.New("http server is not initialized")
	}

	listen, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return err
	}

	httpErr := make(chan error, 1)
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		httpErr <- s.http.ServeTLS(listen, s.certFile, s.keyFile)
	}()

	select {
	case err := <-httpErr:
		_ = s.Shutdown(ctx)
		return err
	case <-ctx.Done():
		_ = s.Shutdown(context.Background())
		return ctx.Err()
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	var err error

	s.shutdownOnce.Do(func() {
		done := make(chan struct{})

		go func() {
			_ = s.http.Shutdown(ctx)
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

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	if s == nil {
		slog.Error("http server is not initialized")
		httpx.ErrorInternalServerError(w)
		return
	}

	if s.relay == nil {
		s.logger.Error("DNS relay is not initialized")
		httpx.ErrorInternalServerError(w)
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodGet:
		s.handleGet(w, r)
	default:
		logx.WarnForHttpRequest(s.logger, r, "unsupported method")
		w.Header().Set(httpx.HeaderAllow, http.MethodGet+", "+http.MethodPost)
		httpx.ErrorMethodNotAllowed(w)
	}
}

func (s *server) handlePost(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	if accept := r.Header.Get(httpx.HeaderAccept); accept != httpx.ContentTypeApplicationDnsMessage {
		logx.WarnForHttpRequest(s.logger, r, "unsupported accept header",
			"accept", accept,
		)
		httpx.ErrorNotAcceptable(w)
		return
	}

	if contentType := r.Header.Get(httpx.HeaderContentType); contentType != httpx.ContentTypeApplicationDnsMessage {
		logx.WarnForHttpRequest(s.logger, r, "unsupported content type",
			"content_type", contentType,
		)
		httpx.ErrorUnsupportedMediaType(w)
		return
	}

	contentLength := r.ContentLength

	if contentLength == 0 {
		logx.WarnForHttpRequest(s.logger, r, "request body is empty")
		httpx.ErrorLengthRequired(w)
		return
	}

	if contentLength > s.readLimit {
		logx.WarnForHttpRequest(s.logger, r, "request body is too large",
			"content_length", contentLength,
			"read_limit", s.readLimit,
		)
		httpx.ErrorRequestEntityTooLarge(w)
		return
	}

	if contentLength == -1 {
		contentLength = s.readLimit
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, contentLength+1))
	if err != nil {
		logx.WarnForHttpRequest(s.logger, r, "failed to read request body",
			"error", err,
		)
		httpx.ErrorBadRequest(w)
		return
	}

	if int64(len(body)) > contentLength {
		logx.WarnForHttpRequest(s.logger, r, "request body is too large",
			"content_length", contentLength,
			"read_limit", s.readLimit,
		)
		httpx.ErrorRequestEntityTooLarge(w)
		return
	}

	s.exchange(body, w, r)
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	dnsParam := r.URL.Query().Get("dns")
	if dnsParam == "" {
		logx.WarnForHttpRequest(s.logger, r, "missing dns query parameter")
		httpx.ErrorBadRequest(w)
		return
	}

	dnsMessage, err := base64.RawURLEncoding.DecodeString(dnsParam)
	if err != nil {
		logx.WarnForHttpRequest(s.logger, r, "failed to decode DNS message",
			"error", err,
		)
		httpx.ErrorBadRequest(w)
		return
	}

	s.exchange(dnsMessage, w, r)
}

func (s *server) exchange(dnsMessage []byte, w http.ResponseWriter, r *http.Request) {
	resp, err := s.relay.Exchange(r.Context(), dnsMessage)
	if err != nil {
		logx.ErrorForHttpRequest(s.logger, r, "failed to relay DNS message",
			"error", err,
		)
		httpx.ErrorBadGateway(w)
		return
	}

	w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeApplicationDnsMessage)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logx.ErrorForHttpRequest(s.logger, r, "failed to write response",
			"error", err,
		)
	}
}
