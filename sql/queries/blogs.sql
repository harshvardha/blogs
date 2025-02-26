-- name: CreateBlog :one
insert into blogs(
    id, 
    title, 
    author_id, 
    thumbnail_url, 
    content, 
    category, 
    created_at, 
    updated_at
)
values (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW()
)
returning *;

-- name: EditBlog :one
update blogs set title = $1, thumbnail_url = $2, content = $3, category = $4, updated_at = NOW() where id = $5
returning *;

-- name: DeleteBlog :one
delete from blogs where id = $1
returning *;

-- name: GetBlogById :one
select blogs.id, 
    blogs.title, 
    blogs.author_id,
    blogs.thumbnail_url,
    blogs.content,
    blogs.category,
    blogs.created_at,
    blogs.updated_at,
    count(likes.blog_id) as likes_count 
    from blogs left join likes on blogs.id = likes.blog_id 
    where blogs.id = $1 group by blogs.id, blogs.title;

-- name: GetBlogsByAuthorId :many
select blogs.id, blogs.title, blogs.author_id, blogs.content, blogs.thumbnail_url, blogs.category, blogs.created_at, blogs.updated_at, count(likes.blog_id) as likes_count from blogs left join likes on blogs.id = likes.blog_id where blogs.author_id = $1 group by blogs.id, blogs.title, blogs.author_id, blogs.thumbnail_url;

-- name: LikeBlog :exec
insert into likes (user_id, blog_id, created_at, updated_at)
values ($1, $2, NOW(), NOW());

-- name: UnlikeBlog :exec
delete from likes where user_id = $1 and blog_id = $2;

-- name: GetNoOfLikes :one
select count(*) from likes where blog_id = $1;

-- name: IsBlogLiked :one
select * from likes where user_id = $1 and blog_id = $2;

-- name: GetBlogsByTitle :many
select id, title, author_id, thumbnail_url from blogs where title = $1;

-- name: GetBlogsByCategory :many
select blogs.id, blogs.title, blogs.author_id, blogs.thumbnail_url, count(likes.blog_id) as likes_count from blogs left join likes on blogs.id = likes.blog_id where blogs.category = $1 group by blogs.id, blogs.title, blogs.author_id, blogs.thumbnail_url;

-- name: GetAuthorNameByBlogId :one
select username from users join blogs on users.id = blogs.author_id where blogs.id = $1;

-- name: GetBlogNameById :one
select title from blogs where id = $1;