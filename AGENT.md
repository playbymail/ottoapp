# Ottoapp

OttoApp implements the updated version of the OttoMap website.

## Project Overview

This project is a **full-stack web application** built with **Ember.js (v6.8+)** on the frontend and a **Go REST server** on the backend. It demonstrates a modern, stable stack focused on maintainability, compatibility, and security — suitable for both local development and production deployment behind **Caddy**.

Ember v6.8 defaults to Vite + Embroider.

### Tech Stack

| Layer                  | Technology                                                          | Notes                                                                                  |
| ---------------------- |---------------------------------------------------------------------|----------------------------------------------------------------------------------------|
| **Frontend**           | Ember.js (v6.8+)                                                    | Using Ember CLI, Ember Octane idioms, and Ember Simple Auth (ESA) for session handling. |
| **Backend**            | Go (`net/http`)                                                     | Pure stdlib REST API with cookie-based auth.                         |
| **Web Server / Proxy** | Caddy                                                               | Handles HTTPS, static file serving, and reverse proxying for API requests.             |
| **Auth**               | Cookie sessions                                                     | Default uses secure, HTTP-only cookies.                                                |
| **Storage**            | SQLite                                                              | Via `modernc.org/sqlite` — single binary deployment.                                   |
| **Environment**        | Identical setup in dev and production; Caddy proxies to Go backend. |                                                                                        |

## Authentication and Authorization

### Overview

* **Preferred Method:** Cookie-based sessions.

* **Flow:**

    1. User logs in via `/api/login` (POST JSON `{ username, password }`).
    2. Backend validates credentials and sets a secure, HTTP-only cookie.
    3. Ember Simple Auth manages session state.
    4. Backend checks session cookie for all protected routes.
    5. `/api/session` lets ESA restore or validate the session.
    6. POST `/api/logout` clears the session cookie.

* **Caddy Integration:**

    * Terminates TLS and proxies `/api/*` requests to Go backend.
    * Serves Ember frontend directly.
    * Handles same-domain cookie scoping automatically.

### Command Line applications

* cmd/gentz - rebuild the IANA timezone database if needed
* cmd/ottoapp - manage the database, start the server, and use the API to run commands

## Development Setup

We have Caddy serving `https://ottoapp.localhost:8443` and forwarding requests to both Ember and Go.
Caddy is configured to forward `/api/` routes to the Go API server at `localhost:8181`.
All other routes are forwarded to EmberJS at `localhost:4200`.

See [CADDY.md](CADDY.md) for Caddy configuration details.

### Development Instance

- Development instance database: `data/alpha`
- Test users: `penguin` and `catbird`
- Setup commands and credentials are in `data/users` file
- **Note:** `data/users` contains live email addresses; treat as private data

### Non-Destructive Testing

To run tests that might alter the database, work on a copy:

```bash
# Option 1: Use db clone command (recommended for testing)
mkdir -p tmp/test
go run ./cmd/ottoapp --db data/alpha db clone tmp/test
go run ./cmd/ottoapp -N --db tmp/test [commands]
rm -rf tmp/test

# Option 2: Use shell copy
mkdir -p tmp/foo
cp -r data/alpha/* tmp/foo/
go run ./cmd/ottoapp -N --db tmp/foo [commands]

# Option 3: Use backup command (creates timestamped backup)
mkdir -p tmp
go run ./cmd/ottoapp --db data/alpha db backup --output tmp
go run ./cmd/ottoapp -N --db tmp [commands]
```

For testing the backend, we can create a new instance on a temporary port.
We can start the server with a timeout to kill it after a delay.
We can pass the server a flag to enable a shutdown route to stop the server.

### Prerequisites

* Go ≥ 1.23
* Node ≥ 20
* Ember CLI (latest LTS)
* Caddy ≥ 2.7

### Running Temporary Instances for Testing

IMPORTANT: Use the `-N` flag when testing to avoid accidentally picking up a configuration from an `ottoapp.json` file.

Use `:memory:` as the database path to create an in-memory database for testing the server.

```bash
$ go run ./cmd/ottoapp -N api serve --db :memory: --port 8181 &
$ sleep 1 # let the database initialize
$ curl http://127.0.0.1:8181/api/ping
{"status":"ok","msg":"pong"}%
```

Set the shutdown timer flag stop the server after a short duration.

```bash
$ go run ./cmd/ottoapp -N api serve --db :memory: --shutdown-timer 1m &
$ sleep 1 # let the database initialize
$ curl http://127.0.0.1:8181/api/ping
{"status":"ok","msg":"pong"}%
$ sleep 60 # let the shutdown timer fire
$ curl http://127.0.0.1:8181/api/ping
curl: (7) Failed to connect to 127.0.0.1 port 8181 after 0 ms: Couldn't connect to server
```

Or set the key for the `/api/shutdown` route to stop the test instance on demand.

```bash
$ go run ./cmd/ottoapp -N api serve --db :memory: --shutdown-key foo &
$ sleep 1 # let the database initialize
$ curl -H "Content-Type: application/json" -d '{"key": "foo"}' http://127.0.0.1:8181/api/shutdown
{"status":"ok","msg":"shutdown initiated"}
$ sleep 5 # let the shutdown complete
$ curl http://127.0.0.1:8181/api/ping
curl: (7) Failed to connect to 127.0.0.1 port 8181 after 0 ms: Couldn't connect to server
```

### Start the Ember Frontend

Run Tailwind watch (if it isn't already running)

```bash
cd frontend && tailwindcss -i ./tailwind.css -o ./app/styles/app.css --watch
```

Start the development Ember server (if it isn't already running)

```bash
cd frontend && npm start
```

## Design Notes

* **No SSR:** Ember runs as SPA; Go only serves JSON and cookies.
* **Stateless API:** Except for cookies, all endpoints are stateless.
* **Security:**

    * Session cookies are `Secure`, `HttpOnly`, `SameSite=Strict`.
    * TLS always terminated by Caddy.

