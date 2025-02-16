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

-- name: GetUser :one
select * from users where email = $1;