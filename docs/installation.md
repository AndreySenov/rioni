# Installation

Rioni provides pre-built binaries for arm64 and amd64 architectures,
supporting Linux, macOS, FreeBSD, and Windows.

Official releases can be downloaded from the [GitHub Releases](https://github.com/AndreySenov/rioni/releases) page.

## Homebrew

The easiest way to install and manage Rioni on Linux and macOS is through [Homebrew](https://brew.sh).
This method automatically handles updates and configures Rioni to run as a system service.

To get started quickly, run:

```sh
$ brew tap andreysenov/tap
$ brew install rioni
$ sudo brew services start rioni
```

For detailed step-by-step instructions,
please refer to the guide [How to install Rioni on Linux and macOS with Homebrew](installation_guide_homebrew_linux_macos.md).

## Linux Packages (DEB/RPM)

For a native experience on Debian-based or Red Hat-based distributions,
you can install Rioni using standard package managers.

To get started quickly, run:

```sh
$ curl -sSfL https://raw.githubusercontent.com/AndreySenov/rioni/main/install/linux_package.sh | sh
$ sudo systemctl enable --now rioni
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
$ tar -xzf rioni-*.tar.gz -C <path>
$ cd <path>
$ chmod +x ./run.sh
$ ./run.sh
```

### Windows

Unzip the archive to the directory of your choice and execute `run.bat` from there.
