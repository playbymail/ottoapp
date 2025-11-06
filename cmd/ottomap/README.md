# OttoApp command line interface

## Database Commands

Note: You may run into problems running some commands if the server is also accessing the database.

### Init
The `ottoapp db init` command initializes a new database.

### Backup
The `ottoapp db backup` command creates a backup of the database.

### Compact
The `ottoapp db compact` command compacts the database and merges any WAL log files into the file.
This reduces the size of the database.

### Create
Create records in the database.

#### User
Create a new user.

```bash
ottoapp db --db testdata create user penguin2@ottomap --tz America/New_York --role chief
```

### Migrate
The `ottoapp db migrate` command applies missing migrations to the database.

### Update

#### User
```bash
ottoapp db update user penguin2@ottomap --password sardines.mmmmm
```