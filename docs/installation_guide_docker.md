# How to install Rioni with Docker

Using Docker is the recommended way to run Rioni if you want a consistent, isolated environment.
The official image is built on top of a minimal base, keeping the footprint small and secure.

## Prerequisites

Ensure you have [Docker](https://www.docker.com) installed on your system.

## Installation

### Docker Container

Pull the latest image and run Rioni using the following command:

```sh
docker run -d \
  --name rioni \
  --restart unless-stopped \
  -p 443:443/tcp \
  -p 53:53/tcp \
  -p 53:53/udp \
  ghcr.io/andreysenov/rioni:latest
```

This command will start Rioni in a Docker container,
mapping ports 443 and 53 on the host to ports 443 and 53 in the container.
The `-d` flag runs the container in detached mode, allowing you to continue using your terminal while Rioni runs in the background.
The `--restart unless-stopped` flag ensures that the container is restarted automatically if it crashes.

Useful commands:
* `docker logs -f rioni` shows the logs of the container
* `docker stop rioni` stops the container
* `docker rm rioni` removes the container

### Docker Compose

For a more manageable setup, especially if you run other services, create a `docker-compose.yml` file:

```yaml
services:
  rioni:
    image: ghcr.io/andreysenov/rioni:latest
    container_name: rioni
    restart: unless-stopped
    ports:
      - "443:443/tcp"
      - "53:53/udp"
      - "53:53/tcp"
```

Run it with the following command:

```sh
docker compose up -d
```

Use the `-d` flag to run in detached mode, similar to the docker run command.

Useful commands:
* `docker compose stop` stops the container
* `docker compose down` removes the container

## Configuration

The preferred way to configure Rioni in a Docker container is through environment variables.
For instance, let's disable the HTTP server:

```sh
docker run -d \
  --name rioni \
  --restart unless-stopped \
  -p 53:53/tcp \
  -p 53:53/udp \
  -e RIONI_SERVER_HTTP_ENABLE=false \
  ghcr.io/andreysenov/rioni:latest
```

Or with Docker Compose:

```yaml
services:
  rioni:
    image: ghcr.io/andreysenov/rioni:latest
    container_name: rioni
    restart: unless-stopped
    environment:
      RIONI_SERVER_HTTP_ENABLE: false
    ports:
      - "53:53/udp"
      - "53:53/tcp"
```

## See Also

[Configuration](configuration.md)
