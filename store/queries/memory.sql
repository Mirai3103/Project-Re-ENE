-- name: CreateMemory :exec
INSERT INTO memories (
  id, user_id, character_id, content, embedding, importance, confidence, source, tags
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetAllEmbeddings :many
SELECT id, embedding
FROM memories;

-- name: GetMemoriesByIDs :many
SELECT *
FROM memories
WHERE id IN (sqlc.slice('ids'));

-- name: DeleteMemory :exec
DELETE FROM memories
WHERE id = ?;