package user_test

import (
	"testing"

	"github.com/fgeck/go-register/internal/service/user"
	"github.com/stretchr/testify/assert"
)

func TestUSerRoleFromString(t *testing.T) {
	t.Run("returns UserRoleUser for USER", func(t *testing.T) {
		role := user.UserRoleFromString("USER")
		assert.Equal(t, user.UserRoleUser, role)
	})
	t.Run("returns UserRoleAdmin for ADMIN", func(t *testing.T) {
		role := user.UserRoleFromString("ADMIN")
		assert.Equal(t, user.UserRoleAdmin, role)
	})
	t.Run("returns UserRoleUser for unknown role", func(t *testing.T) {
		role := user.UserRoleFromString("UNKNOWN")
		assert.Equal(t, user.UserRoleUser, role)
	})
	t.Run("returns UserRoleUser for empty string", func(t *testing.T) {
		role := user.UserRoleFromString("")
		assert.Equal(t, user.UserRoleUser, role)
	})
}

func TestIsAdmin(t *testing.T) {
	t.Run("returns true for UserRoleAdmin", func(t *testing.T) {
		userDto := &user.UserDto{Role: user.UserRoleAdmin}
		assert.True(t, userDto.IsAdmin())
	})
	t.Run("returns false for UserRoleUser", func(t *testing.T) {
		userDto := &user.UserDto{Role: user.UserRoleUser}
		assert.False(t, userDto.IsAdmin())
	})
}
func TestIsUser(t *testing.T) {
	t.Run("returns true for UserRoleUser", func(t *testing.T) {
		userDto := &user.UserDto{Role: user.UserRoleUser}
		assert.True(t, userDto.IsUser())
	})
	t.Run("returns false for UserRoleAdmin", func(t *testing.T) {
		userDto := &user.UserDto{Role: user.UserRoleAdmin}
		assert.False(t, userDto.IsUser())
	})
}
