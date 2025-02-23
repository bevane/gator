-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
returning *;

-- name: GetFeedsWithUserName :many
select feeds.name, feeds.url, users.name from feeds
join users on feeds.user_id = users.id;


-- name: GetFeedByUrl :one
select * from feeds
where url = $1;
