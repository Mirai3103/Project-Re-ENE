-- name: GetConversation :one
SELECT *
FROM conversations
WHERE id = ?
LIMIT 1;

-- name: CreateConversation :exec
INSERT INTO conversations (id, max_window_size, character_id, user_id)
VALUES (?, ?, ?, ?);

-- name: CreateConversationMessage :exec
INSERT INTO conversation_messages (id, conversation_id, content, role, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: ListConversationMessages :many
SELECT *
FROM conversation_messages
WHERE conversation_id = ?;

-- name: GetConversationWindowSize :one
SELECT max_window_size
FROM conversations
WHERE id = ?;

-- name: ListRecentMessages :many
SELECT *
FROM conversation_messages
WHERE conversation_id = ?
ORDER BY created_at DESC
LIMIT ?;