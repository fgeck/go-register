//go:build integrationtest

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	POSTGRES_IMAGE    = "postgres:latest"
	POSTGRES_USER     = "testuser"
	POSTGRES_PASSWORD = "testpassword"
	POSTGRES_DB       = "postgres"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	initScripts, err := getInitScripts()
	require.NoError(t, err)

	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	defer os.Unsetenv("TESTCONTAINERS_RYUK_DISABLED")
	postgresContainer, err := postgres.Run(ctx,
		POSTGRES_IMAGE,
		postgres.WithInitScripts(initScripts...),
		postgres.WithDatabase(POSTGRES_DB),
		postgres.WithUsername(POSTGRES_USER),
		postgres.WithPassword(POSTGRES_PASSWORD),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	// Set environment variables for the app
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", POSTGRES_USER)
	os.Setenv("DB_PASSWORD", POSTGRES_PASSWORD)
	os.Setenv("DB_DATABASE", POSTGRES_DB)
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")

	go func() {
		main()
	}()
	time.Sleep(1 * time.Second)

	// Perform a simple HTTP request to verify the app is running
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/login", "localhost", "8081"))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func getInitScripts() ([]string, error) {
	var initScripts []string
	err := filepath.Walk("../../migrations", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		initScripts = append(initScripts, path)
		return nil
	})

	return initScripts, err
}
