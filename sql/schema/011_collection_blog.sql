-- +goose Up
create table collection_blog (
    collection_id uuid not null references collections(id) on delete cascade,
    blog_id uuid not null references blogs(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(collection_id, blog_id)
);

-- +goose Down
drop table collection_blog;