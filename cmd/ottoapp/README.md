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

### Import

Import users from a CSV file. Creates new users or updates existing ones based on email:

```bash
ottoapp user import data/inputs/players.csv
```

The CSV file must have these exact columns:
- **Clan** - 4-digit number starting with 0 (e.g., "0500")
- **User Name** - username for the user
- **Email** - email address (used as unique identifier)
- **Role** - "admin" or "user" (Penguin must have "admin", all others "user")
- **Timezone** - IANA timezone name (e.g., "America/New_York")
- **Password** - password or blank (blank passwords will be auto-generated)

The import process:
1. Validates all records before making any changes
2. For each user:
    - If email doesn't exist: creates a new user
    - If email exists: updates username and timezone if different
    - Generates passwords for blank entries using the phrases generator
    - Updates or creates user_secrets record
    - Ensures "active" role exists
    - Ensures the role from CSV (admin/user) exists
3. Updates the CSV file with any generated passwords

**Example CSV:**

```csv
Clan,User Name,Email,Role,Timezone,Password
0000,Penguin,penguin@ottoapp,admin,Antarctica/Palmer,joyous.pepper.frogs
0123,Tom Tomtom,tomtom.tom@example.com,user,Europe/London,
0987,Catbird,catbird@ottoapp,user,Europe/London,
```

After import, blank passwords will be filled in with generated values like `aroma.jump.rally.limb.chomp.crown`.

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
