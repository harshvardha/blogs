-- name: CreateCategory :one
insert into categories (id, category_name, created_at, updated_at)
values (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning *;

-- name: EditCategory :one
update categories set category_name = $1, updated_at = NOW() where id = $2
returning *;

-- name: DeleteCategory :one
delete from categories where id = $1
returning *;

-- name: GetCategoryIdByName :one
select id from categories where category_name = $1;

-- name: GetCategoryNameById :one
select category_name from categories where id = $1;