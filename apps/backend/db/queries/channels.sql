-- name: GetChannel :one
SELECT * FROM channels WHERE id = $1;

-- name: GetChannelsFromServer :many
SELECT * FROM channels WHERE server_id = $1 AND id = (SELECT channel_id FROM channel_membership WHERE user_id = $2);

-- name: CreateChannel :one
INSERT INTO channels (
  server_id, name, type, description, x, y
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateChannelName :exec
UPDATE channels SET name = $1 WHERE id = $2;

-- name: UpdateChannelDescription :exec
UPDATE channels SET description = $1 WHERE id = $2;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;
