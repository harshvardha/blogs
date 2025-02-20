-- name: CreateUser :one
insert into users(id, username, email, hashed_password, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
returning id, username, email, created_at, updated_at;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: GetUserById :one
select * from users where id = $1;

-- name: UpdateUserEmailOrUsername :one
update users set email = $1 and username = $2, updated_at = NOW() where id = $3
returning id, username, email, created_at, updated_at;

-- name: DeleteUser :one
delete from users where id = $1
returning *;

-- name: GetUsersByUsername :many
select id, username from users where username = $1;

-- name: GetUserFeed :many
select * from blogs where author_id = (select following_id from users_follow where follower_id = $1) order by created_at;