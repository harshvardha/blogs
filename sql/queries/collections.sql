-- name: CreateCollection :one
insert into collections (id, name, user_id, created_at, updated_at)
values (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: EditCollection :one
update collections set name = $1, updated_at = NOW() where id = $2
returning *;

-- name: DeleteCollection :one
delete from collections where id = $1
returning *;

-- name: GetAllCollectionsByUserId :many
select * from collections where user_id = $1;

-- name: GetAllBlogsByCollectionId :many
select blogs.id, blogs.title, blogs.author_id, users.username as author_name, blogs.thumbnail_url, blogs.content, blogs.category, categories.category_name, blogs.created_at, blogs.updated_at from collections join collection_blog on collections.id = collection_blog.collection_id join blogs on collection_blog.blog_id = blogs.id join users on blogs.author_id = users.id join categories on blogs.category = categories.id where collections.id = $1;

-- name: GetOwnerId :one
select user_id from collections where id = $1;

-- name: GetCollectionNameById :one
select name from collections where id = $1;