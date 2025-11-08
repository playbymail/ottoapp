--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- name: CreateDocument :one
INSERT INTO documents (mime_type, contents_hash, content_length, created_at, updated_at)
VALUES (:mime_type, :contents_hash, :content_length, :created_at, :updated_at)
RETURNING document_id;

-- name: CreateDocumentContent :exec
INSERT INTO document_contents(document_id, contents, created_at, updated_at)
VALUES (:document_id, :contents, :created_at, :updated_at);

-- name: CreateDocumentAcl :exec
INSERT INTO document_acl(document_id, user_id, document_name, created_by, is_owner, can_read, can_write, can_delete, created_at, updated_at)
VALUES (:document_id, :user_id, :document_name, :created_by, :is_owner, :can_read, :can_write, :can_delete, :created_at, :updated_at);
