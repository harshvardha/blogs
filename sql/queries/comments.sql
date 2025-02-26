-- name: CreateComment :one
insert into comments (id, description, blog_id, user_id, created_at, updated_at)
values (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
returning *;

-- name: EditComment :one
update comments set description = $1, updated_at = NOW() where id = $2 and blog_id = $3
returning *;

-- name: DeleteComment :one
delete from comments where id = $1
returning *;

-- name: LikeComment :exec
insert into comment_likes (user_id, comment_id, created_at, updated_at)
values ($1, $2, NOW(), NOW());

-- name: UnlikeComment :exec
delete from comment_likes where user_id = $1 and comment_id = $2;

-- name: GetAllCommentsByBlogId :many
select comments.id, 
    comments.description, 
    comments.blog_id, 
    comments.user_id, 
    comments.created_at, 
    comments.updated_at, 
    count(comment_likes.comment_id) as likes_count 
    from comments left join comment_likes on comments.id = comment_likes.comment_id where comments.blog_id = $1
    group by comments.id;

-- name: GetCommentById :one
select * from comments where id = $1;

-- name: IsCommentLiked :one
select user_id, comment_id from comment_likes where user_id = $1 and comment_id = $2;