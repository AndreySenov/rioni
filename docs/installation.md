# Installation

Rioni is distributed as a Docker image and as pre-built binaries for `arm64` and `amd64` architectures,
supporting Linux, macOS, FreeBSD, and Windows.

* Official Docker images are available on [GitHub Container Registry](https://github.com/AndreySenov/rioni/pkgs/container/rioni).
* Official releases can be downloaded from [GitHub Releases](https://github.com/AndreySenov/rioni/releases).

## Docker

The fastest way to deploy Rioni in an isolated environment is using [Docker](https://www.docker.com).
This method provides a pre-configured environment, making it easy to manage Rioni as a service with standard Docker tools.

To get started quickly, run:

```sh
docker run -d \
  --name rioni \
  --restart always \
  -p 443:443/tcp \
  -p 53:53/tcp \
  -p 53:53/udp \
  ghcr.io/andreysenov/rioni:latest
```

For detailed step-by-step instructions and configuration examples,
please refer to the guide [How to install Rioni with Docker](installation_guide_docker.md).

## Homebrew

The easiest way to install and manage Rioni on Linux and macOS is through [Homebrew](https://brew.sh).
This method automatically handles updates and configures Rioni to run as a system service.

To get started quickly, run:

```sh
brew tap andreysenov/tap
brew install rioni
sudo brew services start rioni
```

For detailed step-by-step instructions,
please refer to the guide [How to install Rioni on Linux and macOS with Homebrew](installation_guide_homebrew_linux_macos.md).

## Linux Packages (DEB/RPM)

For a native experience on Debian-based or Red Hat-based distributions,
you can install Rioni using standard package managers.

To get started quickly, run:

```sh
curl -sSfL https://raw.githubusercontent.com/AndreySenov/rioni/main/install/linux_package.sh | sh
sudo systemctl enable --now rioni
```

For detailed step-by-step instructions,
please refer to the guide [How to install Rioni on Linux from DEB and RPM packages](installation_guide_linux_package.md).

## Generic Archives

Generic distribution archives are available for Linux, FreeBSD, macOS, and Windows.
Download the archive for your system and architecture, then unpack it
from the directory where the file was downloaded into the directory of your choice.


### Linux, FreeBSD, macOS

Replace `*` with the actual name of the archive you downloaded and
`<path>` with the path to the directory where you want to install Rioni
before running the commands below.

```sh
tar -xzf rioni-*.tar.gz -C <path>
cd <path>
chmod +x ./run.sh
./run.sh
```

### Windows

Unzip the archive to the directory of your choice and execute `run.bat` from there.
