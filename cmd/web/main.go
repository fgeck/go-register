package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fgeck/go-register/internal/auth"
	"github.com/fgeck/go-register/internal/handlers"
	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/server"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize context with timeout for startup operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Database setup
	pgxConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
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
		log.Fatalf("Database ping failed: %v\n", err)
	}

	// Run migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Migrations failed: %v\n", err)
	}

	// Initialize repository
	repo := repository.New(pgxConnPool)

	// Auth service setup with config
	authService := auth.NewAuthService(repo, auth.AuthConfig{
		SessionExpiry: 7 * 24 * time.Hour, // 1 week
		Pepper:        cfg.SessionSecret,
		Argon2Params: &auth.Argon2Params{
			Memory:      cfg.Argon2Config.Memory,
			Iterations:  cfg.Argon2Config.Iterations,
			Parallelism: cfg.Argon2Config.Parallelism,
			SaltLength:  cfg.Argon2Config.SaltLength,
			KeyLength:   cfg.Argon2Config.KeyLength,
		},
	})

	// Handlers setup
	authHandler := handlers.NewAuthHandler(authService)
	homeHandler := handlers.NewHomeHandler(authService)

	// HTTP Server setup
	srv := server.NewServer(authHandler, homeHandler)

	// Add middleware stack
	srv.Use(
		server.RequestLogger,
		authService.AuthMiddleware,
		server.CSRFProtection(cfg.CSRFKey),
	)

	// Graceful shutdown setup
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Server starting on %s (HTTPS: %v)\n", cfg.ServerAddress, cfg.HTTPS)

		var err error
		if cfg.HTTPS {
			err = srv.Start(cfg.ServerAddress)
		} else {
			err = srv.Start(cfg.ServerAddress)
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	<-done
	log.Println("Server shutting down...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func runMigrations(databaseURL string) error {
	// In production, use a proper migration tool like golang-migrate
	// This is just a placeholder
	log.Println("Running database migrations...")
	return nil
}
