-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessagesFromChannel :many
SELECT * FROM messages WHERE channel_id = $1;

-- name: CreateMessage :one
INSERT INTO messages (
  author_id, channel_id, content, mentions_users, mentions_channels, attached
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1;
