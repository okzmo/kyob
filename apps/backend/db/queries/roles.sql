-- name: CreateRole :exec
INSERT INTO roles (
  id, server_id, name, color, description, abilities
) VALUES (
  $1, $2, $3, $4, $5, $6
);
