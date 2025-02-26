-- +goose Up
create table comment_likes (
    user_id uuid not null,
    comment_id uuid not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(user_id, comment_id)
);

-- +goose Down
drop table comment_likes;