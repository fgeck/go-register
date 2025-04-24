package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fgeck/go-register/internal/handlers"
	"github.com/fgeck/go-register/internal/repository"
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
	// Load configuration
	// cfg := loadConfig()
	port := "8080"

	// Initialize context with timeout for startup operations
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()

	// Database setup
	pgxConfig, err := pgxpool.ParseConfig("postgres://user:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())

		return nil
	}

	pgxConnPool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		panic(err)
	}
	defer pgxConnPool.Close()

	// Verify database connection
	if err := pgxConnPool.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v\n", err)
		pgxConnPool.Close()
		//nolint
		os.Exit(1)
	}

	// Run migrations
	// if err := runMigrations(cfg.DatabaseURL); err != nil {
	// 	log.Fatalf("Migrations failed: %v\n", err)
	// }

	// Initialize repository
	queries := repository.New(pgxConnPool)

	// Initialize server
	echoServer := echo.New()
	// Middleware
	echoServer.Use(middleware.Logger())
	// create and use render

	handlers.SetupHandlers(echoServer, queries)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := echoServer.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
