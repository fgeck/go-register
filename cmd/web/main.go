package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fgeck/go-register/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	// cfg := loadConfig()
	port := "8080"

	// Initialize context with timeout for startup operations
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// // Database setup
	// pgxConfig, err := pgxpool.ParseConfig("postgres://user:password@localhost:5432/postgres?sslmode=disable")
	// if err != nil {
	// 	panic(err)
	// }

	// pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	// 	pgxUUID.Register(conn.TypeMap())
	// 	return nil
	// }

	// pgxConnPool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	// if err != nil {
	// 	panic(err)
	// }
	// defer pgxConnPool.Close()

	// // Verify database connection
	// if err := pgxConnPool.Ping(ctx); err != nil {
	// 	log.Fatalf("Database ping failed: %v\n", err)
	// }

	// Run migrations
	// if err := runMigrations(cfg.DatabaseURL); err != nil {
	// 	log.Fatalf("Migrations failed: %v\n", err)
	// }

	// Initialize repository
	// repo := repository.New(pgxConnPool)

	// Initialize server
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	// create and use render

	handlers.RegisterAllHandlers(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func runMigrations(databaseURL string) error {
	// In production, use a proper migration tool like golang-migrate
	// This is just a placeholder
	log.Println("Running database migrations...")
	return nil
}
