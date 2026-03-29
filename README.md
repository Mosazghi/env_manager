# env_manager

Are you tired of juggling multiple `.env` files across different projects and devices? Do you wish there was a secure, centralized way to manage your environment variables without relying on third-party services? 

The env_manager is a self-hosted environment variable manager that allows you to securely store, manage, and access your environment variables from anywhere. 

It's meant to be self-hosted on your own server (e.g. on a RPi) and accessed via its CLI app. 


> [!WARNING]
This project is under active development and not yet ready for production use.

---

# Installation

## 1. Run the Server

Download the latest binary for your platform from the [releases page](#) and run it:

```bash
./envm-server
```

The server listens on `http://127.0.0.1:8080` by default.

---

## 2. Enable HTTPS (recommended)

> HTTPS is handled by [Caddy](https://caddyserver.com/), a reverse proxy that automatically obtains and renews TLS certificates via Let's Encrypt.

### Install Caddy

Naigate to the [Caddy download page](https://caddyserver.com/download) and follow the instructions for your platform to install Caddy.

---

### Configure Caddy

**Linux/macOS** — edit `/etc/caddy/Caddyfile`:
```bash
sudo nano /etc/caddy/Caddyfile
```

**Windows** — edit `C:\ProgramData\Caddy\Caddyfile`

Paste this, replacing with your domain:

```
your-domain.com {
    reverse_proxy 127.0.0.1:8080
}
```



---

### Start Caddy


### Prerequisites

Before starting Caddy, make sure:

- Your domain's DNS A record points to this machine's public IP
- Ports **80** and **443** are open and forwarded to this machine on your router
- `envm-server` is already running on port `8080`

--- 

**Linux :**
```bash
sudo systemctl enable --now caddy
```

**macOS:**
```bash
brew services start caddy
```

**Windows:**
```powershell
caddy start
```

Once Caddy starts it will automatically obtain a TLS certificate, no manual certificate management needed.


You should be able to access the env manager securely at `https://your-domain.com` now!
