// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package user

import (
	"context"

	"github.com/fgeck/go-register/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// NewMockUserServiceInterface creates a new instance of MockUserServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserServiceInterface {
	mock := &MockUserServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockUserServiceInterface is an autogenerated mock type for the UserServiceInterface type
type MockUserServiceInterface struct {
	mock.Mock
}

type MockUserServiceInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserServiceInterface) EXPECT() *MockUserServiceInterface_Expecter {
	return &MockUserServiceInterface_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function for the type MockUserServiceInterface
func (_mock *MockUserServiceInterface) CreateUser(ctx context.Context, username string, email string, passwordHash string) (repository.User, error) {
	ret := _mock.Called(ctx, username, email, passwordHash)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 repository.User
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string, string) (repository.User, error)); ok {
		return returnFunc(ctx, username, email, passwordHash)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string, string) repository.User); ok {
		r0 = returnFunc(ctx, username, email, passwordHash)
	} else {
		r0 = ret.Get(0).(repository.User)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = returnFunc(ctx, username, email, passwordHash)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockUserServiceInterface_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockUserServiceInterface_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx
//   - username
//   - email
//   - passwordHash
func (_e *MockUserServiceInterface_Expecter) CreateUser(ctx interface{}, username interface{}, email interface{}, passwordHash interface{}) *MockUserServiceInterface_CreateUser_Call {
	return &MockUserServiceInterface_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, username, email, passwordHash)}
}

func (_c *MockUserServiceInterface_CreateUser_Call) Run(run func(ctx context.Context, username string, email string, passwordHash string)) *MockUserServiceInterface_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockUserServiceInterface_CreateUser_Call) Return(user repository.User, err error) *MockUserServiceInterface_CreateUser_Call {
	_c.Call.Return(user, err)
	return _c
}

func (_c *MockUserServiceInterface_CreateUser_Call) RunAndReturn(run func(ctx context.Context, username string, email string, passwordHash string) (repository.User, error)) *MockUserServiceInterface_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateCreateUserParams provides a mock function for the type MockUserServiceInterface
func (_mock *MockUserServiceInterface) ValidateCreateUserParams(username string, email string, password string) error {
	ret := _mock.Called(username, email, password)

	if len(ret) == 0 {
		panic("no return value specified for ValidateCreateUserParams")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = returnFunc(username, email, password)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockUserServiceInterface_ValidateCreateUserParams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateCreateUserParams'
type MockUserServiceInterface_ValidateCreateUserParams_Call struct {
	*mock.Call
}

// ValidateCreateUserParams is a helper method to define mock.On call
//   - username
//   - email
//   - password
func (_e *MockUserServiceInterface_Expecter) ValidateCreateUserParams(username interface{}, email interface{}, password interface{}) *MockUserServiceInterface_ValidateCreateUserParams_Call {
	return &MockUserServiceInterface_ValidateCreateUserParams_Call{Call: _e.mock.On("ValidateCreateUserParams", username, email, password)}
}

func (_c *MockUserServiceInterface_ValidateCreateUserParams_Call) Run(run func(username string, email string, password string)) *MockUserServiceInterface_ValidateCreateUserParams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockUserServiceInterface_ValidateCreateUserParams_Call) Return(err error) *MockUserServiceInterface_ValidateCreateUserParams_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockUserServiceInterface_ValidateCreateUserParams_Call) RunAndReturn(run func(username string, email string, password string) error) *MockUserServiceInterface_ValidateCreateUserParams_Call {
	_c.Call.Return(run)
	return _c
}
