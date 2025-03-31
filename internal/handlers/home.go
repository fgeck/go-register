package handlers

import (
	"net/http"

	"github.com/fgeck/go-register/internal/auth"
	"github.com/fgeck/go-register/internal/server/render"
)

type HomeHandler struct {
	authService *auth.AuthService
}

func NewHomeHandler(authService *auth.AuthService) *HomeHandler {
	return &HomeHandler{authService: authService}
}

func (h *HomeHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	user := h.authService.UserFromContext(r.Context())

	data := map[string]interface{}{
		"User":    user,
		"Message": "Welcome to our Go + HTMX + Alpine.js app!",
	}

	render.HTML(w, http.StatusOK, "home.html", data)
}

func (h *HomeHandler) HandleProtected(w http.ResponseWriter, r *http.Request) {
	user := h.authService.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"User":    user,
		"Message": "This is a protected page!",
	}

	render.HTML(w, http.StatusOK, "protected.html", data)
}
