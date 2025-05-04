-- name: CreateRole :exec
INSERT INTO roles (
  server_id, name, color, description, abilities
) VALUES (
  $1, $2, $3, $4, $5
);
