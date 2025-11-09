# Deploying OttoApp

## Building for Production

Use the build script to create a production distribution tarball:

```bash
./tools/build.sh
```

This creates:
- `dist/prod/ottoapp-{version}` - Linux AMD64 binary
- `dist/prod/emberjs-{version}/` - Ember production build
- `dist/prod/ottoapp-{version}.tgz` - Deployment tarball

## Deployment Steps

### 1. Upload and extract the tarball

```bash
# Upload the tarball to the remote server
scp dist/prod/ottoapp-{version}.tgz ottomap:/var/www/dev/ottoapp/

# SSH to the server and extract
ssh ottomap
cd /var/www/dev/ottoapp
tar -xzf ottoapp-{version}.tgz
```

This creates:
- `ottoapp-{version}` - the Go binary
- `emberjs-{version}/` - the Ember build directory

### 2. Update symbolic links

```bash
# Update symlinks to point to the new version
ln -sfn ottoapp-{version} ottoapp
ln -sfn emberjs-{version} emberjs
```

The directory structure should look like:
```
/var/www/dev/ottoapp/
├── data/                    # database directory
├── emberjs -> emberjs-0.17.1/
├── emberjs-0.17.1/
├── ottoapp -> ottoapp-0.17.1
└── ottoapp-0.17.1
```

### 3. Restart the service

```bash
systemctl restart ottoapp.service
```

### 4. Create a systemd service

Create a systemd service file at `/etc/systemd/system/ottoapp.service`:

```text
[Unit]
Description=OttoApp dev server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=ottopb
Group=ottopb
WorkingDirectory=/var/www/dev/ottoapp/data
ExecStart=/var/www/dev/ottoapp/ottoapp api serve --db .
Restart=on-failure
RestartSec=13
TimeoutStopSec=30

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ottoapp

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/www/dev/ottoapp/data
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
RestrictRealtime=true
RestrictNamespaces=true
LockPersonality=true

[Install]
WantedBy=multi-user.target
```

### 5. Reload systemd and start the service (first time only)

```bash
systemctl daemon-reload
systemctl enable ottoapp.service
systemctl start ottoapp.service
journalctl --no-page -u ottoapp
```

### 6. Configure Caddy (first time only)

See [CADDY.md](CADDY.md) for production Caddy configuration.

Example production Caddyfile:

```caddy
ottoapp.example.com {
  encode gzip

  # API reverse proxy
  handle /api/* {
    reverse_proxy http://localhost:8181
  }

  # Everything else is static, with SPA treatment for the frontend app
  handle {
    root * /var/www/ottoapp/emberjs
    try_files {path} /index.html
    file_server
  }

  log {
    output file /var/log/caddy/ottoapp.log
    format json
  }
}
```

## Notes

- The systemd service and Caddyfile reference the symlinks (`ottoapp` and `emberjs`), not versioned files
- This allows zero-downtime deployments by updating symlinks and restarting the service
- Because Caddy fronts both Ember and the Go API at one origin, no CORS is needed and cookies flow automatically
- Old versions can be kept for rollback: just update the symlinks and restart
