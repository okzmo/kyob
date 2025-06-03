-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: GetServerWithChannels :one
SELECT DISTINCT s.*, sm.x, sm.y, (SELECT count(id) FROM server_membership smc WHERE smc.server_id=$1) AS member_count
FROM servers s, server_membership sm
WHERE s.id = $1 AND sm.server_id = s.id AND sm.user_id = $2;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: GetServers :many
SELECT * FROM servers;

-- name: IsMember :execresult
SELECT id FROM server_membership WHERE server_id = $1 AND user_id = $2;

-- name: GetServersCountFromUser :one
SELECT count(id)
FROM server_membership
WHERE user_id = $1;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*, sm.x, sm.y, (SELECT count(id) FROM server_membership smc WHERE smc.server_id=s.id) AS member_count
FROM servers s
LEFT JOIN server_membership sm ON sm.server_id = s.id AND sm.user_id = $1
WHERE s.id = 'global' OR sm.user_id IS NOT NULL;

-- name: GetServerMembers :many
SELECT u.id, u.username, u.display_name, u.avatar FROM server_membership sm, users u WHERE sm.server_id = $1 AND sm.user_id = u.id;

-- name: GetMembersFromServers :many
SELECT u.id, u.username, u.display_name, u.avatar, sm.server_id FROM server_membership sm, users u WHERE sm.server_id = ANY($1::text[]) AND sm.user_id = u.id;

-- name: CreateServer :one
INSERT INTO servers (
  id, owner_id, name, avatar, description, private
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: JoinServer :exec
INSERT INTO server_membership (
  id, user_id, server_id, x, y
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: LeaveServer :exec
DELETE FROM server_membership WHERE user_id = $1 AND server_id = $2;

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

