--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- The Document_Types table is used when sending the document contents via HTTP.
-- MIME Types
--  DOCX        application/vnd.openxmlformats-officedocument.wordprocessingml.document
--  TURN_REPORT application/tn-3.0
--  TURN_REPORT application/tn-3.1
--  WXX         application/wxx.xml
--
CREATE TABLE document_types
(
    document_type TEXT NOT NULL,
    document_ext  TEXT NOT NULL,
    descr         TEXT NOT NULL,
    content_type  TEXT NOT NULL, -- for HTTP headers when sending content
    PRIMARY KEY (document_type)
);

insert into document_types (document_type, document_ext, descr, content_type)
values ('turn-report-file', 'docx', 'Word document',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document'),
       ('turn-report-extract', 'txt', 'Turn Report extract',
        'text/plain; charset=UTF-8'),
       ('worldographer-map', 'wxx', 'Worldographer document',
        'application/wxx.xml');

-- The Documents table for documents (e.g., turn reports, maps).
--
-- The Document_Type column is used to categorize or filter
-- documents in the API.
--
-- Note: the upload process should enforce conventions for the document_name.
-- Our convention is {game}.{turn}.{unitId}.{documentType}
-- Example: 0301.0899-12.0987.report.docx
CREATE TABLE documents
(
    document_id   INTEGER PRIMARY KEY AUTOINCREMENT,
    clan_id       INTEGER NOT NULL,

    document_name TEXT    NOT NULL,                         -- clan's name for the doc
    document_type TEXT    NOT NULL,
    modified_at   INTEGER NOT NULL CHECK (modified_at > 0), -- timestamp for file

    -- audit (unix seconds, UTC)
    created_at    INTEGER NOT NULL,                         -- set in app
    updated_at    INTEGER NOT NULL,                         -- set in app

    -- enforce unique file names
    UNIQUE (clan_id, document_name),

    FOREIGN KEY (clan_id)
        REFERENCES clans (clan_id),
    FOREIGN KEY (document_type)
        REFERENCES document_types (document_type)
);

-- index for "show me all docs I own"
CREATE INDEX idx_documents_owner
    ON documents (clan_id);

-- The Document_Contents table stores the data for a document.
CREATE TABLE document_contents
(
    document_id    INTEGER NOT NULL,

    content_length INTEGER NOT NULL, -- size in bytes
    contents_hash  TEXT    NOT NULL, -- hex encoded SHA-256
    contents       BLOB    NOT NULL,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL, -- set in app
    updated_at     INTEGER NOT NULL, -- set in app

    FOREIGN KEY (document_id)
        REFERENCES documents (document_id)
        ON DELETE CASCADE
);

-- The Turn_Reports table contains meta-data for turn reports.
--
-- We impose a constraint on the report - every section in it
-- must be for the same turn and all elements must be in the
-- same Clan.
CREATE TABLE turn_reports
(
    turn_report_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id        INTEGER NOT NULL,                      -- game the report was created for
    user_id        INTEGER NOT NULL,                      -- user that owns the report
    document_id    INTEGER NOT NULL,
    turn_no        INTEGER NOT NULL CHECK (turn_no >= 0), -- turn number from the report
    clan_id        INTEGER NOT NULL,

    -- audit (unix seconds, UTC)
    created_at     INTEGER NOT NULL,                      -- set in app
    updated_at     INTEGER NOT NULL,                      -- set in app

    FOREIGN KEY (clan_id)
        REFERENCES clans (clan_id)
        ON DELETE CASCADE,
    FOREIGN KEY (document_id)
        REFERENCES documents (document_id)
        ON DELETE CASCADE,
    FOREIGN KEY (game_id)
        REFERENCES games (game_id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE CASCADE
);
