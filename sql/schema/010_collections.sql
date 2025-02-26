-- +goose Up
create table collections (
    id uuid not null primary key,
    name text not null,
    user_id uuid not null references users(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table collections;