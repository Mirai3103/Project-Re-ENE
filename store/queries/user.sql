-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserFacts :many
SELECT *
FROM user_facts
WHERE user_id = ?
ORDER BY updated_at DESC
LIMIT ?;

-- name: AddUserFact :exec
INSERT INTO user_facts (id, user_id, name, value, type)
VALUES (?, ?, ?, ?, ?);