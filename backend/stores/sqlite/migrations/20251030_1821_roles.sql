--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

CREATE TABLE roles
(
    role_id     TEXT PRIMARY KEY,
    is_active   BOOL      NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),

    description TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- sqlite timestamp should be UTC
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  -- sqlite timestamp should be UTC
);

insert into roles (role_id, is_active, description)
VALUES ('active', 1, 'active user role'),
       ('sysop', 1, 'sysop role'),
       ('admin', 1, 'administrator role'),
       ('player', 1, 'player role'),
       ('guest', 1, 'guest / anonymous visitor role'),
       ('tn3', 1, 'game TN3 role'),
       ('tn3.1', 1, 'game TN3.1 role')
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
select user_id, role_id
from users
         cross join (select roles.role_id
                     from roles
                     where role_id in ('active', 'sysop'))
where users.handle = 'sysop';

INSERT INTO schema_version (version)
VALUES (4);

UPDATE config
SET VALUE = '20251030_1821',
    updated_at = CURRENT_TIMESTAMP
WHERE key = 'schema_version';
