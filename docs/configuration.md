# Configuration

## Overview

Rioni uses a YAML configuration file to define network parameters, upstream DNS servers, and other settings.

## Configuration File Location

The default configuration file locations are:

- `/etc/rioni/configs/rioni.cfg.yml` (standard Linux package installation)
- `/home/linuxbrew/.linuxbrew/etc/rioni/configs/rioni.cfg.yml` (Linux Homebrew package installation)
- `/opt/homebrew/etc/rioni/configs/rioni.cfg.yml` (macOS Apple Silicon Homebrew package installation)
- `/usr/local/etc/rioni/configs/rioni.cfg.yml` (macOS Intel Homebrew package installation)
- `./configs/rioni.cfg.yml` (generic archive installation)

You can specify a custom configuration file path using the `--config` CLI flag:

```sh
rioni --config /path/to/rioni/config.yml
```

Another way to specify the configuration file is by setting the `RIONI_CONFIG_FILE` environment variable:

```sh
export RIONI_CONFIG_FILE=/path/to/rioni/config.yml
rioni
```

If you specify both the `--config` flag and the `RIONI_CONFIG_FILE` environment variable,
the environment variable will take precedence.

If none of these methods is used, Rioni will fail to start.

## Configuration File Properties

- `rioni` - Top-level configuration section.

- `rioni.relay` - Settings for forwarding DNS requests to upstream DoH resolvers.

- `rioni.relay.upstream` - List of upstream DoH endpoint URLs.
    - Type: array of strings
    - Example values:
        - `https://dns.google/dns-query`
        - `https://cloudflare-dns.com/dns-query`

- `rioni.relay.client` - Settings for outbound requests sent to upstream resolvers.

- `rioni.relay.client.timeout` - Overall timeout for upstream requests, including response processing.
    - Type: duration string
    - Example: `5s`

- `rioni.relay.client.read-limit` - Maximum size of an upstream response that Rioni will read.
    - Type: size string or integer
    - Examples: `64kb`, `1mb`, `1048576`

- `rioni.relay.client.dns` - List of DNS server addresses used to resolve upstream hostnames.
    - Type: array of strings
    - Example values:
        - `"8.8.8.8"`
        - `"1.1.1.1"`

- `rioni.server` - Settings for the servers exposed by Rioni.

- `rioni.server.http` - Configuration for the HTTPS server.

- `rioni.server.http.address` - Address for the HTTP server to listen on.
    - Type: string
    - Examples: `:443`, `0.0.0.0:443`, `127.0.0.1:8443`

- `rioni.server.http.read-header-timeout` - Maximum time allowed to read HTTP request headers.
    - Type: duration string
    - Example: `5s`

- `rioni.server.http.read-timeout` - Maximum time allowed to read the full HTTP request.
    - Type: duration string
    - Example: `10s`

- `rioni.server.http.read-limit` - Maximum HTTP request size that Rioni will read.
    - Type: size string or integer
    - Examples: `64kb`, `1mb`

- `rioni.server.http.write-timeout` - Maximum time allowed to write the HTTP response.
    - Type: duration string
    - Example: `10s`

- `rioni.server.http.idle-timeout` - Maximum time an idle HTTP connection is kept open while waiting for the next request.
    - Type: duration string
    - Example: `30s`

- `rioni.server.http.tls` - TLS configuration for the HTTPS server.

- `rioni.server.http.tls.cert-file` - Path to the TLS certificate file.
    - Type: string
    - Example: `tls/rioni.crt`

- `rioni.server.http.tls.key-file` - Path to the TLS private key file.
    - Type: string
    - Example: `tls/rioni.key`

- `rioni.server.http.tls.build-self-signed` - If set to `true`, Rioni generates a self-signed certificate when the configured certificate and key files do not exist.
    - Type: boolean
    - Example: `true`

- `rioni.server.dns` - Configuration for the DNS server.

- `rioni.server.dns.address` - Address for the DNS server to listen on.
    - Type: string
    - Examples: `:53`, `0.0.0.0:53`, `127.0.0.1:5353`

- `rioni.server.dns.read-timeout` - Maximum time allowed to read a DNS request.
    - Type: duration string
    - Example: `2s`

- `rioni.server.dns.write-timeout` - Maximum time allowed to write a DNS response.
    - Type: duration string
    - Example: `2s`
