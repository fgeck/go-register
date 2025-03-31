package handlers

import (
	"net/http"

	"github.com/fgeck/go-register/internal/auth"
	"github.com/fgeck/go-register/internal/server/render"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) HandleLoginShow(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, http.StatusOK, "auth/login", nil)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	sessionID, err := h.authService.Login(r.Context(), email, password)
	if err != nil {
		render.HTML(w, http.StatusUnauthorized, "auth/login", map[string]interface{}{
			"Error": "Invalid credentials",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // In production set this to true
		SameSite: http.SameSiteLaxMode,
	})

	if isHTMX := r.Header.Get("HX-Request"); isHTMX == "true" {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) HandleRegisterShow(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, http.StatusOK, "auth/register", nil)
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.authService.Register(r.Context(), username, email, password)
	if err != nil {
		render.HTML(w, http.StatusBadRequest, "auth/register", map[string]interface{}{
			"Error": "Registration failed",
		})
		return
	}

	if isHTMX := r.Header.Get("HX-Request"); isHTMX == "true" {
		w.Header().Set("HX-Redirect", "/login")
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err == nil {
		// Delete session from DB
		h.authService.repo.DeleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	if isHTMX := r.Header.Get("HX-Request"); isHTMX == "true" {
		w.Header().Set("HX-Redirect", "/login")
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
