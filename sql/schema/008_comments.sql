-- +goose Up
create table comments(
    id uuid primary key,
    description text not null,
    blog_id uuid not null references blogs(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table comments