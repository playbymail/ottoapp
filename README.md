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

## License

The frontend is built with Tailwind. The styles, components, and javascript are not open source and are not licensed to be used outside of this applicatin.

Proprietary License | https://tailwindcss.com/plus/license
* tailwindplus/blocks
* tailwindplus/elements 
