--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

CREATE TABLE schema_migrations
(
    id         INTEGER PRIMARY KEY,
    name       TEXT      NOT NULL UNIQUE,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- sqlite timestamp should be UTC
);

CREATE TABLE schema_version
(
    id         INTEGER PRIMARY KEY,
    version    INTEGER   NOT NULL UNIQUE,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- sqlite timestamp should be UTC
);

CREATE TABLE config
(
    key        TEXT      NOT NULL,
    value      TEXT      NOT NULL,

    -- columns for auditing
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC

    PRIMARY KEY (key)
);

INSERT INTO schema_version (version)
VALUES (1);

INSERT INTO config (key, value)
VALUES ('schema_version', '20251029_1540');
