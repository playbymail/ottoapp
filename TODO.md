# TODO

## Backend

- Implement "app" commands to load users, reports, etc. for in-memory database configuration
  - Currently we can start the server using an in-memory database but we can't configure it yet

## Command Line

- Add helper function for common database operations in cmd/ottoapp/main.go
  - Extract pattern of getting `--db` flag and opening database connection
  - Many commands repeat this code: `cmd.Flags().GetString("db")` followed by `sqlite.Open()`
  - Helper should handle errors consistently and reduce boilerplate
