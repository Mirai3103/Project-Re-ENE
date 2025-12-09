-- name: GetCharacter :one
SELECT *
FROM characters
WHERE id = :id
LIMIT 1;

-- name: GetCharacters :many
SELECT *
FROM characters;

-- name: GetCharacterFacts :many
SELECT *
FROM character_facts
WHERE character_id = :character_id
ORDER BY updated_at DESC
LIMIT :limit;

-- name: AddCharacterFact :exec
INSERT INTO character_facts (id, character_id, name, value, type)
VALUES (:id, :character_id, :name, :value, :type);