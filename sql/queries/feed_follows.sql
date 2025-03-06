-- name: CreateFeedFollow :one
with inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT
inserted_feed_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users on inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.id, 
    ff.created_at, 
    ff.updated_at, 
    ff.user_id, 
    ff.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows ff
INNER JOIN feeds ON ff.feed_id = feeds.id
INNER JOIN users ON ff.user_id = users.id
WHERE ff.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
USING feeds
WHERE feed_follows.feed_id = feeds.id
AND feed_follows.user_id = $1
AND feeds.url = $2;