# ottoapp
Rewrite of the Ottomap web server for v2

## CORS

Because Caddy fronts both Ember and the Go API at one origin, no CORS is needed and cookies flow automatically.

## Build

Dev & prod usage with Caddy is simple.

Run Ember in development and the Go server in separate terminals.

```bash
# dev
ember serve            # :4200 (Caddy will reverse-proxy it)
go run ./cmd/server    # :8181 (Caddy proxies /api)
# Caddyfile.dev already set to serve https://ottoapp.localhost
```

Or, if you have `air` installed, use it to watch the Go folder and rebuild automatically.

Production build creates files in the `frontend/dist/` directory.

```bash
ember build --environment=production   # outputs to dist/
```

Assuming that we deploy the entire `dist/` folder to `/var/www/ottoapp.mdhenderson.com/dist`, the Caddyfile should be configured to serve `dist/` at `/` and proxy `/api/` to Go.

## Testing

If the server's debug.autoLog flag is set,

```curl
curl https://ottoapp.localhost:8443/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"catbird","password":"secret"}'
```

If the password is `admin`, or `chief`, the session will be created with the same role.
Otherwise, the role will be `guest`.

## Deploying

1. Create a systemd service file at `/etc/systemd/system/ottoapp.service`

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

2. Reload systemd, enable and start the service:

```bash
systemctl daemon-reload
systemctl enable ottoapp.service
systemctl start ottoapp.service
journalctl --no-page -u ottoapp
```

3. Create a Caddyfile at /etc/caddy/Caddyfile

```text
ottomap.playbymailgames.com {
	encode gzip

	# CORS and preflight
	header Access-Control-Allow-Origin "*"
	header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS"
	header Access-Control-Allow-Headers "Content-Type, Authorization"

	@options method OPTIONS
	respond @options 200

	# API reverse proxy
	handle /api/* {
		reverse_proxy http://localhost:8181
	}

	# Everything else is static, with SPA treatment for the frontend app
	handle {
		root * /var/www/dev/ottoapp/emberjs
		try_files {path} /index.html
		file_server
	}

	log {
		output file /var/log/caddy/prd-ottoapp.log
		format json
	}
}
```

## License

The frontend is built with Tailwind. The styles, components, and javascript are not open source and are not licensed to be used outside of this applicatin.

Proprietary License | https://tailwindcss.com/plus/license
* tailwindplus/blocks
* tailwindplus/elements 
