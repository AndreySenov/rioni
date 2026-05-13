// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/AndreySenov/rioni/internal/cert"
	"github.com/AndreySenov/rioni/internal/cfg"
	"github.com/AndreySenov/rioni/internal/dns"
	"github.com/AndreySenov/rioni/internal/doh"
	"github.com/AndreySenov/rioni/internal/logx"
	"github.com/AndreySenov/rioni/internal/relay"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.String("config", "", "path to config file")
	flag.Parse()

	version := readVersion()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	baseLogger := logx.NewJsonLog("rioni", version)
	slog.SetDefault(baseLogger)

	logger := baseLogger.With(logx.ComponentKey, "main")

	logger.Info("starting Rioni")

	configFlag := flag.Lookup("config")

	config, err := readConfig(logger, configFlag)
	if err != nil {
		logger.Error("failed to read config, shutting down", "error", err)
		os.Exit(1)
	}

	err = cfg.Check(config)
	if err != nil {
		logger.Error("invalid config, shutting down", "error", err)
		os.Exit(1)
	}

	serverConfig := config.Rioni.Server
	anyServerEnabled := serverConfig.Http.IsEnabled() || serverConfig.Dns.IsEnabled()
	if !anyServerEnabled {
		logger.Error("no servers are enabled, shutting down")
		os.Exit(1)
	}

	if err := cert.Check(serverConfig.Http.Tls); err != nil {
		logger.Error("TLS certificate error, shutting down", "error", err)
		os.Exit(1)
	}

	dnsRelay := relay.NewRelay(config.Rioni)
	dohServer := doh.NewServer(baseLogger, serverConfig.Http, dnsRelay)
	dnsServer := dns.NewServer(baseLogger, serverConfig.Dns, dnsRelay)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupErr := make(chan error, 2)

	if serverConfig.Http.IsEnabled() {
		go func() {
			logger.Info("starting DNS-over-HTTPS server", "address", serverConfig.Http.Address())
			if err := dohServer.Start(ctx); err != nil {
				startupErr <- fmt.Errorf("DNS-over-HTTPS server error: %w", err)
			}
		}()
	}

	if serverConfig.Dns.IsEnabled() {
		go func() {
			logger.Info("starting DNS server", "address", serverConfig.Dns.Address())
			if err := dnsServer.Start(ctx); err != nil {
				startupErr <- fmt.Errorf("DNS server error: %w", err)
			}
		}()
	}

	select {
	case <-ctx.Done():
		logger.Info("shutdown requested")
	case err := <-startupErr:
		logger.Error("server failed to start", "error", err)
	}

	_ = dohServer.Shutdown(ctx)
	_ = dnsServer.Shutdown(ctx)
	logger.Info("shutdown complete")
}

func readVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			return info.Main.Version
		}
	}
	return "unknown"
}

func readConfig(logger *slog.Logger, configFlag *flag.Flag) (cfg.Config, error) {
	path, ok := os.LookupEnv(cfg.EnvRioniConfigFile)
	if ok && path == "" {
		logger.Warn("environment variable is set but empty, ignored", "name", cfg.EnvRioniConfigFile)
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
			logger.Warn("config file does not exist", "path", path)
			return cfg.ReadEnv()
		}
		return cfg.Config{}, err
	}

	return cfg.ReadEnvToTarget(config)
}
