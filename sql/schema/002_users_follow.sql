-- +goose Up
create table users_follow(
    follower_id UUID not null,
    following_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(follower_id, following_id)
);

-- +goose Down
drop table users_follow;