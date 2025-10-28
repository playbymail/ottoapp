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
