-- +goose Up
create table blogs(
    id uuid primary key,
    title text not null,
    author_id uuid not null references users(id) on delete cascade,
    thumbnail_url text not null unique,
    content text not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table blogs;