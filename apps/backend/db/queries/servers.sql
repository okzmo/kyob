-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: GetServers :many
SELECT * FROM servers;

-- name: IsMember :execresult
SELECT id FROM server_membership WHERE server_id = $1 AND user_id = $2;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*, sm.x, sm.y
FROM servers s, server_membership sm
WHERE sm.server_id = s.id AND sm.user_id = $1;

-- name: CreateServer :one
INSERT INTO servers (
  owner_id, name, avatar, description, private
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: JoinServer :exec
INSERT INTO server_membership (
  user_id, server_id, x, y
) VALUES (
  $1, $2, $3, $4
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

