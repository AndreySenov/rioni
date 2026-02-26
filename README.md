# Rioni

A free and open-source DNS proxy server.

## Overview

Rioni is a lightweight DNS proxy designed to bridge
traditional DNS traffic with modern encrypted protocols.
It acts as a gateway that multiplexes incoming queries from UDP, TCP,
and DoH clients and securely forwards them to upstream DoH resolvers.

## Key Features

- **Downstream:** Handles inbound queries via classic DNS and DoH.
- **Upstream:** Uses DoH for all outbound requests.
- **Performance:** Built with an asynchronous, non-blocking architecture to ensure low latency.
- **No Dependencies:** Compiled as a static binary for effortless deployment across various environments.

## Usage

[Installation](docs/installation.md)
<br>
[Configuration](docs/configuration.md)

## License
Rioni is licensed under the Apache License, Version 2.0. See [NOTICE](NOTICE) and [LICENSE](LICENSE) for details.
