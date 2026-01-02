-- name: CreateDocument :one
INSERT INTO documents (clan_id,
                       document_name, document_type,
                       modified_at,
                       created_at, updated_at)
VALUES (:clan_id,
        :document_name, :document_type,
        :modified_at,
        :created_at, :updated_at)
RETURNING document_id;

-- name: ReadDocumentById :one
SELECT documents.document_id,
       documents.clan_id,
       documents.document_name,
       documents.document_type,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
FROM documents
WHERE documents.document_id = :document_id;

-- name: ReadDocumentByClanAndName :one
SELECT clans.game_id,
       clans.user_id,
       clans.clan,
       documents.document_id,
       documents.clan_id,
       documents.document_name,
       documents.document_type,
       document_contents.contents_hash,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
FROM clans,
     documents,
     document_contents
WHERE clans.clan_id = :clan_id
  AND documents.clan_id = clans.clan_id
  AND documents.document_name = :document_name
  AND document_contents.document_id = documents.document_id;


-- name: DeleteDocumentById :exec
DELETE
FROM documents
WHERE documents.document_id = :document_id;

-- name: DeleteDocumentByIdAuthorized :exec
DELETE
FROM documents
WHERE documents.document_id = :document_id
  AND clan_id = :clan_id;

-- name: DeleteDocumentByClanAndNameAuthorized :exec
DELETE
FROM documents
WHERE clan_id = :clan_id
  AND document_name = :document_name;

-- name: UpdateDocumentById :exec
UPDATE documents
SET document_name = :document_name,
    document_type = :document_type,
    modified_at   = :modified_at,
    updated_at    = :updated_at
WHERE document_id = :document_id
  AND clan_id = :clan_id;

-- name: UpdateDocumentByIdAuthorized :exec
UPDATE documents
SET document_name = :document_name,
    document_type = :document_type,
    modified_at   = :modified_at,
    updated_at    = :updated_at
WHERE document_id = :document_id
  AND clan_id = :clan_id;


-- name: CreateDocumentContents :exec
INSERT INTO document_contents(document_id,
                              content_length,
                              contents_hash,
                              contents,
                              created_at,
                              updated_at)
VALUES (:document_id,
        :content_length,
        :contents_hash,
        :contents,
        :created_at,
        :updated_at);

-- name: ReadDocumentContents :one
SELECT contents
FROM document_contents
WHERE document_id = :document_id;

-- name: ReadDocumentContentsByIdAuthorized :one
SELECT clans.game_id,
       clans.user_id,
       clans.clan,
       documents.document_id,
       documents.clan_id,
       documents.document_name,
       documents.document_type,
       document_types.content_type,
       document_contents.contents,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
FROM documents,
     document_types,
     document_contents,
     clans
WHERE documents.document_id = :document_id
  AND documents.clan_id = :clan_id
  AND document_types.document_type = documents.document_type
  AND document_contents.document_id = documents.document_id
  AND clans.clan_id = documents.clan_id;


-- name: UpdateDocumentContentsById :exec
UPDATE document_contents
SET content_length = :content_length,
    contents_hash  = :contents_hash,
    contents       = :contents,
    updated_at     = :updated_at
WHERE document_contents.document_id = :document_id;

-- name: DeleteDocumentContentsById :exec
DELETE
FROM document_contents
WHERE document_contents.document_id = :document_id;


-- name: ReadDocumentsByHash :one
SELECT clans.game_id,
       clans.user_id,
       clans.clan,
       documents.document_id,
       documents.clan_id,
       documents.document_name,
       documents.document_type,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
FROM document_contents,
     documents,
     clans
WHERE document_contents.contents_hash = :contents_hash
  AND documents.document_id = document_contents.document_id
  AND clans.clan_id = documents.clan_id;

-- name: ReadDocumentsByGameAndClanNo :many
SELECT clans.game_id,
       clans.user_id,
       clans.clan,
       documents.document_id,
       documents.clan_id,
       documents.document_name,
       documents.document_type,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
FROM clans,
     documents
WHERE clans.game_id = :game_id
  AND clans.clan_id = :clan_no
  AND documents.document_id = :document_id
  AND documents.clan_id = clans.clan_id;

-- name: ReadDocumentsByUser :many
SELECT clans.game_id,
       clans.user_id,
       clans.clan,
       documents.document_id,
       documents.clan_id,
       documents.document_name,
       documents.document_type,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
FROM clans,
     documents
WHERE clans.user_id = :user_id
  AND documents.clan_id = clans.clan_id;

-- name: ReadDocumentOwner :one
SELECT clans.game_id,
       clans.user_id,
       clans.clan_id,
       clans.clan,
       clans.setup_turn,
       clans.is_active
FROM documents,
     clans
WHERE documents.document_id = :document_id
  AND clans.clan_id = documents.clan_id;

-- name: ReadReportExtracts :many
select documents.document_id,
       clans.game_id  as game_id,
       substr(documents.document_name, 6, 7)  as turn_no,
       clans.clan as clan,
       substr(documents.document_name, 6)     as document_name,
       documents.modified_at,
       documents.created_at,
       documents.updated_at
from documents, clans
where documents.document_type = 'txt'
and clans.clan_id = documents.clan_id
order by game_id, turn_no, clan;
