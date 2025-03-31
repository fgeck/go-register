package auth

import (
	"context"
	"net/http"

	"github.com/fgeck/go-register/internal/repository"
)

type contextKey string

const (
	sessionCookieName = "session_token"
	userContextKey    = contextKey("user")
)

func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Validate session
		session, err := s.repo.GetSession(r.Context(), cookie.Value)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			next.ServeHTTP(w, r)
			return
		}

		// Get user
		user, err := s.repo.GetUser(r.Context(), session.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AuthService) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := UserFromContext(r.Context()); user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserFromContext(ctx context.Context) *repository.User {
	if user, ok := ctx.Value(userContextKey).(*repository.User); ok {
		return user
	}
	return nil
}
