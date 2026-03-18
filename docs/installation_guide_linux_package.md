# How to install Rioni on Linux from DEB and RPM packages

The fastest way to install Rioni from Linux packages is using the automated installation script.
It detects your system architecture and OS to install the correct `.deb` or `.rpm` package.

```sh
curl -sSfL https://raw.githubusercontent.com/AndreySenov/rioni/main/install/linux_package.sh | sh
```

### Manual Installation

If you prefer to control the installation process manually, follow these steps:

#### Download the Package

Go to the [GitHub Releases](https://github.com/AndreySenov/rioni/releases) page and download the version for your system:

* `.deb` for Ubuntu, Debian, Mint, and other Debian-based Linux distributions
* `.rpm` for Fedora, RHEL, CentOS, and other Red Hat-based Linux distributions

#### Install the Package

Replace `*` with the downloaded package name.

```sh
# For Ubuntu/Debian/Mint
sudo dpkg -i rioni-*.deb

# For Fedora/RHEL/CentOS
sudo dnf install rioni-*.rpm
```

### Start the Service

Rioni is registered as a systemd service during installation.

```sh
sudo systemctl enable --now rioni
```

### Verify the installation

```sh
systemctl status rioni
```

### Test Connection

Verify that Rioni is correctly processing requests using `dig`:

#### Classic DNS (Port 53):

```sh
dig @127.0.0.1 -p 53 example.com
```

#### DoH POST (Port 443):

```sh
dig @127.0.0.1 -p 443 +https example.com
```

#### DoH GET (Port 443):

```sh
dig @127.0.0.1 -p 443 +https-get example.com
```

### Locate the Configuration File

The default path for the configuration file is:

`/etc/rioni/configs/rioni.cfg.yml`

Edit the configuration file with your favorite text editor, e.g.:

```sh
sudo nano /etc/rioni/configs/rioni.cfg.yml
```

Any changes to the configuration require a service restart to take effect:

```sh
sudo systemctl restart rioni
```

## See Also

[Configuration](configuration.md)
