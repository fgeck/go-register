package server

import (
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/fgeck/go-register/internal/server/render"
)

//go:embed templates/*
var templatesFS embed.FS

type Server struct {
	router   *chi.Mux
	renderer *render.Renderer
}

func NewServer(authHandler, homeHandler interface{}) *Server {
	r := chi.NewRouter()
	renderer := render.New(templatesFS)

	if err := renderer.LoadTemplates(); err != nil {
		panic("failed to load templates: " + err.Error())
	}

	// Middleware
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Routes
	r.Get("/", homeHandler.(interface {
		HandleHome(http.ResponseWriter, *http.Request)
	}).HandleHome)
	r.Get("/login", authHandler.(interface {
		HandleLoginShow(http.ResponseWriter, *http.Request)
	}).HandleLoginShow)
	r.Post("/login", authHandler.(interface {
		HandleLogin(http.ResponseWriter, *http.Request)
	}).HandleLogin)
	r.Get("/register", authHandler.(interface {
		HandleRegisterShow(http.ResponseWriter, *http.Request)
	}).HandleRegisterShow)
	r.Post("/register", authHandler.(interface {
		HandleRegister(http.ResponseWriter, *http.Request)
	}).HandleRegister)
	r.Post("/logout", authHandler.(interface {
		HandleLogout(http.ResponseWriter, *http.Request)
	}).HandleLogout)

	return &Server{
		router:   r,
		renderer: renderer,
	}
}

// Use adds middleware to the router
func (s *Server) Use(middlewares ...func(http.Handler) http.Handler) {
	s.router.Use(middlewares...)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Add any cleanup logic here
	return nil
}

// Render renders a template
func (s *Server) Render(w http.ResponseWriter, status int, name string, data interface{}) {
	s.renderer.HTML(w, status, name, data)
}
