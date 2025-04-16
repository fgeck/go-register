package service

import (
	"context"

	"github.com/fgeck/go-register/internal/repository"
	service "github.com/fgeck/go-register/internal/service/validation"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (repository.User, error)
	ValidateCreateUserParams(username, email, password string) error
}

type UserService struct {
	repo      *repository.Queries
	validator *service.ValidationService
}

func NewUserService(repo *repository.Queries) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username, email, passwordHash string) (repository.User, error) {
	user, err := s.repo.CreateUser(
		ctx,
		repository.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		return repository.User{}, err
	}
	return user, nil
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
