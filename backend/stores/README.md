# SQLite Store

This package implements the **SQLite persistence layer** for the backend.

It is the **lowest-level** part of the data stack: it knows how to open the database, run migrations, compact or back up the file, and expose the `sqlc`-generated query methods. It does **not** hold business logic — that lives in feature packages.

This repo is organized so that:

- **`backend/stores/sqlite`** = infrastructure + generated SQL access
- **`backend/domains`** = shared types, errors, and constants used across features
- **`backend/<feature>`** (e.g. `auth`, `users`, `documents`, `sessions`) = business logic for that feature, built on top of the SQLite store and domain types
- **`frontend/`** = EmberJS app

That separation keeps the SQLite package from turning into a giant catch-all and makes "future me" happier when I have to change how auth works.

---

## High-Level Structure

```text
frontend/                 ← EmberJS
backend/
  domains/                ← shared domain types, IDs, errors, enums
  stores/
    sqlite/               ← SQLite infra + sqlc
  auth/                   ← feature logic, uses sqlite + domains
  users/                  ← feature logic, uses sqlite + domains
  documents/              ← feature logic, uses sqlite + domains
  sessions/               ← feature logic, uses sqlite + domains
  ...
````

### `backend/domains`

This package is the glue. It defines things the whole backend agrees on:

* core structs (`User_t`, `Document_t`, etc.)
* value objects / IDs (`type ID string` or whatever you're using)
* shared errors (`ErrInvalidEmail`, `ErrNotFound`, …)
* maybe small helpers for time or status

Because **every** feature package depends on these types, they should live in `backend/domains` and **not** inside the SQLite package. That way you can swap storage or add an in-memory version later without changing the domain model.

Example:

```go
// backend/domains/user.go
package domains

import "time"

type User_t struct {
    ID        ID
    Handle    string
    Email     string
    Timezone  string
    CreatedAt time.Time
}
```

Feature packages (like `backend/users`) import `domains` and `backend/stores/sqlite` and combine them.

---

## What This Package (sqlite) Actually Does

Files that belong here:

| Responsibility               | File(s)                        |
| ---------------------------- | ------------------------------ |
| DB config & paths            | `config.go`                    |
| Open/init SQLite             | `open.go`, `init.go`           |
| Core DB wrapper              | `db.go`                        |
| Migrations                   | `migrations.go`, `migrations/` |
| Maintenance (backup, vacuum) | `backup.go`, `compact.go`      |
| Generated queries            | `sqlc/`                        |
| sqlc config                  | `sqlc.yaml`                    |

The wrapper struct looks like:

```go
type DB struct {
    path string
    name string // ":memory:" for a temporary database
    db   *sql.DB
    q    *sqlc.Queries
}

func (d *DB) Queries() *sqlc.Queries { return d.q }
func (d *DB) Stdlib() *sql.DB        { return d.db }
```

That's intentionally small. Other packages should **take** this DB and do their work there.

---

## Why Separate SQLite From Features?

Feature logic now lives in separate packages:

* `backend/auth/service.go` - authentication logic
* `backend/users/service.go` - user management logic
* `backend/documents/service.go` - document management logic
* `backend/sessions/service.go` - session management logic

This separation has several advantages:

1. **You always know where to look.**
   Auth change? Go to `backend/auth`.
   Document change? Go to `backend/documents`.
   No more "where did I hang this helper on the DB?"

2. **SQLite stays infrastructure.**
   This folder is now about "how we talk to SQLite," not "how we do auth."

3. **Domains stay central.**
   `backend/domains` gives all feature packages a shared vocabulary, so you don't end up with slightly different `User` shapes all over the place.

4. **Clean dependency direction.**

    * `backend/domains` ← no one above it
    * `backend/stores/sqlite` → may *use* `domains` types for scan/return
    * `backend/<feature>` → uses both `domains` and `sqlite`

   Handlers at the top just wire it together.

---

## Wiring Example

Startup code (roughly):

```go
import (
    bsqlite "github.com/playbymail/ottoapp/backend/stores/sqlite"
    "github.com/playbymail/ottoapp/backend/auth"
    "github.com/playbymail/ottoapp/backend/users"
    "github.com/playbymail/ottoapp/backend/documents"
)

func main() {
    db, err := bsqlite.NewFromConfig(...)
    if err != nil {
        panic(err)
    }

    authSvc := auth.New(db)         // uses sqlite + domains
    userSvc := users.New(db)        // uses sqlite + domains
    docSvc := documents.New(db)     // uses sqlite + domains

    // pass these to HTTP handlers, RPC, etc.
}
```

Each service returns or accepts types from `backend/domains`, so everything above the services can stay consistent.

---

## Example: Adding a New Feature

Let's say we want "notifications."

1. **Add SQL** to `backend/stores/sqlite/sqlc/notifications.sql`:

   ```sql
   -- name: CreateNotification :one
   INSERT INTO notifications (user_id, message, created_at)
   VALUES (?, ?, unixepoch())
   RETURNING id, user_id, message, created_at;
   ```

   Run `sqlc generate` to get `notifications.sql.go` in `sqlc/`.

2. **Add a domain type** (optional but recommended) in `backend/domains/notification.go`:

   ```go
   package domains

   import "time"

   type Notification struct {
       ID        ID
       UserID    ID
       Message   string
       CreatedAt time.Time
   }
   ```

3. **Create a feature package**:

   ```text
   backend/notifications/service.go
   ```

   ```go
   package notifications

   import (
       "context"

       "github.com/playbymail/ottoapp/backend/stores/sqlite"
       "github.com/playbymail/ottoapp/backend/domains"
   )

   type Service struct {
       db *sqlite.DB
   }

   func New(db *sqlite.DB) *Service {
       return &Service{db: db}
   }

   func (s *Service) Create(ctx context.Context, userID domains.ID, msg string) (*domains.Notification, error) {
       q := s.db.Queries()

       rec, err := q.CreateNotification(ctx, sqlite.CreateNotificationParams{
           UserID:  string(userID),
           Message: msg,
       })
       if err != nil {
           return nil, err
       }

       return &domains.Notification{
           ID:        domains.ID(rec.ID),
           UserID:    domains.ID(rec.UserID),
           Message:   rec.Message,
           CreatedAt: rec.CreatedAt.Time,
       }, nil
   }
   ```

4. **Wire it up** in `cmd/.../main.go`.

That pattern (SQL → sqlc → domains → feature) is the whole point of this layout.

---

**Bottom line:**

* `backend/domains` defines what the app *talks about*
* `backend/stores/sqlite` defines how we *store* it
* `backend/<feature>` defines what we *do* with it

That's why we're doing it this way.
