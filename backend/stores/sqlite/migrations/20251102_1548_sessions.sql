--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE sessions
(
    session_id TEXT PRIMARY KEY,
    csrf       TEXT      NOT NULL,
    user_id    INTEGER   NOT NULL,
    expires_at TIMESTAMP NOT NULL, -- must always be UTC
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

INSERT INTO schema_version (version, applied_at)
VALUES (5, current_timestamp);
