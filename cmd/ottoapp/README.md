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

**Backup to a specific directory:**

```bash
ottoapp --db data/alpha db backup --output tmp
```

The output directory must exist (will not be created).

### Clone

Clone the database to a working copy for testing:

```bash
mkdir -p tmp/test
ottoapp --db data/alpha db clone tmp/test
```

This creates `tmp/test/ottoapp.db` - a clean, single-file copy with no WAL sidecars.

Safety features:
- Output directory must exist (will not be created)
- Refuses to overwrite if `ottoapp.db` already exists in output directory
- Source and destination paths must be different

**Using clones for testing:**

```bash
mkdir -p tmp/test
ottoapp --db data/alpha db clone tmp/test
ottoapp -N --db tmp/test [test commands]
rm -rf tmp/test
```

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

## Game Commands

### Import

Import users from a JSON file.
Creates new users or updates existing ones based on handle.

```bash
ottoapp game import data/inputs/alpha.json
```

The import process:
1. Validates all records before making any changes
2. For each user:
    - If handle doesn't exist: creates a new user
    - If handle exists: updates remaining fields
    - Generates passwords for blank entries using the phrases generator
    - Updates or creates user_secrets record

### Upload

Upload a document for a game:

```bash
ottoapp game upload 0987.docx --owner catbird --name 0301.0899-12.0987.report.docx
```

## Report Commands

### Extract

## Run Commands

### Genmake

Generate a Makefile for cumulative map generation.

```bash
ottoapp run genmake --root data/tn3.1
```

This scans the data directory for turn reports (both `.docx` and `.txt`) and generates a Makefile that:
1. Extracts text reports from `.docx` files where needed.
2. Builds maps cumulatively, meaning each map depends on its turn report AND all previous turn reports for that clan.
3. Handles allied data (multiple clans in the same folder) by including all text reports found in the folder for each turn.

Options:
- `--root` - root directory containing clan folders (default: "data/tn3.1")
- `--output` - path to write the Makefile (default: "data/tn3.1/maps.mk")

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
