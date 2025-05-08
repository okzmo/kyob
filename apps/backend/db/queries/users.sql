-- name: GetUser :one
SELECT * FROM users WHERE email = $1 OR username = $2;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserFacts :many
SELECT label, value FROM facts WHERE user_id = $1;

-- name: GetUserLinks :many
SELECT label, url FROM links WHERE user_id = $1;

-- name: CreateUser :one
INSERT INTO users (
  email, username, display_name, avatar, password
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserDisplayName :exec
UPDATE users
  set display_name = $2
WHERE id = $1;

-- name: UpdateUserAvatar :exec
UPDATE users
  set avatar = $2
WHERE id = $1;

-- name: UpdateUserAbout :exec
UPDATE users
  set about = $2
WHERE id = $1;

-- name: UpdateUserUsername :exec
UPDATE users
  set username = $2
WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
  set email = $2
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
  set password = $2
WHERE id = $1;
