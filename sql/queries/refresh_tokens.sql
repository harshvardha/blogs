-- name: CreateRefreshToken :exec
insert into refresh_token(
    token,
    user_id,
    expires_at,
    created_at,
    updated_at
)
values (
    $1,
    $2,
    $3,
    NOW(),
    NOW()
);

-- name: GetRefreshToken :one
select expires_at from refresh_token where user_id = $1;