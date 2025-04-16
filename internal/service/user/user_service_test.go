package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fgeck/go-register/internal/repository"
	repositoryMocks "github.com/fgeck/go-register/internal/repository/mocks"
	"github.com/fgeck/go-register/internal/service/user"
	validationMocks "github.com/fgeck/go-register/internal/service/validation/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserServiceTest(t *testing.T) (*repositoryMocks.MockQuerier, *validationMocks.MockValidationServiceInterface, *user.UserService) {
	mockQueries := repositoryMocks.NewMockQuerier(t)
	mockValidator := validationMocks.NewMockValidationServiceInterface(t)
	userService := user.NewUserService(mockQueries, mockValidator)
	return mockQueries, mockValidator, userService
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	username := "testuser"
	email := "testuser@example.com"
	passwordHash := "hashedpassword"

	t.Run("successfully creates a user", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)

		mockQueries.On("UserExistsByEmail", ctx, email).Return(false, nil)
		mockQueries.On("CreateUser", ctx, mock.Anything).Return(repository.User{
			Username: username,
			Email:    email,
		}, nil)

		userDto, err := userService.CreateUser(ctx, username, email, passwordHash)

		assert.NoError(t, err)
		assert.NotNil(t, userDto)
		assert.Equal(t, username, userDto.Username)
		assert.Equal(t, email, userDto.Email)

		mockQueries.AssertExpectations(t)
	})

	t.Run("fails when user already exists", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("UserExistsByEmail", ctx, email).Return(true, nil)

		userDto, err := userService.CreateUser(ctx, username, email, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, userDto)
		assert.Equal(t, "user already exists", err.Error())

		mockQueries.AssertExpectations(t)
	})

	t.Run("fails when database error occurs", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("UserExistsByEmail", ctx, email).Return(false, nil)
		mockQueries.On("CreateUser", ctx, mock.Anything).Return(repository.User{}, errors.New("database error"))

		userDto, err := userService.CreateUser(ctx, username, email, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, userDto)
		assert.Equal(t, "database error", err.Error())

		mockQueries.AssertExpectations(t)
	})
}

func TestValidateCreateUserParams(t *testing.T) {
	username := "testuser"
	email := "testuser@example.com"
	password := "Valid1@"

	t.Run("successfully validates parameters", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(nil)
		mockValidator.On("ValidatePassword", password).Return(nil)
		mockValidator.On("ValidateUsername", username).Return(nil)

		err := userService.ValidateCreateUserParams(username, email, password)

		assert.NoError(t, err)

		mockValidator.AssertExpectations(t)
	})

	t.Run("fails when email validation fails", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(errors.New("invalid email format"))

		err := userService.ValidateCreateUserParams(username, email, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid email format", err.Error())

		mockValidator.AssertExpectations(t)
	})

	t.Run("fails when password validation fails", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(nil)
		mockValidator.On("ValidatePassword", password).Return(errors.New("password too weak"))

		err := userService.ValidateCreateUserParams(username, email, password)

		assert.Error(t, err)
		assert.Equal(t, "password too weak", err.Error())

		mockValidator.AssertExpectations(t)
	})

	t.Run("fails when username validation fails", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(nil)
		mockValidator.On("ValidatePassword", password).Return(nil)
		mockValidator.On("ValidateUsername", username).Return(errors.New("username too short"))

		err := userService.ValidateCreateUserParams(username, email, password)

		assert.Error(t, err)
		assert.Equal(t, "username too short", err.Error())

		mockValidator.AssertExpectations(t)
	})
}
