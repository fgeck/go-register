// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	DropAllUsers(ctx context.Context) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id pgtype.UUID) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}

var _ Querier = (*Queries)(nil)
