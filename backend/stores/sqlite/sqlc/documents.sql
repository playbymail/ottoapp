--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateDocument does
--
-- name: CreateDocument :one
INSERT INTO documents (document_id, document_created_by, document_created_at, document_path)
VALUES (:document_id, :document_created_by, :document_created_at, :document_path)
RETURNING document_id;
