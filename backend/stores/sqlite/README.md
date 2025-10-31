# Sqlite database

## Write-ahead logging (WAL)

WAL creates two sidecar files: `ottoapp.db-wal` and `ottoapp.db-shm`.
They’re safe to delete only when no connections are open and after a checkpoint (SQLite recreates them as needed).

To check current settings:
```sql
PRAGMA journal_mode;
PRAGMA synchronous;
PRAGMA wal_autocheckpoint;
```

See current WAL status:

```sql
-- Returns three integers: (busy, log, checkpointed)
PRAGMA wal_checkpoint;    -- passive status check
```

One-off checkpoints (SQL):
```sql
-- Most common: wait for writers to pause, then copy frames into the db,
-- reset the WAL so it can be reused; size may remain.
PRAGMA wal_checkpoint(FULL);

-- Like FULL, but also resets the WAL file to size 0.
PRAGMA wal_checkpoint(RESTART);

-- Strongest clean-up: copy, reset *and* truncate the WAL to 0 bytes.
PRAGMA wal_checkpoint(TRUNCATE);

-- Non-blocking: do as much as possible without interfering with writers.
PRAGMA wal_checkpoint(PASSIVE);
```

What to use when:

* Routine maintenance / before backups: `FULL`
* Before shipping the DB file (avoid needing the .wal): `RESTART` or `TRUNCATE`
* During live traffic (don’t block writers): `PASSIVE`

Tip: If you plan to copy just `ottoapp.db`, run `PRAGMA wal_checkpoint(TRUNCATE)`; first so you don’t forget the `.wal` file.

## Housekeeping

After a large delete or vacuum-worthy change:

```sql
PRAGMA wal_checkpoint(TRUNCATE);
VACUUM;                 -- compacts the main db file
```
