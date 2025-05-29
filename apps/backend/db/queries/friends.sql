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
SELECT u.id, u.display_name, u.avatar, u.about, f.accepted, f.id AS friendship_id, 
       f.user_id AS friendship_sender_id, c.id AS channel_id
FROM users u
INNER JOIN friends f ON u.id = f.friend_id
LEFT JOIN channels c ON $1 = ANY(c.users) AND u.id::text = ANY(c.users)
WHERE f.user_id = $1

UNION

SELECT u.id, u.display_name, u.avatar, u.about, f.accepted, f.id AS friendship_id, 
       f.user_id AS friendship_sender_id, c.id AS channel_id
FROM users u
INNER JOIN friends f ON u.id = f.user_id  
LEFT JOIN channels c ON $1 = ANY(c.users) AND u.id::text = ANY(c.users)
WHERE f.friend_id = $1;

-- name: GetExistingChannel :one
UPDATE channels SET active = true
WHERE type = 'dm'
  AND array_length(users, 1) = 2
  AND $1::varchar = ANY(users) 
  AND $2::varchar = ANY(users)
RETURNING *;
