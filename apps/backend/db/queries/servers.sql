-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: GetServers :many
SELECT * FROM servers;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*
FROM servers s, server_membership sm
WHERE s.private = false OR (sm.server_id = s.id AND sm.user_id = $1);

-- name: CreateServer :one
INSERT INTO servers (
  owner_id, name, background, description, x, y, private
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: JoinServer :exec
INSERT INTO server_membership (
  user_id, server_id
) VALUES (
  $1, $2
);

-- name: UpdateServerName :exec
UPDATE servers SET name = $1 WHERE id = $2 AND owner_id = $3;

-- name: UpdateServerBackground :exec
UPDATE servers SET background = $1 WHERE id = $2 AND owner_id = $3;

-- name: UpdateServerDescription :exec
UPDATE servers SET description = $1 WHERE id = $2 AND owner_id = $3;

-- name: DeleteServer :execresult
DELETE FROM servers WHERE id = $1 AND owner_id = $2;
