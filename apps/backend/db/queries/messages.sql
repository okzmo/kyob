-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessagesFromChannel :many
SELECT * FROM messages WHERE channel_id = $1;

-- name: CheckChannelMembership :execresult
SELECT c.id FROM channels c, server_membership sm WHERE c.id = $1 and c.server_id = sm.server_id and sm.user_id = $2;

-- name: CreateMessage :one
INSERT INTO messages (
  author_id, channel_id, content, mentions_users, mentions_channels, attached
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateMessageContent :execresult
UPDATE messages SET content = $1 WHERE id = $2 AND author_id = $3;

-- name: UpdateMessageMentionsUsers :execresult
UPDATE messages SET mentions_users = $1 WHERE id = $2 AND author_id = $3;

-- name: UpdateMessageMentionsChannels :execresult
UPDATE messages SET mentions_channels = $1 WHERE id = $2 AND author_id = $3;

-- name: DeleteMessage :execresult
DELETE FROM messages WHERE id = $1 AND author_id = $2;
