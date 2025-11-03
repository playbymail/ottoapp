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

INSERT INTO schema_version (version)
VALUES (1);
