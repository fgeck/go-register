package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/validation"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (*UserCreatedDto, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDto, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	ValidateCreateUserParams(username, email, password string) error
}

type UserService struct {
	queries   repository.Querier
	validator validation.ValidationServiceInterface
}

func NewUserService(queries repository.Querier, validator validation.ValidationServiceInterface) *UserService {
	return &UserService{
		queries:   queries,
		validator: validator,
	}
}

type UserDto struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
}

func NewUserDto(user repository.User) *UserDto {
	return &UserDto{
		ID:           uuid.UUID(user.ID.Bytes),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}

type UserCreatedDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserCreatedDto(username, email string) *UserCreatedDto {
	return &UserCreatedDto{username, email}
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserDto, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return NewUserDto(user), err
}

func (s *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.queries.UserExistsByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, username, email, hashedPassword string) (*UserCreatedDto, error) {
	user, err := s.queries.CreateUser(
		ctx,
		repository.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
		},
	)
	if err != nil {
		// Todo log error
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
