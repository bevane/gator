-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
    values (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    returning *
)
select inserted_feed_follow.*, feeds.name as feed_name, users.name as user_name
from inserted_feed_follow
join users on inserted_feed_follow.user_id = users.id
join feeds on inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
select *, users.name as user_name, feeds.name as feed_name from feed_follows
join users on feed_follows.user_id = users.id
join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1;


-- name: DeleteFeedFollowForUserByUrl :exec
delete from feed_follows
    where feed_follows.id in (
    select feed_follows.id from feed_follows
    join users on feed_follows.user_id = users.id
    join feeds on feed_follows.feed_id = feeds.id
    where feed_follows.user_id = $1
    and feeds.url = $2
);


