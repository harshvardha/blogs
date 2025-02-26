// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: collection_blog.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addBlogToCollection = `-- name: AddBlogToCollection :one
insert into collection_blog (collection_id, blog_id, created_at, updated_at)
values (
    $1,
    $2,
    NOW(),
    NOW()
)
returning collection_id, blog_id, created_at, updated_at
`

type AddBlogToCollectionParams struct {
	CollectionID uuid.UUID
	BlogID       uuid.UUID
}

func (q *Queries) AddBlogToCollection(ctx context.Context, arg AddBlogToCollectionParams) (CollectionBlog, error) {
	row := q.db.QueryRowContext(ctx, addBlogToCollection, arg.CollectionID, arg.BlogID)
	var i CollectionBlog
	err := row.Scan(
		&i.CollectionID,
		&i.BlogID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeBlogFromCollection = `-- name: RemoveBlogFromCollection :one
delete from collection_blog where collection_id = $1 and blog_id = $2
returning collection_id, blog_id, created_at, updated_at
`

type RemoveBlogFromCollectionParams struct {
	CollectionID uuid.UUID
	BlogID       uuid.UUID
}

func (q *Queries) RemoveBlogFromCollection(ctx context.Context, arg RemoveBlogFromCollectionParams) (CollectionBlog, error) {
	row := q.db.QueryRowContext(ctx, removeBlogFromCollection, arg.CollectionID, arg.BlogID)
	var i CollectionBlog
	err := row.Scan(
		&i.CollectionID,
		&i.BlogID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
