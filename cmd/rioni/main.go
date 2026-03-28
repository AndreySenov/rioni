// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
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

	config, err := readConfig(configFlag)
	if err != nil {
		log.Fatal(err)
	}

	serverConfig := config.Rioni.Server
	anyServerEnabled := serverConfig.Http.IsEnabled() || serverConfig.Dns.IsEnabled()
	if !anyServerEnabled {
		log.Fatal("No servers enabled; exiting")
	}

	err = checkConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := cert.Check(
		serverConfig.Http.Tls.CertFile(),
		serverConfig.Http.Tls.KeyFile(),
		serverConfig.Http.Tls.IsBuildSelfSigned()); err != nil {
		log.Fatal(err)
	}

	dnsRelay := relay.NewRelay(config.Rioni)
	dohServer := doh.NewServer(serverConfig.Http, dnsRelay)
	dnsServer := dns.NewServer(serverConfig.Dns, dnsRelay)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupErr := make(chan error, 2)

	if serverConfig.Http.IsEnabled() {
		go func() {
			log.Printf("Starting DNS-over-HTTPS server on %s", serverConfig.Http.Address())
			if err := dohServer.Start(ctx); err != nil {
				startupErr <- err
			}
		}()
	}

	if serverConfig.Dns.IsEnabled() {
		go func() {
			log.Printf("Starting DNS server on %s", serverConfig.Dns.Address())
			if err := dnsServer.Start(ctx); err != nil {
				startupErr <- err
			}
		}()
	}

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

func readConfig(configFlag *flag.Flag) (cfg.Config, error) {
	path, ok := os.LookupEnv(cfg.EnvRioniConfigFile)
	if ok && path == "" {
		log.Printf("environment variable %q is set but empty", cfg.EnvRioniConfigFile)
	}

	if path == "" {
		if configFlag != nil {
			path = configFlag.Value.String()
		}
	}

	if path == "" {
		return cfg.ReadEnv()
	}

	config, err := cfg.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("config file %q does not exist", path)
			return cfg.ReadEnv()
		}
		return cfg.Config{}, err
	}

	return cfg.ReadEnvToTarget(config)
}

func checkConfig(config cfg.Config) error {
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
