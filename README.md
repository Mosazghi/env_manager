# env_manager

Are you tired of juggling multiple `.env` files across different projects and devices? Do you wish there was a secure, centralized way to manage your environment variables without relying on third-party services?

The env_manager is a self-hosted environment variable manager that allows you to securely store, manage, and access your environment variables from anywhere.

> [!WARNING]
> This project is under active development and not yet ready for production use.

---

## Installation

### 1. Run the Server

Download the latest binary for your platform from the [releases page](https://github.com/Mosazghi/env_manager/releases).

### 2. Enable HTTPS (recommended)

See [configure_caddy.md](docs/configure_caddy.md) for instructions on setting up Caddy to enable HTTPS with automatic TLS certificates. This is highly recommended if you plan to access the env manager remotely.

## Usage

### Env Manager Server

Firstly the env manager server needs to be running. You can start it with:

```bash
./envm-server service install
```

This will setup the server as a system service and start it immediately. The server listens on port `8080` by default.

If it doesn't start, you can check the logs with:

```bash
./envm-server service start OR ./envm-server service restart
```

TODO
