# How to install Rioni on Linux and macOS with Homebrew

This guide provides a step-by-step process for installing and running Rioni as a background service.

## Prerequisites

Ensure you have [Homebrew](https://brew.sh) installed on your Linux or macOS.

## Installation

### 1. Add the Tap and Install

Run the following commands to add the repository and install the Rioni formula:

```sh
$ brew tap andreysenov/tap
$ brew install rioni
```

### 2. Start the Service

Since Rioni uses privileged ports (53 and 443) by default,
you may need to start it with root privileges to allow it to bind to these ports:

```sh
$ sudo brew services start rioni
```

_Note: If you configure Rioni to use ports above 1024, you can run it without `sudo`._

### 3. Verify Installation

Check if the service is active and running:

```sh
sudo brew services list | grep rioni
```

_Note: If you start the service without `sudo`, you must also check its status without `sudo`._

### 4. Test Connection

Verify that Rioni is correctly processing requests using `dig`:

#### Classic DNS (Port 53):

```sh
$ dig @127.0.0.1 -p 53 example.com
```

#### DoH POST (Port 443):

```sh
$ dig @127.0.0.1 -p 443 +https example.com
```

#### DoH GET (Port 443):

```sh
$ dig @127.0.0.1 -p 443 +https-get example.com
```

### Locate the Configuration File

The default path for the configuration file is:

| Operating System    | Path to Configuration File                                 |
|---------------------|------------------------------------------------------------|
| Linux               | /home/linuxbrew/.linuxbrew/etc/rioni/configs/rioni.cfg.yml |
| macOS Apple Silicon | /opt/homebrew/etc/rioni/configs/rioni.cfg.yml              |
| macOS Intel         | /usr/local/etc/rioni/configs/rioni.cfg.yml                 |

Edit the configuration file with your favorite text editor, e.g.:

```sh
$ sudo nano $(brew --prefix)/etc/rioni/configs/rioni.cfg.yml
```

Any changes to the configuration require a service restart to take effect:

```sh
$ sudo brew services restart rioni
```

## See Also

[Configuration](configuration.md)
