# Configuration

## Overview

Rioni uses environment variables and a YAML configuration file to define network parameters,
upstream DNS servers, and other settings.

## Configuration Hierarchy

Rioni resolves settings using the following order of precedence (from highest to lowest):

1. Environment Variables

2. YAML Configuration File

3. Default Values

## Configuration File Location

The default configuration file locations are:

- `/etc/rioni/configs/rioni.cfg.yml` (standard Linux package installation)
- `/home/linuxbrew/.linuxbrew/etc/rioni/configs/rioni.cfg.yml` (Linux Homebrew package installation)
- `/opt/homebrew/etc/rioni/configs/rioni.cfg.yml` (macOS Apple Silicon Homebrew package installation)
- `/usr/local/etc/rioni/configs/rioni.cfg.yml` (macOS Intel Homebrew package installation)
- `./configs/rioni.cfg.yml` (generic archive installation)

You can specify an arbitrary configuration file path using the `--config` CLI flag:

```sh
rioni --config /path/to/rioni/config.yml
```

Another way to specify the configuration file path is by setting the `RIONI_CONFIG_FILE` environment variable:

```sh
export RIONI_CONFIG_FILE=/path/to/rioni/config.yml
rioni
```

If you specify both the `--config` flag and the `RIONI_CONFIG_FILE` environment variable,
the environment variable will take precedence.

If the configuration file is missing or some properties are omitted, Rioni falls back to environment variables and default settings.
Should a property exist in both, the environment variable takes precedence.

The configuration file supports the `${VARIABLE:-fallback}` syntax,
allowing you to inject any environment variable into your settings.

## Configuration File Properties

- `rioni` - Top-level configuration section.

- `rioni.relay` - Settings for forwarding DNS requests to upstream DoH resolvers.

- `rioni.relay.upstream` - List of upstream DoH endpoint URLs.
    - Type: array of strings
    - Environment variable: `RIONI_RELAY_UPSTREAM`
    - Default: `[]` (empty array)
    - Example values:
        - `https://dns.google/dns-query`
        - `https://cloudflare-dns.com/dns-query`

- `rioni.relay.client` - Settings for outbound requests sent to upstream resolvers.

- `rioni.relay.client.timeout` - Overall timeout for upstream requests, including response processing.
    - Type: duration string
    - Environment variable: `RIONI_RELAY_CLIENT_TIMEOUT`
    - Default: `5s`
    - Example: `5s`

- `rioni.relay.client.read-limit` - Maximum size of an upstream response that Rioni will read.
    - Type: size string or integer
    - Environment variable: `RIONI_RELAY_CLIENT_READ_LIMIT`
    - Default: `1mb`
    - Examples: `64kb`, `1mb`, `1048576`

- `rioni.relay.client.dns` - List of DNS server addresses used to resolve upstream hostnames.
    - Type: array of strings
    - Environment variable: `RIONI_RELAY_CLIENT_DNS`
    - Default: `[]` (empty array, uses system DNS)
    - Example values:
        - `"8.8.8.8"`
        - `"1.1.1.1"`

- `rioni.server` - Settings for the servers exposed by Rioni.

- `rioni.server.http` - Configuration for the HTTPS server.

- `rioni.server.http.enable` - Enables the HTTP server.
    - Type: boolean
    - Environment variable: `RIONI_SERVER_HTTP_ENABLE`
    - Default: `true`
    - Example: `true`

- `rioni.server.http.address` - Address for the HTTP server to listen on.
    - Type: string
    - Environment variable: `RIONI_SERVER_HTTP_ADDRESS`
    - Default: `:443`
    - Examples: `:443`, `0.0.0.0:443`, `127.0.0.1:8443`

- `rioni.server.http.read-header-timeout` - Maximum time allowed to read HTTP request headers.
    - Type: duration string
    - Environment variable: `RIONI_SERVER_HTTP_READ_HEADER_TIMEOUT`
    - Default: `5s`
    - Example: `5s`

- `rioni.server.http.read-timeout` - Maximum time allowed to read the full HTTP request.
    - Type: duration string
    - Environment variable: `RIONI_SERVER_HTTP_READ_TIMEOUT`
    - Default: `10s`
    - Example: `10s`

- `rioni.server.http.read-limit` - Maximum HTTP request size that Rioni will read.
    - Type: size string or integer
    - Environment variable: `RIONI_SERVER_HTTP_READ_LIMIT`
    - Default: `64kb`
    - Examples: `64kb`, `1mb`

- `rioni.server.http.write-timeout` - Maximum time allowed to write the HTTP response.
    - Type: duration string
    - Environment variable: `RIONI_SERVER_HTTP_WRITE_TIMEOUT`
    - Default: `10s`
    - Example: `10s`

- `rioni.server.http.idle-timeout` - Maximum time an idle HTTP connection is kept open while waiting for the next
  request.
    - Type: duration string
    - Environment variable: `RIONI_SERVER_HTTP_IDLE_TIMEOUT`
    - Default: `30s`
    - Example: `30s`

- `rioni.server.http.tls` - TLS configuration for the HTTPS server.

- `rioni.server.http.tls.cert-file` - Path to the TLS certificate file.
    - Type: string
    - Environment variable: `RIONI_SERVER_HTTP_TLS_CERT_FILE`
    - Default: `""` (empty string)
    - Example: `tls/rioni.crt`

- `rioni.server.http.tls.key-file` - Path to the TLS private key file.
    - Type: string
    - Environment variable: `RIONI_SERVER_HTTP_TLS_KEY_FILE`
    - Default: `""` (empty string)
    - Example: `tls/rioni.key`

- `rioni.server.http.tls.build-self-signed` - If set to `true`, Rioni generates a self-signed certificate when the
  configured certificate and key files do not exist.
    - Type: boolean
    - Environment variable: `RIONI_SERVER_HTTP_TLS_BUILD_SELF_SIGNED`
    - Default: `false`
    - Example: `true`

- `rioni.server.dns` - Configuration for the DNS server.

- `rioni.server.dns.enable` - Enables the DNS server.
    - Type: boolean
    - Environment variable: `RIONI_SERVER_DNS_ENABLE`
    - Default: `true`
    - Example: `true`

- `rioni.server.dns.address` - Address for the DNS server to listen on.
    - Type: string
    - Environment variable: `RIONI_SERVER_DNS_ADDRESS`
    - Default: `:53`
    - Examples: `:53`, `0.0.0.0:53`, `127.0.0.1:5353`

- `rioni.server.dns.read-timeout` - Maximum time allowed to read a DNS request.
    - Type: duration string
    - Environment variable: `RIONI_SERVER_DNS_READ_TIMEOUT`
    - Default: `2s`
    - Example: `2s`

- `rioni.server.dns.write-timeout` - Maximum time allowed to write a DNS response.
    - Type: duration string
    - Environment variable: `RIONI_SERVER_DNS_WRITE_TIMEOUT`
    - Default: `2s`
    - Example: `2s`
