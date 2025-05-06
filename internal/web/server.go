package web

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/config"
	"github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/security/jwt"
	"github.com/fgeck/go-register/internal/service/security/password"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/fgeck/go-register/internal/service/validation"
	"github.com/fgeck/go-register/internal/web/handlers"
	mw "github.com/fgeck/go-register/internal/web/middleware"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const (
	TWENTY_FOUR_HOURS_IN_SECONDS = 24 * 60 * 60
	ISSUER                       = "go-register"
	CONTEXT_TIMEOUT              = 10 * time.Second
)

func InitServer(e *echo.Echo, cfg *config.Config) {
	// Initialize DB
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()
	queries := connectToDatabase(ctx, cfg)
	createAdminUser(ctx, queries, cfg)

	// Services
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJwtService(cfg.App.JwtSecret, ISSUER, TWENTY_FOUR_HOURS_IN_SECONDS)
	loginRegisterService := loginRegister.NewLoginRegisterService(userService, passwordService, jwtService)

	// Handlers
	registerHandler := handlers.NewRegisterHandler(loginRegisterService)
	loginHandler := handlers.NewLoginHandler(loginRegisterService)

	// Middlewares
	authenticationMiddleware := mw.NewAuthenticationMiddleware(cfg.App.JwtSecret)
	authorizationMiddleware := mw.NewAuthorizationMiddleware()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	e.Static("/", "public")
	e.GET("/", handlers.HomeHandler)
	e.GET("/login", loginHandler.LoginRegisterContainerHandler)
	e.GET("/loginForm", loginHandler.LoginFormHandler)
	e.POST("/api/login", loginHandler.LoginHandler)
	e.GET("/registerForm", registerHandler.RegisterFormHandler)
	e.POST("/api/register", registerHandler.RegisterUserHandler)

	// JWT Middleware only
	res := e.Group("/restricted")
	res.Use(authenticationMiddleware.JwtAuthMiddleware())
	// for testing purposes
	res.GET("", func(c echo.Context) error {
		token, ok := c.Get("user").(*gojwt.Token)
		if !ok {
			return echo.ErrForbidden
		}
		claims, ok := token.Claims.(*jwt.JwtCustomClaims)
		if !ok {
			return echo.ErrForbidden
		}
		name := claims.UserId
		role := claims.UserRole

		return c.String(http.StatusOK, "Welcome "+name+" with role: "+role+"!")
	})

	// Admin Routes (requires "UserRole" == "admin")
	adminGroup := e.Group("/api/admin")
	adminGroup.Use(authenticationMiddleware.JwtAuthMiddleware(), authorizationMiddleware.RequireAdminMiddleware())
	// for testing purposes
	adminGroup.GET("/users", func(c echo.Context) error {
		token, ok := c.Get("user").(*gojwt.Token)
		if !ok {
			return echo.ErrForbidden
		}
		claims, ok := token.Claims.(*jwt.JwtCustomClaims)
		if !ok {
			return echo.ErrForbidden
		}
		name := claims.UserId
		role := claims.UserRole

		return c.String(http.StatusOK, "Welcome "+name+" with role: "+role+"!e")
	})
}

func connectToDatabase(ctx context.Context, cfg *config.Config) *repository.Queries {
	pgxConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			cfg.Db.User,
			cfg.Db.Password,
			net.JoinHostPort(cfg.Db.Host, cfg.Db.Port),
			cfg.Db.Database,
		),
	)
	if err != nil {
		panic(err)
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())

		return nil
	}

	pgxConnPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		panic(err)
	}

	// Verify database connection
	if err := pgxConnPool.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v\n", err)
		pgxConnPool.Close()
		os.Exit(1)
	}

	queries := repository.New(pgxConnPool)

	return queries
}

func createAdminUser(ctx context.Context, queries *repository.Queries, cfg *config.Config) {
	adminName := cfg.App.AdminUser
	adminPassword := cfg.App.AdminPassword
	adminEmail := cfg.App.AdminEmail
	hashedPassword, err := password.NewPasswordService().HashAndSaltPassword(adminPassword)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		return
	}

	exists, err := queries.UserExistsByEmail(ctx, cfg.App.AdminEmail)
	if err != nil {
		log.Printf("Error checking if admin user exists: %v\n", err)
		return
	}

	if exists {
		log.Println("Admin user already exists, skipping creation.")
		return
	}

	userParams := repository.CreateUserParams{
		Username:     adminName,
		Email:        adminEmail,
		PasswordHash: hashedPassword,
		UserRole:     "admin",
	}
	user, err := queries.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating admin user: %v\n", err)
		return
	}

	log.Printf("Admin user created successfully:\n"+
		"	id: %q\n	email: %q\n	username: %q\n",
		user.ID,
		user.Username,
		user.Email,
	)
}
