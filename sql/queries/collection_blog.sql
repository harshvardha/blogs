-- name: AddBlogToCollection :one
insert into collection_blog (collection_id, blog_id, created_at, updated_at)
values (
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: RemoveBlogFromCollection :one
delete from collection_blog where collection_id = $1 and blog_id = $2
returning *;