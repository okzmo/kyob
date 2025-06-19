-- name: GetUser :one
SELECT * FROM users WHERE email = $1 OR username = $2;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUsersByIds :many
SELECT id, username, display_name, avatar FROM users WHERE id = ANY($1::text[]);

-- name: GetUserMinimal :one
SELECT id, username, display_name, avatar FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
  id, email, username, display_name, avatar, banner, body, main_color, password
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateEmoji :exec
UPDATE emojis SET shortcode = $1 WHERE user_id = $2 AND id = $3;

-- name: DeleteEmoji :exec
DELETE FROM emojis WHERE user_id = $1 AND id = $2;

-- name: GetEmojis :many
SELECT id, url, shortcode FROM emojis WHERE user_id = $1;

-- name: CreateEmoji :copyfrom
INSERT INTO emojis (
  id, user_id, url, shortcode
) VALUES (
  $1, $2, $3, $4
);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserDisplayName :exec
UPDATE users
  set display_name = $2
WHERE id = $1;

-- name: UpdateUserAvatarNBanner :exec
UPDATE users
  set avatar = $2, banner = $3, main_color = $4
WHERE id = $1;

-- name: UpdateUserAbout :exec
UPDATE users
  set about = $2
WHERE id = $1;

-- name: UpdateUserUsername :execresult
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

-- name: UpdateUserLinks :exec
UPDATE users
  set links = $2
WHERE id = $1;

-- name: UpdateUserFacts :exec
UPDATE users
  set facts = $2
WHERE id = $1;
