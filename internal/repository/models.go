// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           pgtype.UUID        `json:"id"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password_hash"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}
