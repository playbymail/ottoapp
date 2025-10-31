--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE roles
(
    role_id     TEXT PRIMARY KEY,
    is_active   BOOL      NOT NULL,
    description TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

insert into roles (role_id, is_active, description)
VALUES ('active', 1, 'active user role'),
       ('sysop', 1, 'sysop role'),
       ('admin', 1, 'administrator role'),
       ('guest', 1, 'guest / anonymous visitor role')
;

CREATE TABLE user_roles
(
    user_id    INTEGER   NOT NULL,
    role_id    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (role_id) ON DELETE CASCADE
);

insert into user_roles (user_id, role_id)
select user_id, 'active'
from users
where handle = 'sysop';

insert into user_roles (user_id, role_id)
select user_id, 'sysop'
from users
where handle = 'sysop';

INSERT INTO schema_version (version, applied_at)
VALUES (4, current_timestamp);
