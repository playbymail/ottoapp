# Ottoapp

Ottoapp implements the updated version of the Ottomap website (currently hosted at https://ottomap.mdhenderson.com).

## Project Overview

This project is a **full-stack web application** built with **Ember.js (v6.8+)** on the frontend and a **Go REST server** on the backend. It demonstrates a modern, stable stack focused on maintainability, compatibility, and security — suitable for both local development and production deployment behind **Caddy**.

### Tech Stack

| Layer                  | Technology                                                          | Notes                                                                                   |
| ---------------------- |---------------------------------------------------------------------| --------------------------------------------------------------------------------------- |
| **Frontend**           | Ember.js (v6.8+)                                                    | Using Ember CLI, Ember Octane idioms, and Ember Simple Auth (ESA) for session handling. |
| **Backend**            | Go (`net/http`)                                                     | Pure stdlib REST API with cookie-based auth and optional JWTs.                          |
| **Web Server / Proxy** | Caddy                                                               | Handles HTTPS, static file serving, and reverse proxying for API requests.              |
| **Auth**               | Cookie sessions                                                     | JWT optional; default uses secure, HTTP-only cookies.                                   |
| **Storage**            | SQLite                                                              | Via `modernc.org/sqlite` — single binary deployment.                                    |
| **Environment**        | Identical setup in dev and production; Caddy proxies to Go backend. |                                                                                         |

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

* **Alternate Method:** JWT authentication can be enabled later for stateless APIs.

* **Caddy Integration:**

    * Terminates TLS and proxies `/api/*` requests to Go backend.
    * Serves Ember frontend directly.
    * Handles same-domain cookie scoping automatically.

---

### Ember Simple Auth (ESA) Configuration Example

#### 1. Install ESA

```bash
ember install ember-simple-auth
```

#### 2. Create `app/authenticators/cookie.js`

```js
// app/authenticators/cookie.js
import Base from 'ember-simple-auth/authenticators/base';

export default class CookieAuthenticator extends Base {
  async authenticate(username, password) {
    const response = await fetch('/api/login', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      throw new Error('Invalid credentials');
    }

    return { authenticated: true };
  }

  async restore() {
    const res = await fetch('/api/session', { credentials: 'include' });
    if (res.ok) {
      return { authenticated: true };
    }
    throw new Error('Session expired');
  }

  async invalidate() {
    await fetch('/api/logout', { method: 'POST', credentials: 'include' });
  }
}
```

#### 3. Configure ESA in `config/environment.js`

```js
ENV['ember-simple-auth'] = {
  routeAfterAuthentication: 'dashboard',
  routeAfterInvalidation: 'login',
  routeIfAlreadyAuthenticated: 'dashboard',
};
```

#### 4. Protect Routes with ESA Mixins

```js
// app/routes/dashboard.js
import Route from '@ember/routing/route';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default class DashboardRoute extends Route.extend(AuthenticatedRouteMixin) {}
```

#### 5. Login Controller

```js
// app/controllers/login.js
import Controller from '@ember/controller';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class LoginController extends Controller {
  @service session;
  username = '';
  password = '';

  @action
  async login(event) {
    event.preventDefault();
    try {
      await this.session.authenticate('authenticator:cookie', this.username, this.password);
      this.transitionToRoute('dashboard');
    } catch {
      alert('Login failed. Please check your credentials.');
    }
  }
}
```

---

## Go Backend Example for ESA Cookie Auth

Below is a minimal but production-ready example using the Go standard library.

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var sessions = make(map[string]string) // map[sessionID]username

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/session", handleSession)
	mux.HandleFunc("/api/logout", handleLogout)
	mux.HandleFunc("/api/dashboard", handleDashboard) // example protected route

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Simple check; replace with real user validation.
	if creds.Username != "admin" || creds.Password != "secret" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Create session
	sessionID := creds.Username + "-session"
	sessions[sessionID] = creds.Username

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true, // required for HTTPS via Caddy
		Expires:  time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}

func handleSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || sessions[cookie.Value] == "" {
		http.Error(w, "no active session", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(sessions, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || sessions[cookie.Value] == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user := sessions[cookie.Value]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to your dashboard, " + user,
	})
}
```

#### Notes

* Cookies use `Secure`, `HttpOnly`, and `SameSite=Strict`.
* Behind Caddy, this cookie will only be sent on HTTPS.
* Replace the in-memory `sessions` map with a persistent store (SQLite or Redis).
* This implementation matches ESA’s expectations:

    * `POST /api/login` → 200 if success, 401 if fail
    * `GET /api/session` → 200 if valid, 401 if invalid
    * `POST /api/logout` → clears cookie

---

## Development Setup

### Prerequisites

* Go ≥ 1.23
* Node ≥ 20
* Ember CLI (latest LTS)
* Caddy ≥ 2.7

### Running Locally

#### 1. Start the Go Backend

```bash
cd backend
go run ./cmd/server
```

#### 2. Start the Ember Frontend

```bash
cd frontend
ember serve --port 4200
```

#### 3. Run with Caddy Proxy (Recommended)

```bash
caddy run --config Caddyfile.dev
```

Caddy will:

* Serve Ember at `https://localhost/`
* Proxy `/api/*` to the Go backend
* Handle HTTPS via local certs

---

## Build and Deployment

### Frontend

```bash
cd frontend
ember build --environment=production
```

### Backend

```bash
cd backend
go build -o server ./cmd/server
```

### Caddy (Production)

```bash
caddy run --config Caddyfile
```

---

## Design Notes

* **No SSR:** Ember runs as SPA; Go only serves JSON and cookies.
* **Stateless API:** Except for cookies, all endpoints are stateless.
* **Security:**

    * Session cookies are `Secure`, `HttpOnly`, `SameSite=Strict`.
    * TLS always terminated by Caddy.
    * Secrets configurable via `.env` or environment variables.

---

## TODO

* [ ] Add `.env` configuration for session secret and database path.
* [ ] Replace in-memory session store with persistent store.
* [ ] Add CSRF protection middleware.
* [ ] Add ESA test helpers for integration tests.
* [ ] Deploy to DigitalOcean behind Caddy using HTTPS.
