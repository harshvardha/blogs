// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users_follow.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const followUser = `-- name: FollowUser :exec
insert into users_follow (follower_id, following_id, created_at, updated_at)
values ($1, $2, NOW(), NOW())
`

type FollowUserParams struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
}

func (q *Queries) FollowUser(ctx context.Context, arg FollowUserParams) error {
	_, err := q.db.ExecContext(ctx, followUser, arg.FollowerID, arg.FollowingID)
	return err
}

const getPair = `-- name: GetPair :one
select follower_id, following_id, created_at, updated_at from users_follow where follower_id = $1 and following_id = $2
`

type GetPairParams struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
}

func (q *Queries) GetPair(ctx context.Context, arg GetPairParams) (UsersFollow, error) {
	row := q.db.QueryRowContext(ctx, getPair, arg.FollowerID, arg.FollowingID)
	var i UsersFollow
	err := row.Scan(
		&i.FollowerID,
		&i.FollowingID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const unfollowUser = `-- name: UnfollowUser :exec
delete from users_follow where follower_id = $1 and following_id = $2
`

type UnfollowUserParams struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
}

func (q *Queries) UnfollowUser(ctx context.Context, arg UnfollowUserParams) error {
	_, err := q.db.ExecContext(ctx, unfollowUser, arg.FollowerID, arg.FollowingID)
	return err
}
