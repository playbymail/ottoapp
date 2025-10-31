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

## Development Setup

### Prerequisites

* Go ≥ 1.23
* Node ≥ 20
* Ember CLI (latest LTS)
* Caddy ≥ 2.7

### Running Locally

We can assume that Caddy is serving `https://ottoapp.localhost` and forwarding requests to both Ember and Go.

#### 1. Start the Go Backend

```bash
go run ./cmd/ottoapp
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

