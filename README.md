[![Go Tests](https://github.com/Mosazghi/env_manager/actions/workflows/test.yml/badge.svg)](https://github.com/Mosazghi/env_manager/actions/workflows/test.yml) [![Latest Release](https://img.shields.io/github/v/release/Mosazghi/env_manager?display_name=tag)](https://github.com/Mosazghi/env_manager/releases/latest)

# env_manager

Are you tired of juggling multiple `.env` files across different projects and devices? Do you wish there was a secure, centralized way to manage your environment variables without relying on third-party services?

The env_manager is a self-hosted environment variable manager that allows you to securely store, manage, and access your environment variables from anywhere.

> [!WARNING]
> This project is under active development and not yet ready for production use.

---

## Installation

### 1. Download the Env Manager server

Download the latest binary for your platform from the [releases page](https://github.com/Mosazghi/env_manager/releases).

### 2. Enable HTTPS (recommended)

Enabling HTTPS is highly recommended if you plan to access the env manager remotely outside your own
network.

There's two alternatives for achieving this:

1. Use [Cloudflare Tunnel](/docs/configure-cloud-flare-tunnel.md) or
2. [by configuring Caddy](docs/configure-caddy.md)

Both options are easy to setup and most importantly _FREE_ to use.

See their documentation files a step-by-step guide.

## Usage

### Env Manager Server

Firstly the env manager server needs to be installed as _system service_ (i.e., it will be registered
as a background app):

```bash
./envms service install
```

The server listens on port `8080` by default.

If it doesn't start by default, do:

```bash
./envms service start
```

## Env Client

TODO!
