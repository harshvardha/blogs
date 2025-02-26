-- +goose Up
alter table blogs add column category uuid not null references categories(id) on delete cascade;

-- +goose Down
alter table blogs drop column category;