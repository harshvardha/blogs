-- name: FollowUser :exec
insert into users_follow (follower_id, following_id, created_at, updated_at)
values ($1, $2, NOW(), NOW());

-- name: UnfollowUser :exec
delete from users_follow where follower_id = $1 and following_id = $2;

-- name: GetPair :one
select * from users_follow where follower_id = $1 and following_id = $2;