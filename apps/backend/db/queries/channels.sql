-- name: GetChannel :one
SELECT * FROM channels WHERE id = $1;

-- name: GetChannelsFromServer :many
SELECT *
FROM channels
WHERE server_id = $1;

-- name: CreateChannel :one
INSERT INTO channels (
  id, server_id, name, type, description, users, roles, x, y
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateChannelName :exec
UPDATE channels SET name = $1 WHERE id = $2;

-- name: UpdateChannelDescription :exec
UPDATE channels SET description = $1 WHERE id = $2;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;
