-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: GetServers :many
SELECT * FROM servers;

-- name: IsMember :execresult
SELECT id FROM server_membership WHERE server_id = $1 AND user_id = $2;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*
FROM servers s, server_membership sm
WHERE s.private = false OR (sm.server_id = s.id AND sm.user_id = $1);

-- name: CreateServer :one
INSERT INTO servers (
  owner_id, name, avatar, description, x, y, private
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

-- name: UpdateServerAvatar :exec
UPDATE servers SET avatar = $1 WHERE id = $2 AND owner_id = $3;

-- name: UpdateServerBanner :exec
UPDATE servers SET banner = $1 WHERE id = $2 AND owner_id = $3;

-- name: UpdateServerDescription :exec
UPDATE servers SET description = $1 WHERE id = $2 AND owner_id = $3;

-- name: DeleteServer :execresult
DELETE FROM servers WHERE id = $1 AND owner_id = $2;

-- name: CheckServerPosition :execresult
SELECT id FROM servers WHERE (x BETWEEN $1-100 AND $1+100) AND (y BETWEEN $2-100 AND $2+100);
