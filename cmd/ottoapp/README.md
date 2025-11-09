# OttoApp Command Line Interface

Quick guide to using the `ottoapp` command.

## Quick Start

For quick tests, run directly with Go:

```bash
go run ./cmd/ottoapp [command]
```

For production, use the compiled binary:

```bash
./ottoapp [command]
```

## Configuration File

The command looks for `ottoapp.json` in the current directory by default.

## Global Flags

- `--ignore-config-file`, `-N` - ignore any `ottoapp.json` files
- `--db <path>` - path to the database directory (default: ".")
- `--debug` - enable debugging options
- `--dev` - enable development mode

## API Server Commands

### Serve

Start the API server:

```bash
ottoapp api serve
```

Options:
- `--port` - server port (default: 8181)
- `--shutdown-timer` - auto-shutdown after duration (e.g., "1m")
- `--shutdown-key` - enable `/api/shutdown` route with key

## App Runner Commands

These commands use the API to run processes on the server.

### Version

Fetch the server version:

```bash
ottoapp app version
```

## Database Commands

### Backup

Create a timestamped backup of the database using `VACUUM INTO`:

```bash
ottoapp db backup
```

This creates a clean, single-file copy (e.g., `backup-20250109-150405.db`) with no WAL sidecars.

**Using backups for testing:**

```bash
# Create a backup and test against a copy
ottoapp --db data/alpha db backup
mkdir -p tmp
cp data/alpha/backup-*.db tmp/test.db
ottoapp -N --db tmp [test commands]
```

Note: Future enhancement will support an optional output path.

### Compact

Compact the database and merge WAL log files:

```bash
ottoapp db compact
```

### Init

Initialize a new database:

```bash
ottoapp db init
```

### Migrate

Apply missing migrations to the database:

```bash
ottoapp db migrate up
```

### Version

Show database version:

```bash
ottoapp db version
```

## Report Commands

### Upload

Upload a report file:

```bash
ottoapp report upload 0134.docx --owner penguin --name 0301.0899-12.0134.report.docx
```

## User Commands

### Create

Create a new user:

```bash
ottoapp user create penguin --tz America/New_York --email penguin@ottoapp --password happy-feet
```

### Update

Update user password:

```bash
ottoapp user update penguin --password sardines-mmmmm
```

## Version Command

Show the application version:

```bash
ottoapp version
```
