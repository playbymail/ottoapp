--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE user_secrets
(
    user_id         INTEGER PRIMARY KEY,
    hashed_password TEXT      NOT NULL,
    last_login      INTEGER   NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

insert into user_secrets (user_id, hashed_password, last_login)
select user_id,
       '*',
       0
from users;

INSERT INTO schema_version (version, applied_at)
VALUES (3, current_timestamp);
