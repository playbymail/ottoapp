# OttoApp command line interface

## Configuration File

## Global Flags

* `--ignore-config-file`, `-N` to ignore any `ottoapp.json` files
* `--db` directory the database files are stored in

## API Server Commands
```bash
ottoapp api serve
```

## App Runner
These commands use the API to run processes on the server.

### Version
Fetch the server version.

```bash
ottoapp app version
```

## Database Commands

### Backup
The `ottoapp db backup` command creates a backup of the database.

### Compact
The `ottoapp db compact` command compacts the database and merges any WAL log files into the file.
This reduces the size of the database.

### Init
The `ottoapp db init` command initializes a new database.

### Migrate
### Migrate Up
The `ottoapp db migrate up` command applies missing migrations to the database.

### Version
```bash
ottoapp db version
```

## Report Commands
### Upload
```bash
ottoapp report upload 0134.docx --owner penguin --name 0301.0899-12.0134.report.docx
```

## User Commands
### Create
Create a new user.

```bash
ottoapp user create penguin --tz America/New_York --email penguin@ottoapp --password happy-feet
```

### Update
```bash
ottoapp user update penguin --password sardines-mmmmm
```

## Version Commands
Shows the version number.

```bash
ottoapp version
```
