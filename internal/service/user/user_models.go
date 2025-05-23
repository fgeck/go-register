package user

import (
	"strings"

	"github.com/fgeck/gotth-postgres/internal/repository"
	"github.com/google/uuid"
)

type UserCreatedDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserDto struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Role         UserRole  `json:"role"`
}

func NewUserDto(user repository.User) *UserDto {
	return &UserDto{
		ID:           uuid.UUID(user.ID.Bytes),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         UserRoleFromString(user.UserRole),
	}
}

var (
	UserRoleUser  = UserRole{Name: "USER"}
	UserRoleAdmin = UserRole{Name: "ADMIN"}
)

type UserRole struct {
	Name string `json:"name"`
}

func UserRoleFromString(name string) UserRole {
	switch strings.ToUpper(name) {
	case UserRoleUser.Name:
		return UserRoleUser
	case UserRoleAdmin.Name:
		return UserRoleAdmin
	default:
		return UserRoleUser
	}
}

func (u *UserDto) IsAdmin() bool {
	return u.Role.Name == UserRoleAdmin.Name
}

func (u *UserDto) IsUser() bool {
	return u.Role.Name == UserRoleUser.Name
}
