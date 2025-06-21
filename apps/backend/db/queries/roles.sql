-- name: CreateRole :one
INSERT INTO roles (
  id, idx, server_id, name, color, abilities
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUserAbilities :many
SELECT r.abilities 
FROM roles r, server_membership sm 
WHERE sm.server_id = $1 AND sm.user_id = $2 AND r.id = ANY(sm.roles);

-- name: GetRoles :many
SELECT r.*
FROM roles r
WHERE r.server_id = $1;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: MoveRole :exec
UPDATE roles SET idx = $1 WHERE id = $2;

-- name: UpdateRolePositions :exec
UPDATE roles SET idx = idx + 1 WHERE idx >= $1 AND idx < $2;
