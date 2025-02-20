-- +goose Up
create table Users(
    id UUID primary key,
    username text not null,
    email text not null unique,
    hashed_password text not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Users;