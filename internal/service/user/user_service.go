package user

import (
	"context"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/validation"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (repository.User, error)
	ValidateCreateUserParams(username, email, password string) error
}

type UserService struct {
	queries   *repository.Queries
	validator validation.ValidationServiceInterface
}

func NewUserService(queries *repository.Queries, validator validation.ValidationServiceInterface) *UserService {
	return &UserService{
		queries:   queries,
		validator: validator,
	}
}

type UserCreatedDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserCreatedDto(username, email string) *UserCreatedDto {
	return &UserCreatedDto{username, email}
}

func (s *UserService) CreateUser(ctx context.Context, username, email, passwordHash string) (*UserCreatedDto, error) {
	user, err := s.queries.CreateUser(
		ctx,
		repository.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		return nil, err
	}

	return NewUserCreatedDto(user.Username, user.Email), nil
}

func (s *UserService) ValidateCreateUserParams(username, email, password string) error {
	if err := s.validator.ValidateEmail(email); err != nil {
		return err
	}
	if err := s.validator.ValidatePassword(password); err != nil {
		return err
	}
	if err := s.validator.ValidateUsername(username); err != nil {
		return err
	}
	return nil
}
