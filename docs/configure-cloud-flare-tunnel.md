# Cloudflare Tunnel Setup

## Prerequisites

- `cloudflared` installed on the server machine
- Domain managed by Cloudflare DNS OR
- registered by Cloudflare through your DNS provider
- Env manager server up and running at `http://localhost:8080`

## 1. Authenticate

```bash
cloudflared tunnel login
```

## 2. Create the tunnel

```bash
cloudflared tunnel create <tunnel-name>
```

Note the UUID output - used in all following steps.

## 3. Configure

Create `/etc/cloudflared/config.yml`:

```yaml
tunnel: <TUNNEL-UUID>
credentials-file: /home/<user>/.cloudflared/<TUNNEL-UUID>.json

ingress:
  - hostname: <your-hostname> # (e.g. your.domain.com)
    service: http://localhost:8080
  - service: http_status:404
```

## 4. Route DNS

```bash
cloudflared tunnel route dns <TUNNEL-UUID> <your-hostname>
```

Creates a proxied CNAME record: `<your-hostname>` => `<TUNNEL-UUID>.cfargotunnel.com`

## 5. Install as system service

```bash
sudo cloudflared service install
sudo systemctl enable --now cloudflared
```

## Verify

```bash
curl http://localhost:8080                  # origin reachable
cloudflared tunnel info <tunnel-name>       # tunnel healthy
curl -I https://<your-hostname>             # public endpoint works
```

You should be able to access the env manager server securely at `https://your-domain.com` now!
