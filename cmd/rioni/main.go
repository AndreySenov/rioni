// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/AndreySenov/rioni/internal/cert"
	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/AndreySenov/rioni/internal/dns"
	"github.com/AndreySenov/rioni/internal/doh"
	"github.com/AndreySenov/rioni/internal/relay"
)

func main() {
	log.Printf("Rioni %s", readVersion())

	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.String("config", "", "path to config file")
	flag.Parse()

	if *versionFlag {
		return
	}

	configFlag := flag.Lookup("config")

	config := loadConfig(configFlag).Rioni

	if err := cert.Check(
		config.Server.Http.Tls.CertFile(),
		config.Server.Http.Tls.KeyFile(),
		config.Server.Http.Tls.BuildSelfSigned); err != nil {
		log.Fatal(err)
	}

	dnsRelay := relay.NewRelay(config)
	dohServer := doh.NewServer(config.Server.Http, dnsRelay)
	dnsServer := dns.NewServer(config.Server.Dns, dnsRelay)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupErr := make(chan error, 2)

	go func() {
		log.Printf("Starting DNS-over-HTTPS server on %s", config.Server.Http.Address())
		if err := dohServer.Start(ctx); err != nil {
			startupErr <- err
		}
	}()

	go func() {
		log.Printf("Starting DNS server on %s", config.Server.Dns.Address())
		if err := dnsServer.Start(ctx); err != nil {
			startupErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown requested")
	case err := <-startupErr:
		log.Printf("startup error: %v", err)
	}

	_ = dohServer.Shutdown(ctx)
	_ = dnsServer.Shutdown(ctx)
	log.Println("Shut down")
}

func readVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			return info.Main.Version
		}
	}
	return "unknown"
}

func loadConfig(configFlag *flag.Flag) cfg.Config {
	configPath, err := cfg.ResolvePath(configFlag)
	if err != nil {
		log.Fatalf("Failed to resolve config path: %v", err)
	}

	log.Printf("Reading config from %s", configPath)
	config, err := cfg.Read(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("%s\n%+v", "Configuration loaded", config)
	} else {
		log.Printf("%s\n%s", "Configuration loaded", string(b))
	}

	return config
}
