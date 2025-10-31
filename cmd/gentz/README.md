# gentz
Gentz is a tool to extract the list of time zone names from the Go distribution's internal data file.

We use it to normalize the time zone names in the database.

Usage:

```bash
go run ./cmd/gentz
```

This will overwrite the file `backend/iana/normalize.go`.