--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- name: CreateDocumentContents :exec
INSERT INTO document_contents(contents_hash, content_length, mime_type, contents, created_at, updated_at)
VALUES (:contents_hash, :content_length, :mime_type, :contents, :created_at, :updated_at)
ON CONFLICT (contents_hash) DO UPDATE
    SET updated_at = excluded.updated_at;

-- name: DeleteDocumentContents :exec
DELETE
FROM document_contents
WHERE contents_hash = :contents_hash
  AND NOT EXISTS (SELECT 1
                  FROM documents
                  WHERE contents_hash = :contents_hash);

-- name: GetDocumentContents :one
SELECT content_length,
       contents_hash,
       mime_type,
       contents,
       created_at,
       updated_at
FROM document_contents
WHERE contents_hash = :contents_hash;

-- name: CreateDocument :one
INSERT INTO documents (clan_id,
                       can_read, can_write, can_delete, can_share,
                       document_name,
                       document_type,
                       contents_hash,
                       created_at, updated_at)
VALUES (:clan_id,
        :can_read, :can_write, :can_delete, :can_share,
        :document_name,
        :document_type,
        :contents_hash,
        :created_at, :updated_at)
ON CONFLICT (clan_id, contents_hash) DO UPDATE
    SET can_read      = excluded.can_read,
        can_write     = excluded.can_write,
        can_delete    = excluded.can_delete,
        can_share     = excluded.can_share,
        document_name = excluded.document_name,
        updated_at    = excluded.updated_at
RETURNING document_id;

-- name: DeleteDocument :exec
DELETE
FROM documents
WHERE document_id = :document_id
  AND clan_id = :clan_id;

-- name: DeleteDocumentAuthorized :exec
DELETE
FROM documents
WHERE document_id = :document_id
  AND clan_id = :clan_id
  AND can_delete = 1;

-- name: ShareDocumentById :exec
INSERT INTO document_shares (document_id, clan_id, can_read, can_delete, created_at, updated_at)
VALUES (:document_id, :clan_id, :can_read, :can_delete, :created_at, :updated_at);

-- name: DeleteSharedDocumentById :exec
DELETE
FROM document_shares
WHERE document_id = :document_id
  AND clan_id = :clan_id;

-- GetDocumentAcl returns the list of the users with access to the document.
--
-- name: GetDocumentAcl :many
SELECT d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clan_documents_vw AS d
WHERE d.document_id = :document_id;

-- name: GetDocumentById :one
SELECT d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.created_at,
       d.updated_at
FROM documents AS d
WHERE d.document_id = :document_id;

-- name: GetDocumentByIdAuthorized :one
SELECT d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clan_documents_vw AS d
WHERE d.document_id = :document_id
  AND d.clan_id = :clan_id;

-- name: GetDocumentForUserAuthorized :one
SELECT c.user_id,
       c.game_id,
       d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clans AS c,
     clan_documents_vw AS d
WHERE d.document_id = :document_id
  AND c.user_id = :user_id;


-- name: GetAllDocumentsForClan :many
SELECT d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clan_documents_vw AS d
WHERE d.clan_id = :clan_id
   OR d.owner_id = :clan_id;

-- name: GetAllDocumentsForGameAcrossUsers :many
SELECT c.user_id,
       c.game_id,
       d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clans AS c,
     clan_documents_vw AS d
WHERE c.game_id = :game_id
  AND (d.clan_id = c.clan_id OR d.owner_id = c.clan_id);

-- name: GetAllDocumentsForUserAcrossGames :many
SELECT c.user_id,
       c.game_id,
       d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clans AS c,
     clan_documents_vw AS d
WHERE c.user_id = :user_id
  AND (d.clan_id = c.clan_id OR d.owner_id = c.clan_id);

-- name: GetAllDocumentsForUserInGame :many
SELECT c.user_id,
       c.game_id,
       d.document_id,
       d.clan_id,
       d.can_read,
       d.can_write,
       d.can_delete,
       d.can_share,
       d.document_name,
       d.document_type,
       d.contents_hash,
       d.owner_id,
       d.is_shared,
       d.created_at,
       d.updated_at
FROM clans AS c,
     clan_documents_vw AS d
WHERE c.user_id = :user_id
  AND c.game_id = :game_id
  AND (d.clan_id = c.clan_id OR d.owner_id = c.clan_id);

