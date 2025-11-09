# OttoApp

Rewrite of the OttoMap web server for v2

## Quick Start

### Development

Run the development environment (includes Ember with hot reload, Go with auto-rebuild via Air, and Tailwind watch):

```bash
./tools/run-dev.sh
```

This starts:
- Ember dev server at `:4200` (with hot module reload)
- Go API server at `:8181` (uses Air for auto-rebuilds on file changes)
- Tailwind CSS watch (auto-compiles on CSS changes)

Note: the script assumes that `tailwind` is in your `~/bin`; update as needed.

Access the app at `https://ottoapp.localhost:8443` (Caddy proxies both services).

See [CADDY.md](CADDY.md) for Caddy configuration details.

### Production Build

Build a production distribution tarball:

```bash
./tools/build.sh
```

This creates a tarball in `dist/prod/ottoapp-{version}.tgz` containing:
- Linux AMD64 binary
- Ember production build

## Testing

Test the login endpoint:

```bash
curl https://ottoapp.localhost:8443/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"penguin@ottoapp","password":"sardines-mmmmm"}'
```

## Deployment

See [DEPLOYING.md](DEPLOYING.md) for production deployment instructions.

## Acknowledgements

This project was developed with assistance from AI tools including Amp (by Sourcegraph) and ChatGPT, which were used throughout the design, implementation, testing, and documentation phases. While these tools provided valuable support, all final decisions, code review, and any errors or omissions remain the responsibility of the project authors.

## License

The frontend is built with Tailwind. The styles, components, and javascript are not open source and are not licensed to be used outside of this application.

Proprietary License | https://tailwindcss.com/plus/license
* tailwindplus/blocks
* tailwindplus/elements
