# Configure Caddy for HTTPS

> HTTPS is handled by [Caddy](https://caddyserver.com/), a reverse proxy that automatically obtains and renews TLS certificates via Let's Encrypt.

## Install Caddy

Naigate to the [Caddy download page](https://caddyserver.com/download) and follow the instructions for your platform to install Caddy.

---

## Configure Caddy

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

## Start Caddy


## Prerequisites

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
