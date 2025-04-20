package loginRegister_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	userfacing_errors "github.com/fgeck/go-register/internal/service/errors"
	"github.com/fgeck/go-register/internal/service/loginRegister"
	jwt "github.com/fgeck/go-register/internal/service/security/jwt/mocks"
	password "github.com/fgeck/go-register/internal/service/security/password/mocks"
	"github.com/fgeck/go-register/internal/service/user"
	userMocks "github.com/fgeck/go-register/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupLoginRegisterServiceTest(t *testing.T) (*userMocks.MockUserServiceInterface, *password.MockPasswordServiceInterface, *jwt.MockJwtServiceInterface, *loginRegister.LoginRegisterService) {
	mockUserService := userMocks.NewMockUserServiceInterface(t)
	mockPasswordService := password.NewMockPasswordServiceInterface(t)
	mockJwtService := jwt.NewMockJwtServiceInterface(t)
	service := loginRegister.NewLoginRegisterService(mockUserService, mockPasswordService, mockJwtService)
	return mockUserService, mockPasswordService, mockJwtService, service
}

func TestLoginUser(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	email := "testuser@example.com"
	username := "testuser"
	password := "Valid1@"
	hashedPassword := "hashedpassword"
	token := "mockJwtToken"

	t.Run("successfully logs in user", func(t *testing.T) {
		mockUserService, mockPasswordService, mockJwtService, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("GetUserByEmail", ctx, email).Return(&user.UserDto{
			ID:           id,
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
		}, nil)
		mockPasswordService.On("ComparePassword", hashedPassword, password).Return(nil)
		mockJwtService.On("GenerateToken", mock.Anything).Return(token, nil)

		result, err := service.LoginUser(ctx, email, password)

		assert.NoError(t, err)
		assert.Equal(t, token, result)

		mockUserService.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJwtService.AssertExpectations(t)
	})

	t.Run("fails when user does not exist", func(t *testing.T) {
		mockUserService, _, _, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("GetUserByEmail", ctx, email).Return(nil, errors.New("user not found"))

		result, err := service.LoginUser(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "user not found", err.Error())

		mockUserService.AssertExpectations(t)
	})

	t.Run("fails when password is invalid", func(t *testing.T) {
		mockUserService, mockPasswordService, _, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("GetUserByEmail", ctx, email).Return(&user.UserDto{
			Email:        email,
			PasswordHash: hashedPassword,
		}, nil)
		mockPasswordService.On("ComparePassword", hashedPassword, password).Return(errors.New("invalid password"))

		result, err := service.LoginUser(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "invalid password", err.Error())

		mockUserService.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})
}

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	username := "testuser"
	email := "testuser@example.com"
	password := "Valid1@"
	hashedPassword := "hashedpassword"

	t.Run("successfully registers user", func(t *testing.T) {
		mockUserService, mockPasswordService, _, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("UserExistsByEmail", ctx, email).Return(false, nil)
		mockUserService.On("ValidateCreateUserParams", username, email, password).Return(nil)
		mockPasswordService.On("HashAndSaltPassword", password).Return(hashedPassword, nil)
		mockUserService.On("CreateUser", ctx, username, email, hashedPassword).Return(&user.UserCreatedDto{
			Username: username,
			Email:    email,
		}, nil)

		result, err := service.RegisterUser(ctx, username, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, username, result.Username)
		assert.Equal(t, email, result.Email)

		mockUserService.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("fails when user already exists", func(t *testing.T) {
		mockUserService, _, _, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("UserExistsByEmail", ctx, email).Return(true, nil)

		result, err := service.RegisterUser(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, result)

		// Check for UserFacingError
		ufe, ok := err.(*userfacing_errors.UserFacingError)
		assert.True(t, ok, "expected a UserFacingError")
		assert.Equal(t, "user already exists", ufe.Message)
		assert.Equal(t, http.StatusConflict, ufe.Code)

		mockUserService.AssertExpectations(t)
	})

	t.Run("fails when validation fails", func(t *testing.T) {
		mockUserService, _, _, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("UserExistsByEmail", ctx, email).Return(false, nil)
		mockUserService.On("ValidateCreateUserParams", username, email, password).Return(userfacing_errors.New("invalid input", http.StatusBadRequest))

		result, err := service.RegisterUser(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, result)

		// Check for UserFacingError
		ufe, ok := err.(*userfacing_errors.UserFacingError)
		assert.True(t, ok, "expected a UserFacingError")
		assert.Equal(t, "invalid input", ufe.Message)
		assert.Equal(t, http.StatusBadRequest, ufe.Code)

		mockUserService.AssertExpectations(t)
	})

	t.Run("fails when hashing password fails", func(t *testing.T) {
		mockUserService, mockPasswordService, _, service := setupLoginRegisterServiceTest(t)

		mockUserService.On("UserExistsByEmail", ctx, email).Return(false, nil)
		mockUserService.On("ValidateCreateUserParams", username, email, password).Return(nil)
		mockPasswordService.On("HashAndSaltPassword", password).Return("", errors.New("hashing error"))

		result, err := service.RegisterUser(ctx, username, email, password)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "hashing error", err.Error())

		mockUserService.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})
}
