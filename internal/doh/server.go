// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package doh

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/AndreySenov/rioni/internal/httpx"
	"github.com/AndreySenov/rioni/internal/relay"
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type server struct {
	readLimit    int64
	certFile     string
	keyFile      string
	relay        relay.Relay
	http         *http.Server
	waitGroup    sync.WaitGroup
	shutdownOnce sync.Once
}

const uri = "/dns-query"

func NewServer(config cfg.Http, relay relay.Relay) Server {
	server := &server{
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

	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		if err := s.http.ServeTLS(listen, s.certFile, s.keyFile); err != nil {
			log.Println(err)
		}
	}()

	return nil
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
	if s == nil || s.relay == nil {
		log.Println("http server is not initialized")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodGet:
		s.handleGet(w, r)
	default:
		w.Header().Set(httpx.HeaderAllow, http.MethodGet+", "+http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) handlePost(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	if contentType := r.Header.Get(httpx.HeaderContentType); contentType != httpx.ContentTypeApplicationDnsMessage {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	req, err := io.ReadAll(http.MaxBytesReader(w, r.Body, s.readLimit+1))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if int64(len(req)) > s.readLimit {
		http.Error(w, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	s.exchange(req, w, r)
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	dnsParam := r.URL.Query().Get("dns")
	if dnsParam == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	dnsMessage, err := base64.RawURLEncoding.DecodeString(dnsParam)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	s.exchange(dnsMessage, w, r)
}

func (s *server) exchange(dnsMessage []byte, w http.ResponseWriter, r *http.Request) {
	resp, err := s.relay.Exchange(r.Context(), dnsMessage)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeApplicationDnsMessage)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
