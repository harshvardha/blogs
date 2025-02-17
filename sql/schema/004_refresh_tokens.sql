-- +goose Up
create table refresh_token(
    token text primary key,
    user_id uuid not null references users(id) on delete cascade,
    expires_at timestamp not null,
    revoked_at timestamp,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table refresh_token;