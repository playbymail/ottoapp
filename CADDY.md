# Caddy

## Caddy
When installed with `brew`, Caddy's configuration will be set as `/opt/homebrew/etc/Caddyfile`.

Run Caddy in the foreground for debugging:

```bash
XDG_DATA_HOME="/opt/homebrew/var/lib" HOME="/opt/homebrew/var/lib" caddy run --config /opt/homebrew/etc/Caddyfile --watch
```

## Sample development

NOTE: *.localhost automatically resolves to 127.0.0.1 in modern OS/browsers;
you usually don't need an /etc/hosts entry for ottoapp.localhost.


```caddy
{
  # stay unprivileged in dev
  http_port 8080
  https_port 8443
}

ottoapp.localhost {
  tls internal
  encode zstd gzip

  # Optional access log while developing
  log {
    output file /opt/homebrew/var/log/caddy/ottoapp.dev.access.log
    format console
  }

  # --- API -> Go server on :8181 ---
  @api path /api/*
    handle @api {
    reverse_proxy localhost:8181
  }

  # --- Everything else -> Ember dev server on :4200 (HMR, WS/SSE handled automatically) ---
  handle {
    reverse_proxy localhost:4200
  }

  # Safer for dev (don’t pin HSTS)
  header {
    Strict-Transport-Security "max-age=0"
  }
}
```

## Sample file for production

```caddy
ottoapp.mdhenderson.com, www.ottoapp.mdhenderson.com {
  encode zstd gzip

  # Access log (tweak path/format to taste)
  log {
    output file /var/log/caddy/ottoapp.access.log
    format json
  }

  # --- API -> Go (bind your Go service to 127.0.0.1:8181) ---
  @api path /api/*
    handle @api {
    reverse_proxy 127.0.0.1:8181
  }

  # --- Ember build output from disk ---
  handle {
    root * /var/www/ottoapp.mdhenderson.com/app
    try_files {path} /index.html    # SPA fallback
    file_server
  }
  # strong asset caching on fingerprinted files
  @immutable {
    path_regexp hashed ^/(assets|static)/.*\.[0-9a-f]{8,}\.(js|css|png|jpg|webp|svg|woff2?)$
  }
  header @immutable Cache-Control "public, max-age=31536000, immutable"

  # Security headers (safe defaults; tune as needed)
  header {
    # 180 days HSTS; add preload when you’re ready to submit to preload list
    Strict-Transport-Security "max-age=15552000; includeSubDomains"
    X-Content-Type-Options "nosniff"
    Referrer-Policy "strict-origin-when-cross-origin"
    Permissions-Policy "geolocation=(), microphone=(), camera=()"
  }

  # Optional: ACME account email
  # tls you@example.com
}
```
