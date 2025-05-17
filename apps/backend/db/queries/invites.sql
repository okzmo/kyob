-- name: CreateInvite :one
INSERT INTO invites (
  server_id, invite_id, expire_at
) VALUES (
  $1, $2, $3
)
RETURNING invite_id;

-- name: CheckInvite :one
SELECT server_id FROM invites WHERE invite_id = $1 AND expire_at >= NOW();
