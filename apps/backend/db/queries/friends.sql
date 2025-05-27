-- name: AddFriend :one
INSERT INTO friends (
  id, user_id, friend_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: AcceptFriend :exec
UPDATE friends SET accepted=true WHERE id=$1;

-- name: DeleteFriend :exec
DELETE FROM friends WHERE id=$1;

-- name: GetFriends :many
SELECT u.id, u.display_name, u.avatar, u.about, f.accepted, f.id AS friendship_id, f.user_id AS friendship_sender_id
FROM users u, friends f
WHERE f.user_id = $1 AND u.id = f.friend_id
UNION
SELECT u.id, u.display_name, u.avatar, u.about, f.accepted, f.id AS friendship_id, f.user_id AS friendship_sender_id
FROM users u, friends f
WHERE f.friend_id = $1 AND u.id = f.user_id;
