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
SELECT r.id, r.idx, r.name, r.color, r.abilities, array_agg(sm.user_id) FILTER (WHERE sm.user_id IS NOT NULL) AS members
FROM roles r
LEFT JOIN server_membership sm on r.id = ANY(sm.roles)
WHERE r.server_id = $1
GROUP BY r.id;

-- name: GetRole :one
SELECT r.id, r.idx, r.name, r.color, r.abilities, r.server_id
FROM roles r
WHERE r.id = $1;

-- name: AddRoleMember :exec
UPDATE server_membership 
SET roles = array_append(roles, $1) -- role_name
WHERE server_id = $2 AND user_id = $3;

-- name: RemoveRoleMember :exec
UPDATE server_membership SET roles = array_remove(roles, $1) WHERE server_id = $2 AND user_id = $3;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: RemoveRoleFromAllMembers :exec
UPDATE server_membership SET roles = array_remove(roles, $1) WHERE $1 = ANY(roles);

-- name: MoveRole :exec
UPDATE roles SET idx = $1 WHERE id = $2;

-- name: UpdateRolePositions :exec
UPDATE roles SET idx = idx + 1 WHERE idx >= $1 AND idx < $2;
