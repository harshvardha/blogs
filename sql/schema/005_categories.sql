-- +goose Up
create table categories(
    id uuid primary key,
    category_name text not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table categories;