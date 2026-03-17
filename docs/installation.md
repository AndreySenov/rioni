# Installation

Rioni provides pre-built binaries for arm64 and amd64 architectures,
supporting Linux, FreeBSD, macOS, and Windows.

Official releases can be downloaded from the [GitHub Releases](https://github.com/AndreySenov/rioni/releases) page.

| OS      | Package Format            |
|---------|---------------------------|
| Linux   | `.tar.gz`, `.deb`, `.rpm` |
| FreeBSD | `.tar.gz`                 |
| macOS   | `.tar.gz`                 |
| Windows | `.zip`                    |

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

## Linux Packages

Standard installation packages are available for Linux.
Download the package for your system and architecture, then install it
from the directory where the file was downloaded.
After installation, Rioni will be registered as a systemd service.
Replace `*` with the actual name of the package you downloaded before running the commands below.

### Debian-based Linux Distributions
```sh
$ sudo apt install ./rioni-*.deb
```

### RPM-based Linux Distributions
```sh
$ sudo rpm -i ./rioni-*.rpm
```

Run the service using the following command:

```sh
$ sudo systemctl start rioni
```

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
