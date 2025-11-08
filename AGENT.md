# Ottoapp

Ottoapp implements the updated version of the Ottomap website (currently hosted at https://ottomap.mdhenderson.com).

## Project Overview

This project is a **full-stack web application** built with **Ember.js (v6.8+)** on the frontend and a **Go REST server** on the backend. It demonstrates a modern, stable stack focused on maintainability, compatibility, and security — suitable for both local development and production deployment behind **Caddy**.

Ember v6.8 defaults to Vite + Embroider.

### Tech Stack

| Layer                  | Technology                                                          | Notes                                                                                  |
| ---------------------- |---------------------------------------------------------------------|----------------------------------------------------------------------------------------|
| **Frontend**           | Ember.js (v6.8+)                                                    | Using Ember CLI, Ember Octane idioms, and Ember Simple Auth (ESA) for session handling. |
| **Backend**            | Go (`net/http`)                                                     | Pure stdlib REST API with cookie-based auth and optional JWTs.                         |
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
    6. `/api/logout` clears the session cookie.

* **Caddy Integration:**

    * Terminates TLS and proxies `/api/*` requests to Go backend.
    * Serves Ember frontend directly.
    * Handles same-domain cookie scoping automatically.

### Command Line applications

* cmd/gentz - rebuild the IANA timezone database if needed
* cmd/ottoapp - manage the database, start the server, and use the API to run commands
* cmd/ottomap - parse and render files locally

Use `cmd/ottoapp` for testing the database and the API service.

Use `cmd/ottomap` for testing the parsing and rendering packages using inputs on the file system.

## Development Setup

We have Caddy serving `https://ottoapp.localhost:8443` and forwarding requests to both Ember and Go.
Caddy is configured to forward `/api/` routes to the Go API server at `localhost:8181`.
All other routes are forwarded to EmberJS at `localhost:4200`.

For testing the backend, we can create a new instance on a temporary port.
We can start the server with a timeout to kill it after a delay.
We can pass the server a flag to enable a shutdown route to stop the server.

### Prerequisites

* Go ≥ 1.23
* Node ≥ 20
* Ember CLI (latest LTS)
* Caddy ≥ 2.7

### Running Temporary Instances for Testing

Use `:memory:` as the database path to create an in-memory database for testing the server.

```bash
$ go run ./cmd/ottoapp server serve --db :memory: --port 8181 &
$ sleep 1 # let the database initialize
$ curl http://127.0.0.1:8181/api/ping
{"status":"ok","msg":"pong"}%
```

Set the shutdown timer flag stop the server after a short duration.

```bash
$ go run ./cmd/ottoapp server serve --db :memory: --shutdown-timer 1m &
$ sleep 1 # let the database initialize
$ curl http://127.0.0.1:8181/api/ping
{"status":"ok","msg":"pong"}%
$ sleep 60 # let the shutdown timer fire
$ curl http://127.0.0.1:8181/api/ping
curl: (7) Failed to connect to 127.0.0.1 port 8181 after 0 ms: Couldn't connect to server
```

Or set the key for the `/api/shutdown` route to stop the test instance on demand.

```bash
$ go run ./cmd/ottoapp server serve --db :memory: --shutdown-key foo &
$ sleep 1 # let the database initialize
$ curl -H "Content-Type: application/json" -d '{"key": "foo"}' http://127.0.0.1:8181/api/shutdown
{"status":"ok","msg":"shutdown initiated"}
$ sleep 5 # let the shutdown complete
$ curl http://127.0.0.1:8181/api/ping
curl: (7) Failed to connect to 127.0.0.1 port 8181 after 0 ms: Couldn't connect to server
```

#### 2. Start the Ember Frontend

```bash
cd frontend
npm start
```

## Design Notes

* **No SSR:** Ember runs as SPA; Go only serves JSON and cookies.
* **Stateless API:** Except for cookies, all endpoints are stateless.
* **Security:**

    * Session cookies are `Secure`, `HttpOnly`, `SameSite=Strict`.
    * TLS always terminated by Caddy.
    * Secrets configurable via `.env` or environment variables.

