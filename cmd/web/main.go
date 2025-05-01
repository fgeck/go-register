package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/config"
	"github.com/fgeck/go-register/internal/service/security/password"
	"github.com/fgeck/go-register/internal/web/handlers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const (
	CONTEXT_TIMEOUT = 10 * time.Second
)

func main() {
	cfgLoader := config.NewLoader()
	cfg, err := cfgLoader.LoadConfig("")
	if err != nil {
		panic(err)
	}

	// Initialize context with timeout for startup operations
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()

	queries := connectToDatabase(ctx, cfg)
	createAdminUser(ctx, queries, cfg)

	echoServer := echo.New()
	echoServer.Use(middleware.Logger())
	handlers.InitServer(echoServer, queries, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		err := echoServer.Start(net.JoinHostPort(cfg.App.Host, cfg.App.Port))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			echoServer.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)

	defer cancel()

	if err := echoServer.Shutdown(ctx); err != nil {
		echoServer.Logger.Fatal(err)
	}
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
