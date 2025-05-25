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
