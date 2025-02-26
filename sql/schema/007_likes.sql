-- +goose Up
create table likes (
    user_id uuid not null,
    blog_id uuid not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(user_id, blog_id)
);

-- +goose Down
drop table likes;